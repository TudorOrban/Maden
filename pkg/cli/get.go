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
	Use:   "get",
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

func init() {
	rootCmd.AddCommand(getPodsCmd)

}
