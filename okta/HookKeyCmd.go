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

var HookKeyCmd = &cobra.Command{
	Use:  "hookKey",
	Long: "Manage HookKeyAPI",
}

func init() {
	rootCmd.AddCommand(HookKeyCmd)
}

var (
	CreateHookKeydata string

	CreateHookKeyRestoreFile string

	CreateHookKeyQuiet bool
)

func NewCreateHookKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateHookKeyRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateHookKeyRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateHookKeydata = string(processedData)

				if !CreateHookKeyQuiet {
					fmt.Println("Restoring HookKey from:", CreateHookKeyRestoreFile)
				}
			}

			req := apiClient.HookKeyAPI.CreateHookKey(apiClient.GetConfig().Context)

			if CreateHookKeydata != "" {
				req = req.Data(CreateHookKeydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateHookKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateHookKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateHookKeydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateHookKeyRestoreFile, "restore-from", "r", "", "Restore HookKey from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateHookKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateHookKeyCmd := NewCreateHookKeyCmd()
	HookKeyCmd.AddCommand(CreateHookKeyCmd)
}

var (
	ListHookKeysBackupDir string

	ListHookKeysLimit    int32
	ListHookKeysPage     string
	ListHookKeysFetchAll bool

	ListHookKeysQuiet bool
)

func NewListHookKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.HookKeyAPI.ListHookKeys(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListHookKeysQuiet {
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
					if !ListHookKeysQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListHookKeysFetchAll || len(items) == 0 {
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

			if ListHookKeysFetchAll && pageCount > 1 && !ListHookKeysQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListHookKeysBackupDir, "hookkey", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListHookKeysQuiet {
					fmt.Printf("Backing up HookKeys to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListHookKeysQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListHookKeysQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListHookKeysQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListHookKeysQuiet {
					fmt.Printf("Successfully backed up %d/%d HookKeys\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListHookKeysQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListHookKeysLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListHookKeysPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListHookKeysFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple HookKeys to a directory")

	cmd.Flags().StringVarP(&ListHookKeysBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListHookKeysQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListHookKeysCmd := NewListHookKeysCmd()
	HookKeyCmd.AddCommand(ListHookKeysCmd)
}

var (
	GetPublicKeypublicKeyId string

	GetPublicKeyBackupDir string

	GetPublicKeyQuiet bool
)

func NewGetPublicKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPublicKey",
		Long: "Retrieve a public key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.HookKeyAPI.GetPublicKey(apiClient.GetConfig().Context, GetPublicKeypublicKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPublicKeyQuiet {
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
				dirPath := filepath.Join(GetPublicKeyBackupDir, "hookkey", "getPublicKey")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPublicKeypublicKeyId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPublicKeyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPublicKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPublicKeypublicKeyId, "publicKeyId", "", "", "")
	cmd.MarkFlagRequired("publicKeyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the HookKey to a file")

	cmd.Flags().StringVarP(&GetPublicKeyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPublicKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPublicKeyCmd := NewGetPublicKeyCmd()
	HookKeyCmd.AddCommand(GetPublicKeyCmd)
}

var (
	GetHookKeyhookKeyId string

	GetHookKeyBackupDir string

	GetHookKeyQuiet bool
)

func NewGetHookKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.HookKeyAPI.GetHookKey(apiClient.GetConfig().Context, GetHookKeyhookKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetHookKeyQuiet {
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
				dirPath := filepath.Join(GetHookKeyBackupDir, "hookkey", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetHookKeyhookKeyId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetHookKeyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetHookKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetHookKeyhookKeyId, "hookKeyId", "", "", "")
	cmd.MarkFlagRequired("hookKeyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the HookKey to a file")

	cmd.Flags().StringVarP(&GetHookKeyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetHookKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetHookKeyCmd := NewGetHookKeyCmd()
	HookKeyCmd.AddCommand(GetHookKeyCmd)
}

var (
	ReplaceHookKeyhookKeyId string

	ReplaceHookKeydata string

	ReplaceHookKeyQuiet bool
)

func NewReplaceHookKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.HookKeyAPI.ReplaceHookKey(apiClient.GetConfig().Context, ReplaceHookKeyhookKeyId)

			if ReplaceHookKeydata != "" {
				req = req.Data(ReplaceHookKeydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceHookKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceHookKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceHookKeyhookKeyId, "hookKeyId", "", "", "")
	cmd.MarkFlagRequired("hookKeyId")

	cmd.Flags().StringVarP(&ReplaceHookKeydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceHookKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceHookKeyCmd := NewReplaceHookKeyCmd()
	HookKeyCmd.AddCommand(ReplaceHookKeyCmd)
}

var (
	DeleteHookKeyhookKeyId string

	DeleteHookKeyQuiet bool
)

func NewDeleteHookKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.HookKeyAPI.DeleteHookKey(apiClient.GetConfig().Context, DeleteHookKeyhookKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteHookKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteHookKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteHookKeyhookKeyId, "hookKeyId", "", "", "")
	cmd.MarkFlagRequired("hookKeyId")

	cmd.Flags().BoolVarP(&DeleteHookKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteHookKeyCmd := NewDeleteHookKeyCmd()
	HookKeyCmd.AddCommand(DeleteHookKeyCmd)
}
