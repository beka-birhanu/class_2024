package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

func main() {
	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}
	defer client.Close()

	fmt.Println("Connected to RPC server. Type 'exit' to quit.")

	for {
		// Read operation type from user
		fmt.Print("Enter operation (Multiply, Add, Subtract, Divide, Power) or 'exit' to quit: ")
		var operation string
		fmt.Scan(&operation)

		// Exit condition
		if operation == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// Read input for arguments
		var a, b int
		fmt.Print("Enter two integers (A and B) separated by space: ")
		_, err := fmt.Scan(&a, &b)
		if err != nil {
			fmt.Println("Invalid input, please try again.")
			continue
		}

		// Define the arguments
		args := Args{A: a, B: b}
		var result int

		// Call the selected operation method on the server
		err = client.Call(fmt.Sprintf("Calculator.%s", operation), args, &result)
		if err != nil {
			fmt.Printf("Error calling %s method: %s\n", operation, err)
			continue
		}

		fmt.Printf("Result of %s(%d, %d) = %d\n", operation, args.A, args.B, result)
	}
}
