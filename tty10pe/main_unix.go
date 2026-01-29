//go:build !windows
// +build !windows

package tty10pe

func enable(handle int) (func(), error) {
	return func() {}, nil
}
