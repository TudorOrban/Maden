package cli

import (
	"maden/pkg/shared"

	"bytes"
	"strconv"
	"strings"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getDeploymentsCmd = &cobra.Command{
	Use: "deployment",
	Short: "Fetches current Maden deployments",
	Long: `Fetches and displays the currently active Maden deployments along with their details`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8080/deployments")
		if err != nil {
			fmt.Println("Error fetching data: ", err)
			return
		}
		defer response.Body.Close()
		
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response: ", err)
			return
		}
		
		var deployments []shared.Deployment
		if err := json.Unmarshal(body, &deployments); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		displayDeployments(deployments)
	},
}

func displayDeployments(deployments []shared.Deployment) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Replicas", "Container Port"})
	table.SetBorder(false)

	for _, deployment := range deployments {
		table.Append([]string{
			deployment.ID,
			deployment.Name,
			fmt.Sprint(deployment.Replicas),
			fmt.Sprint(deployment.Template.Spec.Containers[0].Ports[0].ContainerPort),
		})
	}

	table.Render()
}


var deleteDeploymentCmd = &cobra.Command{
	Use: "deployment [deploymentID]",
	Short: "Deletes a Maden deployment",
	Long: `Deletes a Maden deployment by name. For example:
	
maden delete deployment 1234

This command will delete the deployment with name 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentName := args[0]

		continueDelete := addDeploymentConfirmationPrompt(deploymentName)
		if !continueDelete {
			return
		}
		
		err := deleteDeployment(deploymentName)
		if err != nil {
			fmt.Printf("Error deleting deployment: %s\n", err)
			return
		}
		fmt.Printf("Deployment %s deleted successfully\n", deploymentName)
	},
}

func addDeploymentConfirmationPrompt(deploymentName string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Warning: This will delete deployment %s and all associated pods. Continue? (y/n): ", deploymentName)
	
	response, err := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		return false
	}
	if strings.ToLower(response) != "y" {
		fmt.Println("Deletion aborted.")
		return false
	}

	return true
}

func deleteDeployment(deploymentName string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/deployments/%s", deploymentName), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete deployment with status: %s", response.Status)
	}

	return nil
}

var rolloutCmd = &cobra.Command{
    Use:   "rollout",
    Short: "Manage rollouts",
    Long:  `Manage rollouts and their configurations.`,
}

var rolloutRestartDeploymentCmd = &cobra.Command{
	Use:   "restart [deploymentName]",
	Short: "Restarts a Maden deployment",
	Long: `Restarts a Maden deployment by name, by deleting and recreating all associated pods. For example:

maden rollout restart example-deployment

This command will restart the deployment named 'example-deployment'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentName := args[0]

		err := rolloutRestartDeployment(deploymentName)
		if err != nil {
			fmt.Printf("Error restarting deployment: %s\n", err)
			return
		}
		fmt.Printf("Deployment '%s' restarted successfully\n", deploymentName)
	},
}

func rolloutRestartDeployment(deploymentName string) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/deployments/%s/rollout-restart", deploymentName), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to rollout restart deployment with status: %s", response.Status)
	}

	return nil
}

var scaleDeploymentCmd = &cobra.Command{
	Use:   "scale [deploymentName] [replicas]",
	Short: "Scales a Maden deployment",
	Long: `Scales a Maden deployment by name to the specified number of replicas, deleting or creating pods as necessary. For example:

maden scale example-deployment 3

This command will scale the deployment named 'example-deployment' to 3 replicas.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentName := args[0]
		replicas, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid number of replicas: %s\n", args[1])
			return
		}

		err = scaleDeployment(deploymentName, replicas)
		if err != nil {
			fmt.Printf("Error scaling deployment: %s\n", err)
			return
		}
		fmt.Printf("Deployment '%s' scaled successfully to %d replicas\n", deploymentName, replicas)
	},
}

func scaleDeployment(deploymentName string, replicas int) error {
	scaleRequest := shared.ScaleRequest{Replicas: replicas}
	requestBody, err := json.Marshal(scaleRequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/deployments/%s/scale", deploymentName), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to scale deployment with status %s: %s", response.Status, string(responseBody))
	}

	return nil
}

func init() {
	getCmd.AddCommand(getDeploymentsCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	rootCmd.AddCommand(rolloutCmd)
	rolloutCmd.AddCommand(rolloutRestartDeploymentCmd)
	rootCmd.AddCommand(scaleDeploymentCmd)
}
