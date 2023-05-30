package types

import (
	furypb "github.com/elysiumstation/fury/protos/fury"
)

type DataSourceDefinitionExternal struct {
	SourceType dataSourceType
}

func (e *DataSourceDefinitionExternal) isDataSourceType() {}

func (e *DataSourceDefinitionExternal) oneOfProto() interface{} {
	return e.IntoProto()
}

// /
// IntoProto tries to return the base proto object from DataSourceDefinitionExternal.
func (e *DataSourceDefinitionExternal) IntoProto() *furypb.DataSourceDefinitionExternal {
	ds := &furypb.DataSourceDefinitionExternal{}

	if e.SourceType != nil {
		switch dsn := e.SourceType.oneOfProto().(type) {
		case *furypb.DataSourceDefinitionExternal_Oracle:
			ds = &furypb.DataSourceDefinitionExternal{
				SourceType: dsn,
			}
		}
	}

	return ds
}

func (e *DataSourceDefinitionExternal) String() string {
	if e.SourceType != nil {
		return e.SourceType.String()
	}

	return ""
}

func (e *DataSourceDefinitionExternal) DeepClone() dataSourceType {
	if e.SourceType != nil {
		return e.SourceType.DeepClone()
	}

	return nil
}

// /
// DataSourceDefinitionExternalFromProto tries to build the DataSourceDefinitionExternal object
// from the given proto object..
func DataSourceDefinitionExternalFromProto(protoConfig *furypb.DataSourceDefinitionExternal) *DataSourceDefinitionExternal {
	ds := &DataSourceDefinitionExternal{
		SourceType: &DataSourceDefinitionExternalOracle{},
	}

	if protoConfig != nil {
		if protoConfig.SourceType != nil {
			switch tp := protoConfig.SourceType.(type) {
			case *furypb.DataSourceDefinitionExternal_Oracle:
				ds.SourceType = DataSourceDefinitionExternalOracleFromProto(tp)
			}
		}
	}

	return ds
}
