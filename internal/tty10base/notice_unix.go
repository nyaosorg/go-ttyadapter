//go:build !windows

package tty10base

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"golang.org/x/sys/unix"
)

func notice(onResize func(int, int)) (func(), error) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, unix.SIGWINCH)

	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				signal.Stop(ch)
				close(ch)
				wg.Done()
				return
			case <-ch:
				w, h, err := getSize()
				if err == nil {
					onResize(w, h)
				}
			}
		}
	}()

	return func() {
		cancel()
		wg.Wait()
	}, nil
}
