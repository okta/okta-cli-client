package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var SessionCmd = &cobra.Command{
	Use:  "session",
	Long: "Manage SessionAPI",
}

func init() {
	rootCmd.AddCommand(SessionCmd)
}

var (
	CreateSessiondata string

	CreateSessionRestoreFile string

	CreateSessionQuiet bool
)

func NewCreateSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Session with session token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateSessionRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateSessionRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateSessiondata = string(processedData)

				if !CreateSessionQuiet {
					fmt.Println("Restoring Session from:", CreateSessionRestoreFile)
				}
			}

			req := apiClient.SessionAPI.CreateSession(apiClient.GetConfig().Context)

			if CreateSessiondata != "" {
				req = req.Data(CreateSessiondata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateSessiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateSessionRestoreFile, "restore-from", "r", "", "Restore Session from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateSessionCmd := NewCreateSessionCmd()
	SessionCmd.AddCommand(CreateSessionCmd)
}

var (
	GetCurrentSessionBackupDir string

	GetCurrentSessionQuiet bool
)

func NewGetCurrentSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCurrent",
		Long: "Retrieve the current Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.GetCurrentSession(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCurrentSessionQuiet {
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
				dirPath := filepath.Join(GetCurrentSessionBackupDir, "session", "getCurrent")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "session.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCurrentSessionQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCurrentSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the Session to a file")

	cmd.Flags().StringVarP(&GetCurrentSessionBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCurrentSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCurrentSessionCmd := NewGetCurrentSessionCmd()
	SessionCmd.AddCommand(GetCurrentSessionCmd)
}

var CloseCurrentSessionQuiet bool

func NewCloseCurrentSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "closeCurrent",
		Long: "Close the current Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.CloseCurrentSession(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CloseCurrentSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CloseCurrentSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&CloseCurrentSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CloseCurrentSessionCmd := NewCloseCurrentSessionCmd()
	SessionCmd.AddCommand(CloseCurrentSessionCmd)
}

var RefreshCurrentSessionQuiet bool

func NewRefreshCurrentSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "refreshCurrent",
		Long: "Refresh the current Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.RefreshCurrentSession(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RefreshCurrentSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RefreshCurrentSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&RefreshCurrentSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RefreshCurrentSessionCmd := NewRefreshCurrentSessionCmd()
	SessionCmd.AddCommand(RefreshCurrentSessionCmd)
}

var (
	GetSessionsessionId string

	GetSessionBackupDir string

	GetSessionQuiet bool
)

func NewGetSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.GetSession(apiClient.GetConfig().Context, GetSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetSessionQuiet {
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
				dirPath := filepath.Join(GetSessionBackupDir, "session", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetSessionsessionId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetSessionQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Session to a file")

	cmd.Flags().StringVarP(&GetSessionBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetSessionCmd := NewGetSessionCmd()
	SessionCmd.AddCommand(GetSessionCmd)
}

var (
	RevokeSessionsessionId string

	RevokeSessionQuiet bool
)

func NewRevokeSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revoke",
		Long: "Revoke a Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.RevokeSession(apiClient.GetConfig().Context, RevokeSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().BoolVarP(&RevokeSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeSessionCmd := NewRevokeSessionCmd()
	SessionCmd.AddCommand(RevokeSessionCmd)
}

var (
	RefreshSessionsessionId string

	RefreshSessionQuiet bool
)

func NewRefreshSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "refresh",
		Long: "Refresh a Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.RefreshSession(apiClient.GetConfig().Context, RefreshSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RefreshSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RefreshSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RefreshSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().BoolVarP(&RefreshSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RefreshSessionCmd := NewRefreshSessionCmd()
	SessionCmd.AddCommand(RefreshSessionCmd)
}
