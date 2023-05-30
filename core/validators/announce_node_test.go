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

package validators_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/elysiumstation/fury/core/nodewallets"
	"github.com/elysiumstation/fury/core/validators"
	vgrand "github.com/elysiumstation/fury/libs/rand"
	vgtesting "github.com/elysiumstation/fury/libs/testing"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestTendermintKey(t *testing.T) {
	t.Parallel()
	notBase64 := "170ffakjde"
	require.Error(t, validators.VerifyTendermintKey(notBase64))

	validKey := "794AFpbqJvHF711mhAK3fvSLnoXuuiig2ecrdeSJ/bk="
	require.NoError(t, validators.VerifyTendermintKey(validKey))
}

func TestAnnounceNode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tt := getTestTopology(t)
	cmd := createSignedAnnounceCommand(t)

	// Now announce it and check the signature verify
	require.NoError(t, tt.Topology.ProcessAnnounceNode(ctx, cmd))

	// Announce it again
	require.ErrorIs(t, tt.Topology.ProcessAnnounceNode(ctx, cmd), validators.ErrFuryNodeAlreadyRegisterForChain)
}

func createSignedAnnounceCommand(t *testing.T) *commandspb.AnnounceNode {
	t.Helper()
	nodeWallets := createTestNodeWallets(t)
	cmd := commandspb.AnnounceNode{
		Id:              nodeWallets.Fury.ID().Hex(),
		FuryPubKey:      nodeWallets.Fury.PubKey().Hex(),
		FuryPubKeyIndex: nodeWallets.Fury.Index(),
		ChainPubKey:     "794AFpbqJvHF711mhAK3fvSLnoXuuiig2ecrdeSJ/bk=",
		EthereumAddress: nodeWallets.Ethereum.PubKey().Hex(),
		FromEpoch:       1,
		InfoUrl:         "www.some.com",
		Name:            "that is not my name",
		AvatarUrl:       "www.avatar.com",
		Country:         "some country",
	}
	err := validators.SignAnnounceNode(&cmd, nodeWallets.Fury, nodeWallets.Ethereum)
	require.NoError(t, err)

	// verify that the expected signature for fury key is there
	messageToSign := cmd.Id + cmd.FuryPubKey + fmt.Sprintf("%d", cmd.FuryPubKeyIndex) + cmd.ChainPubKey + cmd.EthereumAddress + fmt.Sprintf("%d", cmd.FromEpoch) + cmd.InfoUrl + cmd.Name + cmd.AvatarUrl + cmd.Country
	sig, err := nodeWallets.Fury.Sign([]byte(messageToSign))
	sigHex := hex.EncodeToString(sig)
	require.NoError(t, err)
	require.Equal(t, sigHex, cmd.FurySignature.Value)

	// verify that the expected signature for eth key is there
	ethSig, err := nodeWallets.Ethereum.Sign(crypto.Keccak256([]byte(messageToSign)))
	ethSigHex := hex.EncodeToString(ethSig)
	require.NoError(t, err)
	require.Equal(t, ethSigHex, cmd.EthereumSignature.Value)

	return &cmd
}

func createTestNodeWallets(t *testing.T) *nodewallets.NodeWallets {
	t.Helper()
	config := nodewallets.NewDefaultConfig()
	furyPaths, cleanupFn := vgtesting.NewFuryPaths()
	defer cleanupFn()
	registryPass := vgrand.RandomStr(10)
	walletsPass := vgrand.RandomStr(10)

	if _, err := nodewallets.GenerateEthereumWallet(furyPaths, registryPass, walletsPass, "", false); err != nil {
		t.Fatal("couldn't generate Ethereum node wallet for tests")
	}

	if _, err := nodewallets.GenerateFuryWallet(furyPaths, registryPass, walletsPass, false); err != nil {
		t.Fatal("couldn't generate Fury node wallet for tests")
	}
	nw, err := nodewallets.GetNodeWallets(config, furyPaths, registryPass)
	require.NoError(t, err)
	return nw
}
