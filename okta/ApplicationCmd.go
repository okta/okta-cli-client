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

var ApplicationCmd = &cobra.Command{
	Use:  "application",
	Long: "Manage ApplicationAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationCmd)
}

var (
	CreateApplicationdata string

	CreateApplicationRestoreFile string

	CreateApplicationQuiet bool
)

func NewCreateApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateApplicationRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateApplicationRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateApplicationdata = string(processedData)

				if !CreateApplicationQuiet {
					fmt.Println("Restoring Application from:", CreateApplicationRestoreFile)
				}
			}

			req := apiClient.ApplicationAPI.CreateApplication(apiClient.GetConfig().Context)

			if CreateApplicationdata != "" {
				req = req.Data(CreateApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateApplicationRestoreFile, "restore-from", "r", "", "Restore Application from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateApplicationCmd := NewCreateApplicationCmd()
	ApplicationCmd.AddCommand(CreateApplicationCmd)
}

var (
	ListApplicationsBackupDir string

	ListApplicationsLimit    int32
	ListApplicationsPage     string
	ListApplicationsFetchAll bool

	ListApplicationsQuiet bool
)

func NewListApplicationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Applications",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.ListApplications(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApplicationsQuiet {
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
					if !ListApplicationsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApplicationsFetchAll || len(items) == 0 {
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

			if ListApplicationsFetchAll && pageCount > 1 && !ListApplicationsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApplicationsBackupDir, "application", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApplicationsQuiet {
					fmt.Printf("Backing up Applications to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApplicationsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApplicationsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApplicationsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApplicationsQuiet {
					fmt.Printf("Successfully backed up %d/%d Applications\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApplicationsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListApplicationsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApplicationsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApplicationsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Applications to a directory")

	cmd.Flags().StringVarP(&ListApplicationsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApplicationsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApplicationsCmd := NewListApplicationsCmd()
	ApplicationCmd.AddCommand(ListApplicationsCmd)
}

var (
	GetApplicationappId string

	GetApplicationBackupDir string

	GetApplicationQuiet bool
)

func NewGetApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.GetApplication(apiClient.GetConfig().Context, GetApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApplicationQuiet {
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
				dirPath := filepath.Join(GetApplicationBackupDir, "application", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApplicationappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApplicationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Application to a file")

	cmd.Flags().StringVarP(&GetApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApplicationCmd := NewGetApplicationCmd()
	ApplicationCmd.AddCommand(GetApplicationCmd)
}

var (
	ReplaceApplicationappId string

	ReplaceApplicationdata string

	ReplaceApplicationQuiet bool
)

func NewReplaceApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.ReplaceApplication(apiClient.GetConfig().Context, ReplaceApplicationappId)

			if ReplaceApplicationdata != "" {
				req = req.Data(ReplaceApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&ReplaceApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceApplicationCmd := NewReplaceApplicationCmd()
	ApplicationCmd.AddCommand(ReplaceApplicationCmd)
}

var (
	DeleteApplicationappId string

	DeleteApplicationQuiet bool
)

func NewDeleteApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.DeleteApplication(apiClient.GetConfig().Context, DeleteApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&DeleteApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteApplicationCmd := NewDeleteApplicationCmd()
	ApplicationCmd.AddCommand(DeleteApplicationCmd)
}

var (
	ActivateApplicationappId string

	ActivateApplicationQuiet bool
)

func NewActivateApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.ActivateApplication(apiClient.GetConfig().Context, ActivateApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&ActivateApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateApplicationCmd := NewActivateApplicationCmd()
	ApplicationCmd.AddCommand(ActivateApplicationCmd)
}

var (
	DeactivateApplicationappId string

	DeactivateApplicationQuiet bool
)

func NewDeactivateApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.DeactivateApplication(apiClient.GetConfig().Context, DeactivateApplicationappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&DeactivateApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateApplicationCmd := NewDeactivateApplicationCmd()
	ApplicationCmd.AddCommand(DeactivateApplicationCmd)
}
