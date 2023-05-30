package gql

import (
	"context"

	fury "github.com/elysiumstation/fury/protos/fury"
)

type ledgerEntryResolver FuryResolverRoot

func (le ledgerEntryResolver) FromAccountID(ctx context.Context, obj *fury.LedgerEntry) (*fury.AccountDetails, error) {
	return obj.FromAccount, nil
}

func (le ledgerEntryResolver) ToAccountID(ctx context.Context, obj *fury.LedgerEntry) (*fury.AccountDetails, error) {
	return obj.ToAccount, nil
}
