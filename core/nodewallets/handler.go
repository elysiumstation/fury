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

package nodewallets

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/elysiumstation/fury/core/nodewallets/registry"
	"github.com/elysiumstation/fury/core/nodewallets/fury"
	"github.com/elysiumstation/fury/paths"
)

var (
	ErrEthereumWalletAlreadyExists   = errors.New("the Ethereum node wallet already exists")
	ErrFuryWalletAlreadyExists       = errors.New("the Fury node wallet already exists")
	ErrTendermintPubkeyAlreadyExists = errors.New("the Tendermint pubkey already exists")
)

func GetFuryWallet(furyPaths paths.Paths, registryPassphrase string) (*fury.Wallet, error) {
	registryLoader, err := registry.NewLoader(furyPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	registry, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if registry.Fury == nil {
		return nil, ErrFuryWalletIsMissing
	}

	walletLoader, err := fury.InitialiseWalletLoader(furyPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise Fury node wallet loader: %w", err)
	}

	wallet, err := walletLoader.Load(registry.Fury.Name, registry.Fury.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load Ethereum node wallet: %w", err)
	}

	return wallet, nil
}

func GetNodeWallets(config Config, furyPaths paths.Paths, registryPassphrase string) (*NodeWallets, error) {
	nodeWallets := &NodeWallets{}

	registryLoader, err := registry.NewLoader(furyPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if reg.Ethereum != nil {
		w, err := GetEthereumWalletWithRegistry(furyPaths, reg)
		if err != nil {
			return nil, err
		}

		nodeWallets.Ethereum = w
	}

	if reg.Fury != nil {
		furyWalletLoader, err := fury.InitialiseWalletLoader(furyPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise Fury node wallet loader: %w", err)
		}

		nodeWallets.Fury, err = furyWalletLoader.Load(reg.Fury.Name, reg.Fury.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("couldn't load Fury node wallet: %w", err)
		}
	}

	if reg.Tendermint != nil {
		nodeWallets.Tendermint = &TendermintPubkey{
			Pubkey: reg.Tendermint.Pubkey,
		}
	}

	return nodeWallets, nil
}

func GenerateFuryWallet(furyPaths paths.Paths, registryPassphrase, walletPassphrase string, overwrite bool) (map[string]string, error) {
	registryLoader, err := registry.NewLoader(furyPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if !overwrite && reg.Fury != nil {
		return nil, ErrFuryWalletAlreadyExists
	}

	furyWalletLoader, err := fury.InitialiseWalletLoader(furyPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise Fury node wallet loader: %w", err)
	}

	w, data, err := furyWalletLoader.Generate(walletPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate Fury node wallet: %w", err)
	}

	reg.Fury = &registry.RegisteredFuryWallet{
		Name:       w.Name(),
		Passphrase: walletPassphrase,
	}

	if err := registryLoader.Save(reg, registryPassphrase); err != nil {
		return nil, fmt.Errorf("couldn't save registry: %w", err)
	}

	data["registryFilePath"] = registryLoader.RegistryFilePath()
	return data, nil
}

func ImportFuryWallet(furyPaths paths.Paths, registryPassphrase, walletPassphrase, sourceFilePath string, overwrite bool) (map[string]string, error) {
	if !filepath.IsAbs(sourceFilePath) {
		return nil, fmt.Errorf("path to the wallet file need to be absolute")
	}

	registryLoader, err := registry.NewLoader(furyPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if !overwrite && reg.Fury != nil {
		return nil, ErrFuryWalletAlreadyExists
	}

	furyWalletLoader, err := fury.InitialiseWalletLoader(furyPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise Fury node wallet loader: %w", err)
	}

	w, data, err := furyWalletLoader.Import(sourceFilePath, walletPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Fury node wallet: %w", err)
	}

	reg.Fury = &registry.RegisteredFuryWallet{
		Name:       w.Name(),
		Passphrase: walletPassphrase,
	}

	if err := registryLoader.Save(reg, registryPassphrase); err != nil {
		return nil, fmt.Errorf("couldn't save registry: %w", err)
	}

	data["registryFilePath"] = registryLoader.RegistryFilePath()
	return data, nil
}

func ImportTendermintPubkey(
	furyPaths paths.Paths,
	registryPassphrase, pubkey string,
	overwrite bool,
) (map[string]string, error) {
	registryLoader, err := registry.NewLoader(furyPaths, registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise node wallet registry: %v", err)
	}

	reg, err := registryLoader.Get(registryPassphrase)
	if err != nil {
		return nil, fmt.Errorf("couldn't load node wallet registry: %v", err)
	}

	if !overwrite && reg.Tendermint != nil {
		return nil, ErrTendermintPubkeyAlreadyExists
	}

	reg.Tendermint = &registry.RegisteredTendermintPubkey{
		Pubkey: pubkey,
	}

	if err := registryLoader.Save(reg, registryPassphrase); err != nil {
		return nil, fmt.Errorf("couldn't save registry: %w", err)
	}

	return map[string]string{
		"registryFilePath": registryLoader.RegistryFilePath(),
		"tendermintPubkey": pubkey,
	}, nil
}
