package main

import "fmt"

func main() {
	ch := makeChan()

	fmt.Println("Start read from channel")
	for val := range ch {
		fmt.Println("Read from channel, value:", val)
	}
	fmt.Println("End")
}

func makeChan() chan int {
	ch := make(chan int, 5)
	go func() {
		fmt.Println("Start write to channel")
		for i := 0; i < 5; i++ {
			ch <- i
			fmt.Println("Write to channel", i)
		}
		fmt.Println("Close channel")
		close(ch)
	}()
	return ch
}
