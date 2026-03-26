package virtualterminal

func Enable(handle int) (func(), error) {
	return enable(handle)
}
