package tty10pe

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/nyaosorg/go-ttyadapter/internal/virtualterminal"
	"github.com/nyaosorg/go-ttyadapter/internal/winch10"
)

type Tty struct {
	disable  func()
	cancel   func()
	buf      []byte
	OnPrefix func(string)
}

func (m *Tty) SetOnPrefix(f func(string)) (original func(string)) {
	original = m.OnPrefix
	m.OnPrefix = f
	return
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
	for {
		if len(M.buf) <= 0 {
			M.buf, err = M.getKey()
			if err != nil || len(M.buf) <= 0 {
				return
			}
		}
		_, size := utf8.DecodeRune(M.buf)
		key += string(M.buf[:size])
		M.buf = M.buf[size:]
		switch key {
		case "\x1B", "\x1B[", "\x1B[1", "\x1B[15", "\x1B[16", "\x1B[17",
			"\x1B[18", "\x1B[1;", "\x1B[1;5", "\x1B[2", "\x1B[20",
			"\x1B[21", "\x1B[23", "\x1B[24", "\x1B[3", "\x1B[5",
			"\x1B[5;", "\x1B[5;5", "\x1B[6", "\x1B[6;", "\x1B[6;5", "\x1B[O":
			if M.OnPrefix != nil {
				M.OnPrefix(key)
			}
		case "\x1B\x1B":
			key = "\x1B"
		default:
			return
		}
	}
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
