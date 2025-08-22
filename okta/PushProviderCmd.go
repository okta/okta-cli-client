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

var PushProviderCmd = &cobra.Command{
	Use:  "pushProvider",
	Long: "Manage PushProviderAPI",
}

func init() {
	rootCmd.AddCommand(PushProviderCmd)
}

var (
	CreatePushProviderdata string

	CreatePushProviderRestoreFile string

	CreatePushProviderQuiet bool
)

func NewCreatePushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Push Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreatePushProviderRestoreFile != "" {

				jsonData, err := os.ReadFile(CreatePushProviderRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreatePushProviderdata = string(processedData)

				if !CreatePushProviderQuiet {
					fmt.Println("Restoring PushProvider from:", CreatePushProviderRestoreFile)
				}
			}

			req := apiClient.PushProviderAPI.CreatePushProvider(apiClient.GetConfig().Context)

			if CreatePushProviderdata != "" {
				req = req.Data(CreatePushProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreatePushProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreatePushProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreatePushProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreatePushProviderRestoreFile, "restore-from", "r", "", "Restore PushProvider from a JSON backup file")

	cmd.Flags().BoolVarP(&CreatePushProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreatePushProviderCmd := NewCreatePushProviderCmd()
	PushProviderCmd.AddCommand(CreatePushProviderCmd)
}

var (
	ListPushProvidersBackupDir string

	ListPushProvidersLimit    int32
	ListPushProvidersPage     string
	ListPushProvidersFetchAll bool

	ListPushProvidersQuiet bool
)

func NewListPushProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Push Providers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.ListPushProviders(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListPushProvidersQuiet {
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
					if !ListPushProvidersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListPushProvidersFetchAll || len(items) == 0 {
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

			if ListPushProvidersFetchAll && pageCount > 1 && !ListPushProvidersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListPushProvidersBackupDir, "pushprovider", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListPushProvidersQuiet {
					fmt.Printf("Backing up PushProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListPushProvidersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListPushProvidersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListPushProvidersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListPushProvidersQuiet {
					fmt.Printf("Successfully backed up %d/%d PushProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListPushProvidersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListPushProvidersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListPushProvidersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListPushProvidersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple PushProviders to a directory")

	cmd.Flags().StringVarP(&ListPushProvidersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListPushProvidersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListPushProvidersCmd := NewListPushProvidersCmd()
	PushProviderCmd.AddCommand(ListPushProvidersCmd)
}

var (
	GetPushProviderpushProviderId string

	GetPushProviderBackupDir string

	GetPushProviderQuiet bool
)

func NewGetPushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Push Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.GetPushProvider(apiClient.GetConfig().Context, GetPushProviderpushProviderId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPushProviderQuiet {
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
				dirPath := filepath.Join(GetPushProviderBackupDir, "pushprovider", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPushProviderpushProviderId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPushProviderQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPushProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPushProviderpushProviderId, "pushProviderId", "", "", "")
	cmd.MarkFlagRequired("pushProviderId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the PushProvider to a file")

	cmd.Flags().StringVarP(&GetPushProviderBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPushProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPushProviderCmd := NewGetPushProviderCmd()
	PushProviderCmd.AddCommand(GetPushProviderCmd)
}

var (
	ReplacePushProviderpushProviderId string

	ReplacePushProviderdata string

	ReplacePushProviderQuiet bool
)

func NewReplacePushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Push Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.ReplacePushProvider(apiClient.GetConfig().Context, ReplacePushProviderpushProviderId)

			if ReplacePushProviderdata != "" {
				req = req.Data(ReplacePushProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacePushProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacePushProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePushProviderpushProviderId, "pushProviderId", "", "", "")
	cmd.MarkFlagRequired("pushProviderId")

	cmd.Flags().StringVarP(&ReplacePushProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacePushProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacePushProviderCmd := NewReplacePushProviderCmd()
	PushProviderCmd.AddCommand(ReplacePushProviderCmd)
}

var (
	DeletePushProviderpushProviderId string

	DeletePushProviderQuiet bool
)

func NewDeletePushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Push Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.DeletePushProvider(apiClient.GetConfig().Context, DeletePushProviderpushProviderId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeletePushProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeletePushProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeletePushProviderpushProviderId, "pushProviderId", "", "", "")
	cmd.MarkFlagRequired("pushProviderId")

	cmd.Flags().BoolVarP(&DeletePushProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeletePushProviderCmd := NewDeletePushProviderCmd()
	PushProviderCmd.AddCommand(DeletePushProviderCmd)
}
