package main

import (
	"fmt"
)

func main() {
	ch, done := makeChan()

	fmt.Println("Start read from channel")
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println("Read from channel, value:", val)
		case <-done:
			fmt.Println("Done")
			return
		}
	}
}

func makeChan() (chan int, chan struct{}) {
	ch := make(chan int, 5)
	done := make(chan struct{})
	go func() {
		fmt.Println("Start write to channel")
		for i := 0; i < 5; i++ {
			ch <- i
			fmt.Println("Write to channel", i)
		}
		close(done)
		fmt.Println("Close channel")
	}()
	return ch, done
}
