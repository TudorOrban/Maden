package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maden/pkg/mocks"
	"maden/pkg/shared"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentHandlerListDeploymentsHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockDeploymentRepository(ctrl)
	mockUpdateController := mocks.NewMockDeploymentUpdaterController(ctrl)
    handler := NewDeploymentHandler(mockRepo, mockUpdateController)

    // Prepare mock data
    deployments := []shared.Deployment{{ID: "1", Name: "Deployment1"}}
    mockRepo.EXPECT().ListDeployments().Return(deployments, nil)

    // Create a request and response recorder
    req, err := http.NewRequest("GET", "/deployments", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    // Invoke the handler
    handler.listDeploymentsHandler(rr, req)

    // Check the status code and response body
    assert.Equal(t, http.StatusOK, rr.Code)
    expectedBytes, _ := json.Marshal(deployments)
    assert.Equal(t, string(expectedBytes)+"\n", rr.Body.String())
}

func TestDeploymentHandlerDeleteDeploymentHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockDeploymentRepository(ctrl)
	mockUpdateController := mocks.NewMockDeploymentUpdaterController(ctrl)
    handler := NewDeploymentHandler(mockRepo, mockUpdateController)

    deploymentName := "test-dep"

    // Prepare the request and response recorder
    req, err := http.NewRequest("DELETE", "/deployments/"+deploymentName, nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{"name": deploymentName})
    rr := httptest.NewRecorder()

    // Expectations and call
    mockRepo.EXPECT().DeleteDeployment(deploymentName).Return(nil)

    handler.deleteDeploymentHandler(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)

    // Test not found error
    mockRepo.EXPECT().DeleteDeployment(deploymentName).Return(&shared.ErrNotFound{})
    rr = httptest.NewRecorder()

    handler.deleteDeploymentHandler(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeploymentHandlerRolloutRestartDeploymentHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockDeploymentRepository(ctrl)
    mockUpdateController := mocks.NewMockDeploymentUpdaterController(ctrl)
    handler := NewDeploymentHandler(mockRepo, mockUpdateController)

    deploymentName := "test-deployment"
    deployment := &shared.Deployment{Name: deploymentName, Replicas: 2}

    req, err := http.NewRequest("POST", "/deployments/"+deploymentName+"/rollout-restart", nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{"name": deploymentName})
    rr := httptest.NewRecorder()

    // Successful restart
    mockRepo.EXPECT().GetDeploymentByName(deploymentName).Return(deployment, nil)
    mockUpdateController.EXPECT().HandleDeploymentRolloutRestart(deployment).Return(nil)

    handler.rolloutRestartDeploymentHandler(rr, req)
    assert.Equal(t, http.StatusNoContent, rr.Code)

    // Test not found error
    rr = httptest.NewRecorder()
    mockRepo.EXPECT().GetDeploymentByName(deploymentName).Return(nil, &shared.ErrNotFound{})

    handler.rolloutRestartDeploymentHandler(rr, req)
    assert.Equal(t, http.StatusNotFound, rr.Code)

    // Test internal server error
    rr = httptest.NewRecorder()
    mockRepo.EXPECT().GetDeploymentByName(deploymentName).Return(deployment, nil)
    mockUpdateController.EXPECT().HandleDeploymentRolloutRestart(deployment).Return(fmt.Errorf("failed to restart"))

    handler.rolloutRestartDeploymentHandler(rr, req)
    assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeploymentHandlerScaleDeploymentHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockDeploymentRepository(ctrl)
    handler := NewDeploymentHandler(mockRepo, nil) 

    deploymentName := "test-deployment"
    deployment := &shared.Deployment{Name: deploymentName, Replicas: 2}
    newReplicas := 5
    requestBody, _ := json.Marshal(shared.ScaleRequest{Replicas: newReplicas})

    // Test Case: Successful scaling
    t.Run("success", func(t *testing.T) {
        req, err := http.NewRequest("POST", "/deployments/"+deploymentName+"/scale", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req = mux.SetURLVars(req, map[string]string{"name": deploymentName})
        rr := httptest.NewRecorder()

        mockRepo.EXPECT().GetDeploymentByName(deploymentName).Return(deployment, nil).Times(1)
        mockRepo.EXPECT().UpdateDeployment(deployment).Return(nil).Times(1)

        handler.scaleDeploymentHandler(rr, req)
        assert.Equal(t, http.StatusNoContent, rr.Code)
    })

    // Test Case: Not found error
    t.Run("not found", func(t *testing.T) {
        req, err := http.NewRequest("POST", "/deployments/"+deploymentName+"/scale", bytes.NewBuffer(requestBody))
        if err != nil {
            t.Fatal(err)
        }
        req = mux.SetURLVars(req, map[string]string{"name": deploymentName})
        rr := httptest.NewRecorder()

        mockRepo.EXPECT().GetDeploymentByName(deploymentName).Return(nil, &shared.ErrNotFound{}).Times(1)

        handler.scaleDeploymentHandler(rr, req)
        assert.Equal(t, http.StatusNotFound, rr.Code)
    })
}
