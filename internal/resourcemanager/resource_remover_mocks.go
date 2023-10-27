// Code generated by MockGen. DO NOT EDIT.
// Source: resource_remover.go

// Package resourcemanager is a generated GoMock package.
package resourcemanager

import (
	context "context"
	resourcemodel "github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResourceRemover is a mock of ResourceRemover interface.
type MockResourceRemover struct {
	ctrl     *gomock.Controller
	recorder *MockResourceRemoverMockRecorder
}

// MockResourceRemoverMockRecorder is the mock recorder for MockResourceRemover.
type MockResourceRemoverMockRecorder struct {
	mock *MockResourceRemover
}

// NewMockResourceRemover creates a new mock instance.
func NewMockResourceRemover(ctrl *gomock.Controller) *MockResourceRemover {
	mock := &MockResourceRemover{ctrl: ctrl}
	mock.recorder = &MockResourceRemoverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceRemover) EXPECT() *MockResourceRemoverMockRecorder {
	return m.recorder
}

// RemoveReport mocks base method.
func (m *MockResourceRemover) RemoveReport(ctx context.Context, resource resourcemodel.Resource) ResourceReport {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveReport", ctx, resource)
	ret0, _ := ret[0].(ResourceReport)
	return ret0
}

// RemoveReport indicates an expected call of RemoveReport.
func (mr *MockResourceRemoverMockRecorder) RemoveReport(ctx, resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveReport", reflect.TypeOf((*MockResourceRemover)(nil).RemoveReport), ctx, resource)
}

// RemoveResource mocks base method.
func (m *MockResourceRemover) RemoveResource(ctx context.Context, resource resourcemodel.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveResource", ctx, resource)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveResource indicates an expected call of RemoveResource.
func (mr *MockResourceRemoverMockRecorder) RemoveResource(ctx, resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveResource", reflect.TypeOf((*MockResourceRemover)(nil).RemoveResource), ctx, resource)
}
