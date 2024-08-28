package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationConnectionsCmd = &cobra.Command{
	Use:  "applicationConnections",
	Long: "Manage ApplicationConnectionsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationConnectionsCmd)
}

var (
	UpdateDefaultProvisioningConnectionForApplicationappId string

	UpdateDefaultProvisioningConnectionForApplicationdata string
)

func NewUpdateDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateDefaultProvisioningConnectionForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.UpdateDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, UpdateDefaultProvisioningConnectionForApplicationappId)

			if UpdateDefaultProvisioningConnectionForApplicationdata != "" {
				req = req.Data(UpdateDefaultProvisioningConnectionForApplicationdata)
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

	cmd.Flags().StringVarP(&UpdateDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateDefaultProvisioningConnectionForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateDefaultProvisioningConnectionForApplicationCmd := NewUpdateDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(UpdateDefaultProvisioningConnectionForApplicationCmd)
}

var GetDefaultProvisioningConnectionForApplicationappId string

func NewGetDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getDefaultProvisioningConnectionForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.GetDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, GetDefaultProvisioningConnectionForApplicationappId)

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

	cmd.Flags().StringVarP(&GetDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	GetDefaultProvisioningConnectionForApplicationCmd := NewGetDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(GetDefaultProvisioningConnectionForApplicationCmd)
}

var (
	ActivateDefaultProvisioningConnectionForApplicationappId string

	ActivateDefaultProvisioningConnectionForApplicationdata string
)

func NewActivateDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activateDefaultProvisioningConnectionForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.ActivateDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, ActivateDefaultProvisioningConnectionForApplicationappId)

			if ActivateDefaultProvisioningConnectionForApplicationdata != "" {
				req = req.Data(ActivateDefaultProvisioningConnectionForApplicationdata)
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

	cmd.Flags().StringVarP(&ActivateDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&ActivateDefaultProvisioningConnectionForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivateDefaultProvisioningConnectionForApplicationCmd := NewActivateDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(ActivateDefaultProvisioningConnectionForApplicationCmd)
}

var (
	DeactivateDefaultProvisioningConnectionForApplicationappId string

	DeactivateDefaultProvisioningConnectionForApplicationdata string
)

func NewDeactivateDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivateDefaultProvisioningConnectionForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.DeactivateDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, DeactivateDefaultProvisioningConnectionForApplicationappId)

			if DeactivateDefaultProvisioningConnectionForApplicationdata != "" {
				req = req.Data(DeactivateDefaultProvisioningConnectionForApplicationdata)
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

	cmd.Flags().StringVarP(&DeactivateDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&DeactivateDefaultProvisioningConnectionForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	DeactivateDefaultProvisioningConnectionForApplicationCmd := NewDeactivateDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(DeactivateDefaultProvisioningConnectionForApplicationCmd)
}

var (
	VerifyProvisioningConnectionForApplicationappName string

	VerifyProvisioningConnectionForApplicationappId string

	VerifyProvisioningConnectionForApplicationdata string
)

func NewVerifyProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "verifyProvisioningConnectionForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.VerifyProvisioningConnectionForApplication(apiClient.GetConfig().Context, VerifyProvisioningConnectionForApplicationappName, VerifyProvisioningConnectionForApplicationappId)

			if VerifyProvisioningConnectionForApplicationdata != "" {
				req = req.Data(VerifyProvisioningConnectionForApplicationdata)
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

	cmd.Flags().StringVarP(&VerifyProvisioningConnectionForApplicationappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&VerifyProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&VerifyProvisioningConnectionForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	VerifyProvisioningConnectionForApplicationCmd := NewVerifyProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(VerifyProvisioningConnectionForApplicationCmd)
}
