package apiserver

import (
	"encoding/json"
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
