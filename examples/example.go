//go:build run

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nyaosorg/go-ttyadapter"
	"github.com/nyaosorg/go-ttyadapter/auto"
	"github.com/nyaosorg/go-ttyadapter/tty8pe"
)

var flagInterval = flag.Uint("interval", 0, "delay (seconds) between simulated key inputs")

func run(operations []string) error {
	var tty ttyadapter.Tty

	// If command-line arguments are given, simulate key inputs using auto.Pilot.
	if len(operations) > 0 {
		// Append Ctrl-G to end to exit automatically.
		operations = append(operations, "\x07")
		var hook func(*auto.Pilot) error
		if *flagInterval > 0 {
			hook = func(_ *auto.Pilot) error {
				time.Sleep(time.Duration(*flagInterval) * time.Second)
				return nil
			}
		}
		tty = &auto.Pilot{Text: operations, OnGetKey: hook}
	} else {
		tty = &tty8pe.Tty{}
	}

	if err := tty.Open(nil); err != nil {
		return err
	}
	defer tty.Close()

	for {
		key, err := tty.GetKey()
		if err != nil {
			return err
		}
		if key == "\x07" {
			return nil
		}
		// Show typed key, converting ESC to literal name.
		fmt.Printf("<%s>\n", strings.ReplaceAll(key, "\x1B", "ESC"))
	}
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
