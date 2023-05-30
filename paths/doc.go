package paths

import "fmt"

const (
	// LongestPathNameLen is the length of the longest path name. It is used
	// for text formatting.
	LongestPathNameLen = 35
)

type ListPathsResponse struct {
	CachePaths  map[string]string `json:"cachePaths"`
	ConfigPaths map[string]string `json:"configPaths"`
	DataPaths   map[string]string `json:"dataPaths"`
	StatePaths  map[string]string `json:"statePaths"`
}

func List(furyPaths Paths) *ListPathsResponse {
	// Some files don't support custom fury home. So we need the default anyway.
	defaultPath := New("")

	return &ListPathsResponse{
		CachePaths: map[string]string{
			"DataNodeCacheHome": furyPaths.CachePathFor(DataNodeCacheHome),
		},
		ConfigPaths: map[string]string{
			"DataNodeConfigHome":              furyPaths.ConfigPathFor(DataNodeConfigHome),
			"DataNodeDefaultConfigFile":       furyPaths.ConfigPathFor(DataNodeDefaultConfigFile),
			"FaucetConfigHome":                furyPaths.ConfigPathFor(FaucetConfigHome),
			"FaucetDefaultConfigFile":         furyPaths.ConfigPathFor(FaucetDefaultConfigFile),
			"NodeConfigHome":                  furyPaths.ConfigPathFor(NodeConfigHome),
			"NodeDefaultConfigFile":           furyPaths.ConfigPathFor(NodeDefaultConfigFile),
			"NodeWalletsConfigFile":           furyPaths.ConfigPathFor(NodeWalletsConfigFile),
			"WalletCLIConfigHome":             furyPaths.ConfigPathFor(WalletCLIConfigHome),
			"WalletCLIDefaultConfigFile":      furyPaths.ConfigPathFor(WalletCLIDefaultConfigFile),
			"WalletAppConfigHome":             defaultPath.ConfigPathFor(WalletAppConfigHome),
			"WalletAppFairgroundConfigFile":   defaultPath.ConfigPathFor(WalletAppFairgroundConfigFile),
			"WalletAppDefaultConfigFile":      defaultPath.ConfigPathFor(WalletAppDefaultConfigFile),
			"WalletServiceConfigHome":         furyPaths.ConfigPathFor(WalletServiceConfigHome),
			"WalletServiceDefaultConfigFile":  furyPaths.ConfigPathFor(WalletServiceDefaultConfigFile),
			"WalletServiceNetworksConfigHome": furyPaths.ConfigPathFor(WalletServiceNetworksConfigHome),
		},
		DataPaths: map[string]string{
			"NodeDataHome":                       furyPaths.DataPathFor(NodeDataHome),
			"NodeWalletsDataHome":                furyPaths.DataPathFor(NodeWalletsDataHome),
			"FuryNodeWalletsDataHome":            furyPaths.DataPathFor(FuryNodeWalletsDataHome),
			"EthereumNodeWalletsDataHome":        furyPaths.DataPathFor(EthereumNodeWalletsDataHome),
			"FaucetDataHome":                     furyPaths.DataPathFor(FaucetDataHome),
			"FaucetWalletsDataHome":              furyPaths.DataPathFor(FaucetWalletsDataHome),
			"WalletsDataHome":                    furyPaths.DataPathFor(WalletsDataHome),
			"WalletServiceDataHome":              furyPaths.DataPathFor(WalletServiceDataHome),
			"WalletServiceAPITokensDataFile":     furyPaths.DataPathFor(WalletServiceAPITokensDataFile),
			"WalletServiceSessionTokensDataFile": furyPaths.DataPathFor(WalletServiceSessionTokensDataFile),
			"WalletServiceRSAKeysDataHome":       furyPaths.DataPathFor(WalletServiceRSAKeysDataHome),
			"WalletServicePublicRSAKeyDataFile":  furyPaths.DataPathFor(WalletServicePublicRSAKeyDataFile),
			"WalletServicePrivateRSAKeyDataFile": furyPaths.DataPathFor(WalletServicePrivateRSAKeyDataFile),
		},
		StatePaths: map[string]string{
			"DataNodeStateHome":                      furyPaths.StatePathFor(DataNodeStateHome),
			"DataNodeAutoCertHome":                   furyPaths.StatePathFor(DataNodeAutoCertHome),
			"DataNodeLogsHome":                       furyPaths.StatePathFor(DataNodeLogsHome),
			"DataNodeStorageHome":                    furyPaths.StatePathFor(DataNodeStorageHome),
			"DataNodeStorageSQLStoreHome":            furyPaths.StatePathFor(DataNodeStorageSQLStoreHome),
			"DataNodeStorageSQLStoreNodeDataHome":    furyPaths.StatePathFor(DataNodeStorageSQLStoreNodeDataHome),
			"DataNodeEmbeddedPostgresRuntimeDir":     furyPaths.StatePathFor(DataNodeEmbeddedPostgresRuntimeDir),
			"DataNodeNetworkHistoryHome":             furyPaths.StatePathFor(DataNodeNetworkHistoryHome),
			"DataNodeNetworkHistorySnapshotCopyTo":   furyPaths.StatePathFor(DataNodeNetworkHistorySnapshotCopyTo),
			"DataNodeNetworkHistorySnapshotCopyFrom": furyPaths.StatePathFor(DataNodeNetworkHistorySnapshotCopyFrom),
			"DataNodeEventBufferHome":                furyPaths.StatePathFor(DataNodeEventBufferHome),
			"DataNodeArchivedEventBufferHome":        furyPaths.StatePathFor(DataNodeArchivedEventBufferHome),
			"NodeStateHome":                          furyPaths.StatePathFor(NodeStateHome),
			"NodeLogsHome":                           furyPaths.StatePathFor(NodeLogsHome),
			"CheckpointStateHome":                    furyPaths.StatePathFor(CheckpointStateHome),
			"SnapshotStateHome":                      furyPaths.StatePathFor(SnapshotStateHome),
			"SnapshotDBStateFile":                    furyPaths.StatePathFor(SnapshotDBStateFile),
			"WalletCLIStateHome":                     furyPaths.StatePathFor(WalletCLIStateHome),
			"WalletCLILogsHome":                      furyPaths.StatePathFor(WalletCLILogsHome),
			"WalletAppStateHome":                     furyPaths.StatePathFor(WalletAppStateHome),
			"WalletAppLogsHome":                      furyPaths.StatePathFor(WalletAppLogsHome),
			"WalletServiceStateHome":                 furyPaths.StatePathFor(WalletServiceStateHome),
			"WalletServiceLogsHome":                  furyPaths.StatePathFor(WalletServiceLogsHome),
		},
	}
}

func Explain(name string) (string, error) {
	paths := map[string]string{
		"DataNodeCacheHome":                      `This folder contains the cache used by the data-node.`,
		"DataNodeConfigHome":                     `This folder contains the configuration files used by the data-node.`,
		"DataNodeDefaultConfigFile":              `This file contains the configuration used by the data-node.`,
		"FaucetConfigHome":                       `This folder contains the configuration files used by the faucet.`,
		"FaucetDefaultConfigFile":                `This file contains the configuration used by the faucet.`,
		"NodeConfigHome":                         `This folder contains the configuration files used by the node.`,
		"NodeDefaultConfigFile":                  `This file contains the configuration used by the node.`,
		"NodeWalletsConfigFile":                  `This file contains information related to the registered node's wallets used by the node.`,
		"WalletCLIConfigHome":                    `This folder contains the configuration files used by the wallet-cli.`,
		"WalletCLIDefaultConfigFile":             `This file contains the configuration used by the wallet-cli.`,
		"WalletAppConfigHome":                    `This folder contains the configuration files used by the wallet-app.`,
		"WalletAppFairgroundConfigFile":          `This file contains the Fairground configuration used by the wallet-app.`,
		"WalletAppDefaultConfigFile":             `This file contains the default configuration used by the wallet-app.`,
		"WalletServiceConfigHome":                `This folder contains the configuration files used by the wallet's service.`,
		"WalletServiceDefaultConfigFile":         `This file contains the configuration used by the wallet service.`,
		"WalletServiceNetworksConfigHome":        `This folder contains the network configuration files used by the wallet's service.`,
		"NodeDataHome":                           `This folder contains the data managed by the node.`,
		"NodeWalletsDataHome":                    `This folder contains the data managed by the node's wallets.`,
		"FuryNodeWalletsDataHome":                `This folder contains the Fury wallet registered as node's wallet, used by the node to sign Fury commands.`,
		"EthereumNodeWalletsDataHome":            `This folder contains the Ethereum wallet registered as node's wallet, used by the node to interact with the Ethereum blockchain.`,
		"FaucetDataHome":                         `This folder contains the data used by the faucet.`,
		"FaucetWalletsDataHome":                  `This folder contains the Fury wallet used by the faucet to sign its deposit commands.`,
		"WalletsDataHome":                        `This folder contains the "user's" wallets. These wallets are used by the user to issue commands to a Fury network.`,
		"WalletServiceDataHome":                  `This folder contains the data used by the wallet's service.`,
		"WalletServiceRSAKeysDataHome":           `This folder contains the RSA keys used by the wallet's service for authentication.`,
		"WalletServicePublicRSAKeyDataFile":      `This file contains the public RSA key used by the wallet's service for authentication.`,
		"WalletServicePrivateRSAKeyDataFile":     `This file contains the private RSA key used by the wallet's service for authentication.`,
		"DataNodeStateHome":                      `This folder contains the state files used by the data-node.`,
		"DataNodeAutoCertHome":                   `This folder contains the autogenerated SSL certificates.`,
		"DataNodeLogsHome":                       `This folder contains the log files generated by the data-node.`,
		"DataNodeStorageHome":                    `This folder contains the consolidated state, built out of the Fury network events, and served by the data-node's API.`,
		"DataNodeStorageSQLStoreHome":            `This folder contains the files generated by the embedded Postgres database.`,
		"DataNodeStorageSQLStoreNodeDataHome":    `This folder contains the data files generated by the embedded Postgres database.`,
		"DataNodeEmbeddedPostgresRuntimeDir":     `This folder contains the files used at runtime by the embedded Postgres database.`,
		"DataNodeNetworkHistoryHome":             `This folder contains the network history data.`,
		"DataNodeNetworkHistorySnapshotCopyTo":   `This folder contains the snapshots created by the datanode.`,
		"DataNodeNetworkHistorySnapshotCopyFrom": `This folder contains the snapshots read by the datanode.`,
		"DataNodeEventBufferHome":                `This folder contains the event buffer files.`,
		"DataNodeArchivedEventBufferHome":        `This folder contains the archived event buffer files.`,
		"NodeStateHome":                          `This folder contains the state files used by the node.`,
		"NodeLogsHome":                           `This folder contains the log files generated by the node.`,
		"CheckpointStateHome":                    `This folder contains the network checkpoints generated by the node.`,
		"SnapshotStateHome":                      `This folder contains the Tendermint snapshots of the application state generated by the node.`,
		"SnapshotDBStateFile":                    `This file is a database containing the snapshots' data of the of the application state generated by the node`,
		"WalletCLIStateHome":                     `This folder contains the state files used by the wallet-cli.`,
		"WalletCLILogsHome":                      `This folder contains the log files generated by the wallet-cli.`,
		"WalletAppStateHome":                     `This folder contains the state files used by the wallet-app.`,
		"WalletAppLogsHome":                      `This folder contains the log files generated by the wallet-app.`,
		"WalletServiceStateHome":                 `This folder contains the state files used by the wallet's service.`,
		"WalletServiceLogsHome":                  `This folder contains the log files generated by the wallet's service'.`,
	}

	description, ok := paths[name]
	if !ok {
		return "", fmt.Errorf("path \"%s\" has no documentation", name)
	}

	return description, nil
}
