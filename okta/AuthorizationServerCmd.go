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

var AuthorizationServerCmd = &cobra.Command{
	Use:  "authorizationServer",
	Long: "Manage AuthorizationServerAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerCmd)
}

var (
	CreateAuthorizationServerdata string

	CreateAuthorizationServerRestoreFile string

	CreateAuthorizationServerQuiet bool
)

func NewCreateAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateAuthorizationServerRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateAuthorizationServerRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateAuthorizationServerdata = string(processedData)

				if !CreateAuthorizationServerQuiet {
					fmt.Println("Restoring AuthorizationServer from:", CreateAuthorizationServerRestoreFile)
				}
			}

			req := apiClient.AuthorizationServerAPI.CreateAuthorizationServer(apiClient.GetConfig().Context)

			if CreateAuthorizationServerdata != "" {
				req = req.Data(CreateAuthorizationServerdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateAuthorizationServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateAuthorizationServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateAuthorizationServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateAuthorizationServerRestoreFile, "restore-from", "r", "", "Restore AuthorizationServer from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateAuthorizationServerCmd := NewCreateAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(CreateAuthorizationServerCmd)
}

var (
	ListAuthorizationServersBackupDir string

	ListAuthorizationServersLimit    int32
	ListAuthorizationServersPage     string
	ListAuthorizationServersFetchAll bool

	ListAuthorizationServersQuiet bool
)

func NewListAuthorizationServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Authorization Servers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.ListAuthorizationServers(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAuthorizationServersQuiet {
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
					if !ListAuthorizationServersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAuthorizationServersFetchAll || len(items) == 0 {
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

			if ListAuthorizationServersFetchAll && pageCount > 1 && !ListAuthorizationServersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAuthorizationServersBackupDir, "authorizationserver", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAuthorizationServersQuiet {
					fmt.Printf("Backing up AuthorizationServers to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAuthorizationServersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAuthorizationServersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAuthorizationServersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAuthorizationServersQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServers\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAuthorizationServersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListAuthorizationServersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAuthorizationServersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAuthorizationServersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServers to a directory")

	cmd.Flags().StringVarP(&ListAuthorizationServersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAuthorizationServersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAuthorizationServersCmd := NewListAuthorizationServersCmd()
	AuthorizationServerCmd.AddCommand(ListAuthorizationServersCmd)
}

var (
	GetAuthorizationServerauthServerId string

	GetAuthorizationServerBackupDir string

	GetAuthorizationServerQuiet bool
)

func NewGetAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.GetAuthorizationServer(apiClient.GetConfig().Context, GetAuthorizationServerauthServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAuthorizationServerQuiet {
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
				dirPath := filepath.Join(GetAuthorizationServerBackupDir, "authorizationserver", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAuthorizationServerauthServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAuthorizationServerQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAuthorizationServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AuthorizationServer to a file")

	cmd.Flags().StringVarP(&GetAuthorizationServerBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAuthorizationServerCmd := NewGetAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(GetAuthorizationServerCmd)
}

var (
	ReplaceAuthorizationServerauthServerId string

	ReplaceAuthorizationServerdata string

	ReplaceAuthorizationServerQuiet bool
)

func NewReplaceAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.ReplaceAuthorizationServer(apiClient.GetConfig().Context, ReplaceAuthorizationServerauthServerId)

			if ReplaceAuthorizationServerdata != "" {
				req = req.Data(ReplaceAuthorizationServerdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceAuthorizationServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceAuthorizationServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceAuthorizationServerCmd := NewReplaceAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(ReplaceAuthorizationServerCmd)
}

var (
	DeleteAuthorizationServerauthServerId string

	DeleteAuthorizationServerQuiet bool
)

func NewDeleteAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.DeleteAuthorizationServer(apiClient.GetConfig().Context, DeleteAuthorizationServerauthServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteAuthorizationServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteAuthorizationServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().BoolVarP(&DeleteAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteAuthorizationServerCmd := NewDeleteAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(DeleteAuthorizationServerCmd)
}

var (
	ActivateAuthorizationServerauthServerId string

	ActivateAuthorizationServerQuiet bool
)

func NewActivateAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate an Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.ActivateAuthorizationServer(apiClient.GetConfig().Context, ActivateAuthorizationServerauthServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateAuthorizationServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateAuthorizationServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().BoolVarP(&ActivateAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateAuthorizationServerCmd := NewActivateAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(ActivateAuthorizationServerCmd)
}

var (
	DeactivateAuthorizationServerauthServerId string

	DeactivateAuthorizationServerQuiet bool
)

func NewDeactivateAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate an Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAPI.DeactivateAuthorizationServer(apiClient.GetConfig().Context, DeactivateAuthorizationServerauthServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateAuthorizationServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateAuthorizationServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().BoolVarP(&DeactivateAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateAuthorizationServerCmd := NewDeactivateAuthorizationServerCmd()
	AuthorizationServerCmd.AddCommand(DeactivateAuthorizationServerCmd)
}
