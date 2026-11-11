package cli

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/usage"
)

// expiredWork gives work an expired login and the token server answer,
// called after sent is closed.
func expiredWork(t *testing.T, c *testCLI, answer func()) (cred string, sent chan struct{}) {
	t.Helper()
	sent = make(chan struct{})
	var once sync.Once
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(sent) })
		answer()
		fmt.Fprint(w, `{"access_token":"pct-5","refresh_token":"FAKE-RT-NEW","expires_in":3600}`)
	}))
	t.Cleanup(srv.Close)
	orig := newFetcher
	newFetcher = func(dir string) *usage.Fetcher {
		f := orig(dir)
		f.TokenURL = srv.URL
		return f
	}
	t.Cleanup(func() { newFetcher = orig })
	cred = filepath.Join(c.p.ProfileDir("work"), ".credentials.json")
	os.WriteFile(cred, []byte(`{"claudeAiOauth":{"accessToken":"stale","refreshToken":"FAKE-RT-OLD","expiresAt":1}}`), 0o600)
	return cred, sent
}

// The dashboard's r renews an expired login, as `ccam usage --refresh`
// does, so it shows the usage of now instead of the saved one. Opening the
// dashboard renews nothing, and neither does a dry run.
func TestDashboardReloadRenewsLogin(t *testing.T) {
	c := newTestCLI(t, nil)
	cred, _ := expiredWork(t, c, func() {})
	e, _ := c.env()
	b := backend{e: e, renew: newRenewals()}

	if r := b.Usage(context.Background(), "work", false); r.Status != usage.StatusExpired {
		t.Fatalf("opening the dashboard: %+v", r)
	}
	e.dryRun = true
	if r := b.Usage(context.Background(), "work", true); r.Status != usage.StatusExpired {
		t.Fatalf("r in a dry run: %+v", r)
	}
	e.dryRun = false
	r := b.Usage(context.Background(), "work", true)
	if r.Status != usage.StatusOK || r.Source != "live" || r.Limits.FiveHour.Percent != 5 {
		t.Fatalf("r: %+v", r)
	}
	if data, _ := os.ReadFile(cred); !strings.Contains(string(data), "FAKE-RT-NEW") {
		t.Fatalf("the renewed login was not saved: %s", data)
	}
}

// Leaving the dashboard waits for a renewal r sent: the server may have
// rotated the refresh token already, and only its answer holds the login.
func TestDashboardWaitsForRenewal(t *testing.T) {
	c := newTestCLI(t, nil)
	cred, sent := expiredWork(t, c, func() { time.Sleep(300 * time.Millisecond) })
	e, _ := c.env()
	b := backend{e: e, renew: newRenewals()}

	go b.Usage(context.Background(), "work", true)
	<-sent
	b.renew.end()
	if data, _ := os.ReadFile(cred); !strings.Contains(string(data), "FAKE-RT-NEW") {
		t.Fatalf("the dashboard ended before the renewal was saved: %s", data)
	}
}

// The dashboard marks an account that holds an unused limit reset: what
// Claude Code saved in the account's .claude.json, read with its usage.
func TestDashboardUsageSaysWhoHoldsAReset(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 5, "work": 5})
	now := time.Now()
	saved := fmt.Sprintf(`{"oauthAccount":{"accountUuid":"u-1"},"unusedLimitResetStatusCache":{"accountUuid":"u-1","holdsUnusedReset":true,"checkedAt":%d,"validUntil":%d}}`,
		now.Add(-time.Hour).UnixMilli(), now.Add(time.Hour).UnixMilli())
	if err := os.WriteFile(filepath.Join(c.p.ProfileDir("work"), ".claude.json"), []byte(saved), 0o600); err != nil {
		t.Fatal(err)
	}
	e, _ := c.env()
	b := backend{e: e}
	for name, want := range map[string]bool{"work": true, "personal": false} {
		if r := b.Usage(context.Background(), name, false); r.Status != usage.StatusOK || r.HoldsReset != want {
			t.Errorf("%s: %+v", name, r)
		}
	}
}
