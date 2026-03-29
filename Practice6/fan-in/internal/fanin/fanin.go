package fanin

import (
	"context"
	"sync"
)

func FanIn(ctx context.Context, channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	wg.Add(len(channels))

	for _, ch := range channels {
		go func(c <-chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-c:
					if !ok {
						return
					}

					select {
					case out <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
