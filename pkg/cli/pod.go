package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"maden/pkg/shared"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getPodsCmd = &cobra.Command{
	Use: "get pod",
	Short: "Fetches current Maden pods",
	Long: `Fetches the currently active Maden pods by hitting the API server`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8080/pods")
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
		
		var pods []shared.Pod
		if err := json.Unmarshal(body, &pods); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "Node ID", "CPU", "Memory (MB)"})

		for _, pod := range pods {
			table.Append([]string{
				pod.ID,
				pod.Name,
				pod.Status.String(),
				pod.NodeID,
				fmt.Sprint(pod.Resources.CPU),
				fmt.Sprint(pod.Resources.Memory),
			})
		}

		table.Render()
	},
}

var deletePodCmd = &cobra.Command{
	Use: "delete pod [podID]",
	Short: "Deletes a Maden pod",
	Long: `Deletes a Maden pod by ID. For example:
	
maden delete pod 1234

This command will delete the pod with ID 1234 from the system`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		podID := args[1]
		err := deletePod(podID)
		if err != nil {
			fmt.Printf("Error deleting pod: %s\n", err)
			return
		}
		fmt.Printf("Pod %s deleted successfully\n", podID)
	},
}


func deletePod(podID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/pods/%s", podID), nil)
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
		return fmt.Errorf("failed to delete pod with status: %s", response.Status)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getPodsCmd)
	rootCmd.AddCommand(deletePodCmd)

}
