package apiserver

import (
	"maden/pkg/controller"
	"maden/pkg/shared"

	"encoding/json"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

// type ManifestHandler struct {
// 	Controller controller.DeploymentController
// }

// func NewManifestHandler(controller controller.DeploymentController) *ManifestHandler {
// 	return &ManifestHandler{Controller: controller}
// }


func handleMadenResources(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	decoder := yaml.NewDecoder(bytes.NewReader(body))

	for {
		var resource shared.MadenResource
		err := decoder.Decode(&resource)
		if err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, "Failed to parse YAML: " + err.Error(), http.StatusBadRequest)
			return
		}

		switch resource.Kind {
		case "Deployment":
			err := handleDeployment(resource)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "Service":
			err := handleService(resource)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			fmt.Fprintf(w, "Unsupported kind: %s", resource.Kind)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func handleDeployment(resource shared.MadenResource) error {
	var deploymentSpec shared.DeploymentSpec
	specBytes, err := json.Marshal(resource.Spec)
	if err != nil {
		fmt.Println("Error marshaling deployment spec: ", err)
		return err
	}

	err = json.Unmarshal(specBytes, &deploymentSpec)
	if err != nil {
		fmt.Println("Error unmarshaling deployment spec: ", err)
		return err
	}

	fmt.Printf("Handling Deployment: %+v\n", deploymentSpec)

	return controller.HandleIncomingDeployment(deploymentSpec)
}

func handleService(resource shared.MadenResource) error {
	var serviceSpec shared.ServiceSpec
	specBytes, err := json.Marshal(resource.Spec)
	if err != nil {
		fmt.Println("Error marshaling service spec: ", err)
		return err
	}

	err = json.Unmarshal(specBytes, &serviceSpec)
	if err != nil {
		fmt.Println("Error unmarshaling service spec: ", err)
		return err
	}

	fmt.Printf("Handling Service: %+v\n", serviceSpec)

	return controller.HandleIncomingService(serviceSpec)
}