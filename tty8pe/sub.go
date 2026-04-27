package tty8pe

import (
	"unicode/utf16"
)

// xTty is the interface of tty to use GetKey function.
type xTty interface {
	Raw() (func() error, error)
	ReadRune() (rune, int, error)
	Buffered() bool
}

func readRune(tty xTty) (rune, error) {
	var surrogated rune = 0
	for {
		r, _, err := tty.ReadRune()
		if err != nil {
			return 0, err
		}
		if r == 0 {
			continue
		}
		if surrogated > 0 {
			return utf16.DecodeRune(surrogated, r), nil
		} else if !utf16.IsSurrogate(r) {
			return r, nil
		} else {
			// surrogate pair first word.
			surrogated = r
		}
	}
}

func getOneKey(tty xTty, onPrefix func(string)) (string, error) {
	var sequence string
	for {
		r, err := readRune(tty)
		if err != nil {
			return "", err
		}
		sequence += string(r)
		switch sequence {
		case "\x1B", "\x1B[", "\x1B[1", "\x1B[15", "\x1B[16", "\x1B[17",
			"\x1B[18", "\x1B[1;", "\x1B[1;5", "\x1B[2", "\x1B[20",
			"\x1B[21", "\x1B[23", "\x1B[24", "\x1B[3", "\x1B[5",
			"\x1B[5;", "\x1B[5;5", "\x1B[6", "\x1B[6;", "\x1B[6;5", "\x1B[O":
			if onPrefix != nil {
				onPrefix(sequence)
			}
		case "\x1B\x1B":
			sequence = "\x1B"
		default:
			return sequence, nil
		}
	}
}

func getKeys(tty xTty, onPrefix func(string)) ([]string, error) {
	clean, err := tty.Raw()
	if err != nil {
		return nil, err
	}
	defer clean()

	keys := []string{}

	for {
		key1, err := getOneKey(tty, onPrefix)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key1)
		if !tty.Buffered() {
			return keys, nil
		}
	}
}
