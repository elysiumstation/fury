package commands

import (
	"encoding/hex"

	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
)

func CheckValidatorHeartbeat(cmd *commandspb.ValidatorHeartbeat) error {
	return checkValidatorHeartbeat(cmd).ErrorOrNil()
}

func checkValidatorHeartbeat(cmd *commandspb.ValidatorHeartbeat) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("validator_heartbeat", ErrIsRequired)
	}

	if len(cmd.NodeId) == 0 {
		errs.AddForProperty("validator_heartbeat.node_id", ErrIsRequired)
	} else {
		if !IsFuryPubkey(cmd.NodeId) {
			errs.AddForProperty("validator_heartbeat.node_id", ErrShouldBeAValidFuryPubkey)
		}
	}

	if cmd.EthereumSignature == nil || len(cmd.EthereumSignature.Value) == 0 {
		errs.AddForProperty("validator_heartbeat.ethereum_signature.value", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.EthereumSignature.Value)
		if err != nil {
			errs.AddForProperty("validator_heartbeat.ethereum_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if cmd.FurySignature == nil || len(cmd.FurySignature.Value) == 0 {
		errs.AddForProperty("validator_heartbeat.fury_signature.value", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.FurySignature.Value)
		if err != nil {
			errs.AddForProperty("validator_heartbeat.fury_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if len(cmd.FurySignature.Algo) == 0 {
		errs.AddForProperty("validator_heartbeat.fury_signature.algo", ErrIsRequired)
	}

	return errs
}
