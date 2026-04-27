// Package update implements the cubox-cli self-update notification check.
//
// The check is split into two operations:
//
//   - Lookup is a synchronous, network-free read of an on-disk cache.
//     It returns whatever the last successful registry fetch wrote, plus a
//     hint about whether the cache should be refreshed.
//   - RefreshCache performs a single HTTP GET against the npm registry and
//     writes the response (or a negative-cache marker on failure) back to
//     disk. It is meant to be invoked from a goroutine while the main
//     command runs, so its result is observed on the *next* invocation in
//     the common case.
//
// All file and network errors are deliberately swallowed: an update notice
// is best-effort and must never block or break a normal CLI invocation.
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultCacheTTL    = 6 * time.Hour
	defaultNegativeTTL = 1 * time.Hour
	defaultHTTPTimeout = 3 * time.Second
	defaultRegistryURL = "https://registry.npmjs.org/cubox-cli/latest"
	defaultUpdateCmd   = "npm update -g cubox-cli && npx skills add OLCUBO/cubox-cli -g -y"
)

// Info is the payload included as `_notice.update` in JSON output and printed
// to stderr in text mode. The Command field is a human/agent display hint and
// MUST NOT be executed automatically.
type Info struct {
	Current string `json:"current"`
	Latest  string `json:"latest"`
	Message string `json:"message"`
	Command string `json:"command"`
}

// State describes the update information available from a synchronous Lookup.
type State struct {
	Info         *Info // non-nil when a strictly newer version is known
	NeedsRefresh bool  // true when the cache is stale or empty (and not throttled)
	CacheEmpty   bool  // true when no usable cached version was found
}

// Notifier coordinates lookups and background refreshes. A single Notifier per
// process is expected. Methods are not safe for concurrent use within one
// instance unless documented otherwise; the typical pattern is:
//
//	n := update.Default(Version)
//	state := n.Lookup()
//	go n.RefreshCache(ctx) // optional, when state.NeedsRefresh
//
// Lookup and RefreshCache may safely run concurrently across goroutines: the
// cache file is written via a temp-file + rename, so a concurrent reader
// observes either the previous or the next complete file, never a partial
// one.
type Notifier struct {
	// Version is the current binary version. Empty, "dev", or a git-describe
	// suffix disables the check (see shouldCheck).
	Version string

	// RegistryURL is the npm registry endpoint (`/<package>/latest`).
	RegistryURL string

	// CachePath overrides the default cache file location. When empty
	// DefaultCachePath() is used.
	CachePath string

	// HTTPClient allows tests to inject transports. When nil, a client with
	// HTTPTimeout is constructed lazily.
	HTTPClient *http.Client

	// UserAgent is sent on registry requests.
	UserAgent string

	// CacheTTL is how long a successfully-fetched value remains fresh.
	CacheTTL time.Duration

	// NegativeTTL is how long after a failed fetch we suppress retries to
	// avoid hammering an offline / blocked registry on every invocation.
	NegativeTTL time.Duration

	// HTTPTimeout caps the registry request when HTTPClient is nil.
	HTTPTimeout time.Duration

	// UpdateCommand is the install/refresh hint embedded into Info.
	UpdateCommand string

	// Now is the clock source. nil means time.Now.
	Now func() time.Time

	// Getenv overrides os.Getenv for tests. nil means os.Getenv.
	Getenv func(key string) string

	// Debug, when non-nil, receives diagnostic messages. Use Default's
	// CUBOX_DEBUG-aware implementation in production.
	Debug func(format string, args ...any)
}

// Default returns a Notifier wired up for production use.
func Default(version string) *Notifier {
	n := &Notifier{
		Version:       version,
		RegistryURL:   defaultRegistryURL,
		CachePath:     DefaultCachePath(),
		CacheTTL:      defaultCacheTTL,
		NegativeTTL:   defaultNegativeTTL,
		HTTPTimeout:   defaultHTTPTimeout,
		UpdateCommand: defaultUpdateCmd,
		UserAgent:     fmt.Sprintf("cubox-cli/%s (+update-check)", versionForUA(version)),
		Now:           time.Now,
	}
	if os.Getenv("CUBOX_DEBUG") != "" {
		n.Debug = func(format string, args ...any) {
			fmt.Fprintf(os.Stderr, "[update] "+format+"\n", args...)
		}
	}
	return n
}

// DefaultCachePath returns the default on-disk cache location.
// Returns "" when the user home directory cannot be determined.
func DefaultCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "cubox-cli", "update-check.json")
}

func versionForUA(v string) string {
	if v == "" {
		return "dev"
	}
	return v
}

// cacheFile is the on-disk JSON layout. New fields must be additive — older
// binaries should be able to read newer files without crashing.
type cacheFile struct {
	LatestVersion string    `json:"latest_version,omitempty"`
	CheckedAt     time.Time `json:"checked_at,omitempty"`
	LastFailure   time.Time `json:"last_failure,omitempty"`
}

// Lookup returns the current update state read synchronously from cache.
// It performs no network I/O. A "" or "dev" Version, a git-describe build,
// CI environments, or an explicit opt-out short-circuit to a zero State so
// dev / non-interactive contexts never trigger the machinery.
func (n *Notifier) Lookup() State {
	if !n.shouldCheck() {
		return State{}
	}

	c, ok := n.readCache()
	if !ok {
		return State{NeedsRefresh: true, CacheEmpty: true}
	}

	now := n.now()
	fresh := !c.CheckedAt.IsZero() && now.Sub(c.CheckedAt) <= n.cacheTTL()
	suppressed := !c.LastFailure.IsZero() && now.Sub(c.LastFailure) < n.negativeTTL()

	state := State{
		NeedsRefresh: !fresh && !suppressed,
		CacheEmpty:   c.LatestVersion == "",
	}
	if c.LatestVersion != "" && isNewer(c.LatestVersion, n.Version) {
		state.Info = n.newInfo(c.LatestVersion)
	}
	return state
}

// RefreshCache fetches the latest version and writes the result to disk.
// On any error a negative-cache entry is written so subsequent Lookups can
// honor NegativeTTL and avoid retry storms. Honors the supplied context's
// deadline / cancellation.
func (n *Notifier) RefreshCache(ctx context.Context) {
	if !n.shouldCheck() {
		return
	}
	v, err := n.fetchRegistry(ctx)
	if err != nil {
		n.debug("registry fetch failed: %v", err)
		n.writeNegativeCache()
		return
	}
	n.debug("registry fetched: latest=%s", v)
	n.writeCache(v)
}

// newInfo constructs a notice payload for the given latest version. Callers
// generally don't need this — Lookup populates State.Info when appropriate.
func (n *Notifier) newInfo(latest string) *Info {
	return &Info{
		Current: n.Version,
		Latest:  latest,
		Message: fmt.Sprintf("A new version of cubox-cli is available: %s -> %s", n.Version, latest),
		Command: n.updateCommand(),
	}
}

// shouldCheck returns true when the update machinery should run for this
// process. It rejects:
//   - nil receivers and unset/dev versions,
//   - non-release versions (git-describe with commit suffix),
//   - explicit opt-outs (CUBOX_NO_UPDATE_NOTIFIER, NO_UPDATE_NOTIFIER),
//   - common CI environments where a notice is noise.
func (n *Notifier) shouldCheck() bool {
	if n == nil {
		return false
	}
	if n.Version == "" || n.Version == "dev" {
		return false
	}
	if !isReleaseVersion(n.Version) {
		return false
	}
	if n.isOptedOut() {
		return false
	}
	if n.isCIEnv() {
		return false
	}
	return true
}

// gitDescribePattern matches `git describe` style commit-count suffixes,
// e.g. `1.2.3-12-gabc1234` or `v1.2.3-1-g0123456-dirty`. Such versions
// indicate a local dev build and should not nag for updates.
var gitDescribePattern = regexp.MustCompile(`-\d+-g[0-9a-f]{7,}`)

func isReleaseVersion(v string) bool {
	if _, _, ok := parseSemver(v); !ok {
		return false
	}
	return !gitDescribePattern.MatchString(v)
}

func (n *Notifier) getenv(key string) string {
	if n.Getenv != nil {
		return n.Getenv(key)
	}
	return os.Getenv(key)
}

func (n *Notifier) isOptedOut() bool {
	for _, key := range []string{"CUBOX_NO_UPDATE_NOTIFIER", "NO_UPDATE_NOTIFIER"} {
		if n.getenv(key) != "" {
			return true
		}
	}
	return false
}

func (n *Notifier) isCIEnv() bool {
	for _, key := range []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"BUILD_NUMBER",
		"RUN_ID",
		"GITHUB_ACTIONS",
	} {
		if n.getenv(key) != "" {
			return true
		}
	}
	return false
}

func (n *Notifier) cacheTTL() time.Duration {
	if n.CacheTTL <= 0 {
		return defaultCacheTTL
	}
	return n.CacheTTL
}

func (n *Notifier) negativeTTL() time.Duration {
	if n.NegativeTTL <= 0 {
		return defaultNegativeTTL
	}
	return n.NegativeTTL
}

func (n *Notifier) registryURL() string {
	if n.RegistryURL != "" {
		return n.RegistryURL
	}
	return defaultRegistryURL
}

func (n *Notifier) updateCommand() string {
	if n.UpdateCommand != "" {
		return n.UpdateCommand
	}
	return defaultUpdateCmd
}

func (n *Notifier) cachePath() string {
	if n.CachePath != "" {
		return n.CachePath
	}
	return DefaultCachePath()
}

func (n *Notifier) now() time.Time {
	if n.Now != nil {
		return n.Now()
	}
	return time.Now()
}

func (n *Notifier) debug(format string, args ...any) {
	if n.Debug != nil {
		n.Debug(format, args...)
	}
}

func (n *Notifier) httpClient() *http.Client {
	if n.HTTPClient != nil {
		return n.HTTPClient
	}
	timeout := n.HTTPTimeout
	if timeout <= 0 {
		timeout = defaultHTTPTimeout
	}
	return &http.Client{Timeout: timeout}
}

func (n *Notifier) fetchRegistry(ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, n.registryURL(), nil)
	if err != nil {
		return "", err
	}
	if n.UserAgent != "" {
		req.Header.Set("User-Agent", n.UserAgent)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := n.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registry returned HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MiB ceiling
	if err != nil {
		return "", err
	}
	var payload struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.Version == "" {
		return "", fmt.Errorf("registry response missing version field")
	}
	return payload.Version, nil
}

func (n *Notifier) readCache() (cacheFile, bool) {
	var c cacheFile
	p := n.cachePath()
	if p == "" {
		return c, false
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if !os.IsNotExist(err) {
			n.debug("cache read failed: %v", err)
		}
		return c, false
	}
	if err := json.Unmarshal(data, &c); err != nil {
		n.debug("cache parse failed: %v", err)
		return c, false
	}
	return c, true
}

func (n *Notifier) writeCache(version string) {
	n.persist(cacheFile{
		LatestVersion: version,
		CheckedAt:     n.now().UTC(),
	})
}

// writeNegativeCache records a fetch failure while preserving the last known
// good version so existing notices still display.
func (n *Notifier) writeNegativeCache() {
	prev, _ := n.readCache()
	n.persist(cacheFile{
		LatestVersion: prev.LatestVersion,
		CheckedAt:     prev.CheckedAt,
		LastFailure:   n.now().UTC(),
	})
}

// persist writes the cache atomically: marshal → temp file → rename. On
// POSIX rename is atomic; on Windows os.Rename uses MoveFileEx with
// MOVEFILE_REPLACE_EXISTING (Go ≥1.5), giving the same guarantee. A
// concurrent reader therefore sees either the previous or the next complete
// file, never a torn write.
func (n *Notifier) persist(c cacheFile) {
	p := n.cachePath()
	if p == "" {
		return
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		n.debug("cache mkdir failed: %v", err)
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		n.debug("cache marshal failed: %v", err)
		return
	}
	tmp, err := os.CreateTemp(dir, filepath.Base(p)+".tmp-*")
	if err != nil {
		n.debug("cache temp create failed: %v", err)
		return
	}
	tmpName := tmp.Name()
	committed := false
	defer func() {
		if !committed {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		n.debug("cache write failed: %v", err)
		return
	}
	if err := tmp.Close(); err != nil {
		n.debug("cache close failed: %v", err)
		return
	}
	if err := os.Chmod(tmpName, 0600); err != nil {
		n.debug("cache chmod failed: %v", err)
		return
	}
	if err := os.Rename(tmpName, p); err != nil {
		n.debug("cache rename failed: %v", err)
		return
	}
	committed = true
}

// isNewer reports whether candidate is strictly greater than baseline under
// SemVer 2.0.0 precedence. Both inputs may use an optional leading "v" and
// may include pre-release identifiers (e.g. "1.2.0-rc.1") and build
// metadata (which is ignored). Returns false on any parse error so update
// notices stay conservative.
func isNewer(candidate, baseline string) bool {
	return compareSemver(candidate, baseline) > 0
}

func compareSemver(a, b string) int {
	an, ap, aok := parseSemver(a)
	bn, bp, bok := parseSemver(b)
	if !aok || !bok {
		return 0
	}
	for i := 0; i < 3; i++ {
		if an[i] != bn[i] {
			if an[i] > bn[i] {
				return 1
			}
			return -1
		}
	}
	return comparePrerelease(ap, bp)
}

func parseSemver(v string) (nums [3]int, prerelease string, ok bool) {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexByte(v, '+'); i >= 0 {
		v = v[:i]
	}
	if i := strings.IndexByte(v, '-'); i >= 0 {
		prerelease = v[i+1:]
		v = v[:i]
	}
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return [3]int{}, "", false
	}
	for i, p := range parts {
		if p == "" {
			return [3]int{}, "", false
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return [3]int{}, "", false
		}
		nums[i] = n
	}
	return nums, prerelease, true
}

// comparePrerelease implements SemVer 2.0.0 §11.4 precedence rules:
//
//   - "no prerelease" > "any prerelease".
//   - Identifiers are compared left-to-right; numeric < non-numeric.
//   - Numeric identifiers compare numerically; ASCII otherwise.
//   - More fields beat fewer when all shared fields are equal.
func comparePrerelease(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}
	aIDs := strings.Split(a, ".")
	bIDs := strings.Split(b, ".")
	n := len(aIDs)
	if len(bIDs) < n {
		n = len(bIDs)
	}
	for i := 0; i < n; i++ {
		if c := compareIdentifier(aIDs[i], bIDs[i]); c != 0 {
			return c
		}
	}
	switch {
	case len(aIDs) > len(bIDs):
		return 1
	case len(aIDs) < len(bIDs):
		return -1
	}
	return 0
}

func compareIdentifier(a, b string) int {
	aNum, aErr := strconv.Atoi(a)
	bNum, bErr := strconv.Atoi(b)
	switch {
	case aErr == nil && bErr == nil:
		switch {
		case aNum > bNum:
			return 1
		case aNum < bNum:
			return -1
		}
		return 0
	case aErr == nil:
		return -1
	case bErr == nil:
		return 1
	default:
		switch {
		case a > b:
			return 1
		case a < b:
			return -1
		}
		return 0
	}
}
