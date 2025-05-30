package main

import (
	"fmt"
	"sync"
)

func main() {
	res := make(map[struct {
		a int
		b int
	}]int)
	for i := 0; i < 10_000_000; i++ {
		res1, res2 := memoryModel()
		res[struct {
			a int
			b int
		}{a: res1, b: res2}]++
	}
	fmt.Println(res)
}

func memoryModel() (int, int) {
	var x, y int
	var resultX, resultY int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		x = 1
		resultY = y
	}()
	go func() {
		defer wg.Done()
		y = 1
		resultX = x
	}()
	wg.Wait()
	return resultX, resultY
}
