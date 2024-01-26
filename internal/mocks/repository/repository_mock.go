// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository.go

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/retr0h/gilt/pkg/config"
)

// MockRepositoryManager is a mock of RepositoryManager interface.
type MockRepositoryManager struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryManagerMockRecorder
}

// MockRepositoryManagerMockRecorder is the mock recorder for MockRepositoryManager.
type MockRepositoryManagerMockRecorder struct {
	mock *MockRepositoryManager
}

// NewMockRepositoryManager creates a new mock instance.
func NewMockRepositoryManager(ctrl *gomock.Controller) *MockRepositoryManager {
	mock := &MockRepositoryManager{ctrl: ctrl}
	mock.recorder = &MockRepositoryManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryManager) EXPECT() *MockRepositoryManagerMockRecorder {
	return m.recorder
}

// Clone mocks base method.
func (m *MockRepositoryManager) Clone(config config.Repository, cloneDir string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone", config, cloneDir)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Clone indicates an expected call of Clone.
func (mr *MockRepositoryManagerMockRecorder) Clone(config, cloneDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockRepositoryManager)(nil).Clone), config, cloneDir)
}

// CopySources mocks base method.
func (m *MockRepositoryManager) CopySources(config config.Repository, cloneDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopySources", config, cloneDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopySources indicates an expected call of CopySources.
func (mr *MockRepositoryManagerMockRecorder) CopySources(config, cloneDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopySources", reflect.TypeOf((*MockRepositoryManager)(nil).CopySources), config, cloneDir)
}

// Worktree mocks base method.
func (m *MockRepositoryManager) Worktree(config config.Repository, cloneDir, targetDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Worktree", config, cloneDir, targetDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// Worktree indicates an expected call of Worktree.
func (mr *MockRepositoryManagerMockRecorder) Worktree(config, cloneDir, targetDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Worktree", reflect.TypeOf((*MockRepositoryManager)(nil).Worktree), config, cloneDir, targetDir)
}