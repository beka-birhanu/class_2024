package main

import (
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/beka-birhanu/assignment10/dto"
	"github.com/streadway/amqp"
)

type PaperServer struct {
	Papers map[int]dto.Paper // Map for storing papers
	Mu     sync.Mutex        // Mutex for thread-safe operations
	NextID int               // Counter for generating unique IDs
	MQConn *amqp.Connection  // RabbitMQ connection
}

// Initialize RabbitMQ and declare a fanout exchange
func (s *PaperServer) initializeRabbitMQ() error {
	ch, err := s.MQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Declare a fanout exchange for broadcasting paper updates
	err = ch.ExchangeDeclare(
		"new_paper_added", // Exchange name
		"fanout",          // Exchange type
		true,              // Durable
		false,             // Auto-deleted
		false,             // Internal
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	return nil
}

// AddPaper stores a new paper and publishes a message to the RabbitMQ exchange
func (s *PaperServer) AddPaper(args dto.AddPaperArgs, reply *dto.AddPaperReply) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Create and store the new paper
	paper := dto.Paper{
		Number:  s.NextID,
		Author:  args.Author,
		Title:   args.Title,
		Format:  args.Format,
		Content: args.Content,
	}
	s.Papers[s.NextID] = paper
	s.NextID++

	reply.PaperNumber = paper.Number

	// Publish to RabbitMQ exchange
	ch, err := s.MQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	message := fmt.Sprintf("New paper added: %s by %s", paper.Title, paper.Author)
	err = ch.Publish(
		"new_paper_added", // Exchange name
		"",                // Routing key (ignored for fanout exchanges)
		false,             // Mandatory
		false,             // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %v", err)
	}

	return nil
}

// ListPapers returns a list of all stored papers
func (s *PaperServer) ListPapers(args struct{}, reply *dto.ListPapersReply) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	reply.Papers = slices.Collect(maps.Values(s.Papers))
	return nil
}

// GetPaperDetails returns the author and title of a specific paper
func (s *PaperServer) GetPaperDetails(args dto.GetPaperArgs, reply *dto.GetPaperDetailsReply) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	paper, exists := s.Papers[args.Number]
	if !exists {
		return fmt.Errorf("paper with number %d not found", args.Number)
	}

	reply.Author = paper.Author
	reply.Title = paper.Title
	return nil
}

// FetchPaperContent retrieves the full content of a specific paper
func (s *PaperServer) FetchPaperContent(args dto.FetchPaperArgs, reply *dto.FetchPaperReply) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	paper, exists := s.Papers[args.Number]
	if !exists {
		return fmt.Errorf("paper with number %d not found", args.Number)
	}

	reply.Content = paper.Content
	return nil
}
