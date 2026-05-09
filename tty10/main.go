package tty10

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/nyaosorg/go-ttyadapter/internal/tty10base"
	"github.com/nyaosorg/go-ttyadapter/internal/virtualterminal"
	"github.com/nyaosorg/go-ttyadapter/internal/winch10"
)

type Tty struct {
	devTty  *os.File
	disable func()
	cancel  func()
	buf     []byte
}

func (tt *Tty) fd() int {
	return int(tt.devTty.Fd())
}

func (tt *Tty) Open(onResize func(int, int)) error {
	var err error

	tt.devTty, err = os.OpenFile(tty10base.TtyPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	tt.disable, err = virtualterminal.Enable(tt.fd())
	if err != nil {
		return err
	}

	if onResize != nil {
		tt.cancel, err = winch10.Notice(onResize)
		return err
	}
	return nil
}

func (tt *Tty) GetKey() (key string, err error) {
	if len(tt.buf) <= 0 {
		tt.buf, err = tty10base.ReadAllInRawMode(tt.devTty)
		if err != nil || len(tt.buf) <= 0 {
			return
		}
	}
	r, size := utf8.DecodeRune(tt.buf)
	if r == '\x1B' {
		key, tt.buf = string(tt.buf), nil
	} else {
		key, tt.buf = string(tt.buf[:size]), tt.buf[size:]
	}
	return
}

func (tt *Tty) Close() error {
	if tt.cancel != nil {
		tt.cancel()
		tt.cancel = nil
	}
	if tt.disable != nil {
		tt.disable()
		tt.disable = nil
	}
	if tt.devTty != nil {
		tt.devTty.Close()
	}
	return nil
}

func (tt *Tty) Size() (int, int, error) {
	return term.GetSize(int(os.Stderr.Fd()))
}
