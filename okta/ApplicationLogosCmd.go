package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationLogosCmd = &cobra.Command{
	Use:  "applicationLogos",
	Long: "Manage ApplicationLogosAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationLogosCmd)
}

var (
	UploadApplicationLogoappId string

	UploadApplicationLogodata string

	UploadApplicationLogoQuiet bool
)

func NewUploadApplicationLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadApplicationLogo",
		Long: "Upload an application Logo",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationLogosAPI.UploadApplicationLogo(apiClient.GetConfig().Context, UploadApplicationLogoappId)

			if UploadApplicationLogodata != "" {
				req = req.Data(UploadApplicationLogodata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadApplicationLogoQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadApplicationLogoQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadApplicationLogoappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UploadApplicationLogodata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadApplicationLogoQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadApplicationLogoCmd := NewUploadApplicationLogoCmd()
	ApplicationLogosCmd.AddCommand(UploadApplicationLogoCmd)
}
