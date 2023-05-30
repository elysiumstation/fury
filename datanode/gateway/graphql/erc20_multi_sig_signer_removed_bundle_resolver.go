package gql

import (
	"context"
	"fmt"

	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
)

type erc20MultiSigSignerRemovedBundleResolver FuryResolverRoot

func (e erc20MultiSigSignerRemovedBundleResolver) Timestamp(ctx context.Context, obj *v2.ERC20MultiSigSignerRemovedBundle) (string, error) {
	return fmt.Sprint(obj.Timestamp), nil
}
