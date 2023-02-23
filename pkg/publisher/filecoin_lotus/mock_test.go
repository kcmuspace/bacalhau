// Code generated by MockGen. DO NOT EDIT.
// Source: ./api/api.go

package filecoinlotus

import (
	context "context"
	reflect "reflect"

	api "github.com/filecoin-project/bacalhau/pkg/publisher/filecoin_lotus/api"
	address "github.com/filecoin-project/go-address"
	gomock "github.com/golang/mock/gomock"
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ClientDealPieceCID mocks base method.
func (m *MockClient) ClientDealPieceCID(arg0 context.Context, arg1 cid.Cid) (api.DataCIDSize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientDealPieceCID", arg0, arg1)
	ret0, _ := ret[0].(api.DataCIDSize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientDealPieceCID indicates an expected call of ClientDealPieceCID.
func (mr *MockClientMockRecorder) ClientDealPieceCID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientDealPieceCID", reflect.TypeOf((*MockClient)(nil).ClientDealPieceCID), arg0, arg1)
}

// ClientExport mocks base method.
func (m *MockClient) ClientExport(arg0 context.Context, arg1 api.ExportRef, arg2 api.FileRef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientExport", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClientExport indicates an expected call of ClientExport.
func (mr *MockClientMockRecorder) ClientExport(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientExport", reflect.TypeOf((*MockClient)(nil).ClientExport), arg0, arg1, arg2)
}

// ClientGetDealUpdates mocks base method.
func (m *MockClient) ClientGetDealUpdates(ctx context.Context) (<-chan api.DealInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientGetDealUpdates", ctx)
	ret0, _ := ret[0].(<-chan api.DealInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientGetDealUpdates indicates an expected call of ClientGetDealUpdates.
func (mr *MockClientMockRecorder) ClientGetDealUpdates(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientGetDealUpdates", reflect.TypeOf((*MockClient)(nil).ClientGetDealUpdates), ctx)
}

// ClientImport mocks base method.
func (m *MockClient) ClientImport(arg0 context.Context, arg1 api.FileRef) (*api.ImportRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientImport", arg0, arg1)
	ret0, _ := ret[0].(*api.ImportRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientImport indicates an expected call of ClientImport.
func (mr *MockClientMockRecorder) ClientImport(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientImport", reflect.TypeOf((*MockClient)(nil).ClientImport), arg0, arg1)
}

// ClientListImports mocks base method.
func (m *MockClient) ClientListImports(arg0 context.Context) ([]api.Import, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientListImports", arg0)
	ret0, _ := ret[0].([]api.Import)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientListImports indicates an expected call of ClientListImports.
func (mr *MockClientMockRecorder) ClientListImports(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientListImports", reflect.TypeOf((*MockClient)(nil).ClientListImports), arg0)
}

// ClientQueryAsk mocks base method.
func (m *MockClient) ClientQueryAsk(arg0 context.Context, arg1 peer.ID, arg2 address.Address) (*api.StorageAsk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientQueryAsk", arg0, arg1, arg2)
	ret0, _ := ret[0].(*api.StorageAsk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientQueryAsk indicates an expected call of ClientQueryAsk.
func (mr *MockClientMockRecorder) ClientQueryAsk(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientQueryAsk", reflect.TypeOf((*MockClient)(nil).ClientQueryAsk), arg0, arg1, arg2)
}

// ClientStartDeal mocks base method.
func (m *MockClient) ClientStartDeal(arg0 context.Context, arg1 *api.StartDealParams) (*cid.Cid, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientStartDeal", arg0, arg1)
	ret0, _ := ret[0].(*cid.Cid)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientStartDeal indicates an expected call of ClientStartDeal.
func (mr *MockClientMockRecorder) ClientStartDeal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientStartDeal", reflect.TypeOf((*MockClient)(nil).ClientStartDeal), arg0, arg1)
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// StateGetNetworkParams mocks base method.
func (m *MockClient) StateGetNetworkParams(arg0 context.Context) (*api.NetworkParams, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateGetNetworkParams", arg0)
	ret0, _ := ret[0].(*api.NetworkParams)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StateGetNetworkParams indicates an expected call of StateGetNetworkParams.
func (mr *MockClientMockRecorder) StateGetNetworkParams(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateGetNetworkParams", reflect.TypeOf((*MockClient)(nil).StateGetNetworkParams), arg0)
}

// StateListMiners mocks base method.
func (m *MockClient) StateListMiners(arg0 context.Context, arg1 api.TipSetKey) ([]address.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateListMiners", arg0, arg1)
	ret0, _ := ret[0].([]address.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StateListMiners indicates an expected call of StateListMiners.
func (mr *MockClientMockRecorder) StateListMiners(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateListMiners", reflect.TypeOf((*MockClient)(nil).StateListMiners), arg0, arg1)
}

// StateMinerInfo mocks base method.
func (m *MockClient) StateMinerInfo(arg0 context.Context, arg1 address.Address, arg2 api.TipSetKey) (api.MinerInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateMinerInfo", arg0, arg1, arg2)
	ret0, _ := ret[0].(api.MinerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StateMinerInfo indicates an expected call of StateMinerInfo.
func (mr *MockClientMockRecorder) StateMinerInfo(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateMinerInfo", reflect.TypeOf((*MockClient)(nil).StateMinerInfo), arg0, arg1, arg2)
}

// StateMinerPower mocks base method.
func (m *MockClient) StateMinerPower(arg0 context.Context, arg1 address.Address, arg2 api.TipSetKey) (*api.MinerPower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateMinerPower", arg0, arg1, arg2)
	ret0, _ := ret[0].(*api.MinerPower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StateMinerPower indicates an expected call of StateMinerPower.
func (mr *MockClientMockRecorder) StateMinerPower(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateMinerPower", reflect.TypeOf((*MockClient)(nil).StateMinerPower), arg0, arg1, arg2)
}

// Version mocks base method.
func (m *MockClient) Version(arg0 context.Context) (api.APIVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version", arg0)
	ret0, _ := ret[0].(api.APIVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Version indicates an expected call of Version.
func (mr *MockClientMockRecorder) Version(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockClient)(nil).Version), arg0)
}

// WalletDefaultAddress mocks base method.
func (m *MockClient) WalletDefaultAddress(arg0 context.Context) (address.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalletDefaultAddress", arg0)
	ret0, _ := ret[0].(address.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WalletDefaultAddress indicates an expected call of WalletDefaultAddress.
func (mr *MockClientMockRecorder) WalletDefaultAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalletDefaultAddress", reflect.TypeOf((*MockClient)(nil).WalletDefaultAddress), arg0)
}
