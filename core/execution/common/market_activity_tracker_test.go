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

package common_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/elysiumstation/fury/core/execution/common"
	"github.com/elysiumstation/fury/core/integration/stubs"
	snp "github.com/elysiumstation/fury/core/snapshot"
	"github.com/elysiumstation/fury/core/stats"
	"github.com/elysiumstation/fury/core/types"
	vgcontext "github.com/elysiumstation/fury/libs/context"
	"github.com/elysiumstation/fury/libs/num"
	"github.com/elysiumstation/fury/libs/proto"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"
	vgproto "github.com/elysiumstation/fury/protos/fury"
	snapshotpb "github.com/elysiumstation/fury/protos/fury/snapshot/v1"
	"github.com/stretchr/testify/require"
)

type TestEpochEngine struct {
	target func(context.Context, types.Epoch)
}

func (e *TestEpochEngine) NotifyOnEpoch(f func(context.Context, types.Epoch), _ func(context.Context, types.Epoch)) {
	e.target = f
}

type EligibilityChecker struct{}

func (e *EligibilityChecker) IsEligibleForProposerBonus(marketID string, volumeTraded *num.Uint) bool {
	return volumeTraded.GT(num.NewUint(5000))
}

func TestMarketTracker(t *testing.T) {
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	tracker.SetEligibilityChecker(&EligibilityChecker{})

	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")

	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))

	tracker.AddValueTraded("market1", num.NewUint(1000))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))

	tracker.AddValueTraded("market2", num.NewUint(4000))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))

	tracker.AddValueTraded("market2", num.NewUint(1001))
	tracker.AddValueTraded("market1", num.NewUint(4001))

	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))

	// mark as paid
	tracker.MarkPaidProposer("market1", "FURY", []string{}, "zohar")
	tracker.MarkPaidProposer("market2", "FURY", []string{}, "zohar")

	// check if eligible for the same combo, expect false
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))

	// now check for another funder
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "jeremy"))
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "jeremy"))

	// mark as paid
	tracker.MarkPaidProposer("market1", "FURY", []string{}, "jeremy")
	tracker.MarkPaidProposer("market2", "FURY", []string{}, "jeremy")

	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "jeremy"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "jeremy"))

	// check for another payout asset
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{}, "zohar"))
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{}, "zohar"))

	tracker.MarkPaidProposer("market1", "USDC", []string{}, "zohar")
	tracker.MarkPaidProposer("market2", "USDC", []string{}, "zohar")

	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{}, "zohar"))

	// check for another market scope
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{"market1"}, "zohar"))
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{"market2"}, "zohar"))
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{"market1", "market2"}, "zohar"))
	require.Equal(t, true, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{"market2", "market2"}, "zohar"))

	tracker.MarkPaidProposer("market1", "USDC", []string{"market1"}, "zohar")
	tracker.MarkPaidProposer("market2", "USDC", []string{"market2"}, "zohar")
	tracker.MarkPaidProposer("market1", "USDC", []string{"market1", "market2"}, "zohar")
	tracker.MarkPaidProposer("market2", "USDC", []string{"market1", "market2"}, "zohar")

	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{"market1"}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{"market2"}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{"market1", "market2"}, "zohar"))
	require.Equal(t, false, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{"market1", "market2"}, "zohar"))

	// take a snapshot
	key := (&types.PayloadMarketActivityTracker{}).Key()
	state1, _, err := tracker.GetState(key)
	require.NoError(t, err)

	trackerLoad := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	pl := snapshotpb.Payload{}
	require.NoError(t, proto.Unmarshal(state1, &pl))

	trackerLoad.LoadState(context.Background(), types.PayloadFromProto(&pl))

	state2, _, err := trackerLoad.GetState(key)
	require.NoError(t, err)
	require.True(t, bytes.Equal(state1, state2))
}

func TestRemoveMarket(t *testing.T) {
	epochService := &TestEpochEngine{}
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), epochService)
	tracker.SetEligibilityChecker(&EligibilityChecker{})
	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")
	require.Equal(t, 2, len(tracker.GetAllMarketIDs()))
	require.Equal(t, "market1", tracker.GetAllMarketIDs()[0])
	require.Equal(t, "market2", tracker.GetAllMarketIDs()[1])

	// remove the market - this should only mark the market for removal
	tracker.RemoveMarket("market1")
	require.Equal(t, 2, len(tracker.GetAllMarketIDs()))
	require.Equal(t, "market1", tracker.GetAllMarketIDs()[0])
	require.Equal(t, "market2", tracker.GetAllMarketIDs()[1])
	epochService.target(context.Background(), types.Epoch{Action: vgproto.EpochAction_EPOCH_ACTION_START})

	require.Equal(t, 1, len(tracker.GetAllMarketIDs()))
	require.Equal(t, "market2", tracker.GetAllMarketIDs()[0])
}

func TestGetMarketScores(t *testing.T) {
	epochService := &TestEpochEngine{}
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), epochService)
	tracker.SetEligibilityChecker(&EligibilityChecker{})
	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")
	tracker.MarketProposed("asset1", "market4", "me4")
	tracker.MarketProposed("asset2", "market3", "me3")

	// no fees generated expect empty slice
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))

	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))

	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))

	require.Equal(t, 0, len(tracker.GetMarketScores("asset2", []string{"market3"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset2", []string{"market3"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	require.Equal(t, 0, len(tracker.GetMarketScores("asset2", []string{"market3"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))

	// update with a few transfers
	transfersM1 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(100)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party1", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(400)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(300)}},
		{Owner: "party1", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(600)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(900)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(800)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(700)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(600)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1000)}},
	}
	tracker.UpdateFeesFromTransfers("market1", transfersM1)

	transfersM2 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(500)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1500)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1500)}},
	}
	tracker.UpdateFeesFromTransfers("market2", transfersM2)

	transfersM3 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(500)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(450)}},
	}
	tracker.UpdateFeesFromTransfers("market3", transfersM3)

	// in market1: 2500
	// in market2: 1500
	// in market4: 0 => it is not included in the scores.
	require.Equal(t, 2, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	LPMarket1 := &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market1",
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED,
		Score:  num.MustDecimalFromString("0.625"),
	}
	LPMarket2 := &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market2",
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED,
		Score:  num.MustDecimalFromString("0.375"),
	}
	require.Equal(t, 2, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	assertMarketContributionScore(t, LPMarket1, tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)[0])
	assertMarketContributionScore(t, LPMarket2, tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)[1])

	// scope only market1:
	require.Equal(t, 1, len(tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	LPMarket1.Score = num.DecimalFromInt64(1)
	assertMarketContributionScore(t, LPMarket1, tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)[0])

	// scope only market2:
	require.Equal(t, 1, len(tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))
	LPMarket2.Score = num.DecimalFromInt64(1)
	assertMarketContributionScore(t, LPMarket2, tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)[0])

	// try to scope market3: doesn't exist in the asset
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market3"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))

	// try to get the market from the wrong asset
	require.Equal(t, 0, len(tracker.GetMarketScores("asset2", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_LP_FEES_RECEIVED)))

	// in market1: 2000
	// in market2: 500
	require.Equal(t, 2, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	LPMarket1 = &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market1",
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED,
		Score:  num.MustDecimalFromString("0.8"),
	}
	LPMarket2 = &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market2",
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED,
		Score:  num.MustDecimalFromString("0.2"),
	}
	require.Equal(t, 2, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	assertMarketContributionScore(t, LPMarket1, tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)[0])
	assertMarketContributionScore(t, LPMarket2, tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)[1])

	// scope only market1:
	require.Equal(t, 1, len(tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	LPMarket1.Score = num.DecimalFromInt64(1)
	assertMarketContributionScore(t, LPMarket1, tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)[0])

	// scope only market2:
	require.Equal(t, 1, len(tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))
	LPMarket2.Score = num.DecimalFromInt64(1)
	assertMarketContributionScore(t, LPMarket2, tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)[0])

	// try to scope market3: doesn't exist in the asset
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market3"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))

	// try to get the market from the wrong asset
	require.Equal(t, 0, len(tracker.GetMarketScores("asset2", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_RECEIVED)))

	// in market1: 1500
	// in market2: 1500
	require.Equal(t, 2, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))
	LPMarket1 = &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market1",
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID,
		Score:  num.MustDecimalFromString("0.5"),
	}
	LPMarket2 = &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market2",
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID,
		Score:  num.MustDecimalFromString("0.5"),
	}
	require.Equal(t, 2, len(tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))
	assertMarketContributionScore(t, LPMarket1, tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)[0])
	assertMarketContributionScore(t, LPMarket2, tracker.GetMarketScores("asset1", []string{}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)[1])

	// scope only market1:
	require.Equal(t, 1, len(tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))
	LPMarket1.Score = num.DecimalFromInt64(1)
	assertMarketContributionScore(t, LPMarket1, tracker.GetMarketScores("asset1", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)[0])

	// scope only market2:
	require.Equal(t, 1, len(tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))
	LPMarket2.Score = num.DecimalFromInt64(1)
	assertMarketContributionScore(t, LPMarket2, tracker.GetMarketScores("asset1", []string{"market2"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)[0])

	// try to scope market3: doesn't exist in the asset
	require.Equal(t, 0, len(tracker.GetMarketScores("asset1", []string{"market3"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))

	// try to get the market from the wrong asset
	require.Equal(t, 0, len(tracker.GetMarketScores("asset2", []string{"market1"}, vgproto.DispatchMetric_DISPATCH_METRIC_MAKER_FEES_PAID)))
}

func TestGetMarketsWithEligibleProposer(t *testing.T) {
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	tracker.SetEligibilityChecker(&EligibilityChecker{})
	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")
	tracker.MarketProposed("asset1", "market3", "me3")
	tracker.MarketProposed("asset2", "market4", "me4")
	tracker.MarketProposed("asset3", "market5", "me5")

	tracker.AddValueTraded("market2", num.NewUint(1001))
	tracker.AddValueTraded("market1", num.NewUint(4001))

	// the threshold is 5000 so expect at this point no market should be returned
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1"}, "FURY", "zohar")))
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market2"}, "FURY", "zohar")))
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2"}, "FURY", "zohar")))

	// market1 goes above the threshold
	tracker.AddValueTraded("market1", num.NewUint(1000))
	tracker.AddValueTraded("market4", num.NewUint(5001))
	require.Equal(t, 2, len(tracker.GetMarketsWithEligibleProposer("", []string{"market1", "market2", "market3", "market4"}, "FURY", "zohar")))

	expectedScoreMarket1Full := &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market1",
		Score:  num.DecimalFromInt64(1),
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MARKET_VALUE,
	}
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket1Full, tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")[0])
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1"}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket1Full, tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1"}, "FURY", "zohar")[0])
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market2"}, "FURY", "zohar")))
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2"}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket1Full, tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2"}, "FURY", "zohar")[0])

	// now market 2 goes above the threshold as well so expect the scores to be 0.5 for each
	tracker.AddValueTraded("market2", num.NewUint(4000))
	expectedScoreMarket1Half := &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market1",
		Score:  num.MustDecimalFromString("0.5"),
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MARKET_VALUE,
	}
	expectedScoreMarket2Half := &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market2",
		Score:  num.MustDecimalFromString("0.5"),
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MARKET_VALUE,
	}
	expectedScoreMarket2Full := &types.MarketContributionScore{
		Asset:  "asset1",
		Market: "market2",
		Score:  num.DecimalFromInt64(1),
		Metric: vgproto.DispatchMetric_DISPATCH_METRIC_MARKET_VALUE,
	}
	require.Equal(t, 2, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket1Half, tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")[0])
	assertMarketContributionScore(t, expectedScoreMarket2Half, tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")[1])
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1"}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket1Full, tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1"}, "FURY", "zohar")[0])
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market2"}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket2Full, tracker.GetMarketsWithEligibleProposer("asset1", []string{"market2"}, "FURY", "zohar")[0])
	require.Equal(t, 2, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2"}, "FURY", "zohar")))
	assertMarketContributionScore(t, expectedScoreMarket1Half, tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2"}, "FURY", "zohar")[0])
	assertMarketContributionScore(t, expectedScoreMarket2Half, tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2"}, "FURY", "zohar")[1])

	// all asset all markets
	// markets 1, 2, 4
	require.Equal(t, 3, len(tracker.GetMarketsWithEligibleProposer("", []string{}, "FURY", "zohar")))

	// asset with no markets
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset3", []string{}, "FURY", "zohar")))
}

func assertMarketContributionScore(t *testing.T, expected, actual *types.MarketContributionScore) {
	t.Helper()
	require.Equal(t, expected.Asset, actual.Asset)
	require.Equal(t, expected.Market, actual.Market)
	require.Equal(t, expected.Score.String(), actual.Score.String())
	require.Equal(t, expected.Metric, actual.Metric)
}

func TestMarketTrackerStateChange(t *testing.T) {
	key := (&types.PayloadMarketActivityTracker{}).Key()

	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	tracker.SetEligibilityChecker(&EligibilityChecker{})

	state1, _, err := tracker.GetState(key)
	require.NoError(t, err)

	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")

	state2, _, err := tracker.GetState(key)
	require.NoError(t, err)
	require.False(t, bytes.Equal(state1, state2))

	tracker.AddValueTraded("market1", num.NewUint(1000))
	require.False(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))

	state3, _, err := tracker.GetState(key)
	require.NoError(t, err)
	require.False(t, bytes.Equal(state1, state3))
}

func TestFeesTracker(t *testing.T) {
	epochEngine := &TestEpochEngine{}
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), epochEngine)
	epochEngine.target(context.Background(), types.Epoch{Seq: 1})

	partyScores := tracker.GetFeePartyScores("does not exist", types.TransferTypeMakerFeeReceive)
	require.Equal(t, 0, len(partyScores))

	key := (&types.PayloadMarketActivityTracker{}).Key()
	state1, _, err := tracker.GetState(key)
	require.NoError(t, err)

	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")

	// update with a few transfers
	transfersM1 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(100)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party1", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(400)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(300)}},
		{Owner: "party1", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(600)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(900)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(800)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(700)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(600)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1000)}},
	}
	tracker.UpdateFeesFromTransfers("market1", transfersM1)

	transfersM2 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(150)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(150)}},
	}
	tracker.UpdateFeesFromTransfers("market2", transfersM2)

	// asset1, types.TransferTypeMakerFeeReceive
	// party1 received 500
	// party2 received 1500
	scores := tracker.GetFeePartyScores("market1", types.TransferTypeMakerFeeReceive)
	require.Equal(t, "0.25", scores[0].Score.String())
	require.Equal(t, "party1", scores[0].Party)
	require.Equal(t, "0.75", scores[1].Score.String())
	require.Equal(t, "party2", scores[1].Party)

	// asset1 TransferTypeMakerFeePay
	// party1 paid 500
	// party2 paid 1000
	scores = tracker.GetFeePartyScores("market1", types.TransferTypeMakerFeePay)
	require.Equal(t, "0.3333333333333333", scores[0].Score.String())
	require.Equal(t, "party1", scores[0].Party)
	require.Equal(t, "0.6666666666666667", scores[1].Score.String())
	require.Equal(t, "party2", scores[1].Party)

	// asset1 TransferTypeLiquidityFeeDistribute
	// party1 paid 800
	// party2 paid 1700
	scores = tracker.GetFeePartyScores("market1", types.TransferTypeLiquidityFeeDistribute)
	require.Equal(t, "0.32", scores[0].Score.String())
	require.Equal(t, "party1", scores[0].Party)
	require.Equal(t, "0.68", scores[1].Score.String())
	require.Equal(t, "party2", scores[1].Party)

	// asset2 TransferTypeMakerFeePay
	scores = tracker.GetFeePartyScores("market2", types.TransferTypeMakerFeeReceive)
	require.Equal(t, 1, len(scores))
	require.Equal(t, "1", scores[0].Score.String())
	require.Equal(t, "party1", scores[0].Party)

	// asset2 TransferTypeMakerFeePay
	scores = tracker.GetFeePartyScores("market2", types.TransferTypeMakerFeePay)
	require.Equal(t, 1, len(scores))
	require.Equal(t, "1", scores[0].Score.String())
	require.Equal(t, "party2", scores[0].Party)

	// check state has changed
	state2, _, err := tracker.GetState(key)
	require.NoError(t, err)
	require.False(t, bytes.Equal(state1, state2))

	epochEngineLoad := &TestEpochEngine{}
	trackerLoad := common.NewMarketActivityTracker(logging.NewTestLogger(), epochEngineLoad)
	epochEngineLoad.target(context.Background(), types.Epoch{Seq: 1})

	pl := snapshotpb.Payload{}
	require.NoError(t, proto.Unmarshal(state2, &pl))
	trackerLoad.LoadState(context.Background(), types.PayloadFromProto(&pl))

	state3, _, err := trackerLoad.GetState(key)
	require.NoError(t, err)
	require.True(t, bytes.Equal(state2, state3))

	// check a restored party exist in the restored engine
	scores = trackerLoad.GetFeePartyScores("market2", types.TransferTypeMakerFeeReceive)
	require.Equal(t, 1, len(scores))
	require.Equal(t, "1", scores[0].Score.String())
	require.Equal(t, "party1", scores[0].Party)

	// New epoch should scrub the state an produce a difference hash
	epochEngineLoad.target(context.Background(), types.Epoch{Seq: 2, Action: vgproto.EpochAction_EPOCH_ACTION_START})
	state4, _, err := trackerLoad.GetState(key)
	require.NoError(t, err)
	require.False(t, bytes.Equal(state3, state4))

	// new epoch, we expect the metrics to have been reset
	for _, metric := range []types.TransferType{types.TransferTypeMakerFeePay, types.TransferTypeMakerFeeReceive, types.TransferTypeLiquidityFeeDistribute} {
		require.Equal(t, 0, len(trackerLoad.GetFeePartyScores("market1", metric)))
		require.Equal(t, 0, len(trackerLoad.GetFeePartyScores("market2", metric)))
	}
}

func TestSnapshot(t *testing.T) {
	tracker := setupDefaultTrackerForTest(t)

	// take a snapshot
	key := (&types.PayloadMarketActivityTracker{}).Key()
	state1, _, err := tracker.GetState(key)
	require.NoError(t, err)

	trackerLoad := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	pl := snapshotpb.Payload{}
	require.NoError(t, proto.Unmarshal(state1, &pl))

	trackerLoad.LoadState(context.Background(), types.PayloadFromProto(&pl))
	state2, _, err := trackerLoad.GetState(key)
	require.NoError(t, err)
	require.True(t, bytes.Equal(state1, state2))
}

func TestCheckpoint(t *testing.T) {
	tracker := setupDefaultTrackerForTest(t)

	b, err := tracker.Checkpoint()
	require.NoError(t, err)

	trackerLoad := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	trackerLoad.Load(context.Background(), b)

	bLoad, err := trackerLoad.Checkpoint()
	require.NoError(t, err)
	require.True(t, bytes.Equal(b, bLoad))
}

func setupDefaultTrackerForTest(t *testing.T) *common.MarketActivityTracker {
	t.Helper()
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	tracker.SetEligibilityChecker(&EligibilityChecker{})

	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me2")
	tracker.MarketProposed("asset1", "market4", "me4")
	tracker.MarketProposed("asset2", "market3", "me3")

	// update with a few transfers
	transfersM1 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(100)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party1", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(400)}},
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(300)}},
		{Owner: "party1", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(600)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(900)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(800)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(700)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(600)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1000)}},
	}
	tracker.UpdateFeesFromTransfers("market1", transfersM1)

	transfersM2 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(500)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1500)}},
		{Owner: "party2", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(1500)}},
	}
	tracker.UpdateFeesFromTransfers("market2", transfersM2)

	transfersM3 := []*types.Transfer{
		{Owner: "party1", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(500)}},
		{Owner: "party2", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(450)}},
	}
	tracker.UpdateFeesFromTransfers("market3", transfersM3)
	return tracker
}

func TestSnapshotRoundtripViaEngine(t *testing.T) {
	ctx := vgcontext.WithTraceID(vgcontext.WithBlockHeight(context.Background(), 100), "0xDEADBEEF")
	ctx = vgcontext.WithChainID(ctx, "chainid")
	tracker := setupDefaultTrackerForTest(t)
	now := time.Now()
	log := logging.NewTestLogger()
	timeService := stubs.NewTimeStub()
	timeService.SetTime(now)
	statsData := stats.New(log, stats.NewDefaultConfig())
	config := snp.NewDefaultConfig()
	config.Storage = "memory"
	snapshotEngine, _ := snp.New(context.Background(), &paths.DefaultPaths{}, config, log, timeService, statsData.Blockchain)
	snapshotEngine.AddProviders(tracker)
	snapshotEngine.ClearAndInitialise()
	defer snapshotEngine.Close()

	_, err := snapshotEngine.Snapshot(ctx)
	require.NoError(t, err)
	snaps, err := snapshotEngine.List()
	require.NoError(t, err)
	snap1 := snaps[0]

	trackerLoad := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	tracker.SetEligibilityChecker(&EligibilityChecker{})
	snapshotEngineLoad, _ := snp.New(context.Background(), &paths.DefaultPaths{}, config, log, timeService, statsData.Blockchain)
	snapshotEngineLoad.AddProviders(trackerLoad)
	snapshotEngineLoad.ClearAndInitialise()
	snapshotEngineLoad.ReceiveSnapshot(snap1)
	snapshotEngineLoad.ApplySnapshot(ctx)
	snapshotEngineLoad.CheckLoaded()
	defer snapshotEngineLoad.Close()

	b, err := snapshotEngine.Snapshot(ctx)
	require.NoError(t, err)
	bLoad, err := snapshotEngineLoad.Snapshot(ctx)
	require.NoError(t, err)
	require.True(t, bytes.Equal(b, bLoad))

	// now lets get some activity going and verify they still match
	tracker.MarketProposed("asset1", "market5", "meeeee")
	tracker.MarketProposed("asset2", "market6", "meeeeeee")
	trackerLoad.MarketProposed("asset1", "market5", "meeeee")
	trackerLoad.MarketProposed("asset2", "market6", "meeeeeee")

	transfersM5 := []*types.Transfer{
		{Owner: "party3", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(100)}},
		{Owner: "party3", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
		{Owner: "party3", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset1", Amount: num.NewUint(200)}},
	}
	tracker.UpdateFeesFromTransfers("market5", transfersM5)
	trackerLoad.UpdateFeesFromTransfers("market5", transfersM5)

	transfersM6 := []*types.Transfer{
		{Owner: "party4", Type: types.TransferTypeMakerFeeReceive, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(500)}},
		{Owner: "party4", Type: types.TransferTypeMakerFeePay, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(1500)}},
		{Owner: "party4", Type: types.TransferTypeLiquidityFeeDistribute, Amount: &types.FinancialAmount{Asset: "asset2", Amount: num.NewUint(1500)}},
	}
	tracker.UpdateFeesFromTransfers("market6", transfersM6)
	trackerLoad.UpdateFeesFromTransfers("market6", transfersM6)

	b, err = snapshotEngine.Snapshot(ctx)
	require.NoError(t, err)
	bLoad, err = snapshotEngineLoad.Snapshot(ctx)
	require.NoError(t, err)
	require.True(t, bytes.Equal(b, bLoad))
}

func TestMarketProposerBonusScenarios(t *testing.T) {
	tracker := common.NewMarketActivityTracker(logging.NewTestLogger(), &TestEpochEngine{})
	tracker.SetEligibilityChecker(&EligibilityChecker{})

	// setup 4 market for settlement asset1 2 of them proposed by the same proposer, and 2 markets for settlement asset 2
	tracker.MarketProposed("asset1", "market1", "me")
	tracker.MarketProposed("asset1", "market2", "me")
	tracker.MarketProposed("asset1", "market3", "me2")
	tracker.MarketProposed("asset1", "market4", "me3")
	tracker.MarketProposed("asset2", "market5", "me")
	tracker.MarketProposed("asset2", "market6", "me2")

	// no trading done so far so expect no one to be eligible for bonus
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset2", []string{}, "FURY", "zohar")))

	// market1 goes above the threshold only it should be eligible
	tracker.AddValueTraded("market1", num.NewUint(5001))
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2", "market3"}, "FURY", "zohar")))
	require.True(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	tracker.MarkPaidProposer("market1", "FURY", []string{"market1", "market2", "market3"}, "zohar")

	// now market 2 and 3 become eligible
	tracker.AddValueTraded("market2", num.NewUint(5001))
	tracker.AddValueTraded("market3", num.NewUint(5001))
	require.Equal(t, 2, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2", "market3"}, "FURY", "zohar")))

	// show that only markets 2 and 3 are now eligible with this combo
	require.False(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	tracker.MarkPaidProposer("market2", "FURY", []string{"market1", "market2", "market3"}, "zohar")
	tracker.MarkPaidProposer("market3", "FURY", []string{"market1", "market2", "market3"}, "zohar")

	// now market4 goes above the threshold but no one gets paid by this combo
	tracker.AddValueTraded("market4", num.NewUint(5001))
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market2", "market3"}, "FURY", "zohar")))
	require.False(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{"market1", "market2", "market3"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{"market1", "market2", "market3"}, "zohar"))

	// now "all" is funded by zohar
	require.Equal(t, 4, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))
	require.True(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{}, "zohar"))

	tracker.MarkPaidProposer("market1", "FURY", []string{}, "zohar")
	tracker.MarkPaidProposer("market2", "FURY", []string{}, "zohar")
	tracker.MarkPaidProposer("market3", "FURY", []string{}, "zohar")
	tracker.MarkPaidProposer("market4", "FURY", []string{}, "zohar")

	// everyone were paid so next time no one is eligible
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))

	// a new market is proposed and gets over the limit
	tracker.MarketProposed("asset1", "market7", "mememe")
	tracker.AddValueTraded("market7", num.NewUint(5001))

	// only the new market should be eligible for the "all" combo funded by zohar
	require.Equal(t, 1, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))
	require.False(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market7", "FURY", []string{}, "zohar"))
	tracker.MarkPaidProposer("market7", "FURY", []string{}, "zohar")

	// check that they are no longer eligible for this combo of all
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "zohar")))

	// check new combo
	require.Equal(t, 3, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market3", "market7"}, "FURY", "zohar")))
	require.True(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{"market1", "market3", "market7"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{"market1", "market3", "market7"}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{"market1", "market3", "market7"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{"market1", "market3", "market7"}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market7", "FURY", []string{"market1", "market3", "market7"}, "zohar"))

	tracker.MarkPaidProposer("market1", "FURY", []string{"market1", "market3", "market7"}, "zohar")
	tracker.MarkPaidProposer("market3", "FURY", []string{"market1", "market3", "market7"}, "zohar")
	tracker.MarkPaidProposer("market7", "FURY", []string{"market1", "market3", "market7"}, "zohar")

	// now that they're marked as paid check they're no longer eligible
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market3", "market7"}, "FURY", "zohar")))

	// check new asset for the same combo
	require.Equal(t, 3, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market3", "market7"}, "USDC", "zohar")))
	require.True(t, tracker.IsMarketEligibleForBonus("market1", "USDC", []string{"market1", "market3", "market7"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market2", "USDC", []string{"market1", "market3", "market7"}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market3", "USDC", []string{"market1", "market3", "market7"}, "zohar"))
	require.False(t, tracker.IsMarketEligibleForBonus("market4", "USDC", []string{"market1", "market3", "market7"}, "zohar"))
	require.True(t, tracker.IsMarketEligibleForBonus("market7", "USDC", []string{"market1", "market3", "market7"}, "zohar"))

	tracker.MarkPaidProposer("market1", "USDC", []string{"market1", "market3", "market7"}, "zohar")
	tracker.MarkPaidProposer("market3", "USDC", []string{"market1", "market3", "market7"}, "zohar")
	tracker.MarkPaidProposer("market7", "USDC", []string{"market1", "market3", "market7"}, "zohar")

	// now that they're marked as paid check they're no longer eligible
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{"market1", "market3", "market7"}, "USDC", "zohar")))

	// check new funder for the all combo
	require.Equal(t, 5, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "jeremy")))
	require.True(t, tracker.IsMarketEligibleForBonus("market1", "FURY", []string{}, "jeremy"))
	require.True(t, tracker.IsMarketEligibleForBonus("market2", "FURY", []string{}, "jeremy"))
	require.True(t, tracker.IsMarketEligibleForBonus("market3", "FURY", []string{}, "jeremy"))
	require.True(t, tracker.IsMarketEligibleForBonus("market4", "FURY", []string{}, "jeremy"))
	require.True(t, tracker.IsMarketEligibleForBonus("market7", "FURY", []string{}, "jeremy"))

	tracker.MarkPaidProposer("market1", "FURY", []string{}, "jeremy")
	tracker.MarkPaidProposer("market2", "FURY", []string{}, "jeremy")
	tracker.MarkPaidProposer("market3", "FURY", []string{}, "jeremy")
	tracker.MarkPaidProposer("market4", "FURY", []string{}, "jeremy")
	tracker.MarkPaidProposer("market7", "FURY", []string{}, "jeremy")
	require.Equal(t, 0, len(tracker.GetMarketsWithEligibleProposer("asset1", []string{}, "FURY", "jeremy")))
}
