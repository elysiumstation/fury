// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/elysiumstation/fury/wallet/wallet (interfaces: Store)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	wallet "github.com/elysiumstation/fury/wallet/wallet"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetWallet mocks base method.
func (m *MockStore) GetWallet(arg0 context.Context, arg1, arg2 string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWallet", arg0, arg1, arg2)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWallet indicates an expected call of GetWallet.
func (mr *MockStoreMockRecorder) GetWallet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWallet", reflect.TypeOf((*MockStore)(nil).GetWallet), arg0, arg1, arg2)
}

// GetWalletPath mocks base method.
func (m *MockStore) GetWalletPath(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletPath", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetWalletPath indicates an expected call of GetWalletPath.
func (mr *MockStoreMockRecorder) GetWalletPath(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletPath", reflect.TypeOf((*MockStore)(nil).GetWalletPath), arg0)
}

// ListWallets mocks base method.
func (m *MockStore) ListWallets(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWallets", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWallets indicates an expected call of ListWallets.
func (mr *MockStoreMockRecorder) ListWallets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWallets", reflect.TypeOf((*MockStore)(nil).ListWallets), arg0)
}

// SaveWallet mocks base method.
func (m *MockStore) SaveWallet(arg0 context.Context, arg1 wallet.Wallet, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWallet", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveWallet indicates an expected call of SaveWallet.
func (mr *MockStoreMockRecorder) SaveWallet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWallet", reflect.TypeOf((*MockStore)(nil).SaveWallet), arg0, arg1, arg2)
}

// WalletExists mocks base method.
func (m *MockStore) WalletExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalletExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WalletExists indicates an expected call of WalletExists.
func (mr *MockStoreMockRecorder) WalletExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalletExists", reflect.TypeOf((*MockStore)(nil).WalletExists), arg0, arg1)
}
