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

var RoleCmd = &cobra.Command{
	Use:  "role",
	Long: "Manage RoleAPI",
}

func init() {
	rootCmd.AddCommand(RoleCmd)
}

var (
	CreateRoledata string

	CreateRoleRestoreFile string

	CreateRoleQuiet bool
)

func NewCreateRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateRoleRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateRoleRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateRoledata = string(processedData)

				if !CreateRoleQuiet {
					fmt.Println("Restoring Role from:", CreateRoleRestoreFile)
				}
			}

			req := apiClient.RoleAPI.CreateRole(apiClient.GetConfig().Context)

			if CreateRoledata != "" {
				req = req.Data(CreateRoledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateRoledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateRoleRestoreFile, "restore-from", "r", "", "Restore Role from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateRoleCmd := NewCreateRoleCmd()
	RoleCmd.AddCommand(CreateRoleCmd)
}

var (
	ListRolesBackupDir string

	ListRolesLimit    int32
	ListRolesPage     string
	ListRolesFetchAll bool

	ListRolesQuiet bool
)

func NewListRolesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Roles",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.ListRoles(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRolesQuiet {
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
					if !ListRolesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRolesFetchAll || len(items) == 0 {
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

			if ListRolesFetchAll && pageCount > 1 && !ListRolesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRolesBackupDir, "role", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRolesQuiet {
					fmt.Printf("Backing up Roles to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRolesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRolesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRolesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRolesQuiet {
					fmt.Printf("Successfully backed up %d/%d Roles\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRolesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListRolesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRolesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRolesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Roles to a directory")

	cmd.Flags().StringVarP(&ListRolesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRolesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRolesCmd := NewListRolesCmd()
	RoleCmd.AddCommand(ListRolesCmd)
}

var (
	GetRoleroleIdOrLabel string

	GetRoleBackupDir string

	GetRoleQuiet bool
)

func NewGetRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.GetRole(apiClient.GetConfig().Context, GetRoleroleIdOrLabel)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRoleQuiet {
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
				dirPath := filepath.Join(GetRoleBackupDir, "role", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetRoleroleIdOrLabel
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRoleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRoleroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Role to a file")

	cmd.Flags().StringVarP(&GetRoleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRoleCmd := NewGetRoleCmd()
	RoleCmd.AddCommand(GetRoleCmd)
}

var (
	ReplaceRoleroleIdOrLabel string

	ReplaceRoledata string

	ReplaceRoleQuiet bool
)

func NewReplaceRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.ReplaceRole(apiClient.GetConfig().Context, ReplaceRoleroleIdOrLabel)

			if ReplaceRoledata != "" {
				req = req.Data(ReplaceRoledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRoleroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&ReplaceRoledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRoleCmd := NewReplaceRoleCmd()
	RoleCmd.AddCommand(ReplaceRoleCmd)
}

var (
	DeleteRoleroleIdOrLabel string

	DeleteRoleQuiet bool
)

func NewDeleteRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.DeleteRole(apiClient.GetConfig().Context, DeleteRoleroleIdOrLabel)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteRoleroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().BoolVarP(&DeleteRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteRoleCmd := NewDeleteRoleCmd()
	RoleCmd.AddCommand(DeleteRoleCmd)
}

var (
	ListRolePermissionsroleIdOrLabel string

	ListRolePermissionsBackupDir string

	ListRolePermissionsLimit    int32
	ListRolePermissionsPage     string
	ListRolePermissionsFetchAll bool

	ListRolePermissionsQuiet bool
)

func NewListRolePermissionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listPermissions",
		Long: "List all Permissions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.ListRolePermissions(apiClient.GetConfig().Context, ListRolePermissionsroleIdOrLabel)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRolePermissionsQuiet {
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
					if !ListRolePermissionsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRolePermissionsFetchAll || len(items) == 0 {
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

			if ListRolePermissionsFetchAll && pageCount > 1 && !ListRolePermissionsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRolePermissionsBackupDir, "role", "listPermissions")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRolePermissionsQuiet {
					fmt.Printf("Backing up Roles to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRolePermissionsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRolePermissionsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRolePermissionsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRolePermissionsQuiet {
					fmt.Printf("Successfully backed up %d/%d Roles\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRolePermissionsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListRolePermissionsroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().Int32VarP(&ListRolePermissionsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRolePermissionsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRolePermissionsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Roles to a directory")

	cmd.Flags().StringVarP(&ListRolePermissionsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRolePermissionsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRolePermissionsCmd := NewListRolePermissionsCmd()
	RoleCmd.AddCommand(ListRolePermissionsCmd)
}

var (
	CreateRolePermissionroleIdOrLabel string

	CreateRolePermissionpermissionType string

	CreateRolePermissiondata string

	CreateRolePermissionRestoreFile string

	CreateRolePermissionQuiet bool
)

func NewCreateRolePermissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createPermission",
		Long: "Create a Permission",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateRolePermissionRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateRolePermissionRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateRolePermissiondata = string(processedData)

				if !CreateRolePermissionQuiet {
					fmt.Println("Restoring Role from:", CreateRolePermissionRestoreFile)
				}
			}

			req := apiClient.RoleAPI.CreateRolePermission(apiClient.GetConfig().Context, CreateRolePermissionroleIdOrLabel, CreateRolePermissionpermissionType)

			if CreateRolePermissiondata != "" {
				req = req.Data(CreateRolePermissiondata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateRolePermissionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateRolePermissionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&CreateRolePermissionpermissionType, "permissionType", "", "", "")
	cmd.MarkFlagRequired("permissionType")

	cmd.Flags().StringVarP(&CreateRolePermissiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateRolePermissionRestoreFile, "restore-from", "r", "", "Restore Role from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateRolePermissionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateRolePermissionCmd := NewCreateRolePermissionCmd()
	RoleCmd.AddCommand(CreateRolePermissionCmd)
}

var (
	GetRolePermissionroleIdOrLabel string

	GetRolePermissionpermissionType string

	GetRolePermissionBackupDir string

	GetRolePermissionQuiet bool
)

func NewGetRolePermissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPermission",
		Long: "Retrieve a Permission",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.GetRolePermission(apiClient.GetConfig().Context, GetRolePermissionroleIdOrLabel, GetRolePermissionpermissionType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRolePermissionQuiet {
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
				dirPath := filepath.Join(GetRolePermissionBackupDir, "role", "getPermission")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetRolePermissionroleIdOrLabel
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRolePermissionQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRolePermissionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&GetRolePermissionpermissionType, "permissionType", "", "", "")
	cmd.MarkFlagRequired("permissionType")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Role to a file")

	cmd.Flags().StringVarP(&GetRolePermissionBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRolePermissionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRolePermissionCmd := NewGetRolePermissionCmd()
	RoleCmd.AddCommand(GetRolePermissionCmd)
}

var (
	ReplaceRolePermissionroleIdOrLabel string

	ReplaceRolePermissionpermissionType string

	ReplaceRolePermissiondata string

	ReplaceRolePermissionQuiet bool
)

func NewReplaceRolePermissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePermission",
		Long: "Replace a Permission",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.ReplaceRolePermission(apiClient.GetConfig().Context, ReplaceRolePermissionroleIdOrLabel, ReplaceRolePermissionpermissionType)

			if ReplaceRolePermissiondata != "" {
				req = req.Data(ReplaceRolePermissiondata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceRolePermissionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceRolePermissionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&ReplaceRolePermissionpermissionType, "permissionType", "", "", "")
	cmd.MarkFlagRequired("permissionType")

	cmd.Flags().StringVarP(&ReplaceRolePermissiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceRolePermissionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceRolePermissionCmd := NewReplaceRolePermissionCmd()
	RoleCmd.AddCommand(ReplaceRolePermissionCmd)
}

var (
	DeleteRolePermissionroleIdOrLabel string

	DeleteRolePermissionpermissionType string

	DeleteRolePermissionQuiet bool
)

func NewDeleteRolePermissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletePermission",
		Long: "Delete a Permission",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAPI.DeleteRolePermission(apiClient.GetConfig().Context, DeleteRolePermissionroleIdOrLabel, DeleteRolePermissionpermissionType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteRolePermissionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteRolePermissionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&DeleteRolePermissionpermissionType, "permissionType", "", "", "")
	cmd.MarkFlagRequired("permissionType")

	cmd.Flags().BoolVarP(&DeleteRolePermissionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteRolePermissionCmd := NewDeleteRolePermissionCmd()
	RoleCmd.AddCommand(DeleteRolePermissionCmd)
}
