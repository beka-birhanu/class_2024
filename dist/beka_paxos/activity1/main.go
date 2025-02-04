package main

import (
	"fmt"

	"github.com/beka-birhanu/paxos-lab/activity1/paxos"
)

func main() {
	acceptors := []*paxos.Acceptor{{}, {}, {}, {}, {}}

	proposer := paxos.Proposer{ProposalNumber: 1}
	value := proposer.Propose("Distributed Systems is cool!", acceptors)

	if value != nil {
		fmt.Printf("Consensus reached on value: %v\n", value)
	} else {
		fmt.Println("Consensus not reached.")
	}
}
