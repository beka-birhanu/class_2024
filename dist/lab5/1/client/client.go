package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// Args holds the arguments for multiplication
type Args struct {
	A, B int
}

func main() {
	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}

	// Define the arguments for multiplication
	args := Args{A: 3, B: 5}
	var result int

	// Call the Multiply method on the server
	err = client.Call("Calculator.Multiply", args, &result)
	if err != nil {
		log.Fatal("Error calling Multiply method:", err)
	}

	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, result)

	client.Close()
}
