package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	// Connect to the server on localhost at port 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to the server.")

	for {
		// Receive task (number) from server
		fmt.Println("Waiting for task from server...")
		task, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving task from server:", err)
			return
		}
		task = strings.TrimSpace(task)
		fmt.Printf("Received task: %s\n", task)

		// Convert task to integer
		num, err := strconv.Atoi(task)
		if err != nil {
			fmt.Println("Error converting task to integer:", err)
			continue
		}

		// Calculate result
		result := num * num
		fmt.Printf("Calculated result (square): %d\n", result)

		// Send result back to server
		fmt.Printf("Sending result back to server: %d\n", result)
		_, err = fmt.Fprintf(conn, "%d\n", result)
		if err != nil {
			fmt.Println("Error sending result to server:", err)
			return
		}
	}
}
