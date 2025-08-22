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

var CAPTCHACmd = &cobra.Command{
	Use:  "cAPTCHA",
	Long: "Manage CAPTCHAAPI",
}

func init() {
	rootCmd.AddCommand(CAPTCHACmd)
}

var (
	CreateCaptchaInstancedata string

	CreateCaptchaInstanceRestoreFile string

	CreateCaptchaInstanceQuiet bool
)

func NewCreateCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createCaptchaInstance",
		Long: "Create a CAPTCHA instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateCaptchaInstanceRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateCaptchaInstanceRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateCaptchaInstancedata = string(processedData)

				if !CreateCaptchaInstanceQuiet {
					fmt.Println("Restoring CAPTCHA from:", CreateCaptchaInstanceRestoreFile)
				}
			}

			req := apiClient.CAPTCHAAPI.CreateCaptchaInstance(apiClient.GetConfig().Context)

			if CreateCaptchaInstancedata != "" {
				req = req.Data(CreateCaptchaInstancedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateCaptchaInstanceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateCaptchaInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateCaptchaInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateCaptchaInstanceRestoreFile, "restore-from", "r", "", "Restore CAPTCHA from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateCaptchaInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateCaptchaInstanceCmd := NewCreateCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(CreateCaptchaInstanceCmd)
}

var (
	ListCaptchaInstancesBackupDir string

	ListCaptchaInstancesLimit    int32
	ListCaptchaInstancesPage     string
	ListCaptchaInstancesFetchAll bool

	ListCaptchaInstancesQuiet bool
)

func NewListCaptchaInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listCaptchaInstances",
		Long: "List all CAPTCHA Instances",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.ListCaptchaInstances(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListCaptchaInstancesQuiet {
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
					if !ListCaptchaInstancesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListCaptchaInstancesFetchAll || len(items) == 0 {
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

			if ListCaptchaInstancesFetchAll && pageCount > 1 && !ListCaptchaInstancesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListCaptchaInstancesBackupDir, "captcha", "listCaptchaInstances")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListCaptchaInstancesQuiet {
					fmt.Printf("Backing up CAPTCHAs to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListCaptchaInstancesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListCaptchaInstancesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListCaptchaInstancesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListCaptchaInstancesQuiet {
					fmt.Printf("Successfully backed up %d/%d CAPTCHAs\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListCaptchaInstancesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListCaptchaInstancesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListCaptchaInstancesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListCaptchaInstancesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple CAPTCHAs to a directory")

	cmd.Flags().StringVarP(&ListCaptchaInstancesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListCaptchaInstancesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListCaptchaInstancesCmd := NewListCaptchaInstancesCmd()
	CAPTCHACmd.AddCommand(ListCaptchaInstancesCmd)
}

var (
	UpdateCaptchaInstancecaptchaId string

	UpdateCaptchaInstancedata string

	UpdateCaptchaInstanceQuiet bool
)

func NewUpdateCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateCaptchaInstance",
		Long: "Update a CAPTCHA Instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.UpdateCaptchaInstance(apiClient.GetConfig().Context, UpdateCaptchaInstancecaptchaId)

			if UpdateCaptchaInstancedata != "" {
				req = req.Data(UpdateCaptchaInstancedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateCaptchaInstanceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateCaptchaInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	cmd.Flags().StringVarP(&UpdateCaptchaInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateCaptchaInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateCaptchaInstanceCmd := NewUpdateCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(UpdateCaptchaInstanceCmd)
}

var (
	GetCaptchaInstancecaptchaId string

	GetCaptchaInstanceBackupDir string

	GetCaptchaInstanceQuiet bool
)

func NewGetCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCaptchaInstance",
		Long: "Retrieve a CAPTCHA Instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.GetCaptchaInstance(apiClient.GetConfig().Context, GetCaptchaInstancecaptchaId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCaptchaInstanceQuiet {
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
				dirPath := filepath.Join(GetCaptchaInstanceBackupDir, "captcha", "getCaptchaInstance")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCaptchaInstancecaptchaId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCaptchaInstanceQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCaptchaInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the CAPTCHA to a file")

	cmd.Flags().StringVarP(&GetCaptchaInstanceBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCaptchaInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCaptchaInstanceCmd := NewGetCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(GetCaptchaInstanceCmd)
}

var (
	ReplaceCaptchaInstancecaptchaId string

	ReplaceCaptchaInstancedata string

	ReplaceCaptchaInstanceQuiet bool
)

func NewReplaceCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceCaptchaInstance",
		Long: "Replace a CAPTCHA Instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.ReplaceCaptchaInstance(apiClient.GetConfig().Context, ReplaceCaptchaInstancecaptchaId)

			if ReplaceCaptchaInstancedata != "" {
				req = req.Data(ReplaceCaptchaInstancedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceCaptchaInstanceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceCaptchaInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	cmd.Flags().StringVarP(&ReplaceCaptchaInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceCaptchaInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceCaptchaInstanceCmd := NewReplaceCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(ReplaceCaptchaInstanceCmd)
}

var (
	DeleteCaptchaInstancecaptchaId string

	DeleteCaptchaInstanceQuiet bool
)

func NewDeleteCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteCaptchaInstance",
		Long: "Delete a CAPTCHA Instance",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.DeleteCaptchaInstance(apiClient.GetConfig().Context, DeleteCaptchaInstancecaptchaId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteCaptchaInstanceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteCaptchaInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	cmd.Flags().BoolVarP(&DeleteCaptchaInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteCaptchaInstanceCmd := NewDeleteCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(DeleteCaptchaInstanceCmd)
}

var (
	GetOrgCaptchaSettingsBackupDir string

	GetOrgCaptchaSettingsQuiet bool
)

func NewGetOrgCaptchaSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOrgCaptchaSettings",
		Long: "Retrieve the Org-wide CAPTCHA Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.GetOrgCaptchaSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOrgCaptchaSettingsQuiet {
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
				dirPath := filepath.Join(GetOrgCaptchaSettingsBackupDir, "captcha", "getOrgCaptchaSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "captcha.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOrgCaptchaSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOrgCaptchaSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the CAPTCHA to a file")

	cmd.Flags().StringVarP(&GetOrgCaptchaSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOrgCaptchaSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOrgCaptchaSettingsCmd := NewGetOrgCaptchaSettingsCmd()
	CAPTCHACmd.AddCommand(GetOrgCaptchaSettingsCmd)
}

var (
	ReplacesOrgCaptchaSettingsdata string

	ReplacesOrgCaptchaSettingsQuiet bool
)

func NewReplacesOrgCaptchaSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacesOrgCaptchaSettings",
		Long: "Replace the Org-wide CAPTCHA Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.ReplacesOrgCaptchaSettings(apiClient.GetConfig().Context)

			if ReplacesOrgCaptchaSettingsdata != "" {
				req = req.Data(ReplacesOrgCaptchaSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacesOrgCaptchaSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacesOrgCaptchaSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacesOrgCaptchaSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacesOrgCaptchaSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacesOrgCaptchaSettingsCmd := NewReplacesOrgCaptchaSettingsCmd()
	CAPTCHACmd.AddCommand(ReplacesOrgCaptchaSettingsCmd)
}

var DeleteOrgCaptchaSettingsQuiet bool

func NewDeleteOrgCaptchaSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteOrgCaptchaSettings",
		Long: "Delete the Org-wide CAPTCHA Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.DeleteOrgCaptchaSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteOrgCaptchaSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteOrgCaptchaSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&DeleteOrgCaptchaSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteOrgCaptchaSettingsCmd := NewDeleteOrgCaptchaSettingsCmd()
	CAPTCHACmd.AddCommand(DeleteOrgCaptchaSettingsCmd)
}
