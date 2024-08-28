package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationPoliciesCmd = &cobra.Command{
	Use:  "applicationPolicies",
	Long: "Manage ApplicationPoliciesAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationPoliciesCmd)
}

var (
	AssignApplicationPolicyappId string

	AssignApplicationPolicypolicyId string

	AssignApplicationPolicydata string
)

func NewAssignApplicationPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "assignApplicationPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationPoliciesAPI.AssignApplicationPolicy(apiClient.GetConfig().Context, AssignApplicationPolicyappId, AssignApplicationPolicypolicyId)

			if AssignApplicationPolicydata != "" {
				req = req.Data(AssignApplicationPolicydata)
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
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignApplicationPolicyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&AssignApplicationPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&AssignApplicationPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AssignApplicationPolicyCmd := NewAssignApplicationPolicyCmd()
	ApplicationPoliciesCmd.AddCommand(AssignApplicationPolicyCmd)
}
