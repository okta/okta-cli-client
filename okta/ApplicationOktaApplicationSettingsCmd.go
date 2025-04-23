package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationOktaApplicationSettingsCmd = &cobra.Command{
	Use:  "applicationOktaApplicationSettings",
	Long: "Manage ApplicationOktaApplicationSettingsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationOktaApplicationSettingsCmd)
}

var (
	GetFirstPartyAppSettingsappName string

	GetFirstPartyAppSettingsBackupDir string

	GetFirstPartyAppSettingsQuiet bool
)

func NewGetFirstPartyAppSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getFirstPartyAppSettings",
		Long: "Retrieve the Okta app settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationOktaApplicationSettingsAPI.GetFirstPartyAppSettings(apiClient.GetConfig().Context, GetFirstPartyAppSettingsappName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetFirstPartyAppSettingsQuiet {
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
				dirPath := filepath.Join(GetFirstPartyAppSettingsBackupDir, "applicationoktaapplicationsettings", "getFirstPartyAppSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetFirstPartyAppSettingsappName
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetFirstPartyAppSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetFirstPartyAppSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetFirstPartyAppSettingsappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationOktaApplicationSettings to a file")

	cmd.Flags().StringVarP(&GetFirstPartyAppSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetFirstPartyAppSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetFirstPartyAppSettingsCmd := NewGetFirstPartyAppSettingsCmd()
	ApplicationOktaApplicationSettingsCmd.AddCommand(GetFirstPartyAppSettingsCmd)
}

var (
	ReplaceFirstPartyAppSettingsappName string

	ReplaceFirstPartyAppSettingsdata string

	ReplaceFirstPartyAppSettingsQuiet bool
)

func NewReplaceFirstPartyAppSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceFirstPartyAppSettings",
		Long: "Replace the Okta app settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationOktaApplicationSettingsAPI.ReplaceFirstPartyAppSettings(apiClient.GetConfig().Context, ReplaceFirstPartyAppSettingsappName)

			if ReplaceFirstPartyAppSettingsdata != "" {
				req = req.Data(ReplaceFirstPartyAppSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceFirstPartyAppSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceFirstPartyAppSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceFirstPartyAppSettingsappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&ReplaceFirstPartyAppSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceFirstPartyAppSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceFirstPartyAppSettingsCmd := NewReplaceFirstPartyAppSettingsCmd()
	ApplicationOktaApplicationSettingsCmd.AddCommand(ReplaceFirstPartyAppSettingsCmd)
}
