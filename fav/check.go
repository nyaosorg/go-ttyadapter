//go:build run

package main

import (
	"fmt"

	"github.com/nyaosorg/go-ttyadapter/fav"
)

func main() {
	fmt.Printf("type %T\n", new(fav.Tty))
}
