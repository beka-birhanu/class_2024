package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	"github.com/beka-birhanu/assignment10/dto"
	"github.com/streadway/amqp"
)

type PaperClient struct {
	serverAddress string
	rpcClient     *rpc.Client
	mqConn        *amqp.Connection
	wg            sync.WaitGroup
}

func NewPaperClient(serverAddress string) (*PaperClient, error) {
	// Connect to the RPC server
	rpcClient, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC server: %v", err)
	}

	// Connect to RabbitMQ
	mqConn, err := amqp.Dial("amqp://test:test_pass@localhost:5672/")
	if err != nil {
		rpcClient.Close()
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	return &PaperClient{
		serverAddress: serverAddress,
		rpcClient:     rpcClient,
		mqConn:        mqConn,
	}, nil
}

func (c *PaperClient) Close() {
	fmt.Println("Closing client connections...")
	c.rpcClient.Close()
	c.mqConn.Close()
}

func (c *PaperClient) SubscribeToNewPapers() {
	ch, err := c.mqConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Create a temporary queue for this client
	q, err := ch.QueueDeclare(
		"",    // name (empty means generate a unique name)
		false, // durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Bind the queue to the fanout exchange
	err = ch.QueueBind(
		q.Name,            // queue name
		"",                // routing key
		"new_paper_added", // exchange name
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag
		true,   // auto-acknowledge
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	fmt.Println("\nListening for new paper notifications...")

	for msg := range msgs {
		fmt.Printf("New Paper Notification: %s\n", msg.Body)
	}
}

func (c *PaperClient) AddPaper(author, title, filePath string) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	// Determine the file format
	ext := strings.ToLower(filepath.Ext(filePath))
	var format string
	switch ext {
	case ".pdf":
		format = "PDF"
	case ".doc", ".docx":
		format = "DOC"
	default:
		fmt.Println("Unsupported file format. Only PDF and DOC are supported.")
		return
	}

	args := dto.AddPaperArgs{
		Author:  author,
		Title:   title,
		Format:  format,
		Content: content,
	}
	reply := dto.AddPaperReply{}

	err = c.rpcClient.Call("PaperServer.AddPaper", args, &reply)
	if err != nil {
		fmt.Printf("Error adding paper: %v\n", err)
		return
	}

	fmt.Printf("Paper added successfully with ID %d\n", reply.PaperNumber)
}

func (c *PaperClient) ListPapers() {
	args := struct{}{}
	reply := dto.ListPapersReply{}

	err := c.rpcClient.Call("PaperServer.ListPapers", args, &reply)
	if err != nil {
		fmt.Printf("Error listing papers: %v\n", err)
		return
	}

	if len(reply.Papers) == 0 {
		fmt.Println("No papers found.")
		return
	}

	fmt.Println("Papers:")
	for _, paper := range reply.Papers {
		fmt.Printf("  ID: %d | Author: %s | Title: %s\n", paper.Number, paper.Author, paper.Title)
	}
}

func (c *PaperClient) GetPaperDetails(paperNumber int) {
	args := dto.GetPaperArgs{Number: paperNumber}
	reply := dto.GetPaperDetailsReply{}

	err := c.rpcClient.Call("PaperServer.GetPaperDetails", args, &reply)
	if err != nil {
		fmt.Printf("Error getting paper details: %v\n", err)
		return
	}

	fmt.Printf("Paper Details - Author: %s, Title: %s\n", reply.Author, reply.Title)
}

func (c *PaperClient) FetchPaperContent(paperNumber int) {
	args := dto.FetchPaperArgs{Number: paperNumber}
	reply := dto.FetchPaperReply{}

	err := c.rpcClient.Call("PaperServer.FetchPaperContent", args, &reply)
	if err != nil {
		fmt.Printf("Error fetching paper content: %v\n", err)
		return
	}

	fmt.Printf("Paper Content :\n%s\n", string(reply.Content))
}

func (c *PaperClient) Run() {
	// Capture termination signals to clean up gracefully
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)

	// Start listening for notifications
	go c.SubscribeToNewPapers()

	<-terminate
	fmt.Println("\nShutting down client...")
	c.wg.Wait()
}
