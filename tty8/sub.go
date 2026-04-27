package tty8

import (
	"github.com/mattn/go-tty/v2"
	"strings"

	"github.com/nyaosorg/go-ttyadapter/internal/unsurrogate"
)

func getOneKey(tty *tty.TTY) (string, error) {
	var buffer strings.Builder
	escape := false
	for {
		r, _, err := unsurrogate.ReadRune(tty.ReadRune)
		if err != nil {
			return "", err
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

func getKeys(tty *tty.TTY) ([]string, error) {
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
