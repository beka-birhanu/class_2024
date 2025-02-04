package paxos

import (
	"context"
	"fmt"
	"time"
)

type Proposer struct {
	proposalNumber    ProposalNumber
	prepareChan       chan<- Prepare
	promiseChan       <-chan Promise
	acceptChan        chan<- Accept
	acceptedChan      <-chan Accepted
	maxRetry          int
	numberOfAccepters int
}

// NewProposer creates and returns a new Proposer instance.
func NewProposer(
	proposerID string,
	numberOfAccepters,
	maxRetry int,
	prepareChan chan<- Prepare,
	promiseChan <-chan Promise,
	acceptChan chan<- Accept,
	acceptedChan <-chan Accepted) *Proposer {
	return &Proposer{
		proposalNumber:    ProposalNumber{BallotNumber: 0, ProposerID: proposerID},
		numberOfAccepters: numberOfAccepters,
		maxRetry:          maxRetry,
		prepareChan:       prepareChan,
		promiseChan:       promiseChan,
		acceptChan:        acceptChan,
		acceptedChan:      acceptedChan,
	}
}

func (p *Proposer) Propose(ctx context.Context, value interface{}, ballotNumber int) interface{} {
	p.proposalNumber.BallotNumber = ballotNumber
	for range p.maxRetry {
		p.proposalNumber.BallotNumber++
		prepare := Prepare{ProposalNumber: p.proposalNumber}
		p.prepareChan <- prepare

		promises := 0
		for promises <= p.numberOfAccepters/2 {
			select {
			case <-ctx.Done():
				fmt.Printf("\nError: Time out on propose with Prepare:%v\n", prepare)
				return nil
			case <-time.After(time.Millisecond * 300):
				fmt.Printf("\nInfo: Time out on propose with Prepare:%v  retrying...\n", prepare)
				break
			case promise := <-p.promiseChan:
				fmt.Println(promise)
				if promise.ProposalNumber.BallotNumber == p.proposalNumber.BallotNumber &&
					promise.ProposalNumber.ProposerID == p.proposalNumber.ProposerID {
					promises += 1
				}
			default:
				time.Sleep(time.Millisecond)
			}
		}
		if promises > p.numberOfAccepters/2 {
			break
		}
	}

	for range p.maxRetry {
		accept := Accept{Value: value, ProposalNumber: p.proposalNumber}
		p.acceptChan <- accept

		accepts := 0
		for accepts <= p.numberOfAccepters/2 {
			select {
			case <-ctx.Done():
				fmt.Printf("\nError: Time out on propose with Accept:%v\n", accept)
				return nil
			case <-time.After(time.Millisecond * 300):
				fmt.Printf("\nInfo: Time out on propose with Accept:%v  retrying...\n", accept)
				break
			case ack := <-p.acceptedChan:
				if ack.ProposalNumber.BallotNumber == p.proposalNumber.BallotNumber &&
					ack.ProposalNumber.ProposerID == p.proposalNumber.ProposerID {
					accepts++
				}
			default:
				time.Sleep(time.Millisecond)
			}

		}
		if accepts > p.numberOfAccepters/2 {
			return value
		}
	}

	return nil
}
