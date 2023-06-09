package ethcall

import (
	"encoding/hex"
	"fmt"

	"github.com/elysiumstation/fury/protos/fury"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/types/known/structpb"
)

type Spec struct {
	Call
	Trigger
}

func NewSpec(call Call, trigger Trigger) Spec {
	return Spec{
		Call:    call,
		Trigger: trigger,
	}
}

func (s Spec) Hash() []byte {
	hashFunc := sha3.New256()
	hashFunc.Write(s.Call.Hash())
	hashFunc.Write(s.Trigger.Hash())
	return hashFunc.Sum(nil)
}

func (s Spec) HashHex() string {
	return hex.EncodeToString(s.Hash())
}

func (s Spec) ToProto() (*fury.DataSourceDefinition, error) {
	args, err := s.Args()
	if err != nil {
		return nil, fmt.Errorf("failed to get eth call args: %w", err)
	}

	jsonArgs, err := AnyArgsToJson(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal eth call args: %w", err)
	}

	argsPBValue := []*structpb.Value{}
	for _, arg := range jsonArgs {
		v := structpb.Value{}
		err := v.UnmarshalJSON([]byte(arg))
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal arg json '%s': %w", arg, err)
		}
		argsPBValue = append(argsPBValue, &v)
	}

	abiPBList := structpb.ListValue{}
	err = abiPBList.UnmarshalJSON(s.abiJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal abi json: %w", err)
	}

	return &fury.DataSourceDefinition{
		SourceType: &fury.DataSourceDefinition_External{
			External: &fury.DataSourceDefinitionExternal{
				SourceType: &fury.DataSourceDefinitionExternal_EthCall{
					EthCall: &fury.EthCallSpec{
						Address: s.address.Hex(),
						Abi:     &abiPBList,
						Method:  s.method,
						Args:    argsPBValue,
						Trigger: s.Trigger.ToProto(),
					},
				},
			},
		},
	}, nil
}

func NewSpecFromProto(proto *fury.DataSourceDefinition) (Spec, error) {
	if proto == nil {
		return Spec{}, fmt.Errorf("null data source definition")
	}

	externalProto := proto.GetExternal()
	if externalProto == nil {
		return Spec{}, fmt.Errorf("not an external data source")
	}

	ethCallProto := externalProto.GetEthCall()
	if ethCallProto == nil {
		return Spec{}, fmt.Errorf("not an eth call data source")
	}

	// Get args out of proto 'struct' format into JSON
	jsonArgs := []string{}
	for _, protoArg := range ethCallProto.Args {
		jsonArg, err := protoArg.MarshalJSON()
		if err != nil {
			return Spec{}, fmt.Errorf("unable to marshal args from proto to json: %w", err)
		}
		jsonArgs = append(jsonArgs, string(jsonArg))
	}

	abiJson, err := ethCallProto.Abi.MarshalJSON()
	if err != nil {
		return Spec{}, fmt.Errorf("unable to marshal abi: %w", err)
	}

	// Convert JSON args to go types using ABI
	args, err := JsonArgsToAny(ethCallProto.Method, jsonArgs, string(abiJson))
	if err != nil {
		return Spec{}, fmt.Errorf("unable to deserialize args: %w", err)
	}

	call, err := NewCall(ethCallProto.Method, args, ethCallProto.Address, abiJson)
	if err != nil {
		return Spec{}, fmt.Errorf("unable to create call: %w", err)
	}

	trigger, err := TriggerFromProto(ethCallProto.Trigger)
	if err != nil {
		return Spec{}, fmt.Errorf("unable to create trigger: %w", err)
	}

	return Spec{
		Call:    call,
		Trigger: trigger,
	}, nil
}
