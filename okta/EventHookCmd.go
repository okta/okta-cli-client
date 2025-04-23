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

var EventHookCmd = &cobra.Command{
	Use:  "eventHook",
	Long: "Manage EventHookAPI",
}

func init() {
	rootCmd.AddCommand(EventHookCmd)
}

var (
	CreateEventHookdata string

	CreateEventHookRestoreFile string

	CreateEventHookQuiet bool
)

func NewCreateEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateEventHookRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateEventHookRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateEventHookdata = string(processedData)

				if !CreateEventHookQuiet {
					fmt.Println("Restoring EventHook from:", CreateEventHookRestoreFile)
				}
			}

			req := apiClient.EventHookAPI.CreateEventHook(apiClient.GetConfig().Context)

			if CreateEventHookdata != "" {
				req = req.Data(CreateEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateEventHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateEventHookRestoreFile, "restore-from", "r", "", "Restore EventHook from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateEventHookCmd := NewCreateEventHookCmd()
	EventHookCmd.AddCommand(CreateEventHookCmd)
}

var (
	ListEventHooksBackupDir string

	ListEventHooksLimit    int32
	ListEventHooksPage     string
	ListEventHooksFetchAll bool

	ListEventHooksQuiet bool
)

func NewListEventHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Event Hooks",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.ListEventHooks(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListEventHooksQuiet {
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
					if !ListEventHooksQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListEventHooksFetchAll || len(items) == 0 {
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

			if ListEventHooksFetchAll && pageCount > 1 && !ListEventHooksQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListEventHooksBackupDir, "eventhook", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListEventHooksQuiet {
					fmt.Printf("Backing up EventHooks to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListEventHooksQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListEventHooksQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListEventHooksQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListEventHooksQuiet {
					fmt.Printf("Successfully backed up %d/%d EventHooks\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListEventHooksQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListEventHooksLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListEventHooksPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListEventHooksFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple EventHooks to a directory")

	cmd.Flags().StringVarP(&ListEventHooksBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListEventHooksQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListEventHooksCmd := NewListEventHooksCmd()
	EventHookCmd.AddCommand(ListEventHooksCmd)
}

var (
	GetEventHookeventHookId string

	GetEventHookBackupDir string

	GetEventHookQuiet bool
)

func NewGetEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.GetEventHook(apiClient.GetConfig().Context, GetEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEventHookQuiet {
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
				dirPath := filepath.Join(GetEventHookBackupDir, "eventhook", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEventHookeventHookId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEventHookQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the EventHook to a file")

	cmd.Flags().StringVarP(&GetEventHookBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEventHookCmd := NewGetEventHookCmd()
	EventHookCmd.AddCommand(GetEventHookCmd)
}

var (
	ReplaceEventHookeventHookId string

	ReplaceEventHookdata string

	ReplaceEventHookQuiet bool
)

func NewReplaceEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.ReplaceEventHook(apiClient.GetConfig().Context, ReplaceEventHookeventHookId)

			if ReplaceEventHookdata != "" {
				req = req.Data(ReplaceEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceEventHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().StringVarP(&ReplaceEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceEventHookCmd := NewReplaceEventHookCmd()
	EventHookCmd.AddCommand(ReplaceEventHookCmd)
}

var (
	DeleteEventHookeventHookId string

	DeleteEventHookQuiet bool
)

func NewDeleteEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.DeleteEventHook(apiClient.GetConfig().Context, DeleteEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteEventHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().BoolVarP(&DeleteEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteEventHookCmd := NewDeleteEventHookCmd()
	EventHookCmd.AddCommand(DeleteEventHookCmd)
}

var (
	ActivateEventHookeventHookId string

	ActivateEventHookQuiet bool
)

func NewActivateEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.ActivateEventHook(apiClient.GetConfig().Context, ActivateEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateEventHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().BoolVarP(&ActivateEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateEventHookCmd := NewActivateEventHookCmd()
	EventHookCmd.AddCommand(ActivateEventHookCmd)
}

var (
	DeactivateEventHookeventHookId string

	DeactivateEventHookQuiet bool
)

func NewDeactivateEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.DeactivateEventHook(apiClient.GetConfig().Context, DeactivateEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateEventHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().BoolVarP(&DeactivateEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateEventHookCmd := NewDeactivateEventHookCmd()
	EventHookCmd.AddCommand(DeactivateEventHookCmd)
}

var (
	VerifyEventHookeventHookId string

	VerifyEventHookQuiet bool
)

func NewVerifyEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verify",
		Long: "Verify an Event Hook",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.VerifyEventHook(apiClient.GetConfig().Context, VerifyEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !VerifyEventHookQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !VerifyEventHookQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&VerifyEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().BoolVarP(&VerifyEventHookQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	VerifyEventHookCmd := NewVerifyEventHookCmd()
	EventHookCmd.AddCommand(VerifyEventHookCmd)
}
