package unsurrogate

import (
	"errors"
	"syscall"
	"time"
	"unicode/utf16"
)

// shouldRetry returns true if the error is EAGAIN or EWOULDBLOCK, which means
// the file descriptor is in non-blocking mode and no data is available yet.
//
// Since go-tty v0.0.8, it sets the terminal file descriptor to non-blocking
// mode (via syscall.SetNonblock) to avoid being blocked on Close. As a side
// effect, read operations may return EAGAIN instead of blocking on some
// environments (notably macOS). This function is used to detect that case
// and retry the read instead of treating it as a fatal error.
//
// Note: on macOS (and other BSD-derived systems), EWOULDBLOCK is defined as
// the same value as EAGAIN, so the EWOULDBLOCK checks are technically
// redundant on those platforms. They are kept here for portability and
// clarity on Linux, where the two constants are distinct.
func shouldRetry(err error) bool {
	// Cover the case where the error is wrapped (e.g. via fmt.Errorf("%w", err))
	if errors.Is(err, syscall.EAGAIN) {
		return true
	}
	if errors.Is(err, syscall.EWOULDBLOCK) {
		return true
	}
	// Cover the case where syscall.Errno is embedded in a bufio or other
	// intermediate layer that does not propagate it through the errors.Is chain.
	var errno syscall.Errno
	return errors.As(err, &errno) &&
		(errno == syscall.EAGAIN || errno == syscall.EWOULDBLOCK)
}

func ReadRune(read func() (rune, error)) (rune, error) {
	var surrogated rune = 0
	for {
		r, err := read()
		if err != nil {
			if shouldRetry(err) {
				time.Sleep(10 * time.Millisecond)
				continue
			}
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
