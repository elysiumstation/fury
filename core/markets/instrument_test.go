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

package markets_test

import (
	"context"
	"testing"

	"github.com/elysiumstation/fury/core/broker/mocks"
	emock "github.com/elysiumstation/fury/core/execution/common/mocks"
	"github.com/elysiumstation/fury/core/markets"
	"github.com/elysiumstation/fury/core/oracles"
	"github.com/elysiumstation/fury/logging"

	"github.com/elysiumstation/fury/core/products"
	"github.com/elysiumstation/fury/core/types"
	furypb "github.com/elysiumstation/fury/protos/fury"
	datapb "github.com/elysiumstation/fury/protos/fury/data/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstrument(t *testing.T) {
	t.Run("Create a valid new instrument", func(t *testing.T) {
		pinst := getValidInstrumentProto()
		inst, err := markets.NewInstrument(context.Background(), logging.NewTestLogger(), pinst, newOracleEngine(t))
		assert.NotNil(t, inst)
		assert.Nil(t, err)
	})

	t.Run("nil product", func(t *testing.T) {
		pinst := getValidInstrumentProto()
		pinst.Product = nil
		inst, err := markets.NewInstrument(context.Background(), logging.NewTestLogger(), pinst, newOracleEngine(t))
		assert.Nil(t, inst)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "unable to instantiate product from instrument configuration: nil product")
	})

	t.Run("nil oracle spec", func(t *testing.T) {
		pinst := getValidInstrumentProto()
		pinst.Product = &types.InstrumentFuture{
			Future: &types.Future{
				SettlementAsset:                     "Ethereum/Ether",
				DataSourceSpecForSettlementData:     nil,
				DataSourceSpecForTradingTermination: nil,
				DataSourceSpecBinding: &types.DataSourceSpecBindingForFuture{
					SettlementDataProperty:     "prices.ETH.value",
					TradingTerminationProperty: "trading.terminated",
				},
			},
		}
		inst, err := markets.NewInstrument(context.Background(), logging.NewTestLogger(), pinst, newOracleEngine(t))
		require.NotNil(t, err)
		assert.Nil(t, inst)
		assert.Equal(t, "unable to instantiate product from instrument configuration: a data source spec and spec binding are required", err.Error())
	})

	t.Run("nil oracle spec binding", func(t *testing.T) {
		pinst := getValidInstrumentProto()
		pinst.Product = &types.InstrumentFuture{
			Future: &types.Future{
				SettlementAsset: "Ethereum/Ether",
				DataSourceSpecForSettlementData: &types.DataSourceSpec{
					Data: types.NewDataSourceDefinition(
						furypb.DataSourceDefinitionTypeExt,
					).SetOracleConfig(
						&types.DataSourceSpecConfiguration{
							Signers: []*types.Signer{types.CreateSignerFromString("0xDEADBEEF", types.DataSignerTypePubKey)},
							Filters: []*types.DataSourceSpecFilter{
								{
									Key: &types.DataSourceSpecPropertyKey{
										Name: "prices.ETH.value",
										Type: datapb.PropertyKey_TYPE_INTEGER,
									},
									Conditions: []*types.DataSourceSpecCondition{},
								},
							},
						},
					),
				},
				DataSourceSpecForTradingTermination: &types.DataSourceSpec{
					Data: types.NewDataSourceDefinition(
						furypb.DataSourceDefinitionTypeExt,
					).SetOracleConfig(
						&types.DataSourceSpecConfiguration{
							Signers: []*types.Signer{types.CreateSignerFromString("0xDEADBEEF", types.DataSignerTypePubKey)},
							Filters: []*types.DataSourceSpecFilter{
								{
									Key: &types.DataSourceSpecPropertyKey{
										Name: "trading.terminated",
										Type: datapb.PropertyKey_TYPE_BOOLEAN,
									},
									Conditions: []*types.DataSourceSpecCondition{},
								},
							},
						},
					),
				},
				DataSourceSpecBinding: nil,
			},
		}
		inst, err := markets.NewInstrument(context.Background(), logging.NewTestLogger(), pinst, newOracleEngine(t))
		require.NotNil(t, err)
		assert.Nil(t, inst)
		assert.Equal(t, "unable to instantiate product from instrument configuration: a data source spec and spec binding are required", err.Error())
	})
}

func newOracleEngine(t *testing.T) products.OracleEngine {
	t.Helper()
	ctrl := gomock.NewController(t)
	broker := mocks.NewMockBroker(ctrl)
	broker.EXPECT().Send(gomock.Any()).AnyTimes()

	ts := emock.NewMockTimeService(ctrl)
	ts.EXPECT().GetTimeNow().AnyTimes()

	return oracles.NewEngine(
		logging.NewTestLogger(),
		oracles.NewDefaultConfig(),
		ts,
		broker,
	)
}

func getValidInstrumentProto() *types.Instrument {
	return &types.Instrument{
		ID:   "Crypto/BTCUSD/Futures/Dec19",
		Code: "FX:BTCUSD/DEC19",
		Name: "December 2019 BTC vs USD future",
		Metadata: &types.InstrumentMetadata{
			Tags: []string{
				"asset_class:fx/crypto",
				"product:futures",
			},
		},
		Product: &types.InstrumentFuture{
			Future: &types.Future{
				QuoteName:       "USD",
				SettlementAsset: "Ethereum/Ether",
				DataSourceSpecForSettlementData: &types.DataSourceSpec{
					Data: types.NewDataSourceDefinition(
						furypb.DataSourceDefinitionTypeExt,
					).SetOracleConfig(
						&types.DataSourceSpecConfiguration{
							Signers: []*types.Signer{types.CreateSignerFromString("0xDEADBEEF", types.DataSignerTypePubKey)},
							Filters: []*types.DataSourceSpecFilter{
								{
									Key: &types.DataSourceSpecPropertyKey{
										Name: "prices.ETH.value",
										Type: datapb.PropertyKey_TYPE_INTEGER,
									},
									Conditions: []*types.DataSourceSpecCondition{},
								},
							},
						},
					),
				},
				DataSourceSpecForTradingTermination: &types.DataSourceSpec{
					Data: types.NewDataSourceDefinition(
						furypb.DataSourceDefinitionTypeExt,
					).SetOracleConfig(
						&types.DataSourceSpecConfiguration{
							Signers: []*types.Signer{types.CreateSignerFromString("0xDEADBEEF", types.DataSignerTypePubKey)},
							Filters: []*types.DataSourceSpecFilter{
								{
									Key: &types.DataSourceSpecPropertyKey{
										Name: "trading.terminated",
										Type: datapb.PropertyKey_TYPE_BOOLEAN,
									},
									Conditions: []*types.DataSourceSpecCondition{},
								},
							},
						},
					),
				},
				DataSourceSpecBinding: &types.DataSourceSpecBindingForFuture{
					SettlementDataProperty:     "prices.ETH.value",
					TradingTerminationProperty: "trading.terminated",
				},
			},
		},
	}
}
