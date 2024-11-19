package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

func getResult(client *rpc.Client, operation string, args *Args) (string, error) {
	handleCall := func(operation string, args *Args) (string, error) {
		var result int
		var response string
		var err error

		// Call the selected operation method on the server asynchronously
		call := client.Go(fmt.Sprintf("Calculator.%s", operation), args, &result, nil)

		// Wait for the call to complete or timeout (after 5 seconds)
		select {
		case <-call.Done:
			// The call has finished
			if call.Error != nil {
				err = call.Error

			} else {
				// Print the result if no error occurred
				response = fmt.Sprintf("Result of %s(%d, %d) = %d\n", operation, args.A, args.B, result)
			}
		case <-time.After(5 * time.Second):
			// Timeout case: the RPC call took too long
			response = fmt.Sprintf("Timeout: %s operation took too long\n", operation)
		}

		return response, err
	}

	if _, err := handleCall(operation, args); err != nil {
		return "", nil
	}

	response, err := handleCall("GetLastResult", &Args{})
	return response, err
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

		// Validate the operation
		validOperations := map[string]bool{"Multiply": true, "Add": true, "Subtract": true, "Divide": true, "Power": true}
		if !validOperations[operation] {
			fmt.Println("Invalid operation. Please enter one of Multiply, Add, Subtract, Divide, Power.")
			continue
		}

		// Read input for arguments
		var a, b int
		fmt.Print("Enter two integers (A and B) separated by space: ")
		_, err := fmt.Scan(&a, &b)
		if err != nil {
			fmt.Println("Invalid input, please try again.")
			continue
		}

		result, err := getResult(client, operation, &Args{A: a, B: b})

		if err != nil {
			fmt.Printf("Error calling %s method: %s\n", operation, err)
			fmt.Println("Exiting...")
			os.Exit(0)
		} else {
			fmt.Println(result)
		}

	}
}
