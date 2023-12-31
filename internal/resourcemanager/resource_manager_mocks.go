// Code generated by MockGen. DO NOT EDIT.
// Source: resource_manager.go

// Package resourcemanager is a generated GoMock package.
package resourcemanager

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResourceManager is a mock of ResourceManager interface.
type MockResourceManager struct {
	ctrl     *gomock.Controller
	recorder *MockResourceManagerMockRecorder
}

// MockResourceManagerMockRecorder is the mock recorder for MockResourceManager.
type MockResourceManagerMockRecorder struct {
	mock *MockResourceManager
}

// NewMockResourceManager creates a new mock instance.
func NewMockResourceManager(ctrl *gomock.Controller) *MockResourceManager {
	mock := &MockResourceManager{ctrl: ctrl}
	mock.recorder = &MockResourceManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceManager) EXPECT() *MockResourceManagerMockRecorder {
	return m.recorder
}

// AddCPU mocks base method.
func (m *MockResourceManager) AddCPU(ctx context.Context, cpu float32) (ResourceReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCPU", ctx, cpu)
	ret0, _ := ret[0].(ResourceReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCPU indicates an expected call of AddCPU.
func (mr *MockResourceManagerMockRecorder) AddCPU(ctx, cpu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCPU", reflect.TypeOf((*MockResourceManager)(nil).AddCPU), ctx, cpu)
}

// AddRAM mocks base method.
func (m *MockResourceManager) AddRAM(ctx context.Context, ram float32) (ResourceReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRAM", ctx, ram)
	ret0, _ := ret[0].(ResourceReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRAM indicates an expected call of AddRAM.
func (mr *MockResourceManagerMockRecorder) AddRAM(ctx, ram interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRAM", reflect.TypeOf((*MockResourceManager)(nil).AddRAM), ctx, ram)
}

// RemoveCPU mocks base method.
func (m *MockResourceManager) RemoveCPU(ctx context.Context, cpu, ramLimit float32) (ResourceReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCPU", ctx, cpu, ramLimit)
	ret0, _ := ret[0].(ResourceReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveCPU indicates an expected call of RemoveCPU.
func (mr *MockResourceManagerMockRecorder) RemoveCPU(ctx, cpu, ramLimit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCPU", reflect.TypeOf((*MockResourceManager)(nil).RemoveCPU), ctx, cpu, ramLimit)
}

// RemoveRAM mocks base method.
func (m *MockResourceManager) RemoveRAM(ctx context.Context, ram, cpuLimit float32) (ResourceReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRAM", ctx, ram, cpuLimit)
	ret0, _ := ret[0].(ResourceReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveRAM indicates an expected call of RemoveRAM.
func (mr *MockResourceManagerMockRecorder) RemoveRAM(ctx, ram, cpuLimit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRAM", reflect.TypeOf((*MockResourceManager)(nil).RemoveRAM), ctx, ram, cpuLimit)
}
