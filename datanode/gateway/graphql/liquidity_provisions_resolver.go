package gql

import (
	"context"
	"strconv"

	"github.com/elysiumstation/fury/datanode/furytime"
	types "github.com/elysiumstation/fury/protos/fury"
)

// LiquidityProvision resolver

type myLiquidityProvisionResolver FuryResolverRoot

func (r *myLiquidityProvisionResolver) Version(_ context.Context, obj *types.LiquidityProvision) (string, error) {
	return strconv.FormatUint(obj.Version, 10), nil
}

func (r *myLiquidityProvisionResolver) Party(_ context.Context, obj *types.LiquidityProvision) (*types.Party, error) {
	return &types.Party{Id: obj.PartyId}, nil
}

func (r *myLiquidityProvisionResolver) CreatedAt(ctx context.Context, obj *types.LiquidityProvision) (string, error) {
	return furytime.Format(furytime.UnixNano(obj.CreatedAt)), nil
}

func (r *myLiquidityProvisionResolver) UpdatedAt(ctx context.Context, obj *types.LiquidityProvision) (*string, error) {
	var updatedAt *string
	if obj.UpdatedAt > 0 {
		t := furytime.Format(furytime.UnixNano(obj.UpdatedAt))
		updatedAt = &t
	}
	return updatedAt, nil
}

func (r *myLiquidityProvisionResolver) Market(ctx context.Context, obj *types.LiquidityProvision) (*types.Market, error) {
	return r.r.getMarketByID(ctx, obj.MarketId)
}

func (r *myLiquidityProvisionResolver) CommitmentAmount(ctx context.Context, obj *types.LiquidityProvision) (string, error) {
	return obj.CommitmentAmount, nil
}
