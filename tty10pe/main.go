package tty10pe

import (
	"os"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

type Tty struct {
	done    chan struct{}
	disable func()
	buf     []byte
}

func (M *Tty) Open(onResize func(int, int)) error {
	var err error

	stdin := int(os.Stdin.Fd())
	M.disable, err = enable(stdin)
	if err != nil {
		return err
	}

	if onResize != nil {
		w, h, err := M.Size()
		if err != nil {
			return err
		}
		M.done = make(chan struct{})
		go func(lastw, lasth int) {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-M.done:
					ticker.Stop()
					return
				case <-ticker.C:
					w, h, err := M.Size()
					if err == nil && (w != lastw || h != lasth) {
						onResize(w, h)
						lastw = w
						lasth = h
					}
				}
			}
		}(w, h)
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

func (M *Tty) GetKey() (string, error) {
	var sequence string
	for {
		if len(M.buf) <= 0 {
			var err error
			M.buf, err = M.getKey()
			if err != nil || len(M.buf) <= 0 {
				return "", err
			}
		}
		_, size := utf8.DecodeRune(M.buf)
		sequence += string(M.buf[:size])
		M.buf = M.buf[size:]
		switch sequence {
		case "\x1B", "\x1B[", "\x1B[1", "\x1B[15", "\x1B[16", "\x1B[17",
			"\x1B[18", "\x1B[1;", "\x1B[1;5", "\x1B[2", "\x1B[20",
			"\x1B[21", "\x1B[23", "\x1B[24", "\x1B[3", "\x1B[5",
			"\x1B[5;", "\x1B[5;5", "\x1B[6", "\x1B[6;", "\x1B[6;5", "\x1B[O":
		default:
			return sequence, nil
		}
	}
}

func (M *Tty) Close() error {
	if M.done != nil {
		M.done <- struct{}{}
		close(M.done)
		M.done = nil
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
