// Code generated by MockGen. DO NOT EDIT.
// Source: maden/pkg/etcd (interfaces: PodRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	shared "maden/pkg/shared"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPodRepository is a mock of PodRepository interface.
type MockPodRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPodRepositoryMockRecorder
}

// MockPodRepositoryMockRecorder is the mock recorder for MockPodRepository.
type MockPodRepositoryMockRecorder struct {
	mock *MockPodRepository
}

// NewMockPodRepository creates a new mock instance.
func NewMockPodRepository(ctrl *gomock.Controller) *MockPodRepository {
	mock := &MockPodRepository{ctrl: ctrl}
	mock.recorder = &MockPodRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPodRepository) EXPECT() *MockPodRepositoryMockRecorder {
	return m.recorder
}

// CreatePod mocks base method.
func (m *MockPodRepository) CreatePod(arg0 *shared.Pod) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePod", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePod indicates an expected call of CreatePod.
func (mr *MockPodRepositoryMockRecorder) CreatePod(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePod", reflect.TypeOf((*MockPodRepository)(nil).CreatePod), arg0)
}

// DeletePod mocks base method.
func (m *MockPodRepository) DeletePod(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePod", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePod indicates an expected call of DeletePod.
func (mr *MockPodRepositoryMockRecorder) DeletePod(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePod", reflect.TypeOf((*MockPodRepository)(nil).DeletePod), arg0)
}

// GetPodsByDeploymentID mocks base method.
func (m *MockPodRepository) GetPodsByDeploymentID(arg0 string) ([]shared.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPodsByDeploymentID", arg0)
	ret0, _ := ret[0].([]shared.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPodsByDeploymentID indicates an expected call of GetPodsByDeploymentID.
func (mr *MockPodRepositoryMockRecorder) GetPodsByDeploymentID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPodsByDeploymentID", reflect.TypeOf((*MockPodRepository)(nil).GetPodsByDeploymentID), arg0)
}

// ListPods mocks base method.
func (m *MockPodRepository) ListPods() ([]shared.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPods")
	ret0, _ := ret[0].([]shared.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPods indicates an expected call of ListPods.
func (mr *MockPodRepositoryMockRecorder) ListPods() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPods", reflect.TypeOf((*MockPodRepository)(nil).ListPods))
}

// UpdatePod mocks base method.
func (m *MockPodRepository) UpdatePod(arg0 *shared.Pod) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePod", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePod indicates an expected call of UpdatePod.
func (mr *MockPodRepositoryMockRecorder) UpdatePod(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePod", reflect.TypeOf((*MockPodRepository)(nil).UpdatePod), arg0)
}
