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

var GroupOwnerCmd = &cobra.Command{
	Use:  "groupOwner",
	Long: "Manage GroupOwnerAPI",
}

func init() {
	rootCmd.AddCommand(GroupOwnerCmd)
}

var (
	AssignGroupOwnergroupId string

	AssignGroupOwnerdata string

	AssignGroupOwnerQuiet bool
)

func NewAssignGroupOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assign",
		Long: "Assign a Group Owner",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupOwnerAPI.AssignGroupOwner(apiClient.GetConfig().Context, AssignGroupOwnergroupId)

			if AssignGroupOwnerdata != "" {
				req = req.Data(AssignGroupOwnerdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignGroupOwnerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignGroupOwnerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignGroupOwnergroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignGroupOwnerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AssignGroupOwnerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignGroupOwnerCmd := NewAssignGroupOwnerCmd()
	GroupOwnerCmd.AddCommand(AssignGroupOwnerCmd)
}

var (
	ListGroupOwnersgroupId string

	ListGroupOwnersBackupDir string

	ListGroupOwnersLimit    int32
	ListGroupOwnersPage     string
	ListGroupOwnersFetchAll bool

	ListGroupOwnersQuiet bool
)

func NewListGroupOwnersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Group Owners",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupOwnerAPI.ListGroupOwners(apiClient.GetConfig().Context, ListGroupOwnersgroupId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGroupOwnersQuiet {
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
					if !ListGroupOwnersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGroupOwnersFetchAll || len(items) == 0 {
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

			if ListGroupOwnersFetchAll && pageCount > 1 && !ListGroupOwnersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGroupOwnersBackupDir, "groupowner", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGroupOwnersQuiet {
					fmt.Printf("Backing up GroupOwners to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGroupOwnersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGroupOwnersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGroupOwnersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGroupOwnersQuiet {
					fmt.Printf("Successfully backed up %d/%d GroupOwners\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGroupOwnersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListGroupOwnersgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().Int32VarP(&ListGroupOwnersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGroupOwnersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGroupOwnersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple GroupOwners to a directory")

	cmd.Flags().StringVarP(&ListGroupOwnersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGroupOwnersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGroupOwnersCmd := NewListGroupOwnersCmd()
	GroupOwnerCmd.AddCommand(ListGroupOwnersCmd)
}

var (
	DeleteGroupOwnergroupId string

	DeleteGroupOwnerownerId string

	DeleteGroupOwnerQuiet bool
)

func NewDeleteGroupOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Group Owner",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupOwnerAPI.DeleteGroupOwner(apiClient.GetConfig().Context, DeleteGroupOwnergroupId, DeleteGroupOwnerownerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteGroupOwnerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteGroupOwnerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteGroupOwnergroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&DeleteGroupOwnerownerId, "ownerId", "", "", "")
	cmd.MarkFlagRequired("ownerId")

	cmd.Flags().BoolVarP(&DeleteGroupOwnerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteGroupOwnerCmd := NewDeleteGroupOwnerCmd()
	GroupOwnerCmd.AddCommand(DeleteGroupOwnerCmd)
}
