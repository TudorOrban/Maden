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

var getNodesCmd = &cobra.Command{
	Use: "node",
	Short: "Fetches current Maden nodes",
	Long: `Fetches the currently active Maden nodes by hitting the API server`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8080/nodes")
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
		
		var nodes []shared.Node
		if err := json.Unmarshal(body, &nodes); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "CPU (Capacity)", "CPU (Used)", "Memory (Capacity)", "Memory (Used)"})
		table.SetBorder(false)

		for _, node := range nodes {
			table.Append([]string{
				node.ID,
				node.Name,
				node.Status.String(),
				fmt.Sprint(node.Capacity.CPU),
				fmt.Sprint(node.Used.CPU),
				fmt.Sprint(node.Capacity.Memory),
				fmt.Sprint(node.Used.Memory),
			})
		}

		table.Render()
	},
}

var deleteNodeCmd = &cobra.Command{
	Use: "node [nodeID]",
	Short: "Deletes a Maden node",
	Long: `Deletes a Maden node by ID. For example:
	
maden delete node 1234

This command will delete the node with ID 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nodeID := args[0]
		err := deleteNode(nodeID)
		if err != nil {
			fmt.Printf("Error deleting node: %s\n", err)
			return
		}
		fmt.Printf("Node %s deleted successfully\n", nodeID)
	},
}


func deleteNode(nodeID string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/nodes/%s", nodeID), nil)
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
		return fmt.Errorf("failed to delete node with status: %s", response.Status)
	}

	return nil
}

func init() {
	getCmd.AddCommand(getNodesCmd)
	deleteCmd.AddCommand(deleteNodeCmd)

}
