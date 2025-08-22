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

var AuthorizationServerKeysCmd = &cobra.Command{
	Use:  "authorizationServerKeys",
	Long: "Manage AuthorizationServerKeysAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerKeysCmd)
}

var (
	ListAuthorizationServerKeysauthServerId string

	ListAuthorizationServerKeysBackupDir string

	ListAuthorizationServerKeysLimit    int32
	ListAuthorizationServerKeysPage     string
	ListAuthorizationServerKeysFetchAll bool

	ListAuthorizationServerKeysQuiet bool
)

func NewListAuthorizationServerKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List all Credential Keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerKeysAPI.ListAuthorizationServerKeys(apiClient.GetConfig().Context, ListAuthorizationServerKeysauthServerId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAuthorizationServerKeysQuiet {
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
					if !ListAuthorizationServerKeysQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAuthorizationServerKeysFetchAll || len(items) == 0 {
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

			if ListAuthorizationServerKeysFetchAll && pageCount > 1 && !ListAuthorizationServerKeysQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAuthorizationServerKeysBackupDir, "authorizationserverkeys", "list")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAuthorizationServerKeysQuiet {
					fmt.Printf("Backing up AuthorizationServerKeyss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAuthorizationServerKeysQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAuthorizationServerKeysQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAuthorizationServerKeysQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAuthorizationServerKeysQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerKeyss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAuthorizationServerKeysQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAuthorizationServerKeysauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().Int32VarP(&ListAuthorizationServerKeysLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAuthorizationServerKeysPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAuthorizationServerKeysFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerKeyss to a directory")

	cmd.Flags().StringVarP(&ListAuthorizationServerKeysBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAuthorizationServerKeysQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAuthorizationServerKeysCmd := NewListAuthorizationServerKeysCmd()
	AuthorizationServerKeysCmd.AddCommand(ListAuthorizationServerKeysCmd)
}

var (
	RotateAuthorizationServerKeysauthServerId string

	RotateAuthorizationServerKeysdata string

	RotateAuthorizationServerKeysQuiet bool
)

func NewRotateAuthorizationServerKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "rotate",
		Long: "Rotate all Credential Keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerKeysAPI.RotateAuthorizationServerKeys(apiClient.GetConfig().Context, RotateAuthorizationServerKeysauthServerId)

			if RotateAuthorizationServerKeysdata != "" {
				req = req.Data(RotateAuthorizationServerKeysdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RotateAuthorizationServerKeysQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RotateAuthorizationServerKeysQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RotateAuthorizationServerKeysauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&RotateAuthorizationServerKeysdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&RotateAuthorizationServerKeysQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RotateAuthorizationServerKeysCmd := NewRotateAuthorizationServerKeysCmd()
	AuthorizationServerKeysCmd.AddCommand(RotateAuthorizationServerKeysCmd)
}
