package main

import (
	"fmt"
	"sync"
)

const countOfRead = 1_000_000
const countOfWrite = 10

func main() {
	total1 := useMutex()
	fmt.Println(total1)

	total2 := useRWMutex()
	fmt.Println(total2)
}

func useMutex() int {
	var mu sync.Mutex
	var wg sync.WaitGroup
	total := 0
	count := 1

	for i := 0; i < countOfRead; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			total += count
		}()
	}
	for i := 0; i < countOfWrite; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			count++
		}()
	}

	wg.Wait()
	return total
}

func useRWMutex() int {
	var rw sync.RWMutex
	var wg sync.WaitGroup
	total := 0
	count := 1

	for i := 0; i < countOfRead; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rw.RLock()
			defer rw.RUnlock()
			total += count

		}()
	}
	for i := 0; i < countOfWrite; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rw.Lock()
			defer rw.Unlock()
			count++
		}()
	}

	wg.Wait()
	return total
}
