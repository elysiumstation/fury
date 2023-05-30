// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/elysiumstation/fury/core/config"
	"github.com/elysiumstation/fury/core/nodewallets"
	"github.com/elysiumstation/fury/core/txn"
	"github.com/elysiumstation/fury/core/validators"
	"github.com/elysiumstation/fury/libs/crypto"
	vgjson "github.com/elysiumstation/fury/libs/json"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"

	"github.com/jessevdk/go-flags"
)

type RotateEthKeyCmd struct {
	config.FuryHomeFlag
	config.OutputFlag
	config.Passphrase `long:"passphrase-file"`

	TargetBlock      uint64 `short:"b" long:"target-block" description:"The future block height at which the rotation will take place" `
	RotateFrom       string `short:"r" long:"rotate-from" description:"Ethereum address being rotated from" `
	SubmitterAddress string `short:"s" long:"submitter-address" description:"Ethereum address to use as a submitter to contract changes" `
}

var rotateEthKeyCmd RotateEthKeyCmd

func (opts *RotateEthKeyCmd) Execute(_ []string) error {
	log := logging.NewLoggerFromConfig(logging.NewDefaultConfig())
	defer log.AtExit()

	output, err := opts.GetOutput()
	if err != nil {
		return err
	}

	furyPaths := paths.New(opts.FuryHome)

	_, conf, err := config.EnsureNodeConfig(furyPaths)
	if err != nil {
		return err
	}

	if !conf.IsValidator() {
		return errors.New("node is not a validator")
	}

	registryPass, err := opts.Get("node wallet", false)
	if err != nil {
		return err
	}

	nodeWallets, err := nodewallets.GetNodeWallets(conf.NodeWallet, furyPaths, registryPass)
	if err != nil {
		return fmt.Errorf("couldn't get node wallets: %w", err)
	}

	// ensure the nodewallet is setup properly, if not we could not complete the command
	if err := nodeWallets.Verify(); err != nil {
		return fmt.Errorf("node wallet misconfigured: %w", err)
	}

	cmd := commandspb.EthereumKeyRotateSubmission{
		CurrentAddress:   crypto.EthereumChecksumAddress(opts.RotateFrom),
		NewAddress:       nodeWallets.Ethereum.PubKey().Hex(),
		TargetBlock:      opts.TargetBlock,
		SubmitterAddress: opts.SubmitterAddress,
	}

	if len(cmd.SubmitterAddress) != 0 {
		cmd.SubmitterAddress = crypto.EthereumChecksumAddress(cmd.SubmitterAddress)
	}

	if err := validators.SignEthereumKeyRotation(&cmd, nodeWallets.Ethereum); err != nil {
		return err
	}

	commander, _, cfunc, err := getNodeWalletCommander(log, registryPass, furyPaths)
	if err != nil {
		return fmt.Errorf("failed to get commander: %w", err)
	}
	defer cfunc()

	var txHash string
	ch := make(chan struct{})
	commander.Command(
		context.Background(),
		txn.RotateEthereumKeySubmissionCommand,
		&cmd,
		func(h string, e error) {
			txHash, err = h, e
			close(ch)
		}, nil)

	<-ch
	if err != nil {
		return err
	}

	if output.IsHuman() {
		fmt.Printf("ethereum key rotation successfully sent\ntxHash: %s\nethereum signature: %v\nRotating from: %s\nRotating to: %s",
			txHash,
			cmd.EthereumSignature.Value,
			opts.RotateFrom,
			cmd.NewAddress,
		)
	} else if output.IsJSON() {
		return vgjson.Print(struct {
			TxHash            string `json:"txHash"`
			EthereumSignature string `json:"ethereumSignature"`
			RotateFrom        string `json:"rotateFrom"`
			RotateTo          string `json:"rotateTo"`
		}{
			TxHash:            txHash,
			RotateFrom:        opts.RotateFrom,
			RotateTo:          cmd.NewAddress,
			EthereumSignature: cmd.EthereumSignature.Value,
		})
	}

	return nil
}

func RotateEthKey(ctx context.Context, parser *flags.Parser) error {
	announceNodeCmd = AnnounceNodeCmd{}

	var (
		short = "Send a transaction to rotate to current ethereum key"
		long  = "Send a transaction to rotate to current ethereum key"
	)
	_, err := parser.AddCommand("rotate_eth_key", short, long, &rotateEthKeyCmd)
	return err
}
