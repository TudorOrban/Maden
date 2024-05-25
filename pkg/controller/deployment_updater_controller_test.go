package controller

import (
	"encoding/json"
	"maden/pkg/mocks"
	"maden/pkg/shared"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func TestHandleDeploymentCreate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockPodRepository(ctrl)
    mockOrch := mocks.NewMockPodOrchestrator(ctrl)
    controller := NewDefaultDeploymentUpdaterController(mockRepo, mockOrch)

    // Example JSON string with detailed template specifications
    deploymentJSON := `{
        "ID":"dep-1",
        "Name":"test-deployment",
        "Replicas":3,
        "Template":{
            "Spec":{
                "Containers":[
                    {"Image":"nginx:latest", "Ports":[{"ContainerPort":80}]}
                ],
                "Resources":{
                    "Limits":{"cpu":"100", "memory":"200"},
                    "Requests":{"cpu":"100", "memory":"200"},
					"Affinity":{},
					"AntiAffinity":{},
					"Tolerations":[],
					"RestartPolicy":"Always"
                }
            }
        }
    }`

    kv := &mvccpb.KeyValue{Value: []byte(deploymentJSON)}
    var deployment shared.Deployment
    json.Unmarshal(kv.Value, &deployment) // Assuming unmarshal works correctly for this example

    // Set expectation with detailed assertion about the pod properties
    mockOrch.EXPECT().OrchestratePodCreation(gomock.Any()).Times(deployment.Replicas).DoAndReturn(func(pod *shared.Pod) error {
        assert.Equal(t, "dep-1", pod.DeploymentID)
        assert.Equal(t, "test-deployment", pod.Name)
        assert.Equal(t, "nginx:latest", pod.Containers[0].Image)
        assert.Equal(t, 80, pod.Containers[0].Ports[0].ContainerPort)
        return nil
    })

    controller.HandleDeploymentCreate(kv)
}

func TestHandleDeploymentUpdate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockPodRepository(ctrl)
    mockOrch := mocks.NewMockPodOrchestrator(ctrl)
    controller := NewDefaultDeploymentUpdaterController(mockRepo, mockOrch)

    oldDeploymentJSON := `{
        "ID":"dep-1",
        "Name":"test-deployment",
        "Replicas":2,
        "Template":{
            "Spec":{
                "Containers":[
                    {"Image":"nginx:latest", "Ports":[{"ContainerPort":80}]}
                ],
                "Resources":{
                    "Limits":{"cpu":"100", "memory":"200"},
                    "Requests":{"cpu":"100", "memory":"200"},
					"Affinity":{},
					"AntiAffinity":{},
					"Tolerations":[],
					"RestartPolicy":"Always"
                }
            }
        }
    }`
    newDeploymentJSON := `{
        "ID":"dep-1",
        "Name":"test-deployment",
        "Replicas":4,
        "Template":{
            "Spec":{
                "Containers":[
                    {"Image":"nginx:latest", "Ports":[{"ContainerPort":80}]}
                ],
                "Resources":{
                    "Limits":{"cpu":"100", "memory":"200"},
                    "Requests":{"cpu":"100", "memory":"200"},
					"Affinity":{},
					"AntiAffinity":{},
					"Tolerations":[],
					"RestartPolicy":"Always"
                }
            }
        }
    }`

    oldKv := &mvccpb.KeyValue{Value: []byte(oldDeploymentJSON)}
    newKv := &mvccpb.KeyValue{Value: []byte(newDeploymentJSON)}

    var oldDeployment, newDeployment shared.Deployment
    json.Unmarshal(oldKv.Value, &oldDeployment)
    json.Unmarshal(newKv.Value, &newDeployment)

    difference := newDeployment.Replicas - oldDeployment.Replicas
    mockOrch.EXPECT().OrchestratePodCreation(gomock.Any()).Times(difference).DoAndReturn(func(pod *shared.Pod) error {
        assert.Equal(t, "dep-1", pod.DeploymentID)
        assert.Equal(t, "test-deployment", pod.Name)
        assert.Equal(t, "nginx:latest", pod.Containers[0].Image)
        assert.Equal(t, 80, pod.Containers[0].Ports[0].ContainerPort)
        return nil
    })

    controller.HandleDeploymentUpdate(oldKv, newKv)
}