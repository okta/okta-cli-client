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

var ApiServiceIntegrationsCmd = &cobra.Command{
	Use:  "apiServiceIntegrations",
	Long: "Manage ApiServiceIntegrationsAPI",
}

func init() {
	rootCmd.AddCommand(ApiServiceIntegrationsCmd)
}

var (
	CreateApiServiceIntegrationInstancedata string

	CreateApiServiceIntegrationInstanceRestoreFile string

	CreateApiServiceIntegrationInstanceQuiet bool
)

func NewCreateApiServiceIntegrationInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createApiServiceIntegrationInstance",
		Long: "Create an API Service Integration instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateApiServiceIntegrationInstanceRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateApiServiceIntegrationInstanceRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateApiServiceIntegrationInstancedata = string(processedData)

				if !CreateApiServiceIntegrationInstanceQuiet {
					fmt.Println("Restoring ApiServiceIntegrations from:", CreateApiServiceIntegrationInstanceRestoreFile)
				}
			}

			req := apiClient.ApiServiceIntegrationsAPI.CreateApiServiceIntegrationInstance(apiClient.GetConfig().Context)

			if CreateApiServiceIntegrationInstancedata != "" {
				req = req.Data(CreateApiServiceIntegrationInstancedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateApiServiceIntegrationInstanceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateApiServiceIntegrationInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateApiServiceIntegrationInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateApiServiceIntegrationInstanceRestoreFile, "restore-from", "r", "", "Restore ApiServiceIntegrations from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateApiServiceIntegrationInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateApiServiceIntegrationInstanceCmd := NewCreateApiServiceIntegrationInstanceCmd()
	ApiServiceIntegrationsCmd.AddCommand(CreateApiServiceIntegrationInstanceCmd)
}

var (
	ListApiServiceIntegrationInstancesBackupDir string

	ListApiServiceIntegrationInstancesLimit    int32
	ListApiServiceIntegrationInstancesPage     string
	ListApiServiceIntegrationInstancesFetchAll bool

	ListApiServiceIntegrationInstancesQuiet bool
)

func NewListApiServiceIntegrationInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApiServiceIntegrationInstances",
		Long: "List all API Service Integration instances",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.ListApiServiceIntegrationInstances(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApiServiceIntegrationInstancesQuiet {
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
					if !ListApiServiceIntegrationInstancesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApiServiceIntegrationInstancesFetchAll || len(items) == 0 {
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

			if ListApiServiceIntegrationInstancesFetchAll && pageCount > 1 && !ListApiServiceIntegrationInstancesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApiServiceIntegrationInstancesBackupDir, "apiserviceintegrations", "listApiServiceIntegrationInstances")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApiServiceIntegrationInstancesQuiet {
					fmt.Printf("Backing up ApiServiceIntegrationss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApiServiceIntegrationInstancesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApiServiceIntegrationInstancesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApiServiceIntegrationInstancesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApiServiceIntegrationInstancesQuiet {
					fmt.Printf("Successfully backed up %d/%d ApiServiceIntegrationss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApiServiceIntegrationInstancesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListApiServiceIntegrationInstancesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApiServiceIntegrationInstancesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApiServiceIntegrationInstancesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApiServiceIntegrationss to a directory")

	cmd.Flags().StringVarP(&ListApiServiceIntegrationInstancesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApiServiceIntegrationInstancesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApiServiceIntegrationInstancesCmd := NewListApiServiceIntegrationInstancesCmd()
	ApiServiceIntegrationsCmd.AddCommand(ListApiServiceIntegrationInstancesCmd)
}

var (
	GetApiServiceIntegrationInstanceapiServiceId string

	GetApiServiceIntegrationInstanceBackupDir string

	GetApiServiceIntegrationInstanceQuiet bool
)

func NewGetApiServiceIntegrationInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApiServiceIntegrationInstance",
		Long: "Retrieve an API Service Integration instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.GetApiServiceIntegrationInstance(apiClient.GetConfig().Context, GetApiServiceIntegrationInstanceapiServiceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApiServiceIntegrationInstanceQuiet {
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
				dirPath := filepath.Join(GetApiServiceIntegrationInstanceBackupDir, "apiserviceintegrations", "getApiServiceIntegrationInstance")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApiServiceIntegrationInstanceapiServiceId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApiServiceIntegrationInstanceQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApiServiceIntegrationInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApiServiceIntegrationInstanceapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApiServiceIntegrations to a file")

	cmd.Flags().StringVarP(&GetApiServiceIntegrationInstanceBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApiServiceIntegrationInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApiServiceIntegrationInstanceCmd := NewGetApiServiceIntegrationInstanceCmd()
	ApiServiceIntegrationsCmd.AddCommand(GetApiServiceIntegrationInstanceCmd)
}

var (
	DeleteApiServiceIntegrationInstanceapiServiceId string

	DeleteApiServiceIntegrationInstanceQuiet bool
)

func NewDeleteApiServiceIntegrationInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteApiServiceIntegrationInstance",
		Long: "Delete an API Service Integration instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.DeleteApiServiceIntegrationInstance(apiClient.GetConfig().Context, DeleteApiServiceIntegrationInstanceapiServiceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteApiServiceIntegrationInstanceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteApiServiceIntegrationInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteApiServiceIntegrationInstanceapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().BoolVarP(&DeleteApiServiceIntegrationInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteApiServiceIntegrationInstanceCmd := NewDeleteApiServiceIntegrationInstanceCmd()
	ApiServiceIntegrationsCmd.AddCommand(DeleteApiServiceIntegrationInstanceCmd)
}

var (
	CreateApiServiceIntegrationInstanceSecretapiServiceId string

	CreateApiServiceIntegrationInstanceSecretRestoreFile string

	CreateApiServiceIntegrationInstanceSecretQuiet bool
)

func NewCreateApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createApiServiceIntegrationInstanceSecret",
		Long: "Create an API Service Integration instance Secret",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateApiServiceIntegrationInstanceSecretRestoreFile != "" {

				fmt.Println("Warning: This Create operation doesn't support data input. Cannot restore from file.")
				return fmt.Errorf("restore from file not supported for this operation that doesn't accept data input")

				if !CreateApiServiceIntegrationInstanceSecretQuiet {
					fmt.Println("Restoring ApiServiceIntegrations from:", CreateApiServiceIntegrationInstanceSecretRestoreFile)
				}
			}

			req := apiClient.ApiServiceIntegrationsAPI.CreateApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, CreateApiServiceIntegrationInstanceSecretapiServiceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateApiServiceIntegrationInstanceSecretQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateApiServiceIntegrationInstanceSecretQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&CreateApiServiceIntegrationInstanceSecretRestoreFile, "restore-from", "r", "", "Restore ApiServiceIntegrations from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateApiServiceIntegrationInstanceSecretQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateApiServiceIntegrationInstanceSecretCmd := NewCreateApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(CreateApiServiceIntegrationInstanceSecretCmd)
}

var (
	ListApiServiceIntegrationInstanceSecretsapiServiceId string

	ListApiServiceIntegrationInstanceSecretsBackupDir string

	ListApiServiceIntegrationInstanceSecretsLimit    int32
	ListApiServiceIntegrationInstanceSecretsPage     string
	ListApiServiceIntegrationInstanceSecretsFetchAll bool

	ListApiServiceIntegrationInstanceSecretsQuiet bool
)

func NewListApiServiceIntegrationInstanceSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApiServiceIntegrationInstanceSecrets",
		Long: "List all API Service Integration instance Secrets",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.ListApiServiceIntegrationInstanceSecrets(apiClient.GetConfig().Context, ListApiServiceIntegrationInstanceSecretsapiServiceId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApiServiceIntegrationInstanceSecretsQuiet {
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
					if !ListApiServiceIntegrationInstanceSecretsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApiServiceIntegrationInstanceSecretsFetchAll || len(items) == 0 {
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

			if ListApiServiceIntegrationInstanceSecretsFetchAll && pageCount > 1 && !ListApiServiceIntegrationInstanceSecretsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApiServiceIntegrationInstanceSecretsBackupDir, "apiserviceintegrations", "listApiServiceIntegrationInstanceSecrets")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApiServiceIntegrationInstanceSecretsQuiet {
					fmt.Printf("Backing up ApiServiceIntegrationss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApiServiceIntegrationInstanceSecretsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApiServiceIntegrationInstanceSecretsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApiServiceIntegrationInstanceSecretsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApiServiceIntegrationInstanceSecretsQuiet {
					fmt.Printf("Successfully backed up %d/%d ApiServiceIntegrationss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApiServiceIntegrationInstanceSecretsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListApiServiceIntegrationInstanceSecretsapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().Int32VarP(&ListApiServiceIntegrationInstanceSecretsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApiServiceIntegrationInstanceSecretsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApiServiceIntegrationInstanceSecretsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApiServiceIntegrationss to a directory")

	cmd.Flags().StringVarP(&ListApiServiceIntegrationInstanceSecretsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApiServiceIntegrationInstanceSecretsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApiServiceIntegrationInstanceSecretsCmd := NewListApiServiceIntegrationInstanceSecretsCmd()
	ApiServiceIntegrationsCmd.AddCommand(ListApiServiceIntegrationInstanceSecretsCmd)
}

var (
	DeleteApiServiceIntegrationInstanceSecretapiServiceId string

	DeleteApiServiceIntegrationInstanceSecretsecretId string

	DeleteApiServiceIntegrationInstanceSecretQuiet bool
)

func NewDeleteApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteApiServiceIntegrationInstanceSecret",
		Long: "Delete an API Service Integration instance Secret",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.DeleteApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, DeleteApiServiceIntegrationInstanceSecretapiServiceId, DeleteApiServiceIntegrationInstanceSecretsecretId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteApiServiceIntegrationInstanceSecretQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteApiServiceIntegrationInstanceSecretQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&DeleteApiServiceIntegrationInstanceSecretsecretId, "secretId", "", "", "")
	cmd.MarkFlagRequired("secretId")

	cmd.Flags().BoolVarP(&DeleteApiServiceIntegrationInstanceSecretQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteApiServiceIntegrationInstanceSecretCmd := NewDeleteApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(DeleteApiServiceIntegrationInstanceSecretCmd)
}

var (
	ActivateApiServiceIntegrationInstanceSecretapiServiceId string

	ActivateApiServiceIntegrationInstanceSecretsecretId string

	ActivateApiServiceIntegrationInstanceSecretQuiet bool
)

func NewActivateApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateApiServiceIntegrationInstanceSecret",
		Long: "Activate an API Service Integration instance Secret",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.ActivateApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, ActivateApiServiceIntegrationInstanceSecretapiServiceId, ActivateApiServiceIntegrationInstanceSecretsecretId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateApiServiceIntegrationInstanceSecretQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateApiServiceIntegrationInstanceSecretQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&ActivateApiServiceIntegrationInstanceSecretsecretId, "secretId", "", "", "")
	cmd.MarkFlagRequired("secretId")

	cmd.Flags().BoolVarP(&ActivateApiServiceIntegrationInstanceSecretQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateApiServiceIntegrationInstanceSecretCmd := NewActivateApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(ActivateApiServiceIntegrationInstanceSecretCmd)
}

var (
	DeactivateApiServiceIntegrationInstanceSecretapiServiceId string

	DeactivateApiServiceIntegrationInstanceSecretsecretId string

	DeactivateApiServiceIntegrationInstanceSecretQuiet bool
)

func NewDeactivateApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateApiServiceIntegrationInstanceSecret",
		Long: "Deactivate an API Service Integration instance Secret",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.DeactivateApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, DeactivateApiServiceIntegrationInstanceSecretapiServiceId, DeactivateApiServiceIntegrationInstanceSecretsecretId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateApiServiceIntegrationInstanceSecretQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateApiServiceIntegrationInstanceSecretQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&DeactivateApiServiceIntegrationInstanceSecretsecretId, "secretId", "", "", "")
	cmd.MarkFlagRequired("secretId")

	cmd.Flags().BoolVarP(&DeactivateApiServiceIntegrationInstanceSecretQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateApiServiceIntegrationInstanceSecretCmd := NewDeactivateApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(DeactivateApiServiceIntegrationInstanceSecretCmd)
}
