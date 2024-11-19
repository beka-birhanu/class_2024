package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Args holds the arguments for multiplication
type Args struct {
	A, B int
}

// Calculator provides methods for arithmetic operations
type Calculator int

// Multiply multiplies two integers and returns the result
func (c *Calculator) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	// Register the Calculator service
	calc := new(Calculator)
	rpc.Register(calc)

	// Start listening for incoming RPC connections
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	fmt.Println("RPC server is listening on port 1234...")

	// Accept incoming connections and block to serve requests
	rpc.Accept(listener)
}
