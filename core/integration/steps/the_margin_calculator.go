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
	"github.com/cucumber/godog"

	"github.com/elysiumstation/fury/core/integration/steps/market"
	types "github.com/elysiumstation/fury/protos/fury"
)

func TheMarginCalculator(config *market.Config, name string, table *godog.Table) error {
	row := marginCalculatorRow{row: parseMarginCalculatorTable(table)}

	return config.MarginCalculators.Add(name, &types.MarginCalculator{
		ScalingFactors: &types.ScalingFactors{
			SearchLevel:       row.searchLevelFactor(),
			InitialMargin:     row.initialMarginFactor(),
			CollateralRelease: row.collateralReleaseFactor(),
		},
	})
}

func parseMarginCalculatorTable(table *godog.Table) RowWrapper {
	return StrictParseFirstRow(table, []string{
		"release factor",
		"initial factor",
		"search factor",
	}, []string{})
}

type marginCalculatorRow struct {
	row RowWrapper
}

func (r marginCalculatorRow) collateralReleaseFactor() float64 {
	return r.row.MustF64("release factor")
}

func (r marginCalculatorRow) initialMarginFactor() float64 {
	return r.row.MustF64("initial factor")
}

func (r marginCalculatorRow) searchLevelFactor() float64 {
	return r.row.MustF64("search factor")
}