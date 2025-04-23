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

var ApplicationUsersCmd = &cobra.Command{
	Use:  "applicationUsers",
	Long: "Manage ApplicationUsersAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationUsersCmd)
}

var (
	AssignUserToApplicationappId string

	AssignUserToApplicationdata string

	AssignUserToApplicationQuiet bool
)

func NewAssignUserToApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignUserToApplication",
		Long: "Assign an Application User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.AssignUserToApplication(apiClient.GetConfig().Context, AssignUserToApplicationappId)

			if AssignUserToApplicationdata != "" {
				req = req.Data(AssignUserToApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignUserToApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignUserToApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignUserToApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&AssignUserToApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AssignUserToApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignUserToApplicationCmd := NewAssignUserToApplicationCmd()
	ApplicationUsersCmd.AddCommand(AssignUserToApplicationCmd)
}

var (
	ListApplicationUsersappId string

	ListApplicationUsersBackupDir string

	ListApplicationUsersLimit    int32
	ListApplicationUsersPage     string
	ListApplicationUsersFetchAll bool

	ListApplicationUsersQuiet bool
)

func NewListApplicationUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List all Application Users",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.ListApplicationUsers(apiClient.GetConfig().Context, ListApplicationUsersappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListApplicationUsersQuiet {
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
					if !ListApplicationUsersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListApplicationUsersFetchAll || len(items) == 0 {
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

			if ListApplicationUsersFetchAll && pageCount > 1 && !ListApplicationUsersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListApplicationUsersBackupDir, "applicationusers", "list")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListApplicationUsersQuiet {
					fmt.Printf("Backing up ApplicationUserss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListApplicationUsersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListApplicationUsersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListApplicationUsersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListApplicationUsersQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationUserss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListApplicationUsersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListApplicationUsersappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListApplicationUsersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListApplicationUsersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListApplicationUsersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationUserss to a directory")

	cmd.Flags().StringVarP(&ListApplicationUsersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListApplicationUsersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListApplicationUsersCmd := NewListApplicationUsersCmd()
	ApplicationUsersCmd.AddCommand(ListApplicationUsersCmd)
}

var (
	UpdateApplicationUserappId string

	UpdateApplicationUseruserId string

	UpdateApplicationUserdata string

	UpdateApplicationUserQuiet bool
)

func NewUpdateApplicationUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateApplicationUser",
		Long: "Update an Application User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.UpdateApplicationUser(apiClient.GetConfig().Context, UpdateApplicationUserappId, UpdateApplicationUseruserId)

			if UpdateApplicationUserdata != "" {
				req = req.Data(UpdateApplicationUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateApplicationUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateApplicationUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateApplicationUserappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateApplicationUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UpdateApplicationUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateApplicationUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateApplicationUserCmd := NewUpdateApplicationUserCmd()
	ApplicationUsersCmd.AddCommand(UpdateApplicationUserCmd)
}

var (
	GetApplicationUserappId string

	GetApplicationUseruserId string

	GetApplicationUserBackupDir string

	GetApplicationUserQuiet bool
)

func NewGetApplicationUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationUser",
		Long: "Retrieve an Application User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.GetApplicationUser(apiClient.GetConfig().Context, GetApplicationUserappId, GetApplicationUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApplicationUserQuiet {
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
				dirPath := filepath.Join(GetApplicationUserBackupDir, "applicationusers", "getApplicationUser")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApplicationUserappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApplicationUserQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApplicationUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApplicationUserappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetApplicationUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationUsers to a file")

	cmd.Flags().StringVarP(&GetApplicationUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApplicationUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApplicationUserCmd := NewGetApplicationUserCmd()
	ApplicationUsersCmd.AddCommand(GetApplicationUserCmd)
}

var (
	UnassignUserFromApplicationappId string

	UnassignUserFromApplicationuserId string

	UnassignUserFromApplicationQuiet bool
)

func NewUnassignUserFromApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignUserFromApplication",
		Long: "Unassign an Application User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.UnassignUserFromApplication(apiClient.GetConfig().Context, UnassignUserFromApplicationappId, UnassignUserFromApplicationuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignUserFromApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignUserFromApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignUserFromApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UnassignUserFromApplicationuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&UnassignUserFromApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignUserFromApplicationCmd := NewUnassignUserFromApplicationCmd()
	ApplicationUsersCmd.AddCommand(UnassignUserFromApplicationCmd)
}
