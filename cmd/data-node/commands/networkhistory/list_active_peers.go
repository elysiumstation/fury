package networkhistory

import (
	"context"
	"fmt"
	"os"

	coreConfig "github.com/elysiumstation/fury/core/config"
	"github.com/elysiumstation/fury/datanode/config"
	vgjson "github.com/elysiumstation/fury/libs/json"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"
	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
)

type listActivePeers struct {
	config.FuryHomeFlag
	config.Config
	coreConfig.OutputFlag
}

type listActivePeersOutput struct {
	ActivePeers []string
}

func (o *listActivePeersOutput) printHuman() {
	if len(o.ActivePeers) == 0 {
		fmt.Printf("No active peers found\n")
		return
	}
	fmt.Printf("Active Peers:\n\n")

	for _, peer := range o.ActivePeers {
		fmt.Printf("Active Peer:  %s\n", peer)
	}
}

func (cmd *listActivePeers) Execute(_ []string) error {
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.InfoLevel
	cfg.Environment = "custom"
	log := logging.NewLoggerFromConfig(
		cfg,
	)
	defer log.AtExit()

	furyPaths := paths.New(cmd.FuryHome)
	err := fixConfig(&cmd.Config, furyPaths)
	if err != nil {
		handleErr(log,
			cmd.Output.IsJSON(),
			"failed to fix config",
			err)
	}

	if !datanodeLive(cmd.Config) {
		handleErr(log,
			cmd.Output.IsJSON(),
			"datanode must be running for this command to work",
			fmt.Errorf("couldn't connect to datanode on %v:%v", cmd.Config.API.IP, cmd.Config.API.Port))
		os.Exit(1)
	}

	client, conn, err := getDatanodeClient(cmd.Config)
	if err != nil {
		handleErr(log,
			cmd.Output.IsJSON(),
			"failed to get datanode client",
			err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	resp, err := client.GetActiveNetworkHistoryPeerAddresses(ctx, &v2.GetActiveNetworkHistoryPeerAddressesRequest{})
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get active peer addresses", errorFromGrpcError("", err))
		os.Exit(1)
	}

	output := listActivePeersOutput{ActivePeers: resp.IpAddresses}

	if cmd.Output.IsJSON() {
		if err := vgjson.Print(&output); err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to marshal output", err)
			os.Exit(1)
		}
	} else {
		output.printHuman()
	}

	return nil
}
