package networkhistory

import (
	"context"
	"fmt"
	"os"

	coreConfig "github.com/elysiumstation/fury/core/config"
	vgjson "github.com/elysiumstation/fury/libs/json"
	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"

	"github.com/elysiumstation/fury/datanode/config"
	"github.com/elysiumstation/fury/datanode/networkhistory/store"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"
)

var errNoHistorySegmentFound = fmt.Errorf("no history segments found")

type latestHistorySegment struct {
	config.FuryHomeFlag
	coreConfig.OutputFlag
	config.Config
}

type latestHistoryOutput struct {
	LatestSegment *v2.HistorySegment
}

func (o *latestHistoryOutput) printHuman() {
	fmt.Printf("Latest segment to use data {%s}\n\n", o.LatestSegment)
}

func (cmd *latestHistorySegment) Execute(_ []string) error {
	ctx, cfunc := context.WithCancel(context.Background())
	defer cfunc()
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.ErrorLevel
	cfg.Environment = "custom"
	log := logging.NewLoggerFromConfig(
		cfg,
	)
	defer log.AtExit()

	furyPaths := paths.New(cmd.FuryHome)
	err := fixConfig(&cmd.Config, furyPaths)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to fix config", err)
		os.Exit(1)
	}

	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	var latestSegment *v2.HistorySegment
	if datanodeLive(cmd.Config) {
		client, conn, err := getDatanodeClient(cmd.Config)
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to get datanode client", err)
			os.Exit(1)
		}
		defer func() { _ = conn.Close() }()

		response, err := client.ListAllNetworkHistorySegments(ctx, &v2.ListAllNetworkHistorySegmentsRequest{})
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to list all network history segments", errorFromGrpcError("", err))
			os.Exit(1)
		}

		if len(response.Segments) < 1 {
			handleErr(log, cmd.Output.IsJSON(), errNoHistorySegmentFound.Error(), errNoHistorySegmentFound)
			os.Exit(1)
		}

		latestSegment = response.Segments[0]
	} else {
		networkHistoryStore, err := store.New(ctx, log, cmd.Config.ChainID, cmd.Config.NetworkHistory.Store,
			furyPaths.StatePathFor(paths.DataNodeNetworkHistoryHome), cmd.Config.MaxMemoryPercent)
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to create network history store", err)
			os.Exit(1)
		}
		defer networkHistoryStore.Stop()

		segments, err := networkHistoryStore.ListAllIndexEntriesOldestFirst()
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to list all network history segments", err)
			os.Exit(1)
		}

		if len(segments) < 1 {
			handleErr(log, cmd.Output.IsJSON(), errNoHistorySegmentFound.Error(), errNoHistorySegmentFound)
			os.Exit(1)
		}

		latestSegmentIndex := segments[len(segments)-1]

		latestSegment = &v2.HistorySegment{
			FromHeight:               latestSegmentIndex.GetFromHeight(),
			ToHeight:                 latestSegmentIndex.GetToHeight(),
			HistorySegmentId:         latestSegmentIndex.GetHistorySegmentId(),
			PreviousHistorySegmentId: latestSegmentIndex.GetPreviousHistorySegmentId(),
		}
	}

	output := latestHistoryOutput{
		LatestSegment: latestSegment,
	}

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
