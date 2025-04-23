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

var CustomDomainCmd = &cobra.Command{
	Use:  "customDomain",
	Long: "Manage CustomDomainAPI",
}

func init() {
	rootCmd.AddCommand(CustomDomainCmd)
}

var (
	CreateCustomDomaindata string

	CreateCustomDomainRestoreFile string

	CreateCustomDomainQuiet bool
)

func NewCreateCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Custom Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateCustomDomainRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateCustomDomainRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateCustomDomaindata = string(processedData)

				if !CreateCustomDomainQuiet {
					fmt.Println("Restoring CustomDomain from:", CreateCustomDomainRestoreFile)
				}
			}

			req := apiClient.CustomDomainAPI.CreateCustomDomain(apiClient.GetConfig().Context)

			if CreateCustomDomaindata != "" {
				req = req.Data(CreateCustomDomaindata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateCustomDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateCustomDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateCustomDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateCustomDomainRestoreFile, "restore-from", "r", "", "Restore CustomDomain from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateCustomDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateCustomDomainCmd := NewCreateCustomDomainCmd()
	CustomDomainCmd.AddCommand(CreateCustomDomainCmd)
}

var (
	ListCustomDomainsBackupDir string

	ListCustomDomainsLimit    int32
	ListCustomDomainsPage     string
	ListCustomDomainsFetchAll bool

	ListCustomDomainsQuiet bool
)

func NewListCustomDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Custom Domains",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.ListCustomDomains(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListCustomDomainsQuiet {
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
					if !ListCustomDomainsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListCustomDomainsFetchAll || len(items) == 0 {
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

			if ListCustomDomainsFetchAll && pageCount > 1 && !ListCustomDomainsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListCustomDomainsBackupDir, "customdomain", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListCustomDomainsQuiet {
					fmt.Printf("Backing up CustomDomains to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListCustomDomainsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListCustomDomainsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListCustomDomainsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListCustomDomainsQuiet {
					fmt.Printf("Successfully backed up %d/%d CustomDomains\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListCustomDomainsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListCustomDomainsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListCustomDomainsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListCustomDomainsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple CustomDomains to a directory")

	cmd.Flags().StringVarP(&ListCustomDomainsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListCustomDomainsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListCustomDomainsCmd := NewListCustomDomainsCmd()
	CustomDomainCmd.AddCommand(ListCustomDomainsCmd)
}

var (
	GetCustomDomaindomainId string

	GetCustomDomainBackupDir string

	GetCustomDomainQuiet bool
)

func NewGetCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Custom Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.GetCustomDomain(apiClient.GetConfig().Context, GetCustomDomaindomainId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCustomDomainQuiet {
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
				dirPath := filepath.Join(GetCustomDomainBackupDir, "customdomain", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCustomDomaindomainId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCustomDomainQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCustomDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCustomDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the CustomDomain to a file")

	cmd.Flags().StringVarP(&GetCustomDomainBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCustomDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCustomDomainCmd := NewGetCustomDomainCmd()
	CustomDomainCmd.AddCommand(GetCustomDomainCmd)
}

var (
	ReplaceCustomDomaindomainId string

	ReplaceCustomDomaindata string

	ReplaceCustomDomainQuiet bool
)

func NewReplaceCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Custom Domain's Brand",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.ReplaceCustomDomain(apiClient.GetConfig().Context, ReplaceCustomDomaindomainId)

			if ReplaceCustomDomaindata != "" {
				req = req.Data(ReplaceCustomDomaindata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceCustomDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceCustomDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceCustomDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().StringVarP(&ReplaceCustomDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceCustomDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceCustomDomainCmd := NewReplaceCustomDomainCmd()
	CustomDomainCmd.AddCommand(ReplaceCustomDomainCmd)
}

var (
	DeleteCustomDomaindomainId string

	DeleteCustomDomainQuiet bool
)

func NewDeleteCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Custom Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.DeleteCustomDomain(apiClient.GetConfig().Context, DeleteCustomDomaindomainId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteCustomDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteCustomDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteCustomDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().BoolVarP(&DeleteCustomDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteCustomDomainCmd := NewDeleteCustomDomainCmd()
	CustomDomainCmd.AddCommand(DeleteCustomDomainCmd)
}

var (
	UpsertCertificatedomainId string

	UpsertCertificatedata string

	UpsertCertificateQuiet bool
)

func NewUpsertCertificateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "upsertCertificate",
		Long: "Upsert the Custom Domain's Certificate",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.UpsertCertificate(apiClient.GetConfig().Context, UpsertCertificatedomainId)

			if UpsertCertificatedata != "" {
				req = req.Data(UpsertCertificatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpsertCertificateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpsertCertificateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpsertCertificatedomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().StringVarP(&UpsertCertificatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpsertCertificateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpsertCertificateCmd := NewUpsertCertificateCmd()
	CustomDomainCmd.AddCommand(UpsertCertificateCmd)
}

var (
	VerifyDomaindomainId string

	VerifyDomainQuiet bool
)

func NewVerifyDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verifyDomain",
		Long: "Verify a Custom Domain",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.VerifyDomain(apiClient.GetConfig().Context, VerifyDomaindomainId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !VerifyDomainQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !VerifyDomainQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&VerifyDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().BoolVarP(&VerifyDomainQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	VerifyDomainCmd := NewVerifyDomainCmd()
	CustomDomainCmd.AddCommand(VerifyDomainCmd)
}
