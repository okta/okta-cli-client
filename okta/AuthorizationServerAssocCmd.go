package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthorizationServerAssocCmd = &cobra.Command{
	Use:  "authorizationServerAssoc",
	Long: "Manage AuthorizationServerAssocAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerAssocCmd)
}

var (
	CreateAssociatedServersauthServerId string

	CreateAssociatedServersdata string
)

func NewCreateAssociatedServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createAssociatedServers",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAssocAPI.CreateAssociatedServers(apiClient.GetConfig().Context, CreateAssociatedServersauthServerId)

			if CreateAssociatedServersdata != "" {
				req = req.Data(CreateAssociatedServersdata)
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

	cmd.Flags().StringVarP(&CreateAssociatedServersauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateAssociatedServersdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateAssociatedServersCmd := NewCreateAssociatedServersCmd()
	AuthorizationServerAssocCmd.AddCommand(CreateAssociatedServersCmd)
}

var ListAssociatedServersByTrustedTypeauthServerId string

func NewListAssociatedServersByTrustedTypeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listAssociatedServersByTrustedType",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAssocAPI.ListAssociatedServersByTrustedType(apiClient.GetConfig().Context, ListAssociatedServersByTrustedTypeauthServerId)

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

	cmd.Flags().StringVarP(&ListAssociatedServersByTrustedTypeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	ListAssociatedServersByTrustedTypeCmd := NewListAssociatedServersByTrustedTypeCmd()
	AuthorizationServerAssocCmd.AddCommand(ListAssociatedServersByTrustedTypeCmd)
}

var (
	DeleteAssociatedServerauthServerId string

	DeleteAssociatedServerassociatedServerId string
)

func NewDeleteAssociatedServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteAssociatedServer",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAssocAPI.DeleteAssociatedServer(apiClient.GetConfig().Context, DeleteAssociatedServerauthServerId, DeleteAssociatedServerassociatedServerId)

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

	cmd.Flags().StringVarP(&DeleteAssociatedServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteAssociatedServerassociatedServerId, "associatedServerId", "", "", "")
	cmd.MarkFlagRequired("associatedServerId")

	return cmd
}

func init() {
	DeleteAssociatedServerCmd := NewDeleteAssociatedServerCmd()
	AuthorizationServerAssocCmd.AddCommand(DeleteAssociatedServerCmd)
}
