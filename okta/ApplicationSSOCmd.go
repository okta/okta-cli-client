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

var (
	PreviewSAMLmetadataForApplicationappId string

	PreviewSAMLmetadataForApplicationQuiet bool
)

func NewPreviewSAMLmetadataForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "previewSAMLmetadataForApplication",
		Long: "Preview the application SAML metadata",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationSSOAPI.PreviewSAMLmetadataForApplication(apiClient.GetConfig().Context, PreviewSAMLmetadataForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !PreviewSAMLmetadataForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !PreviewSAMLmetadataForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&PreviewSAMLmetadataForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&PreviewSAMLmetadataForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	PreviewSAMLmetadataForApplicationCmd := NewPreviewSAMLmetadataForApplicationCmd()
	ApplicationSSOCmd.AddCommand(PreviewSAMLmetadataForApplicationCmd)
}
