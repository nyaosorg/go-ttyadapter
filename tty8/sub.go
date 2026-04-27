package tty8

import (
	"strings"
	"unicode/utf16"
)

// xTty is the interface of tty to use GetKey function.
type xTty interface {
	Raw() (func() error, error)
	ReadRune() (rune, int, error)
	Buffered() bool
}

func getOneKey(tty xTty) (string, error) {
	var buffer strings.Builder
	escape := false
	var surrogated rune = 0
	for {
		r, _, err := tty.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 0 {
			continue
		}
		if surrogated > 0 {
			r = utf16.DecodeRune(surrogated, r)
			surrogated = 0
		} else if utf16.IsSurrogate(r) { // surrogate pair first word.
			surrogated = r
			continue
		}
		buffer.WriteRune(r)
		if r == '\x1B' {
			escape = true
		}
		if !(escape && tty.Buffered()) && buffer.Len() > 0 {
			return buffer.String(), nil
		}
	}
}

func getKeys(tty xTty) ([]string, error) {
	clean, err := tty.Raw()
	if err != nil {
		return nil, err
	}
	defer clean()

	keys := []string{}

	for {
		key1, err := getOneKey(tty)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key1)
		if !tty.Buffered() {
			return keys, nil
		}
	}
}
