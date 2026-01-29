//go:build run

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nyaosorg/go-ttyadapter"
	"github.com/nyaosorg/go-ttyadapter/tty10pe"
)

func mains() error {
	var tty ttyadapter.Tty = &tty10pe.Tty{}
	var err error

	if len(os.Args) > 1 {
		err = tty.Open(func(w, h int) {
			fmt.Printf("Change size: %d, %d\n", w, h)
		})
	} else {
		err = tty.Open(nil)
	}

	if err != nil {
		return err
	}
	defer tty.Close()

	w, h, err := tty.Size()
	if err != nil {
		return err
	}
	fmt.Printf("(%d,%d)\n", w, h)

	key, err := tty.GetKey()
	if err != nil {
		return err
	}

	fmt.Printf("\"%s\"\n", strings.ReplaceAll(key, "\x1B", "ESC "))
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
