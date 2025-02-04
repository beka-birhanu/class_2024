package paxos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	acceptorQueueKey    = "FOR_ACCEPTORS"
	proposerQueueKey    = "FOR_PROPOSERS"
	prepareMessageType  = "PREPARE"
	promiseMessageType  = "PROMISE"
	acceptMessageType   = "ACCEPT"
	acceptedMessageType = "ACCEPTED"
)

type Server struct {
	acceptor             *Acceptor
	proposer             *Proposer
	mqConn               *amqp.Connection
	proposing            bool
	acceptorPrepareChan  chan Prepare
	acceptorPromiseChan  chan Promise
	acceptorAcceptChan   chan Accept
	acceptorAcceptedChan chan Accepted
	proposerPrepareChan  chan Prepare
	proposerPromiseChan  chan Promise
	proposerAcceptChan   chan Accept
	proposerAcceptedChan chan Accepted
	mu                   sync.RWMutex
}

func NewServer(mqConn *amqp.Connection, serverID string, numberOfAccepters int) *Server {
	log.Println("Initializing server...")
	server := &Server{mqConn: mqConn,
		acceptorPrepareChan:  make(chan Prepare),
		acceptorPromiseChan:  make(chan Promise),
		acceptorAcceptChan:   make(chan Accept),
		acceptorAcceptedChan: make(chan Accepted),
		proposerPrepareChan:  make(chan Prepare),
		proposerPromiseChan:  make(chan Promise),
		proposerAcceptChan:   make(chan Accept),
		proposerAcceptedChan: make(chan Accepted),
	}

	server.acceptor = NewAcceptor(
		server.acceptorPrepareChan,
		server.acceptorPromiseChan,
		server.acceptorAcceptChan,
		server.acceptorAcceptedChan,
	)
	server.proposer = NewProposer(
		serverID,
		numberOfAccepters,
		3,
		server.proposerPrepareChan,
		server.proposerPromiseChan,
		server.proposerAcceptChan,
		server.proposerAcceptedChan,
	)

	log.Println("Server initialized.")
	return server
}

func (s *Server) Serve() {
	log.Println("Starting server...")
	http.HandleFunc("/porpose", s.proposeHandler)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Error: while starting HTTP server: %s", err)
		}
		log.Println("HTTP server listening on :8080")
	}()

	ch, err := s.mqConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	err = createFanout(ch, proposerQueueKey)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ proposer fanout: %v", err)
		return
	}
	log.Println("Proposer fanout exchange created.")

	proposerQueue, err := getMessageQueue(ch, proposerQueueKey)
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ proposer queue: %v", err)
		return
	}
	log.Println("Proposer queue initialized.")

	err = createFanout(ch, acceptorQueueKey)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ acceptor fanout: %v", err)
		return
	}
	log.Println("Acceptor fanout exchange created.")

	acceptorQueue, err := getMessageQueue(ch, acceptorQueueKey)
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ acceptor queue: %v", err)
		return
	}

	log.Println("Acceptor queue initialized.")
	go s.acceptor.Start()

	for {
		select {
		case prepare := <-s.proposerPrepareChan:
			log.Printf("Publishing PREPARE message: %+v", prepare)
			body, err := json.Marshal(prepare)
			if err != nil {
				log.Printf("Error marshaling PREPARE message: %v", err)
				continue
			}
			message := QueueMessage{
				Type: prepareMessageType,
				Body: body,
			}
			publish(ch, acceptorQueueKey, message)

		case promise := <-s.acceptorPromiseChan:
			log.Printf("Publishing PROMISE message: %+v", promise)
			body, err := json.Marshal(promise)
			if err != nil {
				log.Printf("Error marshaling PROMISE message: %v", err)
				continue
			}
			message := QueueMessage{
				Type: promiseMessageType,
				Body: body,
			}
			publish(ch, proposerQueueKey, message)

		case accept := <-s.proposerAcceptChan:
			log.Printf("Publishing ACCEPT message: %+v", accept)
			body, err := json.Marshal(accept)
			if err != nil {
				log.Printf("Error marshaling ACCEPT message: %v", err)
				continue
			}
			message := QueueMessage{
				Type: acceptMessageType,
				Body: body,
			}
			publish(ch, acceptorQueueKey, message)

		case accepted := <-s.acceptorAcceptedChan:
			log.Printf("Publishing ACCEPTED message: %+v", accepted)
			body, err := json.Marshal(accepted)
			if err != nil {
				log.Printf("Error marshaling ACCEPTED message: %v", err)
				continue
			}
			message := QueueMessage{
				Type: acceptedMessageType,
				Body: body,
			}
			publish(ch, proposerQueueKey, message)

		case messageForProposer := <-proposerQueue:
			s.mu.RLock()
			if !s.proposing {
				s.mu.RUnlock()
				log.Println("Skipping message for proposer, not proposing.")
				continue
			}
			log.Println("Handling message for proposer.")
			s.handleMessageForProposer(messageForProposer.Body)
			s.mu.RUnlock()

		case messageForAcceptor := <-acceptorQueue:
			log.Println("Handling message for acceptor.")
			s.handleMessageForAcceptor(messageForAcceptor.Body)
		}
	}
}

func (s *Server) handleMessageForProposer(msgByte []byte) {
	log.Println("Processing message for proposer...")
	var message QueueMessage
	if err := json.Unmarshal(msgByte, &message); err != nil {
		log.Printf("Error: Failed to unmarshal message: %s", err)
		return
	}

	switch message.Type {
	case promiseMessageType:
		var promise Promise
		if err := json.Unmarshal(message.Body, &promise); err != nil {
			log.Printf("Error: Failed to unmarshal Promise: %s", err)
			return
		}
		log.Printf("Received PROMISE message: %+v", promise)
		s.proposerPromiseChan <- promise

	case acceptedMessageType:
		var accepted Accepted
		if err := json.Unmarshal(message.Body, &accepted); err != nil {
			log.Printf("Error: Failed to unmarshal Accepted: %s", err)
			return
		}
		log.Printf("Received ACCEPTED message: %+v", accepted)
		s.proposerAcceptedChan <- accepted

	default:
		log.Printf("Info: Unknown message type received: %s", message.Type)
	}
}

func (s *Server) handleMessageForAcceptor(msgByte []byte) {
	log.Println("Processing message for acceptor...")
	var message QueueMessage
	if err := json.Unmarshal(msgByte, &message); err != nil {
		log.Printf("Error: Failed to unmarshal message: %s", err)
		return
	}

	log.Printf("Received message for acceptor: %+v", message)
	switch message.Type {
	case prepareMessageType:
		var prepare Prepare
		if err := json.Unmarshal(message.Body, &prepare); err != nil {
			log.Printf("Error: Failed to unmarshal Prepare: %s", err)
			return
		}
		log.Printf("Processing PREPARE message: %+v", prepare)
		s.acceptorPrepareChan <- prepare

	case acceptMessageType:
		var accept Accept
		if err := json.Unmarshal(message.Body, &accept); err != nil {
			log.Printf("Error: Failed to unmarshal Accept: %s", err)
			return
		}
		log.Printf("Processing ACCEPT message: %+v", accept)
		s.acceptorAcceptChan <- accept

	default:
		log.Printf("Info: Unknown message type received: %s", message.Type)
	}
}

func publish(ch *amqp.Channel, queueKey string, message QueueMessage) {
	log.Printf("Publishing message to %s: %+v", queueKey, message)
	messageBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error: Failed to marshal message: %v", err)
		return
	}

	err = ch.Publish(
		queueKey,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Printf("Message published to %s", queueKey)
	}
}

func createFanout(ch *amqp.Channel, exchangeName string) error {
	log.Printf("Creating fanout exchange: %s", exchangeName)
	err := ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare exchange %s: %v", exchangeName, err)
	}
	return err
}

func getMessageQueue(ch *amqp.Channel, queueKey string) (<-chan amqp.Delivery, error) {
	log.Printf("Initializing queue for key: %s", queueKey)
	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		queueKey,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to bind queue: %v", err)
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consuming messages: %v", err)
		return nil, err
	}
	return msgs, nil
}

func (s *Server) proposeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received propose request.")
	var body struct {
		Message string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Error: while decoding propose request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	log.Printf("Proposing value: %s", body.Message)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	s.mu.Lock()
	s.proposing = true
	s.mu.Unlock()
	defer func(cancel func()) {
		cancel()
		s.mu.Lock()
		s.proposing = false
		s.mu.Unlock()
		log.Println("Proposing completed.")
	}(cancel)

	value := s.proposer.Propose(ctx, body.Message, s.acceptor.GetBallotNumber())
	if value != nil {
		log.Printf("Consensus reached: %v", value)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consensus reached: %v", value)
	} else {
		log.Println("Consensus not reached.")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Consensus not reached")
	}
}
