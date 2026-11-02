// Package transcript finds Claude Code's saved conversations in a config
// folder: <config>/projects/<slug of the working folder>/<session id>.jsonl.
// The entries inside are internal to Claude Code and change between
// versions, so only the few fields used here are read, and anything else
// is skipped.
package transcript

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

// maxSlug is where Claude Code (2.1.304) cuts long folder names.
const maxSlug = 200

// Slug is the folder name Claude Code uses under projects/ for the working
// folder dir (with symbolic links resolved): every UTF-16 unit that is not
// an ASCII letter or digit becomes '-'; a name over 200 characters keeps
// the first 200 and a hash of the whole path.
func Slug(dir string) string {
	units := utf16.Encode([]rune(dir))
	var b strings.Builder
	for _, u := range units {
		if u < 0x80 && (u >= 'a' && u <= 'z' || u >= 'A' && u <= 'Z' || u >= '0' && u <= '9') {
			b.WriteByte(byte(u))
		} else {
			b.WriteByte('-')
		}
	}
	s := b.String()
	if len(s) <= maxSlug {
		return s
	}
	// Java's String.hashCode over the UTF-16 units, as a 32-bit integer.
	var h int32
	for _, u := range units {
		h = 31*h + int32(u)
	}
	v := int64(h)
	if v < 0 {
		v = -v
	}
	return s[:maxSlug] + "-" + strconv.FormatInt(v, 36)
}

// ProjectDir is where config folder configDir keeps the conversations of
// the working folder cwd.
func ProjectDir(configDir, cwd string) string {
	if real, err := filepath.EvalSymlinks(cwd); err == nil {
		cwd = real
	}
	return filepath.Join(configDir, "projects", Slug(cwd))
}

// Latest is the conversation in projectDir that went on last, and when:
// the one with the newest prompt or reply. Claude Code also writes rows
// that are neither, such as cost-state and last-prompt when a session
// ends, so the file's time says less. A session that only ran local
// commands (/mcp, /usage, /status), or an empty file, has no conversation
// and is skipped.
func Latest(projectDir string) (path string, at time.Time, ok bool) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return "", time.Time{}, false
	}
	type file struct {
		path string
		mod  time.Time
	}
	var files []file
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		if info, err := e.Info(); err == nil && info.Size() > 0 {
			files = append(files, file{filepath.Join(projectDir, e.Name()), info.ModTime()})
		}
	}
	// A row is written after it is made, so no message is newer than its
	// file: newest files first, until one is older than the best message
	// (with room for file times kept to the second or two).
	slices.SortFunc(files, func(a, b file) int { return b.mod.Compare(a.mod) })
	for _, f := range files {
		if ok && f.mod.Add(2*time.Second).Before(at) {
			break
		}
		if t, found := lastMessage(f.path, f.mod); found && (!ok || t.After(at)) {
			path, at, ok = f.path, t, true
		}
	}
	return path, at, ok
}

// lastMessage is when the conversation at path last had a prompt or reply:
// that row's timestamp, or mod (the file's time) when it has none. A file
// whose end has neither but that is longer than the part read counts as a
// conversation at mod: its end may be one long row.
func lastMessage(path string, mod time.Time) (time.Time, bool) {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		if at, ok := message(lines[i]); ok {
			if at.IsZero() {
				return mod, true
			}
			return at, true
		}
	}
	if info, err := os.Stat(path); err == nil && info.Size() > tailSize {
		return mod, true
	}
	return time.Time{}, false
}

// message reports whether l is a prompt or reply of the main conversation,
// and when it was made (zero when the row does not say). A local command
// (/mcp, /usage) leaves user rows that start no request (later requests
// carry them as context, after a caveat); they do not count, and neither
// do rows Claude Code adds itself (isMeta).
func message(l []byte) (time.Time, bool) {
	if !bytes.Contains(l, []byte(`"user"`)) && !bytes.Contains(l, []byte(`"assistant"`)) {
		return time.Time{}, false
	}
	var e struct {
		Type      string    `json:"type"`
		Sidechain bool      `json:"isSidechain"`
		Meta      bool      `json:"isMeta"`
		At        time.Time `json:"timestamp"`
		Message   struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal(l, &e) != nil || e.Sidechain || e.Meta {
		return time.Time{}, false
	}
	switch e.Type {
	case "assistant":
		return e.At, true
	case "user":
		var text string
		if json.Unmarshal(e.Message.Content, &text) == nil {
			text = strings.TrimSpace(text)
			return e.At, text != "" && !localCommand(text)
		}
		var blocks []json.RawMessage
		return e.At, json.Unmarshal(e.Message.Content, &blocks) == nil && len(blocks) > 0
	}
	return time.Time{}, false
}

// ID is the session id of a conversation file.
func ID(path string) string { return strings.TrimSuffix(filepath.Base(path), ".jsonl") }

// tailSize is how much of the end of a conversation is read for its last
// reply, limit, permission mode and agent, which are all near the end.
const tailSize = 1 << 20

// synthetic is the model Claude Code names on messages it makes itself,
// such as the one for a request refused at a usage limit.
const synthetic = "<synthetic>"

// tail returns the lines of the last tailSize bytes of the conversation at
// path, oldest first; nil when it cannot be read. The first line may be cut
// and the last one half written: neither parses, so both are skipped.
func tail(path string) [][]byte {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	if info, err := f.Stat(); err == nil && info.Size() > tailSize {
		if _, err := f.Seek(info.Size()-tailSize, io.SeekStart); err != nil {
			return nil
		}
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	return bytes.Split(data, []byte("\n"))
}

// Size is about how many tokens the next reply in the conversation reads:
// the prompt of its last reply (input, cache read and cache creation) plus
// that reply. 0 when unknown.
func Size(path string) int {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if !bytes.Contains(l, []byte(`"usage"`)) || !bytes.Contains(l, []byte(`"assistant"`)) {
			continue
		}
		var e struct {
			Type      string `json:"type"`
			Sidechain bool   `json:"isSidechain"`
			Message   struct {
				Model string `json:"model"`
				Usage *struct {
					Input       int `json:"input_tokens"`
					CacheRead   int `json:"cache_read_input_tokens"`
					CacheCreate int `json:"cache_creation_input_tokens"`
					Output      int `json:"output_tokens"`
				} `json:"usage"`
			} `json:"message"`
		}
		if json.Unmarshal(l, &e) != nil || e.Type != "assistant" || e.Sidechain || e.Message.Model == synthetic || e.Message.Usage == nil {
			continue
		}
		u := e.Message.Usage
		return u.Input + u.CacheRead + u.CacheCreate + u.Output
	}
	return 0
}

// Model is the model of the conversation's last reply; "" when unknown.
func Model(path string) string {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		if m, ok := replyModel(lines[i]); ok {
			return m
		}
	}
	return ""
}

// replyModel reads the model of a main-thread reply row; messages Claude
// Code made itself do not count.
func replyModel(l []byte) (string, bool) { return rowModel(l, false) }

// rowModel is replyModel for the rows whose isSidechain is side.
func rowModel(l []byte, side bool) (string, bool) {
	if !bytes.Contains(l, []byte(`"model"`)) || !bytes.Contains(l, []byte(`"assistant"`)) {
		return "", false
	}
	var e struct {
		Type      string `json:"type"`
		Sidechain bool   `json:"isSidechain"`
		Message   struct {
			Model string `json:"model"`
		} `json:"message"`
	}
	if json.Unmarshal(l, &e) != nil || e.Type != "assistant" || e.Sidechain != side || e.Message.Model == "" || e.Message.Model == synthetic {
		return "", false
	}
	return e.Message.Model, true
}

// Limit is a usage limit at which Claude Code refused a request.
type Limit struct {
	// Window is the server's name for the window that is full:
	// five_hour, seven_day, seven_day_opus, seven_day_sonnet,
	// seven_day_overage_included (the Fable allowance) or overage
	// (usage credits).
	Window   string
	ResetsAt time.Time // when the window starts over; zero when not given
	At       time.Time // when the request was refused
	Model    string    // the model of the reply before it; "" when unknown
	Text     string    // the message Claude Code showed for it (its text blocks, a line each)
}

// LastLimit finds the usage limit that ended the last turn of the
// conversation at path, if it was hit at since or later: the row Claude
// Code (2.1.304-2.1.306) writes for a request it refused, an assistant
// message with isApiErrorMessage, apiError "usage_limit_reached" and the
// server's rate_limit_info at status "rejected". A subagent's row
// (isSidechain) and other API errors, such as a model that requires usage
// credits, do not count, and neither does a limit a person's prompt came
// after: that prompt began a turn of its own.
func LastLimit(path string, since time.Time) (Limit, bool) { return lastLimit(path, since, false) }

// SubagentLimit is LastLimit for the conversation of a subagent at path
// (SubagentPath), whose rows are all isSidechain: Claude Code writes a
// subagent's refusal there, not into its session's conversation. Model is
// the subagent's.
func SubagentLimit(path string, since time.Time) (Limit, bool) { return lastLimit(path, since, true) }

// lastLimit reads the limit of LastLimit from the rows whose isSidechain is
// side.
func lastLimit(path string, since time.Time, side bool) (Limit, bool) {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if !bytes.Contains(l, []byte(`"usage_limit_reached"`)) {
			if prompt(l) {
				return Limit{}, false
			}
			continue
		}
		var row struct {
			Type      string    `json:"type"`
			Sidechain bool      `json:"isSidechain"`
			APIError  bool      `json:"isApiErrorMessage"`
			Kind      string    `json:"apiError"`
			At        time.Time `json:"timestamp"`
			Params    struct {
				Info *struct {
					Status   string  `json:"status"`
					ResetsAt float64 `json:"resetsAt"` // Unix seconds
					Window   string  `json:"rateLimitType"`
				} `json:"rate_limit_info"`
			} `json:"apiErrorParams"`
			Message struct {
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(l, &row) != nil || row.Type != "assistant" || row.Sidechain != side || !row.APIError ||
			row.Kind != "usage_limit_reached" || row.Params.Info == nil || row.Params.Info.Status != "rejected" {
			continue
		}
		if row.At.Before(since) {
			return Limit{}, false // the ones before it are older still
		}
		lim := Limit{Window: row.Params.Info.Window, At: row.At}
		var texts []string
		for _, c := range row.Message.Content {
			if c.Type == "text" {
				texts = append(texts, c.Text)
			}
		}
		lim.Text = strings.Join(texts, "\n") // as Claude Code passes it to hooks
		if s := row.Params.Info.ResetsAt; s > 0 {
			lim.ResetsAt = time.Unix(int64(s), 0)
		}
		for j := i - 1; j >= 0; j-- {
			if m, ok := rowModel(lines[j], side); ok {
				lim.Model = m
				break
			}
		}
		return lim, true
	}
	return Limit{}, false
}

// PermissionMode is the permission mode the conversation was last in: the
// newer of the permission-mode rows Claude Code writes when the mode
// changes (and again at exit) and the mode a person's prompt was sent in;
// "" when the end of the conversation has neither.
func PermissionMode(path string) string {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if !bytes.Contains(l, []byte(`"permissionMode"`)) {
			continue
		}
		var e struct {
			Type string `json:"type"`
			Mode string `json:"permissionMode"`
		}
		if json.Unmarshal(l, &e) != nil || e.Mode == "" {
			continue
		}
		if e.Type == "permission-mode" || prompt(l) {
			return e.Mode
		}
	}
	return ""
}

// PromptedAfter reports whether a person sent a prompt in the conversation
// at path after Claude Code refused a request there at a usage limit at
// the time at (a Limit's At). The rows Claude Code adds after a limit by
// itself (a finished task, a queued prompt, Remote Control, file history,
// system rows, those of a local command such as /usage) do not count. A
// limit that is no longer in the end of the conversation counts as
// prompted after: ccam cannot tell.
func PromptedAfter(path string, at time.Time) bool {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if prompt(l) {
			return true
		}
		if !bytes.Contains(l, []byte(`"usage_limit_reached"`)) {
			continue
		}
		var row struct {
			Sidechain bool      `json:"isSidechain"`
			At        time.Time `json:"timestamp"`
		}
		if json.Unmarshal(l, &row) == nil && !row.Sidechain && row.At.Equal(at) {
			return false
		}
	}
	return true
}

// PromptedSince reports whether a person sent a prompt in the conversation
// at path after the time at: for a limit that is not in that conversation
// (a subagent's, SubagentLimit), so PromptedAfter cannot find it there.
// Only the end of the conversation is read.
func PromptedSince(path string, at time.Time) bool {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		if !prompt(lines[i]) {
			continue
		}
		var row struct {
			At time.Time `json:"timestamp"`
		}
		// The newest prompt says it: those before it are older.
		return json.Unmarshal(lines[i], &row) == nil && row.At.After(at)
	}
	return false
}

// RemoteControl reports whether the conversation at path was connected to
// Remote Control when Claude Code last wrote it: the last bridge-session
// row near the end names a bridge session. Claude Code (2.1.308-2.1.310)
// writes that row when Remote Control connects, again with the session's
// state every 32 KiB it writes while connected, and one with an empty
// bridgeSessionId when Remote Control is turned off.
func RemoteControl(path string) bool {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if !bytes.Contains(l, []byte(`"bridge-session"`)) {
			continue
		}
		var e struct {
			Type string `json:"type"`
			ID   string `json:"bridgeSessionId"`
		}
		if json.Unmarshal(l, &e) == nil && e.Type == "bridge-session" {
			return e.ID != ""
		}
	}
	return false
}

// Agent is the agent the newest agent-setting row near the end of the
// conversation at path names; "" when there is none. Claude Code (2.1.314)
// writes these rows with the session's metadata (after every 32 KiB of the
// conversation, and as the session ends) for the agent the session started
// as (--agent or the agent setting), or, once it resumed the conversation
// by id, for the conversation's own agent: also when --agent or the agent
// setting runs the session as another.
func Agent(path string) string {
	lines := tail(path)
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		if !bytes.Contains(l, []byte(`"agent-setting"`)) {
			continue
		}
		var e struct {
			Type  string `json:"type"`
			Agent string `json:"agentSetting"`
		}
		if json.Unmarshal(l, &e) == nil && e.Type == "agent-setting" {
			return e.Agent
		}
	}
	return ""
}

// prompt reports whether l is a prompt a person sent (entry.prompt).
func prompt(l []byte) bool {
	if !bytes.Contains(l, []byte(`"user"`)) {
		return false
	}
	var e entry
	if json.Unmarshal(l, &e) != nil {
		return false
	}
	_, ok := e.prompt()
	return ok
}

// entry is what ccam reads of a user or assistant row of a conversation.
type entry struct {
	Type      string  `json:"type"`
	Sidechain bool    `json:"isSidechain"`
	Meta      bool    `json:"isMeta"`
	Compact   bool    `json:"isCompactSummary"`
	Hidden    bool    `json:"isVisibleInTranscriptOnly"`
	Origin    *origin `json:"origin"`
	Message   struct {
		ID      string          `json:"id"`
		Model   string          `json:"model"`
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

// origin says who sent a user row: "human", or one of Claude Code's own
// (task-notification, peer …). Older versions write none.
type origin struct {
	Kind string `json:"kind"`
}

// prompt is the prompt a person sent in e, a main-thread user row that is
// not a tool's result, as they sent it (personText); ok is false for any
// other row. Prompts Claude Code sends itself (a finished task, a peer's
// message, a local command's caveat, the summary it writes after
// compacting) do not count.
func (e *entry) prompt() (text string, ok bool) {
	if e.Type != "user" || e.Meta || e.Sidechain || e.Compact || e.Hidden || e.Origin != nil && e.Origin.Kind != "human" {
		return "", false
	}
	text, tool := contentText(e.Message.Content)
	if tool {
		return "", false
	}
	return personText(text)
}

// personText is text, of a user row or a queued prompt, as a person sent
// it, and whether it is a person's prompt at all: not empty, and not one of
// the texts Claude Code (2.1.303-2.1.310) writes without isMeta or origin
// for what never reaches the model as a request (a local command or its
// output, a ! command's output, Esc), nor a finished task's notification
// that came without its origin. A slash command the person ran is as they
// typed it, and a ! command is "! <command>", as Claude Code quotes it: it
// queries the model for one unless respondToBashCommands is off (2.1.310).
func personText(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if text == "" || localCommand(text) {
		return "", false
	}
	if cmd, ok := inside(text, "bash-input"); ok && strings.HasPrefix(text, "<bash-input>") {
		return "! " + cmd, true
	}
	for _, own := range []string{"<bash-stdout>", "<bash-stderr>", "<task-notification>", "[Request interrupted by user"} {
		if strings.HasPrefix(text, own) {
			return "", false
		}
	}
	if strings.HasPrefix(text, "<command-message>") {
		return command(text), true
	}
	return text, true
}

// command is the slash command a person ran, as they typed it, from the
// text Claude Code writes for it: <command-message>name</command-message>,
// <command-name>/name</command-name> and, with arguments,
// <command-args>…</command-args>.
func command(text string) string {
	name, ok := inside(text, "command-name")
	if !ok {
		return text
	}
	if args, _ := inside(text, "command-args"); args != "" {
		name += " " + args
	}
	return name
}

// inside is the text of the first <tag>…</tag> in s, trimmed.
func inside(s, tag string) (string, bool) {
	_, after, ok := strings.Cut(s, "<"+tag+">")
	if !ok {
		return "", false
	}
	in, _, ok := strings.Cut(after, "</"+tag+">")
	return strings.TrimSpace(in), ok
}

// contentText is the text of a message's content, a string or blocks: their
// text blocks with a line break between, and a line such as [image] for
// each other block; tool reports a tool's result among the blocks.
func contentText(content json.RawMessage) (text string, tool bool) {
	if json.Unmarshal(content, &text) == nil {
		return text, false
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(content, &blocks) != nil {
		return "", false
	}
	lines := make([]string, 0, len(blocks))
	for _, b := range blocks {
		switch b.Type {
		case "text":
			lines = append(lines, b.Text)
		case "tool_result":
			tool = true
		case "":
		default:
			lines = append(lines, "["+b.Type+"]")
		}
	}
	return strings.Join(lines, "\n"), tool
}

// localCommand reports whether the text of a user row is a local command
// (/mcp, /usage) or its output, which Claude Code (2.1.306-2.1.310) writes
// without isMeta or origin.
func localCommand(text string) bool {
	return strings.HasPrefix(text, "<command-name>") || strings.HasPrefix(text, "<local-command-")
}

// Title is the conversation's name: the one given with /rename, else the
// one Claude Code made up; "" when it has neither.
func Title(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	custom, ai := "", ""
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64<<10), 16<<20)
	for sc.Scan() {
		l := sc.Bytes()
		if !bytes.Contains(l, []byte(`-title"`)) {
			continue
		}
		var e struct {
			Type   string `json:"type"`
			Custom string `json:"customTitle"`
			AI     string `json:"aiTitle"`
		}
		if json.Unmarshal(l, &e) != nil {
			continue
		}
		switch e.Type {
		case "custom-title":
			custom = e.Custom
		case "ai-title":
			ai = e.AI
		}
	}
	if custom != "" {
		return custom
	}
	return ai
}
