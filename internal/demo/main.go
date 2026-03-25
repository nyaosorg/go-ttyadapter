package demo

import (
	"fmt"
	"github.com/nyaosorg/go-ttyadapter"
)

func Exec(tty ttyadapter.Tty) error {
	err := tty.Open(func(w, h int) {
		fmt.Printf("Change size: %d, %d\r\n", w, h)
	})

	if err != nil {
		return err
	}
	defer tty.Close()

	w, h, err := tty.Size()
	if err != nil {
		return err
	}
	fmt.Printf("(%d,%d)\n", w, h)

	for {
		key, err := tty.GetKey()
		if err != nil || key == "\a" || key == "\x03" || key == "\x04" {
			return err
		}
		fmt.Printf("%q\n", key)
	}
	return nil
}
