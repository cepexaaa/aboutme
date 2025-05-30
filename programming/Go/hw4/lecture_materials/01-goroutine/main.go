package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start")
	go hello()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("End")
}

func hello() {
	fmt.Println("Hello world")
}
