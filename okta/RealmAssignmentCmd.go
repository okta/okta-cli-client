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

var RealmAssignmentCmd = &cobra.Command{
	Use:  "realmAssignment",
	Long: "Manage RealmAssignmentAPI",
}

func init() {
	rootCmd.AddCommand(RealmAssignmentCmd)
}

var (
	CreateRealmAssignmentdata string

	CreateRealmAssignmentRestoreFile string

	CreateRealmAssignmentQuiet bool
)

func NewCreateRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateRealmAssignmentRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateRealmAssignmentRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateRealmAssignmentdata = string(processedData)

				if !CreateRealmAssignmentQuiet {
					fmt.Println("Restoring RealmAssignment from:", CreateRealmAssignmentRestoreFile)
				}
			}

			req := apiClient.RealmAssignmentAPI.CreateRealmAssignment(apiClient.GetConfig().Context)

			if CreateRealmAssignmentdata != "" {
				req = req.Data(CreateRealmAssignmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateRealmAssignmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateRealmAssignmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateRealmAssignmentRestoreFile, "restore-from", "r", "", "Restore RealmAssignment from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateRealmAssignmentCmd := NewCreateRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(CreateRealmAssignmentCmd)
}

var (
	ListRealmAssignmentsBackupDir string

	ListRealmAssignmentsLimit    int32
	ListRealmAssignmentsPage     string
	ListRealmAssignmentsFetchAll bool

	ListRealmAssignmentsQuiet bool
)

func NewListRealmAssignmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Realm Assignments",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ListRealmAssignments(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRealmAssignmentsQuiet {
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
					if !ListRealmAssignmentsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRealmAssignmentsFetchAll || len(items) == 0 {
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

			if ListRealmAssignmentsFetchAll && pageCount > 1 && !ListRealmAssignmentsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRealmAssignmentsBackupDir, "realmassignment", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRealmAssignmentsQuiet {
					fmt.Printf("Backing up RealmAssignments to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRealmAssignmentsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRealmAssignmentsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRealmAssignmentsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRealmAssignmentsQuiet {
					fmt.Printf("Successfully backed up %d/%d RealmAssignments\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRealmAssignmentsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListRealmAssignmentsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRealmAssignmentsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRealmAssignmentsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RealmAssignments to a directory")

	cmd.Flags().StringVarP(&ListRealmAssignmentsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRealmAssignmentsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRealmAssignmentsCmd := NewListRealmAssignmentsCmd()
	RealmAssignmentCmd.AddCommand(ListRealmAssignmentsCmd)
}

var (
	ExecuteRealmAssignmentdata string

	ExecuteRealmAssignmentQuiet bool
)

func NewExecuteRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "execute",
		Long: "Execute a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ExecuteRealmAssignment(apiClient.GetConfig().Context)

			if ExecuteRealmAssignmentdata != "" {
				req = req.Data(ExecuteRealmAssignmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ExecuteRealmAssignmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ExecuteRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ExecuteRealmAssignmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ExecuteRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ExecuteRealmAssignmentCmd := NewExecuteRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(ExecuteRealmAssignmentCmd)
}

var (
	ListRealmAssignmentOperationsBackupDir string

	ListRealmAssignmentOperationsLimit    int32
	ListRealmAssignmentOperationsPage     string
	ListRealmAssignmentOperationsFetchAll bool

	ListRealmAssignmentOperationsQuiet bool
)

func NewListRealmAssignmentOperationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listOperations",
		Long: "List all Realm Assignment operations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ListRealmAssignmentOperations(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRealmAssignmentOperationsQuiet {
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
					if !ListRealmAssignmentOperationsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRealmAssignmentOperationsFetchAll || len(items) == 0 {
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

			if ListRealmAssignmentOperationsFetchAll && pageCount > 1 && !ListRealmAssignmentOperationsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRealmAssignmentOperationsBackupDir, "realmassignment", "listOperations")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRealmAssignmentOperationsQuiet {
					fmt.Printf("Backing up RealmAssignments to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRealmAssignmentOperationsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRealmAssignmentOperationsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRealmAssignmentOperationsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRealmAssignmentOperationsQuiet {
					fmt.Printf("Successfully backed up %d/%d RealmAssignments\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRealmAssignmentOperationsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListRealmAssignmentOperationsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRealmAssignmentOperationsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRealmAssignmentOperationsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RealmAssignments to a directory")

	cmd.Flags().StringVarP(&ListRealmAssignmentOperationsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRealmAssignmentOperationsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRealmAssignmentOperationsCmd := NewListRealmAssignmentOperationsCmd()
	RealmAssignmentCmd.AddCommand(ListRealmAssignmentOperationsCmd)
}

var (
	GetRealmAssignmentassignmentId string

	GetRealmAssignmentBackupDir string

	GetRealmAssignmentQuiet bool
)

func NewGetRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.GetRealmAssignment(apiClient.GetConfig().Context, GetRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRealmAssignmentQuiet {
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
				dirPath := filepath.Join(GetRealmAssignmentBackupDir, "realmassignment", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetRealmAssignmentassignmentId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRealmAssignmentQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the RealmAssignment to a file")

	cmd.Flags().StringVarP(&GetRealmAssignmentBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRealmAssignmentCmd := NewGetRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(GetRealmAssignmentCmd)
}

var (
	ReplaceRealmAssignmentassignmentId string

	ReplaceRealmAssignmentdata string

	ReplaceRealmAssignmentQuiet bool
)

func NewReplaceRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ReplaceRealmAssignment(apiClient.GetConfig().Context, ReplaceRealmAssignmentassignmentId)

			if ReplaceRealmAssignmentdata != "" {
				req = req.Data(ReplaceRealmAssignmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRealmAssignmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	cmd.Flags().StringVarP(&ReplaceRealmAssignmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRealmAssignmentCmd := NewReplaceRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(ReplaceRealmAssignmentCmd)
}

var (
	DeleteRealmAssignmentassignmentId string

	DeleteRealmAssignmentQuiet bool
)

func NewDeleteRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.DeleteRealmAssignment(apiClient.GetConfig().Context, DeleteRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteRealmAssignmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	cmd.Flags().BoolVarP(&DeleteRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteRealmAssignmentCmd := NewDeleteRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(DeleteRealmAssignmentCmd)
}

var (
	ActivateRealmAssignmentassignmentId string

	ActivateRealmAssignmentQuiet bool
)

func NewActivateRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ActivateRealmAssignment(apiClient.GetConfig().Context, ActivateRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateRealmAssignmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	cmd.Flags().BoolVarP(&ActivateRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateRealmAssignmentCmd := NewActivateRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(ActivateRealmAssignmentCmd)
}

var (
	DeactivateRealmAssignmentassignmentId string

	DeactivateRealmAssignmentQuiet bool
)

func NewDeactivateRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Realm Assignment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.DeactivateRealmAssignment(apiClient.GetConfig().Context, DeactivateRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateRealmAssignmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateRealmAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	cmd.Flags().BoolVarP(&DeactivateRealmAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateRealmAssignmentCmd := NewDeactivateRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(DeactivateRealmAssignmentCmd)
}
