package okta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/okta/okta-cli-client/sdk"
	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var EmailServerCmd = &cobra.Command{
	Use:  "emailServer",
	Long: "Manage EmailServerAPI",
}

func init() {
	rootCmd.AddCommand(EmailServerCmd)
}

var (
	CreateEmailServerdata string

	CreateEmailServerRestoreFile string

	CreateEmailServerQuiet bool
)

func NewCreateEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a custom SMTP server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateEmailServerRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateEmailServerRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateEmailServerdata = string(processedData)

				if !CreateEmailServerQuiet {
					fmt.Println("Restoring EmailServer from:", CreateEmailServerRestoreFile)
				}
			}

			req := apiClient.EmailServerAPI.CreateEmailServer(apiClient.GetConfig().Context)

			if CreateEmailServerdata != "" {
				req = req.Data(CreateEmailServerdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateEmailServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateEmailServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateEmailServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateEmailServerRestoreFile, "restore-from", "r", "", "Restore EmailServer from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateEmailServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateEmailServerCmd := NewCreateEmailServerCmd()
	EmailServerCmd.AddCommand(CreateEmailServerCmd)
}

var (
	ListEmailServersBackupDir string

	ListEmailServersLimit    int32
	ListEmailServersPage     string
	ListEmailServersFetchAll bool

	ListEmailServersQuiet bool
)

func NewListEmailServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all enrolled SMTP servers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.ListEmailServers(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListEmailServersQuiet {
							utils.PrettyPrintByte(d)
						}
					}
					return err
				}

				d, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				var items []map[string]interface{}
				if err := json.Unmarshal(d, &items); err != nil {
					if !ListEmailServersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListEmailServersFetchAll || len(items) == 0 {
					break
				}

				nextURL := ""
				if resp != nil {
					links := resp.Header["Link"]
					for _, link := range links {
						if strings.Contains(link, `rel="next"`) {
							parts := strings.Split(link, ";")
							if len(parts) > 0 {
								urlPart := strings.TrimSpace(parts[0])
								urlPart = strings.TrimPrefix(urlPart, "<")
								urlPart = strings.TrimSuffix(urlPart, ">")
								nextURL = urlPart
								break
							}
						}
					}
				}

				if nextURL == "" {
					break
				}

				nextReq, err := http.NewRequest("GET", nextURL, nil)
				if err != nil {
					break
				}

				token := ""
				cfg := apiClient.GetConfig()
				if cfg != nil {
					apiKeys, ok := cfg.Context.Value(sdk.ContextAPIKeys).(map[string]sdk.APIKey)
					if ok {
						apiKey, exists := apiKeys["API_Token"]
						if exists {
							token = apiKey.Prefix + " " + apiKey.Key
						}
					}
				}

				if token != "" {
					nextReq.Header.Add("Authorization", token)
				}

				nextReq.Header.Add("Accept", "application/json")

				respNext, err := http.DefaultClient.Do(nextReq)
				if err != nil {
					break
				}

				dNext, err := io.ReadAll(respNext.Body)
				respNext.Body.Close()
				if err != nil {
					break
				}

				var nextItems []map[string]interface{}
				if err := json.Unmarshal(dNext, &nextItems); err != nil {
					break
				}

				allItems = append(allItems, nextItems...)
				pageCount++
			}

			if ListEmailServersFetchAll && pageCount > 1 && !ListEmailServersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListEmailServersBackupDir, "emailserver", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListEmailServersQuiet {
					fmt.Printf("Backing up EmailServers to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListEmailServersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListEmailServersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListEmailServersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListEmailServersQuiet {
					fmt.Printf("Successfully backed up %d/%d EmailServers\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListEmailServersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListEmailServersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListEmailServersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListEmailServersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple EmailServers to a directory")

	cmd.Flags().StringVarP(&ListEmailServersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListEmailServersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListEmailServersCmd := NewListEmailServersCmd()
	EmailServerCmd.AddCommand(ListEmailServersCmd)
}

var (
	GetEmailServeremailServerId string

	GetEmailServerBackupDir string

	GetEmailServerQuiet bool
)

func NewGetEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an SMTP Server configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.GetEmailServer(apiClient.GetConfig().Context, GetEmailServeremailServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailServerQuiet {
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
				dirPath := filepath.Join(GetEmailServerBackupDir, "emailserver", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailServeremailServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailServerQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the EmailServer to a file")

	cmd.Flags().StringVarP(&GetEmailServerBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailServerCmd := NewGetEmailServerCmd()
	EmailServerCmd.AddCommand(GetEmailServerCmd)
}

var (
	DeleteEmailServeremailServerId string

	DeleteEmailServerQuiet bool
)

func NewDeleteEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an SMTP Server configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.DeleteEmailServer(apiClient.GetConfig().Context, DeleteEmailServeremailServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteEmailServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteEmailServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	cmd.Flags().BoolVarP(&DeleteEmailServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteEmailServerCmd := NewDeleteEmailServerCmd()
	EmailServerCmd.AddCommand(DeleteEmailServerCmd)
}

var (
	UpdateEmailServeremailServerId string

	UpdateEmailServerdata string

	UpdateEmailServerQuiet bool
)

func NewUpdateEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: "Update an SMTP Server configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.UpdateEmailServer(apiClient.GetConfig().Context, UpdateEmailServeremailServerId)

			if UpdateEmailServerdata != "" {
				req = req.Data(UpdateEmailServerdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateEmailServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateEmailServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	cmd.Flags().StringVarP(&UpdateEmailServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateEmailServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateEmailServerCmd := NewUpdateEmailServerCmd()
	EmailServerCmd.AddCommand(UpdateEmailServerCmd)
}

var (
	TestEmailServeremailServerId string

	TestEmailServerdata string

	TestEmailServerQuiet bool
)

func NewTestEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "test",
		Long: "Test an SMTP Server configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.TestEmailServer(apiClient.GetConfig().Context, TestEmailServeremailServerId)

			if TestEmailServerdata != "" {
				req = req.Data(TestEmailServerdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !TestEmailServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !TestEmailServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&TestEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	cmd.Flags().StringVarP(&TestEmailServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&TestEmailServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	TestEmailServerCmd := NewTestEmailServerCmd()
	EmailServerCmd.AddCommand(TestEmailServerCmd)
}
