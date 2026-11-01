package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// exeSuffix is the file name suffix of programs.
var exeSuffix = map[bool]string{true: ".exe"}[runtime.GOOS == "windows"]

// fakeRunners makes PATH a folder of empty stand-ins for the programs named,
// so that a dry run shows the one ccusage would start with, and returns the
// folder.
func fakeRunners(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n+exeSuffix), nil, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir)
	return dir
}

// withConversations gives the accounts the projects folder Claude Code makes
// with their first conversation.
func (c *testCLI) withConversations(names ...string) {
	c.t.Helper()
	for _, n := range names {
		if err := os.MkdirAll(filepath.Join(c.p.ProfileDir(n), "projects"), 0o700); err != nil {
			c.t.Fatal(err)
		}
	}
}

// configDirOf is the CLAUDE_CONFIG_DIR a dry run sets: VAR=value for sh,
// $env:VAR='value' for PowerShell.
func configDirOf(out string) string {
	_, v, _ := strings.Cut(out, "CLAUDE_CONFIG_DIR=")
	if v, ok := strings.CutPrefix(v, "'"); ok {
		v, _, _ = strings.Cut(v, "'")
		return v
	}
	v, _, _ = strings.Cut(v, " ")
	return v
}

func TestCcusageReadsEveryAccountWithConversations(t *testing.T) {
	c := newTestCLI(t, nil)
	bin := filepath.Join(fakeRunners(t, "ccusage", "bun", "npx"), "ccusage"+exeSuffix)
	c.withConversations("side", "work")
	t.Setenv("CLAUDE_CONFIG_DIR", filepath.Join(t.TempDir(), "inherited"))
	out, errOut, code := c.run("", "--dry-run", "ccusage", "daily", "--json")
	dirs := c.p.ProfileDir("work") + "," + c.p.ProfileDir("side") // the config's order, not the names'
	// personal has no conversations; the inherited folder is replaced.
	if code != 0 || configDirOf(out) != dirs || !strings.Contains(out, bin) || !strings.Contains(out, "daily --json") ||
		strings.Contains(out, "inherited") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
}

func TestCcusageAccounts(t *testing.T) {
	c := newTestCLI(t, nil)
	fakeRunners(t, "ccusage")
	c.withConversations("personal", "work", "side")
	dir := c.p.ProfileDir
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"-a", "work", "monthly"}, dir("work")},
		{[]string{"-a", "3,work"}, dir("side") + "," + dir("work")},
		{[]string{"--account=personal", "--account", " side , ", "--dry-run"}, dir("personal") + "," + dir("side")},
		{[]string{"-a", "work,2", "-a", "work"}, dir("work")},
		{[]string{"-a=side"}, dir("side")},
	} {
		out, errOut, code := c.run("", append([]string{"--dry-run", "ccusage"}, tc.args...)...)
		if code != 0 || configDirOf(out) != tc.want {
			t.Fatalf("%q: %d %q %q", tc.args, code, out, errOut)
		}
	}
	// What follows the options goes to ccusage, -a included.
	out, _, _ := c.run("", "--dry-run", "ccusage", "-a", "work", "session", "-a")
	if !strings.Contains(out, "session -a") {
		t.Fatal(out)
	}
	// --dry-run among ccam's options is ccam's, wherever it is among them.
	out, errOut, code := c.run("", "ccusage", "-a", "work", "--dry-run", "daily")
	if code != 0 || configDirOf(out) != dir("work") || !strings.Contains(out, "daily") || strings.Contains(out, "dry-run") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
}

// A group that shares history links one projects folder into its members:
// ccusage reads it through each, and -a names whom else it brings.
func TestCcusageSharedHistory(t *testing.T) {
	c := newTestCLI(t, nil)
	fakeRunners(t, "ccusage")
	shared := filepath.Join(c.p.GroupDir("team"), "projects")
	if err := os.MkdirAll(shared, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"personal", "work", "side"} {
		if err := fsutil.Link(shared, filepath.Join(c.p.ProfileDir(n), "projects"), true); err != nil {
			t.Fatal(err)
		}
	}
	dir := c.p.ProfileDir
	out, errOut, code := c.run("", "--dry-run", "ccusage")
	if code != 0 || configDirOf(out) != dir("personal")+","+dir("work")+","+dir("side") || errOut != "" {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
	out, errOut, code = c.run("", "--dry-run", "ccusage", "-a", "work")
	if code != 0 || configDirOf(out) != dir("work") ||
		errOut != "ccam: work shares its conversation history with personal and side, so the report counts their conversations too\n" {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
	// One note for the one folder, however many chosen accounts share it.
	if _, errOut, code = c.run("", "--dry-run", "ccusage", "-a", "work,personal"); code != 0 ||
		errOut != "ccam: work shares its conversation history with side, so the report counts their conversations too\n" {
		t.Fatalf("%d %q", code, errOut)
	}
	if _, errOut, code = c.run("", "--dry-run", "ccusage", "-a", "1,2,3"); code != 0 || errOut != "" {
		t.Fatalf("%d %q", code, errOut)
	}
	// A file named projects is no conversations.
	os.Remove(filepath.Join(dir("side"), "projects"))
	os.WriteFile(filepath.Join(dir("side"), "projects"), nil, 0o600)
	if out, errOut, code = c.run("", "--dry-run", "ccusage"); code != 0 || configDirOf(out) != dir("personal")+","+dir("work") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
}

// Under WSL, a runner on a Windows drive would run as a Windows program,
// without the account folders: ccam passes over it, wherever it is on PATH
// and however the drive is mounted.
func TestCcusageSkipsWindowsRunnersUnderWSL(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("WSL is Linux")
	}
	c := newTestCLI(t, nil)
	c.withConversations("work")
	win, virtio, lin := t.TempDir(), t.TempDir(), t.TempDir()
	for _, f := range []string{filepath.Join(win, "ccusage"), filepath.Join(win, "npx"), filepath.Join(virtio, "npx"), filepath.Join(lin, "npx")} {
		if err := os.WriteFile(f, nil, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(t.TempDir(), "winnode") // a PATH folder that links into a drive
	if err := os.Symlink(win, link); err != nil {
		t.Fatal(err)
	}
	mounts := filepath.Join(t.TempDir(), "mounts")
	os.WriteFile(mounts, []byte("/dev/sdc / ext4 rw,relatime 0 0\n"+
		"C:\\134 "+win+" 9p rw,noatime,aname=drvfs;path=C:\\;uid=1000 0 0\n"+
		"drvfsD "+virtio+" virtiofs rw,noatime 0 0\n"), 0o600)
	orig := mountsFile
	mountsFile = mounts
	t.Cleanup(func() { mountsFile = orig })
	sep := string(os.PathListSeparator)
	for _, tc := range []struct{ path, want string }{
		{win + sep + lin, filepath.Join(lin, "npx")},    // a Linux npx after the Windows ccusage and npx
		{virtio + sep + lin, filepath.Join(lin, "npx")}, // a drive mounted with virtiofs
		{link + sep + lin, filepath.Join(lin, "npx")},   // reached through a link
	} {
		t.Setenv("PATH", tc.path)
		out, errOut, code := c.run("", "--dry-run", "ccusage", "daily")
		if code != 0 || !strings.Contains(out, tc.want+" ccusage@latest daily") {
			t.Fatalf("PATH=%s: %d %q %q", tc.path, code, out, errOut)
		}
	}
	// Only Windows programs: ccam names the first instead of running it.
	for _, tc := range []struct{ path, want string }{
		{win, filepath.Join(win, "ccusage")},
		{virtio, filepath.Join(virtio, "npx")},
		{link, filepath.Join(link, "ccusage")},
	} {
		t.Setenv("PATH", tc.path)
		if out, errOut, code := c.run("", "--dry-run", "ccusage", "daily"); code == 0 ||
			!strings.Contains(errOut, tc.want+" is a Windows program") || !strings.Contains(errOut, "inside WSL") {
			t.Fatalf("PATH=%s: %d %q %q", tc.path, code, out, errOut)
		}
	}
}

func TestCcusageAccountErrors(t *testing.T) {
	c := newTestCLI(t, nil)
	fakeRunners(t, "ccusage")
	c.withConversations("personal")
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"-a"}, "-a needs account names"},
		{[]string{"--account", ","}, "--account needs account names"},
		{[]string{"--account="}, "--account needs account names"},
		{[]string{"-a", "nope", "daily"}, `no account "nope"`},
		{[]string{"-a", "work"}, "work has no Claude Code conversations yet"},
		{[]string{"-a", "work,side"}, "work and side have no Claude Code conversations yet"},
	} {
		out, errOut, code := c.run("", append([]string{"--dry-run", "ccusage"}, tc.args...)...)
		if code == 0 || !strings.Contains(errOut, tc.want) || strings.Contains(out, "CLAUDE_CONFIG_DIR") {
			t.Fatalf("%q: %d %q %q", tc.args, code, out, errOut)
		}
	}
	if out, errOut, code := c.run("", "--dry-run", "ccusage", "-a", "work,personal"); code != 0 || configDirOf(out) != c.p.ProfileDir("personal") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
	os.RemoveAll(filepath.Join(c.p.ProfileDir("personal"), "projects"))
	if _, errOut, code := c.run("", "--dry-run", "ccusage"); code == 0 || !strings.Contains(errOut, "no account has Claude Code conversations yet") {
		t.Fatalf("%d %q", code, errOut)
	}
}

func TestCcusageWithoutAccounts(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	fakeRunners(t, "ccusage")
	var out, errOut strings.Builder
	if code := ExecuteWith([]string{"--dry-run", "ccusage"}, paths.FromEnv(), strings.NewReader(""), &out, &errOut); code == 0 ||
		!strings.Contains(errOut.String(), "no accounts yet") {
		t.Fatalf("%d %q %q", code, out.String(), errOut.String())
	}
}

// A comma would split the folder in two for ccusage.
func TestCcusageFolderWithAComma(t *testing.T) {
	c := newTestCLI(t, nil)
	fakeRunners(t, "ccusage")
	c.withConversations("personal")
	odd := filepath.Join(t.TempDir(), "a,b")
	os.MkdirAll(odd, 0o700)
	cfg := c.cfg()
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: "odd", Slot: 4, Dir: odd})
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
	// Without conversations the folder is not given to ccusage at all.
	if out, errOut, code := c.run("", "--dry-run", "ccusage"); code != 0 || configDirOf(out) != c.p.ProfileDir("personal") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
	os.MkdirAll(filepath.Join(odd, "projects"), 0o700)
	if _, errOut, code := c.run("", "--dry-run", "ccusage"); code == 0 || !strings.Contains(errOut, "the folder of odd") || !strings.Contains(errOut, "comma") {
		t.Fatalf("%d %q", code, errOut)
	}
	if out, errOut, code := c.run("", "--dry-run", "ccusage", "-a", "personal"); code != 0 || configDirOf(out) != c.p.ProfileDir("personal") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
}

// An installed ccusage first, then bun, then npm, as ccusage's guide runs it.
func TestCcusageRunner(t *testing.T) {
	c := newTestCLI(t, nil)
	c.withConversations("work")
	for _, tc := range []struct {
		have []string
		bin  string
		args string
	}{
		{[]string{"bun", "npx"}, "bun", "x ccusage daily"},
		{[]string{"npx"}, "npx", "ccusage@latest daily"},
	} {
		dir := fakeRunners(t, tc.have...)
		out, errOut, code := c.run("", "--dry-run", "ccusage", "daily")
		if code != 0 || !strings.Contains(out, filepath.Join(dir, tc.bin+exeSuffix)) || !strings.Contains(out, tc.args) {
			t.Fatalf("%q: %d %q %q", tc.have, code, out, errOut)
		}
	}
	fakeRunners(t)
	if out, errOut, code := c.run("", "--dry-run", "ccusage", "daily"); code == 0 || !strings.Contains(errOut, "install Node.js or Bun, or ccusage itself") {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
}

// bun x and npx start in ccam's own folder, so the settings ccusage would
// read in the working folder, and a --config relative to it, are named by
// their full path. A --config without a value, and what follows --, are
// left for ccusage.
func TestCcusageConfigArgs(t *testing.T) {
	t.Chdir(t.TempDir())
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	mine, team := filepath.Join(here, "my.json"), filepath.Join(here, "team", "c.json")
	check := func(args, want []string) {
		t.Helper()
		in := slices.Clone(args)
		if got, err := ccusageConfigArgs(args); err != nil || !slices.Equal(got, want) || !slices.Equal(args, in) {
			t.Errorf("%q: got %q (%v), want %q", in, got, err, want)
		}
	}
	// No settings here: only a relative --config changes.
	check(nil, nil)
	check([]string{"daily", "--json"}, []string{"daily", "--json"})
	check([]string{"daily", "--config", "my.json"}, []string{"daily", "--config", mine})
	check([]string{"--config=my.json", "daily"}, []string{"--config=" + mine, "daily"})
	check([]string{"daily", "--config", team, "--config=" + team}, []string{"daily", "--config", team, "--config=" + team})
	check([]string{"daily", "--config"}, []string{"daily", "--config"})
	check([]string{"daily", "--config="}, []string{"daily", "--config="})
	check([]string{"daily", "--config", "--json"}, []string{"daily", "--config", "--json"})
	check([]string{"daily", "--config", "--", "my.json"}, []string{"daily", "--config", "--", "my.json"})
	check([]string{"daily", "--", "--config", "my.json"}, []string{"daily", "--", "--config", "my.json"})

	settings := filepath.Join(here, ".ccusage", "ccusage.json")
	os.MkdirAll(filepath.Dir(settings), 0o700)
	if err := os.WriteFile(settings, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	check(nil, []string{"--config", settings})
	check([]string{"daily", "--json"}, []string{"daily", "--json", "--config", settings})
	check([]string{"daily", "--", "--config", "my.json"}, []string{"daily", "--config", settings, "--", "--config", "my.json"})
	// A --config, even one without a value, is the person's choice.
	check([]string{"daily", "--config", "my.json"}, []string{"daily", "--config", mine})
	check([]string{"daily", "--config=" + team}, []string{"daily", "--config=" + team})
	check([]string{"daily", "--config"}, []string{"daily", "--config"})
}

// In a folder reached through a link, the .. of a relative --config leads
// out of the folder the link names, as the system takes it, not back to the
// link's own folder.
func TestCcusageConfigArgsThroughALink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows takes .. by name")
	}
	t.Chdir(t.TempDir())
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(here, "real", "deep"), 0o700); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(here, "real", "x.json"), []byte("on disk"), 0o600)
	os.WriteFile(filepath.Join(here, "x.json"), []byte("by name"), 0o600)
	if err := os.Symlink(filepath.Join(here, "real", "deep"), filepath.Join(here, "link")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(filepath.Join(here, "link"))
	got, err := ccusageConfigArgs([]string{"--config", "../x.json"})
	if err != nil || len(got) != 2 || !filepath.IsAbs(got[1]) {
		t.Fatalf("%q %v", got, err)
	}
	if data, err := os.ReadFile(got[1]); err != nil || string(data) != "on disk" {
		t.Errorf("--config ../x.json in the linked folder became %s, which holds %q (%v)", got[1], data, err)
	}
}

// The dry run shows bun x and npx starting in ccam's own folder, and the
// full path they get for the settings in the working folder, and makes no
// folder for them; an installed ccusage starts in the working folder with
// the arguments as typed.
func TestCcusageDryRunNamesTheSettingsHere(t *testing.T) {
	c := newTestCLI(t, nil)
	c.withConversations("work")
	t.Chdir(t.TempDir())
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	settings := filepath.Join(here, ".ccusage", "ccusage.json")
	os.MkdirAll(filepath.Dir(settings), 0o700)
	if err := os.WriteFile(settings, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(c.p.State, "ccusage")
	startsThere := func(out string) bool {
		if runtime.GOOS == "windows" {
			return strings.Contains(out, " -WorkingDirectory "+launcher.PSString(dir)+" ")
		}
		return strings.HasPrefix(out, "(cd "+launcher.ShellQuote(dir)+" && ")
	}
	for _, tc := range []struct {
		have []string
		args string
	}{
		{[]string{"bun", "npx"}, "x ccusage daily --config " + settings},
		{[]string{"npx"}, "ccusage@latest daily --config " + settings},
	} {
		fakeRunners(t, tc.have...)
		if out, errOut, code := c.run("", "--dry-run", "ccusage", "daily"); code != 0 || !strings.Contains(out, tc.args) || !startsThere(out) {
			t.Errorf("%q: %d %q %q", tc.have, code, out, errOut)
		}
	}
	if _, err := os.Stat(dir); err == nil {
		t.Error("the dry run made ccam's folder for the runners")
	}
	fakeRunners(t, "ccusage", "npx")
	if out, errOut, code := c.run("", "--dry-run", "ccusage", "daily"); code != 0 || !strings.Contains(out, "daily") || strings.Contains(out, "--config") ||
		strings.Contains(out, dir) {
		t.Errorf("%d %q %q", code, out, errOut)
	}
}

func TestCcusageHelp(t *testing.T) {
	c := newTestCLI(t, nil)
	fakeRunners(t, "ccusage")
	c.withConversations("work")
	for _, help := range []string{"-h", "--help"} {
		out, errOut, code := c.run("", "ccusage", help)
		if code != 0 || !strings.Contains(out, "ccam ccusage [-a name[,name…]] [ccusage args]") || strings.Contains(out, "CLAUDE_CONFIG_DIR=") {
			t.Fatalf("%s: %d %q %q", help, code, out, errOut)
		}
	}
	// After -- (which ccam drops) or after ccusage's own command, --help is ccusage's.
	for _, args := range [][]string{{"--", "--help"}, {"daily", "--help"}} {
		out, errOut, code := c.run("", append([]string{"--dry-run", "ccusage"}, args...)...)
		want := strings.TrimPrefix(strings.Join(args, " "), "-- ")
		if code != 0 || configDirOf(out) != c.p.ProfileDir("work") || !strings.Contains(out, want) || strings.Contains(out, "-- --help") {
			t.Fatalf("%q: %d %q %q", args, code, out, errOut)
		}
	}
}
