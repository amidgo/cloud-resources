// Code generated by MockGen. DO NOT EDIT.
// Source: executor.go

// Package executor is a generated GoMock package.
package executor

import (
	context "context"
	resourcemodel "github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	statiscticsmodel "github.com/amidgo/cloud-resources/internal/model/statiscticsmodel"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExecuteInitter is a mock of ExecuteInitter interface.
type MockExecuteInitter struct {
	ctrl     *gomock.Controller
	recorder *MockExecuteInitterMockRecorder
}

// MockExecuteInitterMockRecorder is the mock recorder for MockExecuteInitter.
type MockExecuteInitterMockRecorder struct {
	mock *MockExecuteInitter
}

// NewMockExecuteInitter creates a new mock instance.
func NewMockExecuteInitter(ctrl *gomock.Controller) *MockExecuteInitter {
	mock := &MockExecuteInitter{ctrl: ctrl}
	mock.recorder = &MockExecuteInitterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecuteInitter) EXPECT() *MockExecuteInitterMockRecorder {
	return m.recorder
}

// Init mocks base method.
func (m *MockExecuteInitter) Init(ctx context.Context) ([]*resourcemodel.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx)
	ret0, _ := ret[0].([]*resourcemodel.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Init indicates an expected call of Init.
func (mr *MockExecuteInitterMockRecorder) Init(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockExecuteInitter)(nil).Init), ctx)
}

// MockExecutor is a mock of Executor interface.
type MockExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExecutorMockRecorder
}

// MockExecutorMockRecorder is the mock recorder for MockExecutor.
type MockExecutorMockRecorder struct {
	mock *MockExecutor
}

// NewMockExecutor creates a new mock instance.
func NewMockExecutor(ctrl *gomock.Controller) *MockExecutor {
	mock := &MockExecutor{ctrl: ctrl}
	mock.recorder = &MockExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecutor) EXPECT() *MockExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockExecutor) Execute(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockExecutorMockRecorder) Execute(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockExecutor)(nil).Execute), ctx)
}

// MockStatiscticsStorage is a mock of StatiscticsStorage interface.
type MockStatiscticsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStatiscticsStorageMockRecorder
}

// MockStatiscticsStorageMockRecorder is the mock recorder for MockStatiscticsStorage.
type MockStatiscticsStorageMockRecorder struct {
	mock *MockStatiscticsStorage
}

// NewMockStatiscticsStorage creates a new mock instance.
func NewMockStatiscticsStorage(ctrl *gomock.Controller) *MockStatiscticsStorage {
	mock := &MockStatiscticsStorage{ctrl: ctrl}
	mock.recorder = &MockStatiscticsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatiscticsStorage) EXPECT() *MockStatiscticsStorageMockRecorder {
	return m.recorder
}

// Statisctics mocks base method.
func (m *MockStatiscticsStorage) Statisctics(ctx context.Context) (statiscticsmodel.Statistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Statisctics", ctx)
	ret0, _ := ret[0].(statiscticsmodel.Statistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Statisctics indicates an expected call of Statisctics.
func (mr *MockStatiscticsStorageMockRecorder) Statisctics(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Statisctics", reflect.TypeOf((*MockStatiscticsStorage)(nil).Statisctics), ctx)
}

// MockResourceStorage is a mock of ResourceStorage interface.
type MockResourceStorage struct {
	ctrl     *gomock.Controller
	recorder *MockResourceStorageMockRecorder
}

// MockResourceStorageMockRecorder is the mock recorder for MockResourceStorage.
type MockResourceStorageMockRecorder struct {
	mock *MockResourceStorage
}

// NewMockResourceStorage creates a new mock instance.
func NewMockResourceStorage(ctrl *gomock.Controller) *MockResourceStorage {
	mock := &MockResourceStorage{ctrl: ctrl}
	mock.recorder = &MockResourceStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceStorage) EXPECT() *MockResourceStorageMockRecorder {
	return m.recorder
}

// ResourceList mocks base method.
func (m *MockResourceStorage) ResourceList(ctx context.Context) ([]*resourcemodel.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResourceList", ctx)
	ret0, _ := ret[0].([]*resourcemodel.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResourceList indicates an expected call of ResourceList.
func (mr *MockResourceStorageMockRecorder) ResourceList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResourceList", reflect.TypeOf((*MockResourceStorage)(nil).ResourceList), ctx)
}
