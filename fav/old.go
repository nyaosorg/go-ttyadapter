//go:build windows && !go1.21

package fav

import (
	"github.com/nyaosorg/go-ttyadapter/tty8pe"
)

type Tty = tty8pe.Tty
