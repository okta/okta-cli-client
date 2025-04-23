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

var RoleTargetCmd = &cobra.Command{
	Use:  "roleTarget",
	Long: "Manage RoleTargetAPI",
}

func init() {
	rootCmd.AddCommand(RoleTargetCmd)
}

var (
	ListApplicationTargetsForApplicationAdministratorRoleForGroupgroupId string

	ListApplicationTargetsForApplicationAdministratorRoleForGrouproleId string

	ListApplicationTargetsForApplicationAdministratorRoleForGroupBackupDir string

	ListApplicationTargetsForApplicationAdministratorRoleForGroupLimit    int32
	ListApplicationTargetsForApplicationAdministratorRoleForGroupPage     string
	ListApplicationTargetsForApplicationAdministratorRoleForGroupFetchAll bool

	ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet bool
)

func NewListApplicationTargetsForApplicationAdministratorRoleForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApplicationTargetsForApplicationAdministratorRoleForGroup",
		Long: "List all Application Targets for an Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.ListApplicationTargetsForApplicationAdministratorRoleForGroup(apiClient.GetConfig().Context, ListApplicationTargetsForApplicationAdministratorRoleForGroupgroupId, ListApplicationTargetsForApplicationAdministratorRoleForGrouproleId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
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
					if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApplicationTargetsForApplicationAdministratorRoleForGroupFetchAll || len(items) == 0 {
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

			if ListApplicationTargetsForApplicationAdministratorRoleForGroupFetchAll && pageCount > 1 && !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApplicationTargetsForApplicationAdministratorRoleForGroupBackupDir, "roletarget", "listApplicationTargetsForApplicationAdministratorRoleForGroup")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
					fmt.Printf("Backing up RoleTargets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleTargets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForGrouproleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().Int32VarP(&ListApplicationTargetsForApplicationAdministratorRoleForGroupLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForGroupPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApplicationTargetsForApplicationAdministratorRoleForGroupFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleTargets to a directory")

	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForGroupBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApplicationTargetsForApplicationAdministratorRoleForGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApplicationTargetsForApplicationAdministratorRoleForGroupCmd := NewListApplicationTargetsForApplicationAdministratorRoleForGroupCmd()
	RoleTargetCmd.AddCommand(ListApplicationTargetsForApplicationAdministratorRoleForGroupCmd)
}

var (
	AssignAppTargetToAdminRoleForGroupgroupId string

	AssignAppTargetToAdminRoleForGrouproleId string

	AssignAppTargetToAdminRoleForGroupappName string

	AssignAppTargetToAdminRoleForGroupQuiet bool
)

func NewAssignAppTargetToAdminRoleForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignAppTargetToAdminRoleForGroup",
		Long: "Assign an Application Target to Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignAppTargetToAdminRoleForGroup(apiClient.GetConfig().Context, AssignAppTargetToAdminRoleForGroupgroupId, AssignAppTargetToAdminRoleForGrouproleId, AssignAppTargetToAdminRoleForGroupappName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignAppTargetToAdminRoleForGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignAppTargetToAdminRoleForGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignAppTargetToAdminRoleForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignAppTargetToAdminRoleForGrouproleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&AssignAppTargetToAdminRoleForGroupappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().BoolVarP(&AssignAppTargetToAdminRoleForGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignAppTargetToAdminRoleForGroupCmd := NewAssignAppTargetToAdminRoleForGroupCmd()
	RoleTargetCmd.AddCommand(AssignAppTargetToAdminRoleForGroupCmd)
}

var (
	UnassignAppTargetToAdminRoleForGroupgroupId string

	UnassignAppTargetToAdminRoleForGrouproleId string

	UnassignAppTargetToAdminRoleForGroupappName string

	UnassignAppTargetToAdminRoleForGroupQuiet bool
)

func NewUnassignAppTargetToAdminRoleForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignAppTargetToAdminRoleForGroup",
		Long: "Unassign an Application Target from Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.UnassignAppTargetToAdminRoleForGroup(apiClient.GetConfig().Context, UnassignAppTargetToAdminRoleForGroupgroupId, UnassignAppTargetToAdminRoleForGrouproleId, UnassignAppTargetToAdminRoleForGroupappName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignAppTargetToAdminRoleForGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignAppTargetToAdminRoleForGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignAppTargetToAdminRoleForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&UnassignAppTargetToAdminRoleForGrouproleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&UnassignAppTargetToAdminRoleForGroupappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().BoolVarP(&UnassignAppTargetToAdminRoleForGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignAppTargetToAdminRoleForGroupCmd := NewUnassignAppTargetToAdminRoleForGroupCmd()
	RoleTargetCmd.AddCommand(UnassignAppTargetToAdminRoleForGroupCmd)
}

var (
	AssignAppInstanceTargetToAppAdminRoleForGroupgroupId string

	AssignAppInstanceTargetToAppAdminRoleForGrouproleId string

	AssignAppInstanceTargetToAppAdminRoleForGroupappName string

	AssignAppInstanceTargetToAppAdminRoleForGroupappId string

	AssignAppInstanceTargetToAppAdminRoleForGroupQuiet bool
)

func NewAssignAppInstanceTargetToAppAdminRoleForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignAppInstanceTargetToAppAdminRoleForGroup",
		Long: "Assign an Application Instance Target to Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignAppInstanceTargetToAppAdminRoleForGroup(apiClient.GetConfig().Context, AssignAppInstanceTargetToAppAdminRoleForGroupgroupId, AssignAppInstanceTargetToAppAdminRoleForGrouproleId, AssignAppInstanceTargetToAppAdminRoleForGroupappName, AssignAppInstanceTargetToAppAdminRoleForGroupappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignAppInstanceTargetToAppAdminRoleForGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignAppInstanceTargetToAppAdminRoleForGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForGrouproleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForGroupappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForGroupappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&AssignAppInstanceTargetToAppAdminRoleForGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignAppInstanceTargetToAppAdminRoleForGroupCmd := NewAssignAppInstanceTargetToAppAdminRoleForGroupCmd()
	RoleTargetCmd.AddCommand(AssignAppInstanceTargetToAppAdminRoleForGroupCmd)
}

var (
	UnassignAppInstanceTargetToAppAdminRoleForGroupgroupId string

	UnassignAppInstanceTargetToAppAdminRoleForGrouproleId string

	UnassignAppInstanceTargetToAppAdminRoleForGroupappName string

	UnassignAppInstanceTargetToAppAdminRoleForGroupappId string

	UnassignAppInstanceTargetToAppAdminRoleForGroupQuiet bool
)

func NewUnassignAppInstanceTargetToAppAdminRoleForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignAppInstanceTargetToAppAdminRoleForGroup",
		Long: "Unassign an Application Instance Target from an Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.UnassignAppInstanceTargetToAppAdminRoleForGroup(apiClient.GetConfig().Context, UnassignAppInstanceTargetToAppAdminRoleForGroupgroupId, UnassignAppInstanceTargetToAppAdminRoleForGrouproleId, UnassignAppInstanceTargetToAppAdminRoleForGroupappName, UnassignAppInstanceTargetToAppAdminRoleForGroupappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignAppInstanceTargetToAppAdminRoleForGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignAppInstanceTargetToAppAdminRoleForGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetToAppAdminRoleForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetToAppAdminRoleForGrouproleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetToAppAdminRoleForGroupappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetToAppAdminRoleForGroupappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&UnassignAppInstanceTargetToAppAdminRoleForGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignAppInstanceTargetToAppAdminRoleForGroupCmd := NewUnassignAppInstanceTargetToAppAdminRoleForGroupCmd()
	RoleTargetCmd.AddCommand(UnassignAppInstanceTargetToAppAdminRoleForGroupCmd)
}

var (
	ListGroupTargetsForGroupRolegroupId string

	ListGroupTargetsForGroupRoleroleId string

	ListGroupTargetsForGroupRoleBackupDir string

	ListGroupTargetsForGroupRoleLimit    int32
	ListGroupTargetsForGroupRolePage     string
	ListGroupTargetsForGroupRoleFetchAll bool

	ListGroupTargetsForGroupRoleQuiet bool
)

func NewListGroupTargetsForGroupRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGroupTargetsForGroupRole",
		Long: "List all Group Targets for a Group Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.ListGroupTargetsForGroupRole(apiClient.GetConfig().Context, ListGroupTargetsForGroupRolegroupId, ListGroupTargetsForGroupRoleroleId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupTargetsForGroupRoleQuiet {
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
					if !ListGroupTargetsForGroupRoleQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupTargetsForGroupRoleFetchAll || len(items) == 0 {
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

			if ListGroupTargetsForGroupRoleFetchAll && pageCount > 1 && !ListGroupTargetsForGroupRoleQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupTargetsForGroupRoleBackupDir, "roletarget", "listGroupTargetsForGroupRole")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupTargetsForGroupRoleQuiet {
					fmt.Printf("Backing up RoleTargets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupTargetsForGroupRoleQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupTargetsForGroupRoleQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupTargetsForGroupRoleQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupTargetsForGroupRoleQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleTargets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupTargetsForGroupRoleQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListGroupTargetsForGroupRolegroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&ListGroupTargetsForGroupRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().Int32VarP(&ListGroupTargetsForGroupRoleLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupTargetsForGroupRolePage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupTargetsForGroupRoleFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleTargets to a directory")

	cmd.Flags().StringVarP(&ListGroupTargetsForGroupRoleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupTargetsForGroupRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupTargetsForGroupRoleCmd := NewListGroupTargetsForGroupRoleCmd()
	RoleTargetCmd.AddCommand(ListGroupTargetsForGroupRoleCmd)
}

var (
	AssignGroupTargetToGroupAdminRolegroupId string

	AssignGroupTargetToGroupAdminRoleroleId string

	AssignGroupTargetToGroupAdminRoletargetGroupId string

	AssignGroupTargetToGroupAdminRoleQuiet bool
)

func NewAssignGroupTargetToGroupAdminRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignGroupTargetToGroupAdminRole",
		Long: "Assign a Group Target to a Group Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignGroupTargetToGroupAdminRole(apiClient.GetConfig().Context, AssignGroupTargetToGroupAdminRolegroupId, AssignGroupTargetToGroupAdminRoleroleId, AssignGroupTargetToGroupAdminRoletargetGroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignGroupTargetToGroupAdminRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignGroupTargetToGroupAdminRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignGroupTargetToGroupAdminRolegroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignGroupTargetToGroupAdminRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&AssignGroupTargetToGroupAdminRoletargetGroupId, "targetGroupId", "", "", "")
	cmd.MarkFlagRequired("targetGroupId")

	cmd.Flags().BoolVarP(&AssignGroupTargetToGroupAdminRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignGroupTargetToGroupAdminRoleCmd := NewAssignGroupTargetToGroupAdminRoleCmd()
	RoleTargetCmd.AddCommand(AssignGroupTargetToGroupAdminRoleCmd)
}

var (
	UnassignGroupTargetFromGroupAdminRolegroupId string

	UnassignGroupTargetFromGroupAdminRoleroleId string

	UnassignGroupTargetFromGroupAdminRoletargetGroupId string

	UnassignGroupTargetFromGroupAdminRoleQuiet bool
)

func NewUnassignGroupTargetFromGroupAdminRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignGroupTargetFromGroupAdminRole",
		Long: "Unassign a Group Target from a Group Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.UnassignGroupTargetFromGroupAdminRole(apiClient.GetConfig().Context, UnassignGroupTargetFromGroupAdminRolegroupId, UnassignGroupTargetFromGroupAdminRoleroleId, UnassignGroupTargetFromGroupAdminRoletargetGroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignGroupTargetFromGroupAdminRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignGroupTargetFromGroupAdminRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignGroupTargetFromGroupAdminRolegroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&UnassignGroupTargetFromGroupAdminRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&UnassignGroupTargetFromGroupAdminRoletargetGroupId, "targetGroupId", "", "", "")
	cmd.MarkFlagRequired("targetGroupId")

	cmd.Flags().BoolVarP(&UnassignGroupTargetFromGroupAdminRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignGroupTargetFromGroupAdminRoleCmd := NewUnassignGroupTargetFromGroupAdminRoleCmd()
	RoleTargetCmd.AddCommand(UnassignGroupTargetFromGroupAdminRoleCmd)
}

var (
	ListApplicationTargetsForApplicationAdministratorRoleForUseruserId string

	ListApplicationTargetsForApplicationAdministratorRoleForUserroleId string

	ListApplicationTargetsForApplicationAdministratorRoleForUserBackupDir string

	ListApplicationTargetsForApplicationAdministratorRoleForUserLimit    int32
	ListApplicationTargetsForApplicationAdministratorRoleForUserPage     string
	ListApplicationTargetsForApplicationAdministratorRoleForUserFetchAll bool

	ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet bool
)

func NewListApplicationTargetsForApplicationAdministratorRoleForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApplicationTargetsForApplicationAdministratorRoleForUser",
		Long: "List all Application Targets for Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.ListApplicationTargetsForApplicationAdministratorRoleForUser(apiClient.GetConfig().Context, ListApplicationTargetsForApplicationAdministratorRoleForUseruserId, ListApplicationTargetsForApplicationAdministratorRoleForUserroleId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
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
					if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApplicationTargetsForApplicationAdministratorRoleForUserFetchAll || len(items) == 0 {
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

			if ListApplicationTargetsForApplicationAdministratorRoleForUserFetchAll && pageCount > 1 && !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApplicationTargetsForApplicationAdministratorRoleForUserBackupDir, "roletarget", "listApplicationTargetsForApplicationAdministratorRoleForUser")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
					fmt.Printf("Backing up RoleTargets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleTargets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().Int32VarP(&ListApplicationTargetsForApplicationAdministratorRoleForUserLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForUserPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApplicationTargetsForApplicationAdministratorRoleForUserFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleTargets to a directory")

	cmd.Flags().StringVarP(&ListApplicationTargetsForApplicationAdministratorRoleForUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApplicationTargetsForApplicationAdministratorRoleForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApplicationTargetsForApplicationAdministratorRoleForUserCmd := NewListApplicationTargetsForApplicationAdministratorRoleForUserCmd()
	RoleTargetCmd.AddCommand(ListApplicationTargetsForApplicationAdministratorRoleForUserCmd)
}

var (
	AssignAllAppsAsTargetToRoleForUseruserId string

	AssignAllAppsAsTargetToRoleForUserroleId string

	AssignAllAppsAsTargetToRoleForUserQuiet bool
)

func NewAssignAllAppsAsTargetToRoleForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignAllAppsAsTargetToRoleForUser",
		Long: "Assign all Apps as Target to Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignAllAppsAsTargetToRoleForUser(apiClient.GetConfig().Context, AssignAllAppsAsTargetToRoleForUseruserId, AssignAllAppsAsTargetToRoleForUserroleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignAllAppsAsTargetToRoleForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignAllAppsAsTargetToRoleForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignAllAppsAsTargetToRoleForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&AssignAllAppsAsTargetToRoleForUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().BoolVarP(&AssignAllAppsAsTargetToRoleForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignAllAppsAsTargetToRoleForUserCmd := NewAssignAllAppsAsTargetToRoleForUserCmd()
	RoleTargetCmd.AddCommand(AssignAllAppsAsTargetToRoleForUserCmd)
}

var (
	AssignAppTargetToAdminRoleForUseruserId string

	AssignAppTargetToAdminRoleForUserroleId string

	AssignAppTargetToAdminRoleForUserappName string

	AssignAppTargetToAdminRoleForUserQuiet bool
)

func NewAssignAppTargetToAdminRoleForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignAppTargetToAdminRoleForUser",
		Long: "Assign an Application Target to Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignAppTargetToAdminRoleForUser(apiClient.GetConfig().Context, AssignAppTargetToAdminRoleForUseruserId, AssignAppTargetToAdminRoleForUserroleId, AssignAppTargetToAdminRoleForUserappName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignAppTargetToAdminRoleForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignAppTargetToAdminRoleForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignAppTargetToAdminRoleForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&AssignAppTargetToAdminRoleForUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&AssignAppTargetToAdminRoleForUserappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().BoolVarP(&AssignAppTargetToAdminRoleForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignAppTargetToAdminRoleForUserCmd := NewAssignAppTargetToAdminRoleForUserCmd()
	RoleTargetCmd.AddCommand(AssignAppTargetToAdminRoleForUserCmd)
}

var (
	UnassignAppTargetFromAppAdminRoleForUseruserId string

	UnassignAppTargetFromAppAdminRoleForUserroleId string

	UnassignAppTargetFromAppAdminRoleForUserappName string

	UnassignAppTargetFromAppAdminRoleForUserQuiet bool
)

func NewUnassignAppTargetFromAppAdminRoleForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignAppTargetFromAppAdminRoleForUser",
		Long: "Unassign an Application Target from an Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.UnassignAppTargetFromAppAdminRoleForUser(apiClient.GetConfig().Context, UnassignAppTargetFromAppAdminRoleForUseruserId, UnassignAppTargetFromAppAdminRoleForUserroleId, UnassignAppTargetFromAppAdminRoleForUserappName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignAppTargetFromAppAdminRoleForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignAppTargetFromAppAdminRoleForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignAppTargetFromAppAdminRoleForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UnassignAppTargetFromAppAdminRoleForUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&UnassignAppTargetFromAppAdminRoleForUserappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().BoolVarP(&UnassignAppTargetFromAppAdminRoleForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignAppTargetFromAppAdminRoleForUserCmd := NewUnassignAppTargetFromAppAdminRoleForUserCmd()
	RoleTargetCmd.AddCommand(UnassignAppTargetFromAppAdminRoleForUserCmd)
}

var (
	AssignAppInstanceTargetToAppAdminRoleForUseruserId string

	AssignAppInstanceTargetToAppAdminRoleForUserroleId string

	AssignAppInstanceTargetToAppAdminRoleForUserappName string

	AssignAppInstanceTargetToAppAdminRoleForUserappId string

	AssignAppInstanceTargetToAppAdminRoleForUserQuiet bool
)

func NewAssignAppInstanceTargetToAppAdminRoleForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignAppInstanceTargetToAppAdminRoleForUser",
		Long: "Assign an Application Instance Target to an Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignAppInstanceTargetToAppAdminRoleForUser(apiClient.GetConfig().Context, AssignAppInstanceTargetToAppAdminRoleForUseruserId, AssignAppInstanceTargetToAppAdminRoleForUserroleId, AssignAppInstanceTargetToAppAdminRoleForUserappName, AssignAppInstanceTargetToAppAdminRoleForUserappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignAppInstanceTargetToAppAdminRoleForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignAppInstanceTargetToAppAdminRoleForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForUserappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&AssignAppInstanceTargetToAppAdminRoleForUserappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&AssignAppInstanceTargetToAppAdminRoleForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignAppInstanceTargetToAppAdminRoleForUserCmd := NewAssignAppInstanceTargetToAppAdminRoleForUserCmd()
	RoleTargetCmd.AddCommand(AssignAppInstanceTargetToAppAdminRoleForUserCmd)
}

var (
	UnassignAppInstanceTargetFromAdminRoleForUseruserId string

	UnassignAppInstanceTargetFromAdminRoleForUserroleId string

	UnassignAppInstanceTargetFromAdminRoleForUserappName string

	UnassignAppInstanceTargetFromAdminRoleForUserappId string

	UnassignAppInstanceTargetFromAdminRoleForUserQuiet bool
)

func NewUnassignAppInstanceTargetFromAdminRoleForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignAppInstanceTargetFromAdminRoleForUser",
		Long: "Unassign an Application Instance Target from an Application Administrator Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.UnassignAppInstanceTargetFromAdminRoleForUser(apiClient.GetConfig().Context, UnassignAppInstanceTargetFromAdminRoleForUseruserId, UnassignAppInstanceTargetFromAdminRoleForUserroleId, UnassignAppInstanceTargetFromAdminRoleForUserappName, UnassignAppInstanceTargetFromAdminRoleForUserappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignAppInstanceTargetFromAdminRoleForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignAppInstanceTargetFromAdminRoleForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetFromAdminRoleForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetFromAdminRoleForUserroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetFromAdminRoleForUserappName, "appName", "", "", "")
	cmd.MarkFlagRequired("appName")

	cmd.Flags().StringVarP(&UnassignAppInstanceTargetFromAdminRoleForUserappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolVarP(&UnassignAppInstanceTargetFromAdminRoleForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignAppInstanceTargetFromAdminRoleForUserCmd := NewUnassignAppInstanceTargetFromAdminRoleForUserCmd()
	RoleTargetCmd.AddCommand(UnassignAppInstanceTargetFromAdminRoleForUserCmd)
}

var (
	ListGroupTargetsForRoleuserId string

	ListGroupTargetsForRoleroleId string

	ListGroupTargetsForRoleBackupDir string

	ListGroupTargetsForRoleLimit    int32
	ListGroupTargetsForRolePage     string
	ListGroupTargetsForRoleFetchAll bool

	ListGroupTargetsForRoleQuiet bool
)

func NewListGroupTargetsForRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGroupTargetsForRole",
		Long: "List all Group Targets for Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.ListGroupTargetsForRole(apiClient.GetConfig().Context, ListGroupTargetsForRoleuserId, ListGroupTargetsForRoleroleId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupTargetsForRoleQuiet {
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
					if !ListGroupTargetsForRoleQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupTargetsForRoleFetchAll || len(items) == 0 {
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

			if ListGroupTargetsForRoleFetchAll && pageCount > 1 && !ListGroupTargetsForRoleQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupTargetsForRoleBackupDir, "roletarget", "listGroupTargetsForRole")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupTargetsForRoleQuiet {
					fmt.Printf("Backing up RoleTargets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupTargetsForRoleQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupTargetsForRoleQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupTargetsForRoleQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupTargetsForRoleQuiet {
					fmt.Printf("Successfully backed up %d/%d RoleTargets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupTargetsForRoleQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListGroupTargetsForRoleuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListGroupTargetsForRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().Int32VarP(&ListGroupTargetsForRoleLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupTargetsForRolePage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupTargetsForRoleFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple RoleTargets to a directory")

	cmd.Flags().StringVarP(&ListGroupTargetsForRoleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupTargetsForRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupTargetsForRoleCmd := NewListGroupTargetsForRoleCmd()
	RoleTargetCmd.AddCommand(ListGroupTargetsForRoleCmd)
}

var (
	AssignGroupTargetToUserRoleuserId string

	AssignGroupTargetToUserRoleroleId string

	AssignGroupTargetToUserRolegroupId string

	AssignGroupTargetToUserRoleQuiet bool
)

func NewAssignGroupTargetToUserRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignGroupTargetToUserRole",
		Long: "Assign a Group Target to Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.AssignGroupTargetToUserRole(apiClient.GetConfig().Context, AssignGroupTargetToUserRoleuserId, AssignGroupTargetToUserRoleroleId, AssignGroupTargetToUserRolegroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignGroupTargetToUserRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignGroupTargetToUserRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignGroupTargetToUserRoleuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&AssignGroupTargetToUserRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&AssignGroupTargetToUserRolegroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().BoolVarP(&AssignGroupTargetToUserRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignGroupTargetToUserRoleCmd := NewAssignGroupTargetToUserRoleCmd()
	RoleTargetCmd.AddCommand(AssignGroupTargetToUserRoleCmd)
}

var (
	UnassignGroupTargetFromUserAdminRoleuserId string

	UnassignGroupTargetFromUserAdminRoleroleId string

	UnassignGroupTargetFromUserAdminRolegroupId string

	UnassignGroupTargetFromUserAdminRoleQuiet bool
)

func NewUnassignGroupTargetFromUserAdminRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignGroupTargetFromUserAdminRole",
		Long: "Unassign a Group Target from Role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RoleTargetAPI.UnassignGroupTargetFromUserAdminRole(apiClient.GetConfig().Context, UnassignGroupTargetFromUserAdminRoleuserId, UnassignGroupTargetFromUserAdminRoleroleId, UnassignGroupTargetFromUserAdminRolegroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignGroupTargetFromUserAdminRoleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignGroupTargetFromUserAdminRoleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignGroupTargetFromUserAdminRoleuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UnassignGroupTargetFromUserAdminRoleroleId, "roleId", "", "", "")
	cmd.MarkFlagRequired("roleId")

	cmd.Flags().StringVarP(&UnassignGroupTargetFromUserAdminRolegroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().BoolVarP(&UnassignGroupTargetFromUserAdminRoleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignGroupTargetFromUserAdminRoleCmd := NewUnassignGroupTargetFromUserAdminRoleCmd()
	RoleTargetCmd.AddCommand(UnassignGroupTargetFromUserAdminRoleCmd)
}
