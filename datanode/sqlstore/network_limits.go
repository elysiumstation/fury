// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package sqlstore

import (
	"context"

	"github.com/elysiumstation/fury/datanode/entities"
	"github.com/elysiumstation/fury/datanode/metrics"
	"github.com/georgysavva/scany/pgxscan"
)

type NetworkLimits struct {
	*ConnectionSource
}

func NewNetworkLimits(connectionSource *ConnectionSource) *NetworkLimits {
	return &NetworkLimits{ConnectionSource: connectionSource}
}

// Add inserts a row into the network limits table. If a row with the same fury time
// exists, that row is updated instead. (i.e. there are multiple updates of the limits
// in one block, does occur).
func (nl *NetworkLimits) Add(ctx context.Context, limits entities.NetworkLimits) error {
	defer metrics.StartSQLQuery("NetworkLimits", "Add")()
	_, err := nl.Connection.Exec(ctx, `
	INSERT INTO network_limits(
		tx_hash,
		fury_time,
		can_propose_market,
		can_propose_asset,
		propose_market_enabled,
		propose_asset_enabled,
		genesis_loaded,
		propose_market_enabled_from,
		propose_asset_enabled_from)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (fury_time) DO UPDATE SET
		can_propose_market=EXCLUDED.can_propose_market,
		can_propose_asset=EXCLUDED.can_propose_asset,
		propose_market_enabled=EXCLUDED.propose_market_enabled,
		propose_asset_enabled=EXCLUDED.propose_asset_enabled,
		genesis_loaded=EXCLUDED.genesis_loaded,
		propose_market_enabled_from=EXCLUDED.propose_market_enabled_from,
		propose_asset_enabled_from=EXCLUDED.propose_asset_enabled_from
	`,
		limits.TxHash,
		limits.FuryTime,
		limits.CanProposeMarket,
		limits.CanProposeAsset,
		limits.ProposeMarketEnabled,
		limits.ProposeAssetEnabled,
		limits.GenesisLoaded,
		limits.ProposeMarketEnabledFrom,
		limits.ProposeAssetEnabledFrom)
	return err
}

// GetLatest returns the most recent network limits.
func (nl *NetworkLimits) GetLatest(ctx context.Context) (entities.NetworkLimits, error) {
	networkLimits := entities.NetworkLimits{}
	defer metrics.StartSQLQuery("NetworkLimits", "GetLatest")()
	return networkLimits, nl.wrapE(pgxscan.Get(ctx, nl.Connection, &networkLimits,
		`SELECT * FROM network_limits ORDER BY fury_time DESC limit 1;`))
}
