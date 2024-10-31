package main

import "fmt"

func sendData(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	fmt.Println("program finished.")
}

func main6() {
	ch := make(chan int)
	go sendData(ch)

	for val := range ch {
		fmt.Println("Recieved:", val)
	}

	fmt.Println("program finished.")
}
