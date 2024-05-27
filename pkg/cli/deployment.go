package cli

import (
	"maden/pkg/shared"
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
	Long: `Fetches the currently active Maden deployments by hitting the API server`,
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
	Long: `Deletes a Maden deployment by ID. For example:
	
maden delete deployment 1234

This command will delete the deployment with ID 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentID := args[0]

		continueDelete := addDeploymentConfirmationPrompt(deploymentID)
		if !continueDelete {
			return
		}
		
		err := deleteDeployment(deploymentID)
		if err != nil {
			fmt.Printf("Error deleting deployment: %s\n", err)
			return
		}
		fmt.Printf("Deployment %s deleted successfully\n", deploymentID)
	},
}

func addDeploymentConfirmationPrompt(deploymentID string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Warning: This will delete deployment %s and all associated pods. Continue? (y/n): ", deploymentID)
	
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

func deleteDeployment(deploymentID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/deployments/%s", deploymentID), nil)
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
	Long: `Restarts a Maden deployment by name. For example:

maden rollout restart example-deployment

This command will restart the deployment named 'example-deployment' in the system.`,
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

func init() {
	getCmd.AddCommand(getDeploymentsCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	rootCmd.AddCommand(rolloutCmd)
	rolloutCmd.AddCommand(rolloutRestartDeploymentCmd)
}
