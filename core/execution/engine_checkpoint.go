package execution

import (
	"context"
	"sort"

	"github.com/elysiumstation/fury/core/types"
	"github.com/elysiumstation/fury/libs/proto"
	checkpoint "github.com/elysiumstation/fury/protos/fury/checkpoint/v1"
)

func (e *Engine) Name() types.CheckpointName {
	return types.ExecutionCheckpoint
}

func (e *Engine) Checkpoint() ([]byte, error) {
	for id, mkt := range e.markets {
		state := mkt.GetCPState()
		e.marketCPStates[id] = state
	}
	data := make([]*types.CPMarketState, 0, len(e.marketCPStates))
	// @TODO prune the marketCPStates
	for _, s := range e.marketCPStates {
		data = append(data, s)
	}
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].ID < data[j].ID
	})
	wrapped := types.ExecutionState{
		Data: data,
	}
	cpData, err := proto.Marshal(wrapped.IntoProto())
	return cpData, err
}

func (e *Engine) Load(ctx context.Context, data []byte) error {
	wrapper := checkpoint.ExecutionState{}
	if err := proto.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	e.marketCPStates = make(map[string]*types.CPMarketState, len(wrapper.Data))
	// for now, restore all pending markets as though their state is valid for the full TTL window
	for _, mcp := range wrapper.Data {
		e.marketCPStates[mcp.Id] = types.NewMarketStateFromProto(mcp)
	}
	return nil
}
