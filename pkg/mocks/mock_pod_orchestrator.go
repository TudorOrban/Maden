// Code generated by MockGen. DO NOT EDIT.
// Source: maden/pkg/orchestrator (interfaces: PodOrchestrator)

// Package mocks is a generated GoMock package.
package mocks

import (
	shared "maden/pkg/shared"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPodOrchestrator is a mock of PodOrchestrator interface.
type MockPodOrchestrator struct {
	ctrl     *gomock.Controller
	recorder *MockPodOrchestratorMockRecorder
}

// MockPodOrchestratorMockRecorder is the mock recorder for MockPodOrchestrator.
type MockPodOrchestratorMockRecorder struct {
	mock *MockPodOrchestrator
}

// NewMockPodOrchestrator creates a new mock instance.
func NewMockPodOrchestrator(ctrl *gomock.Controller) *MockPodOrchestrator {
	mock := &MockPodOrchestrator{ctrl: ctrl}
	mock.recorder = &MockPodOrchestratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPodOrchestrator) EXPECT() *MockPodOrchestratorMockRecorder {
	return m.recorder
}

// OrchestratePodCreation mocks base method.
func (m *MockPodOrchestrator) OrchestratePodCreation(arg0 *shared.Pod) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrchestratePodCreation", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// OrchestratePodCreation indicates an expected call of OrchestratePodCreation.
func (mr *MockPodOrchestratorMockRecorder) OrchestratePodCreation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrchestratePodCreation", reflect.TypeOf((*MockPodOrchestrator)(nil).OrchestratePodCreation), arg0)
}