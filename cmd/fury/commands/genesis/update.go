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

	"github.com/elysiumstation/fury/core/genesis"
	"github.com/elysiumstation/fury/core/nodewallets"
	vgtm "github.com/elysiumstation/fury/core/tendermint"
	"github.com/elysiumstation/fury/core/validators"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"

	"github.com/jessevdk/go-flags"
)

type updateCmd struct {
	Config nodewallets.Config

	DryRun bool   `long:"dry-run" description:"Display the genesis file without writing it"`
	TmHome string `short:"t" long:"tm-home" description:"The home path of tendermint"`
}

func (opts *updateCmd) Execute(_ []string) error {
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
	genesisState.Validators[b64TmPubKey] = validators.ValidatorData{
		ID:              walletID,
		FuryPubKey:      furyKey.value,
		FuryPubKeyIndex: furyKey.index,
		EthereumAddress: ethAddress,
		TmPubKey:        b64TmPubKey,
	}

	genesisDoc, _, err := tmConfig.Genesis()
	if err != nil {
		return fmt.Errorf("couldn't get genesis file: %w", err)
	}

	if err := vgtm.AddAppStateToGenesis(genesisDoc, &genesisState); err != nil {
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
