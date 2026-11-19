package cli

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
)

// agentEntry is one session `claude agents --json` lists (Claude Code
// 2.1.304); background sessions have an id and a state.
type agentEntry struct {
	Account    string `json:"account"`
	ID         string `json:"id,omitempty"`
	Kind       string `json:"kind"`
	Name       string `json:"name,omitempty"`
	State      string `json:"state,omitempty"`
	Status     string `json:"status,omitempty"`
	WaitingFor string `json:"waitingFor,omitempty"`
	Cwd        string `json:"cwd"`
	StartedAt  int64  `json:"startedAt"`
	SessionID  string `json:"sessionId,omitempty"`
	PID        int    `json:"pid,omitempty"`
}

// agentsTimeout bounds one account's `claude agents --json` or `claude stop`.
const agentsTimeout = 15 * time.Second

// agentsWaitDelay is how long asAccount waits, once claude has ended, for a
// process it left behind to close claude's output.
var agentsWaitDelay = 2 * time.Second

func agentsCmd(e *env) *cobra.Command {
	var all, asJSON bool
	cmd := &cobra.Command{
		Use:   "agents [--all] [--json]",
		Short: "List the Claude Code sessions of every account: background agents and open sessions",
		Long: "Runs `claude agents --json` in every account at once and lists the sessions\n" +
			"together, with the account each belongs to. Background sessions belong to the\n" +
			"account that started them; open one with `ccam attach <id>`.\n\n" +
			"The dashboard (ccam on a terminal) shows the same list on its Agents tab, where\n" +
			"enter attaches to a background session and x stops one. To run a task as another\n" +
			"account, see `ccam agent`.",
		Example: "  ccam agents\n  ccam agents --all     # also finished background sessions\n  ccam attach 49ec0028",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			list, failed := listAgents(e, cfg, all)
			if asJSON {
				if list == nil {
					list = []agentEntry{}
				}
				data, _ := json.MarshalIndent(list, "", "  ")
				e.printf("%s\n", data)
			} else if len(list) > 0 || len(failed) == 0 {
				printAgents(e, list)
			}
			for _, f := range failed {
				fmt.Fprintln(e.stderr, "ccam:", f)
			}
			if len(failed) > 0 {
				// A script must not take a failure for an empty list.
				return fmt.Errorf("%d of %d accounts could not be asked", len(failed), max(len(cfg.Profiles), 1))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "also list finished background sessions")
	cmd.Flags().BoolVar(&asJSON, "json", false, "print as JSON")
	return cmd
}

// listAgents asks every account in parallel; an account that fails is
// named in failed and left out.
func listAgents(e *env, cfg *config.Config, all bool) (list []agentEntry, failed []string) {
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return nil, []string{err.Error()}
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := accountAgents(e, bin, p, all)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failed = append(failed, fmt.Sprintf("%s: could not list its sessions: %v", p.Name, err))
				return
			}
			list = append(list, got...)
		}()
	}
	wg.Wait()
	slices.SortFunc(list, func(a, b agentEntry) int { return cmp.Compare(b.StartedAt, a.StartedAt) })
	slices.Sort(failed)
	return list, failed
}

func accountAgents(e *env, bin string, p *config.Profile, all bool) ([]agentEntry, error) {
	args := []string{"agents", "--json"}
	if all {
		args = append(args, "--all")
	}
	out, err := asAccount(e, bin, p, args...)
	if err != nil {
		return nil, err
	}
	var list []agentEntry
	if err := json.Unmarshal(out, &list); err != nil {
		return nil, fmt.Errorf("unexpected output of claude agents --json: %w", err)
	}
	for i := range list {
		list[i].Account = p.Name
	}
	return list, nil
}

// asAccount runs `claude args…` as account p within agentsTimeout and
// returns what it printed; an error carries what claude wrote on stderr.
func asAccount(e *env, bin string, p *config.Profile, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), agentsTimeout)
	defer cancel()
	c, err := launcher.ClaudeCommand(ctx, bin, args...)
	if err != nil {
		return nil, err
	}
	c.Env = launcher.ProfileEnv(e.paths, os.Environ(), profile.Dir(e.paths, p), p.Name)
	var stderr bytes.Buffer
	c.Stderr = &stderr
	// Without WaitDelay, Output would wait for a process claude leaves
	// behind, and the dashboard with it.
	c.WaitDelay = agentsWaitDelay
	out, err := c.Output()
	if errors.Is(err, exec.ErrWaitDelay) {
		err = nil // claude itself ended well
	}
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// stateText is the State column of `ccam agents` and the dashboard's Agents
// tab: a background session's state and what its process does ("blocked,
// idle"), or "open (busy)" for an interactive one, then what it waits for.
func (a agentEntry) stateText() string {
	state := cmp.Or(a.State, a.Status, "-")
	if a.Kind == "interactive" {
		state = "open (" + cmp.Or(a.Status, "interactive") + ")"
	} else if a.State != "" && a.Status != "" && a.Status != a.State {
		state += ", " + a.Status
	}
	if a.WaitingFor != "" {
		state += ": " + a.WaitingFor
	}
	return state
}

// tildeDir writes a folder under home as ~/….
func tildeDir(dir, home string) string {
	if rel, err := filepath.Rel(home, dir); err == nil && home != "" && !strings.HasPrefix(rel, "..") {
		return filepath.Join("~", rel)
	}
	return dir
}

func printAgents(e *env, list []agentEntry) {
	if len(list) == 0 {
		e.printf("No sessions. Start one in the background with `ccam <account> --bg \"…\"`.\n")
		return
	}
	home, _ := os.UserHomeDir()
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Account\tID\tName\tState\tFolder\tStarted")
	for _, a := range list {
		started := "-"
		if a.StartedAt > 0 {
			started = Ago(time.Since(time.UnixMilli(a.StartedAt)))
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", a.Account, cmp.Or(a.ID, "-"), cmp.Or(a.Name, "-"), a.stateText(), tildeDir(a.Cwd, home), started)
	}
	_ = tw.Flush()
	e.printf("Open a background session with `ccam attach <id>`.\n")
}

// stopSession stops account's background session id as `ccam <account>
// stop <id>` does; Claude Code keeps its conversation.
func stopSession(e *env, account, id string) error {
	cfg, p, err := e.loadProfile(account)
	if err != nil {
		return err
	}
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return err
	}
	if _, err := asAccount(e, bin, p, "stop", id); err != nil {
		return fmt.Errorf("could not stop %s: %w", id, err)
	}
	return nil
}

func attachCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach <id|name> [claude args]",
		Short: "Open a background session of whichever account started it",
		Long: "Finds the account whose folder has the background session (jobs/<id>) and runs\n" +
			"`claude attach <id>` as that account. A unique start of the id is enough, and so\n" +
			"is part of the session name. Everything after the id goes to claude unchanged.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			hits := findJobs(e, cfg, args[0])
			switch len(hits) {
			case 0:
				return fmt.Errorf("no account has a background session %q (see `ccam agents --all`)", args[0])
			case 1:
			default:
				return fmt.Errorf("%q fits more than one session: %s; give more of the id or the name", args[0], strings.Join(hits, ", "))
			}
			account, id, _ := strings.Cut(hits[0], "/")
			fmt.Fprintf(e.stderr, "→ %s (%s)\n", id, account)
			return launcher.LaunchWith(e.paths, cfg, account, append([]string{"attach", id}, args[1:]...), e.dryRun, e.stdout)
		},
	}
	// Flags after the id are claude's, as for `ccam <name>`.
	cmd.Flags().SetInterspersed(false)
	return cmd
}

// findJobs lists account/id of the background sessions sel names: by the
// start of the id, else, as `claude attach` does, by part of the session
// name (any case). Like Claude Code it never reads what may start an id
// (up to 8 of 0-9 and a-f) as a name.
func findJobs(e *env, cfg *config.Config, sel string) []string {
	var byID, byName []string
	name := strings.ToLower(strings.TrimSpace(sel))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		jobs := filepath.Join(profile.Dir(e.paths, p), "jobs")
		entries, _ := os.ReadDir(jobs)
		for _, en := range entries {
			if !en.IsDir() || sel == "" {
				continue
			}
			if strings.HasPrefix(en.Name(), sel) {
				byID = append(byID, p.Name+"/"+en.Name())
				continue
			}
			if name == "" || idStart(name) {
				continue
			}
			var st struct {
				Name string `json:"name"`
			}
			if b, err := os.ReadFile(filepath.Join(jobs, en.Name(), "state.json")); err == nil &&
				json.Unmarshal(b, &st) == nil && strings.Contains(strings.ToLower(st.Name), name) {
				byName = append(byName, p.Name+"/"+en.Name())
			}
		}
	}
	if len(byID) > 0 {
		return byID
	}
	return byName
}

// idStart reports whether s may be the start of a background session id.
func idStart(s string) bool {
	if len(s) == 0 || len(s) > 8 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if c := s[i]; !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
