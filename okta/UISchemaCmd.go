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

var UISchemaCmd = &cobra.Command{
	Use:  "uISchema",
	Long: "Manage UISchemaAPI",
}

func init() {
	rootCmd.AddCommand(UISchemaCmd)
}

var (
	CreateUISchemadata string

	CreateUISchemaRestoreFile string

	CreateUISchemaQuiet bool
)

func NewCreateUISchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a UI Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateUISchemaRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateUISchemaRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateUISchemadata = string(processedData)

				if !CreateUISchemaQuiet {
					fmt.Println("Restoring UISchema from:", CreateUISchemaRestoreFile)
				}
			}

			req := apiClient.UISchemaAPI.CreateUISchema(apiClient.GetConfig().Context)

			if CreateUISchemadata != "" {
				req = req.Data(CreateUISchemadata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateUISchemaQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateUISchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateUISchemadata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateUISchemaRestoreFile, "restore-from", "r", "", "Restore UISchema from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateUISchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateUISchemaCmd := NewCreateUISchemaCmd()
	UISchemaCmd.AddCommand(CreateUISchemaCmd)
}

var (
	ListUISchemasBackupDir string

	ListUISchemasLimit    int32
	ListUISchemasPage     string
	ListUISchemasFetchAll bool

	ListUISchemasQuiet bool
)

func NewListUISchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all UI Schemas",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.ListUISchemas(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUISchemasQuiet {
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
					if !ListUISchemasQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUISchemasFetchAll || len(items) == 0 {
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

			if ListUISchemasFetchAll && pageCount > 1 && !ListUISchemasQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUISchemasBackupDir, "uischema", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUISchemasQuiet {
					fmt.Printf("Backing up UISchemas to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUISchemasQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUISchemasQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUISchemasQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUISchemasQuiet {
					fmt.Printf("Successfully backed up %d/%d UISchemas\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUISchemasQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListUISchemasLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUISchemasPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUISchemasFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple UISchemas to a directory")

	cmd.Flags().StringVarP(&ListUISchemasBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUISchemasQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUISchemasCmd := NewListUISchemasCmd()
	UISchemaCmd.AddCommand(ListUISchemasCmd)
}

var (
	GetUISchemaid string

	GetUISchemaBackupDir string

	GetUISchemaQuiet bool
)

func NewGetUISchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a UI Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.GetUISchema(apiClient.GetConfig().Context, GetUISchemaid)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetUISchemaQuiet {
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
				dirPath := filepath.Join(GetUISchemaBackupDir, "uischema", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetUISchemaid
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetUISchemaQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetUISchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetUISchemaid, "id", "", "", "")
	cmd.MarkFlagRequired("id")

	cmd.Flags().BoolP("backup", "b", false, "Backup the UISchema to a file")

	cmd.Flags().StringVarP(&GetUISchemaBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetUISchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetUISchemaCmd := NewGetUISchemaCmd()
	UISchemaCmd.AddCommand(GetUISchemaCmd)
}

var (
	ReplaceUISchemasid string

	ReplaceUISchemasdata string

	ReplaceUISchemasQuiet bool
)

func NewReplaceUISchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaces",
		Long: "Replace a UI Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.ReplaceUISchemas(apiClient.GetConfig().Context, ReplaceUISchemasid)

			if ReplaceUISchemasdata != "" {
				req = req.Data(ReplaceUISchemasdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceUISchemasQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceUISchemasQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceUISchemasid, "id", "", "", "")
	cmd.MarkFlagRequired("id")

	cmd.Flags().StringVarP(&ReplaceUISchemasdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceUISchemasQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceUISchemasCmd := NewReplaceUISchemasCmd()
	UISchemaCmd.AddCommand(ReplaceUISchemasCmd)
}

var (
	DeleteUISchemasid string

	DeleteUISchemasQuiet bool
)

func NewDeleteUISchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletes",
		Long: "Delete a UI Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.DeleteUISchemas(apiClient.GetConfig().Context, DeleteUISchemasid)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteUISchemasQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteUISchemasQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteUISchemasid, "id", "", "", "")
	cmd.MarkFlagRequired("id")

	cmd.Flags().BoolVarP(&DeleteUISchemasQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteUISchemasCmd := NewDeleteUISchemasCmd()
	UISchemaCmd.AddCommand(DeleteUISchemasCmd)
}
