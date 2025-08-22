package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ThreatInsightCmd = &cobra.Command{
	Use:  "threatInsight",
	Long: "Manage ThreatInsightAPI",
}

func init() {
	rootCmd.AddCommand(ThreatInsightCmd)
}

var (
	UpdateConfigurationdata string

	UpdateConfigurationQuiet bool
)

func NewUpdateConfigurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateConfiguration",
		Long: "Update the ThreatInsight Configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ThreatInsightAPI.UpdateConfiguration(apiClient.GetConfig().Context)

			if UpdateConfigurationdata != "" {
				req = req.Data(UpdateConfigurationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateConfigurationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateConfigurationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateConfigurationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateConfigurationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateConfigurationCmd := NewUpdateConfigurationCmd()
	ThreatInsightCmd.AddCommand(UpdateConfigurationCmd)
}

var (
	GetCurrentConfigurationBackupDir string

	GetCurrentConfigurationQuiet bool
)

func NewGetCurrentConfigurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCurrentConfiguration",
		Long: "Retrieve the ThreatInsight Configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ThreatInsightAPI.GetCurrentConfiguration(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCurrentConfigurationQuiet {
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
				dirPath := filepath.Join(GetCurrentConfigurationBackupDir, "threatinsight", "getCurrentConfiguration")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "threatinsight.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCurrentConfigurationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCurrentConfigurationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the ThreatInsight to a file")

	cmd.Flags().StringVarP(&GetCurrentConfigurationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCurrentConfigurationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCurrentConfigurationCmd := NewGetCurrentConfigurationCmd()
	ThreatInsightCmd.AddCommand(GetCurrentConfigurationCmd)
}
