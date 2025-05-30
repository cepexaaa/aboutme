package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		fmt.Println("Start write to channel")
		for i := 0; i < 5; i++ {
			ch <- i
			fmt.Println("Write to channel", i)
		}
		fmt.Println("Close channel")
		close(ch) // Идиома go: закрывает канал та горутина, которая пишет в него
	}()

	fmt.Println("Start read from channel")
	for i := 0; i < 5; i++ {
		val := <-ch
		fmt.Println("Read from channel, value:", val)
	}
	fmt.Println("End")
}
