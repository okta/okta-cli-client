package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var SystemLogCmd = &cobra.Command{
	Use:  "systemLog",
	Long: "Manage SystemLogAPI",
}

func init() {
	rootCmd.AddCommand(SystemLogCmd)
}

func NewListLogEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listLogEvents",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SystemLogAPI.ListLogEvents(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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

	return cmd
}

func init() {
	ListLogEventsCmd := NewListLogEventsCmd()
	SystemLogCmd.AddCommand(ListLogEventsCmd)
}
