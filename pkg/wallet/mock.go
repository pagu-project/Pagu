// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/wallet/interface.go
//
// Generated by this command:
//
//	mockgen -source=./pkg/wallet/interface.go -destination=./pkg/wallet/mock.go -package=wallet
//

// Package wallet is a generated GoMock package.
package wallet

import (
	reflect "reflect"

	amount "github.com/pagu-project/Pagu/pkg/amount"
	gomock "go.uber.org/mock/gomock"
)

// MockIWallet is a mock of IWallet interface.
type MockIWallet struct {
	ctrl     *gomock.Controller
	recorder *MockIWalletMockRecorder
}

// MockIWalletMockRecorder is the mock recorder for MockIWallet.
type MockIWalletMockRecorder struct {
	mock *MockIWallet
}

// NewMockIWallet creates a new mock instance.
func NewMockIWallet(ctrl *gomock.Controller) *MockIWallet {
	mock := &MockIWallet{ctrl: ctrl}
	mock.recorder = &MockIWalletMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIWallet) EXPECT() *MockIWalletMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockIWallet) Address() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(string)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockIWalletMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockIWallet)(nil).Address))
}

// Balance mocks base method.
func (m *MockIWallet) Balance() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Balance")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Balance indicates an expected call of Balance.
func (mr *MockIWalletMockRecorder) Balance() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Balance", reflect.TypeOf((*MockIWallet)(nil).Balance))
}

// BondTransaction mocks base method.
func (m *MockIWallet) BondTransaction(pubKey, toAddress, memo string, amt amount.Amount) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BondTransaction", pubKey, toAddress, memo, amt)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BondTransaction indicates an expected call of BondTransaction.
func (mr *MockIWalletMockRecorder) BondTransaction(pubKey, toAddress, memo, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BondTransaction", reflect.TypeOf((*MockIWallet)(nil).BondTransaction), pubKey, toAddress, memo, amt)
}

// NewAddress mocks base method.
func (m *MockIWallet) NewAddress(lb string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAddress", lb)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewAddress indicates an expected call of NewAddress.
func (mr *MockIWalletMockRecorder) NewAddress(lb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAddress", reflect.TypeOf((*MockIWallet)(nil).NewAddress), lb)
}

// TransferTransaction mocks base method.
func (m *MockIWallet) TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferTransaction", toAddress, memo, amt)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferTransaction indicates an expected call of TransferTransaction.
func (mr *MockIWalletMockRecorder) TransferTransaction(toAddress, memo, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferTransaction", reflect.TypeOf((*MockIWallet)(nil).TransferTransaction), toAddress, memo, amt)
}
