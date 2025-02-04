package paxos

import (
	"sync"
	"time"
)

type Acceptor struct {
	mu             sync.Mutex
	promisedNumber ProposalNumber
	acceptedValue  interface{}
	prepareChan    <-chan Prepare
	promiseChan    chan<- Promise
	acceptChan     <-chan Accept
	acceptedChan   chan<- Accepted
}

// NewAcceptor creates and initializes a new Acceptor with the provided channels
func NewAcceptor(
	prepareChan <-chan Prepare,
	promiseChan chan<- Promise,
	acceptChan <-chan Accept,
	acceptedChan chan<- Accepted,
) *Acceptor {
	return &Acceptor{
		promisedNumber: ProposalNumber{},
		acceptedValue:  nil,
		prepareChan:    prepareChan,
		promiseChan:    promiseChan,
		acceptChan:     acceptChan,
		acceptedChan:   acceptedChan,
	}
}

func (a *Acceptor) Start() {
	for {
		select {
		case p := <-a.prepareChan:
			a.mu.Lock()
			if p.ProposalNumber.BallotNumber > a.promisedNumber.BallotNumber {
				a.promisedNumber = p.ProposalNumber
				a.promiseChan <- Promise{ProposalNumber: a.promisedNumber}
			}
			a.mu.Unlock()

		case ac := <-a.acceptChan:
			a.mu.Lock()
			if ac.ProposalNumber.BallotNumber > a.promisedNumber.BallotNumber ||
				(ac.ProposalNumber.BallotNumber == a.promisedNumber.BallotNumber &&
					ac.ProposalNumber.ProposerID == a.promisedNumber.ProposerID) {
				a.promisedNumber = ac.ProposalNumber
				a.acceptedValue = ac.Value
				a.acceptedChan <- Accepted{ProposalNumber: a.promisedNumber, Value: a.acceptedValue}
			}
			a.mu.Unlock()

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (a *Acceptor) GetBallotNumber() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.promisedNumber.BallotNumber
}
