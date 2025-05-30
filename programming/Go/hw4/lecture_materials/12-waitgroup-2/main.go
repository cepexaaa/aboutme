package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("start")
	go hello(&wg)
	go world(&wg)
	wg.Wait()
}

func hello(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello")
}

func world(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("world")
}
