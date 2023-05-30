package service

//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks github.com/elysiumstation/fury/datanode/service OrderStore,ChainStore,MarketStore,MarketDataStore,PositionStore
