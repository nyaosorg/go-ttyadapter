package tty10base

func EnableVirtualTerminal(handle int) (func(), error) {
	return enable(handle)
}
