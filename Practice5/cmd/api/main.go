package main

import (
	"awesomeProject/internal/app"
	"fmt"
)

func main() {

	if err := app.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		return
	}
}
