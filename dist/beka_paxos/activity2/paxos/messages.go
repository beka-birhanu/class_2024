package paxos

type Prepare struct {
	ProposalNumber ProposalNumber `json:"proposal_number"`
}

type Promise struct {
	ProposalNumber ProposalNumber `json:"proposal_number"`
}

type Accept struct {
	ProposalNumber ProposalNumber `json:"proposal_number"`
	Value          interface{}    `json:"value"`
}

type Accepted struct {
	ProposalNumber ProposalNumber `json:"proposal_number"`
	Value          interface{}    `json:"value"`
}

type ProposalNumber struct {
	BallotNumber int    `json:"ballot_number"`
	ProposerID   string `json:"proposer_ID"`
}
