package cli

import (
	"maden/pkg/shared"

	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getPersistentVolumeClaimsCmd = &cobra.Command{
	Use:   "persistentVolumeClaim",
	Short: "Fetches current Maden persistentVolumeClaims",
	Long:  `Fetches and displays the currently active Maden persistentVolumeClaims along with their details`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8080/persistentVolumeClaims")
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

		var persistentVolumeClaims []shared.PersistentVolumeClaim
		if err := json.Unmarshal(body, &persistentVolumeClaims); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		displayPersistentVolumeClaims(persistentVolumeClaims)
	},
}

func displayPersistentVolumeClaims(persistentVolumeClaims []shared.PersistentVolumeClaim) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Replicas", "Container Port"})
	table.SetBorder(false)

	for _, persistentVolumeClaim := range persistentVolumeClaims {
		table.Append([]string{
			persistentVolumeClaim.ID,
			persistentVolumeClaim.Name,
			fmt.Sprint(persistentVolumeClaim.VolumeName),
			fmt.Sprint(persistentVolumeClaim.AccessModes),
		})
	}

	table.Render()
}

var deletePersistentVolumeClaimCmd = &cobra.Command{
	Use:   "persistentVolumeClaim [persistentVolumeClaimID]",
	Short: "Deletes a Maden persistentVolumeClaim",
	Long: `Deletes a Maden persistentVolumeClaim by name. For example:
	
maden delete persistentVolumeClaim 1234

This command will delete the persistentVolumeClaim with name 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		persistentVolumeClaimName := args[0]

		continueDelete := addPersistentVolumeClaimConfirmationPrompt(persistentVolumeClaimName)
		if !continueDelete {
			return
		}

		err := deletePersistentVolumeClaim(persistentVolumeClaimName)
		if err != nil {
			fmt.Printf("Error deleting persistentVolumeClaim: %s\n", err)
			return
		}
		fmt.Printf("PersistentVolumeClaim %s deleted successfully\n", persistentVolumeClaimName)
	},
}

func addPersistentVolumeClaimConfirmationPrompt(persistentVolumeClaimName string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Warning: This will delete persistentVolumeClaim %s and all associated pods. Continue? (y/n): ", persistentVolumeClaimName)

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

func deletePersistentVolumeClaim(persistentVolumeClaimName string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/persistentVolumeClaims/%s", persistentVolumeClaimName), nil)
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
		return fmt.Errorf("failed to delete persistentVolumeClaim with status: %s", response.Status)
	}

	return nil
}

func init() {
	getCmd.AddCommand(getPersistentVolumeClaimsCmd)
	deleteCmd.AddCommand(deletePersistentVolumeClaimCmd)
	rootCmd.AddCommand(rolloutCmd)
}
