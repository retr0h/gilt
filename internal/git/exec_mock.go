// Code generated by MockGen. DO NOT EDIT.
// Source: internal/git/types.go

// Package git is a generated GoMock package.
package git

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExecManager is a mock of ExecManager interface.
type MockExecManager struct {
	ctrl     *gomock.Controller
	recorder *MockExecManagerMockRecorder
}

// MockExecManagerMockRecorder is the mock recorder for MockExecManager.
type MockExecManagerMockRecorder struct {
	mock *MockExecManager
}

// NewMockExecManager creates a new mock instance.
func NewMockExecManager(ctrl *gomock.Controller) *MockExecManager {
	mock := &MockExecManager{ctrl: ctrl}
	mock.recorder = &MockExecManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExecManager) EXPECT() *MockExecManagerMockRecorder {
	return m.recorder
}

// RunCmd mocks base method.
func (m *MockExecManager) RunCmd(name string, args ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunCmd", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunCmd indicates an expected call of RunCmd.
func (mr *MockExecManagerMockRecorder) RunCmd(name interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCmd", reflect.TypeOf((*MockExecManager)(nil).RunCmd), varargs...)
}