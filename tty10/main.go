package tty10

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/nyaosorg/go-ttyadapter/internal/virtualterminal"
	"github.com/nyaosorg/go-ttyadapter/internal/winch10"
)

type Tty struct {
	disable func()
	cancel  func()
	buf     []byte
}

func (M *Tty) Open(onResize func(int, int)) error {
	var err error

	stdin := int(os.Stdin.Fd())
	M.disable, err = virtualterminal.Enable(stdin)
	if err != nil {
		return err
	}

	if onResize != nil {
		M.cancel, err = winch10.Notice(onResize)
		return err
	}
	return nil
}

func (M *Tty) getKey() ([]byte, error) {
	stdin := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(stdin)
	if err != nil {
		return nil, err
	}
	defer term.Restore(stdin, oldState)

	var buffer [1024]byte
	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (M *Tty) GetKey() (key string, err error) {
	if len(M.buf) <= 0 {
		M.buf, err = M.getKey()
		if err != nil || len(M.buf) <= 0 {
			return
		}
	}
	r, size := utf8.DecodeRune(M.buf)
	if r == '\x1B' {
		key, M.buf = string(M.buf), nil
	} else {
		key, M.buf = string(M.buf[:size]), M.buf[size:]
	}
	return
}

func (M *Tty) Close() error {
	if M.cancel != nil {
		M.cancel()
		M.cancel = nil
	}
	if M.disable != nil {
		M.disable()
		M.disable = nil
	}
	return nil
}

func (M *Tty) Size() (int, int, error) {
	return term.GetSize(int(os.Stderr.Fd()))
}
