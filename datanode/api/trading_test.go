// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package api_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/jackc/pgconn"

	"github.com/elysiumstation/fury/libs/subscribers"

	"github.com/elysiumstation/fury/datanode/api"
	"github.com/elysiumstation/fury/datanode/api/mocks"
	"github.com/elysiumstation/fury/datanode/broker"
	"github.com/elysiumstation/fury/datanode/candlesv2"
	"github.com/elysiumstation/fury/datanode/config"
	vgtesting "github.com/elysiumstation/fury/datanode/libs/testing"
	"github.com/elysiumstation/fury/datanode/service"
	"github.com/elysiumstation/fury/datanode/sqlstore"
	"github.com/elysiumstation/fury/logging"
	v2 "github.com/elysiumstation/fury/protos/data-node/api/v2"
	furyprotoapi "github.com/elysiumstation/fury/protos/fury/api/v1"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const connBufSize = 1024 * 1024

func waitForNode(ctx context.Context, t *testing.T, conn *grpc.ClientConn) {
	t.Helper()
	const maxSleep = 2000 // milliseconds

	c := v2.NewTradingDataServiceClient(conn)

	sleepTime := 10 // milliseconds
	for sleepTime < maxSleep {
		_, err := c.Ping(ctx, &v2.PingRequest{})
		if err == nil {
			return
		}

		fmt.Println(err)

		fmt.Printf("Sleeping for %d milliseconds\n", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		sleepTime *= 2
	}
	if sleepTime >= maxSleep {
		t.Fatalf("Gave up waiting for gRPC server to respond properly.")
	}
}

func getTestGRPCServer(t *testing.T, ctx context.Context) (tidy func(), conn *grpc.ClientConn, mockCoreServiceClient *mocks.MockCoreServiceClient, err error) {
	t.Helper()
	_, cleanupFn := vgtesting.NewFuryPaths()

	conf := config.NewDefaultConfig()
	conf.API.IP = "127.0.0.1"
	conf.API.Port = 64201

	// Mock BlockchainClient
	mockCtrl := gomock.NewController(t)

	mockCoreServiceClient = mocks.NewMockCoreServiceClient(mockCtrl)

	mockNetworkHistoryService := mocks.NewMockNetworkHistoryService(mockCtrl)

	rawEventSource, err := broker.NewEventReceiverSender(conf.Broker, logging.NewTestLogger(), "")
	if err != nil {
		t.Fatalf("failed to create raw event source: %v", err)
	}

	eventSource := broker.NewDeserializer(rawEventSource)
	if err != nil {
		t.Fatalf("failed to create event source: %v", err)
	}

	conf.CandlesV2.CandleStore.DefaultCandleIntervals = ""

	sqlConn := &sqlstore.ConnectionSource{
		Connection: dummyConnection{},
	}

	bro, err := broker.New(ctx, logging.NewTestLogger(), conf.Broker, "", eventSource)
	if err != nil {
		err = errors.Wrap(err, "failed to create broker")
		return
	}

	logger := logging.NewTestLogger()
	eventService := subscribers.NewService(logger, bro, conf.Broker.EventBusClientBufferSize)
	sqlOrderStore := sqlstore.NewOrders(sqlConn)
	sqlOrderService := service.NewOrder(sqlOrderStore, logger)
	sqlNetworkLimitsService := service.NewNetworkLimits(sqlstore.NewNetworkLimits(sqlConn))
	sqlMarketDataService := service.NewMarketData(sqlstore.NewMarketData(sqlConn), logger)
	sqlCandleStore := sqlstore.NewCandles(ctx, sqlConn, conf.CandlesV2.CandleStore)
	sqlCandlesService := candlesv2.NewService(ctx, logger, conf.CandlesV2, sqlCandleStore)
	sqlTradeService := service.NewTrade(sqlstore.NewTrades(sqlConn), logger)
	sqlPositionService := service.NewPosition(sqlstore.NewPositions(sqlConn), logger)
	sqlAssetService := service.NewAsset(sqlstore.NewAssets(sqlConn))
	sqlAccountService := service.NewAccount(sqlstore.NewAccounts(sqlConn), sqlstore.NewBalances(sqlConn), logger)
	sqlRewardsService := service.NewReward(sqlstore.NewRewards(sqlConn), logger)
	sqlMarketsService := service.NewMarkets(sqlstore.NewMarkets(sqlConn))
	sqlDelegationService := service.NewDelegation(sqlstore.NewDelegations(sqlConn), logger)
	sqlEpochService := service.NewEpoch(sqlstore.NewEpochs(sqlConn))
	sqlDepositService := service.NewDeposit(sqlstore.NewDeposits(sqlConn))
	sqlWithdrawalService := service.NewWithdrawal(sqlstore.NewWithdrawals(sqlConn))
	sqlGovernanceService := service.NewGovernance(sqlstore.NewProposals(sqlConn), sqlstore.NewVotes(sqlConn), logger)
	sqlRiskFactorsService := service.NewRiskFactor(sqlstore.NewRiskFactors(sqlConn))
	sqlMarginLevelsService := service.NewRisk(sqlstore.NewMarginLevels(sqlConn), sqlAccountService, logger)
	sqlNetParamService := service.NewNetworkParameter(sqlstore.NewNetworkParameters(sqlConn))
	sqlBlockService := service.NewBlock(sqlstore.NewBlocks(sqlConn))
	sqlCheckpointService := service.NewCheckpoint(sqlstore.NewCheckpoints(sqlConn))
	sqlPartyService := service.NewParty(sqlstore.NewParties(sqlConn))
	sqlOracleSpecService := service.NewOracleSpec(sqlstore.NewOracleSpec(sqlConn))
	sqlOracleDataService := service.NewOracleData(sqlstore.NewOracleData(sqlConn))
	sqlLPDataService := service.NewLiquidityProvision(sqlstore.NewLiquidityProvision(sqlConn, logger))
	sqlTransferService := service.NewTransfer(sqlstore.NewTransfers(sqlConn))
	sqlStakeLinkingService := service.NewStakeLinking(sqlstore.NewStakeLinking(sqlConn))
	sqlNotaryService := service.NewNotary(sqlstore.NewNotary(sqlConn))
	sqlMultiSigService := service.NewMultiSig(sqlstore.NewERC20MultiSigSignerEvent(sqlConn))
	sqlKeyRotationsService := service.NewKeyRotations(sqlstore.NewKeyRotations(sqlConn))
	sqlEthereumKeyRotationService := service.NewEthereumKeyRotation(sqlstore.NewEthereumKeyRotations(sqlConn), logger)
	sqlNodeService := service.NewNode(sqlstore.NewNode(sqlConn))
	sqlMarketDepthService := service.NewMarketDepth(sqlOrderService, logger)
	sqlLedgerService := service.NewLedger(sqlstore.NewLedger(sqlConn), logger)
	sqlProtocolUpgradeService := service.NewProtocolUpgrade(sqlstore.NewProtocolUpgradeProposals(sqlConn), logger)
	sqlCoreSnapshotService := service.NewSnapshotData(sqlstore.NewCoreSnapshotData(sqlConn))

	g := api.NewGRPCServer(
		logger,
		conf.API,
		mockCoreServiceClient,
		eventService,
		sqlOrderService,
		sqlNetworkLimitsService,
		sqlMarketDataService,
		sqlTradeService,
		sqlAssetService,
		sqlAccountService,
		sqlRewardsService,
		sqlMarketsService,
		sqlDelegationService,
		sqlEpochService,
		sqlDepositService,
		sqlWithdrawalService,
		sqlGovernanceService,
		sqlRiskFactorsService,
		sqlMarginLevelsService,
		sqlNetParamService,
		sqlBlockService,
		sqlCheckpointService,
		sqlPartyService,
		sqlCandlesService,
		sqlOracleSpecService,
		sqlOracleDataService,
		sqlLPDataService,
		sqlPositionService,
		sqlTransferService,
		sqlStakeLinkingService,
		sqlNotaryService,
		sqlMultiSigService,
		sqlKeyRotationsService,
		sqlEthereumKeyRotationService,
		sqlNodeService,
		sqlMarketDepthService,
		sqlLedgerService,
		sqlProtocolUpgradeService,
		mockNetworkHistoryService,
		sqlCoreSnapshotService,
	)
	if g == nil {
		err = fmt.Errorf("failed to create gRPC server")
		return
	}

	tidy = func() {
		mockCtrl.Finish()
		cleanupFn()
	}

	lis := bufconn.Listen(connBufSize)
	ctxDialer := func(context.Context, string) (net.Conn, error) { return lis.Dial() }

	// Start the gRPC server, then wait for it to be ready.
	go func() {
		_ = g.Start(ctx, lis)
	}()

	conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(ctxDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to create connection to gRPC server")
	}

	waitForNode(ctx, t, conn)

	return tidy, conn, mockCoreServiceClient, err
}

type dummyConnection struct {
	sqlstore.Connection
}

func (d dummyConnection) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, pgx.ErrNoRows
}

func (d dummyConnection) CopyTo(context.Context, io.Writer, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func TestSubmitTransaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("proxy call is successful", func(t *testing.T) {
		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("Failed to get test gRPC server: %s", err.Error())
		}
		defer tidy()

		req := &furyprotoapi.SubmitTransactionRequest{
			Type: furyprotoapi.SubmitTransactionRequest_TYPE_UNSPECIFIED,
			Tx: &commandspb.Transaction{
				InputData: []byte("input data"),
				Signature: &commandspb.Signature{
					Value:   "value",
					Algo:    "algo",
					Version: 1,
				},
			},
		}

		expectedRes := &furyprotoapi.SubmitTransactionResponse{Success: true}

		furyReq := &furyprotoapi.SubmitTransactionRequest{
			Type: furyprotoapi.SubmitTransactionRequest_TYPE_UNSPECIFIED,
			Tx: &commandspb.Transaction{
				InputData: []byte("input data"),
				Signature: &commandspb.Signature{
					Value:   "value",
					Algo:    "algo",
					Version: 1,
				},
			},
		}

		mockTradingServiceClient.EXPECT().
			SubmitTransaction(gomock.Any(), vgtesting.ProtosEq(furyReq)).
			Return(&furyprotoapi.SubmitTransactionResponse{Success: true}, nil).Times(1)

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		actualResp, err := proxyClient.SubmitTransaction(ctx, req)
		assert.NoError(t, err)
		vgtesting.AssertProtoEqual(t, expectedRes, actualResp)
	})

	t.Run("proxy propagates an error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("Failed to get test gRPC server: %s", err.Error())
		}
		defer tidy()

		req := &furyprotoapi.SubmitTransactionRequest{
			Type: furyprotoapi.SubmitTransactionRequest_TYPE_COMMIT,
			Tx: &commandspb.Transaction{
				InputData: []byte("input data"),
				Signature: &commandspb.Signature{
					Value:   "value",
					Algo:    "algo",
					Version: 1,
				},
			},
		}

		furyReq := &furyprotoapi.SubmitTransactionRequest{
			Type: furyprotoapi.SubmitTransactionRequest_TYPE_COMMIT,
			Tx: &commandspb.Transaction{
				InputData: []byte("input data"),
				Signature: &commandspb.Signature{
					Value:   "value",
					Algo:    "algo",
					Version: 1,
				},
			},
		}

		mockTradingServiceClient.EXPECT().
			SubmitTransaction(gomock.Any(), vgtesting.ProtosEq(furyReq)).
			Return(nil, errors.New("Critical error"))

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		actualResp, err := proxyClient.SubmitTransaction(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, actualResp)
		assert.Contains(t, err.Error(), "Critical error")
	})
}

func TestSubmitRawTransaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("proxy call is successful", func(t *testing.T) {
		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("Failed to get test gRPC server: %s", err.Error())
		}
		defer tidy()

		tx := &commandspb.Transaction{
			InputData: []byte("input data"),
			Signature: &commandspb.Signature{
				Value:   "value",
				Algo:    "algo",
				Version: 1,
			},
		}

		bs, err := proto.Marshal(tx)
		assert.NoError(t, err)

		req := &furyprotoapi.SubmitRawTransactionRequest{
			Type: furyprotoapi.SubmitRawTransactionRequest_TYPE_UNSPECIFIED,
			Tx:   bs,
		}

		expectedRes := &furyprotoapi.SubmitRawTransactionResponse{Success: true}

		furyReq := &furyprotoapi.SubmitRawTransactionRequest{
			Type: furyprotoapi.SubmitRawTransactionRequest_TYPE_UNSPECIFIED,
			Tx:   bs,
		}

		mockTradingServiceClient.EXPECT().
			SubmitRawTransaction(gomock.Any(), vgtesting.ProtosEq(furyReq)).
			Return(&furyprotoapi.SubmitRawTransactionResponse{Success: true}, nil).Times(1)

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		actualResp, err := proxyClient.SubmitRawTransaction(ctx, req)
		assert.NoError(t, err)
		vgtesting.AssertProtoEqual(t, expectedRes, actualResp)
	})

	t.Run("proxy propagates an error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("Failed to get test gRPC server: %s", err.Error())
		}
		defer tidy()
		tx := &commandspb.Transaction{
			InputData: []byte("input data"),
			Signature: &commandspb.Signature{
				Value:   "value",
				Algo:    "algo",
				Version: 1,
			},
		}

		bs, err := proto.Marshal(tx)
		assert.NoError(t, err)

		req := &furyprotoapi.SubmitRawTransactionRequest{
			Type: furyprotoapi.SubmitRawTransactionRequest_TYPE_COMMIT,
			Tx:   bs,
		}

		furyReq := &furyprotoapi.SubmitRawTransactionRequest{
			Type: furyprotoapi.SubmitRawTransactionRequest_TYPE_COMMIT,
			Tx:   bs,
		}

		mockTradingServiceClient.EXPECT().
			SubmitRawTransaction(gomock.Any(), vgtesting.ProtosEq(furyReq)).
			Return(nil, errors.New("Critical error"))

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		actualResp, err := proxyClient.SubmitRawTransaction(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, actualResp)
		assert.Contains(t, err.Error(), "Critical error")
	})
}

func TestLastBlockHeight(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("proxy call is successful", func(t *testing.T) {
		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("Failed to get test gRPC server: %s", err.Error())
		}
		defer tidy()

		req := &furyprotoapi.LastBlockHeightRequest{}
		expectedRes := &furyprotoapi.LastBlockHeightResponse{Height: 20}

		furyReq := &furyprotoapi.LastBlockHeightRequest{}

		mockTradingServiceClient.EXPECT().
			LastBlockHeight(gomock.Any(), vgtesting.ProtosEq(furyReq)).
			Return(&furyprotoapi.LastBlockHeightResponse{Height: 20}, nil).Times(1)

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		actualResp, err := proxyClient.LastBlockHeight(ctx, req)
		assert.NoError(t, err)
		vgtesting.AssertProtoEqual(t, expectedRes, actualResp)
	})

	t.Run("proxy propagates an error", func(t *testing.T) {
		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("Failed to get test gRPC server: %s", err.Error())
		}
		defer tidy()

		req := &furyprotoapi.LastBlockHeightRequest{}
		furyReq := &furyprotoapi.LastBlockHeightRequest{}

		mockTradingServiceClient.EXPECT().
			LastBlockHeight(gomock.Any(), vgtesting.ProtosEq(furyReq)).
			Return(nil, fmt.Errorf("Critical error")).Times(1)

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		actualResp, err := proxyClient.LastBlockHeight(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, actualResp)
		assert.Contains(t, err.Error(), "Critical error")
	})
}

func TestGetSpamStatistics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("proxy call is successful", func(t *testing.T) {
		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("failed to get test gRPC server: %v", err)
		}
		defer tidy()

		req := &furyprotoapi.GetSpamStatisticsRequest{
			PartyId: "DEADBEEF",
		}

		wantReq := &furyprotoapi.GetSpamStatisticsRequest{
			PartyId: "DEADBEEF",
		}

		wantResp := &furyprotoapi.GetSpamStatisticsResponse{
			Statistics: &furyprotoapi.SpamStatistics{
				Proposals:         nil,
				Delegations:       nil,
				Transfers:         nil,
				NodeAnnouncements: nil,
				Votes:             nil,
			},
		}

		mockTradingServiceClient.EXPECT().
			GetSpamStatistics(gomock.Any(), vgtesting.ProtosEq(wantReq)).
			Return(&furyprotoapi.GetSpamStatisticsResponse{
				Statistics: &furyprotoapi.SpamStatistics{
					Proposals:         nil,
					Delegations:       nil,
					Transfers:         nil,
					NodeAnnouncements: nil,
					Votes:             nil,
				},
			}, nil).Times(1)

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		resp, err := proxyClient.GetSpamStatistics(ctx, req)
		assert.NoError(t, err)
		vgtesting.AssertProtoEqual(t, wantResp, resp)
	})

	t.Run("proxy propagates an error", func(t *testing.T) {
		tidy, conn, mockTradingServiceClient, err := getTestGRPCServer(t, ctx)
		if err != nil {
			t.Fatalf("failed to get test gRPC server: %v", err)
		}
		defer tidy()

		req := &furyprotoapi.GetSpamStatisticsRequest{
			PartyId: "DEADBEEF",
		}

		wantReq := &furyprotoapi.GetSpamStatisticsRequest{
			PartyId: "DEADBEEF",
		}

		mockTradingServiceClient.EXPECT().
			GetSpamStatistics(gomock.Any(), vgtesting.ProtosEq(wantReq)).
			Return(nil, fmt.Errorf("Critical error")).Times(1)

		proxyClient := furyprotoapi.NewCoreServiceClient(conn)
		assert.NotNil(t, proxyClient)

		resp, err := proxyClient.GetSpamStatistics(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Critical error")
	})
}
