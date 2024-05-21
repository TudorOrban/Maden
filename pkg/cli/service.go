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

var getServicesCmd = &cobra.Command{
	Use: "service",
	Short: "Fetches current Maden services",
	Long: `Fetches the currently active Maden services by hitting the API server`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8080/services")
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
		
		var services []shared.Service
		if err := json.Unmarshal(body, &services); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		displayServices(services)
	},
}

func displayServices(services []shared.Service) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Port", "Target Port"})
	table.SetBorder(false)

	for _, service := range services {
		table.Append([]string{
			service.ID,
			service.Name,
			fmt.Sprint(service.Ports[0].Port),
			fmt.Sprint(service.Ports[0].TargetPort),
		})
	}

	table.Render()
}


var deleteServiceCmd = &cobra.Command{
	Use: "service [serviceID]",
	Short: "Deletes a Maden service",
	Long: `Deletes a Maden service by ID. For example:
	
maden delete service 1234

This command will delete the service with ID 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]

		continueDelete := addServiceConfirmationPrompt(serviceID)
		if !continueDelete {
			return
		}
		
		err := deleteService(serviceID)
		if err != nil {
			fmt.Printf("Error deleting service: %s\n", err)
			return
		}
		fmt.Printf("Service %s deleted successfully\n", serviceID)
	},
}

func addServiceConfirmationPrompt(serviceID string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Warning: This will delete service %s. Continue? (y/n): ", serviceID)
	
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

func deleteService(serviceID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/services/%s", serviceID), nil)
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
		return fmt.Errorf("failed to delete service with status: %s", response.Status)
	}

	return nil
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	deleteCmd.AddCommand(deleteServiceCmd)

}
