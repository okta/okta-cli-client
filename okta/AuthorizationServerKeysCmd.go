package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthorizationServerKeysCmd = &cobra.Command{
	Use:  "authorizationServerKeys",
	Long: "Manage AuthorizationServerKeysAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerKeysCmd)
}

var ListAuthorizationServerKeysauthServerId string

func NewListAuthorizationServerKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerKeysAPI.ListAuthorizationServerKeys(apiClient.GetConfig().Context, ListAuthorizationServerKeysauthServerId)

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

	cmd.Flags().StringVarP(&ListAuthorizationServerKeysauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	ListAuthorizationServerKeysCmd := NewListAuthorizationServerKeysCmd()
	AuthorizationServerKeysCmd.AddCommand(ListAuthorizationServerKeysCmd)
}

var (
	RotateAuthorizationServerKeysauthServerId string

	RotateAuthorizationServerKeysdata string
)

func NewRotateAuthorizationServerKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "rotate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerKeysAPI.RotateAuthorizationServerKeys(apiClient.GetConfig().Context, RotateAuthorizationServerKeysauthServerId)

			if RotateAuthorizationServerKeysdata != "" {
				req = req.Data(RotateAuthorizationServerKeysdata)
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

	cmd.Flags().StringVarP(&RotateAuthorizationServerKeysauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&RotateAuthorizationServerKeysdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	RotateAuthorizationServerKeysCmd := NewRotateAuthorizationServerKeysCmd()
	AuthorizationServerKeysCmd.AddCommand(RotateAuthorizationServerKeysCmd)
}
