package commands

import (
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
)

func CheckOrderCancellation(cmd *commandspb.OrderCancellation) error {
	return checkOrderCancellation(cmd).ErrorOrNil()
}

func checkOrderCancellation(cmd *commandspb.OrderCancellation) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("order_cancellation", ErrIsRequired)
	}

	if len(cmd.MarketId) > 0 && !IsFuryPubkey(cmd.MarketId) {
		errs.AddForProperty("order_cancellation.market_id", ErrShouldBeAValidFuryID)
	}

	if len(cmd.OrderId) > 0 && !IsFuryPubkey(cmd.OrderId) {
		errs.AddForProperty("order_cancellation.order_id", ErrShouldBeAValidFuryID)
	}

	return errs
}
