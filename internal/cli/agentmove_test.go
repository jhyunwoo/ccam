package cli

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// An account folder loads the agents of another the same way when its
// agents and plugins folders and its settings.json are the very ones of the
// other, through ccam's group links, or both lack them.
func TestSameAgents(t *testing.T) {
	root := t.TempDir()
	group, empty := filepath.Join(root, "group"), filepath.Join(root, "empty")
	for _, dir := range []string{filepath.Join(group, "agents"), filepath.Join(group, "plugins"), empty} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(group, "settings.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	// account makes an account folder: spec lists its entries, "name" for
	// one of its own and "name>group" for a link into group (a group folder
	// for "empty", which has none of them yet).
	n := 0
	account := func(spec string) string {
		t.Helper()
		n++
		dir := filepath.Join(root, "account"+strconv.Itoa(n))
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		for _, entry := range strings.Fields(spec) {
			name, into, linked := strings.Cut(entry, ">")
			path, isDir := filepath.Join(dir, name), name != "settings.json"
			var err error
			switch {
			case linked:
				err = fsutil.Link(filepath.Join(root, into, name), path, isDir)
			case isDir:
				err = os.MkdirAll(path, 0o700)
			default:
				err = os.WriteFile(path, []byte("{}"), 0o600)
			}
			if err != nil {
				t.Fatal(err)
			}
		}
		return dir
	}
	folders := "agents>group plugins>group"
	for _, c := range []struct {
		name, from, to string
		files          bool // links a file, which Windows does only in Developer Mode
		want           bool
	}{
		{"both missing", "", "", false, true},
		{"shared folders", folders, folders, false, true},
		{"one missing", folders, "plugins>group", false, false},
		{"folders of their own", "agents plugins", "agents plugins", false, false},
		{"one of their own", folders, "agents plugins>group", false, false},
		{"settings of its own", folders + " settings.json", folders, false, false},
		{"shared links", folders + " settings.json>group", folders + " settings.json>group", true, true},
		{"a shared settings.json not made yet", folders + " settings.json>empty", folders + " settings.json>empty", true, true},
		{"shared settings and settings of its own", folders + " settings.json>group", folders + " settings.json", true, false},
	} {
		if c.files {
			if err := fsutil.CanLinkFiles(root); err != nil {
				t.Logf("%s: %v", c.name, err)
				continue
			}
		}
		from, to := account(c.from), account(c.to)
		if got := sameAgents(from, to); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
		if got := sameAgents(to, from); got != c.want {
			t.Errorf("%s, the other way: %v, want %v", c.name, got, c.want)
		}
	}
	if dir := account(folders + " settings.json"); !sameAgents(dir, dir) {
		t.Error("an account folder and itself")
	}
}

// When its arguments named no agent, a session ran its conversation as the
// agent setting says, the project's local settings over its shared ones
// over the account's; else as the conversation's newest agent-setting row
// says, which a session that resumed it by id goes on writing whatever the
// setting.
func TestConversationAgent(t *testing.T) {
	row := `{"type":"agent-setting","agentSetting":"writer","sessionId":"s1"}` + "\n"
	for _, c := range []struct {
		name               string
		account, own, team string // settings.json of the account, and the project's local and shared one
		rows               string
		want               string
	}{
		{"nothing", "", "", "", "", ""},
		{"the conversation's", "", "", "", row, "writer"},
		{"the account's setting", `{"agent":"reviewer"}`, "", "", row, "reviewer"},
		{"the project's over the account's", `{"agent":"reviewer"}`, "", `{"agent":"tester"}`, row, "tester"},
		{"the project's local over its shared", `{"agent":"reviewer"}`, `{"agent":"local"}`, `{"agent":"tester"}`, row, "local"},
		{"settings without one", `{"model":"opus"}`, `{}`, `{"agent":""}`, row, "writer"},
		{"a setting of another type", `{"agent":["reviewer"]}`, "", "", row, "writer"},
	} {
		root := t.TempDir()
		dir, project, conv := filepath.Join(root, "account"), filepath.Join(root, "app"), filepath.Join(root, "s1.jsonl")
		for path, body := range map[string]string{
			filepath.Join(dir, "settings.json"):                      c.account,
			filepath.Join(project, ".claude", "settings.local.json"): c.own,
			filepath.Join(project, ".claude", "settings.json"):       c.team,
			conv: c.rows,
		} {
			if body == "" {
				continue
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if got := conversationAgent(dir, project, conv); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}
