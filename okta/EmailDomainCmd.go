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

var EmailDomainCmd = &cobra.Command{
	Use:  "emailDomain",
	Long: "Manage EmailDomainAPI",
}

func init() {
	rootCmd.AddCommand(EmailDomainCmd)
}

var (
	CreateEmailDomaindata string

	CreateEmailDomainRestoreFile string

	CreateEmailDomainQuiet bool
)

func NewCreateEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create an Email Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateEmailDomainRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateEmailDomainRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateEmailDomaindata = string(processedData)

				if !CreateEmailDomainQuiet {
					fmt.Println("Restoring EmailDomain from:", CreateEmailDomainRestoreFile)
				}
			}

			req := apiClient.EmailDomainAPI.CreateEmailDomain(apiClient.GetConfig().Context)

			if CreateEmailDomaindata != "" {
				req = req.Data(CreateEmailDomaindata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateEmailDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateEmailDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateEmailDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateEmailDomainRestoreFile, "restore-from", "r", "", "Restore EmailDomain from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateEmailDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateEmailDomainCmd := NewCreateEmailDomainCmd()
	EmailDomainCmd.AddCommand(CreateEmailDomainCmd)
}

var (
	ListEmailDomainsBackupDir string

	ListEmailDomainsLimit    int32
	ListEmailDomainsPage     string
	ListEmailDomainsFetchAll bool

	ListEmailDomainsQuiet bool
)

func NewListEmailDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Email Domains",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.ListEmailDomains(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListEmailDomainsQuiet {
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
					if !ListEmailDomainsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListEmailDomainsFetchAll || len(items) == 0 {
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

			if ListEmailDomainsFetchAll && pageCount > 1 && !ListEmailDomainsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListEmailDomainsBackupDir, "emaildomain", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListEmailDomainsQuiet {
					fmt.Printf("Backing up EmailDomains to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListEmailDomainsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListEmailDomainsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListEmailDomainsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListEmailDomainsQuiet {
					fmt.Printf("Successfully backed up %d/%d EmailDomains\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListEmailDomainsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListEmailDomainsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListEmailDomainsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListEmailDomainsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple EmailDomains to a directory")

	cmd.Flags().StringVarP(&ListEmailDomainsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListEmailDomainsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListEmailDomainsCmd := NewListEmailDomainsCmd()
	EmailDomainCmd.AddCommand(ListEmailDomainsCmd)
}

var (
	GetEmailDomainemailDomainId string

	GetEmailDomainBackupDir string

	GetEmailDomainQuiet bool
)

func NewGetEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an Email Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.GetEmailDomain(apiClient.GetConfig().Context, GetEmailDomainemailDomainId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailDomainQuiet {
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
				dirPath := filepath.Join(GetEmailDomainBackupDir, "emaildomain", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailDomainemailDomainId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailDomainQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the EmailDomain to a file")

	cmd.Flags().StringVarP(&GetEmailDomainBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailDomainCmd := NewGetEmailDomainCmd()
	EmailDomainCmd.AddCommand(GetEmailDomainCmd)
}

var (
	ReplaceEmailDomainemailDomainId string

	ReplaceEmailDomaindata string

	ReplaceEmailDomainQuiet bool
)

func NewReplaceEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace an Email Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.ReplaceEmailDomain(apiClient.GetConfig().Context, ReplaceEmailDomainemailDomainId)

			if ReplaceEmailDomaindata != "" {
				req = req.Data(ReplaceEmailDomaindata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceEmailDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceEmailDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	cmd.Flags().StringVarP(&ReplaceEmailDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceEmailDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceEmailDomainCmd := NewReplaceEmailDomainCmd()
	EmailDomainCmd.AddCommand(ReplaceEmailDomainCmd)
}

var (
	DeleteEmailDomainemailDomainId string

	DeleteEmailDomainQuiet bool
)

func NewDeleteEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an Email Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.DeleteEmailDomain(apiClient.GetConfig().Context, DeleteEmailDomainemailDomainId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteEmailDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteEmailDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	cmd.Flags().BoolVarP(&DeleteEmailDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteEmailDomainCmd := NewDeleteEmailDomainCmd()
	EmailDomainCmd.AddCommand(DeleteEmailDomainCmd)
}

var (
	VerifyEmailDomainemailDomainId string

	VerifyEmailDomainQuiet bool
)

func NewVerifyEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verify",
		Long: "Verify an Email Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.VerifyEmailDomain(apiClient.GetConfig().Context, VerifyEmailDomainemailDomainId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !VerifyEmailDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !VerifyEmailDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&VerifyEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	cmd.Flags().BoolVarP(&VerifyEmailDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	VerifyEmailDomainCmd := NewVerifyEmailDomainCmd()
	EmailDomainCmd.AddCommand(VerifyEmailDomainCmd)
}
