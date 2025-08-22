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

var InlineHookCmd = &cobra.Command{
	Use:  "inlineHook",
	Long: "Manage InlineHookAPI",
}

func init() {
	rootCmd.AddCommand(InlineHookCmd)
}

var (
	CreateInlineHookdata string

	CreateInlineHookRestoreFile string

	CreateInlineHookQuiet bool
)

func NewCreateInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateInlineHookRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateInlineHookRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateInlineHookdata = string(processedData)

				if !CreateInlineHookQuiet {
					fmt.Println("Restoring InlineHook from:", CreateInlineHookRestoreFile)
				}
			}

			req := apiClient.InlineHookAPI.CreateInlineHook(apiClient.GetConfig().Context)

			if CreateInlineHookdata != "" {
				req = req.Data(CreateInlineHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateInlineHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateInlineHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateInlineHookRestoreFile, "restore-from", "r", "", "Restore InlineHook from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateInlineHookCmd := NewCreateInlineHookCmd()
	InlineHookCmd.AddCommand(CreateInlineHookCmd)
}

var (
	ListInlineHooksBackupDir string

	ListInlineHooksLimit    int32
	ListInlineHooksPage     string
	ListInlineHooksFetchAll bool

	ListInlineHooksQuiet bool
)

func NewListInlineHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Inline Hooks",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ListInlineHooks(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListInlineHooksQuiet {
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
					if !ListInlineHooksQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListInlineHooksFetchAll || len(items) == 0 {
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

			if ListInlineHooksFetchAll && pageCount > 1 && !ListInlineHooksQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListInlineHooksBackupDir, "inlinehook", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListInlineHooksQuiet {
					fmt.Printf("Backing up InlineHooks to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListInlineHooksQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListInlineHooksQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListInlineHooksQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListInlineHooksQuiet {
					fmt.Printf("Successfully backed up %d/%d InlineHooks\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListInlineHooksQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListInlineHooksLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListInlineHooksPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListInlineHooksFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple InlineHooks to a directory")

	cmd.Flags().StringVarP(&ListInlineHooksBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListInlineHooksQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListInlineHooksCmd := NewListInlineHooksCmd()
	InlineHookCmd.AddCommand(ListInlineHooksCmd)
}

var (
	GetInlineHookinlineHookId string

	GetInlineHookBackupDir string

	GetInlineHookQuiet bool
)

func NewGetInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.GetInlineHook(apiClient.GetConfig().Context, GetInlineHookinlineHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetInlineHookQuiet {
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
				dirPath := filepath.Join(GetInlineHookBackupDir, "inlinehook", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetInlineHookinlineHookId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetInlineHookQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the InlineHook to a file")

	cmd.Flags().StringVarP(&GetInlineHookBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetInlineHookCmd := NewGetInlineHookCmd()
	InlineHookCmd.AddCommand(GetInlineHookCmd)
}

var (
	ReplaceInlineHookinlineHookId string

	ReplaceInlineHookdata string

	ReplaceInlineHookQuiet bool
)

func NewReplaceInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ReplaceInlineHook(apiClient.GetConfig().Context, ReplaceInlineHookinlineHookId)

			if ReplaceInlineHookdata != "" {
				req = req.Data(ReplaceInlineHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceInlineHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().StringVarP(&ReplaceInlineHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceInlineHookCmd := NewReplaceInlineHookCmd()
	InlineHookCmd.AddCommand(ReplaceInlineHookCmd)
}

var (
	DeleteInlineHookinlineHookId string

	DeleteInlineHookQuiet bool
)

func NewDeleteInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.DeleteInlineHook(apiClient.GetConfig().Context, DeleteInlineHookinlineHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteInlineHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().BoolVarP(&DeleteInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteInlineHookCmd := NewDeleteInlineHookCmd()
	InlineHookCmd.AddCommand(DeleteInlineHookCmd)
}

var (
	ExecuteInlineHookinlineHookId string

	ExecuteInlineHookdata string

	ExecuteInlineHookQuiet bool
)

func NewExecuteInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "execute",
		Long: "Execute an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ExecuteInlineHook(apiClient.GetConfig().Context, ExecuteInlineHookinlineHookId)

			if ExecuteInlineHookdata != "" {
				req = req.Data(ExecuteInlineHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ExecuteInlineHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ExecuteInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ExecuteInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().StringVarP(&ExecuteInlineHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ExecuteInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ExecuteInlineHookCmd := NewExecuteInlineHookCmd()
	InlineHookCmd.AddCommand(ExecuteInlineHookCmd)
}

var (
	ActivateInlineHookinlineHookId string

	ActivateInlineHookQuiet bool
)

func NewActivateInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ActivateInlineHook(apiClient.GetConfig().Context, ActivateInlineHookinlineHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateInlineHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().BoolVarP(&ActivateInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateInlineHookCmd := NewActivateInlineHookCmd()
	InlineHookCmd.AddCommand(ActivateInlineHookCmd)
}

var (
	DeactivateInlineHookinlineHookId string

	DeactivateInlineHookQuiet bool
)

func NewDeactivateInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate an Inline Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.DeactivateInlineHook(apiClient.GetConfig().Context, DeactivateInlineHookinlineHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateInlineHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateInlineHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().BoolVarP(&DeactivateInlineHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateInlineHookCmd := NewDeactivateInlineHookCmd()
	InlineHookCmd.AddCommand(DeactivateInlineHookCmd)
}
