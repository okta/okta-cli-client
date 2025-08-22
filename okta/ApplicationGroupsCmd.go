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

var ApplicationGroupsCmd = &cobra.Command{
	Use:  "applicationGroups",
	Long: "Manage ApplicationGroupsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationGroupsCmd)
}

var (
	ListApplicationGroupAssignmentsappId string

	ListApplicationGroupAssignmentsBackupDir string

	ListApplicationGroupAssignmentsLimit    int32
	ListApplicationGroupAssignmentsPage     string
	ListApplicationGroupAssignmentsFetchAll bool

	ListApplicationGroupAssignmentsQuiet bool
)

func NewListApplicationGroupAssignmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApplicationGroupAssignments",
		Long: "List all Assigned Groups",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.ListApplicationGroupAssignments(apiClient.GetConfig().Context, ListApplicationGroupAssignmentsappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApplicationGroupAssignmentsQuiet {
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
					if !ListApplicationGroupAssignmentsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApplicationGroupAssignmentsFetchAll || len(items) == 0 {
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

			if ListApplicationGroupAssignmentsFetchAll && pageCount > 1 && !ListApplicationGroupAssignmentsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApplicationGroupAssignmentsBackupDir, "applicationgroups", "listApplicationGroupAssignments")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApplicationGroupAssignmentsQuiet {
					fmt.Printf("Backing up ApplicationGroupss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApplicationGroupAssignmentsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApplicationGroupAssignmentsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApplicationGroupAssignmentsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApplicationGroupAssignmentsQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationGroupss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApplicationGroupAssignmentsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListApplicationGroupAssignmentsappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListApplicationGroupAssignmentsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApplicationGroupAssignmentsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApplicationGroupAssignmentsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationGroupss to a directory")

	cmd.Flags().StringVarP(&ListApplicationGroupAssignmentsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApplicationGroupAssignmentsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApplicationGroupAssignmentsCmd := NewListApplicationGroupAssignmentsCmd()
	ApplicationGroupsCmd.AddCommand(ListApplicationGroupAssignmentsCmd)
}

var (
	GetApplicationGroupAssignmentappId string

	GetApplicationGroupAssignmentgroupId string

	GetApplicationGroupAssignmentBackupDir string

	GetApplicationGroupAssignmentQuiet bool
)

func NewGetApplicationGroupAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationGroupAssignment",
		Long: "Retrieve an Assigned Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.GetApplicationGroupAssignment(apiClient.GetConfig().Context, GetApplicationGroupAssignmentappId, GetApplicationGroupAssignmentgroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApplicationGroupAssignmentQuiet {
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
				dirPath := filepath.Join(GetApplicationGroupAssignmentBackupDir, "applicationgroups", "getApplicationGroupAssignment")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApplicationGroupAssignmentappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApplicationGroupAssignmentQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApplicationGroupAssignmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApplicationGroupAssignmentappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetApplicationGroupAssignmentgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationGroups to a file")

	cmd.Flags().StringVarP(&GetApplicationGroupAssignmentBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApplicationGroupAssignmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApplicationGroupAssignmentCmd := NewGetApplicationGroupAssignmentCmd()
	ApplicationGroupsCmd.AddCommand(GetApplicationGroupAssignmentCmd)
}

var (
	AssignGroupToApplicationappId string

	AssignGroupToApplicationgroupId string

	AssignGroupToApplicationdata string

	AssignGroupToApplicationQuiet bool
)

func NewAssignGroupToApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignGroupToApplication",
		Long: "Assign a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.AssignGroupToApplication(apiClient.GetConfig().Context, AssignGroupToApplicationappId, AssignGroupToApplicationgroupId)

			if AssignGroupToApplicationdata != "" {
				req = req.Data(AssignGroupToApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignGroupToApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignGroupToApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignGroupToApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&AssignGroupToApplicationgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignGroupToApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AssignGroupToApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignGroupToApplicationCmd := NewAssignGroupToApplicationCmd()
	ApplicationGroupsCmd.AddCommand(AssignGroupToApplicationCmd)
}

var (
	UnassignApplicationFromGroupappId string

	UnassignApplicationFromGroupgroupId string

	UnassignApplicationFromGroupQuiet bool
)

func NewUnassignApplicationFromGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignApplicationFromGroup",
		Long: "Unassign a Group",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.UnassignApplicationFromGroup(apiClient.GetConfig().Context, UnassignApplicationFromGroupappId, UnassignApplicationFromGroupgroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignApplicationFromGroupQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignApplicationFromGroupQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignApplicationFromGroupappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UnassignApplicationFromGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().BoolVarP(&UnassignApplicationFromGroupQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignApplicationFromGroupCmd := NewUnassignApplicationFromGroupCmd()
	ApplicationGroupsCmd.AddCommand(UnassignApplicationFromGroupCmd)
}
