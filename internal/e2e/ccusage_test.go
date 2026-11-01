package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// The fake claude stands in for ccusage: it records CLAUDE_CONFIG_DIR and
// its arguments, and ends with FAKE_EXIT.
func TestCcusageRunsOverTheAccounts(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	for _, n := range []string{"personal", "work"} {
		if err := os.MkdirAll(filepath.Join(h.profileDir(n), "projects"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(fakeBin)
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Join(h.home, "bin"), 0o700)
	if err := os.WriteFile(filepath.Join(h.home, "bin", "ccusage"+exe), data, 0o755); err != nil {
		t.Fatal(err)
	}

	r := h.run("", []string{"FAKE_EXIT=3", "FAKE_STDOUT=the report"}, "ccam", "ccusage", "daily", "--json")
	if r.code != 3 || !strings.Contains(r.stdout, "the report") {
		t.Fatalf("ccusage's exit code or output is not ccam's: %+v", r)
	}
	dirs, args := h.called()
	if dirs != h.profileDir("personal")+","+h.profileDir("work") || !slices.Equal(args, []string{"daily", "--json"}) {
		t.Fatalf("CLAUDE_CONFIG_DIR=%s %q", dirs, args)
	}

	h.ok("ccam", "ccusage", "-a", "work", "monthly", "-a")
	dirs, args = h.called()
	if dirs != h.profileDir("work") || !slices.Equal(args, []string{"monthly", "-a"}) {
		t.Fatalf("CLAUDE_CONFIG_DIR=%s %q", dirs, args)
	}
}

// bun x and npx take the program they run (node_modules/.bin, a
// package.json bin) and npm its registry (.npmrc) from the folder they start
// in, so ccam starts them in its own folder, and names the settings ccusage
// reads in the working folder, and a --config relative to it, by their full
// path. The stand-in npx records where it started, its $PWD and its
// arguments and then, as npm does, runs the program the folder has under
// node_modules/.bin. The dry run's line, pasted into sh in the working
// folder, starts npx where the launch does. An installed ccusage still
// starts in the working folder, with the arguments as typed.
func TestCcusageRunnerStartsOutsideTheWorkingFolder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell scripts stand in for npx and ccusage")
	}
	h := newHarness(t)
	h.addAccounts()
	os.MkdirAll(filepath.Join(h.profileDir("work"), "projects"), 0o700)
	write := func(path, data string, mode os.FileMode) {
		t.Helper()
		os.MkdirAll(filepath.Dir(path), 0o700)
		if err := os.WriteFile(path, []byte(data), mode); err != nil {
			t.Fatal(err)
		}
	}
	repo, _ := filepath.EvalSymlinks(t.TempDir())
	planted := filepath.Join(repo, "planted-ran")
	write(filepath.Join(repo, "node_modules", ".bin", "ccusage@latest"), "#!/bin/sh\n: > '"+planted+"'\n", 0o755)
	write(filepath.Join(repo, ".npmrc"), "registry=http://127.0.0.1:9/\nyes=true\n", 0o600)
	settings := filepath.Join(repo, ".ccusage", "ccusage.json")
	write(settings, "{}\n", 0o600)
	runners, rec := filepath.Join(h.home, "runners"), filepath.Join(h.home, "runner-ran")
	record := "#!/bin/sh\n{ pwd -P; printf '%s\\n' \"$PWD\" \"$@\"; } > '" + rec + "'\n"
	write(filepath.Join(runners, "npx"), record+"if [ -x \"node_modules/.bin/$1\" ]; then exec \"node_modules/.bin/$1\"; fi\n", 0o755)
	// CCAM_HOME through a link: ccam names its folder through the link and
	// pwd -P names the folder itself, so a $PWD left naming the repository,
	// which sh would replace with the folder itself, shows.
	link := filepath.Join(t.TempDir(), "ccam-home")
	if err := os.Symlink(h.home, link); err != nil {
		t.Fatal(err)
	}
	home, _ := filepath.EvalSymlinks(h.home)
	state, named := filepath.Join(home, "state", "ccusage"), filepath.Join(link, "state", "ccusage")
	env := []string{"PATH=" + runners, "CCAM_HOME=" + link, "PWD=" + repo}

	// The dry run starts npx in ccam's folder with the settings named, and
	// makes no folder. Its line, pasted into sh in the repository before
	// there is one, stops at cd and runs nothing.
	r := h.runIn(repo, "", env, "ccam", "--dry-run", "ccusage", "daily")
	line := r.stdout
	if r.code != 0 || !strings.HasPrefix(line, "(cd "+named+" && ") || !strings.HasSuffix(line, " ccusage@latest daily --config "+settings+")\n") {
		t.Errorf("dry run: exit %d %q %q", r.code, line, r.stderr)
	}
	if _, err := os.Stat(state); err == nil {
		t.Error("the dry run made ccam's folder for the runners")
	}
	paste := func() error {
		c := exec.Command("sh", "-c", line)
		c.Dir, c.Env = repo, append(slices.Clone(h.env), env...)
		return c.Run()
	}
	if err := paste(); err == nil {
		t.Error("the dry run's line ran without ccam's folder")
	}
	if _, err := os.Stat(rec); err == nil {
		t.Error("the dry run's line started npx without ccam's folder")
	}
	for _, tc := range []struct{ args, want []string }{
		{[]string{"daily"}, []string{"ccusage@latest", "daily", "--config", settings}},
		{[]string{"daily", "--config", "my.json"}, []string{"ccusage@latest", "daily", "--config", filepath.Join(repo, "my.json")}},
		{[]string{"daily", "--config=/abs/c.json"}, []string{"ccusage@latest", "daily", "--config=/abs/c.json"}},
	} {
		os.Remove(rec)
		r := h.runIn(repo, "", env, append([]string{"ccam", "ccusage"}, tc.args...)...)
		got, _ := os.ReadFile(rec)
		want := strings.Join(append([]string{state, named}, tc.want...), "\n") + "\n"
		if r.code != 0 || string(got) != want {
			t.Errorf("%q: npx recorded (folder, $PWD, arguments)\n%s\nwant\n%s(exit %d %s)", tc.args, got, want, r.code, r.stderr)
		}
	}
	// With the folder there, the line starts npx in it, as the launch does.
	os.Remove(rec)
	err := paste()
	got, _ := os.ReadFile(rec)
	if want := strings.Join([]string{state, named, "ccusage@latest", "daily", "--config", settings}, "\n") + "\n"; err != nil || string(got) != want {
		t.Errorf("the dry run's line: npx recorded (folder, $PWD, arguments)\n%s\nwant\n%s(%v)", got, want, err)
	}
	if _, err := os.Stat(planted); err == nil {
		t.Error("the repository's node_modules/.bin/ccusage@latest ran")
	}
	if fi, err := os.Stat(state); err != nil || fi.Mode().Perm() != 0o700 {
		t.Errorf("ccam's folder for the runners: %v %v", fi, err)
	}

	write(filepath.Join(runners, "ccusage"), record, 0o755)
	os.Remove(rec)
	r = h.runIn(repo, "", env, "ccam", "ccusage", "daily")
	if got, _ := os.ReadFile(rec); r.code != 0 || string(got) != repo+"\n"+repo+"\ndaily\n" {
		t.Errorf("an installed ccusage recorded\n%s(exit %d %s)", got, r.code, r.stderr)
	}
}

// npm installs ccusage as ccusage.cmd, and npx is npx.cmd: ccam starts the
// batch file through cmd.exe with the arguments as typed, as it does an npm
// claude.cmd.
func TestBatchCcusageGetsTheArgumentsAsTyped(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("batch files are Windows'")
	}
	h := newHarness(t)
	h.addAccounts()
	os.MkdirAll(filepath.Join(h.profileDir("work"), "projects"), 0o700)
	shim := filepath.Join(h.home, "bin", "ccusage.cmd")
	os.MkdirAll(filepath.Dir(shim), 0o700)
	// The last line is the one npm's cmd-shim writes, with the fake claude
	// for node and its script.
	os.WriteFile(shim, []byte("@ECHO off\r\nSETLOCAL\r\nSET \"_prog="+fakeBin+"\"\r\n"+
		"endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & \"%_prog%\" %*\r\n"), 0o600)
	proj := t.TempDir()
	args := []string{"daily", `a & echo INJECTED> injected.txt`, `100% %PATH% !OS!`, `a"b`, `C:\dir with space\`, "é"}
	r := h.runIn(proj, "", []string{"FAKE_ARGV=ucrt"}, append([]string{"ccam", "ccusage"}, args...)...)
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	dir, got := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(got, args) {
		t.Fatalf("%s got\n%q\nwant\n%q", dir, got, args)
	}
	if _, err := os.Stat(filepath.Join(proj, "injected.txt")); err == nil {
		t.Fatal("cmd.exe ran something else")
	}
	// A line break would end the command in cmd.exe: refused, ccusage does
	// not run.
	r = h.runIn(proj, "", nil, "ccam", "ccusage", "daily", "two\nlines & echo INJECTED> injected.txt")
	if r.code == 0 || !strings.Contains(r.stderr, "line break") {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("ccusage ran with a line break")
	}
	// The dry run says the same, not the advice for a claude.cmd.
	r = h.runIn(proj, "", nil, "ccam", "--dry-run", "ccusage", "daily", "two\nlines")
	if r.code == 0 || !strings.Contains(r.stderr, "ccusage.cmd is a batch file") || strings.Contains(r.stderr, "claude") {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
}

// On Windows npx is npx.cmd, which ccam starts through cmd.exe as a child
// process: the child starts in ccam's own folder too. The stand-in npx.cmd
// records the folder it starts in, runs, as npm does, the folder's
// node_modules\.bin\ccusage@latest.cmd if there is one, and then the fake
// claude, which records the arguments.
func TestBatchNpxStartsOutsideTheWorkingFolder(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("batch files are Windows'")
	}
	h := newHarness(t)
	h.addAccounts()
	os.MkdirAll(filepath.Join(h.profileDir("work"), "projects"), 0o700)
	repo := t.TempDir()
	planted := filepath.Join(repo, "planted-ran")
	os.MkdirAll(filepath.Join(repo, "node_modules", ".bin"), 0o700)
	os.WriteFile(filepath.Join(repo, "node_modules", ".bin", "ccusage@latest.cmd"), []byte("@type nul > \""+planted+"\"\r\n"), 0o600)
	settings := filepath.Join(repo, ".ccusage", "ccusage.json")
	os.MkdirAll(filepath.Dir(settings), 0o700)
	os.WriteFile(settings, []byte("{}\n"), 0o600)
	rec := filepath.Join(h.home, "npx-folder")
	npx := filepath.Join(h.home, "bin", "npx.cmd")
	os.MkdirAll(filepath.Dir(npx), 0o700)
	os.WriteFile(npx, []byte("@echo off\r\ncd > \""+rec+"\"\r\n"+
		"if exist \"node_modules\\.bin\\%1.cmd\" call \"node_modules\\.bin\\%1.cmd\"\r\n"+
		"\""+fakeBin+"\" %*\r\n"), 0o600)
	same := func(a, b string) bool {
		fa, err1 := os.Stat(a)
		fb, err2 := os.Stat(b)
		return err1 == nil && err2 == nil && os.SameFile(fa, fb)
	}

	r := h.runIn(repo, "", nil, "ccam", "ccusage", "daily")
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	if _, args := h.called(); len(args) != 4 || !slices.Equal(args[:3], []string{"ccusage@latest", "daily", "--config"}) || !same(args[3], settings) {
		t.Errorf("npx.cmd got %q, want ccusage@latest daily --config %s", args, settings)
	}
	data, _ := os.ReadFile(rec)
	if folder := strings.TrimSpace(string(data)); !same(folder, filepath.Join(h.home, "state", "ccusage")) {
		t.Errorf("npx.cmd started in %q, not in ccam's folder", folder)
	}
	if _, err := os.Stat(planted); err == nil {
		t.Error(`the folder's node_modules\.bin\ccusage@latest.cmd ran`)
	}
}
