package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/beka-birhanu/assignment10/dto"
	"github.com/streadway/amqp"
)

func main() {
	// Initialize the paper server
	paperServer := &PaperServer{
		Papers: make(map[int]dto.Paper),
		Mu:     sync.Mutex{},
		NextID: 1,
	}

	// Establish RabbitMQ connection
	conn, err := amqp.Dial("amqp://test:test_pass@localhost:5672/")
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()
	paperServer.MQConn = conn

	// Initialize RabbitMQ fanout exchange
	err = paperServer.initializeRabbitMQ()
	if err != nil {
		fmt.Println("Failed to initialize RabbitMQ:", err)
		return
	}

	// Register the PaperServer service
	err = rpc.Register(paperServer)
	if err != nil {
		fmt.Println("Failed to register PaperServer:", err)
		return
	}

	// Start listening for incoming RPC connections
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("RPC server is listening on port 1234...")

	// Accept and handle connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
