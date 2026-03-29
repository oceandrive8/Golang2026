package main

import (
	"fmt"

	"thread-safe-map/mutexmap"
	"thread-safe-map/syncmap"
)

func main() {
	fmt.Println("Running sync.Map version:")
	syncValue := syncmap.Run()
	fmt.Printf("Final Value (sync.Map): %d\n\n", syncValue)

	fmt.Println("Running RWMutex map version:")
	mutexValue := mutexmap.Run()
	fmt.Printf("Final Value (RWMutex): %d\n", mutexValue)
}
