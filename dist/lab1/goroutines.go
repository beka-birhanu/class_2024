package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func printLetters() {
	for i := 'A'; i <= 'E'; i++ {
		fmt.Printf("%c\n", i)
		time.Sleep(1 * time.Second)
	}
}

func main5() {
	go printNumbers()

	printLetters()

	time.Sleep(6 * time.Second)
	fmt.Println("Main function finished")
}
