// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.FURY file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package genesis

import (
	"github.com/elysiumstation/fury/core/assets"
	"github.com/elysiumstation/fury/core/blockchain/abci"
	"github.com/elysiumstation/fury/core/checkpoint"
	"github.com/elysiumstation/fury/core/limits"
	"github.com/elysiumstation/fury/core/netparams"
	"github.com/elysiumstation/fury/core/validators"
)

type State struct {
	Assets             assets.GenesisState             `json:"assets"`
	Validators         validators.GenesisState         `json:"validators"`
	Network            abci.GenesisState               `json:"network"`
	NetParams          netparams.GenesisState          `json:"network_parameters"`
	NetParamsOverwrite netparams.GenesisStateOverwrite `json:"network_parameters_checkpoint_overwrite"`
	Limits             limits.GenesisState             `json:"network_limits"`
	Checkpoint         checkpoint.GenesisState         `json:"checkpoint"`
}

func DefaultState() State {
	return State{
		Limits:     limits.DefaultGenesisState(),
		Assets:     assets.DefaultGenesisState(),
		Validators: validators.DefaultGenesisState(),
		Network:    abci.DefaultGenesis(),
		NetParams:  netparams.DefaultGenesisState(),
		Checkpoint: checkpoint.DefaultGenesisState(),
	}
}
