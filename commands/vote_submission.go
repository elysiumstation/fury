package commands

import (
	types "github.com/elysiumstation/fury/protos/fury"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
)

func CheckVoteSubmission(cmd *commandspb.VoteSubmission) error {
	return checkVoteSubmission(cmd).ErrorOrNil()
}

func checkVoteSubmission(cmd *commandspb.VoteSubmission) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("vote_submission", ErrIsRequired)
	}

	if len(cmd.ProposalId) <= 0 {
		errs.AddForProperty("vote_submission.proposal_id", ErrIsRequired)
	} else if !IsFuryPubkey(cmd.ProposalId) {
		errs.AddForProperty("vote_submission.proposal_id", ErrShouldBeAValidFuryID)
	}

	if cmd.Value == types.Vote_VALUE_UNSPECIFIED {
		errs.AddForProperty("vote_submission.value", ErrIsRequired)
	}

	if _, ok := types.Vote_Value_name[int32(cmd.Value)]; !ok {
		errs.AddForProperty("vote_submission.value", ErrIsNotValid)
	}

	return errs
}
