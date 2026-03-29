package counter

import (
	"fmt"
	"sync"
)

func RunUnsafe() {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("Unsafe initial Counter:", counter)
}

//The final value is not 1000 because counter++ is not atomic and causes a data race
//when multiple goroutines concurrently read and write the same memory location.
