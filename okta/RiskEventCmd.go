package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var RiskEventCmd = &cobra.Command{
	Use:  "riskEvent",
	Long: "Manage RiskEventAPI",
}

func init() {
	rootCmd.AddCommand(RiskEventCmd)
}

var (
	SendRiskEventsdata string

	SendRiskEventsQuiet bool
)

func NewSendRiskEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "sends",
		Long: "Send multiple Risk Events",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RiskEventAPI.SendRiskEvents(apiClient.GetConfig().Context)

			if SendRiskEventsdata != "" {
				req = req.Data(SendRiskEventsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !SendRiskEventsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !SendRiskEventsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&SendRiskEventsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&SendRiskEventsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	SendRiskEventsCmd := NewSendRiskEventsCmd()
	RiskEventCmd.AddCommand(SendRiskEventsCmd)
}
