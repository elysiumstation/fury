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

package steps

import (
	"strconv"

	"github.com/elysiumstation/fury/core/integration/steps/market"
)

func OracleSpecSettlementDataDecimals(config *market.Config, name string, settlementDP string) error {
	dp, err := strconv.ParseUint(settlementDP, 10, 32)
	if err != nil {
		return err
	}
	config.OracleConfigs.SetSettlementDataDP(name, uint32(dp))
	return nil
}
