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

var BehaviorCmd = &cobra.Command{
	Use:  "behavior",
	Long: "Manage BehaviorAPI",
}

func init() {
	rootCmd.AddCommand(BehaviorCmd)
}

var (
	CreateBehaviorDetectionRuledata string

	CreateBehaviorDetectionRuleRestoreFile string

	CreateBehaviorDetectionRuleQuiet bool
)

func NewCreateBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createDetectionRule",
		Long: "Create a Behavior Detection Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateBehaviorDetectionRuleRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateBehaviorDetectionRuleRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateBehaviorDetectionRuledata = string(processedData)

				if !CreateBehaviorDetectionRuleQuiet {
					fmt.Println("Restoring Behavior from:", CreateBehaviorDetectionRuleRestoreFile)
				}
			}

			req := apiClient.BehaviorAPI.CreateBehaviorDetectionRule(apiClient.GetConfig().Context)

			if CreateBehaviorDetectionRuledata != "" {
				req = req.Data(CreateBehaviorDetectionRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateBehaviorDetectionRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateBehaviorDetectionRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateBehaviorDetectionRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateBehaviorDetectionRuleRestoreFile, "restore-from", "r", "", "Restore Behavior from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateBehaviorDetectionRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateBehaviorDetectionRuleCmd := NewCreateBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(CreateBehaviorDetectionRuleCmd)
}

var (
	ListBehaviorDetectionRulesBackupDir string

	ListBehaviorDetectionRulesLimit    int32
	ListBehaviorDetectionRulesPage     string
	ListBehaviorDetectionRulesFetchAll bool

	ListBehaviorDetectionRulesQuiet bool
)

func NewListBehaviorDetectionRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listDetectionRules",
		Long: "List all Behavior Detection Rules",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.ListBehaviorDetectionRules(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListBehaviorDetectionRulesQuiet {
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
					if !ListBehaviorDetectionRulesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListBehaviorDetectionRulesFetchAll || len(items) == 0 {
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

			if ListBehaviorDetectionRulesFetchAll && pageCount > 1 && !ListBehaviorDetectionRulesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListBehaviorDetectionRulesBackupDir, "behavior", "listDetectionRules")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListBehaviorDetectionRulesQuiet {
					fmt.Printf("Backing up Behaviors to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListBehaviorDetectionRulesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListBehaviorDetectionRulesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListBehaviorDetectionRulesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListBehaviorDetectionRulesQuiet {
					fmt.Printf("Successfully backed up %d/%d Behaviors\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListBehaviorDetectionRulesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListBehaviorDetectionRulesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListBehaviorDetectionRulesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListBehaviorDetectionRulesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Behaviors to a directory")

	cmd.Flags().StringVarP(&ListBehaviorDetectionRulesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListBehaviorDetectionRulesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListBehaviorDetectionRulesCmd := NewListBehaviorDetectionRulesCmd()
	BehaviorCmd.AddCommand(ListBehaviorDetectionRulesCmd)
}

var (
	GetBehaviorDetectionRulebehaviorId string

	GetBehaviorDetectionRuleBackupDir string

	GetBehaviorDetectionRuleQuiet bool
)

func NewGetBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDetectionRule",
		Long: "Retrieve a Behavior Detection Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.GetBehaviorDetectionRule(apiClient.GetConfig().Context, GetBehaviorDetectionRulebehaviorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetBehaviorDetectionRuleQuiet {
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
				dirPath := filepath.Join(GetBehaviorDetectionRuleBackupDir, "behavior", "getDetectionRule")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetBehaviorDetectionRulebehaviorId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetBehaviorDetectionRuleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetBehaviorDetectionRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Behavior to a file")

	cmd.Flags().StringVarP(&GetBehaviorDetectionRuleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetBehaviorDetectionRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetBehaviorDetectionRuleCmd := NewGetBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(GetBehaviorDetectionRuleCmd)
}

var (
	ReplaceBehaviorDetectionRulebehaviorId string

	ReplaceBehaviorDetectionRuledata string

	ReplaceBehaviorDetectionRuleQuiet bool
)

func NewReplaceBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceDetectionRule",
		Long: "Replace a Behavior Detection Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.ReplaceBehaviorDetectionRule(apiClient.GetConfig().Context, ReplaceBehaviorDetectionRulebehaviorId)

			if ReplaceBehaviorDetectionRuledata != "" {
				req = req.Data(ReplaceBehaviorDetectionRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceBehaviorDetectionRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceBehaviorDetectionRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	cmd.Flags().StringVarP(&ReplaceBehaviorDetectionRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceBehaviorDetectionRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceBehaviorDetectionRuleCmd := NewReplaceBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(ReplaceBehaviorDetectionRuleCmd)
}

var (
	DeleteBehaviorDetectionRulebehaviorId string

	DeleteBehaviorDetectionRuleQuiet bool
)

func NewDeleteBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteDetectionRule",
		Long: "Delete a Behavior Detection Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.DeleteBehaviorDetectionRule(apiClient.GetConfig().Context, DeleteBehaviorDetectionRulebehaviorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteBehaviorDetectionRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteBehaviorDetectionRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	cmd.Flags().BoolVarP(&DeleteBehaviorDetectionRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteBehaviorDetectionRuleCmd := NewDeleteBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(DeleteBehaviorDetectionRuleCmd)
}

var (
	ActivateBehaviorDetectionRulebehaviorId string

	ActivateBehaviorDetectionRuleQuiet bool
)

func NewActivateBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateDetectionRule",
		Long: "Activate a Behavior Detection Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.ActivateBehaviorDetectionRule(apiClient.GetConfig().Context, ActivateBehaviorDetectionRulebehaviorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateBehaviorDetectionRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateBehaviorDetectionRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	cmd.Flags().BoolVarP(&ActivateBehaviorDetectionRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateBehaviorDetectionRuleCmd := NewActivateBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(ActivateBehaviorDetectionRuleCmd)
}

var (
	DeactivateBehaviorDetectionRulebehaviorId string

	DeactivateBehaviorDetectionRuleQuiet bool
)

func NewDeactivateBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateDetectionRule",
		Long: "Deactivate a Behavior Detection Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.DeactivateBehaviorDetectionRule(apiClient.GetConfig().Context, DeactivateBehaviorDetectionRulebehaviorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateBehaviorDetectionRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateBehaviorDetectionRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	cmd.Flags().BoolVarP(&DeactivateBehaviorDetectionRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateBehaviorDetectionRuleCmd := NewDeactivateBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(DeactivateBehaviorDetectionRuleCmd)
}
