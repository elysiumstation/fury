package commands

import (
	"encoding/hex"

	"github.com/elysiumstation/fury/libs/crypto"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
)

func CheckAnnounceNode(cmd *commandspb.AnnounceNode) error {
	return checkAnnounceNode(cmd).ErrorOrNil()
}

func checkAnnounceNode(cmd *commandspb.AnnounceNode) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("announce_node", ErrIsRequired)
	}

	if len(cmd.FuryPubKey) == 0 {
		errs.AddForProperty("announce_node.fury_pub_key", ErrIsRequired)
	} else if !IsFuryPubkey(cmd.FuryPubKey) {
		errs.AddForProperty("announce_node.fury_pub_key", ErrShouldBeAValidFuryPubkey)
	}

	if len(cmd.Id) == 0 {
		errs.AddForProperty("announce_node.id", ErrIsRequired)
	} else if !IsFuryPubkey(cmd.Id) {
		errs.AddForProperty("announce_node.id", ErrShouldBeAValidFuryPubkey)
	}

	if len(cmd.EthereumAddress) == 0 {
		errs.AddForProperty("announce_node.ethereum_address", ErrIsRequired)
	} else if !crypto.EthereumIsValidAddress(cmd.EthereumAddress) {
		errs.AddForProperty("announce_node.ethereum_address", ErrIsNotValidEthereumAddress)
	}

	if len(cmd.ChainPubKey) == 0 {
		errs.AddForProperty("announce_node.chain_pub_key", ErrIsRequired)
	}

	if cmd.EthereumSignature == nil || len(cmd.EthereumSignature.Value) == 0 {
		errs.AddForProperty("announce_node.ethereum_signature", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.EthereumSignature.Value)
		if err != nil {
			errs.AddForProperty("announce_node.ethereum_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if cmd.FurySignature == nil || len(cmd.FurySignature.Value) == 0 {
		errs.AddForProperty("announce_node.fury_signature", ErrIsRequired)
	} else {
		_, err := hex.DecodeString(cmd.FurySignature.Value)
		if err != nil {
			errs.AddForProperty("announce_node.fury_signature.value", ErrShouldBeHexEncoded)
		}
	}

	if len(cmd.SubmitterAddress) != 0 && !crypto.EthereumIsValidAddress(cmd.SubmitterAddress) {
		errs.AddForProperty("announce_node.submitter_address", ErrIsNotValidEthereumAddress)
	}

	return errs
}
