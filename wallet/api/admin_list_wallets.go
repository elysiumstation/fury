package api

import (
	"context"
	"fmt"

	"github.com/elysiumstation/fury/libs/jsonrpc"
)

type AdminListWalletsResult struct {
	Wallets []string `json:"wallets"`
}

type AdminListWallets struct {
	walletStore WalletStore
}

// Handle list all the wallets present on the computer.
func (h *AdminListWallets) Handle(ctx context.Context, rawParams jsonrpc.Params) (jsonrpc.Result, *jsonrpc.ErrorDetails) {
	wallets, err := h.walletStore.ListWallets(ctx)
	if err != nil {
		return nil, InternalError(fmt.Errorf("could not list the wallets: %w", err))
	}

	return AdminListWalletsResult{
		Wallets: wallets,
	}, nil
}

func NewAdminListWallets(
	walletStore WalletStore,
) *AdminListWallets {
	return &AdminListWallets{
		walletStore: walletStore,
	}
}
