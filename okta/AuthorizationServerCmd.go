package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthorizationServerCmd = &cobra.Command{
	Use:  "authorizationServer",
	Long: "Manage AuthorizationServerAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerCmd)
}

var CreateAuthorizationServerdata string

func NewCreateAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.CreateAuthorizationServer(apiClient.GetConfig().Context)

			if CreateAuthorizationServerdata != "" {
				req = req.Data(CreateAuthorizationServerdata)
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

	cmd.Flags().StringVarP(&CreateAuthorizationServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateAuthorizationServerCmd := NewCreateAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(CreateAuthorizationServerCmd)
}

func NewListAuthorizationServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.ListAuthorizationServers(apiClient.GetConfig().Context)

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

	return cmd
}

func init() {
	ListAuthorizationServersCmd := NewListAuthorizationServersCmd()
	AuthorizationServerCmd.AddCommand(ListAuthorizationServersCmd)
}

var GetAuthorizationServerauthServerId string

func NewGetAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.GetAuthorizationServer(apiClient.GetConfig().Context, GetAuthorizationServerauthServerId)

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

	cmd.Flags().StringVarP(&GetAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	GetAuthorizationServerCmd := NewGetAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(GetAuthorizationServerCmd)
}

var (
	ReplaceAuthorizationServerauthServerId string

	ReplaceAuthorizationServerdata string
)

func NewReplaceAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.ReplaceAuthorizationServer(apiClient.GetConfig().Context, ReplaceAuthorizationServerauthServerId)

			if ReplaceAuthorizationServerdata != "" {
				req = req.Data(ReplaceAuthorizationServerdata)
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

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceAuthorizationServerCmd := NewReplaceAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(ReplaceAuthorizationServerCmd)
}

var DeleteAuthorizationServerauthServerId string

func NewDeleteAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.DeleteAuthorizationServer(apiClient.GetConfig().Context, DeleteAuthorizationServerauthServerId)

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

	cmd.Flags().StringVarP(&DeleteAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	DeleteAuthorizationServerCmd := NewDeleteAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(DeleteAuthorizationServerCmd)
}

var (
	ActivateAuthorizationServerauthServerId string

	ActivateAuthorizationServerdata string
)

func NewActivateAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.ActivateAuthorizationServer(apiClient.GetConfig().Context, ActivateAuthorizationServerauthServerId)

			if ActivateAuthorizationServerdata != "" {
				req = req.Data(ActivateAuthorizationServerdata)
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

	cmd.Flags().StringVarP(&ActivateAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ActivateAuthorizationServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivateAuthorizationServerCmd := NewActivateAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(ActivateAuthorizationServerCmd)
}

var (
	DeactivateAuthorizationServerauthServerId string

	DeactivateAuthorizationServerdata string
)

func NewDeactivateAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.DeactivateAuthorizationServer(apiClient.GetConfig().Context, DeactivateAuthorizationServerauthServerId)

			if DeactivateAuthorizationServerdata != "" {
				req = req.Data(DeactivateAuthorizationServerdata)
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

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	DeactivateAuthorizationServerCmd := NewDeactivateAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(DeactivateAuthorizationServerCmd)
}
