package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationGrantsCmd = &cobra.Command{
	Use:  "applicationGrants",
	Long: "Manage ApplicationGrantsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationGrantsCmd)
}

var (
	GrantConsentToScopeappId string

	GrantConsentToScopedata string
)

func NewGrantConsentToScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "grantConsentToScope",
		Long: "Grant consent to scope",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.GrantConsentToScope(apiClient.GetConfig().Context, GrantConsentToScopeappId)

			if GrantConsentToScopedata != "" {
				req = req.Data(GrantConsentToScopedata)
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

	cmd.Flags().StringVarP(&GrantConsentToScopeappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GrantConsentToScopedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	GrantConsentToScopeCmd := NewGrantConsentToScopeCmd()
	ApplicationGrantsCmd.AddCommand(GrantConsentToScopeCmd)
}

var ListScopeConsentGrantsappId string

func NewListScopeConsentGrantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listScopeConsentGrants",
		Long: "List all app Grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.ListScopeConsentGrants(apiClient.GetConfig().Context, ListScopeConsentGrantsappId)

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

	cmd.Flags().StringVarP(&ListScopeConsentGrantsappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	ListScopeConsentGrantsCmd := NewListScopeConsentGrantsCmd()
	ApplicationGrantsCmd.AddCommand(ListScopeConsentGrantsCmd)
}

var (
	GetScopeConsentGrantappId string

	GetScopeConsentGrantgrantId string
)

func NewGetScopeConsentGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getScopeConsentGrant",
		Long: "Retrieve an app Grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.GetScopeConsentGrant(apiClient.GetConfig().Context, GetScopeConsentGrantappId, GetScopeConsentGrantgrantId)

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

	cmd.Flags().StringVarP(&GetScopeConsentGrantappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetScopeConsentGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	return cmd
}

func init() {
	GetScopeConsentGrantCmd := NewGetScopeConsentGrantCmd()
	ApplicationGrantsCmd.AddCommand(GetScopeConsentGrantCmd)
}

var (
	RevokeScopeConsentGrantappId string

	RevokeScopeConsentGrantgrantId string
)

func NewRevokeScopeConsentGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeScopeConsentGrant",
		Long: "Revoke an app Grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.RevokeScopeConsentGrant(apiClient.GetConfig().Context, RevokeScopeConsentGrantappId, RevokeScopeConsentGrantgrantId)

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

	cmd.Flags().StringVarP(&RevokeScopeConsentGrantappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&RevokeScopeConsentGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	return cmd
}

func init() {
	RevokeScopeConsentGrantCmd := NewRevokeScopeConsentGrantCmd()
	ApplicationGrantsCmd.AddCommand(RevokeScopeConsentGrantCmd)
}
