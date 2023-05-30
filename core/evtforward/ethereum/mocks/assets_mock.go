// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/elysiumstation/fury/core/evtforward/ethereum (interfaces: Assets)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAssets is a mock of Assets interface.
type MockAssets struct {
	ctrl     *gomock.Controller
	recorder *MockAssetsMockRecorder
}

// MockAssetsMockRecorder is the mock recorder for MockAssets.
type MockAssetsMockRecorder struct {
	mock *MockAssets
}

// NewMockAssets creates a new mock instance.
func NewMockAssets(ctrl *gomock.Controller) *MockAssets {
	mock := &MockAssets{ctrl: ctrl}
	mock.recorder = &MockAssetsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssets) EXPECT() *MockAssetsMockRecorder {
	return m.recorder
}

// GetFuryIDFromEthereumAddress mocks base method.
func (m *MockAssets) GetFuryIDFromEthereumAddress(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFuryIDFromEthereumAddress", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFuryIDFromEthereumAddress indicates an expected call of GetFuryIDFromEthereumAddress.
func (mr *MockAssetsMockRecorder) GetFuryIDFromEthereumAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFuryIDFromEthereumAddress", reflect.TypeOf((*MockAssets)(nil).GetFuryIDFromEthereumAddress), arg0)
}