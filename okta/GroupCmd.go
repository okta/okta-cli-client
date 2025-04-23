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

var GroupCmd = &cobra.Command{
	Use:  "group",
	Long: "Manage GroupAPI",
}

func init() {
	rootCmd.AddCommand(GroupCmd)
}

var (
	CreateGroupdata string

	CreateGroupRestoreFile string

	CreateGroupQuiet bool
)

func NewCreateGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateGroupRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateGroupRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateGroupdata = string(processedData)

				if !CreateGroupQuiet {
					fmt.Println("Restoring Group from:", CreateGroupRestoreFile)
				}
			}

			req := apiClient.GroupAPI.CreateGroup(apiClient.GetConfig().Context)

			if CreateGroupdata != "" {
				req = req.Data(CreateGroupdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateGroupdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateGroupRestoreFile, "restore-from", "r", "", "Restore Group from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateGroupCmd := NewCreateGroupCmd()
	GroupCmd.AddCommand(CreateGroupCmd)
}

var (
	ListGroupsBackupDir string

	ListGroupsLimit    int32
	ListGroupsPage     string
	ListGroupsFetchAll bool

	ListGroupsQuiet bool
)

func NewListGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Groups",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListGroups(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupsQuiet {
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
					if !ListGroupsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupsFetchAll || len(items) == 0 {
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

			if ListGroupsFetchAll && pageCount > 1 && !ListGroupsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupsBackupDir, "group", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupsQuiet {
					fmt.Printf("Backing up Groups to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupsQuiet {
					fmt.Printf("Successfully backed up %d/%d Groups\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListGroupsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Groups to a directory")

	cmd.Flags().StringVarP(&ListGroupsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupsCmd := NewListGroupsCmd()
	GroupCmd.AddCommand(ListGroupsCmd)
}

var (
	CreateGroupRuledata string

	CreateGroupRuleRestoreFile string

	CreateGroupRuleQuiet bool
)

func NewCreateGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createRule",
		Long: "Create a Group Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateGroupRuleRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateGroupRuleRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateGroupRuledata = string(processedData)

				if !CreateGroupRuleQuiet {
					fmt.Println("Restoring Group from:", CreateGroupRuleRestoreFile)
				}
			}

			req := apiClient.GroupAPI.CreateGroupRule(apiClient.GetConfig().Context)

			if CreateGroupRuledata != "" {
				req = req.Data(CreateGroupRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateGroupRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateGroupRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateGroupRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateGroupRuleRestoreFile, "restore-from", "r", "", "Restore Group from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateGroupRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateGroupRuleCmd := NewCreateGroupRuleCmd()
	GroupCmd.AddCommand(CreateGroupRuleCmd)
}

var (
	ListGroupRulesBackupDir string

	ListGroupRulesLimit    int32
	ListGroupRulesPage     string
	ListGroupRulesFetchAll bool

	ListGroupRulesQuiet bool
)

func NewListGroupRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listRules",
		Long: "List all Group Rules",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListGroupRules(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupRulesQuiet {
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
					if !ListGroupRulesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupRulesFetchAll || len(items) == 0 {
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

			if ListGroupRulesFetchAll && pageCount > 1 && !ListGroupRulesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupRulesBackupDir, "group", "listRules")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupRulesQuiet {
					fmt.Printf("Backing up Groups to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupRulesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupRulesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupRulesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupRulesQuiet {
					fmt.Printf("Successfully backed up %d/%d Groups\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupRulesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListGroupRulesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupRulesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupRulesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Groups to a directory")

	cmd.Flags().StringVarP(&ListGroupRulesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupRulesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupRulesCmd := NewListGroupRulesCmd()
	GroupCmd.AddCommand(ListGroupRulesCmd)
}

var (
	GetGroupRulegroupRuleId string

	GetGroupRuleBackupDir string

	GetGroupRuleQuiet bool
)

func NewGetGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getRule",
		Long: "Retrieve a Group Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.GetGroupRule(apiClient.GetConfig().Context, GetGroupRulegroupRuleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetGroupRuleQuiet {
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
				dirPath := filepath.Join(GetGroupRuleBackupDir, "group", "getRule")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetGroupRulegroupRuleId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetGroupRuleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetGroupRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Group to a file")

	cmd.Flags().StringVarP(&GetGroupRuleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetGroupRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetGroupRuleCmd := NewGetGroupRuleCmd()
	GroupCmd.AddCommand(GetGroupRuleCmd)
}

var (
	ReplaceGroupRulegroupRuleId string

	ReplaceGroupRuledata string

	ReplaceGroupRuleQuiet bool
)

func NewReplaceGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceRule",
		Long: "Replace a Group Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ReplaceGroupRule(apiClient.GetConfig().Context, ReplaceGroupRulegroupRuleId)

			if ReplaceGroupRuledata != "" {
				req = req.Data(ReplaceGroupRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceGroupRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceGroupRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().StringVarP(&ReplaceGroupRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceGroupRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceGroupRuleCmd := NewReplaceGroupRuleCmd()
	GroupCmd.AddCommand(ReplaceGroupRuleCmd)
}

var (
	DeleteGroupRulegroupRuleId string

	DeleteGroupRuleQuiet bool
)

func NewDeleteGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteRule",
		Long: "Delete a group Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.DeleteGroupRule(apiClient.GetConfig().Context, DeleteGroupRulegroupRuleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteGroupRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteGroupRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().BoolVarP(&DeleteGroupRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteGroupRuleCmd := NewDeleteGroupRuleCmd()
	GroupCmd.AddCommand(DeleteGroupRuleCmd)
}

var (
	ActivateGroupRulegroupRuleId string

	ActivateGroupRuleQuiet bool
)

func NewActivateGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateRule",
		Long: "Activate a Group Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ActivateGroupRule(apiClient.GetConfig().Context, ActivateGroupRulegroupRuleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateGroupRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateGroupRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().BoolVarP(&ActivateGroupRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateGroupRuleCmd := NewActivateGroupRuleCmd()
	GroupCmd.AddCommand(ActivateGroupRuleCmd)
}

var (
	DeactivateGroupRulegroupRuleId string

	DeactivateGroupRuleQuiet bool
)

func NewDeactivateGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateRule",
		Long: "Deactivate a Group Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.DeactivateGroupRule(apiClient.GetConfig().Context, DeactivateGroupRulegroupRuleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateGroupRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateGroupRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().BoolVarP(&DeactivateGroupRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateGroupRuleCmd := NewDeactivateGroupRuleCmd()
	GroupCmd.AddCommand(DeactivateGroupRuleCmd)
}

var (
	GetGroupgroupId string

	GetGroupBackupDir string

	GetGroupQuiet bool
)

func NewGetGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.GetGroup(apiClient.GetConfig().Context, GetGroupgroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetGroupQuiet {
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
				dirPath := filepath.Join(GetGroupBackupDir, "group", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetGroupgroupId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetGroupQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Group to a file")

	cmd.Flags().StringVarP(&GetGroupBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetGroupCmd := NewGetGroupCmd()
	GroupCmd.AddCommand(GetGroupCmd)
}

var (
	ReplaceGroupgroupId string

	ReplaceGroupdata string

	ReplaceGroupQuiet bool
)

func NewReplaceGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ReplaceGroup(apiClient.GetConfig().Context, ReplaceGroupgroupId)

			if ReplaceGroupdata != "" {
				req = req.Data(ReplaceGroupdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&ReplaceGroupdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceGroupCmd := NewReplaceGroupCmd()
	GroupCmd.AddCommand(ReplaceGroupCmd)
}

var (
	DeleteGroupgroupId string

	DeleteGroupQuiet bool
)

func NewDeleteGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.DeleteGroup(apiClient.GetConfig().Context, DeleteGroupgroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().BoolVarP(&DeleteGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteGroupCmd := NewDeleteGroupCmd()
	GroupCmd.AddCommand(DeleteGroupCmd)
}

var (
	ListAssignedApplicationsForGroupgroupId string

	ListAssignedApplicationsForGroupBackupDir string

	ListAssignedApplicationsForGroupLimit    int32
	ListAssignedApplicationsForGroupPage     string
	ListAssignedApplicationsForGroupFetchAll bool

	ListAssignedApplicationsForGroupQuiet bool
)

func NewListAssignedApplicationsForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAssignedApplicationsFor",
		Long: "List all Assigned Applications",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListAssignedApplicationsForGroup(apiClient.GetConfig().Context, ListAssignedApplicationsForGroupgroupId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAssignedApplicationsForGroupQuiet {
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
					if !ListAssignedApplicationsForGroupQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAssignedApplicationsForGroupFetchAll || len(items) == 0 {
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

			if ListAssignedApplicationsForGroupFetchAll && pageCount > 1 && !ListAssignedApplicationsForGroupQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAssignedApplicationsForGroupBackupDir, "group", "listAssignedApplicationsFor")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAssignedApplicationsForGroupQuiet {
					fmt.Printf("Backing up Groups to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAssignedApplicationsForGroupQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAssignedApplicationsForGroupQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAssignedApplicationsForGroupQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAssignedApplicationsForGroupQuiet {
					fmt.Printf("Successfully backed up %d/%d Groups\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAssignedApplicationsForGroupQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAssignedApplicationsForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().Int32VarP(&ListAssignedApplicationsForGroupLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAssignedApplicationsForGroupPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAssignedApplicationsForGroupFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Groups to a directory")

	cmd.Flags().StringVarP(&ListAssignedApplicationsForGroupBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAssignedApplicationsForGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAssignedApplicationsForGroupCmd := NewListAssignedApplicationsForGroupCmd()
	GroupCmd.AddCommand(ListAssignedApplicationsForGroupCmd)
}

var (
	ListGroupUsersgroupId string

	ListGroupUsersBackupDir string

	ListGroupUsersLimit    int32
	ListGroupUsersPage     string
	ListGroupUsersFetchAll bool

	ListGroupUsersQuiet bool
)

func NewListGroupUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listUsers",
		Long: "List all Member Users",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListGroupUsers(apiClient.GetConfig().Context, ListGroupUsersgroupId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupUsersQuiet {
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
					if !ListGroupUsersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupUsersFetchAll || len(items) == 0 {
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

			if ListGroupUsersFetchAll && pageCount > 1 && !ListGroupUsersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupUsersBackupDir, "group", "listUsers")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupUsersQuiet {
					fmt.Printf("Backing up Groups to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupUsersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupUsersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupUsersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupUsersQuiet {
					fmt.Printf("Successfully backed up %d/%d Groups\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupUsersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListGroupUsersgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().Int32VarP(&ListGroupUsersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupUsersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupUsersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Groups to a directory")

	cmd.Flags().StringVarP(&ListGroupUsersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupUsersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupUsersCmd := NewListGroupUsersCmd()
	GroupCmd.AddCommand(ListGroupUsersCmd)
}

var (
	AssignUserToGroupgroupId string

	AssignUserToGroupuserId string

	AssignUserToGroupQuiet bool
)

func NewAssignUserToGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignUserTo",
		Long: "Assign a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.AssignUserToGroup(apiClient.GetConfig().Context, AssignUserToGroupgroupId, AssignUserToGroupuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignUserToGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignUserToGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignUserToGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignUserToGroupuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&AssignUserToGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignUserToGroupCmd := NewAssignUserToGroupCmd()
	GroupCmd.AddCommand(AssignUserToGroupCmd)
}

var (
	UnassignUserFromGroupgroupId string

	UnassignUserFromGroupuserId string

	UnassignUserFromGroupQuiet bool
)

func NewUnassignUserFromGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignUserFrom",
		Long: "Unassign a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.UnassignUserFromGroup(apiClient.GetConfig().Context, UnassignUserFromGroupgroupId, UnassignUserFromGroupuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignUserFromGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignUserFromGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignUserFromGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&UnassignUserFromGroupuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&UnassignUserFromGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignUserFromGroupCmd := NewUnassignUserFromGroupCmd()
	GroupCmd.AddCommand(UnassignUserFromGroupCmd)
}
