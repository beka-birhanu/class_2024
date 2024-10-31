package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server on localhost at port 8080
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for {

		// Prompt the user to enter a message
		fmt.Print("Enter a message: ")
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		// Send the message to the server
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// Read the server's response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Print("Server response: " + response)
	}
}
