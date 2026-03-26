//go:build run

package main

import (
	"fmt"
	"os"

	"github.com/nyaosorg/go-ttyadapter/tty8pe"

	"github.com/nyaosorg/go-ttyadapter/internal/demo"
)

func main() {
	if err := demo.Exec(&tty8pe.Tty{}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
