package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApiTokenCmd = &cobra.Command{
	Use:  "apiToken",
	Long: "Manage ApiTokenAPI",
}

func init() {
	rootCmd.AddCommand(ApiTokenCmd)
}

func NewListApiTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.ListApiTokens(apiClient.GetConfig().Context)

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
	ListApiTokensCmd := NewListApiTokensCmd()
	ApiTokenCmd.AddCommand(ListApiTokensCmd)
}

func NewRevokeCurrentApiTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "revokeCurrent",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.RevokeCurrentApiToken(apiClient.GetConfig().Context)

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
	RevokeCurrentApiTokenCmd := NewRevokeCurrentApiTokenCmd()
	ApiTokenCmd.AddCommand(RevokeCurrentApiTokenCmd)
}

var GetApiTokenapiTokenId string

func NewGetApiTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.GetApiToken(apiClient.GetConfig().Context, GetApiTokenapiTokenId)

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

	cmd.Flags().StringVarP(&GetApiTokenapiTokenId, "apiTokenId", "", "", "")
	cmd.MarkFlagRequired("apiTokenId")

	return cmd
}

func init() {
	GetApiTokenCmd := NewGetApiTokenCmd()
	ApiTokenCmd.AddCommand(GetApiTokenCmd)
}

var RevokeApiTokenapiTokenId string

func NewRevokeApiTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "revoke",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.RevokeApiToken(apiClient.GetConfig().Context, RevokeApiTokenapiTokenId)

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

	cmd.Flags().StringVarP(&RevokeApiTokenapiTokenId, "apiTokenId", "", "", "")
	cmd.MarkFlagRequired("apiTokenId")

	return cmd
}

func init() {
	RevokeApiTokenCmd := NewRevokeApiTokenCmd()
	ApiTokenCmd.AddCommand(RevokeApiTokenCmd)
}
