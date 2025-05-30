package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Waiting...")
			time.Sleep(300 * time.Millisecond)
		}
	}
}
