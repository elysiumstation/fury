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

package governance

import (
	"context"

	"github.com/elysiumstation/fury/core/events"
	"github.com/elysiumstation/fury/core/execution"
	"github.com/elysiumstation/fury/core/netparams"
	"github.com/elysiumstation/fury/core/types"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/protos/fury"
	checkpointpb "github.com/elysiumstation/fury/protos/fury/checkpoint/v1"

	"github.com/elysiumstation/fury/libs/proto"
)

type enactmentTime struct {
	current         int64
	shouldNotVerify bool
}

func (e *Engine) Name() types.CheckpointName {
	return types.GovernanceCheckpoint
}

func (e *Engine) Checkpoint() ([]byte, error) {
	if len(e.enactedProposals) == 0 {
		return nil, nil
	}
	cp := &checkpointpb.Proposals{
		Proposals: e.getCheckpointProposals(),
	}
	return proto.Marshal(cp)
}

func (e *Engine) Load(ctx context.Context, data []byte) error {
	cp := &checkpointpb.Proposals{}
	if err := proto.Unmarshal(data, cp); err != nil {
		return err
	}

	evts := make([]events.Event, 0, len(cp.Proposals))
	now := e.timeService.GetTimeNow()
	minEnact, err := e.netp.GetDuration(netparams.GovernanceProposalMarketMinEnact)
	if err != nil {
		e.log.Panic("failed to get proposal market min enactment duration from network parameter")
	}
	minAuctionDuration, err := e.netp.GetDuration(netparams.MarketAuctionMinimumDuration)
	if err != nil {
		e.log.Panic("failed to get proposal market min auction duration from network parameter")
	}
	duration := minEnact
	// we have to choose the max between minEnact and minAuctionDuration otherwise we won't be able to submit the market successfully
	if int64(minEnact) < int64(minAuctionDuration) {
		duration = minAuctionDuration
	}

	latestUpdateMarketProposals := map[string]*types.Proposal{}
	updatedMarketIDs := []string{}
	for _, p := range cp.Proposals {
		prop, err := types.ProposalFromProto(p)
		if err != nil {
			return err
		}

		switch prop.Terms.Change.GetTermType() {
		case types.ProposalTermsTypeNewMarket:
			enct := &enactmentTime{}
			// if the proposal is for a new market we want to restore it such that it will be in opening auction
			if p.Terms.EnactmentTimestamp <= now.Unix() {
				prop.Terms.EnactmentTimestamp = now.Add(duration).Unix()
				enct.shouldNotVerify = true
			}
			enct.current = prop.Terms.EnactmentTimestamp
			toSubmit, err := e.intoToSubmit(ctx, prop, enct)
			if err != nil {
				e.log.Panic("Failed to convert proposal into market", logging.Error(err))
			}
			nm := toSubmit.NewMarket()
			err = e.markets.RestoreMarket(ctx, nm.Market())
			if err != nil {
				if err == execution.ErrMarketDoesNotExist {
					// market has been settled, we don't care
					continue
				}
				// any other error, panic
				e.log.Panic("failed to restore market from checkpoint", logging.Market(*nm.Market()), logging.Error(err))
			}

			if err := e.markets.StartOpeningAuction(ctx, prop.ID); err != nil {
				e.log.Panic("failed to start opening auction for market", logging.String("market-id", prop.ID), logging.Error(err))
			}
		case types.ProposalTermsTypeUpdateMarket:
			marketID := prop.Terms.GetUpdateMarket().MarketID
			updatedMarketIDs = append(updatedMarketIDs, marketID)
			last, ok := latestUpdateMarketProposals[marketID]
			if !ok || prop.Terms.EnactmentTimestamp > last.Terms.EnactmentTimestamp {
				latestUpdateMarketProposals[marketID] = prop
			}
		}

		evts = append(evts, events.NewProposalEvent(ctx, *prop))
		e.enactedProposals = append(e.enactedProposals, &proposal{
			Proposal: prop,
		})
	}

	for _, v := range updatedMarketIDs {
		p := latestUpdateMarketProposals[v]
		mkt, _, err := e.updatedMarketFromProposal(&proposal{Proposal: p})
		if err != nil {
			continue
		}
		e.markets.UpdateMarket(ctx, mkt)
	}

	// send events for restored proposals
	e.broker.SendBatch(evts)
	// @TODO ensure OnTick is called
	return nil
}

func (e *Engine) getCheckpointProposals() []*fury.Proposal {
	ret := make([]*fury.Proposal, 0, len(e.enactedProposals))
	for _, p := range e.enactedProposals {
		switch p.Terms.Change.GetTermType() {
		case types.ProposalTermsTypeNewMarket:
			mktState, err := e.markets.GetMarketState(p.ID)
			// if the market is missing from the execution engine it means it's been already cancelled or settled or rejected
			if err == types.ErrInvalidMarketID {
				e.log.Info("not saving market proposal to checkpoint - market has already been removed", logging.String("market-id", p.ID))
				continue
			}
			if mktState == types.MarketStateTradingTerminated {
				e.log.Info("not saving market proposal to checkpoint ", logging.String("market-id", p.ID), logging.String("market-state", mktState.String()))
				continue
			}
		case types.ProposalTermsTypeUpdateMarket:
			mktState, err := e.markets.GetMarketState(p.MarketUpdate().MarketID)
			// if the market is missing from the execution engine it means it's been already cancelled or settled or rejected
			if err == types.ErrInvalidMarketID {
				e.log.Info("not saving market update proposal to checkpoint - market has already been removed", logging.String("market-id", p.ID))
				continue
			}
			if mktState == types.MarketStateTradingTerminated {
				e.log.Info("not saving market update proposal to checkpoint ", logging.String("market-id", p.ID), logging.String("market-state", mktState.String()))
				continue
			}
		}

		ret = append(ret, p.IntoProto())
	}
	return ret
}
