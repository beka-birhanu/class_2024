package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	client, err := NewPaperClient("localhost:1234")
	if err != nil {
		fmt.Printf("Error initializing client: %v\n", err)
		return
	}
	defer client.Close()

	go client.Run()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("PaperClient is running. Type 'quit' to exit.")

	for {
		fmt.Print("Enter command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "quit" {
			fmt.Println("Exiting PaperClient.")
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "add":
			if len(parts) < 4 {
				fmt.Println("Usage: add <Author> <Title> <FilePath>")
				continue
			}
			client.AddPaper(parts[1], parts[2], parts[3])

		case "list":
			client.ListPapers()

		case "details":
			if len(parts) != 2 {
				fmt.Println("Usage: details <PaperNumber>")
				continue
			}
			paperNumber := atoi(parts[1])
			client.GetPaperDetails(paperNumber)

		case "fetch":
			if len(parts) != 2 {
				fmt.Println("Usage: fetch <PaperNumber>")
				continue
			}
			paperNumber := atoi(parts[1])
			client.FetchPaperContent(paperNumber)

		default:
			fmt.Println("Invalid command. Use 'add', 'list', 'details', 'fetch', or 'quit'.")
		}
	}
}

func atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
