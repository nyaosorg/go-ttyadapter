package winch8

import (
	"sync"

	"github.com/mattn/go-tty/v2"
)

func Notice(tt *tty.TTY, onResize func(w, h int)) (func(), error) {
	_lastw, _lasth, err := tt.Size()
	if err != nil || onResize == nil {
		return func() {}, err
	}
	var wg sync.WaitGroup

	ws := tt.SIGWINCH()
	wg.Add(1)
	go func(lastw, lasth int) {
		for wh := range ws {
			if lastw != wh.W || lasth != wh.H {
				onResize(wh.W, wh.H)
				lastw = wh.W
				lasth = wh.H
			}
		}
		wg.Done()
	}(_lastw, _lasth)

	return wg.Wait, nil
}
