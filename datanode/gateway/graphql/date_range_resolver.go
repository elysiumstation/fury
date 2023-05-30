package gql

import (
	"context"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
)

type dateRangeResolver FuryResolverRoot

func (r *dateRangeResolver) Start(ctx context.Context, obj *v2.DateRange, data *int64) error {
	obj.StartTimestamp = data
	return nil
}

func (r *dateRangeResolver) End(ctx context.Context, obj *v2.DateRange, data *int64) error {
	obj.EndTimestamp = data
	return nil
}
