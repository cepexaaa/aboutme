package main

import (
	"fmt"
)

func main() {
	ch := make(chan string) // создаем канал
	go func() {
		ch <- "Hello world" // пишем в канал
		close(ch)           // закрываем канал
	}()
	res := <-ch // читаем из канала
	fmt.Println(res)
}
