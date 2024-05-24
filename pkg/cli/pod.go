package cli

import (
	"bufio"
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
	Use: "pod",
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

		displayPods(pods)
	},
}

func displayPods(pods []shared.Pod) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Status", "Node ID", "CPU", "Memory (MB)"})
	table.SetBorder(false)

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
}

var deletePodCmd = &cobra.Command{
	Use: "pod [podID]",
	Short: "Deletes a Maden pod",
	Long: `Deletes a Maden pod by ID. For example:
	
maden delete pod 1234

This command will delete the pod with ID 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podID := args[0]
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

var logsCmd = &cobra.Command{
	Use: "logs [podID] [containerID]",
	Short: "Fetch logs for a specific pod and container",
	Long: `Fetches logs for a specific pod and container by ID. For example:

maden get logs 1234 5678

This command will fetch logs for container 5678 in pod 1234.
If the pod contains only one container, the containerID can be omitted.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		podID := args[0]
		containerID := ""
		if len(args) > 1 {
			containerID = args[1]
		}
		follow, _ := cmd.Flags().GetBool("follow")
	
		url := fmt.Sprintf("http://localhost:8080/pods/%s/logs?containerID=%s&follow=%t", podID, containerID, follow)
		client := &http.Client{
			Timeout: 0,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request: ", err)
			return
		}
	
		// Set the request to close the connection upon completion
		req.Close = false
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching logs: ", err)
			return
		}
		defer response.Body.Close()
	
		// Use bufio.Reader to handle streaming more robustly
		reader := bufio.NewReader(response.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Error reading logs: %v\n", err)
				}
				break // Exit the loop if we reach the end of the stream or an error occurs
			}
			fmt.Print(line) // Print each line of logs
		}

		fmt.Println("Stream ended. Press Ctrl+C to terminate.")
		select {} 
	},
}

func init() {
	getCmd.AddCommand(getPodsCmd)
	deleteCmd.AddCommand(deletePodCmd)
	getCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolP("follow", "f", false, "Follow the logs")
}
