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

type ManifestHandler struct {
	DController controller.DeploymentController
	SController controller.ServiceController
}

func NewManifestHandler(
	dController controller.DeploymentController,
	sController controller.ServiceController,
) *ManifestHandler {
	return &ManifestHandler{DController: dController, SController: sController}
}


func (h *ManifestHandler) handleMadenResources(w http.ResponseWriter, r *http.Request) {
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
			err := h.handleDeployment(resource)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "Service":
			err := h.handleService(resource)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			fmt.Fprintf(w, "Unsupported kind: %s", resource.Kind)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ManifestHandler) handleDeployment(resource shared.MadenResource) error {
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

	return h.DController.HandleIncomingDeployment(deploymentSpec)
}

func (h *ManifestHandler) handleService(resource shared.MadenResource) error {
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

	return h.SController.HandleIncomingService(serviceSpec)
}