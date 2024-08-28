package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationSSOCmd = &cobra.Command{
	Use:  "applicationSSO",
	Long: "Manage ApplicationSSOAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationSSOCmd)
}

var PreviewSAMLmetadataForApplicationappId string

func NewPreviewSAMLmetadataForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "previewSAMLmetadataForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationSSOAPI.PreviewSAMLmetadataForApplication(apiClient.GetConfig().Context, PreviewSAMLmetadataForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&PreviewSAMLmetadataForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	PreviewSAMLmetadataForApplicationCmd := NewPreviewSAMLmetadataForApplicationCmd()
	ApplicationSSOCmd.AddCommand(PreviewSAMLmetadataForApplicationCmd)
}
