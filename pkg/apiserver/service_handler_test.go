package apiserver

import (
	"encoding/json"
	"errors"
	"maden/pkg/mocks"
	"maden/pkg/shared"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestServiceHandlerListServicesHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockServiceRepository(ctrl)
    handler := NewServiceHandler(mockRepo, nil)

    // Prepare mock data
    services := []shared.Service{{ID: "1", Name: "Service1"}}
    mockRepo.EXPECT().ListServices().Return(services, nil)

    // Create a request and response recorder
    req, err := http.NewRequest("GET", "/services", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    // Invoke the handler
    handler.listServicesHandler(rr, req)

    // Check the status code and response body
    assert.Equal(t, http.StatusOK, rr.Code)
    expectedBytes, _ := json.Marshal(services)
    assert.Equal(t, string(expectedBytes)+"\n", rr.Body.String())
}

func TestServiceHandlerDeleteServiceHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockOrchestrator := mocks.NewMockServiceOrchestrator(ctrl)
    handler := NewServiceHandler(nil, mockOrchestrator)

    serviceName := "example-service"

    // Test case 1: Successful deletion
    mockOrchestrator.EXPECT().OrchestrateServiceDeletion(serviceName).Return(nil)

    req, err := http.NewRequest("DELETE", "/services/"+serviceName, nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{"name": serviceName})

    rr := httptest.NewRecorder()

    handler.deleteServiceHandler(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)

    // Test case 2: Service not found
    mockOrchestrator.EXPECT().OrchestrateServiceDeletion(serviceName).Return(&shared.ErrNotFound{})

    rr = httptest.NewRecorder()

    handler.deleteServiceHandler(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)

    // Test case 3: Other errors
    mockOrchestrator.EXPECT().OrchestrateServiceDeletion(serviceName).Return(errors.New("internal error"))

    rr = httptest.NewRecorder()

    handler.deleteServiceHandler(rr, req)

    assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
