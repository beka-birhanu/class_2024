package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server on localhost at port 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start a goroutine to listen for incoming messages from the server
	go receiveMessages(conn)

	// Read messages from stdin and send them to the server
	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(conn, message)
	}
}

// Function to receive messages from the server
func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Print("Message from server: ", message)
	}
}
