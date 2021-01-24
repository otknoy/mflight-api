// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_mflight is a generated GoMock package.
package mock_mflight

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	mflight "mflight-api/infrastructure/mflight"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetSensorMonitor mocks base method
func (m *MockClient) GetSensorMonitor(ctx context.Context) (*mflight.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSensorMonitor", ctx)
	ret0, _ := ret[0].(*mflight.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSensorMonitor indicates an expected call of GetSensorMonitor
func (mr *MockClientMockRecorder) GetSensorMonitor(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSensorMonitor", reflect.TypeOf((*MockClient)(nil).GetSensorMonitor), ctx)
}
