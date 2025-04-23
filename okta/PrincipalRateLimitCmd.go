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

var PrincipalRateLimitCmd = &cobra.Command{
	Use:  "principalRateLimit",
	Long: "Manage PrincipalRateLimitAPI",
}

func init() {
	rootCmd.AddCommand(PrincipalRateLimitCmd)
}

var (
	CreatePrincipalRateLimitEntitydata string

	CreatePrincipalRateLimitEntityRestoreFile string

	CreatePrincipalRateLimitEntityQuiet bool
)

func NewCreatePrincipalRateLimitEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createEntity",
		Long: "Create a Principal Rate Limit",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreatePrincipalRateLimitEntityRestoreFile != "" {

				jsonData, err := os.ReadFile(CreatePrincipalRateLimitEntityRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreatePrincipalRateLimitEntitydata = string(processedData)

				if !CreatePrincipalRateLimitEntityQuiet {
					fmt.Println("Restoring PrincipalRateLimit from:", CreatePrincipalRateLimitEntityRestoreFile)
				}
			}

			req := apiClient.PrincipalRateLimitAPI.CreatePrincipalRateLimitEntity(apiClient.GetConfig().Context)

			if CreatePrincipalRateLimitEntitydata != "" {
				req = req.Data(CreatePrincipalRateLimitEntitydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreatePrincipalRateLimitEntityQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreatePrincipalRateLimitEntityQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreatePrincipalRateLimitEntitydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreatePrincipalRateLimitEntityRestoreFile, "restore-from", "r", "", "Restore PrincipalRateLimit from a JSON backup file")

	cmd.Flags().BoolVarP(&CreatePrincipalRateLimitEntityQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreatePrincipalRateLimitEntityCmd := NewCreatePrincipalRateLimitEntityCmd()
	PrincipalRateLimitCmd.AddCommand(CreatePrincipalRateLimitEntityCmd)
}

var (
	ListPrincipalRateLimitEntitiesBackupDir string

	ListPrincipalRateLimitEntitiesLimit    int32
	ListPrincipalRateLimitEntitiesPage     string
	ListPrincipalRateLimitEntitiesFetchAll bool

	ListPrincipalRateLimitEntitiesQuiet bool
)

func NewListPrincipalRateLimitEntitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listEntities",
		Long: "List all Principal Rate Limits",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.ListPrincipalRateLimitEntities(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListPrincipalRateLimitEntitiesQuiet {
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
					if !ListPrincipalRateLimitEntitiesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListPrincipalRateLimitEntitiesFetchAll || len(items) == 0 {
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

			if ListPrincipalRateLimitEntitiesFetchAll && pageCount > 1 && !ListPrincipalRateLimitEntitiesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListPrincipalRateLimitEntitiesBackupDir, "principalratelimit", "listEntities")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListPrincipalRateLimitEntitiesQuiet {
					fmt.Printf("Backing up PrincipalRateLimits to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListPrincipalRateLimitEntitiesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListPrincipalRateLimitEntitiesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListPrincipalRateLimitEntitiesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListPrincipalRateLimitEntitiesQuiet {
					fmt.Printf("Successfully backed up %d/%d PrincipalRateLimits\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListPrincipalRateLimitEntitiesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListPrincipalRateLimitEntitiesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListPrincipalRateLimitEntitiesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListPrincipalRateLimitEntitiesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple PrincipalRateLimits to a directory")

	cmd.Flags().StringVarP(&ListPrincipalRateLimitEntitiesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListPrincipalRateLimitEntitiesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListPrincipalRateLimitEntitiesCmd := NewListPrincipalRateLimitEntitiesCmd()
	PrincipalRateLimitCmd.AddCommand(ListPrincipalRateLimitEntitiesCmd)
}

var (
	GetPrincipalRateLimitEntityprincipalRateLimitId string

	GetPrincipalRateLimitEntityBackupDir string

	GetPrincipalRateLimitEntityQuiet bool
)

func NewGetPrincipalRateLimitEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEntity",
		Long: "Retrieve a Principal Rate Limit",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.GetPrincipalRateLimitEntity(apiClient.GetConfig().Context, GetPrincipalRateLimitEntityprincipalRateLimitId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPrincipalRateLimitEntityQuiet {
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
				dirPath := filepath.Join(GetPrincipalRateLimitEntityBackupDir, "principalratelimit", "getEntity")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPrincipalRateLimitEntityprincipalRateLimitId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPrincipalRateLimitEntityQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPrincipalRateLimitEntityQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPrincipalRateLimitEntityprincipalRateLimitId, "principalRateLimitId", "", "", "")
	cmd.MarkFlagRequired("principalRateLimitId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the PrincipalRateLimit to a file")

	cmd.Flags().StringVarP(&GetPrincipalRateLimitEntityBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPrincipalRateLimitEntityQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPrincipalRateLimitEntityCmd := NewGetPrincipalRateLimitEntityCmd()
	PrincipalRateLimitCmd.AddCommand(GetPrincipalRateLimitEntityCmd)
}

var (
	ReplacePrincipalRateLimitEntityprincipalRateLimitId string

	ReplacePrincipalRateLimitEntitydata string

	ReplacePrincipalRateLimitEntityQuiet bool
)

func NewReplacePrincipalRateLimitEntityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceEntity",
		Long: "Replace a Principal Rate Limit",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PrincipalRateLimitAPI.ReplacePrincipalRateLimitEntity(apiClient.GetConfig().Context, ReplacePrincipalRateLimitEntityprincipalRateLimitId)

			if ReplacePrincipalRateLimitEntitydata != "" {
				req = req.Data(ReplacePrincipalRateLimitEntitydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacePrincipalRateLimitEntityQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacePrincipalRateLimitEntityQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePrincipalRateLimitEntityprincipalRateLimitId, "principalRateLimitId", "", "", "")
	cmd.MarkFlagRequired("principalRateLimitId")

	cmd.Flags().StringVarP(&ReplacePrincipalRateLimitEntitydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacePrincipalRateLimitEntityQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacePrincipalRateLimitEntityCmd := NewReplacePrincipalRateLimitEntityCmd()
	PrincipalRateLimitCmd.AddCommand(ReplacePrincipalRateLimitEntityCmd)
}
