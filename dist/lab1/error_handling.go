package main

import (
	"fmt"
	"os"
)

func readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("File opened successfully:", filename)
	return nil
}

func main() {
	err := readFile("test.txt")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("File read successfully.")
	}
}
