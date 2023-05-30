package v1

import types "github.com/elysiumstation/fury/protos/fury"

func ProposalSubmissionFromProposal(p *types.Proposal) *ProposalSubmission {
	return &ProposalSubmission{
		Reference: p.Reference,
		Terms:     p.Terms,
	}
}
