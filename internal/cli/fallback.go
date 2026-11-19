package cli

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// Fallback runs `ccam --fallback [name] -p …`: claude runs as a child with
// its output held back. When it fails because the account reached its
// usage limit, the output is dropped and the same command runs again as
// another account, at most once per account. When every account is at its
// limit it exits with 75, or with wait waits for one to reset first. It
// returns the exit code.
func Fallback(p paths.Paths, sel string, args []string, dry bool, wait string) int {
	e := &env{paths: p, stdin: os.Stdin, in: bufio.NewReader(os.Stdin), stdout: os.Stdout, stderr: os.Stderr, dryRun: dry}
	err := fallback(e, sel, args, wait)
	var code launcher.ExitCode
	var coded codedError
	switch {
	case err == nil:
		return 0
	case errors.As(err, &code):
		return int(code)
	case errors.As(err, &coded):
		fmt.Fprintln(e.stderr, "ccam:", coded.msg)
		return coded.code
	}
	fmt.Fprintln(e.stderr, "ccam:", err)
	return 1
}

func fallback(e *env, sel string, args []string, wait string) error {
	// What claude gets (a leading "--" dropped, as for any launch) is what
	// says whether the run is -p and which model it uses. args itself goes
	// to Ready, which drops the separator again.
	claude := launcher.ClaudeArgs(args)
	if !printMode(claude) {
		return fmt.Errorf("--fallback works for non-interactive runs with -p or --print. When an interactive session reaches its limit, use `ccam next --auto --continue`")
	}
	cfg, err := e.load()
	if err != nil {
		return err
	}
	if _, found := cfg.Find(sel); !found && wait == DefaultWait {
		// `--wait 30m work`: a bare --wait waits 5 h, and 30m became the
		// account (names such as 30m are allowed, hence only without one).
		if err := waitHint(sel); err != nil {
			return err
		}
	}
	p, err := choiceFlags{}.policy(cfg, claude)
	if err != nil {
		return err
	}
	waitFor, err := waitFor(wait)
	if err != nil {
		return err
	}
	pl, err := launcher.Ready(e.paths, cfg, sel, args, e.dryRun)
	if err != nil {
		return err
	}
	if e.dryRun {
		if err := pl.DryRun(e.stdout); err != nil {
			return err
		}
		fmt.Fprintln(e.stderr, "ccam: if this account reaches its usage limit, the command runs again as the logged-in account with the most room")
		return nil
	}
	// Each attempt gets the same input, so it is read once up front.
	var input []byte
	if !isTerminal(e.stdin) {
		wait := stdinWait
		if streamInput(claude) {
			wait = 0 // Claude Code reads a stream-json input as long as it runs
		}
		if input, err = readInput(e.stdin, e.stderr, wait); err != nil {
			return err
		}
	}
	tried := map[string]bool{}
	left := waitFor // --wait bounds the waiting of the whole run, not of each wait
	waited := false
	for {
		tried[pl.Profile.Name] = true
		stdin := e.stdin
		if input != nil {
			stdin = bytes.NewReader(input)
		}
		var out bytes.Buffer
		errTail := &tail{max: 64 << 10}
		code, err := pl.Run(e.paths, stdin, &out, io.MultiWriter(e.stderr, errTail))
		pl.Done()
		if err != nil {
			return err
		}
		// A run stopped by a signal (ccam's, passed on, or claude's own:
		// 128+n, 143 for SIGTERM) was asked to end, not refused: it is not
		// run again as another account.
		if code == 0 || code > 128 || pl.Signaled() || !limited(e, cfg, pl.Profile, p, out.Bytes(), errTail.buf) {
			_, _ = e.stdout.Write(out.Bytes())
			return exitWith(code)
		}
		next, _ := fallbackTarget(e, cfg, pl.Profile, tried, claude, p)
		if next == nil && waitFor > 0 {
			// Every account was tried or is full: wait for one to reset, and
			// then any account may be tried again.
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit, and no other logged-in account has room left\n", pl.Profile.Name)
			// The first look only finds out how long to wait: an account that
			// just hit its limit may still show room in saved usage, and must
			// not be run again at once.
			from, first, start := pl.Profile, true, clock()
			next, err = untilRoom(e, p, left, "no logged-in account has room left", func() (*config.Profile, []*usage.Report, error) {
				prof, reports := fallbackTarget(e, cfg, from, map[string]bool{}, claude, p)
				if first {
					first = false
					return nil, reports, nil
				}
				return prof, reports, nil
			})
			left -= clock().Sub(start)
			if err != nil {
				_, _ = e.stdout.Write(out.Bytes())
				return err
			}
			tried, waited = map[string]bool{}, true
		}
		if next == nil {
			_, _ = e.stdout.Write(out.Bytes())
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit, and no other logged-in account has room left\n", pl.Profile.Name)
			return launcher.ExitCode(exitNoRoom)
		}
		if waited {
			// The limit was said before the wait.
			fmt.Fprintf(e.stderr, "ccam: running the command again as %s\n", next.Name)
			waited = false
		} else {
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit; running the command again as %s\n", pl.Profile.Name, next.Name)
		}
		if pl, err = launcher.Ready(e.paths, cfg, next.Name, args, false); err != nil {
			return err
		}
	}
}

// stdinWait is how long piped input may take to start, as in Claude Code.
var stdinWait = 3 * time.Second

// maxInput is the most piped input claude -p takes (Claude Code refuses
// more than 10 MB).
const maxInput = 10 << 20

// readInput reads piped input to its end, so that every attempt gets it.
// Like Claude Code it waits at most wait for the first byte: a pipe that
// stays open and sends nothing (ssh without -n, Node's spawn, Python's
// Popen with stdin=PIPE) gives empty input instead of a run that never
// starts. Once input starts it is read to its end, up to maxInput. With
// wait 0 it waits for the first byte as long as it takes.
func readInput(r io.Reader, stderr io.Writer, wait time.Duration) ([]byte, error) {
	type chunk struct {
		b   []byte
		err error
	}
	got := make(chan chunk, 1)
	go func() {
		b := make([]byte, 32<<10)
		for {
			n, err := r.Read(b)
			if n > 0 || err != nil {
				got <- chunk{b[:n], err}
				return
			}
		}
	}()
	var first chunk
	var deadline <-chan time.Time // none for wait 0
	if wait > 0 {
		deadline = time.After(wait)
	}
	select {
	case first = <-got:
	case <-deadline:
		// The reader stays blocked until ccam exits; what comes later is
		// dropped, as Claude Code drops it.
		fmt.Fprintf(stderr, "ccam: no input on stdin in %s, so the command runs without it. Redirect stdin from /dev/null to skip the wait\n", wait)
		return []byte{}, nil // not nil: every attempt gets this empty input
	}
	if first.err == io.EOF {
		return first.b, nil
	}
	if first.err != nil {
		return nil, first.err
	}
	rest, err := io.ReadAll(io.LimitReader(r, int64(maxInput+1-len(first.b))))
	if err != nil {
		return nil, err
	}
	input := append(first.b, rest...)
	if len(input) > maxInput {
		return nil, fmt.Errorf("piped input is over 10 MiB, more than claude -p takes; name a file in the prompt instead")
	}
	return input, nil
}

func exitWith(code int) error {
	if code != 0 {
		return launcher.ExitCode(code)
	}
	return nil
}

// streamInput reports whether claude args take stream-json input, which
// Claude Code reads for the whole run, not only at the start.
func streamInput(args []string) bool {
	for i, a := range args {
		switch {
		case a == "--":
			return false
		case a == "--input-format=stream-json", a == "--input-format" && i+1 < len(args) && args[i+1] == "stream-json":
			return true
		}
	}
	return false
}

// printMode reports whether claude args ask for a non-interactive run.
func printMode(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if a == "-p" || a == "--print" {
			return true
		}
	}
	return false
}

// resumes reports whether claude args continue an earlier conversation,
// which only an account sharing its history can do.
func resumes(args []string) bool {
	return slices.ContainsFunc(args, func(a string) bool {
		return a == "-c" || a == "--continue" || a == "-r" || a == "--resume"
	})
}

// limited reports whether a failed run hit the account's usage limit: by
// Claude Code's message, or by the account's usage when it can be looked
// up, judged as p judges it (with --model, that model's week too).
func limited(e *env, cfg *config.Config, prof *config.Profile, p usage.Policy, out, errOut []byte) bool {
	if limitMessage(out) || limitMessage(errOut) {
		return true
	}
	for _, a := range usageAccounts(e, cfg) {
		if a.Name == prof.Name {
			r := e.fetcher().Fetch(context.Background(), a, usage.Options{NoCache: true})
			return p.Full(&r)
		}
	}
	return false
}

// limitMessage reports whether claude printed a usage-limit message
// (handoff.Classify): as a line, or in a line of JSON (--output-format json
// or stream-json) as the result or as a message of Claude Code's own. The
// lines and the result hold Claude's answer too (text and json output), so
// an answer line that reads like one counts as well; only stream-json
// messages other than Claude Code's own (<synthetic>) are left out. It is
// asked only about a run that failed.
func limitMessage(b []byte) bool {
	for _, line := range bytes.Split(b, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if handoff.Classify(string(line)) != "" {
			return true
		}
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var v struct {
			Result  string `json:"result"`
			Message struct {
				Model   string `json:"model"`
				Content []struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		}
		_ = json.Unmarshal(line, &v) // a field of another shape is left empty
		if handoff.Classify(v.Result) != "" {
			return true
		}
		if v.Message.Model == "<synthetic>" {
			for _, c := range v.Message.Content {
				if handoff.Classify(c.Text) != "" {
					return true
				}
			}
		}
	}
	return false
}

// fallbackTarget picks the next account: logged in, not tried yet, fitting
// from and every account tried (each reached its limit, which a second
// folder of its Claude account shares), sharing history when the args
// resume a conversation. p chooses among those with known usage (full ones
// never); accounts whose usage is unknown (offline, lookup refused) come
// after them in number order. It also returns the reports, for waiting.
func fallbackTarget(e *env, cfg *config.Config, from *config.Profile, tried map[string]bool, args []string, p usage.Policy) (*config.Profile, []*usage.Report) {
	return targetIn(e, cfg, newOrgs(e, cfg, false), from, tried, args, p)
}

// targetIn is fallbackTarget with the organization rule org.
func targetIn(e *env, cfg *config.Config, org *orgs, from *config.Profile, tried map[string]bool, args []string, p usage.Policy) (*config.Profile, []*usage.Report) {
	fits := func(prof *config.Profile) bool {
		for name := range tried {
			if t, ok := cfg.Find(name); ok && !org.fits(t, prof) {
				return false
			}
		}
		return org.fits(from, prof)
	}
	var profs []*config.Profile
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		if tried[prof.Name] || !creds.Has(profile.Dir(e.paths, prof), e.paths.Home) || !fits(prof) {
			continue
		}
		if resumes(args) {
			if h := from.Share.Get("history"); h == config.Local || prof.Share.Get("history") != h {
				continue
			}
		}
		profs = append(profs, prof)
	}
	if len(profs) == 0 {
		return nil, nil
	}
	reports := reportsOf(e, cfg, profs, usage.Options{})
	if best := p.Choose(reports, clock()); best != nil {
		prof, _ := cfg.Find(best.Name)
		return prof, reports
	}
	var unknown []*usage.Report
	for _, r := range reports {
		// A window known to be full (a limit Claude Code reported, say)
		// rules the account out even when its 5-hour use is unknown.
		if (r.Limits == nil || !r.Limits.FiveHour.Known) && !p.Full(r) {
			unknown = append(unknown, r)
		}
	}
	if len(unknown) == 0 {
		return nil, reports
	}
	slices.SortFunc(unknown, func(a, b *usage.Report) int { return cmp.Compare(a.Slot, b.Slot) })
	prof, _ := cfg.Find(unknown[0].Name)
	return prof, reports
}

// tail keeps the last max bytes written to it.
type tail struct {
	buf []byte
	max int
}

func (t *tail) Write(p []byte) (int, error) {
	t.buf = append(t.buf, p...)
	if len(t.buf) > t.max {
		t.buf = t.buf[len(t.buf)-t.max:]
	}
	return len(p), nil
}
