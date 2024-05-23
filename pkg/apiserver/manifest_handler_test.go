package apiserver

import (
	"bytes"
	"maden/pkg/mocks"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestManifestHandlerHandleMadenResources(t *testing.T) {
	// Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDeploymentController := mocks.NewMockDeploymentController(ctrl)
    mockServiceController := mocks.NewMockServiceController(ctrl)
    handler := NewManifestHandler(mockDeploymentController, mockServiceController)

    deploymentYAML := `
kind: Deployment
spec:
  name: test-deployment
  replicas: 3
`
    serviceYAML := `
kind: Service
spec:
  name: test-service
`
    malformedYAML := `kind: Unknown\n`

    // Test valid deployment handling
    req, _ := http.NewRequest("POST", "/maden-resources", bytes.NewBufferString(deploymentYAML))
    rr := httptest.NewRecorder()

    mockDeploymentController.EXPECT().
        HandleIncomingDeployment(gomock.Any()).
        Return(nil).Times(1)

    handler.handleMadenResources(rr, req)
    assert.Equal(t, http.StatusCreated, rr.Code)

    // Test valid service handling
    req, _ = http.NewRequest("POST", "/maden-resources", bytes.NewBufferString(serviceYAML))
    rr = httptest.NewRecorder()

    mockServiceController.EXPECT().
        HandleIncomingService(gomock.Any()).
        Return(nil).Times(1)

    handler.handleMadenResources(rr, req)
    assert.Equal(t, http.StatusCreated, rr.Code)

    // Test malformed input
    req, _ = http.NewRequest("POST", "/maden-resources", bytes.NewBufferString(malformedYAML))
    rr = httptest.NewRecorder()

    handler.handleMadenResources(rr, req)

    assert.Equal(t, http.StatusBadRequest, rr.Code)
}
