package commands

import (
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
)

func CheckCancelTransfer(cmd *commandspb.CancelTransfer) error {
	return checkCancelTransfer(cmd).ErrorOrNil()
}

func checkCancelTransfer(cmd *commandspb.CancelTransfer) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("cancel_transfer", ErrIsRequired)
	}

	if len(cmd.TransferId) <= 0 {
		errs.AddForProperty("cancel_transfer.transfer_id", ErrIsRequired)
	} else if !IsFuryPubkey(cmd.TransferId) {
		errs.AddForProperty("cancel_transfer.transfer_id", ErrShouldBeAValidFuryID)
	}

	return errs
}
