package assets_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/elysiumstation/fury/core/assets"
	erc20mocks "github.com/elysiumstation/fury/core/assets/erc20/mocks"
	"github.com/elysiumstation/fury/core/assets/mocks"
	bmocks "github.com/elysiumstation/fury/core/broker/mocks"
	nwethmocks "github.com/elysiumstation/fury/core/nodewallets/eth/mocks"
	"github.com/elysiumstation/fury/core/types"
	"github.com/elysiumstation/fury/libs/num"
	vgrand "github.com/elysiumstation/fury/libs/rand"
	"github.com/elysiumstation/fury/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testService struct {
	*assets.Service
	broker     *bmocks.MockInterface
	bridgeView *mocks.MockERC20BridgeView
	notary     *mocks.MockNotary
	ctrl       *gomock.Controller
	ethClient  *erc20mocks.MockETHClient
	ethWallet  *nwethmocks.MockEthereumWallet
}

func TestAssets(t *testing.T) {
	t.Run("Staging asset update for unknown asset fails", testStagingAssetUpdateForUnknownAssetFails)
	t.Run("Offers signature on tick success", testOffersSignaturesOnTickSuccess)
}

func testOffersSignaturesOnTickSuccess(t *testing.T) {
	service := getTestService(t)

	assetID := hex.EncodeToString([]byte("asset_id"))
	nodeSignature := []byte("node_signature")

	service.broker.EXPECT().Send(gomock.Any()).Times(2)
	service.ethClient.EXPECT().CollateralBridgeAddress().Times(1)
	service.ethWallet.EXPECT().Algo().Times(1)
	service.ethWallet.EXPECT().Sign(gomock.Any()).Return(nodeSignature, nil).Times(1)

	service.notary.EXPECT().
		OfferSignatures(gomock.Any(), gomock.Any()).DoAndReturn(
		func(kind types.NodeSignatureKind, f func(id string) []byte) {
			require.Equal(t, kind, types.NodeSignatureKindAssetNew)
			require.Equal(t, nodeSignature, f(assetID))
		},
	)

	assetDetails := &types.AssetDetails{
		Name:     vgrand.RandomStr(5),
		Symbol:   vgrand.RandomStr(3),
		Decimals: 10,
		Quantum:  num.DecimalFromInt64(42),
		Source: &types.AssetDetailsErc20{
			ERC20: &types.ERC20{
				ContractAddress:   vgrand.RandomStr(5),
				LifetimeLimit:     num.NewUint(42),
				WithdrawThreshold: num.NewUint(84),
			},
		},
	}

	ctx := context.Background()

	_, err := service.NewAsset(ctx, assetID, assetDetails)
	require.NoError(t, err)

	err = service.Enable(ctx, assetID)
	require.NoError(t, err)

	service.OnTick(ctx, time.Now())
}

func testStagingAssetUpdateForUnknownAssetFails(t *testing.T) {
	service := getTestService(t)

	// given
	asset := &types.Asset{
		ID: vgrand.RandomStr(5),
		Details: &types.AssetDetails{
			Name:     vgrand.RandomStr(5),
			Symbol:   vgrand.RandomStr(3),
			Decimals: 10,
			Quantum:  num.DecimalFromInt64(42),
			Source: &types.AssetDetailsErc20{
				ERC20: &types.ERC20{
					ContractAddress:   vgrand.RandomStr(5),
					LifetimeLimit:     num.NewUint(42),
					WithdrawThreshold: num.NewUint(84),
				},
			},
		},
	}

	// when
	err := service.StageAssetUpdate(asset)

	// then
	require.ErrorIs(t, err, assets.ErrAssetDoesNotExist)
}

func getTestService(t *testing.T) *testService {
	t.Helper()
	conf := assets.NewDefaultConfig()
	logger := logging.NewTestLogger()
	ctrl := gomock.NewController(t)
	ethClient := erc20mocks.NewMockETHClient(ctrl)
	broker := bmocks.NewMockInterface(ctrl)
	bridgeView := mocks.NewMockERC20BridgeView(ctrl)
	notary := mocks.NewMockNotary(ctrl)
	ethWallet := nwethmocks.NewMockEthereumWallet(ctrl)

	service := assets.New(logger, conf, ethWallet, ethClient, broker, bridgeView, notary, true)
	return &testService{
		Service:    service,
		broker:     broker,
		ctrl:       ctrl,
		bridgeView: bridgeView,
		notary:     notary,
		ethClient:  ethClient,
		ethWallet:  ethWallet,
	}
}
