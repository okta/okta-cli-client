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

var ApplicationCredentialsCmd = &cobra.Command{
	Use:  "applicationCredentials",
	Long: "Manage ApplicationCredentialsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationCredentialsCmd)
}

var (
	GenerateCsrForApplicationappId string

	GenerateCsrForApplicationdata string

	GenerateCsrForApplicationQuiet bool
)

func NewGenerateCsrForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateCsrForApplication",
		Long: "Generate a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GenerateCsrForApplication(apiClient.GetConfig().Context, GenerateCsrForApplicationappId)

			if GenerateCsrForApplicationdata != "" {
				req = req.Data(GenerateCsrForApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GenerateCsrForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GenerateCsrForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GenerateCsrForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GenerateCsrForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&GenerateCsrForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GenerateCsrForApplicationCmd := NewGenerateCsrForApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(GenerateCsrForApplicationCmd)
}

var (
	ListCsrsForApplicationappId string

	ListCsrsForApplicationBackupDir string

	ListCsrsForApplicationLimit    int32
	ListCsrsForApplicationPage     string
	ListCsrsForApplicationFetchAll bool

	ListCsrsForApplicationQuiet bool
)

func NewListCsrsForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listCsrsForApplication",
		Long: "List all Certificate Signing Requests",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.ListCsrsForApplication(apiClient.GetConfig().Context, ListCsrsForApplicationappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListCsrsForApplicationQuiet {
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
					if !ListCsrsForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListCsrsForApplicationFetchAll || len(items) == 0 {
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

			if ListCsrsForApplicationFetchAll && pageCount > 1 && !ListCsrsForApplicationQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListCsrsForApplicationBackupDir, "applicationcredentials", "listCsrsForApplication")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListCsrsForApplicationQuiet {
					fmt.Printf("Backing up ApplicationCredentialss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListCsrsForApplicationQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListCsrsForApplicationQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListCsrsForApplicationQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListCsrsForApplicationQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationCredentialss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListCsrsForApplicationQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListCsrsForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListCsrsForApplicationLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListCsrsForApplicationPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListCsrsForApplicationFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationCredentialss to a directory")

	cmd.Flags().StringVarP(&ListCsrsForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListCsrsForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListCsrsForApplicationCmd := NewListCsrsForApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(ListCsrsForApplicationCmd)
}

var (
	GetCsrForApplicationappId string

	GetCsrForApplicationcsrId string

	GetCsrForApplicationBackupDir string

	GetCsrForApplicationQuiet bool
)

func NewGetCsrForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCsrForApplication",
		Long: "Retrieve a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GetCsrForApplication(apiClient.GetConfig().Context, GetCsrForApplicationappId, GetCsrForApplicationcsrId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCsrForApplicationQuiet {
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
				dirPath := filepath.Join(GetCsrForApplicationBackupDir, "applicationcredentials", "getCsrForApplication")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCsrForApplicationappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCsrForApplicationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCsrForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCsrForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetCsrForApplicationcsrId, "csrId", "", "", "")
	cmd.MarkFlagRequired("csrId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationCredentials to a file")

	cmd.Flags().StringVarP(&GetCsrForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCsrForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCsrForApplicationCmd := NewGetCsrForApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(GetCsrForApplicationCmd)
}

var (
	RevokeCsrFromApplicationappId string

	RevokeCsrFromApplicationcsrId string

	RevokeCsrFromApplicationQuiet bool
)

func NewRevokeCsrFromApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeCsrFromApplication",
		Long: "Revoke a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.RevokeCsrFromApplication(apiClient.GetConfig().Context, RevokeCsrFromApplicationappId, RevokeCsrFromApplicationcsrId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeCsrFromApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeCsrFromApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeCsrFromApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&RevokeCsrFromApplicationcsrId, "csrId", "", "", "")
	cmd.MarkFlagRequired("csrId")

	cmd.Flags().BoolVarP(&RevokeCsrFromApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeCsrFromApplicationCmd := NewRevokeCsrFromApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(RevokeCsrFromApplicationCmd)
}

var (
	PublishCsrFromApplicationappId string

	PublishCsrFromApplicationcsrId string

	PublishCsrFromApplicationdata string

	PublishCsrFromApplicationQuiet bool
)

func NewPublishCsrFromApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "publishCsrFromApplication",
		Long: "Publish a Certificate Signing Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.PublishCsrFromApplication(apiClient.GetConfig().Context, PublishCsrFromApplicationappId, PublishCsrFromApplicationcsrId)

			if PublishCsrFromApplicationdata != "" {
				req = req.Data(PublishCsrFromApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !PublishCsrFromApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !PublishCsrFromApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&PublishCsrFromApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&PublishCsrFromApplicationcsrId, "csrId", "", "", "")
	cmd.MarkFlagRequired("csrId")

	cmd.Flags().StringVarP(&PublishCsrFromApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&PublishCsrFromApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	PublishCsrFromApplicationCmd := NewPublishCsrFromApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(PublishCsrFromApplicationCmd)
}

var (
	ListApplicationKeysappId string

	ListApplicationKeysBackupDir string

	ListApplicationKeysLimit    int32
	ListApplicationKeysPage     string
	ListApplicationKeysFetchAll bool

	ListApplicationKeysQuiet bool
)

func NewListApplicationKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApplicationKeys",
		Long: "List all Key Credentials",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.ListApplicationKeys(apiClient.GetConfig().Context, ListApplicationKeysappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApplicationKeysQuiet {
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
					if !ListApplicationKeysQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApplicationKeysFetchAll || len(items) == 0 {
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

			if ListApplicationKeysFetchAll && pageCount > 1 && !ListApplicationKeysQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApplicationKeysBackupDir, "applicationcredentials", "listApplicationKeys")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApplicationKeysQuiet {
					fmt.Printf("Backing up ApplicationCredentialss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApplicationKeysQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApplicationKeysQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApplicationKeysQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApplicationKeysQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationCredentialss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApplicationKeysQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListApplicationKeysappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListApplicationKeysLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApplicationKeysPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApplicationKeysFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationCredentialss to a directory")

	cmd.Flags().StringVarP(&ListApplicationKeysBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApplicationKeysQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApplicationKeysCmd := NewListApplicationKeysCmd()
	ApplicationCredentialsCmd.AddCommand(ListApplicationKeysCmd)
}

var (
	GenerateApplicationKeyappId string

	GenerateApplicationKeyQuiet bool
)

func NewGenerateApplicationKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateApplicationKey",
		Long: "Generate a Key Credential",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GenerateApplicationKey(apiClient.GetConfig().Context, GenerateApplicationKeyappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GenerateApplicationKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GenerateApplicationKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GenerateApplicationKeyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&GenerateApplicationKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GenerateApplicationKeyCmd := NewGenerateApplicationKeyCmd()
	ApplicationCredentialsCmd.AddCommand(GenerateApplicationKeyCmd)
}

var (
	GetApplicationKeyappId string

	GetApplicationKeykeyId string

	GetApplicationKeyBackupDir string

	GetApplicationKeyQuiet bool
)

func NewGetApplicationKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationKey",
		Long: "Retrieve a Key Credential",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GetApplicationKey(apiClient.GetConfig().Context, GetApplicationKeyappId, GetApplicationKeykeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApplicationKeyQuiet {
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
				dirPath := filepath.Join(GetApplicationKeyBackupDir, "applicationcredentials", "getApplicationKey")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApplicationKeyappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApplicationKeyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApplicationKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApplicationKeyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetApplicationKeykeyId, "keyId", "", "", "")
	cmd.MarkFlagRequired("keyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationCredentials to a file")

	cmd.Flags().StringVarP(&GetApplicationKeyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApplicationKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApplicationKeyCmd := NewGetApplicationKeyCmd()
	ApplicationCredentialsCmd.AddCommand(GetApplicationKeyCmd)
}

var (
	CloneApplicationKeyappId string

	CloneApplicationKeykeyId string

	CloneApplicationKeyQuiet bool
)

func NewCloneApplicationKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "cloneApplicationKey",
		Long: "Clone a Key Credential",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.CloneApplicationKey(apiClient.GetConfig().Context, CloneApplicationKeyappId, CloneApplicationKeykeyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CloneApplicationKeyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CloneApplicationKeyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CloneApplicationKeyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&CloneApplicationKeykeyId, "keyId", "", "", "")
	cmd.MarkFlagRequired("keyId")

	cmd.Flags().BoolVarP(&CloneApplicationKeyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CloneApplicationKeyCmd := NewCloneApplicationKeyCmd()
	ApplicationCredentialsCmd.AddCommand(CloneApplicationKeyCmd)
}
