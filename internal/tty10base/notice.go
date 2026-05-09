package tty10base

import (
	"os"

	"golang.org/x/term"
)

func getSize() (int, int, error) {
	return term.GetSize(int(os.Stderr.Fd()))
}

func Notice(onResize func(int, int)) (func(), error) {
	return notice(onResize)
}
