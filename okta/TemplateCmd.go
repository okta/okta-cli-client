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

var TemplateCmd = &cobra.Command{
	Use:  "template",
	Long: "Manage TemplateAPI",
}

func init() {
	rootCmd.AddCommand(TemplateCmd)
}

var (
	CreateSmsTemplatedata string

	CreateSmsTemplateRestoreFile string

	CreateSmsTemplateQuiet bool
)

func NewCreateSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createSms",
		Long: "Create an SMS Template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateSmsTemplateRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateSmsTemplateRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateSmsTemplatedata = string(processedData)

				if !CreateSmsTemplateQuiet {
					fmt.Println("Restoring Template from:", CreateSmsTemplateRestoreFile)
				}
			}

			req := apiClient.TemplateAPI.CreateSmsTemplate(apiClient.GetConfig().Context)

			if CreateSmsTemplatedata != "" {
				req = req.Data(CreateSmsTemplatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateSmsTemplateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateSmsTemplateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateSmsTemplatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateSmsTemplateRestoreFile, "restore-from", "r", "", "Restore Template from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateSmsTemplateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateSmsTemplateCmd := NewCreateSmsTemplateCmd()
	TemplateCmd.AddCommand(CreateSmsTemplateCmd)
}

var (
	ListSmsTemplatesBackupDir string

	ListSmsTemplatesLimit    int32
	ListSmsTemplatesPage     string
	ListSmsTemplatesFetchAll bool

	ListSmsTemplatesQuiet bool
)

func NewListSmsTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSmss",
		Long: "List all SMS Templates",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.ListSmsTemplates(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListSmsTemplatesQuiet {
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
					if !ListSmsTemplatesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListSmsTemplatesFetchAll || len(items) == 0 {
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

			if ListSmsTemplatesFetchAll && pageCount > 1 && !ListSmsTemplatesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListSmsTemplatesBackupDir, "template", "listSmss")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListSmsTemplatesQuiet {
					fmt.Printf("Backing up Templates to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListSmsTemplatesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListSmsTemplatesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListSmsTemplatesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListSmsTemplatesQuiet {
					fmt.Printf("Successfully backed up %d/%d Templates\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListSmsTemplatesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListSmsTemplatesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListSmsTemplatesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListSmsTemplatesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Templates to a directory")

	cmd.Flags().StringVarP(&ListSmsTemplatesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListSmsTemplatesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListSmsTemplatesCmd := NewListSmsTemplatesCmd()
	TemplateCmd.AddCommand(ListSmsTemplatesCmd)
}

var (
	UpdateSmsTemplatetemplateId string

	UpdateSmsTemplatedata string

	UpdateSmsTemplateQuiet bool
)

func NewUpdateSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateSms",
		Long: "Update an SMS Template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.UpdateSmsTemplate(apiClient.GetConfig().Context, UpdateSmsTemplatetemplateId)

			if UpdateSmsTemplatedata != "" {
				req = req.Data(UpdateSmsTemplatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateSmsTemplateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateSmsTemplateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	cmd.Flags().StringVarP(&UpdateSmsTemplatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateSmsTemplateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateSmsTemplateCmd := NewUpdateSmsTemplateCmd()
	TemplateCmd.AddCommand(UpdateSmsTemplateCmd)
}

var (
	GetSmsTemplatetemplateId string

	GetSmsTemplateBackupDir string

	GetSmsTemplateQuiet bool
)

func NewGetSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSms",
		Long: "Retrieve an SMS Template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.GetSmsTemplate(apiClient.GetConfig().Context, GetSmsTemplatetemplateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetSmsTemplateQuiet {
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
				dirPath := filepath.Join(GetSmsTemplateBackupDir, "template", "getSms")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetSmsTemplatetemplateId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetSmsTemplateQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetSmsTemplateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Template to a file")

	cmd.Flags().StringVarP(&GetSmsTemplateBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetSmsTemplateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetSmsTemplateCmd := NewGetSmsTemplateCmd()
	TemplateCmd.AddCommand(GetSmsTemplateCmd)
}

var (
	ReplaceSmsTemplatetemplateId string

	ReplaceSmsTemplatedata string

	ReplaceSmsTemplateQuiet bool
)

func NewReplaceSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceSms",
		Long: "Replace an SMS Template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.ReplaceSmsTemplate(apiClient.GetConfig().Context, ReplaceSmsTemplatetemplateId)

			if ReplaceSmsTemplatedata != "" {
				req = req.Data(ReplaceSmsTemplatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceSmsTemplateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceSmsTemplateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	cmd.Flags().StringVarP(&ReplaceSmsTemplatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceSmsTemplateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceSmsTemplateCmd := NewReplaceSmsTemplateCmd()
	TemplateCmd.AddCommand(ReplaceSmsTemplateCmd)
}

var (
	DeleteSmsTemplatetemplateId string

	DeleteSmsTemplateQuiet bool
)

func NewDeleteSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteSms",
		Long: "Delete an SMS Template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.DeleteSmsTemplate(apiClient.GetConfig().Context, DeleteSmsTemplatetemplateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteSmsTemplateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteSmsTemplateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	cmd.Flags().BoolVarP(&DeleteSmsTemplateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteSmsTemplateCmd := NewDeleteSmsTemplateCmd()
	TemplateCmd.AddCommand(DeleteSmsTemplateCmd)
}
