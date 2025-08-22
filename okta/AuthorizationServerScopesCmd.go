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

var AuthorizationServerScopesCmd = &cobra.Command{
	Use:  "authorizationServerScopes",
	Long: "Manage AuthorizationServerScopesAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerScopesCmd)
}

var (
	CreateOAuth2ScopeauthServerId string

	CreateOAuth2Scopedata string

	CreateOAuth2ScopeRestoreFile string

	CreateOAuth2ScopeQuiet bool
)

func NewCreateOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createOAuth2Scope",
		Long: "Create a Custom Token Scope",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateOAuth2ScopeRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateOAuth2ScopeRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateOAuth2Scopedata = string(processedData)

				if !CreateOAuth2ScopeQuiet {
					fmt.Println("Restoring AuthorizationServerScopes from:", CreateOAuth2ScopeRestoreFile)
				}
			}

			req := apiClient.AuthorizationServerScopesAPI.CreateOAuth2Scope(apiClient.GetConfig().Context, CreateOAuth2ScopeauthServerId)

			if CreateOAuth2Scopedata != "" {
				req = req.Data(CreateOAuth2Scopedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateOAuth2ScopeQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateOAuth2ScopeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateOAuth2Scopedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateOAuth2ScopeRestoreFile, "restore-from", "r", "", "Restore AuthorizationServerScopes from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateOAuth2ScopeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateOAuth2ScopeCmd := NewCreateOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(CreateOAuth2ScopeCmd)
}

var (
	ListOAuth2ScopesauthServerId string

	ListOAuth2ScopesBackupDir string

	ListOAuth2ScopesLimit    int32
	ListOAuth2ScopesPage     string
	ListOAuth2ScopesFetchAll bool

	ListOAuth2ScopesQuiet bool
)

func NewListOAuth2ScopesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listOAuth2Scopes",
		Long: "List all Custom Token Scopes",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.ListOAuth2Scopes(apiClient.GetConfig().Context, ListOAuth2ScopesauthServerId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListOAuth2ScopesQuiet {
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
					if !ListOAuth2ScopesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListOAuth2ScopesFetchAll || len(items) == 0 {
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

			if ListOAuth2ScopesFetchAll && pageCount > 1 && !ListOAuth2ScopesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListOAuth2ScopesBackupDir, "authorizationserverscopes", "listOAuth2Scopes")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListOAuth2ScopesQuiet {
					fmt.Printf("Backing up AuthorizationServerScopess to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListOAuth2ScopesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListOAuth2ScopesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListOAuth2ScopesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListOAuth2ScopesQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerScopess\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListOAuth2ScopesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListOAuth2ScopesauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().Int32VarP(&ListOAuth2ScopesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListOAuth2ScopesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListOAuth2ScopesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerScopess to a directory")

	cmd.Flags().StringVarP(&ListOAuth2ScopesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListOAuth2ScopesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListOAuth2ScopesCmd := NewListOAuth2ScopesCmd()
	AuthorizationServerScopesCmd.AddCommand(ListOAuth2ScopesCmd)
}

var (
	GetOAuth2ScopeauthServerId string

	GetOAuth2ScopescopeId string

	GetOAuth2ScopeBackupDir string

	GetOAuth2ScopeQuiet bool
)

func NewGetOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOAuth2Scope",
		Long: "Retrieve a Custom Token Scope",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.GetOAuth2Scope(apiClient.GetConfig().Context, GetOAuth2ScopeauthServerId, GetOAuth2ScopescopeId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOAuth2ScopeQuiet {
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
				dirPath := filepath.Join(GetOAuth2ScopeBackupDir, "authorizationserverscopes", "getOAuth2Scope")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetOAuth2ScopeauthServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOAuth2ScopeQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOAuth2ScopeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetOAuth2ScopescopeId, "scopeId", "", "", "")
	cmd.MarkFlagRequired("scopeId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AuthorizationServerScopes to a file")

	cmd.Flags().StringVarP(&GetOAuth2ScopeBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOAuth2ScopeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOAuth2ScopeCmd := NewGetOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(GetOAuth2ScopeCmd)
}

var (
	ReplaceOAuth2ScopeauthServerId string

	ReplaceOAuth2ScopescopeId string

	ReplaceOAuth2Scopedata string

	ReplaceOAuth2ScopeQuiet bool
)

func NewReplaceOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceOAuth2Scope",
		Long: "Replace a Custom Token Scope",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.ReplaceOAuth2Scope(apiClient.GetConfig().Context, ReplaceOAuth2ScopeauthServerId, ReplaceOAuth2ScopescopeId)

			if ReplaceOAuth2Scopedata != "" {
				req = req.Data(ReplaceOAuth2Scopedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceOAuth2ScopeQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceOAuth2ScopeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceOAuth2ScopescopeId, "scopeId", "", "", "")
	cmd.MarkFlagRequired("scopeId")

	cmd.Flags().StringVarP(&ReplaceOAuth2Scopedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceOAuth2ScopeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceOAuth2ScopeCmd := NewReplaceOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(ReplaceOAuth2ScopeCmd)
}

var (
	DeleteOAuth2ScopeauthServerId string

	DeleteOAuth2ScopescopeId string

	DeleteOAuth2ScopeQuiet bool
)

func NewDeleteOAuth2ScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteOAuth2Scope",
		Long: "Delete a Custom Token Scope",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerScopesAPI.DeleteOAuth2Scope(apiClient.GetConfig().Context, DeleteOAuth2ScopeauthServerId, DeleteOAuth2ScopescopeId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteOAuth2ScopeQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteOAuth2ScopeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteOAuth2ScopeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteOAuth2ScopescopeId, "scopeId", "", "", "")
	cmd.MarkFlagRequired("scopeId")

	cmd.Flags().BoolVarP(&DeleteOAuth2ScopeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteOAuth2ScopeCmd := NewDeleteOAuth2ScopeCmd()
	AuthorizationServerScopesCmd.AddCommand(DeleteOAuth2ScopeCmd)
}
