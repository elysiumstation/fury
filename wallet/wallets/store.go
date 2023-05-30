package wallets

import (
	"fmt"

	"github.com/elysiumstation/fury/paths"
	wstorev1 "github.com/elysiumstation/fury/wallet/wallet/store/v1"
)

// InitialiseStore builds a wallet Store specifically for users wallets.
func InitialiseStore(furyHome string, withFileWatcher bool) (*wstorev1.FileStore, error) {
	p := paths.New(furyHome)
	return InitialiseStoreFromPaths(p, withFileWatcher)
}

// InitialiseStoreFromPaths builds a wallet Store specifically for users wallets.
func InitialiseStoreFromPaths(furyPaths paths.Paths, withFileWatcher bool) (*wstorev1.FileStore, error) {
	walletsHome, err := furyPaths.CreateDataPathFor(paths.WalletsDataHome)
	if err != nil {
		return nil, fmt.Errorf("couldn't get wallets data home path: %w", err)
	}
	return wstorev1.InitialiseStore(walletsHome, withFileWatcher)
}
