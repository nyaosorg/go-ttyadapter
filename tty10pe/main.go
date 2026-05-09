package tty10pe

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/nyaosorg/go-ttyadapter/internal/tty10base"
)

type Tty struct {
	devTty   *os.File
	disable  func()
	cancel   func()
	buf      []byte
	OnPrefix func(string)
}

func (tt *Tty) fd() int {
	return int(tt.devTty.Fd())
}

func (m *Tty) SetOnPrefix(f func(string)) (original func(string)) {
	original = m.OnPrefix
	m.OnPrefix = f
	return
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
	for {
		if len(tt.buf) <= 0 {
			tt.buf, err = tty10base.ReadAllInRawMode(tt.devTty)
			if err != nil || len(tt.buf) <= 0 {
				return
			}
		}
		_, size := utf8.DecodeRune(tt.buf)
		key += string(tt.buf[:size])
		tt.buf = tt.buf[size:]
		switch key {
		case "\x1B", "\x1B[", "\x1B[1", "\x1B[15", "\x1B[16", "\x1B[17",
			"\x1B[18", "\x1B[1;", "\x1B[1;5", "\x1B[2", "\x1B[20",
			"\x1B[21", "\x1B[23", "\x1B[24", "\x1B[3", "\x1B[5",
			"\x1B[5;", "\x1B[5;5", "\x1B[6", "\x1B[6;", "\x1B[6;5", "\x1B[O":
			if tt.OnPrefix != nil {
				tt.OnPrefix(key)
			}
		case "\x1B\x1B":
			key = "\x1B"
		default:
			return
		}
	}
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
