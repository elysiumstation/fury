// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package start

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elysiumstation/fury/libs/subscribers"

	"github.com/elysiumstation/fury/datanode/admin"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"golang.org/x/sync/errgroup"

	"github.com/elysiumstation/fury/datanode/api"
	"github.com/elysiumstation/fury/datanode/broker"
	"github.com/elysiumstation/fury/datanode/config"
	"github.com/elysiumstation/fury/datanode/gateway/server"
	"github.com/elysiumstation/fury/datanode/metrics"
	"github.com/elysiumstation/fury/datanode/networkhistory"
	"github.com/elysiumstation/fury/datanode/networkhistory/snapshot"
	"github.com/elysiumstation/fury/datanode/sqlstore"
	"github.com/elysiumstation/fury/libs/pprof"
	"github.com/elysiumstation/fury/logging"
	"github.com/elysiumstation/fury/paths"
)

// NodeCommand use to implement 'node' command.
type NodeCommand struct {
	SQLSubscribers
	ctx    context.Context
	cancel context.CancelFunc

	embeddedPostgres              *embeddedpostgres.EmbeddedPostgres
	transactionalConnectionSource *sqlstore.ConnectionSource

	networkHistoryService *networkhistory.Service
	snapshotService       *snapshot.Service

	furyCoreServiceClient api.CoreServiceClient

	broker    *broker.Broker
	sqlBroker broker.SQLStoreEventBroker

	eventService *subscribers.Service

	pproffhandlr  *pprof.Pprofhandler
	Log           *logging.Logger
	furyPaths     paths.Paths
	configWatcher *config.Watcher
	conf          config.Config

	Version     string
	VersionHash string
}

func (l *NodeCommand) Run(ctx context.Context, cfgwatchr *config.Watcher, furyPaths paths.Paths, args []string) error {
	l.configWatcher = cfgwatchr

	l.conf = cfgwatchr.Get()
	l.furyPaths = furyPaths
	if l.cancel != nil {
		l.cancel()
	}
	l.ctx, l.cancel = context.WithCancel(ctx)

	stages := []func([]string) error{
		l.persistentPre,
		l.preRun,
		l.runNode,
		l.postRun,
		l.persistentPost,
	}
	for _, fn := range stages {
		if err := fn(args); err != nil {
			return err
		}
	}

	return nil
}

// Stop is for graceful shutdown.
func (l *NodeCommand) Stop() {
	l.cancel()
}

// runNode is the entry of node command.
func (l *NodeCommand) runNode([]string) error {
	nodeLog := l.Log.Named("start.runNode")
	var eg *errgroup.Group
	eg, l.ctx = errgroup.WithContext(l.ctx)

	// gRPC server
	grpcServer := l.createGRPCServer(l.conf.API)

	// Admin server
	adminServer := admin.NewServer(l.Log, l.conf.Admin, l.furyPaths, admin.NewNetworkHistoryAdminService(l.networkHistoryService))

	// watch configs
	l.configWatcher.OnConfigUpdate(
		func(cfg config.Config) {
			grpcServer.ReloadConf(cfg.API)
			adminServer.ReloadConf(cfg.Admin)
		},
	)

	// start the grpc server
	eg.Go(func() error { return grpcServer.Start(l.ctx, nil) })

	// start the admin server
	eg.Go(func() error {
		if err := adminServer.Start(l.ctx); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// start gateway
	if l.conf.GatewayEnabled {
		gty := server.New(l.conf.Gateway, l.Log, l.furyPaths)
		eg.Go(func() error { return gty.Start(l.ctx) })
	}

	eg.Go(func() error {
		return l.broker.Receive(l.ctx)
	})

	eg.Go(func() error {
		defer func() {
			if l.conf.NetworkHistory.Enabled {
				l.networkHistoryService.Stop()
			}
		}()

		return l.sqlBroker.Receive(l.ctx)
	})

	// waitSig will wait for a sigterm or sigint interrupt.
	eg.Go(func() error {
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

		select {
		case sig := <-gracefulStop:
			nodeLog.Info("Caught signal", logging.String("name", fmt.Sprintf("%+v", sig)))
			l.cancel()
		case <-l.ctx.Done():
			return l.ctx.Err()
		}
		return nil
	})

	metrics.Start(l.conf.Metrics)

	nodeLog.Info("Fury data node startup complete")

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		nodeLog.Error("Fury data node stopped with error", logging.Error(err))
		return fmt.Errorf("fury data node stopped with error: %w", err)
	}

	return nil
}

func (l *NodeCommand) createGRPCServer(config api.Config) *api.GRPCServer {
	grpcServer := api.NewGRPCServer(
		l.Log,
		config,
		l.furyCoreServiceClient,
		l.eventService,
		l.orderService,
		l.networkLimitsService,
		l.marketDataService,
		l.tradeService,
		l.assetService,
		l.accountService,
		l.rewardService,
		l.marketsService,
		l.delegationService,
		l.epochService,
		l.depositService,
		l.withdrawalService,
		l.governanceService,
		l.riskFactorService,
		l.riskService,
		l.networkParameterService,
		l.blockService,
		l.checkpointService,
		l.partyService,
		l.candleService,
		l.oracleSpecService,
		l.oracleDataService,
		l.liquidityProvisionService,
		l.positionService,
		l.transferService,
		l.stakeLinkingService,
		l.notaryService,
		l.multiSigService,
		l.keyRotationsService,
		l.ethereumKeyRotationsService,
		l.nodeService,
		l.marketDepthService,
		l.ledgerService,
		l.protocolUpgradeService,
		l.networkHistoryService,
		l.coreSnapshotService,
	)
	return grpcServer
}
