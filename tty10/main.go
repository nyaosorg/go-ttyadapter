package tty10

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/nyaosorg/go-ttyadapter/internal/tty10base"
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

func (tt *Tty) IsOpen() bool {
	return tt.devTty != nil
}

func (tt *Tty) Open(onResize func(int, int)) (err error) {
	defer func() {
		if err != nil {
			tt.Close()
		}
	}()

	if tt.devTty == nil {
		tt.devTty, err = os.OpenFile(tty10base.TtyPath, os.O_RDWR, 0666)
		if err != nil {
			return
		}
	}

	if tt.disable == nil {
		tt.disable, err = tty10base.EnableVirtualTerminal(tt.fd())
		if err != nil {
			return
		}
	}

	if onResize != nil {
		tt.cancel, err = tty10base.Notice(onResize)
		return
	}
	return
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
		tt.devTty = nil
	}
	return nil
}

func (tt *Tty) Size() (int, int, error) {
	return term.GetSize(int(os.Stderr.Fd()))
}
