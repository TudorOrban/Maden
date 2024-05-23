package apiserver

import (
	"bytes"
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

func TestNodeHandlerListNodesHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockNodeRepository(ctrl)
    handler := NewNodeHandler(mockRepo)

    // Prepare mock data
    nodes := []shared.Node{{ID: "1", Name: "Node1"}}
    mockRepo.EXPECT().ListNodes().Return(nodes, nil)

    // Create a request and response recorder
    req, err := http.NewRequest("GET", "/nodes", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    // Invoke the handler
    handler.listNodesHandler(rr, req)

    // Check the status code and response body
    assert.Equal(t, http.StatusOK, rr.Code)
    expectedBytes, _ := json.Marshal(nodes)
    assert.Equal(t, string(expectedBytes)+"\n", rr.Body.String())
}

func TestNodeHandlerCreateNodeHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockNodeRepository(ctrl)
    handler := NewNodeHandler(mockRepo)

    node := shared.Node{ID: "1", Name: "Node1"}
    nodeBytes, _ := json.Marshal(node)
    reader := bytes.NewReader(nodeBytes)

    // Test successful creation
    req, err := http.NewRequest("POST", "/nodes", reader)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    mockRepo.EXPECT().CreateNode(gomock.Any()).Return(nil)
    
    handler.createNodeHandler(rr, req)
    
    assert.Equal(t, http.StatusCreated, rr.Code)
    assert.Equal(t, string(nodeBytes)+"\n", rr.Body.String())

    // Test error for bad request
    req, err = http.NewRequest("POST", "/nodes", bytes.NewReader([]byte("invalid")))
    if err != nil {
        t.Fatal(err)
    }
    rr = httptest.NewRecorder()

    handler.createNodeHandler(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestNodeHandlerDeleteNodeHandler(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockNodeRepository(ctrl)
    handler := NewNodeHandler(mockRepo)

    nodeID := "1"

    // Prepare the request and response recorder
    req, err := http.NewRequest("DELETE", "/nodes/"+nodeID, nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{"id": nodeID})
    rr := httptest.NewRecorder()

    // Expectations and call
    mockRepo.EXPECT().DeleteNode(nodeID).Return(nil)

    handler.deleteNodeHandler(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)

    // Test not found error
    mockRepo.EXPECT().DeleteNode(nodeID).Return(&shared.ErrNotFound{})
    rr = httptest.NewRecorder()

    handler.deleteNodeHandler(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}
