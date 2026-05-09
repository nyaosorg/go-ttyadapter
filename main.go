package ttyadapter

type Tty interface {
	IsOpen() bool
	Open(onResize func(int, int)) error
	GetKey() (string, error)
	Size() (int, int, error)
	Close() error
}
