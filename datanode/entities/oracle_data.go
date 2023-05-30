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
	"time"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	furypb "github.com/elysiumstation/fury/protos/fury"
)

type OracleData struct {
	ExternalData *ExternalData
}

func OracleDataFromProto(data *furypb.OracleData, txHash TxHash, furyTime time.Time, seqNum uint64) (*OracleData, error) {
	extData, err := ExternalDataFromProto(data.ExternalData, txHash, furyTime, seqNum)
	if err != nil {
		return nil, err
	}

	return &OracleData{
		ExternalData: extData,
	}, nil
}

func (od OracleData) ToProto() *furypb.OracleData {
	return &furypb.OracleData{
		ExternalData: od.ExternalData.ToProto(),
	}
}

func (od OracleData) Cursor() *Cursor {
	return od.ExternalData.Cursor()
}

func (od OracleData) ToProtoEdge(_ ...any) (*v2.OracleDataEdge, error) {
	tp, err := od.ExternalData.ToOracleProtoEdge()
	if err != nil {
		return nil, err
	}

	return tp, nil
}

type OracleDataCursor = ExternalDataCursor
