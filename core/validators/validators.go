package validators

//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks github.com/elysiumstation/fury/core/validators NodeWallets,TimeService,Commander,ValidatorTopology,Wallet,ValidatorPerformance,Notary,Signatures,MultiSigTopology
