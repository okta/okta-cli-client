package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var RateLimitSettingsCmd = &cobra.Command{
	Use:  "rateLimitSettings",
	Long: "Manage RateLimitSettingsAPI",
}

func init() {
	rootCmd.AddCommand(RateLimitSettingsCmd)
}

var (
	GetRateLimitSettingsAdminNotificationsBackupDir string

	GetRateLimitSettingsAdminNotificationsQuiet bool
)

func NewGetRateLimitSettingsAdminNotificationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getAdminNotifications",
		Long: "Retrieve the Rate Limit Admin Notification Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RateLimitSettingsAPI.GetRateLimitSettingsAdminNotifications(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRateLimitSettingsAdminNotificationsQuiet {
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
				dirPath := filepath.Join(GetRateLimitSettingsAdminNotificationsBackupDir, "ratelimitsettings", "getAdminNotifications")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "ratelimitsettings.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRateLimitSettingsAdminNotificationsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRateLimitSettingsAdminNotificationsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the RateLimitSettings to a file")

	cmd.Flags().StringVarP(&GetRateLimitSettingsAdminNotificationsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRateLimitSettingsAdminNotificationsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRateLimitSettingsAdminNotificationsCmd := NewGetRateLimitSettingsAdminNotificationsCmd()
	RateLimitSettingsCmd.AddCommand(GetRateLimitSettingsAdminNotificationsCmd)
}

var (
	ReplaceRateLimitSettingsAdminNotificationsdata string

	ReplaceRateLimitSettingsAdminNotificationsQuiet bool
)

func NewReplaceRateLimitSettingsAdminNotificationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceAdminNotifications",
		Long: "Replace the Rate Limit Admin Notification Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RateLimitSettingsAPI.ReplaceRateLimitSettingsAdminNotifications(apiClient.GetConfig().Context)

			if ReplaceRateLimitSettingsAdminNotificationsdata != "" {
				req = req.Data(ReplaceRateLimitSettingsAdminNotificationsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRateLimitSettingsAdminNotificationsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRateLimitSettingsAdminNotificationsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRateLimitSettingsAdminNotificationsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRateLimitSettingsAdminNotificationsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRateLimitSettingsAdminNotificationsCmd := NewReplaceRateLimitSettingsAdminNotificationsCmd()
	RateLimitSettingsCmd.AddCommand(ReplaceRateLimitSettingsAdminNotificationsCmd)
}

var (
	GetRateLimitSettingsPerClientBackupDir string

	GetRateLimitSettingsPerClientQuiet bool
)

func NewGetRateLimitSettingsPerClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPerClient",
		Long: "Retrieve the Per-Client Rate Limit Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RateLimitSettingsAPI.GetRateLimitSettingsPerClient(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRateLimitSettingsPerClientQuiet {
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
				dirPath := filepath.Join(GetRateLimitSettingsPerClientBackupDir, "ratelimitsettings", "getPerClient")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "ratelimitsettings.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRateLimitSettingsPerClientQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRateLimitSettingsPerClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the RateLimitSettings to a file")

	cmd.Flags().StringVarP(&GetRateLimitSettingsPerClientBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRateLimitSettingsPerClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRateLimitSettingsPerClientCmd := NewGetRateLimitSettingsPerClientCmd()
	RateLimitSettingsCmd.AddCommand(GetRateLimitSettingsPerClientCmd)
}

var (
	ReplaceRateLimitSettingsPerClientdata string

	ReplaceRateLimitSettingsPerClientQuiet bool
)

func NewReplaceRateLimitSettingsPerClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePerClient",
		Long: "Replace the Per-Client Rate Limit Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RateLimitSettingsAPI.ReplaceRateLimitSettingsPerClient(apiClient.GetConfig().Context)

			if ReplaceRateLimitSettingsPerClientdata != "" {
				req = req.Data(ReplaceRateLimitSettingsPerClientdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRateLimitSettingsPerClientQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRateLimitSettingsPerClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRateLimitSettingsPerClientdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRateLimitSettingsPerClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRateLimitSettingsPerClientCmd := NewReplaceRateLimitSettingsPerClientCmd()
	RateLimitSettingsCmd.AddCommand(ReplaceRateLimitSettingsPerClientCmd)
}

var (
	GetRateLimitSettingsWarningThresholdBackupDir string

	GetRateLimitSettingsWarningThresholdQuiet bool
)

func NewGetRateLimitSettingsWarningThresholdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getWarningThreshold",
		Long: "Retrieve the Rate Limit Warning Threshold Percentage",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RateLimitSettingsAPI.GetRateLimitSettingsWarningThreshold(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRateLimitSettingsWarningThresholdQuiet {
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
				dirPath := filepath.Join(GetRateLimitSettingsWarningThresholdBackupDir, "ratelimitsettings", "getWarningThreshold")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "ratelimitsettings.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRateLimitSettingsWarningThresholdQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRateLimitSettingsWarningThresholdQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the RateLimitSettings to a file")

	cmd.Flags().StringVarP(&GetRateLimitSettingsWarningThresholdBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRateLimitSettingsWarningThresholdQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRateLimitSettingsWarningThresholdCmd := NewGetRateLimitSettingsWarningThresholdCmd()
	RateLimitSettingsCmd.AddCommand(GetRateLimitSettingsWarningThresholdCmd)
}

var (
	ReplaceRateLimitSettingsWarningThresholddata string

	ReplaceRateLimitSettingsWarningThresholdQuiet bool
)

func NewReplaceRateLimitSettingsWarningThresholdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceWarningThreshold",
		Long: "Replace the Rate Limit Warning Threshold Percentage",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RateLimitSettingsAPI.ReplaceRateLimitSettingsWarningThreshold(apiClient.GetConfig().Context)

			if ReplaceRateLimitSettingsWarningThresholddata != "" {
				req = req.Data(ReplaceRateLimitSettingsWarningThresholddata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRateLimitSettingsWarningThresholdQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRateLimitSettingsWarningThresholdQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRateLimitSettingsWarningThresholddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRateLimitSettingsWarningThresholdQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRateLimitSettingsWarningThresholdCmd := NewReplaceRateLimitSettingsWarningThresholdCmd()
	RateLimitSettingsCmd.AddCommand(ReplaceRateLimitSettingsWarningThresholdCmd)
}
