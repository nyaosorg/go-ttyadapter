//go:build !windows
// +build !windows

package virtualterminal

func enable(handle int) (func(), error) {
	return func() {}, nil
}
