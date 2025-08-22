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

var IdentityProviderCmd = &cobra.Command{
	Use:  "identityProvider",
	Long: "Manage IdentityProviderAPI",
}

func init() {
	rootCmd.AddCommand(IdentityProviderCmd)
}

var (
	CreateIdentityProviderdata string

	CreateIdentityProviderRestoreFile string

	CreateIdentityProviderQuiet bool
)

func NewCreateIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateIdentityProviderRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateIdentityProviderRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateIdentityProviderdata = string(processedData)

				if !CreateIdentityProviderQuiet {
					fmt.Println("Restoring IdentityProvider from:", CreateIdentityProviderRestoreFile)
				}
			}

			req := apiClient.IdentityProviderAPI.CreateIdentityProvider(apiClient.GetConfig().Context)

			if CreateIdentityProviderdata != "" {
				req = req.Data(CreateIdentityProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateIdentityProviderRestoreFile, "restore-from", "r", "", "Restore IdentityProvider from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateIdentityProviderCmd := NewCreateIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(CreateIdentityProviderCmd)
}

var (
	ListIdentityProvidersBackupDir string

	ListIdentityProvidersLimit    int32
	ListIdentityProvidersPage     string
	ListIdentityProvidersFetchAll bool

	ListIdentityProvidersQuiet bool
)

func NewListIdentityProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Identity Providers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviders(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListIdentityProvidersQuiet {
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
					if !ListIdentityProvidersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListIdentityProvidersFetchAll || len(items) == 0 {
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

			if ListIdentityProvidersFetchAll && pageCount > 1 && !ListIdentityProvidersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListIdentityProvidersBackupDir, "identityprovider", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListIdentityProvidersQuiet {
					fmt.Printf("Backing up IdentityProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListIdentityProvidersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListIdentityProvidersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListIdentityProvidersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListIdentityProvidersQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentityProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListIdentityProvidersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListIdentityProvidersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListIdentityProvidersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListIdentityProvidersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentityProviders to a directory")

	cmd.Flags().StringVarP(&ListIdentityProvidersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListIdentityProvidersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListIdentityProvidersCmd := NewListIdentityProvidersCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProvidersCmd)
}

var (
	CreateIdentityProviderKeydata string

	CreateIdentityProviderKeyRestoreFile string

	CreateIdentityProviderKeyQuiet bool
)

func NewCreateIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createKey",
		Long: "Create an X.509 Certificate Public Key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateIdentityProviderKeyRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateIdentityProviderKeyRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateIdentityProviderKeydata = string(processedData)

				if !CreateIdentityProviderKeyQuiet {
					fmt.Println("Restoring IdentityProvider from:", CreateIdentityProviderKeyRestoreFile)
				}
			}

			req := apiClient.IdentityProviderAPI.CreateIdentityProviderKey(apiClient.GetConfig().Context)

			if CreateIdentityProviderKeydata != "" {
				req = req.Data(CreateIdentityProviderKeydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateIdentityProviderKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateIdentityProviderKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateIdentityProviderKeydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateIdentityProviderKeyRestoreFile, "restore-from", "r", "", "Restore IdentityProvider from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateIdentityProviderKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateIdentityProviderKeyCmd := NewCreateIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(CreateIdentityProviderKeyCmd)
}

var (
	ListIdentityProviderKeysBackupDir string

	ListIdentityProviderKeysLimit    int32
	ListIdentityProviderKeysPage     string
	ListIdentityProviderKeysFetchAll bool

	ListIdentityProviderKeysQuiet bool
)

func NewListIdentityProviderKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listKeys",
		Long: "List all Credential Keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviderKeys(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListIdentityProviderKeysQuiet {
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
					if !ListIdentityProviderKeysQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListIdentityProviderKeysFetchAll || len(items) == 0 {
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

			if ListIdentityProviderKeysFetchAll && pageCount > 1 && !ListIdentityProviderKeysQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListIdentityProviderKeysBackupDir, "identityprovider", "listKeys")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListIdentityProviderKeysQuiet {
					fmt.Printf("Backing up IdentityProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListIdentityProviderKeysQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListIdentityProviderKeysQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListIdentityProviderKeysQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListIdentityProviderKeysQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentityProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListIdentityProviderKeysQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListIdentityProviderKeysLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListIdentityProviderKeysPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListIdentityProviderKeysFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentityProviders to a directory")

	cmd.Flags().StringVarP(&ListIdentityProviderKeysBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListIdentityProviderKeysQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListIdentityProviderKeysCmd := NewListIdentityProviderKeysCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProviderKeysCmd)
}

var (
	GetIdentityProviderKeyidpKeyId string

	GetIdentityProviderKeyBackupDir string

	GetIdentityProviderKeyQuiet bool
)

func NewGetIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getKey",
		Long: "Retrieve an Credential Key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProviderKey(apiClient.GetConfig().Context, GetIdentityProviderKeyidpKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetIdentityProviderKeyQuiet {
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
				dirPath := filepath.Join(GetIdentityProviderKeyBackupDir, "identityprovider", "getKey")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetIdentityProviderKeyidpKeyId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetIdentityProviderKeyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetIdentityProviderKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetIdentityProviderKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the IdentityProvider to a file")

	cmd.Flags().StringVarP(&GetIdentityProviderKeyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetIdentityProviderKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetIdentityProviderKeyCmd := NewGetIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderKeyCmd)
}

var (
	DeleteIdentityProviderKeyidpKeyId string

	DeleteIdentityProviderKeyQuiet bool
)

func NewDeleteIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteKey",
		Long: "Delete a Signing Credential Key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.DeleteIdentityProviderKey(apiClient.GetConfig().Context, DeleteIdentityProviderKeyidpKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteIdentityProviderKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteIdentityProviderKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteIdentityProviderKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	cmd.Flags().BoolVarP(&DeleteIdentityProviderKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteIdentityProviderKeyCmd := NewDeleteIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(DeleteIdentityProviderKeyCmd)
}

var (
	GetIdentityProvideridpId string

	GetIdentityProviderBackupDir string

	GetIdentityProviderQuiet bool
)

func NewGetIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProvider(apiClient.GetConfig().Context, GetIdentityProvideridpId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetIdentityProviderQuiet {
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
				dirPath := filepath.Join(GetIdentityProviderBackupDir, "identityprovider", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetIdentityProvideridpId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetIdentityProviderQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the IdentityProvider to a file")

	cmd.Flags().StringVarP(&GetIdentityProviderBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetIdentityProviderCmd := NewGetIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderCmd)
}

var (
	ReplaceIdentityProvideridpId string

	ReplaceIdentityProviderdata string

	ReplaceIdentityProviderQuiet bool
)

func NewReplaceIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ReplaceIdentityProvider(apiClient.GetConfig().Context, ReplaceIdentityProvideridpId)

			if ReplaceIdentityProviderdata != "" {
				req = req.Data(ReplaceIdentityProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&ReplaceIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceIdentityProviderCmd := NewReplaceIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(ReplaceIdentityProviderCmd)
}

var (
	DeleteIdentityProvideridpId string

	DeleteIdentityProviderQuiet bool
)

func NewDeleteIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.DeleteIdentityProvider(apiClient.GetConfig().Context, DeleteIdentityProvideridpId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().BoolVarP(&DeleteIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteIdentityProviderCmd := NewDeleteIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(DeleteIdentityProviderCmd)
}

var (
	GenerateCsrForIdentityProvideridpId string

	GenerateCsrForIdentityProviderdata string

	GenerateCsrForIdentityProviderQuiet bool
)

func NewGenerateCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateCsrFor",
		Long: "Generate a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GenerateCsrForIdentityProvider(apiClient.GetConfig().Context, GenerateCsrForIdentityProvideridpId)

			if GenerateCsrForIdentityProviderdata != "" {
				req = req.Data(GenerateCsrForIdentityProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GenerateCsrForIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GenerateCsrForIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GenerateCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GenerateCsrForIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&GenerateCsrForIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GenerateCsrForIdentityProviderCmd := NewGenerateCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(GenerateCsrForIdentityProviderCmd)
}

var (
	ListCsrsForIdentityProvideridpId string

	ListCsrsForIdentityProviderBackupDir string

	ListCsrsForIdentityProviderLimit    int32
	ListCsrsForIdentityProviderPage     string
	ListCsrsForIdentityProviderFetchAll bool

	ListCsrsForIdentityProviderQuiet bool
)

func NewListCsrsForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listCsrsFor",
		Long: "List all Certificate Signing Requests",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListCsrsForIdentityProvider(apiClient.GetConfig().Context, ListCsrsForIdentityProvideridpId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListCsrsForIdentityProviderQuiet {
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
					if !ListCsrsForIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListCsrsForIdentityProviderFetchAll || len(items) == 0 {
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

			if ListCsrsForIdentityProviderFetchAll && pageCount > 1 && !ListCsrsForIdentityProviderQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListCsrsForIdentityProviderBackupDir, "identityprovider", "listCsrsFor")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListCsrsForIdentityProviderQuiet {
					fmt.Printf("Backing up IdentityProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListCsrsForIdentityProviderQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListCsrsForIdentityProviderQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListCsrsForIdentityProviderQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListCsrsForIdentityProviderQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentityProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListCsrsForIdentityProviderQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListCsrsForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().Int32VarP(&ListCsrsForIdentityProviderLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListCsrsForIdentityProviderPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListCsrsForIdentityProviderFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentityProviders to a directory")

	cmd.Flags().StringVarP(&ListCsrsForIdentityProviderBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListCsrsForIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListCsrsForIdentityProviderCmd := NewListCsrsForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(ListCsrsForIdentityProviderCmd)
}

var (
	GetCsrForIdentityProvideridpId string

	GetCsrForIdentityProvideridpCsrId string

	GetCsrForIdentityProviderBackupDir string

	GetCsrForIdentityProviderQuiet bool
)

func NewGetCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCsrFor",
		Long: "Retrieve a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetCsrForIdentityProvider(apiClient.GetConfig().Context, GetCsrForIdentityProvideridpId, GetCsrForIdentityProvideridpCsrId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCsrForIdentityProviderQuiet {
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
				dirPath := filepath.Join(GetCsrForIdentityProviderBackupDir, "identityprovider", "getCsrFor")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCsrForIdentityProvideridpId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCsrForIdentityProviderQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCsrForIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GetCsrForIdentityProvideridpCsrId, "idpCsrId", "", "", "")
	cmd.MarkFlagRequired("idpCsrId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the IdentityProvider to a file")

	cmd.Flags().StringVarP(&GetCsrForIdentityProviderBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCsrForIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCsrForIdentityProviderCmd := NewGetCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(GetCsrForIdentityProviderCmd)
}

var (
	RevokeCsrForIdentityProvideridpId string

	RevokeCsrForIdentityProvideridpCsrId string

	RevokeCsrForIdentityProviderQuiet bool
)

func NewRevokeCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeCsrFor",
		Long: "Revoke a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.RevokeCsrForIdentityProvider(apiClient.GetConfig().Context, RevokeCsrForIdentityProvideridpId, RevokeCsrForIdentityProvideridpCsrId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeCsrForIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeCsrForIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&RevokeCsrForIdentityProvideridpCsrId, "idpCsrId", "", "", "")
	cmd.MarkFlagRequired("idpCsrId")

	cmd.Flags().BoolVarP(&RevokeCsrForIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeCsrForIdentityProviderCmd := NewRevokeCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(RevokeCsrForIdentityProviderCmd)
}

var (
	PublishCsrForIdentityProvideridpId string

	PublishCsrForIdentityProvideridpCsrId string

	PublishCsrForIdentityProviderdata string

	PublishCsrForIdentityProviderQuiet bool
)

func NewPublishCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "publishCsrFor",
		Long: "Publish a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.PublishCsrForIdentityProvider(apiClient.GetConfig().Context, PublishCsrForIdentityProvideridpId, PublishCsrForIdentityProvideridpCsrId)

			if PublishCsrForIdentityProviderdata != "" {
				req = req.Data(PublishCsrForIdentityProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !PublishCsrForIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !PublishCsrForIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&PublishCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&PublishCsrForIdentityProvideridpCsrId, "idpCsrId", "", "", "")
	cmd.MarkFlagRequired("idpCsrId")

	cmd.Flags().StringVarP(&PublishCsrForIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&PublishCsrForIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	PublishCsrForIdentityProviderCmd := NewPublishCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(PublishCsrForIdentityProviderCmd)
}

var (
	ListIdentityProviderSigningKeysidpId string

	ListIdentityProviderSigningKeysBackupDir string

	ListIdentityProviderSigningKeysLimit    int32
	ListIdentityProviderSigningKeysPage     string
	ListIdentityProviderSigningKeysFetchAll bool

	ListIdentityProviderSigningKeysQuiet bool
)

func NewListIdentityProviderSigningKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSigningKeys",
		Long: "List all Signing Credential Keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviderSigningKeys(apiClient.GetConfig().Context, ListIdentityProviderSigningKeysidpId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListIdentityProviderSigningKeysQuiet {
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
					if !ListIdentityProviderSigningKeysQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListIdentityProviderSigningKeysFetchAll || len(items) == 0 {
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

			if ListIdentityProviderSigningKeysFetchAll && pageCount > 1 && !ListIdentityProviderSigningKeysQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListIdentityProviderSigningKeysBackupDir, "identityprovider", "listSigningKeys")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListIdentityProviderSigningKeysQuiet {
					fmt.Printf("Backing up IdentityProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListIdentityProviderSigningKeysQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListIdentityProviderSigningKeysQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListIdentityProviderSigningKeysQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListIdentityProviderSigningKeysQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentityProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListIdentityProviderSigningKeysQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListIdentityProviderSigningKeysidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().Int32VarP(&ListIdentityProviderSigningKeysLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListIdentityProviderSigningKeysPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListIdentityProviderSigningKeysFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentityProviders to a directory")

	cmd.Flags().StringVarP(&ListIdentityProviderSigningKeysBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListIdentityProviderSigningKeysQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListIdentityProviderSigningKeysCmd := NewListIdentityProviderSigningKeysCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProviderSigningKeysCmd)
}

var (
	GenerateIdentityProviderSigningKeyidpId string

	GenerateIdentityProviderSigningKeyQuiet bool
)

func NewGenerateIdentityProviderSigningKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateSigningKey",
		Long: "Generate a new Signing Credential Key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GenerateIdentityProviderSigningKey(apiClient.GetConfig().Context, GenerateIdentityProviderSigningKeyidpId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GenerateIdentityProviderSigningKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GenerateIdentityProviderSigningKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GenerateIdentityProviderSigningKeyidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().BoolVarP(&GenerateIdentityProviderSigningKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GenerateIdentityProviderSigningKeyCmd := NewGenerateIdentityProviderSigningKeyCmd()
	IdentityProviderCmd.AddCommand(GenerateIdentityProviderSigningKeyCmd)
}

var (
	GetIdentityProviderSigningKeyidpId string

	GetIdentityProviderSigningKeyidpKeyId string

	GetIdentityProviderSigningKeyBackupDir string

	GetIdentityProviderSigningKeyQuiet bool
)

func NewGetIdentityProviderSigningKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSigningKey",
		Long: "Retrieve a Signing Credential Key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProviderSigningKey(apiClient.GetConfig().Context, GetIdentityProviderSigningKeyidpId, GetIdentityProviderSigningKeyidpKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetIdentityProviderSigningKeyQuiet {
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
				dirPath := filepath.Join(GetIdentityProviderSigningKeyBackupDir, "identityprovider", "getSigningKey")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetIdentityProviderSigningKeyidpId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetIdentityProviderSigningKeyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetIdentityProviderSigningKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetIdentityProviderSigningKeyidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GetIdentityProviderSigningKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the IdentityProvider to a file")

	cmd.Flags().StringVarP(&GetIdentityProviderSigningKeyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetIdentityProviderSigningKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetIdentityProviderSigningKeyCmd := NewGetIdentityProviderSigningKeyCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderSigningKeyCmd)
}

var (
	CloneIdentityProviderKeyidpId string

	CloneIdentityProviderKeyidpKeyId string

	CloneIdentityProviderKeyQuiet bool
)

func NewCloneIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "cloneKey",
		Long: "Clone a Signing Credential Key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.CloneIdentityProviderKey(apiClient.GetConfig().Context, CloneIdentityProviderKeyidpId, CloneIdentityProviderKeyidpKeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CloneIdentityProviderKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CloneIdentityProviderKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CloneIdentityProviderKeyidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&CloneIdentityProviderKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	cmd.Flags().BoolVarP(&CloneIdentityProviderKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CloneIdentityProviderKeyCmd := NewCloneIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(CloneIdentityProviderKeyCmd)
}

var (
	ActivateIdentityProvideridpId string

	ActivateIdentityProviderQuiet bool
)

func NewActivateIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate an Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ActivateIdentityProvider(apiClient.GetConfig().Context, ActivateIdentityProvideridpId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().BoolVarP(&ActivateIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateIdentityProviderCmd := NewActivateIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(ActivateIdentityProviderCmd)
}

var (
	DeactivateIdentityProvideridpId string

	DeactivateIdentityProviderQuiet bool
)

func NewDeactivateIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate an Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.DeactivateIdentityProvider(apiClient.GetConfig().Context, DeactivateIdentityProvideridpId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().BoolVarP(&DeactivateIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateIdentityProviderCmd := NewDeactivateIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(DeactivateIdentityProviderCmd)
}

var (
	ListIdentityProviderApplicationUsersidpId string

	ListIdentityProviderApplicationUsersBackupDir string

	ListIdentityProviderApplicationUsersLimit    int32
	ListIdentityProviderApplicationUsersPage     string
	ListIdentityProviderApplicationUsersFetchAll bool

	ListIdentityProviderApplicationUsersQuiet bool
)

func NewListIdentityProviderApplicationUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApplicationUsers",
		Long: "List all Users",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviderApplicationUsers(apiClient.GetConfig().Context, ListIdentityProviderApplicationUsersidpId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListIdentityProviderApplicationUsersQuiet {
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
					if !ListIdentityProviderApplicationUsersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListIdentityProviderApplicationUsersFetchAll || len(items) == 0 {
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

			if ListIdentityProviderApplicationUsersFetchAll && pageCount > 1 && !ListIdentityProviderApplicationUsersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListIdentityProviderApplicationUsersBackupDir, "identityprovider", "listApplicationUsers")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListIdentityProviderApplicationUsersQuiet {
					fmt.Printf("Backing up IdentityProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListIdentityProviderApplicationUsersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListIdentityProviderApplicationUsersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListIdentityProviderApplicationUsersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListIdentityProviderApplicationUsersQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentityProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListIdentityProviderApplicationUsersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListIdentityProviderApplicationUsersidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().Int32VarP(&ListIdentityProviderApplicationUsersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListIdentityProviderApplicationUsersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListIdentityProviderApplicationUsersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentityProviders to a directory")

	cmd.Flags().StringVarP(&ListIdentityProviderApplicationUsersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListIdentityProviderApplicationUsersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListIdentityProviderApplicationUsersCmd := NewListIdentityProviderApplicationUsersCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProviderApplicationUsersCmd)
}

var (
	LinkUserToIdentityProvideridpId string

	LinkUserToIdentityProvideruserId string

	LinkUserToIdentityProviderdata string

	LinkUserToIdentityProviderQuiet bool
)

func NewLinkUserToIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "linkUserTo",
		Long: "Link a User to a Social IdP",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.LinkUserToIdentityProvider(apiClient.GetConfig().Context, LinkUserToIdentityProvideridpId, LinkUserToIdentityProvideruserId)

			if LinkUserToIdentityProviderdata != "" {
				req = req.Data(LinkUserToIdentityProviderdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !LinkUserToIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !LinkUserToIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&LinkUserToIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&LinkUserToIdentityProvideruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&LinkUserToIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&LinkUserToIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	LinkUserToIdentityProviderCmd := NewLinkUserToIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(LinkUserToIdentityProviderCmd)
}

var (
	GetIdentityProviderApplicationUseridpId string

	GetIdentityProviderApplicationUseruserId string

	GetIdentityProviderApplicationUserBackupDir string

	GetIdentityProviderApplicationUserQuiet bool
)

func NewGetIdentityProviderApplicationUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationUser",
		Long: "Retrieve a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProviderApplicationUser(apiClient.GetConfig().Context, GetIdentityProviderApplicationUseridpId, GetIdentityProviderApplicationUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetIdentityProviderApplicationUserQuiet {
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
				dirPath := filepath.Join(GetIdentityProviderApplicationUserBackupDir, "identityprovider", "getApplicationUser")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetIdentityProviderApplicationUseridpId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetIdentityProviderApplicationUserQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetIdentityProviderApplicationUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetIdentityProviderApplicationUseridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GetIdentityProviderApplicationUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the IdentityProvider to a file")

	cmd.Flags().StringVarP(&GetIdentityProviderApplicationUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetIdentityProviderApplicationUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetIdentityProviderApplicationUserCmd := NewGetIdentityProviderApplicationUserCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderApplicationUserCmd)
}

var (
	UnlinkUserFromIdentityProvideridpId string

	UnlinkUserFromIdentityProvideruserId string

	UnlinkUserFromIdentityProviderQuiet bool
)

func NewUnlinkUserFromIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unlinkUserFrom",
		Long: "Unlink a User from IdP",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.UnlinkUserFromIdentityProvider(apiClient.GetConfig().Context, UnlinkUserFromIdentityProvideridpId, UnlinkUserFromIdentityProvideruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnlinkUserFromIdentityProviderQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnlinkUserFromIdentityProviderQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnlinkUserFromIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&UnlinkUserFromIdentityProvideruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&UnlinkUserFromIdentityProviderQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnlinkUserFromIdentityProviderCmd := NewUnlinkUserFromIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(UnlinkUserFromIdentityProviderCmd)
}

var (
	ListSocialAuthTokensidpId string

	ListSocialAuthTokensuserId string

	ListSocialAuthTokensBackupDir string

	ListSocialAuthTokensLimit    int32
	ListSocialAuthTokensPage     string
	ListSocialAuthTokensFetchAll bool

	ListSocialAuthTokensQuiet bool
)

func NewListSocialAuthTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSocialAuthTokens",
		Long: "List all Tokens from a OIDC Identity Provider",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListSocialAuthTokens(apiClient.GetConfig().Context, ListSocialAuthTokensidpId, ListSocialAuthTokensuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListSocialAuthTokensQuiet {
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
					if !ListSocialAuthTokensQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListSocialAuthTokensFetchAll || len(items) == 0 {
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

			if ListSocialAuthTokensFetchAll && pageCount > 1 && !ListSocialAuthTokensQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListSocialAuthTokensBackupDir, "identityprovider", "listSocialAuthTokens")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListSocialAuthTokensQuiet {
					fmt.Printf("Backing up IdentityProviders to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListSocialAuthTokensQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListSocialAuthTokensQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListSocialAuthTokensQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListSocialAuthTokensQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentityProviders\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListSocialAuthTokensQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListSocialAuthTokensidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&ListSocialAuthTokensuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListSocialAuthTokensLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListSocialAuthTokensPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListSocialAuthTokensFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentityProviders to a directory")

	cmd.Flags().StringVarP(&ListSocialAuthTokensBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListSocialAuthTokensQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListSocialAuthTokensCmd := NewListSocialAuthTokensCmd()
	IdentityProviderCmd.AddCommand(ListSocialAuthTokensCmd)
}
