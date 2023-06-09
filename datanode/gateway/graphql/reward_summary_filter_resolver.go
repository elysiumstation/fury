package gql

import (
	"context"
	"fmt"

	"github.com/elysiumstation/fury/libs/ptr"
	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
)

type rewardSummaryFilterResolver FuryResolverRoot

func (r *rewardSummaryFilterResolver) FromEpoch(_ context.Context, obj *v2.RewardSummaryFilter, data *int) (err error) {
	obj.FromEpoch, err = intPtrToUint64Ptr(data)
	return
}

func (r *rewardSummaryFilterResolver) ToEpoch(_ context.Context, obj *v2.RewardSummaryFilter, data *int) (err error) {
	obj.ToEpoch, err = intPtrToUint64Ptr(data)
	return
}

func intPtrToUint64Ptr(i *int) (*uint64, error) {
	iVal := ptr.UnBox(i)
	if iVal < 0 {
		return nil, fmt.Errorf("cannot convert to uint - must be positive")
	}
	return ptr.From(uint64(iVal)), nil
}
