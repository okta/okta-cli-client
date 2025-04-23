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

var ApiTokenCmd = &cobra.Command{
	Use:  "apiToken",
	Long: "Manage ApiTokenAPI",
}

func init() {
	rootCmd.AddCommand(ApiTokenCmd)
}

var (
	ListApiTokensBackupDir string

	ListApiTokensLimit    int32
	ListApiTokensPage     string
	ListApiTokensFetchAll bool

	ListApiTokensQuiet bool
)

func NewListApiTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all API Token Metadata",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.ListApiTokens(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApiTokensQuiet {
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
					if !ListApiTokensQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApiTokensFetchAll || len(items) == 0 {
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

			if ListApiTokensFetchAll && pageCount > 1 && !ListApiTokensQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApiTokensBackupDir, "apitoken", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApiTokensQuiet {
					fmt.Printf("Backing up ApiTokens to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApiTokensQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApiTokensQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApiTokensQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApiTokensQuiet {
					fmt.Printf("Successfully backed up %d/%d ApiTokens\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApiTokensQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListApiTokensLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApiTokensPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApiTokensFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApiTokens to a directory")

	cmd.Flags().StringVarP(&ListApiTokensBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApiTokensQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApiTokensCmd := NewListApiTokensCmd()
	ApiTokenCmd.AddCommand(ListApiTokensCmd)
}

var RevokeCurrentApiTokenQuiet bool

func NewRevokeCurrentApiTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeCurrent",
		Long: "Revoke the Current API Token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.RevokeCurrentApiToken(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeCurrentApiTokenQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeCurrentApiTokenQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&RevokeCurrentApiTokenQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeCurrentApiTokenCmd := NewRevokeCurrentApiTokenCmd()
	ApiTokenCmd.AddCommand(RevokeCurrentApiTokenCmd)
}

var (
	GetApiTokenapiTokenId string

	GetApiTokenBackupDir string

	GetApiTokenQuiet bool
)

func NewGetApiTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an API Token's Metadata",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.GetApiToken(apiClient.GetConfig().Context, GetApiTokenapiTokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApiTokenQuiet {
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
				dirPath := filepath.Join(GetApiTokenBackupDir, "apitoken", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApiTokenapiTokenId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApiTokenQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApiTokenQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApiTokenapiTokenId, "apiTokenId", "", "", "")
	cmd.MarkFlagRequired("apiTokenId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApiToken to a file")

	cmd.Flags().StringVarP(&GetApiTokenBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApiTokenQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApiTokenCmd := NewGetApiTokenCmd()
	ApiTokenCmd.AddCommand(GetApiTokenCmd)
}

var (
	RevokeApiTokenapiTokenId string

	RevokeApiTokenQuiet bool
)

func NewRevokeApiTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revoke",
		Long: "Revoke an API Token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiTokenAPI.RevokeApiToken(apiClient.GetConfig().Context, RevokeApiTokenapiTokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeApiTokenQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeApiTokenQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeApiTokenapiTokenId, "apiTokenId", "", "", "")
	cmd.MarkFlagRequired("apiTokenId")

	cmd.Flags().BoolVarP(&RevokeApiTokenQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeApiTokenCmd := NewRevokeApiTokenCmd()
	ApiTokenCmd.AddCommand(RevokeApiTokenCmd)
}
