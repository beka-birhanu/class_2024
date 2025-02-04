package paxos

type Proposer struct {
	ProposalNumber int
}

func (p *Proposer) Propose(value interface{}, acceptors []*Acceptor) interface{} {
	promises := 0
	prepare := Prepare{ProposalNumber: p.ProposalNumber}
	for _, acceptor := range acceptors {
		promise := acceptor.HandlePrepare(prepare)
		if promise != nil && promise.ProposalNumber == p.ProposalNumber {
			promises += 1
		}
	}

	if promises > len(acceptors)/2 { // majored accpted
		accepted := 0
		accept := Accept{Value: value, ProposalNumber: p.ProposalNumber}
		for _, acceptors := range acceptors {
			ack := acceptors.HandleAccept(accept)
			if ack != nil && ack.ProposalNumber == p.ProposalNumber {
				accepted++
			}
		}

		if accepted > len(acceptors)/2 {
			return value
		}
	}

	return nil
}
