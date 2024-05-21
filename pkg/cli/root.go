package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "madencli",
	Short: "Maden is a container orchestration tool",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:....`,
}

var getCmd = &cobra.Command{
	Use: "get",
	Short: "Get resources",
	Long: `Get resources from the Maden API server`,
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete resources",
	Long: `Delete resources from the Maden API server`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.madencli.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".madencli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
