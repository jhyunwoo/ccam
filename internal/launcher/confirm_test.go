package launcher

import (
	"io"
	"os"
	"testing"
)

// ConfirmMove moves a conversation only on a y or yes read from the
// terminal. Without a terminal it may ask on, it reads nothing and the
// answer is no; a terminal that is gone by the time the answer is read
// (the job put in the background, say) gives no yes either.
func TestConfirmMoveNeedsAYes(t *testing.T) {
	oldIn, oldAttended := os.Stdin, attended
	t.Cleanup(func() { os.Stdin, attended = oldIn, oldAttended })
	for _, c := range []struct {
		typed    string
		terminal int // how many looks find the terminal there; -1: every one
		want     bool
		left     string // what stays unread
	}{
		{"y\n", -1, true, ""},
		{"Yes\n", -1, true, ""},
		{"n\n", -1, false, ""},
		{"\n", -1, false, ""},
		{"y\n", 0, false, "y\n"},
		{"n\n", 1, false, ""},
	} {
		looks := 0
		attended = func() bool { looks++; return c.terminal < 0 || looks <= c.terminal }
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		w.WriteString(c.typed)
		w.Close()
		os.Stdin = r
		got := ConfirmMove("ccam: move this conversation?")
		rest, _ := io.ReadAll(r)
		r.Close()
		if got != c.want || string(rest) != c.left {
			t.Errorf("%q, terminal for %d looks: %v, %q left unread", c.typed, c.terminal, got, rest)
		}
	}
}
