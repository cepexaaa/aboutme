package main

import "fmt"

func main() {
	done := make(chan struct{})
	go func() {
		fmt.Println("Hello world")
		close(done)
	}()
	<-done
}
