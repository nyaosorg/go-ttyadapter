package tty8

import (
	"strings"

	"github.com/mattn/go-tty"

	"github.com/nyaosorg/go-ttyadapter/internal/tty8base"
)

// Tty is a wrapper around github.com/mattn/go-tty.
// While go-tty reads input per rune, control keys may be sent as multiple
// runes. To handle this, Tty buffers input and provides the GetKey method
// to retrieve keys per physical key press. It also replaces the terminal
// size notification mechanism from a channel of go-tty's WINSIZE to a
// callback function, making it easier to abstract single-character input
// handling.
type Tty struct {
	*tty.TTY
	stopNotice func()
	buf        []string
}

func (tt *Tty) IsOpen() bool {
	return tt.TTY != nil
}

// Open calls go-tty's Open method to initialize the Tty instance.
// It also starts a goroutine that listens for terminal resize notifications.
// The goroutine receives events from go-tty's SIGWINCH channel and,
// if onResize is not nil, invokes the provided callback function.
func (tt *Tty) Open(onResize func(width, height int)) (err error) {
	defer func() {
		if err != nil {
			tt.Close()
		}
	}()

	if tt.TTY == nil {
		tt.TTY, err = tty.Open()
		if err != nil {
			return
		}
	}
	if tt.stopNotice == nil {
		tt.stopNotice, err = tty8base.Notice(tt.TTY, onResize)
	}
	return
}

func getKeys(tty *tty.TTY) ([]string, error) {
	keys := []string{}

	var buffer strings.Builder
	escape := false
	for {
		r, err := tty8base.ReadRune(tty.ReadRune)
		if err != nil {
			return nil, err
		}
		buffer.WriteRune(r)
		if r == '\x1B' {
			escape = true
		}
		if !(escape && tty.Buffered()) && buffer.Len() > 0 {
			keys = append(keys, buffer.String())
			if !tty.Buffered() {
				return keys, nil
			}
		}
	}
}

// GetKey switches the terminal to raw mode and reads a single key input.
// Since control keys may consist of multiple runes, the result is returned
// as a string. Any unread input is buffered internally and returned on
// subsequent calls. After processing, the terminal is restored to cooked
// mode.
func (tt *Tty) GetKey() (key string, err error) {
	if len(tt.buf) <= 0 {
		clean, _err := tt.TTY.Raw()
		if err != nil {
			return "", _err
		}
		defer clean()

		tt.buf, err = getKeys(tt.TTY)
		if err != nil || len(tt.buf) <= 0 {
			return
		}
	}
	key, tt.buf = tt.buf[0], tt.buf[1:]
	return
}

// Close calls go-tty's `Close` method to shut down the Tty instance.
// It clears internal references (by overwriting them with nil) to prevent
// reuse. Since go-tty closes the SIGWINCH channel, the goroutine started
// by Open detects the channel closure and terminates automatically.
func (tt *Tty) Close() (err error) {
	if tt.TTY != nil {
		err = tt.TTY.Close()
		tt.TTY = nil
	}
	if tt.stopNotice != nil {
		tt.stopNotice()
		tt.stopNotice = nil
	}
	return
}
