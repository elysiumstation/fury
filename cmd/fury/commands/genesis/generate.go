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

package genesis

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/elysiumstation/fury/core/genesis"
	"github.com/elysiumstation/fury/core/nodewallets"
	vgtm "github.com/elysiumstation/fury/core/tendermint"
	"github.com/elysiumstation/fury/core/validators"
	vgrand "github.com/elysiumstation/fury/libs/rand"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"

	"github.com/jessevdk/go-flags"
	tmtypes "github.com/tendermint/tendermint/types"
)

type generateCmd struct {
	Config nodewallets.Config

	DryRun bool   `long:"dry-run" description:"Display the genesis file without writing it"`
	TmHome string `short:"t" long:"tm-home" description:"The home path of tendermint"`
}

func (opts *generateCmd) Execute(_ []string) error {
	log := logging.NewLoggerFromConfig(
		logging.NewDefaultConfig(),
	)
	defer log.AtExit()

	pass, err := genesisCmd.PassphraseFile.Get("node wallet", false)
	if err != nil {
		return err
	}

	furyPaths := paths.New(genesisCmd.FuryHome)

	if _, err := flags.NewParser(opts, flags.Default|flags.IgnoreUnknown).Parse(); err != nil {
		return err
	}

	furyKey, ethAddress, walletID, err := loadNodeWalletPubKey(opts.Config, furyPaths, pass)
	if err != nil {
		return err
	}

	tmConfig := vgtm.NewConfig(opts.TmHome)

	pubKey, err := tmConfig.PublicValidatorKey()
	if err != nil {
		return err
	}

	b64TmPubKey := base64.StdEncoding.EncodeToString(pubKey.Bytes())
	genesisState := genesis.DefaultState()
	genesisState.Validators[base64.StdEncoding.EncodeToString(pubKey.Bytes())] = validators.ValidatorData{
		ID:              walletID,
		FuryPubKey:      furyKey.value,
		FuryPubKeyIndex: furyKey.index,
		EthereumAddress: ethAddress,
		TmPubKey:        b64TmPubKey,
	}

	genesisDoc := &tmtypes.GenesisDoc{
		ChainID:         fmt.Sprintf("test-chain-%v", vgrand.RandomStr(6)),
		GenesisTime:     time.Now().Round(0).UTC(),
		ConsensusParams: tmtypes.DefaultConsensusParams(),
		Validators: []tmtypes.GenesisValidator{
			{
				Address: pubKey.Address(),
				PubKey:  pubKey,
				Power:   10,
			},
		},
	}

	if err = vgtm.AddAppStateToGenesis(genesisDoc, &genesisState); err != nil {
		return fmt.Errorf("couldn't add app_state to genesis: %w", err)
	}

	if !opts.DryRun {
		if err := tmConfig.SaveGenesis(genesisDoc); err != nil {
			return fmt.Errorf("couldn't save genesis: %w", err)
		}
	}

	prettifiedDoc, err := vgtm.Prettify(genesisDoc)
	if err != nil {
		return err
	}
	fmt.Println(prettifiedDoc)
	return nil
}

type furyPubKey struct {
	index uint32
	value string
}

func loadNodeWalletPubKey(config nodewallets.Config, furyPaths paths.Paths, registryPass string) (furyKey *furyPubKey, ethAddr, walletID string, err error) {
	nw, err := nodewallets.GetNodeWallets(config, furyPaths, registryPass)
	if err != nil {
		return nil, "", "", fmt.Errorf("couldn't get node wallets: %w", err)
	}

	if err := nw.Verify(); err != nil {
		return nil, "", "", err
	}

	furyPubKey := &furyPubKey{
		index: nw.Fury.Index(),
		value: nw.Fury.PubKey().Hex(),
	}

	return furyPubKey, nw.Ethereum.PubKey().Hex(), nw.Fury.ID().Hex(), nil
}
