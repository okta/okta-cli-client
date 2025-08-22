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

var ApplicationTokensCmd = &cobra.Command{
	Use:  "applicationTokens",
	Long: "Manage ApplicationTokensAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationTokensCmd)
}

var (
	ListOAuth2TokensForApplicationappId string

	ListOAuth2TokensForApplicationBackupDir string

	ListOAuth2TokensForApplicationLimit    int32
	ListOAuth2TokensForApplicationPage     string
	ListOAuth2TokensForApplicationFetchAll bool

	ListOAuth2TokensForApplicationQuiet bool
)

func NewListOAuth2TokensForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listOAuth2TokensForApplication",
		Long: "List all application refresh Tokens",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationTokensAPI.ListOAuth2TokensForApplication(apiClient.GetConfig().Context, ListOAuth2TokensForApplicationappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListOAuth2TokensForApplicationQuiet {
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
					if !ListOAuth2TokensForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListOAuth2TokensForApplicationFetchAll || len(items) == 0 {
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

			if ListOAuth2TokensForApplicationFetchAll && pageCount > 1 && !ListOAuth2TokensForApplicationQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListOAuth2TokensForApplicationBackupDir, "applicationtokens", "listOAuth2TokensForApplication")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListOAuth2TokensForApplicationQuiet {
					fmt.Printf("Backing up ApplicationTokenss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListOAuth2TokensForApplicationQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListOAuth2TokensForApplicationQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListOAuth2TokensForApplicationQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListOAuth2TokensForApplicationQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationTokenss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListOAuth2TokensForApplicationQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListOAuth2TokensForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListOAuth2TokensForApplicationLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListOAuth2TokensForApplicationPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListOAuth2TokensForApplicationFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationTokenss to a directory")

	cmd.Flags().StringVarP(&ListOAuth2TokensForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListOAuth2TokensForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListOAuth2TokensForApplicationCmd := NewListOAuth2TokensForApplicationCmd()
	ApplicationTokensCmd.AddCommand(ListOAuth2TokensForApplicationCmd)
}

var (
	RevokeOAuth2TokensForApplicationappId string

	RevokeOAuth2TokensForApplicationQuiet bool
)

func NewRevokeOAuth2TokensForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeOAuth2TokensForApplication",
		Long: "Revoke all application Tokens",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationTokensAPI.RevokeOAuth2TokensForApplication(apiClient.GetConfig().Context, RevokeOAuth2TokensForApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeOAuth2TokensForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeOAuth2TokensForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeOAuth2TokensForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&RevokeOAuth2TokensForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeOAuth2TokensForApplicationCmd := NewRevokeOAuth2TokensForApplicationCmd()
	ApplicationTokensCmd.AddCommand(RevokeOAuth2TokensForApplicationCmd)
}

var (
	GetOAuth2TokenForApplicationappId string

	GetOAuth2TokenForApplicationtokenId string

	GetOAuth2TokenForApplicationBackupDir string

	GetOAuth2TokenForApplicationQuiet bool
)

func NewGetOAuth2TokenForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOAuth2TokenForApplication",
		Long: "Retrieve an application Token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationTokensAPI.GetOAuth2TokenForApplication(apiClient.GetConfig().Context, GetOAuth2TokenForApplicationappId, GetOAuth2TokenForApplicationtokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOAuth2TokenForApplicationQuiet {
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
				dirPath := filepath.Join(GetOAuth2TokenForApplicationBackupDir, "applicationtokens", "getOAuth2TokenForApplication")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetOAuth2TokenForApplicationappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOAuth2TokenForApplicationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOAuth2TokenForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetOAuth2TokenForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetOAuth2TokenForApplicationtokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationTokens to a file")

	cmd.Flags().StringVarP(&GetOAuth2TokenForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOAuth2TokenForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOAuth2TokenForApplicationCmd := NewGetOAuth2TokenForApplicationCmd()
	ApplicationTokensCmd.AddCommand(GetOAuth2TokenForApplicationCmd)
}

var (
	RevokeOAuth2TokenForApplicationappId string

	RevokeOAuth2TokenForApplicationtokenId string

	RevokeOAuth2TokenForApplicationQuiet bool
)

func NewRevokeOAuth2TokenForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeOAuth2TokenForApplication",
		Long: "Revoke an application Token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationTokensAPI.RevokeOAuth2TokenForApplication(apiClient.GetConfig().Context, RevokeOAuth2TokenForApplicationappId, RevokeOAuth2TokenForApplicationtokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeOAuth2TokenForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeOAuth2TokenForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeOAuth2TokenForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&RevokeOAuth2TokenForApplicationtokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	cmd.Flags().BoolVarP(&RevokeOAuth2TokenForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeOAuth2TokenForApplicationCmd := NewRevokeOAuth2TokenForApplicationCmd()
	ApplicationTokensCmd.AddCommand(RevokeOAuth2TokenForApplicationCmd)
}
