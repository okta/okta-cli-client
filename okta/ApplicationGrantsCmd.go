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

var ApplicationGrantsCmd = &cobra.Command{
	Use:  "applicationGrants",
	Long: "Manage ApplicationGrantsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationGrantsCmd)
}

var (
	GrantConsentToScopeappId string

	GrantConsentToScopedata string

	GrantConsentToScopeQuiet bool
)

func NewGrantConsentToScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "grantConsentToScope",
		Long: "Grant consent to scope",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.GrantConsentToScope(apiClient.GetConfig().Context, GrantConsentToScopeappId)

			if GrantConsentToScopedata != "" {
				req = req.Data(GrantConsentToScopedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GrantConsentToScopeQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GrantConsentToScopeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GrantConsentToScopeappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GrantConsentToScopedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&GrantConsentToScopeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GrantConsentToScopeCmd := NewGrantConsentToScopeCmd()
	ApplicationGrantsCmd.AddCommand(GrantConsentToScopeCmd)
}

var (
	ListScopeConsentGrantsappId string

	ListScopeConsentGrantsBackupDir string

	ListScopeConsentGrantsLimit    int32
	ListScopeConsentGrantsPage     string
	ListScopeConsentGrantsFetchAll bool

	ListScopeConsentGrantsQuiet bool
)

func NewListScopeConsentGrantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listScopeConsentGrants",
		Long: "List all app Grants",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.ListScopeConsentGrants(apiClient.GetConfig().Context, ListScopeConsentGrantsappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListScopeConsentGrantsQuiet {
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
					if !ListScopeConsentGrantsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListScopeConsentGrantsFetchAll || len(items) == 0 {
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

			if ListScopeConsentGrantsFetchAll && pageCount > 1 && !ListScopeConsentGrantsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListScopeConsentGrantsBackupDir, "applicationgrants", "listScopeConsentGrants")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListScopeConsentGrantsQuiet {
					fmt.Printf("Backing up ApplicationGrantss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListScopeConsentGrantsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListScopeConsentGrantsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListScopeConsentGrantsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListScopeConsentGrantsQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationGrantss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListScopeConsentGrantsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListScopeConsentGrantsappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListScopeConsentGrantsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListScopeConsentGrantsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListScopeConsentGrantsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationGrantss to a directory")

	cmd.Flags().StringVarP(&ListScopeConsentGrantsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListScopeConsentGrantsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListScopeConsentGrantsCmd := NewListScopeConsentGrantsCmd()
	ApplicationGrantsCmd.AddCommand(ListScopeConsentGrantsCmd)
}

var (
	GetScopeConsentGrantappId string

	GetScopeConsentGrantgrantId string

	GetScopeConsentGrantBackupDir string

	GetScopeConsentGrantQuiet bool
)

func NewGetScopeConsentGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getScopeConsentGrant",
		Long: "Retrieve an app Grant",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.GetScopeConsentGrant(apiClient.GetConfig().Context, GetScopeConsentGrantappId, GetScopeConsentGrantgrantId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetScopeConsentGrantQuiet {
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
				dirPath := filepath.Join(GetScopeConsentGrantBackupDir, "applicationgrants", "getScopeConsentGrant")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetScopeConsentGrantappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetScopeConsentGrantQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetScopeConsentGrantQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetScopeConsentGrantappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetScopeConsentGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationGrants to a file")

	cmd.Flags().StringVarP(&GetScopeConsentGrantBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetScopeConsentGrantQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetScopeConsentGrantCmd := NewGetScopeConsentGrantCmd()
	ApplicationGrantsCmd.AddCommand(GetScopeConsentGrantCmd)
}

var (
	RevokeScopeConsentGrantappId string

	RevokeScopeConsentGrantgrantId string

	RevokeScopeConsentGrantQuiet bool
)

func NewRevokeScopeConsentGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeScopeConsentGrant",
		Long: "Revoke an app Grant",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGrantsAPI.RevokeScopeConsentGrant(apiClient.GetConfig().Context, RevokeScopeConsentGrantappId, RevokeScopeConsentGrantgrantId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeScopeConsentGrantQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeScopeConsentGrantQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeScopeConsentGrantappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&RevokeScopeConsentGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	cmd.Flags().BoolVarP(&RevokeScopeConsentGrantQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeScopeConsentGrantCmd := NewRevokeScopeConsentGrantCmd()
	ApplicationGrantsCmd.AddCommand(RevokeScopeConsentGrantCmd)
}
