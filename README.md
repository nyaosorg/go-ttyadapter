go-ttyadapter
=============

This package provides an abstraction layer for reading keyboard input from terminal devices.

Features
--------

- Backends for common terminal implementations  
  - `tty8` → wrapper around [github.com/mattn/go-tty][go-tty]
  - `tty10` → wrapper around [golang.org/x/term][xterm]
  - `auto` → pseudo terminal that feeds programmed key sequences (useful for tests)

- Automatic handling of terminal modes  
  The library switches terminal modes (raw/cooked) internally as needed so callers do not have to manage mode changes manually.

- Test-friendly pseudo terminal for automated tests  
  Use `auto.Pilot` to simulate key sequences and verify interactive behavior in `go test` or CI.

- Designed for integration with readline-like tools  
  Intended to be a reusable terminal-input layer for libraries that implement line-editing or interactive selection UIs.

- Lightweight wrappers (not a pure-zero-dependency package)  
  This module provides small adapters over existing terminal libraries rather than reimplementing low-level terminal handling. Backends depend on [github.com/mattn/go-tty][go-tty] or [golang.org/x/term][xterm] as appropriate.

- Optional "pending escape" backends for robust key sequence handling  
  Some terminal environments may deliver escape sequences (such as arrow keys) in multiple chunks.  
  To handle this reliably, the package provides alternative backends that treat `ESC` as a prefix key rather than a standalone input.

  - `tty8pe` → variant of `tty8` with pending-escape handling  
    (`ESC` is interpreted as a prefix, similar to traditional UNIX terminals)
  - `tty10pe` → planned variant of `tty10` with the same behavior

  These backends avoid ambiguity between standalone `ESC` and escape sequences, at the cost of not supporting `ESC` as an independent key.

Example
-------

```examples/example.go
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
```

This example prints pressed keys enclosed in &lt; &gt; and exits when Ctrl-G is pressed.
The demo below shows two modes in sequence:

- Simulated input using auto.Pilot (numbers 1–5 are automatically typed)
- Manual input where actual keypresses (e.g. arrow keys) are captured from the terminal

![demo.gif](demo.gif)

[go-tty]: https://github.com/mattn/go-tty
[xterm]: https://pkg.go.dev/golang.org/x/term

Release notes
-------------

- [English](release_note.md)
- [Japanese](release_note_ja.md)

License
-------

MIT License

Author
------

- [hymkor (HAYAMA Kaoru)](https://github.com/hymkor)
