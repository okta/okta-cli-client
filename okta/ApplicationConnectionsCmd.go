package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	UpdateDefaultProvisioningConnectionForApplicationQuiet bool
)

func NewUpdateDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateDefaultProvisioningConnectionForApplication",
		Long: "Update the default Provisioning Connection",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.UpdateDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, UpdateDefaultProvisioningConnectionForApplicationappId)

			if UpdateDefaultProvisioningConnectionForApplicationdata != "" {
				req = req.Data(UpdateDefaultProvisioningConnectionForApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateDefaultProvisioningConnectionForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateDefaultProvisioningConnectionForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateDefaultProvisioningConnectionForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateDefaultProvisioningConnectionForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateDefaultProvisioningConnectionForApplicationCmd := NewUpdateDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(UpdateDefaultProvisioningConnectionForApplicationCmd)
}

var (
	GetDefaultProvisioningConnectionForApplicationappId string

	GetDefaultProvisioningConnectionForApplicationBackupDir string

	GetDefaultProvisioningConnectionForApplicationQuiet bool
)

func NewGetDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDefaultProvisioningConnectionForApplication",
		Long: "Retrieve the default Provisioning Connection",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.GetDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, GetDefaultProvisioningConnectionForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetDefaultProvisioningConnectionForApplicationQuiet {
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
				dirPath := filepath.Join(GetDefaultProvisioningConnectionForApplicationBackupDir, "applicationconnections", "getDefaultProvisioningConnectionForApplication")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetDefaultProvisioningConnectionForApplicationappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetDefaultProvisioningConnectionForApplicationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetDefaultProvisioningConnectionForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationConnections to a file")

	cmd.Flags().StringVarP(&GetDefaultProvisioningConnectionForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetDefaultProvisioningConnectionForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetDefaultProvisioningConnectionForApplicationCmd := NewGetDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(GetDefaultProvisioningConnectionForApplicationCmd)
}

var (
	ActivateDefaultProvisioningConnectionForApplicationappId string

	ActivateDefaultProvisioningConnectionForApplicationQuiet bool
)

func NewActivateDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateDefaultProvisioningConnectionForApplication",
		Long: "Activate the default Provisioning Connection",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.ActivateDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, ActivateDefaultProvisioningConnectionForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateDefaultProvisioningConnectionForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateDefaultProvisioningConnectionForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&ActivateDefaultProvisioningConnectionForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateDefaultProvisioningConnectionForApplicationCmd := NewActivateDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(ActivateDefaultProvisioningConnectionForApplicationCmd)
}

var (
	DeactivateDefaultProvisioningConnectionForApplicationappId string

	DeactivateDefaultProvisioningConnectionForApplicationQuiet bool
)

func NewDeactivateDefaultProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateDefaultProvisioningConnectionForApplication",
		Long: "Deactivate the default Provisioning Connection",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.DeactivateDefaultProvisioningConnectionForApplication(apiClient.GetConfig().Context, DeactivateDefaultProvisioningConnectionForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateDefaultProvisioningConnectionForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateDefaultProvisioningConnectionForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateDefaultProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&DeactivateDefaultProvisioningConnectionForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateDefaultProvisioningConnectionForApplicationCmd := NewDeactivateDefaultProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(DeactivateDefaultProvisioningConnectionForApplicationCmd)
}

var (
	VerifyProvisioningConnectionForApplicationappName string

	VerifyProvisioningConnectionForApplicationappId string

	VerifyProvisioningConnectionForApplicationQuiet bool
)

func NewVerifyProvisioningConnectionForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verifyProvisioningConnectionForApplication",
		Long: "Verify the Provisioning Connection",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationConnectionsAPI.VerifyProvisioningConnectionForApplication(apiClient.GetConfig().Context, VerifyProvisioningConnectionForApplicationappName, VerifyProvisioningConnectionForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !VerifyProvisioningConnectionForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !VerifyProvisioningConnectionForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&VerifyProvisioningConnectionForApplicationappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&VerifyProvisioningConnectionForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&VerifyProvisioningConnectionForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	VerifyProvisioningConnectionForApplicationCmd := NewVerifyProvisioningConnectionForApplicationCmd()
	ApplicationConnectionsCmd.AddCommand(VerifyProvisioningConnectionForApplicationCmd)
}
