package gql

import (
	"context"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	"github.com/elysiumstation/fury/protos/fury"
)

type orderFilterResolver FuryResolverRoot

func (o orderFilterResolver) Status(ctx context.Context, obj *v2.OrderFilter, data []fury.Order_Status) error {
	obj.Statuses = data
	return nil
}

func (o orderFilterResolver) TimeInForce(ctx context.Context, obj *v2.OrderFilter, data []fury.Order_TimeInForce) error {
	obj.TimeInForces = data
	return nil
}
