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

func TestPodHandlerListPodsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPodRepository(ctrl)
	handler := NewPodHandler(mockRepo, nil)

	pods := []shared.Pod{{ID: "1", Name: "test-pod"}}
	mockRepo.EXPECT().ListPods().Return(pods, nil)

	req, err := http.NewRequest("GET", "/pods", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.listPodsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBytes, _ := json.Marshal(pods)
	expected := string(expectedBytes) + "\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPodHandlerCreatePodHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPodRepository(ctrl)
	mockOrchestrator := mocks.NewMockPodOrchestrator(ctrl)
	handler := NewPodHandler(mockRepo, mockOrchestrator)

	pod := shared.Pod{ID: "1", Name: "test-pod"}
	podBytes, _ := json.Marshal(pod)
	reader := bytes.NewReader(podBytes)
	
	req, err := http.NewRequest("POST", "/pods", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockOrchestrator.EXPECT().OrchestratePodCreation(gomock.Any()).Return(nil)

	handler.createPodHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	expected := string(podBytes) + "\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPodHandlerDeletePodHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPodRepository(ctrl)
	mockOrchestrator := mocks.NewMockPodOrchestrator(ctrl)
	handler := NewPodHandler(mockRepo, mockOrchestrator)

	pod := shared.Pod{ID: "1", Name: "test-pod"}
	mockRepo.EXPECT().GetPodByID(pod.ID).Return(&pod, nil)
	mockOrchestrator.EXPECT().OrchestratePodDeletion(&pod).Return(nil)

	req, err := http.NewRequest("DELETE", "/pods/"+pod.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": pod.ID})

	rr := httptest.NewRecorder()

	handler.deletePodHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
