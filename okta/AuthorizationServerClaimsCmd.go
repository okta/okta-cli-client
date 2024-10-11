package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthorizationServerClaimsCmd = &cobra.Command{
	Use:  "authorizationServerClaims",
	Long: "Manage AuthorizationServerClaimsAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerClaimsCmd)
}

var (
	CreateOAuth2ClaimauthServerId string

	CreateOAuth2Claimdata string
)

func NewCreateOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createOAuth2Claim",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.CreateOAuth2Claim(apiClient.GetConfig().Context, CreateOAuth2ClaimauthServerId)

			if CreateOAuth2Claimdata != "" {
				req = req.Data(CreateOAuth2Claimdata)
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

	cmd.Flags().StringVarP(&CreateOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateOAuth2Claimdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateOAuth2ClaimCmd := NewCreateOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(CreateOAuth2ClaimCmd)
}

var ListOAuth2ClaimsauthServerId string

func NewListOAuth2ClaimsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listOAuth2Claims",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.ListOAuth2Claims(apiClient.GetConfig().Context, ListOAuth2ClaimsauthServerId)

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

	cmd.Flags().StringVarP(&ListOAuth2ClaimsauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	return cmd
}

func init() {
	ListOAuth2ClaimsCmd := NewListOAuth2ClaimsCmd()
	AuthorizationServerClaimsCmd.AddCommand(ListOAuth2ClaimsCmd)
}

var (
	GetOAuth2ClaimauthServerId string

	GetOAuth2ClaimclaimId string
)

func NewGetOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOAuth2Claim",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.GetOAuth2Claim(apiClient.GetConfig().Context, GetOAuth2ClaimauthServerId, GetOAuth2ClaimclaimId)

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

	cmd.Flags().StringVarP(&GetOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetOAuth2ClaimclaimId, "claimId", "", "", "")
	cmd.MarkFlagRequired("claimId")

	return cmd
}

func init() {
	GetOAuth2ClaimCmd := NewGetOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(GetOAuth2ClaimCmd)
}

var (
	ReplaceOAuth2ClaimauthServerId string

	ReplaceOAuth2ClaimclaimId string

	ReplaceOAuth2Claimdata string
)

func NewReplaceOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceOAuth2Claim",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.ReplaceOAuth2Claim(apiClient.GetConfig().Context, ReplaceOAuth2ClaimauthServerId, ReplaceOAuth2ClaimclaimId)

			if ReplaceOAuth2Claimdata != "" {
				req = req.Data(ReplaceOAuth2Claimdata)
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

	cmd.Flags().StringVarP(&ReplaceOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceOAuth2ClaimclaimId, "claimId", "", "", "")
	cmd.MarkFlagRequired("claimId")

	cmd.Flags().StringVarP(&ReplaceOAuth2Claimdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceOAuth2ClaimCmd := NewReplaceOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(ReplaceOAuth2ClaimCmd)
}

var (
	DeleteOAuth2ClaimauthServerId string

	DeleteOAuth2ClaimclaimId string
)

func NewDeleteOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteOAuth2Claim",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.DeleteOAuth2Claim(apiClient.GetConfig().Context, DeleteOAuth2ClaimauthServerId, DeleteOAuth2ClaimclaimId)

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

	cmd.Flags().StringVarP(&DeleteOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteOAuth2ClaimclaimId, "claimId", "", "", "")
	cmd.MarkFlagRequired("claimId")

	return cmd
}

func init() {
	DeleteOAuth2ClaimCmd := NewDeleteOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(DeleteOAuth2ClaimCmd)
}
