//go:build run

package main

import (
	"fmt"
	"os"

	"github.com/nyaosorg/go-ttyadapter/auto"

	"github.com/nyaosorg/go-ttyadapter/internal/demo"
)

func main() {
	ap := &auto.Pilot{
		Text: []string{
			"a",
			"b",
			"c",
		},
	}
	if err := demo.Exec(ap); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
