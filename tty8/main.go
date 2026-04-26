package tty8

import (
	"fmt"
	"sync"

	"github.com/mattn/go-tty/v2"
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
	wg  sync.WaitGroup
	buf []string
}

// Open calls go-tty's Open method to initialize the Tty instance.
// It also starts a goroutine that listens for terminal resize notifications.
// The goroutine receives events from go-tty's SIGWINCH channel and,
// if onResize is not nil, invokes the provided callback function.
func (m *Tty) Open(onResize func(width, height int)) error {
	var err error
	m.TTY, err = tty.Open()
	if err != nil {
		return fmt.Errorf("go-tty.Open: %w", err)
	}
	_lastw, _lasth, err := m.TTY.Size()
	if err != nil {
		return fmt.Errorf("go-tty.Size: %w", err)
	}
	if onResize != nil {
		ws := m.TTY.SIGWINCH()
		m.wg.Add(1)
		go func(lastw, lasth int) {
			for wh := range ws {
				if lastw != wh.W || lasth != wh.H {
					onResize(wh.W, wh.H)
					lastw = wh.W
					lasth = wh.H
				}
			}
			m.wg.Done()
		}(_lastw, _lasth)
	}
	return nil
}

// GetKey switches the terminal to raw mode and reads a single key input.
// Since control keys may consist of multiple runes, the result is returned
// as a string. Any unread input is buffered internally and returned on
// subsequent calls. After processing, the terminal is restored to cooked
// mode.
func (m *Tty) GetKey() (string, error) {
	if len(m.buf) <= 0 {
		var err error
		m.buf, err = getKeys(m.TTY)
		if err != nil || len(m.buf) <= 0 {
			return "", err
		}
	}
	var top string
	top, m.buf = m.buf[0], m.buf[1:]
	return top, nil
}

// Close calls go-tty's `Close` method to shut down the Tty instance.
// It clears internal references (by overwriting them with nil) to prevent
// reuse. Since go-tty closes the SIGWINCH channel, the goroutine started
// by Open detects the channel closure and terminates automatically.
func (m *Tty) Close() error {
	if m.TTY != nil {
		m.TTY.Close()
		m.TTY = nil
	}
	m.wg.Wait()
	return nil
}
