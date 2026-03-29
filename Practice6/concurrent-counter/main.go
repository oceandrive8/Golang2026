package main

import (
	"fmt"

	"concurrent-counter/counter"
)

func main() {
	//fmt.Println("Running unsafe version:")
	//counter.RunUnsafe()

	fmt.Println("\nRunning mutex version:")
	counter.RunMutex()

	fmt.Println("\nRunning atomic version:")
	counter.RunAtomic()
}
