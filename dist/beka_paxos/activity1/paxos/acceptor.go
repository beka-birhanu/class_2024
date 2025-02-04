package paxos

import "sync"

type Acceptor struct {
	mu             sync.Mutex
	promisedNumber int
	acceptedNumber int
	acceptedValue  interface{}
}

func (a *Acceptor) HandlePrepare(p Prepare) *Promise {
	a.mu.Lock()
	defer a.mu.Unlock()

	if p.ProposalNumber > a.promisedNumber {
		a.promisedNumber = p.ProposalNumber
		return &Promise{ProposalNumber: a.promisedNumber, AcceptedValue: a.acceptedNumber}
	}

	return nil
}

func (a *Acceptor) HandleAccept(ac Accept) *Accepted {
	a.mu.Lock()
	defer a.mu.Unlock()

	if ac.ProposalNumber >= a.promisedNumber {
		a.promisedNumber = ac.ProposalNumber
		a.acceptedNumber = ac.ProposalNumber
		a.acceptedValue = ac.Value
		return &Accepted{ProposalNumber: a.promisedNumber, Value: a.acceptedValue}
	}
	return nil
}
