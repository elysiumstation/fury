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
	"encoding/json"
	"fmt"
	"time"

	"github.com/elysiumstation/fury/core/types"
	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	datapb "github.com/elysiumstation/fury/protos/fury/data/v1"

	furypb "github.com/elysiumstation/fury/protos/fury"
)

type ExternalDataSourceSpec struct {
	Spec *DataSourceSpec
}

func (s *ExternalDataSourceSpec) ToProto() *furypb.ExternalDataSourceSpec {
	return &furypb.ExternalDataSourceSpec{
		Spec: s.Spec.ToProto(),
	}
}

func ExternalDataSourceSpecFromProto(spec *furypb.ExternalDataSourceSpec, txHash TxHash, furyTime time.Time) (*ExternalDataSourceSpec, error) {
	if spec.Spec != nil {
		ds, err := DataSourceSpecFromProto(spec.Spec, txHash, furyTime)
		if err != nil {
			return nil, err
		}

		return &ExternalDataSourceSpec{
			Spec: ds,
		}, nil
	}

	return &ExternalDataSourceSpec{
		Spec: &DataSourceSpec{},
	}, nil
}

type (
	_Spec  struct{}
	SpecID = ID[_Spec]
)

type (
	Signer  []byte
	Signers = []Signer
)

type DataSourceSpecConfiguration struct {
	Signers Signers
	Filters []Filter
}

type DataSourceSpec struct {
	ID        SpecID
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *DataSourceDefinition
	Status    DataSourceSpecStatus
	TxHash    TxHash
	FuryTime  time.Time
}

type DataSourceSpecRaw struct {
	ID        SpecID
	CreatedAt time.Time
	UpdatedAt time.Time
	Signers   Signers
	Filters   []Filter
	Status    DataSourceSpecStatus
	TxHash    TxHash
	FuryTime  time.Time
}

func DataSourceSpecFromProto(spec *furypb.DataSourceSpec, txHash TxHash, furyTime time.Time) (*DataSourceSpec, error) {
	if spec != nil {
		id := SpecID(spec.Id)

		data := &DataSourceDefinition{}
		if spec.Data != nil {
			filters := FiltersFromProto(spec.Data.GetFilters())
			signers, err := SerializeSigners(types.SignersFromProto(spec.Data.GetSigners()))
			if err != nil {
				return nil, err
			}

			data.External = &DataSourceDefinitionExternal{
				Signers: signers,
				Filters: filters,
			}
		}

		return &DataSourceSpec{
			ID:        id,
			CreatedAt: time.Unix(0, spec.CreatedAt),
			UpdatedAt: time.Unix(0, spec.UpdatedAt),
			Data:      data,
			Status:    DataSourceSpecStatus(spec.Status),
			TxHash:    txHash,
			FuryTime:  furyTime,
		}, nil
	}
	return nil, nil
}

func (ds *DataSourceSpec) ToProto() *furypb.DataSourceSpec {
	filters := []*datapb.Filter{}
	signers := []*datapb.Signer{}

	if ds.Data != nil {
		desSigners := DeserializeSigners(ds.Data.External.Signers)
		signers = types.SignersIntoProto(desSigners)
		filters = filtersToProto(ds.Data.External.Filters)
	}

	return &furypb.DataSourceSpec{
		Id:        ds.ID.String(),
		CreatedAt: ds.CreatedAt.UnixNano(),
		UpdatedAt: ds.UpdatedAt.UnixNano(),
		Data: furypb.NewDataSourceDefinition(
			furypb.DataSourceDefinitionTypeExt,
		).SetOracleConfig(
			&furypb.DataSourceSpecConfiguration{
				Signers: signers,
				Filters: filters,
			},
		),
		Status: furypb.DataSourceSpec_Status(ds.Status),
	}
}

func (ds *DataSourceSpec) ToOracleProto() *furypb.OracleSpec {
	return &furypb.OracleSpec{
		ExternalDataSourceSpec: &furypb.ExternalDataSourceSpec{
			Spec: ds.ToProto(),
		},
	}
}

func (ds DataSourceSpec) Cursor() *Cursor {
	return NewCursor(DataSourceSpecCursor{ds.FuryTime, ds.ID}.String())
}

func (ds DataSourceSpec) ToOracleProtoEdge(_ ...any) (*v2.OracleSpecEdge, error) {
	return &v2.OracleSpecEdge{
		Node:   ds.ToOracleProto(),
		Cursor: ds.Cursor().Encode(),
	}, nil
}

func SerializeSigners(signers []*types.Signer) (Signers, error) {
	if len(signers) > 0 {
		sigList := Signers{}

		for _, signer := range signers {
			data, err := signer.Serialize()
			if err != nil {
				return nil, err
			}
			sigList = append(sigList, data)
		}

		return sigList, nil
	}

	return Signers{}, nil
}

func DeserializeSigners(data Signers) []*types.Signer {
	if len(data) > 0 {
		signers := []*types.Signer{}
		for _, s := range data {
			signer := types.DeserializeSigner(s)
			signers = append(signers, signer)
		}

		return signers
	}

	return nil
}

type DataSourceSpecCursor struct {
	FuryTime time.Time `json:"furyTime"`
	ID       SpecID    `json:"id"`
}

func (ds DataSourceSpecCursor) String() string {
	bs, err := json.Marshal(ds)
	if err != nil {
		panic(fmt.Errorf("could not marshal oracle spec cursor: %w", err))
	}
	return string(bs)
}

func (ds *DataSourceSpecCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}
	return json.Unmarshal([]byte(cursorString), ds)
}
