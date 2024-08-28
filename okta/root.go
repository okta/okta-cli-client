package okta

import (
	"fmt"
	"os"

	"github.com/okta/okta-cli-client/sdk"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "okta-cli-client",
	Long: "A command line tool for management API",
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "okta-cli-client",
		Long: "A command line tool for management API",
	}
	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var apiClient *sdk.APIClient

func init() {
	configuration, err := sdk.NewConfiguration(sdk.WithCache(false))
	if err != nil {
		fmt.Printf("Create new config should not be error %v", err)
	}
	configuration.Debug = false

	apiClient = sdk.NewAPIClient(configuration)
}
