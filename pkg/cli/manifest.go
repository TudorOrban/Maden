package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var filePath string

var applyCmd = &cobra.Command{
	Use:   "apply -f [yaml-filepath]",
	Short: "Apply a manifest to the Maden cluster",
	Long: `Apply a manifest yaml file to the Maden cluster to create or update resources.
For example:

maden apply -f deployment.yaml

This command will create deployments and services in the Maden cluster based on the deployment.yaml manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			fmt.Println("Error: YAML file path must be provided using the -f flag")
			os.Exit(1)
		}
		
		// Read YAML file
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			os.Exit(1)
		}

		// Send request to API server
		err = applyResources(fileContent)
		if err != nil {
			fmt.Printf("Error applying resources: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Resources applied successfully")
	},
}

func applyResources(fileContent []byte) error {
	request, err := http.NewRequest("POST", "http://localhost:8080/manifests", bytes.NewBuffer(fileContent))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-yaml")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	
	if response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to apply resources with status: %s, response: %s", response.Status, string(body))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
	
	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "YAML file path")
}