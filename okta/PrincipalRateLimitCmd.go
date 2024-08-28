package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var PrincipalRateLimitCmd = &cobra.Command{
	Use:  "principalRateLimit",
	Long: "Manage PrincipalRateLimitAPI",
}

func init() {
	rootCmd.AddCommand(PrincipalRateLimitCmd)
}

var CreatePrincipalRateLimitEntitydata string

func NewCreatePrincipalRateLimitEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createEntity",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.CreatePrincipalRateLimitEntity(apiClient.GetConfig().Context)

			if CreatePrincipalRateLimitEntitydata != "" {
				req = req.Data(CreatePrincipalRateLimitEntitydata)
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

	cmd.Flags().StringVarP(&CreatePrincipalRateLimitEntitydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreatePrincipalRateLimitEntityCmd := NewCreatePrincipalRateLimitEntityCmd()
	PrincipalRateLimitCmd.AddCommand(CreatePrincipalRateLimitEntityCmd)
}

func NewListPrincipalRateLimitEntitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listEntities",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.ListPrincipalRateLimitEntities(apiClient.GetConfig().Context)

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

	return cmd
}

func init() {
	ListPrincipalRateLimitEntitiesCmd := NewListPrincipalRateLimitEntitiesCmd()
	PrincipalRateLimitCmd.AddCommand(ListPrincipalRateLimitEntitiesCmd)
}

var GetPrincipalRateLimitEntityprincipalRateLimitId string

func NewGetPrincipalRateLimitEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getEntity",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.GetPrincipalRateLimitEntity(apiClient.GetConfig().Context, GetPrincipalRateLimitEntityprincipalRateLimitId)

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

	cmd.Flags().StringVarP(&GetPrincipalRateLimitEntityprincipalRateLimitId, "principalRateLimitId", "", "", "")
	cmd.MarkFlagRequired("principalRateLimitId")

	return cmd
}

func init() {
	GetPrincipalRateLimitEntityCmd := NewGetPrincipalRateLimitEntityCmd()
	PrincipalRateLimitCmd.AddCommand(GetPrincipalRateLimitEntityCmd)
}

var (
	ReplacePrincipalRateLimitEntityprincipalRateLimitId string

	ReplacePrincipalRateLimitEntitydata string
)

func NewReplacePrincipalRateLimitEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceEntity",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.ReplacePrincipalRateLimitEntity(apiClient.GetConfig().Context, ReplacePrincipalRateLimitEntityprincipalRateLimitId)

			if ReplacePrincipalRateLimitEntitydata != "" {
				req = req.Data(ReplacePrincipalRateLimitEntitydata)
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

	cmd.Flags().StringVarP(&ReplacePrincipalRateLimitEntityprincipalRateLimitId, "principalRateLimitId", "", "", "")
	cmd.MarkFlagRequired("principalRateLimitId")

	cmd.Flags().StringVarP(&ReplacePrincipalRateLimitEntitydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacePrincipalRateLimitEntityCmd := NewReplacePrincipalRateLimitEntityCmd()
	PrincipalRateLimitCmd.AddCommand(ReplacePrincipalRateLimitEntityCmd)
}
