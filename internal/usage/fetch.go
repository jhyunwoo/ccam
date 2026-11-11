package usage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/profile"
)

// Endpoint is the undocumented usage API (docs/claude-layout.md).
const Endpoint = "https://api.anthropic.com/api/oauth/usage"

const (
	CacheTTL       = 60 * time.Second
	connectTimeout = 2 * time.Second
	accountTimeout = 3 * time.Second
)

// LookupTimeout bounds one account's live lookup.
const LookupTimeout = accountTimeout

type Status string

const (
	StatusOK        Status = "ok"
	StatusExpired   Status = "expired"    // token expired; values may come from cache
	StatusLoggedOut Status = "logged_out" // no credentials
	StatusAPI       Status = "api"        // API key account, no subscription limits
	StatusError     Status = "error"      // request failed; values may come from cache
	StatusPending   Status = "pending"    // not answered yet (live view only)
)

// Account is one profile to query.
type Account struct {
	Name       string
	Slot       int
	Dir        string // CLAUDE_CONFIG_DIR
	SessionDir string // markers of running sessions started by ccam
	Imported   bool   // registered in place; plain `claude` may use it without ccam's lock
	API        bool   // an API key account (ccam add --api): no subscription limits
}

// Report is one account's result.
type Report struct {
	Name   string    `json:"name"`
	Slot   int       `json:"slot"`
	Status Status    `json:"status"`
	Limits *Limits   `json:"limits,omitempty"`
	Source string    `json:"source,omitempty"` // live | cache | statusline | limit
	At     time.Time `json:"at,omitzero"`      // when the values were measured
	Err    string    `json:"error,omitempty"`
	// Limited is when Claude Code reported the usage limit that set the
	// windows LimitedWindows names ("5-hour", "Weekly", a model) to 100%
	// until their reset, over what the source said; zero when none.
	Limited        time.Time `json:"limited,omitzero"`
	LimitedWindows []string  `json:"limited_windows,omitempty"`
	// FormatChanged means the endpoint answered but the parser did not
	// understand it; doctor reports this.
	FormatChanged bool `json:"format_changed,omitempty"`
	// Sessions is how many ccam sessions of the account run now; the
	// caller fills it in for the balance strategy.
	Sessions int `json:"sessions,omitempty"`
	// HoldsReset means the account holds an unused limit reset, as Claude
	// Code last saw it (profile.HoldsReset); the dashboard fills it in.
	HoldsReset bool `json:"holds_reset,omitempty"`
}

type Options struct {
	NoCache bool
	// Source is "live" (ccam asks the usage endpoint), "claude" (Claude Code
	// does, see ProbeRequest), "cache" (saved values only), or "" for the
	// Fetcher's Source.
	Source  string
	Refresh bool // refresh expired tokens (only when the profile is not running)
}

// Fetcher queries accounts in parallel over one shared HTTP client.
type Fetcher struct {
	Client *http.Client
	// TokenClient sends token renewals: Client without its 3 s wait for an
	// answer, which refreshTimeout bounds instead (nil: Client).
	TokenClient *http.Client
	URL         string
	TokenURL    string
	StateDir    string // ~/.local/state/ccam/usage
	Home        string // names the macOS Keychain item of an imported ~/.claude
	Now         func() time.Time
	// Source is what Options.Source "" means (config usage_source).
	Source string
	// Probe runs Claude Code in the account's folder with ProbeRequest on
	// its standard input and returns its output.
	Probe func(ctx context.Context, a Account) ([]byte, error)
}

func NewFetcher(stateDir string) *Fetcher {
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: connectTimeout}).DialContext,
		TLSHandshakeTimeout:   connectTimeout,
		ResponseHeaderTimeout: accountTimeout,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   8,
	}
	// Neither endpoint redirects; following one would resend the token (or
	// the refresh token in a POST body) elsewhere.
	noRedirect := func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	token := tr.Clone()
	token.ResponseHeaderTimeout = 0
	return &Fetcher{
		Client:      &http.Client{Transport: tr, CheckRedirect: noRedirect},
		TokenClient: &http.Client{Transport: token, CheckRedirect: noRedirect},
		URL:         Endpoint,
		TokenURL:    TokenURL,
		StateDir:    stateDir,
		Now:         time.Now,
	}
}

// FetchAll queries every account at once and calls onResult as each one
// finishes (from multiple goroutines, serialized). It returns all reports
// in input order.
func (f *Fetcher) FetchAll(ctx context.Context, accts []Account, opt Options, onResult func(int, Report)) []Report {
	out := make([]Report, len(accts))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i, a := range accts {
		wg.Go(func() {
			r := f.Fetch(ctx, a, opt)
			mu.Lock()
			out[i] = r
			if onResult != nil {
				onResult(i, r)
			}
			mu.Unlock()
		})
	}
	wg.Wait()
	return out
}

// Fetch queries one account, falling back to cached values on failure. A
// window Claude Code has reported full (SaveLimit) counts as full until its
// reset, whatever the values say.
func (f *Fetcher) Fetch(ctx context.Context, a Account, opt Options) Report {
	return f.withLimit(a, f.fetch(ctx, a, opt))
}

func (f *Fetcher) fetch(ctx context.Context, a Account, opt Options) Report {
	r := Report{Name: a.Name, Slot: a.Slot}
	if a.API {
		r.Status = StatusAPI
		return r
	}
	if opt.Source == "" && f.Source == "claude" {
		opt.Source = f.Source // the config picks live or claude; cache is asked for per call
	}
	if opt.Source == "claude" {
		return f.fetchClaude(ctx, a, opt)
	}
	if opt.Source == "cache" {
		r = f.fallback(a, r, StatusOK, "")
		if r.Limits == nil {
			// Nothing saved: say why when the account has no working login.
			switch _, _, err := readCredentials(a.Dir, f.Home, f.Now()); {
			case errors.Is(err, errLoginExpired):
				r.Status, r.Err = StatusExpired, loginExpired(a.Name)
			case errors.Is(err, errNoOAuth):
				r.Err = ""
				r.Status = StatusLoggedOut
				if creds.APIKey(a.Dir) {
					r.Status = StatusAPI
				}
			}
		}
		return r
	}
	if !opt.NoCache {
		if c, ok := f.readCache(a.Name, "live"); ok && fresh(c.At, f.Now()) {
			r.Status, r.Limits, r.Source, r.At = StatusOK, c.Limits, "live", c.At
			return r
		}
	}
	cred, src, err := readCredentials(a.Dir, f.Home, f.Now())
	switch {
	case errors.Is(err, errLoginExpired):
		return f.fallback(a, r, StatusExpired, loginExpired(a.Name))
	case errors.Is(err, errNoOAuth):
		if creds.APIKey(a.Dir) {
			r.Status = StatusAPI
			return r
		}
		r.Status = StatusLoggedOut
		return r
	case errors.Is(err, creds.ErrCredentialManager):
		return f.fallback(a, r, StatusError, err.Error()+`; usage_source = "claude" asks Claude Code instead`)
	case err != nil:
		return f.fallback(a, r, StatusError, err.Error())
	}
	if cred.expired(f.Now()) {
		if !opt.Refresh {
			return f.fallback(a, r, StatusExpired, fmt.Sprintf("access token expired; `ccam %s` renews it when it starts", a.Name))
		}
		if src == creds.Keychain {
			return f.fallback(a, r, StatusExpired, fmt.Sprintf("token not renewed: ccam does not renew macOS Keychain logins; `ccam %s` renews it when it starts", a.Name))
		}
		// Not bounded by the lookup's 3 s: a renewal may take as long as
		// Claude Code lets its own (refreshTimeout).
		if cred, err = f.refresh(ctx, a); err != nil {
			// Most refusals leave a working login (Claude Code renews it);
			// a renewal ccam lost says so itself.
			why := err.Error()
			if !errors.Is(err, errLoginLost) {
				why = "token not renewed: " + why
			}
			return f.fallback(a, r, StatusExpired, why)
		}
	}
	if until, ok := f.paused(a.Name); ok {
		return f.fallback(a, r, StatusError, pausedMsg(until))
	}
	ctx, cancel := context.WithTimeout(ctx, accountTimeout)
	defer cancel()
	limits, status, err := f.live(ctx, cred.AccessToken)
	var limited rateLimited
	if errors.As(err, &limited) {
		until := f.Now().Add(retryAfter(limited.retryAfter, f.Now()))
		_ = fsutil.WriteFast(f.pausePath(a.Name), []byte(until.UTC().Format(time.RFC3339)))
		return f.fallback(a, r, StatusError, pausedMsg(until))
	}
	if err != nil {
		r.FormatChanged = errors.Is(err, errFormat)
		return f.fallback(a, r, status, err.Error())
	}
	r.Status, r.Limits, r.Source, r.At = StatusOK, &limits, "live", f.Now()
	f.writeCache(a.Name, "live", cached{At: r.At, Limits: &limits})
	return r
}

func (f *Fetcher) live(ctx context.Context, token string) (Limits, Status, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.URL, nil)
	if err != nil {
		return Limits{}, StatusError, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	req.Header.Set("User-Agent", "ccam")
	resp, err := f.Client.Do(req)
	if err != nil {
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			return Limits{}, StatusError, errors.New("cancelled")
		case ctx.Err() != nil:
			return Limits{}, StatusError, fmt.Errorf("no response (over %s)", accountTimeout)
		}
		return Limits{}, StatusError, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Limits{}, StatusError, err
	}
	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		return Limits{}, StatusError, rateLimited{retryAfter: resp.Header.Get("Retry-After")}
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return Limits{}, StatusExpired, fmt.Errorf("login expired: authentication failed (%d)", resp.StatusCode)
	case resp.StatusCode != http.StatusOK:
		return Limits{}, StatusError, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	l, err := ParseLive(body)
	if err != nil {
		return Limits{}, StatusError, fmt.Errorf("%w: %v", errFormat, err)
	}
	return l, StatusOK, nil
}

// fallback fills r with the saved values: the newest lookup (live or
// claude), with the windows a newer status line knows on top. The status
// line has no model windows and leaves out a window that is not running,
// so it never replaces a lookup's values wholesale.
func (f *Fetcher) fallback(a Account, r Report, status Status, errMsg string) Report {
	r.Status, r.Err = status, errMsg
	var base cached
	found := false
	for _, src := range []string{"live", "claude"} {
		if c, ok := f.readCache(a.Name, src); ok && (!found || c.At.After(base.At)) {
			base, found = c, true
		}
	}
	line, haveLine := f.readCache(a.Name, "statusline")
	switch {
	case !found && !haveLine:
		if status == StatusOK {
			r.Status, r.Err = StatusError, "no saved usage"
		}
		return r
	case !found:
		r.Limits, r.At, r.Source = line.Limits, line.At, "statusline"
		return r
	}
	r.Limits, r.At, r.Source = base.Limits, base.At, "cache"
	if haveLine && line.At.After(base.At) {
		if line.Limits.FiveHour.Known {
			r.Limits.FiveHour = line.Limits.FiveHour
		}
		if line.Limits.Weekly.Known {
			r.Limits.Weekly = line.Limits.Weekly
		}
		r.At, r.Source = line.At, "statusline"
	}
	return r
}

type cached struct {
	At     time.Time `json:"at"`
	Limits *Limits   `json:"limits"`
}

// fresh reports whether values saved at at stand in for a lookup at now. A
// time after now (the clock was set back) is not fresh.
func fresh(at, now time.Time) bool {
	age := now.Sub(at)
	return age >= 0 && age < CacheTTL
}

func (f *Fetcher) cachePath(name, src string) string {
	if src == "live" {
		return filepath.Join(f.StateDir, name+".json")
	}
	return filepath.Join(f.StateDir, name+"."+src+".json")
}

// CachePaths lists the cache files of account name in stateDir, its limit
// record among them.
func CachePaths(stateDir, name string) []string {
	f := Fetcher{StateDir: stateDir}
	return []string{f.cachePath(name, "live"), f.cachePath(name, "claude"), f.cachePath(name, "statusline"), f.pausePath(name), f.limitPath(name)}
}

// The endpoint answers HTTP 429 when asked too often. ccam then leaves the
// account alone for Retry-After (or defaultPause), up to maxPause, and shows
// saved numbers meanwhile, so retries do not keep the limit going.
const (
	defaultPause = 5 * time.Minute
	maxPause     = time.Hour
)

type rateLimited struct{ retryAfter string }

func (rateLimited) Error() string { return "rate limited" }

// retryAfter reads a Retry-After value: seconds or an HTTP date.
func retryAfter(v string, now time.Time) time.Duration {
	d := defaultPause
	if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
		d = time.Duration(secs) * time.Second
	} else if at, err := http.ParseTime(v); err == nil && at.After(now) {
		d = at.Sub(now)
	}
	return min(d, maxPause)
}

func (f *Fetcher) pausePath(name string) string {
	return filepath.Join(f.StateDir, name+".pause")
}

// paused reports whether the endpoint asked to wait for this account.
func (f *Fetcher) paused(name string) (time.Time, bool) {
	data, err := os.ReadFile(f.pausePath(name))
	if err != nil {
		return time.Time{}, false
	}
	until, err := time.Parse(time.RFC3339, strings.TrimSpace(string(data)))
	return until, err == nil && f.Now().Before(until)
}

func pausedMsg(until time.Time) string {
	return "rate limited by the usage endpoint; asking again after " + until.Local().Format("15:04")
}

func (f *Fetcher) readCache(name, src string) (cached, bool) {
	data, err := os.ReadFile(f.cachePath(name, src))
	if err != nil {
		return cached{}, false
	}
	var c cached
	if json.Unmarshal(data, &c) != nil || c.Limits == nil {
		return cached{}, false
	}
	// A window that has reset since the values were saved has started over.
	c.Limits = c.Limits.as(f.Now())
	return c, true
}

func (f *Fetcher) writeCache(name, src string, c cached) {
	data, _ := json.Marshal(c)
	_ = fsutil.WriteFast(f.cachePath(name, src), data)
}

// SaveStatusLine stores values captured by the status line collector.
func SaveStatusLine(stateDir, name string, l Limits, at time.Time) error {
	data, err := json.Marshal(cached{At: at, Limits: &l})
	if err != nil {
		return err
	}
	return fsutil.WriteFast(filepath.Join(stateDir, name+".statusline.json"), data)
}

var (
	errNoOAuth = errors.New("no subscription login")
	// errLoginExpired is a login Claude Code cannot renew (creds.Judge).
	errLoginExpired = errors.New("login expired")
	errFormat       = errors.New("the usage response format changed")
)

type oauthCred struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"` // Unix ms
	ClientID     string `json:"clientId"`  // the OAuth client the login was made with ("" = Claude Code's)
	// refreshExpiresAt is when a renewed refresh token expires (Unix ms),
	// when the token endpoint said so; only a renewal sets it.
	refreshExpiresAt int64
}

// loginExpired is the error of a login Claude Code cannot renew.
func loginExpired(name string) string {
	return fmt.Sprintf("login expired → `ccam login %s`", name)
}

func (c oauthCred) expired(now time.Time) bool {
	return c.ExpiresAt > 0 && now.UnixMilli() >= c.ExpiresAt
}

// readCredentials reads only; ccam never rewrites the token file except
// for an explicit --refresh.
func readCredentials(dir, home string, now time.Time) (oauthCred, creds.Source, error) {
	var c oauthCred
	data, src, err := creds.Read(dir, home)
	if errors.Is(err, creds.ErrNone) {
		return c, src, errNoOAuth
	}
	if err != nil {
		return c, src, err
	}
	if creds.Judge(data, now) == "login expired" {
		return c, src, errLoginExpired
	}
	raw, ok, err := profile.TopLevelValue(data, "claudeAiOauth")
	if err == nil && ok {
		err = json.Unmarshal(raw, &c)
	}
	if err == nil && c.AccessToken == "" {
		err = errNoOAuth
	}
	return c, src, err
}
