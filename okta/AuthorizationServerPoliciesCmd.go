package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthorizationServerPoliciesCmd = &cobra.Command{
	Use:  "authorizationServerPolicies",
	Long: "Manage AuthorizationServerPoliciesAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerPoliciesCmd)
}

var (
	CreateAuthorizationServerPolicyauthServerId string

	CreateAuthorizationServerPolicydata string
)

func NewCreateAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createAuthorizationServerPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.CreateAuthorizationServerPolicy(apiClient.GetConfig().Context, CreateAuthorizationServerPolicyauthServerId)

			if CreateAuthorizationServerPolicydata != "" {
				req = req.Data(CreateAuthorizationServerPolicydata)
			}

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

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateAuthorizationServerPolicyCmd := NewCreateAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(CreateAuthorizationServerPolicyCmd)
}

var ListAuthorizationServerPoliciesauthServerId string

func NewListAuthorizationServerPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.ListAuthorizationServerPolicies(apiClient.GetConfig().Context, ListAuthorizationServerPoliciesauthServerId)

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

	cmd.Flags().StringVarP(&ListAuthorizationServerPoliciesauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	ListAuthorizationServerPoliciesCmd := NewListAuthorizationServerPoliciesCmd()
	AuthorizationServerPoliciesCmd.AddCommand(ListAuthorizationServerPoliciesCmd)
}

var (
	GetAuthorizationServerPolicyauthServerId string

	GetAuthorizationServerPolicypolicyId string
)

func NewGetAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getAuthorizationServerPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.GetAuthorizationServerPolicy(apiClient.GetConfig().Context, GetAuthorizationServerPolicyauthServerId, GetAuthorizationServerPolicypolicyId)

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

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	GetAuthorizationServerPolicyCmd := NewGetAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(GetAuthorizationServerPolicyCmd)
}

var (
	ReplaceAuthorizationServerPolicyauthServerId string

	ReplaceAuthorizationServerPolicypolicyId string

	ReplaceAuthorizationServerPolicydata string
)

func NewReplaceAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceAuthorizationServerPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.ReplaceAuthorizationServerPolicy(apiClient.GetConfig().Context, ReplaceAuthorizationServerPolicyauthServerId, ReplaceAuthorizationServerPolicypolicyId)

			if ReplaceAuthorizationServerPolicydata != "" {
				req = req.Data(ReplaceAuthorizationServerPolicydata)
			}

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

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceAuthorizationServerPolicyCmd := NewReplaceAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(ReplaceAuthorizationServerPolicyCmd)
}

var (
	DeleteAuthorizationServerPolicyauthServerId string

	DeleteAuthorizationServerPolicypolicyId string
)

func NewDeleteAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteAuthorizationServerPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.DeleteAuthorizationServerPolicy(apiClient.GetConfig().Context, DeleteAuthorizationServerPolicyauthServerId, DeleteAuthorizationServerPolicypolicyId)

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

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	DeleteAuthorizationServerPolicyCmd := NewDeleteAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(DeleteAuthorizationServerPolicyCmd)
}

var (
	ActivateAuthorizationServerPolicyauthServerId string

	ActivateAuthorizationServerPolicypolicyId string
)

func NewActivateAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activateAuthorizationServerPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.ActivateAuthorizationServerPolicy(apiClient.GetConfig().Context, ActivateAuthorizationServerPolicyauthServerId, ActivateAuthorizationServerPolicypolicyId)

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

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	ActivateAuthorizationServerPolicyCmd := NewActivateAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(ActivateAuthorizationServerPolicyCmd)
}

var (
	DeactivateAuthorizationServerPolicyauthServerId string

	DeactivateAuthorizationServerPolicypolicyId string
)

func NewDeactivateAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivateAuthorizationServerPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.DeactivateAuthorizationServerPolicy(apiClient.GetConfig().Context, DeactivateAuthorizationServerPolicyauthServerId, DeactivateAuthorizationServerPolicypolicyId)

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

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	DeactivateAuthorizationServerPolicyCmd := NewDeactivateAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(DeactivateAuthorizationServerPolicyCmd)
}
