package apiserver

import (
	"maden/pkg/controller"
	"maden/pkg/shared"

	"bytes"
	"encoding/json"
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
			http.Error(w, "Failed to parse YAML: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = h.handleIncomingResource(resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ManifestHandler) handleIncomingResource(resource shared.MadenResource) error {
	switch resource.Kind {
	case "Deployment":
		err := h.handleIncomingDeployment(resource)
		if err != nil {
			return err
		}
	case "Service":
		err := h.handleIncomingService(resource)
		if err != nil {
			return err
		}
	case "PersistentVolume":
		err := h.handleIncomingPersistentVolume(resource)
		if err != nil {
			return err
		}
	case "PersistentVolumeClaim":
		err := h.handleIncomingPersistentVolumeClaim(resource)
		if err != nil {
			return err
		}
	default:
		errorMsg := fmt.Sprintf("Unsupported kind: %s", resource.Kind)
		return fmt.Errorf(errorMsg)
	}

	return nil
}

func (h *ManifestHandler) handleIncomingDeployment(resource shared.MadenResource) error {
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

func (h *ManifestHandler) handleIncomingService(resource shared.MadenResource) error {
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

func (h *ManifestHandler) handleIncomingPersistentVolume(resource shared.MadenResource) error {
	var pvSpec shared.PersistentVolumeSpec
	specBytes, err := json.Marshal(resource.Spec)
	if err != nil {
		return err
	}

	err = json.Unmarshal(specBytes, &pvSpec)
	if err != nil {
		return err
	}

	// Implement logic to handle the persistent volume lifecycle
	return nil
}

func (h *ManifestHandler) handleIncomingPersistentVolumeClaim(resource shared.MadenResource) error {
	var pvcSpec shared.PersistentVolumeClaimSpec
	specBytes, err := json.Marshal(resource.Spec)
	if err != nil {
		return err
	}

	err = json.Unmarshal(specBytes, &pvcSpec)
	if err != nil {
		return err
	}

	// Implement logic to bind PVC to PV and manage the lifecycle
	return nil
}
