package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

// List to hold active clients
var clients = make(map[net.Conn]bool)
var mu sync.Mutex

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is ready to assign tasks...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Add client to the list of active clients
		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		// Handle each client connection in a separate goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		// Remove the client from active clients list and close the connection
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	for {
		// Generate a task: a random number based on the current time
		task := time.Now().Unix() % 100

		// Send the task to the client
		fmt.Fprintf(conn, "%d\n", task)

		// Receive response from the client
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}
		fmt.Println("Received result from client:", response)

		// Simulate a 5-second interval between tasks
		time.Sleep(5 * time.Second)
	}
}
