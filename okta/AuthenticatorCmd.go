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

var AuthenticatorCmd = &cobra.Command{
	Use:  "authenticator",
	Long: "Manage AuthenticatorAPI",
}

func init() {
	rootCmd.AddCommand(AuthenticatorCmd)
}

var (
	GetWellKnownAppAuthenticatorConfigurationBackupDir string

	GetWellKnownAppAuthenticatorConfigurationQuiet bool
)

func NewGetWellKnownAppAuthenticatorConfigurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getWellKnownAppConfiguration",
		Long: "Retrieve the Well-Known App Authenticator Configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.GetWellKnownAppAuthenticatorConfiguration(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetWellKnownAppAuthenticatorConfigurationQuiet {
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
				dirPath := filepath.Join(GetWellKnownAppAuthenticatorConfigurationBackupDir, "authenticator", "getWellKnownAppConfiguration")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "authenticator.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetWellKnownAppAuthenticatorConfigurationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetWellKnownAppAuthenticatorConfigurationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the Authenticator to a file")

	cmd.Flags().StringVarP(&GetWellKnownAppAuthenticatorConfigurationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetWellKnownAppAuthenticatorConfigurationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetWellKnownAppAuthenticatorConfigurationCmd := NewGetWellKnownAppAuthenticatorConfigurationCmd()
	AuthenticatorCmd.AddCommand(GetWellKnownAppAuthenticatorConfigurationCmd)
}

var (
	CreateAuthenticatordata string

	CreateAuthenticatorRestoreFile string

	CreateAuthenticatorQuiet bool
)

func NewCreateAuthenticatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Authenticator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateAuthenticatorRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateAuthenticatorRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateAuthenticatordata = string(processedData)

				if !CreateAuthenticatorQuiet {
					fmt.Println("Restoring Authenticator from:", CreateAuthenticatorRestoreFile)
				}
			}

			req := apiClient.AuthenticatorAPI.CreateAuthenticator(apiClient.GetConfig().Context)

			if CreateAuthenticatordata != "" {
				req = req.Data(CreateAuthenticatordata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateAuthenticatorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateAuthenticatorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateAuthenticatordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateAuthenticatorRestoreFile, "restore-from", "r", "", "Restore Authenticator from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateAuthenticatorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateAuthenticatorCmd := NewCreateAuthenticatorCmd()
	AuthenticatorCmd.AddCommand(CreateAuthenticatorCmd)
}

var (
	ListAuthenticatorsBackupDir string

	ListAuthenticatorsLimit    int32
	ListAuthenticatorsPage     string
	ListAuthenticatorsFetchAll bool

	ListAuthenticatorsQuiet bool
)

func NewListAuthenticatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Authenticators",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.ListAuthenticators(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAuthenticatorsQuiet {
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
					if !ListAuthenticatorsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAuthenticatorsFetchAll || len(items) == 0 {
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

			if ListAuthenticatorsFetchAll && pageCount > 1 && !ListAuthenticatorsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAuthenticatorsBackupDir, "authenticator", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAuthenticatorsQuiet {
					fmt.Printf("Backing up Authenticators to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAuthenticatorsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAuthenticatorsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAuthenticatorsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAuthenticatorsQuiet {
					fmt.Printf("Successfully backed up %d/%d Authenticators\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAuthenticatorsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListAuthenticatorsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAuthenticatorsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAuthenticatorsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Authenticators to a directory")

	cmd.Flags().StringVarP(&ListAuthenticatorsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAuthenticatorsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAuthenticatorsCmd := NewListAuthenticatorsCmd()
	AuthenticatorCmd.AddCommand(ListAuthenticatorsCmd)
}

var (
	GetAuthenticatorauthenticatorId string

	GetAuthenticatorBackupDir string

	GetAuthenticatorQuiet bool
)

func NewGetAuthenticatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Authenticator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.GetAuthenticator(apiClient.GetConfig().Context, GetAuthenticatorauthenticatorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAuthenticatorQuiet {
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
				dirPath := filepath.Join(GetAuthenticatorBackupDir, "authenticator", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAuthenticatorauthenticatorId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAuthenticatorQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAuthenticatorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Authenticator to a file")

	cmd.Flags().StringVarP(&GetAuthenticatorBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAuthenticatorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAuthenticatorCmd := NewGetAuthenticatorCmd()
	AuthenticatorCmd.AddCommand(GetAuthenticatorCmd)
}

var (
	ReplaceAuthenticatorauthenticatorId string

	ReplaceAuthenticatordata string

	ReplaceAuthenticatorQuiet bool
)

func NewReplaceAuthenticatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Authenticator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.ReplaceAuthenticator(apiClient.GetConfig().Context, ReplaceAuthenticatorauthenticatorId)

			if ReplaceAuthenticatordata != "" {
				req = req.Data(ReplaceAuthenticatordata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceAuthenticatorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceAuthenticatorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().StringVarP(&ReplaceAuthenticatordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceAuthenticatorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceAuthenticatorCmd := NewReplaceAuthenticatorCmd()
	AuthenticatorCmd.AddCommand(ReplaceAuthenticatorCmd)
}

var (
	ActivateAuthenticatorauthenticatorId string

	ActivateAuthenticatorQuiet bool
)

func NewActivateAuthenticatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate an Authenticator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.ActivateAuthenticator(apiClient.GetConfig().Context, ActivateAuthenticatorauthenticatorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateAuthenticatorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateAuthenticatorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().BoolVarP(&ActivateAuthenticatorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateAuthenticatorCmd := NewActivateAuthenticatorCmd()
	AuthenticatorCmd.AddCommand(ActivateAuthenticatorCmd)
}

var (
	DeactivateAuthenticatorauthenticatorId string

	DeactivateAuthenticatorQuiet bool
)

func NewDeactivateAuthenticatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate an Authenticator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.DeactivateAuthenticator(apiClient.GetConfig().Context, DeactivateAuthenticatorauthenticatorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateAuthenticatorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateAuthenticatorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().BoolVarP(&DeactivateAuthenticatorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateAuthenticatorCmd := NewDeactivateAuthenticatorCmd()
	AuthenticatorCmd.AddCommand(DeactivateAuthenticatorCmd)
}

var (
	ListAuthenticatorMethodsauthenticatorId string

	ListAuthenticatorMethodsBackupDir string

	ListAuthenticatorMethodsLimit    int32
	ListAuthenticatorMethodsPage     string
	ListAuthenticatorMethodsFetchAll bool

	ListAuthenticatorMethodsQuiet bool
)

func NewListAuthenticatorMethodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listMethods",
		Long: "List all Methods of an Authenticator",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.ListAuthenticatorMethods(apiClient.GetConfig().Context, ListAuthenticatorMethodsauthenticatorId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAuthenticatorMethodsQuiet {
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
					if !ListAuthenticatorMethodsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAuthenticatorMethodsFetchAll || len(items) == 0 {
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

			if ListAuthenticatorMethodsFetchAll && pageCount > 1 && !ListAuthenticatorMethodsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAuthenticatorMethodsBackupDir, "authenticator", "listMethods")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAuthenticatorMethodsQuiet {
					fmt.Printf("Backing up Authenticators to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAuthenticatorMethodsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAuthenticatorMethodsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAuthenticatorMethodsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAuthenticatorMethodsQuiet {
					fmt.Printf("Successfully backed up %d/%d Authenticators\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAuthenticatorMethodsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAuthenticatorMethodsauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().Int32VarP(&ListAuthenticatorMethodsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAuthenticatorMethodsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAuthenticatorMethodsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Authenticators to a directory")

	cmd.Flags().StringVarP(&ListAuthenticatorMethodsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAuthenticatorMethodsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAuthenticatorMethodsCmd := NewListAuthenticatorMethodsCmd()
	AuthenticatorCmd.AddCommand(ListAuthenticatorMethodsCmd)
}

var (
	GetAuthenticatorMethodauthenticatorId string

	GetAuthenticatorMethodmethodType string

	GetAuthenticatorMethodBackupDir string

	GetAuthenticatorMethodQuiet bool
)

func NewGetAuthenticatorMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getMethod",
		Long: "Retrieve a Method",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.GetAuthenticatorMethod(apiClient.GetConfig().Context, GetAuthenticatorMethodauthenticatorId, GetAuthenticatorMethodmethodType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAuthenticatorMethodQuiet {
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
				dirPath := filepath.Join(GetAuthenticatorMethodBackupDir, "authenticator", "getMethod")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAuthenticatorMethodauthenticatorId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAuthenticatorMethodQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAuthenticatorMethodQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().StringVarP(&GetAuthenticatorMethodmethodType, "methodType", "", "", "")
	cmd.MarkFlagRequired("methodType")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Authenticator to a file")

	cmd.Flags().StringVarP(&GetAuthenticatorMethodBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAuthenticatorMethodQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAuthenticatorMethodCmd := NewGetAuthenticatorMethodCmd()
	AuthenticatorCmd.AddCommand(GetAuthenticatorMethodCmd)
}

var (
	ReplaceAuthenticatorMethodauthenticatorId string

	ReplaceAuthenticatorMethodmethodType string

	ReplaceAuthenticatorMethoddata string

	ReplaceAuthenticatorMethodQuiet bool
)

func NewReplaceAuthenticatorMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceMethod",
		Long: "Replace a Method",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.ReplaceAuthenticatorMethod(apiClient.GetConfig().Context, ReplaceAuthenticatorMethodauthenticatorId, ReplaceAuthenticatorMethodmethodType)

			if ReplaceAuthenticatorMethoddata != "" {
				req = req.Data(ReplaceAuthenticatorMethoddata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceAuthenticatorMethodQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceAuthenticatorMethodQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().StringVarP(&ReplaceAuthenticatorMethodmethodType, "methodType", "", "", "")
	cmd.MarkFlagRequired("methodType")

	cmd.Flags().StringVarP(&ReplaceAuthenticatorMethoddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceAuthenticatorMethodQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceAuthenticatorMethodCmd := NewReplaceAuthenticatorMethodCmd()
	AuthenticatorCmd.AddCommand(ReplaceAuthenticatorMethodCmd)
}

var (
	ActivateAuthenticatorMethodauthenticatorId string

	ActivateAuthenticatorMethodmethodType string

	ActivateAuthenticatorMethodQuiet bool
)

func NewActivateAuthenticatorMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateMethod",
		Long: "Activate an Authenticator Method",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.ActivateAuthenticatorMethod(apiClient.GetConfig().Context, ActivateAuthenticatorMethodauthenticatorId, ActivateAuthenticatorMethodmethodType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateAuthenticatorMethodQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateAuthenticatorMethodQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().StringVarP(&ActivateAuthenticatorMethodmethodType, "methodType", "", "", "")
	cmd.MarkFlagRequired("methodType")

	cmd.Flags().BoolVarP(&ActivateAuthenticatorMethodQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateAuthenticatorMethodCmd := NewActivateAuthenticatorMethodCmd()
	AuthenticatorCmd.AddCommand(ActivateAuthenticatorMethodCmd)
}

var (
	DeactivateAuthenticatorMethodauthenticatorId string

	DeactivateAuthenticatorMethodmethodType string

	DeactivateAuthenticatorMethodQuiet bool
)

func NewDeactivateAuthenticatorMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateMethod",
		Long: "Deactivate an Authenticator Method",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthenticatorAPI.DeactivateAuthenticatorMethod(apiClient.GetConfig().Context, DeactivateAuthenticatorMethodauthenticatorId, DeactivateAuthenticatorMethodmethodType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateAuthenticatorMethodQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateAuthenticatorMethodQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
	cmd.MarkFlagRequired("authenticatorId")

	cmd.Flags().StringVarP(&DeactivateAuthenticatorMethodmethodType, "methodType", "", "", "")
	cmd.MarkFlagRequired("methodType")

	cmd.Flags().BoolVarP(&DeactivateAuthenticatorMethodQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateAuthenticatorMethodCmd := NewDeactivateAuthenticatorMethodCmd()
	AuthenticatorCmd.AddCommand(DeactivateAuthenticatorMethodCmd)
}
