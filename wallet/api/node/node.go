package node

import (
	"context"

	apipb "github.com/elysiumstation/fury/protos/fury/api/v1"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
	nodetypes "github.com/elysiumstation/fury/wallet/api/node/types"
)

// Generates mocks
//go:generate go run github.com/golang/mock/mockgen -destination mocks/nodes_mocks.go -package mocks github.com/elysiumstation/fury/wallet/api/node Node,Selector

// Node is the component used to get network information and send transactions.
type Node interface {
	Host() string
	Stop() error
	CheckTransaction(context.Context, *commandspb.Transaction) error
	SendTransaction(context.Context, *commandspb.Transaction, apipb.SubmitTransactionRequest_Type) (string, error)
	Statistics(ctx context.Context) (nodetypes.Statistics, error)
	LastBlock(context.Context) (nodetypes.LastBlock, error)
	SpamStatistics(ctx context.Context, pubKey string) (nodetypes.SpamStatistics, error)
}

// ReportType defines the type of event that occurred.
type ReportType string

var (
	InfoEvent    ReportType = "Info"
	WarningEvent ReportType = "Warning"
	ErrorEvent   ReportType = "Error"
	SuccessEvent ReportType = "Success"
)

type SelectionReporter func(ReportType, string)

// Selector implementing the strategy for node selection.
type Selector interface {
	Node(ctx context.Context, reporterFn SelectionReporter) (Node, error)
	Stop()
}
