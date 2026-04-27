package winch10

import (
	"context"
	"sync"
	"time"
)

func notice(onResize func(int, int)) (func(), error) {
	w, h, err := getSize()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func(lastw, lasth int) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				wg.Done()
				return
			case <-ticker.C:
				w, h, err := getSize()
				if err == nil && (w != lastw || h != lasth) {
					onResize(w, h)
					lastw = w
					lasth = h
				}
			}
		}
	}(w, h)

	return func() {
		cancel()
		wg.Wait()
	}, nil
}
