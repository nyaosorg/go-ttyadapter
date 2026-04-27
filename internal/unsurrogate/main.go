package unsurrogate

import (
	"unicode/utf16"
	"unicode/utf8"
)

func ReadRune(read func() (rune, int, error)) (rune, int, error) {
	var surrogated rune = 0
	for {
		r, siz, err := read()
		if err != nil {
			return 0, siz, err
		}
		if r == 0 {
			continue
		}
		if surrogated > 0 {
			r := utf16.DecodeRune(surrogated, r)
			return r, utf8.RuneLen(r), nil
		} else if !utf16.IsSurrogate(r) {
			return r, siz, nil
		} else {
			// surrogate pair first word.
			surrogated = r
		}
	}
}
