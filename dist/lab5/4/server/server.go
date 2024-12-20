package main

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/rpc"
	"sync"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

// Calculator provides methods for arithmetic operations
type Calculator struct {
	lastResult int
	mu         sync.Mutex
}

// Multiply multiplies two integers and returns the result
func (c *Calculator) Multiply(args *Args, result *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = args.A * args.B
	return nil
}

// Add adds two integers and returns the result
func (c *Calculator) Add(args *Args, result *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = args.A + args.B
	return nil
}

// Subtract subtracts the second integer from the first and returns the result
func (c *Calculator) Subtract(args *Args, result *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = args.A - args.B
	return nil
}

// Divide divides the first integer by the second and returns the result
// Returns an error if division by zero is attempted
func (c *Calculator) Divide(args *Args, result *int) error {
	if args.B == 0 {
		return errors.New("division by zero is not allowed")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = args.A / args.B
	return nil
}

// Power raises the first integer to the power of the second and returns the result
func (c *Calculator) Power(args *Args, result *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastResult = int(math.Pow(float64(args.A), float64(args.B)))
	return nil
}

// GetLastResult returns the last result
func (c *Calculator) GetLastResult(args *Args, result *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	*result = c.lastResult
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

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go rpc.ServeConn(conn)
	}
}