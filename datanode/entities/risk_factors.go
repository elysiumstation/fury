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

package entities

import (
	"fmt"
	"time"

	"github.com/elysiumstation/fury/protos/fury"

	"github.com/shopspring/decimal"
)

type RiskFactor struct {
	MarketID MarketID
	Short    decimal.Decimal
	Long     decimal.Decimal
	TxHash   TxHash
	FuryTime time.Time
}

func RiskFactorFromProto(factor *fury.RiskFactor, txHash TxHash, furyTime time.Time) (*RiskFactor, error) {
	var short, long decimal.Decimal
	var err error

	if short, err = decimal.NewFromString(factor.Short); err != nil {
		return nil, fmt.Errorf("invalid value for short: %s - %v", factor.Short, err)
	}

	if long, err = decimal.NewFromString(factor.Long); err != nil {
		return nil, fmt.Errorf("invalid value for long: %s - %v", factor.Long, err)
	}

	return &RiskFactor{
		MarketID: MarketID(factor.Market),
		Short:    short,
		Long:     long,
		TxHash:   txHash,
		FuryTime: furyTime,
	}, nil
}

func (rf *RiskFactor) ToProto() *fury.RiskFactor {
	return &fury.RiskFactor{
		Market: rf.MarketID.String(),
		Short:  rf.Short.String(),
		Long:   rf.Long.String(),
	}
}
