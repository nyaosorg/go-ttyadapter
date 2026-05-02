//go:build !windows || go1.21

package fav

import (
	"github.com/nyaosorg/go-ttyadapter/tty10pe"
)

type Tty = tty10pe.Tty
