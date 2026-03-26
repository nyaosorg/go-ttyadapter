//go:build run

package main

import (
	"fmt"
	"os"

	"github.com/nyaosorg/go-ttyadapter/tty10pe"

	"github.com/nyaosorg/go-ttyadapter/internal/demo"
)

func onPrefix(s string) {
	fmt.Printf("onPrefix: %q\r\n", s)
}

func main() {
	if err := demo.Exec(&tty10pe.Tty{OnPrefix: onPrefix}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
