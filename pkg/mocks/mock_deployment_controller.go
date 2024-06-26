// Code generated by MockGen. DO NOT EDIT.
// Source: maden/pkg/controller (interfaces: DeploymentController)

// Package mocks is a generated GoMock package.
package mocks

import (
	shared "maden/pkg/shared"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDeploymentController is a mock of DeploymentController interface.
type MockDeploymentController struct {
	ctrl     *gomock.Controller
	recorder *MockDeploymentControllerMockRecorder
}

// MockDeploymentControllerMockRecorder is the mock recorder for MockDeploymentController.
type MockDeploymentControllerMockRecorder struct {
	mock *MockDeploymentController
}

// NewMockDeploymentController creates a new mock instance.
func NewMockDeploymentController(ctrl *gomock.Controller) *MockDeploymentController {
	mock := &MockDeploymentController{ctrl: ctrl}
	mock.recorder = &MockDeploymentControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeploymentController) EXPECT() *MockDeploymentControllerMockRecorder {
	return m.recorder
}

// HandleIncomingDeployment mocks base method.
func (m *MockDeploymentController) HandleIncomingDeployment(arg0 shared.DeploymentSpec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleIncomingDeployment", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleIncomingDeployment indicates an expected call of HandleIncomingDeployment.
func (mr *MockDeploymentControllerMockRecorder) HandleIncomingDeployment(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIncomingDeployment", reflect.TypeOf((*MockDeploymentController)(nil).HandleIncomingDeployment), arg0)
}
