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

var RiskProviderCmd = &cobra.Command{
	Use:  "riskProvider",
	Long: "Manage RiskProviderAPI",
}

func init() {
	rootCmd.AddCommand(RiskProviderCmd)
}

var (
	CreateRiskProviderdata string

	CreateRiskProviderRestoreFile string

	CreateRiskProviderQuiet bool
)

func NewCreateRiskProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Risk Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateRiskProviderRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateRiskProviderRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateRiskProviderdata = string(processedData)

				if !CreateRiskProviderQuiet {
					fmt.Println("Restoring RiskProvider from:", CreateRiskProviderRestoreFile)
				}
			}

			req := apiClient.RiskProviderAPI.CreateRiskProvider(apiClient.GetConfig().Context)

			if CreateRiskProviderdata != "" {
				req = req.Data(CreateRiskProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateRiskProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateRiskProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateRiskProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateRiskProviderRestoreFile, "restore-from", "r", "", "Restore RiskProvider from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateRiskProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateRiskProviderCmd := NewCreateRiskProviderCmd()
	RiskProviderCmd.AddCommand(CreateRiskProviderCmd)
}

var (
	ListRiskProvidersBackupDir string

	ListRiskProvidersLimit    int32
	ListRiskProvidersPage     string
	ListRiskProvidersFetchAll bool

	ListRiskProvidersQuiet bool
)

func NewListRiskProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Risk Providers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RiskProviderAPI.ListRiskProviders(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRiskProvidersQuiet {
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
					if !ListRiskProvidersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRiskProvidersFetchAll || len(items) == 0 {
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

			if ListRiskProvidersFetchAll && pageCount > 1 && !ListRiskProvidersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRiskProvidersBackupDir, "riskprovider", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRiskProvidersQuiet {
					fmt.Printf("Backing up RiskProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRiskProvidersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRiskProvidersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRiskProvidersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRiskProvidersQuiet {
					fmt.Printf("Successfully backed up %d/%d RiskProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRiskProvidersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListRiskProvidersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRiskProvidersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRiskProvidersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RiskProviders to a directory")

	cmd.Flags().StringVarP(&ListRiskProvidersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRiskProvidersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRiskProvidersCmd := NewListRiskProvidersCmd()
	RiskProviderCmd.AddCommand(ListRiskProvidersCmd)
}

var (
	GetRiskProviderriskProviderId string

	GetRiskProviderBackupDir string

	GetRiskProviderQuiet bool
)

func NewGetRiskProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Risk Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RiskProviderAPI.GetRiskProvider(apiClient.GetConfig().Context, GetRiskProviderriskProviderId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRiskProviderQuiet {
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
				dirPath := filepath.Join(GetRiskProviderBackupDir, "riskprovider", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetRiskProviderriskProviderId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRiskProviderQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRiskProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRiskProviderriskProviderId, "riskProviderId", "", "", "")
	cmd.MarkFlagRequired("riskProviderId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the RiskProvider to a file")

	cmd.Flags().StringVarP(&GetRiskProviderBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRiskProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRiskProviderCmd := NewGetRiskProviderCmd()
	RiskProviderCmd.AddCommand(GetRiskProviderCmd)
}

var (
	ReplaceRiskProviderriskProviderId string

	ReplaceRiskProviderdata string

	ReplaceRiskProviderQuiet bool
)

func NewReplaceRiskProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Risk Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RiskProviderAPI.ReplaceRiskProvider(apiClient.GetConfig().Context, ReplaceRiskProviderriskProviderId)

			if ReplaceRiskProviderdata != "" {
				req = req.Data(ReplaceRiskProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRiskProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRiskProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRiskProviderriskProviderId, "riskProviderId", "", "", "")
	cmd.MarkFlagRequired("riskProviderId")

	cmd.Flags().StringVarP(&ReplaceRiskProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRiskProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRiskProviderCmd := NewReplaceRiskProviderCmd()
	RiskProviderCmd.AddCommand(ReplaceRiskProviderCmd)
}

var (
	DeleteRiskProviderriskProviderId string

	DeleteRiskProviderQuiet bool
)

func NewDeleteRiskProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Risk Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RiskProviderAPI.DeleteRiskProvider(apiClient.GetConfig().Context, DeleteRiskProviderriskProviderId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteRiskProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteRiskProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteRiskProviderriskProviderId, "riskProviderId", "", "", "")
	cmd.MarkFlagRequired("riskProviderId")

	cmd.Flags().BoolVarP(&DeleteRiskProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteRiskProviderCmd := NewDeleteRiskProviderCmd()
	RiskProviderCmd.AddCommand(DeleteRiskProviderCmd)
}
