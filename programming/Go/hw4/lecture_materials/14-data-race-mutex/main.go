package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var count int

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count = count + 1
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
