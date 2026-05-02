package unsurrogate

import (
	"unicode/utf16"
)

func ReadRune(read func() (rune, error)) (rune, error) {
	var surrogated rune = 0
	for {
		r, err := read()
		if err != nil {
			return 0, err
		}
		if r == 0 {
			continue
		}
		if surrogated > 0 {
			r := utf16.DecodeRune(surrogated, r)
			return r, nil
		} else if !utf16.IsSurrogate(r) {
			return r, nil
		} else {
			// surrogate pair first word.
			surrogated = r
		}
	}
}
