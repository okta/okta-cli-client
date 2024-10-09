package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var TrustedOriginCmd = &cobra.Command{
	Use:  "trustedOrigin",
	Long: "Manage TrustedOriginAPI",
}

func init() {
	rootCmd.AddCommand(TrustedOriginCmd)
}

var CreateTrustedOrigindata string

func NewCreateTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.CreateTrustedOrigin(apiClient.GetConfig().Context)

			if CreateTrustedOrigindata != "" {
				req = req.Data(CreateTrustedOrigindata)
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

	cmd.Flags().StringVarP(&CreateTrustedOrigindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateTrustedOriginCmd := NewCreateTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(CreateTrustedOriginCmd)
}

func NewListTrustedOriginsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.ListTrustedOrigins(apiClient.GetConfig().Context)

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
	ListTrustedOriginsCmd := NewListTrustedOriginsCmd()
	TrustedOriginCmd.AddCommand(ListTrustedOriginsCmd)
}

var GetTrustedOrigintrustedOriginId string

func NewGetTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.GetTrustedOrigin(apiClient.GetConfig().Context, GetTrustedOrigintrustedOriginId)

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

	cmd.Flags().StringVarP(&GetTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	return cmd
}

func init() {
	GetTrustedOriginCmd := NewGetTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(GetTrustedOriginCmd)
}

var (
	ReplaceTrustedOrigintrustedOriginId string

	ReplaceTrustedOrigindata string
)

func NewReplaceTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.ReplaceTrustedOrigin(apiClient.GetConfig().Context, ReplaceTrustedOrigintrustedOriginId)

			if ReplaceTrustedOrigindata != "" {
				req = req.Data(ReplaceTrustedOrigindata)
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

	cmd.Flags().StringVarP(&ReplaceTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	cmd.Flags().StringVarP(&ReplaceTrustedOrigindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceTrustedOriginCmd := NewReplaceTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(ReplaceTrustedOriginCmd)
}

var DeleteTrustedOrigintrustedOriginId string

func NewDeleteTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.DeleteTrustedOrigin(apiClient.GetConfig().Context, DeleteTrustedOrigintrustedOriginId)

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

	cmd.Flags().StringVarP(&DeleteTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	return cmd
}

func init() {
	DeleteTrustedOriginCmd := NewDeleteTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(DeleteTrustedOriginCmd)
}

var ActivateTrustedOrigintrustedOriginId string

func NewActivateTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.ActivateTrustedOrigin(apiClient.GetConfig().Context, ActivateTrustedOrigintrustedOriginId)

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

	cmd.Flags().StringVarP(&ActivateTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	return cmd
}

func init() {
	ActivateTrustedOriginCmd := NewActivateTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(ActivateTrustedOriginCmd)
}

var DeactivateTrustedOrigintrustedOriginId string

func NewDeactivateTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.DeactivateTrustedOrigin(apiClient.GetConfig().Context, DeactivateTrustedOrigintrustedOriginId)

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

	cmd.Flags().StringVarP(&DeactivateTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	return cmd
}

func init() {
	DeactivateTrustedOriginCmd := NewDeactivateTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(DeactivateTrustedOriginCmd)
}
