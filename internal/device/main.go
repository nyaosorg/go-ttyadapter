package device

import (
	"os"

	"golang.org/x/term"
)

func ReadAllInRawMode(devTty *os.File) ([]byte, error) {
	fd := int(devTty.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	defer term.Restore(fd, oldState)

	var buffer [1024]byte
	n, err := devTty.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}
