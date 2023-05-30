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

//go:build !race
// +build !race

package nodewallets_test

import (
	"testing"

	"github.com/elysiumstation/fury/core/nodewallets"
	vgrand "github.com/elysiumstation/fury/libs/rand"
	vgtesting "github.com/elysiumstation/fury/libs/testing"
	"github.com/elysiumstation/fury/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Run("Getting node wallets succeeds", testHandlerGettingNodeWalletsSucceeds)
	t.Run("Getting node wallets with wrong registry passphrase fails", testHandlerGettingNodeWalletsWithWrongRegistryPassphraseFails)
	t.Run("Getting Ethereum wallet succeeds", testHandlerGettingEthereumWalletSucceeds)
	t.Run("Getting Ethereum wallet succeeds", testHandlerGettingEthereumWalletWithWrongRegistryPassphraseFails)
	t.Run("Getting Fury wallet succeeds", testHandlerGettingFuryWalletSucceeds)
	t.Run("Getting Fury wallet succeeds", testHandlerGettingFuryWalletWithWrongRegistryPassphraseFails)
	t.Run("Generating Ethereum wallet succeeds", testHandlerGeneratingEthereumWalletSucceeds)
	t.Run("Generating an already existing Ethereum wallet fails", testHandlerGeneratingAlreadyExistingEthereumWalletFails)
	t.Run("Generating Ethereum wallet with overwrite succeeds", testHandlerGeneratingEthereumWalletWithOverwriteSucceeds)
	t.Run("Generating Fury wallet succeeds", testHandlerGeneratingFuryWalletSucceeds)
	t.Run("Generating an already existing Fury wallet fails", testHandlerGeneratingAlreadyExistingFuryWalletFails)
	t.Run("Generating Fury wallet with overwrite succeeds", testHandlerGeneratingFuryWalletWithOverwriteSucceeds)
	t.Run("Importing Ethereum wallet succeeds", testHandlerImportingEthereumWalletSucceeds)
	t.Run("Importing an already existing Ethereum wallet fails", testHandlerImportingAlreadyExistingEthereumWalletFails)
	t.Run("Importing Ethereum wallet with overwrite succeeds", testHandlerImportingEthereumWalletWithOverwriteSucceeds)
	t.Run("Importing Fury wallet succeeds", testHandlerImportingFuryWalletSucceeds)
	t.Run("Importing an already existing Fury wallet fails", testHandlerImportingAlreadyExistingFuryWalletFails)
	t.Run("Importing Fury wallet with overwrite succeeds", testHandlerImportingFuryWalletWithOverwriteSucceeds)
}

func testHandlerGettingNodeWalletsSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)
	config := nodewallets.NewDefaultConfig()

	// setup
	createTestNodeWallets(furyPaths, registryPass, walletsPass)

	// when
	nw, err := nodewallets.GetNodeWallets(config, furyPaths, registryPass)

	// assert
	require.NoError(t, err)
	require.NotNil(t, nw)
	require.NotNil(t, nw.Ethereum)
	require.NotNil(t, nw.Fury)
}

func testHandlerGettingNodeWalletsWithWrongRegistryPassphraseFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	wrongRegistryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)
	config := nodewallets.NewDefaultConfig()

	// setup
	createTestNodeWallets(furyPaths, registryPass, walletsPass)

	// when
	nw, err := nodewallets.GetNodeWallets(config, furyPaths, wrongRegistryPass)

	// assert
	require.Error(t, err)
	assert.Nil(t, nw)
}

func testHandlerGettingEthereumWalletSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(furyPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetEthereumWallet(furyPaths, registryPass)

	// assert
	require.NoError(t, err)
	assert.NotNil(t, wallet)
}

func testHandlerGettingEthereumWalletWithWrongRegistryPassphraseFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	wrongRegistryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(furyPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetEthereumWallet(furyPaths, wrongRegistryPass)

	// assert
	require.Error(t, err)
	assert.Nil(t, wallet)
}

func testHandlerGettingFuryWalletSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(furyPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetFuryWallet(furyPaths, registryPass)

	// then
	require.NoError(t, err)
	assert.NotNil(t, wallet)
}

func testHandlerGettingFuryWalletWithWrongRegistryPassphraseFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	wrongRegistryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	// setup
	createTestNodeWallets(furyPaths, registryPass, walletsPass)

	// when
	wallet, err := nodewallets.GetFuryWallet(furyPaths, wrongRegistryPass)

	// assert
	require.Error(t, err)
	assert.Nil(t, wallet)
}

func testHandlerGeneratingEthereumWalletSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	data, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data["registryFilePath"])
	assert.NotEmpty(t, data["walletFilePath"])
}

func testHandlerGeneratingAlreadyExistingEthereumWalletFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass1, "", false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass2, "", false)

	// then
	require.EqualError(t, err, nodewallets.ErrEthereumWalletAlreadyExists.Error())
	assert.Empty(t, data2)
}

func testHandlerGeneratingEthereumWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass1, "", false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass2, "", true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data2["registryFilePath"])
	assert.Equal(t, data1["registryFilePath"], data2["registryFilePath"])
	assert.NotEmpty(t, data2["walletFilePath"])
	assert.NotEqual(t, data1["walletFilePath"], data2["walletFilePath"])
}

func testHandlerGeneratingFuryWalletSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	data, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data["registryFilePath"])
	assert.NotEmpty(t, data["walletFilePath"])
	assert.NotEmpty(t, data["mnemonic"])
}

func testHandlerGeneratingAlreadyExistingFuryWalletFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass1, false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])
	assert.NotEmpty(t, data1["mnemonic"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass2, false)

	// then
	require.EqualError(t, err, nodewallets.ErrFuryWalletAlreadyExists.Error())
	assert.Empty(t, data2)
}

func testHandlerGeneratingFuryWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass1 := vgrand.RandomStr(10)

	// when
	data1, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass1, false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data1["registryFilePath"])
	assert.NotEmpty(t, data1["walletFilePath"])

	// given
	walletPass2 := vgrand.RandomStr(10)

	// when
	data2, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass2, true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, data2["registryFilePath"])
	assert.Equal(t, data1["registryFilePath"], data2["registryFilePath"])
	assert.NotEmpty(t, data2["walletFilePath"])
	assert.NotEqual(t, data1["walletFilePath"], data2["walletFilePath"])
	assert.NotEmpty(t, data2["mnemonic"])
	assert.NotEqual(t, data1["mnemonic"], data2["mnemonic"])
}

func testHandlerImportingEthereumWalletSucceeds(t *testing.T) {
	// given
	genFuryPaths, genCleanupFn := vgtesting.NewFuryPaths()
	defer genCleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateEthereumWallet(genFuryPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)

	// given
	importFuryPaths, importCleanupFn := vgtesting.NewFuryPaths()
	defer importCleanupFn()

	// when
	importData, err := nodewallets.ImportEthereumWallet(importFuryPaths, registryPass, walletPass, "", "", genData["walletFilePath"], false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, importData["registryFilePath"])
	assert.NotEqual(t, genData["registryFilePath"], importData["registryFilePath"])
	assert.NotEmpty(t, importData["walletFilePath"])
	assert.NotEqual(t, genData["walletFilePath"], importData["walletFilePath"])
}

func testHandlerImportingAlreadyExistingEthereumWalletFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportEthereumWallet(furyPaths, registryPass, walletPass, "", genData["walletFilePath"], "", false)

	// then
	require.EqualError(t, err, nodewallets.ErrEthereumWalletAlreadyExists.Error())
	assert.Empty(t, importData)
}

func testHandlerImportingEthereumWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass, "", false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportEthereumWallet(furyPaths, registryPass, walletPass, "", "", genData["walletFilePath"], true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, genData["registryFilePath"])
	assert.Equal(t, importData["registryFilePath"], genData["registryFilePath"])
	assert.NotEmpty(t, genData["walletFilePath"])
	assert.Equal(t, importData["walletFilePath"], genData["walletFilePath"])
}

func testHandlerImportingFuryWalletSucceeds(t *testing.T) {
	// given
	genFuryPaths, genCleanupFn := vgtesting.NewFuryPaths()
	defer genCleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateFuryWallet(genFuryPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)

	// given
	importFuryPaths, importCleanupFn := vgtesting.NewFuryPaths()
	defer importCleanupFn()

	// when
	importData, err := nodewallets.ImportFuryWallet(importFuryPaths, registryPass, walletPass, genData["walletFilePath"], false)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, importData["registryFilePath"])
	assert.NotEqual(t, genData["registryFilePath"], importData["registryFilePath"])
	assert.NotEmpty(t, importData["walletFilePath"])
	assert.NotEqual(t, genData["walletFilePath"], importData["walletFilePath"])
}

func testHandlerImportingAlreadyExistingFuryWalletFails(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportFuryWallet(furyPaths, registryPass, walletPass, genData["walletFilePath"], false)

	// then
	require.EqualError(t, err, nodewallets.ErrFuryWalletAlreadyExists.Error())
	assert.Empty(t, importData)
}

func testHandlerImportingFuryWalletWithOverwriteSucceeds(t *testing.T) {
	// given
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletPass := vgrand.RandomStr(10)

	// when
	genData, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass, false)

	// then
	require.NoError(t, err)

	// when
	importData, err := nodewallets.ImportFuryWallet(furyPaths, registryPass, walletPass, genData["walletFilePath"], true)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, importData["registryFilePath"])
	assert.Equal(t, genData["registryFilePath"], importData["registryFilePath"])
	assert.NotEmpty(t, importData["walletFilePath"])
	assert.NotEqual(t, genData["walletFilePath"], importData["walletFilePath"])
}

func createTestNodeWallets(furyPaths paths.Paths, registryPass, walletPass string) {
	if _, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletPass, "", false); err != nil {
		panic("couldn't generate Ethereum node wallet for tests")
	}

	if _, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletPass, false); err != nil {
		panic("couldn't generate Fury node wallet for tests")
	}
}
