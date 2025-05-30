package main

import (
	"fmt"
)

func main() {
	fmt.Println("start")
	go hello()
	go world()
}

func hello() {
	fmt.Println("Hello")
}

func world() {
	fmt.Println("world")
}
