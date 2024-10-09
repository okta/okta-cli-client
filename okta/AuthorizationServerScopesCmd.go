package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthorizationServerScopesCmd = &cobra.Command{
	Use:  "authorizationServerScopes",
	Long: "Manage AuthorizationServerScopesAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerScopesCmd)
}

var (
	CreateOAuth2ScopeauthServerId string

	CreateOAuth2Scopedata string
)

func NewCreateOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createOAuth2Scope",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.CreateOAuth2Scope(apiClient.GetConfig().Context, CreateOAuth2ScopeauthServerId)

			if CreateOAuth2Scopedata != "" {
				req = req.Data(CreateOAuth2Scopedata)
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

	cmd.Flags().StringVarP(&CreateOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateOAuth2Scopedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateOAuth2ScopeCmd := NewCreateOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(CreateOAuth2ScopeCmd)
}

var ListOAuth2ScopesauthServerId string

func NewListOAuth2ScopesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listOAuth2Scopes",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.ListOAuth2Scopes(apiClient.GetConfig().Context, ListOAuth2ScopesauthServerId)

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

	cmd.Flags().StringVarP(&ListOAuth2ScopesauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	ListOAuth2ScopesCmd := NewListOAuth2ScopesCmd()
	AuthorizationServerScopesCmd.AddCommand(ListOAuth2ScopesCmd)
}

var (
	GetOAuth2ScopeauthServerId string

	GetOAuth2ScopescopeId string
)

func NewGetOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOAuth2Scope",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.GetOAuth2Scope(apiClient.GetConfig().Context, GetOAuth2ScopeauthServerId, GetOAuth2ScopescopeId)

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

	cmd.Flags().StringVarP(&GetOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetOAuth2ScopescopeId, "scopeId", "", "", "")
	cmd.MarkFlagRequired("scopeId")

	return cmd
}

func init() {
	GetOAuth2ScopeCmd := NewGetOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(GetOAuth2ScopeCmd)
}

var (
	ReplaceOAuth2ScopeauthServerId string

	ReplaceOAuth2ScopescopeId string

	ReplaceOAuth2Scopedata string
)

func NewReplaceOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceOAuth2Scope",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.ReplaceOAuth2Scope(apiClient.GetConfig().Context, ReplaceOAuth2ScopeauthServerId, ReplaceOAuth2ScopescopeId)

			if ReplaceOAuth2Scopedata != "" {
				req = req.Data(ReplaceOAuth2Scopedata)
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

	cmd.Flags().StringVarP(&ReplaceOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceOAuth2ScopescopeId, "scopeId", "", "", "")
	cmd.MarkFlagRequired("scopeId")

	cmd.Flags().StringVarP(&ReplaceOAuth2Scopedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceOAuth2ScopeCmd := NewReplaceOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(ReplaceOAuth2ScopeCmd)
}

var (
	DeleteOAuth2ScopeauthServerId string

	DeleteOAuth2ScopescopeId string
)

func NewDeleteOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteOAuth2Scope",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.DeleteOAuth2Scope(apiClient.GetConfig().Context, DeleteOAuth2ScopeauthServerId, DeleteOAuth2ScopescopeId)

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

	cmd.Flags().StringVarP(&DeleteOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteOAuth2ScopescopeId, "scopeId", "", "", "")
	cmd.MarkFlagRequired("scopeId")

	return cmd
}

func init() {
	DeleteOAuth2ScopeCmd := NewDeleteOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(DeleteOAuth2ScopeCmd)
}
