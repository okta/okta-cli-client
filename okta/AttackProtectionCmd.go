package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AttackProtectionCmd = &cobra.Command{
	Use:  "attackProtection",
	Long: "Manage AttackProtectionAPI",
}

func init() {
	rootCmd.AddCommand(AttackProtectionCmd)
}

func NewGetAuthenticatorSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getAuthenticatorSettings",
		Long: "Retrieve the Authenticator Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.GetAuthenticatorSettings(apiClient.GetConfig().Context)

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
	GetAuthenticatorSettingsCmd := NewGetAuthenticatorSettingsCmd()
	AttackProtectionCmd.AddCommand(GetAuthenticatorSettingsCmd)
}

var ReplaceAuthenticatorSettingsdata string

func NewReplaceAuthenticatorSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceAuthenticatorSettings",
		Long: "Replace the Authenticator Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.ReplaceAuthenticatorSettings(apiClient.GetConfig().Context)

			if ReplaceAuthenticatorSettingsdata != "" {
				req = req.Data(ReplaceAuthenticatorSettingsdata)
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

	cmd.Flags().StringVarP(&ReplaceAuthenticatorSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceAuthenticatorSettingsCmd := NewReplaceAuthenticatorSettingsCmd()
	AttackProtectionCmd.AddCommand(ReplaceAuthenticatorSettingsCmd)
}

func NewGetUserLockoutSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getUserLockoutSettings",
		Long: "Retrieve the User Lockout Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.GetUserLockoutSettings(apiClient.GetConfig().Context)

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
	GetUserLockoutSettingsCmd := NewGetUserLockoutSettingsCmd()
	AttackProtectionCmd.AddCommand(GetUserLockoutSettingsCmd)
}

var ReplaceUserLockoutSettingsdata string

func NewReplaceUserLockoutSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceUserLockoutSettings",
		Long: "Replace the User Lockout Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.ReplaceUserLockoutSettings(apiClient.GetConfig().Context)

			if ReplaceUserLockoutSettingsdata != "" {
				req = req.Data(ReplaceUserLockoutSettingsdata)
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

	cmd.Flags().StringVarP(&ReplaceUserLockoutSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceUserLockoutSettingsCmd := NewReplaceUserLockoutSettingsCmd()
	AttackProtectionCmd.AddCommand(ReplaceUserLockoutSettingsCmd)
}
