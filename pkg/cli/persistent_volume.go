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

var getPersistentVolumesCmd = &cobra.Command{
	Use:   "persistentVolume",
	Short: "Fetches current Maden persistentVolumes",
	Long:  `Fetches and displays the currently active Maden persistentVolumes along with their details`,
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8080/persistentVolumes")
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

		var persistentVolumes []shared.PersistentVolume
		if err := json.Unmarshal(body, &persistentVolumes); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		displayPersistentVolumes(persistentVolumes)
	},
}

func displayPersistentVolumes(persistentVolumes []shared.PersistentVolume) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Replicas", "Container Port"})
	table.SetBorder(false)

	for _, persistentVolume := range persistentVolumes {
		table.Append([]string{
			persistentVolume.ID,
			persistentVolume.Name,
			fmt.Sprint(persistentVolume.Capacity),
			fmt.Sprint(persistentVolume.StorageClassName),
		})
	}

	table.Render()
}

var deletePersistentVolumeCmd = &cobra.Command{
	Use:   "persistentVolume [persistentVolumeID]",
	Short: "Deletes a Maden persistentVolume",
	Long: `Deletes a Maden persistentVolume by name. For example:
	
maden delete persistentVolume 1234

This command will delete the persistentVolume with name 1234 from the system`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		persistentVolumeName := args[0]

		continueDelete := addPersistentVolumeConfirmationPrompt(persistentVolumeName)
		if !continueDelete {
			return
		}

		err := deletePersistentVolume(persistentVolumeName)
		if err != nil {
			fmt.Printf("Error deleting persistentVolume: %s\n", err)
			return
		}
		fmt.Printf("PersistentVolume %s deleted successfully\n", persistentVolumeName)
	},
}

func addPersistentVolumeConfirmationPrompt(persistentVolumeName string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Warning: This will delete persistentVolume %s and all associated pods. Continue? (y/n): ", persistentVolumeName)

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

func deletePersistentVolume(persistentVolumeName string) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/persistentVolumes/%s", persistentVolumeName), nil)
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
		return fmt.Errorf("failed to delete persistentVolume with status: %s", response.Status)
	}

	return nil
}

func init() {
	getCmd.AddCommand(getPersistentVolumesCmd)
	deleteCmd.AddCommand(deletePersistentVolumeCmd)
	rootCmd.AddCommand(rolloutCmd)
}
