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

var TrustedOriginCmd = &cobra.Command{
	Use:  "trustedOrigin",
	Long: "Manage TrustedOriginAPI",
}

func init() {
	rootCmd.AddCommand(TrustedOriginCmd)
}

var (
	CreateTrustedOrigindata string

	CreateTrustedOriginRestoreFile string

	CreateTrustedOriginQuiet bool
)

func NewCreateTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Trusted Origin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateTrustedOriginRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateTrustedOriginRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateTrustedOrigindata = string(processedData)

				if !CreateTrustedOriginQuiet {
					fmt.Println("Restoring TrustedOrigin from:", CreateTrustedOriginRestoreFile)
				}
			}

			req := apiClient.TrustedOriginAPI.CreateTrustedOrigin(apiClient.GetConfig().Context)

			if CreateTrustedOrigindata != "" {
				req = req.Data(CreateTrustedOrigindata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateTrustedOriginQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateTrustedOriginQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateTrustedOrigindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateTrustedOriginRestoreFile, "restore-from", "r", "", "Restore TrustedOrigin from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateTrustedOriginQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateTrustedOriginCmd := NewCreateTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(CreateTrustedOriginCmd)
}

var (
	ListTrustedOriginsBackupDir string

	ListTrustedOriginsLimit    int32
	ListTrustedOriginsPage     string
	ListTrustedOriginsFetchAll bool

	ListTrustedOriginsQuiet bool
)

func NewListTrustedOriginsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Trusted Origins",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.ListTrustedOrigins(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListTrustedOriginsQuiet {
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
					if !ListTrustedOriginsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListTrustedOriginsFetchAll || len(items) == 0 {
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

			if ListTrustedOriginsFetchAll && pageCount > 1 && !ListTrustedOriginsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListTrustedOriginsBackupDir, "trustedorigin", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListTrustedOriginsQuiet {
					fmt.Printf("Backing up TrustedOrigins to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListTrustedOriginsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListTrustedOriginsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListTrustedOriginsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListTrustedOriginsQuiet {
					fmt.Printf("Successfully backed up %d/%d TrustedOrigins\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListTrustedOriginsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListTrustedOriginsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListTrustedOriginsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListTrustedOriginsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple TrustedOrigins to a directory")

	cmd.Flags().StringVarP(&ListTrustedOriginsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListTrustedOriginsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListTrustedOriginsCmd := NewListTrustedOriginsCmd()
	TrustedOriginCmd.AddCommand(ListTrustedOriginsCmd)
}

var (
	GetTrustedOrigintrustedOriginId string

	GetTrustedOriginBackupDir string

	GetTrustedOriginQuiet bool
)

func NewGetTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Trusted Origin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.GetTrustedOrigin(apiClient.GetConfig().Context, GetTrustedOrigintrustedOriginId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetTrustedOriginQuiet {
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
				dirPath := filepath.Join(GetTrustedOriginBackupDir, "trustedorigin", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetTrustedOrigintrustedOriginId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetTrustedOriginQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetTrustedOriginQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the TrustedOrigin to a file")

	cmd.Flags().StringVarP(&GetTrustedOriginBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetTrustedOriginQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetTrustedOriginCmd := NewGetTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(GetTrustedOriginCmd)
}

var (
	ReplaceTrustedOrigintrustedOriginId string

	ReplaceTrustedOrigindata string

	ReplaceTrustedOriginQuiet bool
)

func NewReplaceTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Trusted Origin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.ReplaceTrustedOrigin(apiClient.GetConfig().Context, ReplaceTrustedOrigintrustedOriginId)

			if ReplaceTrustedOrigindata != "" {
				req = req.Data(ReplaceTrustedOrigindata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceTrustedOriginQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceTrustedOriginQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	cmd.Flags().StringVarP(&ReplaceTrustedOrigindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceTrustedOriginQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceTrustedOriginCmd := NewReplaceTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(ReplaceTrustedOriginCmd)
}

var (
	DeleteTrustedOrigintrustedOriginId string

	DeleteTrustedOriginQuiet bool
)

func NewDeleteTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Trusted Origin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.DeleteTrustedOrigin(apiClient.GetConfig().Context, DeleteTrustedOrigintrustedOriginId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteTrustedOriginQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteTrustedOriginQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	cmd.Flags().BoolVarP(&DeleteTrustedOriginQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteTrustedOriginCmd := NewDeleteTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(DeleteTrustedOriginCmd)
}

var (
	ActivateTrustedOrigintrustedOriginId string

	ActivateTrustedOriginQuiet bool
)

func NewActivateTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Trusted Origin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.ActivateTrustedOrigin(apiClient.GetConfig().Context, ActivateTrustedOrigintrustedOriginId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateTrustedOriginQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateTrustedOriginQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	cmd.Flags().BoolVarP(&ActivateTrustedOriginQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateTrustedOriginCmd := NewActivateTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(ActivateTrustedOriginCmd)
}

var (
	DeactivateTrustedOrigintrustedOriginId string

	DeactivateTrustedOriginQuiet bool
)

func NewDeactivateTrustedOriginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Trusted Origin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TrustedOriginAPI.DeactivateTrustedOrigin(apiClient.GetConfig().Context, DeactivateTrustedOrigintrustedOriginId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateTrustedOriginQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateTrustedOriginQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateTrustedOrigintrustedOriginId, "trustedOriginId", "", "", "")
	cmd.MarkFlagRequired("trustedOriginId")

	cmd.Flags().BoolVarP(&DeactivateTrustedOriginQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateTrustedOriginCmd := NewDeactivateTrustedOriginCmd()
	TrustedOriginCmd.AddCommand(DeactivateTrustedOriginCmd)
}
