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

package nodewallet

import (
	"fmt"

	"github.com/elysiumstation/fury/core/config"
	"github.com/elysiumstation/fury/core/nodewallets"
	vgfmt "github.com/elysiumstation/fury/libs/fmt"
	vgjson "github.com/elysiumstation/fury/libs/json"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"

	"github.com/jessevdk/go-flags"
)

type generateCmd struct {
	config.OutputFlag

	Config nodewallets.Config

	WalletPassphrase config.Passphrase `long:"wallet-passphrase-file"`

	Chain string `short:"c" long:"chain" required:"true" description:"The chain to be imported" choice:"fury" choice:"ethereum"`
	Force bool   `long:"force" description:"Should the command generate a new wallet on top of an existing one"`

	// clef options
	EthereumClefAddress string `long:"ethereum-clef-address" description:"The URL to the clef instance that Fury will use to generate a clef wallet."`
}

const (
	ethereumChain   = "ethereum"
	furyChain       = "fury"
	tendermintChain = "tendermint"
)

func (opts *generateCmd) Execute(_ []string) error {
	output, err := opts.GetOutput()
	if err != nil {
		return err
	}

	if output.IsHuman() && opts.EthereumClefAddress != "" {
		fmt.Println(yellow("Warning: Generating a new account in Clef has to be manually approved, and only the Key Store backend is supported. \nPlease consider using the 'import' command instead."))
	}

	log := logging.NewLoggerFromConfig(logging.NewDefaultConfig())
	defer log.AtExit()

	registryPass, err := rootCmd.PassphraseFile.Get("node wallet", false)
	if err != nil {
		return err
	}

	furyPaths := paths.New(rootCmd.FuryHome)

	_, conf, err := config.EnsureNodeConfig(furyPaths)
	if err != nil {
		return err
	}

	opts.Config = conf.NodeWallet

	if _, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse(); err != nil {
		return err
	}

	var data map[string]string
	switch opts.Chain {
	case ethereumChain:
		var walletPass string
		if opts.EthereumClefAddress == "" {
			walletPass, err = opts.WalletPassphrase.Get("blockchain wallet", true)
			if err != nil {
				return err
			}
		}

		data, err = nodewallets.GenerateEthereumWallet(
			furyPaths,
			registryPass,
			walletPass,
			opts.EthereumClefAddress,
			opts.Force,
		)
		if err != nil {
			return fmt.Errorf("couldn't generate Ethereum node wallet: %w", err)
		}
	case furyChain:
		walletPass, err := opts.WalletPassphrase.Get("blockchain wallet", true)
		if err != nil {
			return err
		}

		data, err = nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass, opts.Force)
		if err != nil {
			return fmt.Errorf("couldn't generate Fury node wallet: %w", err)
		}
	default:
		return fmt.Errorf("chain %q is not supported", opts.Chain)
	}

	if output.IsHuman() {
		fmt.Println(green("generation successful:"))
		vgfmt.PrettyPrint(data)
	} else if output.IsJSON() {
		if err := vgjson.Print(data); err != nil {
			return err
		}
	}

	return nil
}
