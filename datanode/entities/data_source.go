package entities

import (
	"time"

	"github.com/elysiumstation/fury/core/types"
	v1 "github.com/elysiumstation/fury/protos/fury/data/v1"
)

type DataSourceDefinitionExternal struct {
	Signers Signers
	Filters []Filter
}

type DataSourceDefinitionInternal struct {
	Time time.Time
}

type DataSourceDefinition struct {
	Type     int
	External *DataSourceDefinitionExternal
	Internal *DataSourceDefinitionInternal
}

func (s *DataSourceDefinition) GetSigners() []*v1.Signer {
	return types.SignersIntoProto(DeserializeSigners(s.External.Signers))
}

func (s *DataSourceDefinition) GetFilters() []*v1.Filter {
	return filtersToProto(s.External.Filters)
}
