package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationOktaApplicationSettingsCmd = &cobra.Command{
	Use:  "applicationOktaApplicationSettings",
	Long: "Manage ApplicationOktaApplicationSettingsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationOktaApplicationSettingsCmd)
}

var GetFirstPartyAppSettingsappName string

func NewGetFirstPartyAppSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getFirstPartyAppSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationOktaApplicationSettingsAPI.GetFirstPartyAppSettings(apiClient.GetConfig().Context, GetFirstPartyAppSettingsappName)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetFirstPartyAppSettingsappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	return cmd
}

func init() {
	GetFirstPartyAppSettingsCmd := NewGetFirstPartyAppSettingsCmd()
	ApplicationOktaApplicationSettingsCmd.AddCommand(GetFirstPartyAppSettingsCmd)
}

var (
	ReplaceFirstPartyAppSettingsappName string

	ReplaceFirstPartyAppSettingsdata string
)

func NewReplaceFirstPartyAppSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceFirstPartyAppSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationOktaApplicationSettingsAPI.ReplaceFirstPartyAppSettings(apiClient.GetConfig().Context, ReplaceFirstPartyAppSettingsappName)

			if ReplaceFirstPartyAppSettingsdata != "" {
				req = req.Data(ReplaceFirstPartyAppSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceFirstPartyAppSettingsappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&ReplaceFirstPartyAppSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceFirstPartyAppSettingsCmd := NewReplaceFirstPartyAppSettingsCmd()
	ApplicationOktaApplicationSettingsCmd.AddCommand(ReplaceFirstPartyAppSettingsCmd)
}
