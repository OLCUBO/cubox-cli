package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	registryURL = "https://registry.npmjs.org/cubox-cli/latest"
	cacheTTL    = 6 * time.Hour
	httpTimeout = 3 * time.Second
	updateCmd   = "npm update -g cubox-cli && npx skills add OLCUBO/cubox-cli -g -y"
)

// Info holds the update notification payload included in _notice.update.
type Info struct {
	Current string `json:"current"`
	Latest  string `json:"latest"`
	Message string `json:"message"`
	Command string `json:"command"`
}

type cacheFile struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// Check fetches the latest published version from the npm registry and
// compares it with currentVersion. Returns nil when up-to-date, when the
// check fails, or when the cached value is still fresh. Never panics.
func Check(currentVersion string) *Info {
	if currentVersion == "" || currentVersion == "dev" {
		return nil
	}
	latest := latestVersion()
	if latest == "" {
		return nil
	}
	if !isNewer(latest, currentVersion) {
		return nil
	}
	return &Info{
		Current: currentVersion,
		Latest:  latest,
		Message: fmt.Sprintf("A new version of cubox-cli is available: %s -> %s", currentVersion, latest),
		Command: updateCmd,
	}
}

// latestVersion returns the latest version string from cache or the npm registry.
// Returns "" on any error.
func latestVersion() string {
	if v := readCache(); v != "" {
		return v
	}
	v := fetchRegistry()
	if v != "" {
		writeCache(v)
	}
	return v
}

func fetchRegistry() string {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(registryURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var payload struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	return payload.Version
}

func cachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "cubox-cli", "update-check.json")
}

func readCache() string {
	p := cachePath()
	if p == "" {
		return ""
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	var c cacheFile
	if err := json.Unmarshal(data, &c); err != nil {
		return ""
	}
	if time.Since(c.CheckedAt) > cacheTTL {
		return ""
	}
	return c.LatestVersion
}

func writeCache(version string) {
	p := cachePath()
	if p == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(p), 0700)
	data, err := json.Marshal(cacheFile{
		LatestVersion: version,
		CheckedAt:     time.Now().UTC(),
	})
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0600)
}

// isNewer returns true when candidate is strictly greater than baseline.
// Both must be dot-separated numeric version strings (e.g. "0.2.0").
// Returns false on any parse error.
func isNewer(candidate, baseline string) bool {
	cv := parseSemver(candidate)
	bv := parseSemver(baseline)
	if cv == nil || bv == nil {
		return false
	}
	for i := 0; i < 3; i++ {
		if cv[i] > bv[i] {
			return true
		}
		if cv[i] < bv[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) []int {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return nil
	}
	nums := make([]int, 3)
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil
		}
		nums[i] = n
	}
	return nums
}
