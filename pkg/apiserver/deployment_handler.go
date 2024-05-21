package apiserver

import (
	"maden/pkg/controller"
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"encoding/json"
	// "errors"
	"bytes"
	"fmt"
	"io"
	"net/http"

	// "github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func listDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	deployments, err := etcd.ListDeployments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}


// func createDeploymentHandler(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

	
//     fmt.Println("Received YAML:", string(body))

// 	var deploymentSpec shared.DeploymentSpec
//     if err := yaml.Unmarshal(body, &deploymentSpec); err != nil {
//         http.Error(w, "Failed to parse YAML: "+err.Error(), http.StatusBadRequest)
//         return
//     }

// 	if err := etcd.CreateDeployment(&deploymentSpec.Deployment); err != nil {
// 		var dupErr *shared.ErrDuplicateResource
// 		if errors.As(err, &dupErr) {
// 			http.Error(w, dupErr.Error(), http.StatusConflict)
// 		} else {
// 			http.Error(w, fmt.Sprintf("Error storing deployment data in etcd: %v", err), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(deploymentSpec.Deployment)
// }

// func deleteDeploymentHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	deploymentID := vars["id"]

// 	if err := etcd.DeleteDeployment(deploymentID); err != nil {
// 		if err == shared.ErrNotFound {
// 			w.WriteHeader(http.StatusNotFound)
// 		} else {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
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
			handleService(resource)
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

func handleService(resource shared.MadenResource) {
	var serviceSpec shared.ServiceSpec
	specBytes, err := json.Marshal(resource.Spec)
	if err != nil {
		fmt.Println("Error marshaling service spec: ", err)
		return
	}

	err = json.Unmarshal(specBytes, &serviceSpec)
	if err != nil {
		fmt.Println("Error unmarshaling service spec: ", err)
		return
	}

	fmt.Printf("Handling Service: %+v\n", serviceSpec)

}