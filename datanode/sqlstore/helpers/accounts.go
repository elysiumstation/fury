package helpers

import (
	"context"
	"testing"

	"github.com/elysiumstation/fury/core/types"
	"github.com/elysiumstation/fury/datanode/entities"
	"github.com/elysiumstation/fury/datanode/sqlstore"
	"github.com/stretchr/testify/require"
)

func AddTestAccount(t *testing.T,
	ctx context.Context,
	accountStore *sqlstore.Accounts,
	party entities.Party,
	asset entities.Asset,
	accountType types.AccountType,
	block entities.Block,
) entities.Account {
	t.Helper()
	account := entities.Account{
		PartyID:  party.ID,
		AssetID:  asset.ID,
		MarketID: entities.MarketID(GenerateID()),
		Type:     accountType,
		FuryTime: block.FuryTime,
	}

	err := accountStore.Add(ctx, &account)
	require.NoError(t, err)
	return account
}

func AddTestAccountWithTxHash(t *testing.T,
	ctx context.Context,
	accountStore *sqlstore.Accounts,
	party entities.Party,
	asset entities.Asset,
	accountType types.AccountType,
	block entities.Block,
	txHash entities.TxHash,
) entities.Account {
	t.Helper()
	account := entities.Account{
		PartyID:  party.ID,
		AssetID:  asset.ID,
		MarketID: entities.MarketID(GenerateID()),
		Type:     accountType,
		FuryTime: block.FuryTime,
		TxHash:   txHash,
	}

	err := accountStore.Add(ctx, &account)
	require.NoError(t, err)
	return account
}

func AddTestAccountWithMarketAndType(t *testing.T,
	ctx context.Context,
	accountStore *sqlstore.Accounts,
	party entities.Party,
	asset entities.Asset,
	block entities.Block,
	market entities.MarketID,
	accountType types.AccountType,
) entities.Account {
	t.Helper()
	account := entities.Account{
		PartyID:  party.ID,
		AssetID:  asset.ID,
		MarketID: market,
		Type:     accountType,
		FuryTime: block.FuryTime,
	}

	err := accountStore.Add(ctx, &account)
	require.NoError(t, err)
	return account
}
