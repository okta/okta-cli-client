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

var RoleAssignmentCmd = &cobra.Command{
	Use:  "roleAssignment",
	Long: "Manage RoleAssignmentAPI",
}

func init() {
	rootCmd.AddCommand(RoleAssignmentCmd)
}

var (
	AssignRoleToGroupgroupId string

	AssignRoleToGroupdata string

	AssignRoleToGroupQuiet bool
)

func NewAssignRoleToGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignRoleToGroup",
		Long: "Assign a Role to a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.AssignRoleToGroup(apiClient.GetConfig().Context, AssignRoleToGroupgroupId)

			if AssignRoleToGroupdata != "" {
				req = req.Data(AssignRoleToGroupdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignRoleToGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignRoleToGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignRoleToGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignRoleToGroupdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AssignRoleToGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignRoleToGroupCmd := NewAssignRoleToGroupCmd()
	RoleAssignmentCmd.AddCommand(AssignRoleToGroupCmd)
}

var (
	ListGroupAssignedRolesgroupId string

	ListGroupAssignedRolesBackupDir string

	ListGroupAssignedRolesLimit    int32
	ListGroupAssignedRolesPage     string
	ListGroupAssignedRolesFetchAll bool

	ListGroupAssignedRolesQuiet bool
)

func NewListGroupAssignedRolesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGroupAssignedRoles",
		Long: "List all Assigned Roles of Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.ListGroupAssignedRoles(apiClient.GetConfig().Context, ListGroupAssignedRolesgroupId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupAssignedRolesQuiet {
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
					if !ListGroupAssignedRolesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupAssignedRolesFetchAll || len(items) == 0 {
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

			if ListGroupAssignedRolesFetchAll && pageCount > 1 && !ListGroupAssignedRolesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupAssignedRolesBackupDir, "roleassignment", "listGroupAssignedRoles")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupAssignedRolesQuiet {
					fmt.Printf("Backing up RoleAssignments to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupAssignedRolesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupAssignedRolesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupAssignedRolesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupAssignedRolesQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleAssignments\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupAssignedRolesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListGroupAssignedRolesgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().Int32VarP(&ListGroupAssignedRolesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupAssignedRolesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupAssignedRolesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleAssignments to a directory")

	cmd.Flags().StringVarP(&ListGroupAssignedRolesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupAssignedRolesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupAssignedRolesCmd := NewListGroupAssignedRolesCmd()
	RoleAssignmentCmd.AddCommand(ListGroupAssignedRolesCmd)
}

var (
	GetGroupAssignedRolegroupId string

	GetGroupAssignedRoleroleId string

	GetGroupAssignedRoleBackupDir string

	GetGroupAssignedRoleQuiet bool
)

func NewGetGroupAssignedRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getGroupAssignedRole",
		Long: "Retrieve a Role assigned to Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.GetGroupAssignedRole(apiClient.GetConfig().Context, GetGroupAssignedRolegroupId, GetGroupAssignedRoleroleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetGroupAssignedRoleQuiet {
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
				dirPath := filepath.Join(GetGroupAssignedRoleBackupDir, "roleassignment", "getGroupAssignedRole")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetGroupAssignedRolegroupId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetGroupAssignedRoleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetGroupAssignedRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetGroupAssignedRolegroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&GetGroupAssignedRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the RoleAssignment to a file")

	cmd.Flags().StringVarP(&GetGroupAssignedRoleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetGroupAssignedRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetGroupAssignedRoleCmd := NewGetGroupAssignedRoleCmd()
	RoleAssignmentCmd.AddCommand(GetGroupAssignedRoleCmd)
}

var (
	UnassignRoleFromGroupgroupId string

	UnassignRoleFromGrouproleId string

	UnassignRoleFromGroupQuiet bool
)

func NewUnassignRoleFromGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignRoleFromGroup",
		Long: "Unassign a Role from a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.UnassignRoleFromGroup(apiClient.GetConfig().Context, UnassignRoleFromGroupgroupId, UnassignRoleFromGrouproleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignRoleFromGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignRoleFromGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignRoleFromGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&UnassignRoleFromGrouproleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().BoolVarP(&UnassignRoleFromGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignRoleFromGroupCmd := NewUnassignRoleFromGroupCmd()
	RoleAssignmentCmd.AddCommand(UnassignRoleFromGroupCmd)
}

var (
	ListUsersWithRoleAssignmentsBackupDir string

	ListUsersWithRoleAssignmentsLimit    int32
	ListUsersWithRoleAssignmentsPage     string
	ListUsersWithRoleAssignmentsFetchAll bool

	ListUsersWithRoleAssignmentsQuiet bool
)

func NewListUsersWithRoleAssignmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listUsersWiths",
		Long: "List all Users with Role Assignments",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.ListUsersWithRoleAssignments(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUsersWithRoleAssignmentsQuiet {
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
					if !ListUsersWithRoleAssignmentsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUsersWithRoleAssignmentsFetchAll || len(items) == 0 {
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

			if ListUsersWithRoleAssignmentsFetchAll && pageCount > 1 && !ListUsersWithRoleAssignmentsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUsersWithRoleAssignmentsBackupDir, "roleassignment", "listUsersWiths")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUsersWithRoleAssignmentsQuiet {
					fmt.Printf("Backing up RoleAssignments to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUsersWithRoleAssignmentsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUsersWithRoleAssignmentsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUsersWithRoleAssignmentsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUsersWithRoleAssignmentsQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleAssignments\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUsersWithRoleAssignmentsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListUsersWithRoleAssignmentsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUsersWithRoleAssignmentsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUsersWithRoleAssignmentsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleAssignments to a directory")

	cmd.Flags().StringVarP(&ListUsersWithRoleAssignmentsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUsersWithRoleAssignmentsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUsersWithRoleAssignmentsCmd := NewListUsersWithRoleAssignmentsCmd()
	RoleAssignmentCmd.AddCommand(ListUsersWithRoleAssignmentsCmd)
}

var (
	AssignRoleToUseruserId string

	AssignRoleToUserdata string

	AssignRoleToUserQuiet bool
)

func NewAssignRoleToUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignRoleToUser",
		Long: "Assign a Role to a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.AssignRoleToUser(apiClient.GetConfig().Context, AssignRoleToUseruserId)

			if AssignRoleToUserdata != "" {
				req = req.Data(AssignRoleToUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignRoleToUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignRoleToUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignRoleToUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&AssignRoleToUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AssignRoleToUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignRoleToUserCmd := NewAssignRoleToUserCmd()
	RoleAssignmentCmd.AddCommand(AssignRoleToUserCmd)
}

var (
	ListAssignedRolesForUseruserId string

	ListAssignedRolesForUserBackupDir string

	ListAssignedRolesForUserLimit    int32
	ListAssignedRolesForUserPage     string
	ListAssignedRolesForUserFetchAll bool

	ListAssignedRolesForUserQuiet bool
)

func NewListAssignedRolesForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAssignedRolesForUser",
		Long: "List all Roles assigned to a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.ListAssignedRolesForUser(apiClient.GetConfig().Context, ListAssignedRolesForUseruserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAssignedRolesForUserQuiet {
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
					if !ListAssignedRolesForUserQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAssignedRolesForUserFetchAll || len(items) == 0 {
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

			if ListAssignedRolesForUserFetchAll && pageCount > 1 && !ListAssignedRolesForUserQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAssignedRolesForUserBackupDir, "roleassignment", "listAssignedRolesForUser")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAssignedRolesForUserQuiet {
					fmt.Printf("Backing up RoleAssignments to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAssignedRolesForUserQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAssignedRolesForUserQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAssignedRolesForUserQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAssignedRolesForUserQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleAssignments\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAssignedRolesForUserQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAssignedRolesForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListAssignedRolesForUserLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAssignedRolesForUserPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAssignedRolesForUserFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleAssignments to a directory")

	cmd.Flags().StringVarP(&ListAssignedRolesForUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAssignedRolesForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAssignedRolesForUserCmd := NewListAssignedRolesForUserCmd()
	RoleAssignmentCmd.AddCommand(ListAssignedRolesForUserCmd)
}

var (
	GetUserAssignedRoleuserId string

	GetUserAssignedRoleroleId string

	GetUserAssignedRoleBackupDir string

	GetUserAssignedRoleQuiet bool
)

func NewGetUserAssignedRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getUserAssignedRole",
		Long: "Retrieve a Role assigned to a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.GetUserAssignedRole(apiClient.GetConfig().Context, GetUserAssignedRoleuserId, GetUserAssignedRoleroleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetUserAssignedRoleQuiet {
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
				dirPath := filepath.Join(GetUserAssignedRoleBackupDir, "roleassignment", "getUserAssignedRole")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetUserAssignedRoleuserId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetUserAssignedRoleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetUserAssignedRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetUserAssignedRoleuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetUserAssignedRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the RoleAssignment to a file")

	cmd.Flags().StringVarP(&GetUserAssignedRoleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetUserAssignedRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetUserAssignedRoleCmd := NewGetUserAssignedRoleCmd()
	RoleAssignmentCmd.AddCommand(GetUserAssignedRoleCmd)
}

var (
	UnassignRoleFromUseruserId string

	UnassignRoleFromUserroleId string

	UnassignRoleFromUserQuiet bool
)

func NewUnassignRoleFromUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignRoleFromUser",
		Long: "Unassign a Role from a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleAssignmentAPI.UnassignRoleFromUser(apiClient.GetConfig().Context, UnassignRoleFromUseruserId, UnassignRoleFromUserroleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignRoleFromUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignRoleFromUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignRoleFromUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UnassignRoleFromUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().BoolVarP(&UnassignRoleFromUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignRoleFromUserCmd := NewUnassignRoleFromUserCmd()
	RoleAssignmentCmd.AddCommand(UnassignRoleFromUserCmd)
}
