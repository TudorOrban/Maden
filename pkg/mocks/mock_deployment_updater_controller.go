// Code generated by MockGen. DO NOT EDIT.
// Source: maden/pkg/controller (interfaces: DeploymentUpdaterController)

// Package mocks is a generated GoMock package.
package mocks

import (
	shared "maden/pkg/shared"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
)

// MockDeploymentUpdaterController is a mock of DeploymentUpdaterController interface.
type MockDeploymentUpdaterController struct {
	ctrl     *gomock.Controller
	recorder *MockDeploymentUpdaterControllerMockRecorder
}

// MockDeploymentUpdaterControllerMockRecorder is the mock recorder for MockDeploymentUpdaterController.
type MockDeploymentUpdaterControllerMockRecorder struct {
	mock *MockDeploymentUpdaterController
}

// NewMockDeploymentUpdaterController creates a new mock instance.
func NewMockDeploymentUpdaterController(ctrl *gomock.Controller) *MockDeploymentUpdaterController {
	mock := &MockDeploymentUpdaterController{ctrl: ctrl}
	mock.recorder = &MockDeploymentUpdaterControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeploymentUpdaterController) EXPECT() *MockDeploymentUpdaterControllerMockRecorder {
	return m.recorder
}

// HandleDeploymentCreate mocks base method.
func (m *MockDeploymentUpdaterController) HandleDeploymentCreate(arg0 *mvccpb.KeyValue) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleDeploymentCreate", arg0)
}

// HandleDeploymentCreate indicates an expected call of HandleDeploymentCreate.
func (mr *MockDeploymentUpdaterControllerMockRecorder) HandleDeploymentCreate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDeploymentCreate", reflect.TypeOf((*MockDeploymentUpdaterController)(nil).HandleDeploymentCreate), arg0)
}

// HandleDeploymentDelete mocks base method.
func (m *MockDeploymentUpdaterController) HandleDeploymentDelete(arg0 *mvccpb.KeyValue) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleDeploymentDelete", arg0)
}

// HandleDeploymentDelete indicates an expected call of HandleDeploymentDelete.
func (mr *MockDeploymentUpdaterControllerMockRecorder) HandleDeploymentDelete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDeploymentDelete", reflect.TypeOf((*MockDeploymentUpdaterController)(nil).HandleDeploymentDelete), arg0)
}

// HandleDeploymentRolloutRestart mocks base method.
func (m *MockDeploymentUpdaterController) HandleDeploymentRolloutRestart(arg0 *shared.Deployment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleDeploymentRolloutRestart", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleDeploymentRolloutRestart indicates an expected call of HandleDeploymentRolloutRestart.
func (mr *MockDeploymentUpdaterControllerMockRecorder) HandleDeploymentRolloutRestart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDeploymentRolloutRestart", reflect.TypeOf((*MockDeploymentUpdaterController)(nil).HandleDeploymentRolloutRestart), arg0)
}

// HandleDeploymentUpdate mocks base method.
func (m *MockDeploymentUpdaterController) HandleDeploymentUpdate(arg0, arg1 *mvccpb.KeyValue) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleDeploymentUpdate", arg0, arg1)
}

// HandleDeploymentUpdate indicates an expected call of HandleDeploymentUpdate.
func (mr *MockDeploymentUpdaterControllerMockRecorder) HandleDeploymentUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDeploymentUpdate", reflect.TypeOf((*MockDeploymentUpdaterController)(nil).HandleDeploymentUpdate), arg0, arg1)
}
