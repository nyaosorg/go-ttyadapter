//go:build run

package main

import (
	"fmt"
	"os"

	"github.com/nyaosorg/go-ttyadapter/tty8"

	"github.com/nyaosorg/go-ttyadapter/internal/demo"
)

func main() {
	if err := demo.Exec(&tty8.Tty{}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
