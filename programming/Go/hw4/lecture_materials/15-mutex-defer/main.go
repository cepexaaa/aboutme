package main

import (
	"fmt"
	"sync"
)

var count int

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			if count >= 7 {
				return
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
