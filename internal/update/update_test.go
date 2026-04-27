package update

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestNotifier returns a Notifier wired to a per-test temp cache file and
// a frozen clock, with all environment lookups returning empty by default
// (so host CI vars / opt-out vars don't leak into test outcomes).
func newTestNotifier(t *testing.T, version string) *Notifier {
	t.Helper()
	dir := t.TempDir()
	return &Notifier{
		Version:       version,
		RegistryURL:   defaultRegistryURL,
		CachePath:     filepath.Join(dir, "update-check.json"),
		CacheTTL:      defaultCacheTTL,
		NegativeTTL:   defaultNegativeTTL,
		HTTPTimeout:   defaultHTTPTimeout,
		UpdateCommand: defaultUpdateCmd,
		UserAgent:     "cubox-cli/test",
		Now:           time.Now,
		Getenv:        func(string) string { return "" },
	}
}

func writeCacheFile(t *testing.T, n *Notifier, c cacheFile) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(n.CachePath), 0700); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(n.CachePath, data, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestLookupFreshCacheReturnsNoticeWithoutRefresh(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.8",
		CheckedAt:     time.Now().UTC(),
	})

	state := n.Lookup()

	if state.Info == nil {
		t.Fatal("expected update info")
	}
	if state.Info.Latest != "1.0.8" || state.Info.Current != "1.0.0" {
		t.Fatalf("unexpected update info: %+v", state.Info)
	}
	if state.NeedsRefresh {
		t.Fatal("fresh cache should not need refresh")
	}
	if state.CacheEmpty {
		t.Fatal("fresh cache should not be empty")
	}
}

func TestLookupStaleCacheReturnsNoticeAndNeedsRefresh(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.8",
		CheckedAt:     time.Now().Add(-defaultCacheTTL - time.Minute).UTC(),
	})

	state := n.Lookup()

	if state.Info == nil {
		t.Fatal("expected stale update info to be usable")
	}
	if !state.NeedsRefresh {
		t.Fatal("stale cache should need refresh")
	}
	if state.CacheEmpty {
		t.Fatal("stale cache with a value should not be empty")
	}
}

func TestLookupStaleCacheWithoutNoticeStillNeedsRefresh(t *testing.T) {
	n := newTestNotifier(t, "1.0.6")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.5",
		CheckedAt:     time.Now().Add(-defaultCacheTTL - time.Minute).UTC(),
	})

	state := n.Lookup()

	if state.Info != nil {
		t.Fatalf("did not expect update info: %+v", state.Info)
	}
	if !state.NeedsRefresh {
		t.Fatal("stale cache should need refresh even when it has no notice")
	}
}

func TestLookupEmptyCacheNeedsRefresh(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")

	state := n.Lookup()

	if state.Info != nil {
		t.Fatalf("did not expect update info: %+v", state.Info)
	}
	if !state.NeedsRefresh {
		t.Fatal("empty cache should need refresh")
	}
	if !state.CacheEmpty {
		t.Fatal("empty cache should be marked empty")
	}
}

func TestLookupNegativeCacheSuppressesRefresh(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.8",
		CheckedAt:     time.Now().Add(-defaultCacheTTL - time.Minute).UTC(),
		LastFailure:   time.Now().Add(-time.Minute).UTC(),
	})

	state := n.Lookup()

	if state.NeedsRefresh {
		t.Fatal("recent failure should suppress refresh inside NegativeTTL")
	}
	if state.Info == nil {
		t.Fatal("stale-but-known notice should still surface during negative-TTL window")
	}
}

func TestLookupNegativeCacheExpiresAfterTTL(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.8",
		CheckedAt:     time.Now().Add(-defaultCacheTTL - time.Minute).UTC(),
		LastFailure:   time.Now().Add(-defaultNegativeTTL - time.Minute).UTC(),
	})

	state := n.Lookup()

	if !state.NeedsRefresh {
		t.Fatal("after NegativeTTL elapses, refresh should be allowed again")
	}
}

func TestLookupSkipsForDevAndUnsetVersions(t *testing.T) {
	for _, v := range []string{"", "dev"} {
		n := newTestNotifier(t, v)
		state := n.Lookup()
		if state.Info != nil || state.NeedsRefresh || state.CacheEmpty {
			t.Fatalf("version %q should skip update lookup, got %+v", v, state)
		}
	}
}

func TestLookupSkipsForGitDescribeBuilds(t *testing.T) {
	cases := []string{
		"1.0.0-12-g9b933f1",
		"v1.0.0-1-g0123456-dirty",
		"v0.2.0-7-gabcdef0",
	}
	for _, v := range cases {
		n := newTestNotifier(t, v)
		writeCacheFile(t, n, cacheFile{
			LatestVersion: "9.9.9",
			CheckedAt:     time.Now().UTC(),
		})
		state := n.Lookup()
		if state.Info != nil || state.NeedsRefresh {
			t.Fatalf("git-describe build %q should skip lookup, got %+v", v, state)
		}
	}
}

func TestLookupSkipsWhenOptedOut(t *testing.T) {
	for _, key := range []string{"CUBOX_NO_UPDATE_NOTIFIER", "NO_UPDATE_NOTIFIER"} {
		n := newTestNotifier(t, "1.0.0")
		n.Getenv = func(k string) string {
			if k == key {
				return "1"
			}
			return ""
		}
		writeCacheFile(t, n, cacheFile{
			LatestVersion: "9.9.9",
			CheckedAt:     time.Now().UTC(),
		})
		if state := n.Lookup(); state.Info != nil || state.NeedsRefresh {
			t.Fatalf("%s should suppress checks, got %+v", key, state)
		}
	}
}

func TestLookupSkipsInCIEnvironments(t *testing.T) {
	for _, key := range []string{"CI", "CONTINUOUS_INTEGRATION", "BUILD_NUMBER", "RUN_ID", "GITHUB_ACTIONS"} {
		n := newTestNotifier(t, "1.0.0")
		n.Getenv = func(k string) string {
			if k == key {
				return "true"
			}
			return ""
		}
		writeCacheFile(t, n, cacheFile{
			LatestVersion: "9.9.9",
			CheckedAt:     time.Now().UTC(),
		})
		if state := n.Lookup(); state.Info != nil || state.NeedsRefresh {
			t.Fatalf("%s=true should suppress checks, got %+v", key, state)
		}
	}
}

func TestRefreshCacheWritesRegistryVersion(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "cubox-cli/test" {
			t.Errorf("expected User-Agent header, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.0.8"}`))
	}))
	t.Cleanup(server.Close)
	n.RegistryURL = server.URL
	n.HTTPClient = server.Client()

	n.RefreshCache(context.Background())

	c, ok := n.readCache()
	if !ok {
		t.Fatal("expected cache to exist after RefreshCache")
	}
	if c.LatestVersion != "1.0.8" {
		t.Fatalf("expected refreshed version 1.0.8, got %q", c.LatestVersion)
	}
	if c.CheckedAt.IsZero() {
		t.Fatal("expected CheckedAt to be set on success")
	}
	if !c.LastFailure.IsZero() {
		t.Fatal("LastFailure should be cleared on success")
	}
}

func TestRefreshCacheWritesNegativeOnHTTPError(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.5",
		CheckedAt:     time.Now().Add(-defaultCacheTTL - time.Minute).UTC(),
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)
	n.RegistryURL = server.URL
	n.HTTPClient = server.Client()

	n.RefreshCache(context.Background())

	c, ok := n.readCache()
	if !ok {
		t.Fatal("expected cache file after negative refresh")
	}
	if c.LatestVersion != "1.0.5" {
		t.Fatalf("negative refresh must preserve last known version, got %q", c.LatestVersion)
	}
	if c.LastFailure.IsZero() {
		t.Fatal("expected LastFailure to be recorded on HTTP error")
	}
}

func TestRefreshCacheRespectsContextCancel(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")

	gotRequest := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(gotRequest)
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)
	n.RegistryURL = server.URL
	n.HTTPClient = server.Client()
	n.HTTPTimeout = 30 * time.Second // ensure cancellation, not the client deadline, is what stops us

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		n.RefreshCache(ctx)
	}()

	select {
	case <-gotRequest:
	case <-time.After(2 * time.Second):
		t.Fatal("server never received the request")
	}
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("RefreshCache did not return after context cancel")
	}

	c, ok := n.readCache()
	if !ok {
		t.Fatal("expected negative cache after cancel")
	}
	if c.LastFailure.IsZero() {
		t.Fatal("cancellation should be recorded as a failure")
	}
}

func TestRefreshCacheWriteIsAtomic(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	writeCacheFile(t, n, cacheFile{
		LatestVersion: "1.0.5",
		CheckedAt:     time.Now().UTC(),
	})

	// Background reader pounds the cache file while we rewrite it. Any torn
	// write would surface as a JSON unmarshal error here.
	stop := make(chan struct{})
	defer close(stop)
	var torn atomic.Int32
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			data, err := os.ReadFile(n.CachePath)
			if err != nil {
				continue
			}
			var c cacheFile
			if err := json.Unmarshal(data, &c); err != nil {
				torn.Add(1)
			}
		}
	}()

	for i := 0; i < 100; i++ {
		n.writeCache("1.0." + strings.Repeat("9", 1+i%5))
	}

	if got := torn.Load(); got != 0 {
		t.Fatalf("observed %d torn reads — persist is not atomic", got)
	}
}

func TestRefreshCacheLeavesNoTempFiles(t *testing.T) {
	n := newTestNotifier(t, "1.0.0")
	for i := 0; i < 10; i++ {
		n.writeCache("1.0.8")
	}

	entries, err := os.ReadDir(filepath.Dir(n.CachePath))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".tmp-") {
			t.Fatalf("leftover temp file %q in cache dir", e.Name())
		}
	}
}

func TestIsNewerSemantics(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.0.1", "1.0.0", true},
		{"1.0.0", "1.0.0", false},
		{"1.0.0", "1.0.1", false},
		{"v1.2.3", "1.2.2", true},
		{"1.2.3+build", "1.2.3", false},
		// SemVer 2.0.0 §11.4: GA > any prerelease of the same core
		{"1.2.3", "1.2.3-rc.1", true},
		{"1.2.3-rc.2", "1.2.3-rc.1", true},
		{"1.2.3-rc.1", "1.2.3-rc.10", false},
		{"1.2.3-alpha", "1.2.3-beta", false},
		// unparseable inputs → conservative false
		{"garbage", "1.0.0", false},
		{"1.0.0", "garbage", false},
	}
	for _, tc := range cases {
		if got := isNewer(tc.a, tc.b); got != tc.want {
			t.Errorf("isNewer(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestIsReleaseVersion(t *testing.T) {
	cases := []struct {
		v    string
		want bool
	}{
		{"1.0.0", true},
		{"v1.0.0", true},
		{"1.0.0-rc.1", true},
		{"1.0.0+build.123", true},
		{"1.0.0-12-g9b933f1", false},
		{"v1.0.0-1-g0123456-dirty", false},
		{"dev", false},
		{"", false},
		{"1.0", false},
	}
	for _, tc := range cases {
		if got := isReleaseVersion(tc.v); got != tc.want {
			t.Errorf("isReleaseVersion(%q) = %v, want %v", tc.v, got, tc.want)
		}
	}
}
