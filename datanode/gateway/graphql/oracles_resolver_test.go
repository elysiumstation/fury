package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	furypb "github.com/elysiumstation/fury/protos/fury"
	v1 "github.com/elysiumstation/fury/protos/fury/data/v1"
)

func Test_oracleSpecResolver_DataSourceSpec(t *testing.T) {
	type args struct {
		in0 context.Context
		obj *furypb.OracleSpec
	}
	tests := []struct {
		name    string
		o       oracleSpecResolver
		args    args
		wantJsn string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success: DataSourceDefinition_External",
			args: args{
				obj: &furypb.OracleSpec{
					ExternalDataSourceSpec: &furypb.ExternalDataSourceSpec{
						Spec: &furypb.DataSourceSpec{
							Status: furypb.DataSourceSpec_STATUS_ACTIVE,
							Data: &furypb.DataSourceDefinition{
								SourceType: &furypb.DataSourceDefinition_External{
									External: &furypb.DataSourceDefinitionExternal{
										SourceType: &furypb.DataSourceDefinitionExternal_Oracle{
											Oracle: &furypb.DataSourceSpecConfiguration{
												Signers: []*v1.Signer{
													{
														Signer: &v1.Signer_PubKey{
															PubKey: &v1.PubKey{
																Key: "key",
															},
														},
													}, {
														Signer: &v1.Signer_EthAddress{
															EthAddress: &v1.ETHAddress{
																Address: "address",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantJsn: `{"spec":{"id":"","createdAt":0,"updatedAt":null,"data":{"sourceType":{"sourceType":{"signers":[{"Signer":{"PubKey":{"key":"key"}}},{"Signer":{"EthAddress":{"address":"address"}}}]}}},"status":"STATUS_ACTIVE"}}`,
			wantErr: assert.NoError,
		}, {
			name: "success: DataSourceDefinition_Internal",
			args: args{
				obj: &furypb.OracleSpec{
					ExternalDataSourceSpec: &furypb.ExternalDataSourceSpec{
						Spec: &furypb.DataSourceSpec{
							Status: furypb.DataSourceSpec_STATUS_ACTIVE,
							Data: &furypb.DataSourceDefinition{
								SourceType: &furypb.DataSourceDefinition_Internal{
									Internal: &furypb.DataSourceDefinitionInternal{
										SourceType: &furypb.DataSourceDefinitionInternal_Time{
											Time: &furypb.DataSourceSpecConfigurationTime{
												Conditions: []*v1.Condition{
													{
														Operator: 12,
														Value:    "blah",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantJsn: `{"spec":{"id":"","createdAt":0,"updatedAt":null,"data":{"sourceType":{"sourceType":{"conditions":[{"operator":12,"value":"blah"}]}}},"status":"STATUS_ACTIVE"}}`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.DataSourceSpec(tt.args.in0, tt.args.obj)
			if !tt.wantErr(t, err, fmt.Sprintf("DataSourceSpec(%v, %v)", tt.args.in0, tt.args.obj)) {
				return
			}

			gotJsn, _ := json.Marshal(got)
			assert.JSONEqf(t, tt.wantJsn, string(gotJsn), "mismatch:\n\twant: %s \n\tgot: %s", tt.wantJsn, string(gotJsn))
		})
	}
}
