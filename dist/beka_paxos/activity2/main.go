package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/beka-birhanu/paxos-lab-activity2/paxos"
	"github.com/streadway/amqp"
)

func main() {
	// Load serverID from the environment and convert to integer
	serverID := os.Getenv("SERVER_ID")
	if serverID == "" {
		fmt.Println("SERVER_ID environment variable is not set")
		return
	}

	// Load numberOfAcceptor from the environment and convert to integer
	numberOfAcceptorStr := os.Getenv("NUMBER_OF_ACCEPTOR")
	if numberOfAcceptorStr == "" {
		fmt.Println("NUMBER_OF_ACCEPTOR environment variable is not set")
		return
	}

	numberOfAcceptor, err := strconv.Atoi(numberOfAcceptorStr)
	if err != nil {
		fmt.Printf("Invalid NUMBER_OF_ACCEPTOR value: %s\n", numberOfAcceptorStr)
		return
	}

	// Load RabbitMQ configuration from environment variables
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	if rabbitmqHost == "" {
		fmt.Println("RABBITMQ_HOST environment variable is not set")
		return
	}

	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	if rabbitmqPort == "" {
		fmt.Println("RABBITMQ_PORT environment variable is not set")
		return
	}

	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	if rabbitmqUser == "" {
		fmt.Println("RABBITMQ_USER environment variable is not set")
		return
	}

	rabbitmqPass := os.Getenv("RABBITMQ_PASS")
	if rabbitmqPass == "" {
		fmt.Println("RABBITMQ_PASS environment variable is not set")
		return
	}

	// Construct RabbitMQ connection URL
	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPass, rabbitmqHost, rabbitmqPort)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	// Create and start the Paxos server
	server := paxos.NewServer(conn, serverID, numberOfAcceptor)
	server.Serve()
}
