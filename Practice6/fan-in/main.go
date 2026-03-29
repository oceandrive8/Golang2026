package main

import (
	"context"
	"fmt"
	"time"

	"fan-in/internal/fanin"
	"fan-in/internal/server"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := server.StartServer(ctx, "Alpha")
	ch2 := server.StartServer(ctx, "Beta")
	ch3 := server.StartServer(ctx, "Gamma")
	merged := fanin.FanIn(ctx, ch1, ch2, ch3)

	for val := range merged {
		fmt.Println(val)
	}
}
