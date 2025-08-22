package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

var (
	GetAuthenticatorSettingsBackupDir string

	GetAuthenticatorSettingsQuiet bool
)

func NewGetAuthenticatorSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getAuthenticatorSettings",
		Long: "Retrieve the Authenticator Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.GetAuthenticatorSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAuthenticatorSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if cmd.Flags().Changed("backup") {
				dirPath := filepath.Join(GetAuthenticatorSettingsBackupDir, "attackprotection", "getAuthenticatorSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "attackprotection.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAuthenticatorSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAuthenticatorSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the AttackProtection to a file")

	cmd.Flags().StringVarP(&GetAuthenticatorSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAuthenticatorSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAuthenticatorSettingsCmd := NewGetAuthenticatorSettingsCmd()
	AttackProtectionCmd.AddCommand(GetAuthenticatorSettingsCmd)
}

var (
	ReplaceAuthenticatorSettingsdata string

	ReplaceAuthenticatorSettingsQuiet bool
)

func NewReplaceAuthenticatorSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceAuthenticatorSettings",
		Long: "Replace the Authenticator Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.ReplaceAuthenticatorSettings(apiClient.GetConfig().Context)

			if ReplaceAuthenticatorSettingsdata != "" {
				req = req.Data(ReplaceAuthenticatorSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceAuthenticatorSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceAuthenticatorSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceAuthenticatorSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceAuthenticatorSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceAuthenticatorSettingsCmd := NewReplaceAuthenticatorSettingsCmd()
	AttackProtectionCmd.AddCommand(ReplaceAuthenticatorSettingsCmd)
}

var (
	GetUserLockoutSettingsBackupDir string

	GetUserLockoutSettingsQuiet bool
)

func NewGetUserLockoutSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getUserLockoutSettings",
		Long: "Retrieve the User Lockout Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.GetUserLockoutSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetUserLockoutSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if cmd.Flags().Changed("backup") {
				dirPath := filepath.Join(GetUserLockoutSettingsBackupDir, "attackprotection", "getUserLockoutSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "attackprotection.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetUserLockoutSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetUserLockoutSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the AttackProtection to a file")

	cmd.Flags().StringVarP(&GetUserLockoutSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetUserLockoutSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetUserLockoutSettingsCmd := NewGetUserLockoutSettingsCmd()
	AttackProtectionCmd.AddCommand(GetUserLockoutSettingsCmd)
}

var (
	ReplaceUserLockoutSettingsdata string

	ReplaceUserLockoutSettingsQuiet bool
)

func NewReplaceUserLockoutSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceUserLockoutSettings",
		Long: "Replace the User Lockout Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AttackProtectionAPI.ReplaceUserLockoutSettings(apiClient.GetConfig().Context)

			if ReplaceUserLockoutSettingsdata != "" {
				req = req.Data(ReplaceUserLockoutSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceUserLockoutSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceUserLockoutSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceUserLockoutSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceUserLockoutSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceUserLockoutSettingsCmd := NewReplaceUserLockoutSettingsCmd()
	AttackProtectionCmd.AddCommand(ReplaceUserLockoutSettingsCmd)
}
