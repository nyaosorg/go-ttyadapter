package tty8pe

import (
	"github.com/mattn/go-tty/v2"

	"github.com/nyaosorg/go-ttyadapter/internal/unsurrogate"
)

func getOneKey(tty *tty.TTY, onPrefix func(string)) (string, error) {
	var sequence string
	for {
		r, _, err := unsurrogate.ReadRune(tty.ReadRune)
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

func getKeys(tty *tty.TTY, onPrefix func(string)) ([]string, error) {
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
