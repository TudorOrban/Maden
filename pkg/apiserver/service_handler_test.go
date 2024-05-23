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

func TestServiceHandlerListServicesHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockServiceRepository(ctrl)
    handler := NewServiceHandler(mockRepo)

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

    mockRepo := mocks.NewMockServiceRepository(ctrl)
    handler := NewServiceHandler(mockRepo)

    serviceName := "test-dep"

    // Prepare the request and response recorder
    req, err := http.NewRequest("DELETE", "/services/"+serviceName, nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{"name": serviceName})
    rr := httptest.NewRecorder()

    // Expectations and call
    mockRepo.EXPECT().DeleteService(serviceName).Return(nil)

    handler.deleteServiceHandler(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)

    // Test not found error
    mockRepo.EXPECT().DeleteService(serviceName).Return(&shared.ErrNotFound{})
    rr = httptest.NewRecorder()

    handler.deleteServiceHandler(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}
