package v1

import (
	"context"

	api "github.com/elysiumstation/fury/protos/fury/api/v1"
	commandspb "github.com/elysiumstation/fury/protos/fury/commands/v1"
	walletpb "github.com/elysiumstation/fury/protos/fury/wallet/v1"
	nodetypes "github.com/elysiumstation/fury/wallet/api/node/types"
	"github.com/elysiumstation/fury/wallet/wallet"
)

// Generates mocks
//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks github.com/elysiumstation/fury/wallet/service/v1 WalletHandler,Auth,NodeForward,RSAStore,SpamHandler

//nolint:interfacebloat
type WalletHandler interface {
	CreateWallet(name, passphrase string) (string, error)
	ImportWallet(name, passphrase, recoveryPhrase string, version uint32) error
	LoginWallet(name, passphrase string) error
	SecureGenerateKeyPair(name, passphrase string, meta []wallet.Metadata) (string, error)
	GetPublicKey(name, pubKey string) (wallet.PublicKey, error)
	ListPublicKeys(name string) ([]wallet.PublicKey, error)
	SignTx(name string, req *walletpb.SubmitTransactionRequest, height uint64, chainID string) (*commandspb.Transaction, error)
	SignAny(name string, inputData []byte, pubKey string) ([]byte, error)
	VerifyAny(inputData, sig []byte, pubKey string) (bool, error)
	TaintKey(name, pubKey, passphrase string) error
	UpdateMeta(name, pubKey, passphrase string, meta []wallet.Metadata) error
}

type Auth interface {
	NewSession(name string) (string, error)
	VerifyToken(token string) (string, error)
	Revoke(token string) (string, error)
	RevokeAllToken()
}

type NodeForward interface {
	SendTx(context.Context, *commandspb.Transaction, api.SubmitTransactionRequest_Type, int) (*api.SubmitTransactionResponse, error)
	CheckTx(context.Context, *commandspb.Transaction, int) (*api.CheckTransactionResponse, error)
	HealthCheck(context.Context) error
	LastBlockHeightAndHash(context.Context) (*api.LastBlockHeightResponse, int, error)
	SpamStatistics(context.Context, string) (*api.GetSpamStatisticsResponse, int, error)
	Stop()
}

type RSAStore interface {
	GetRsaKeys() (*RSAKeys, error)
}

type SpamHandler interface {
	GenerateProofOfWork(pubKey string, stats *nodetypes.SpamStatistics) (*commandspb.ProofOfWork, error)
	CheckSubmission(req *walletpb.SubmitTransactionRequest, stats *nodetypes.SpamStatistics) error
}
