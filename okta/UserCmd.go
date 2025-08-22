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

var UserCmd = &cobra.Command{
	Use:  "user",
	Long: "Manage UserAPI",
}

func init() {
	rootCmd.AddCommand(UserCmd)
}

var (
	CreateUserdata string

	CreateUserRestoreFile string

	CreateUserQuiet bool
)

func NewCreateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateUserRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateUserRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateUserdata = string(processedData)

				if !CreateUserQuiet {
					fmt.Println("Restoring User from:", CreateUserRestoreFile)
				}
			}

			req := apiClient.UserAPI.CreateUser(apiClient.GetConfig().Context)

			if CreateUserdata != "" {
				req = req.Data(CreateUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateUserRestoreFile, "restore-from", "r", "", "Restore User from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateUserCmd := NewCreateUserCmd()
	UserCmd.AddCommand(CreateUserCmd)
}

var (
	ListUsersBackupDir string

	ListUsersLimit    int32
	ListUsersPage     string
	ListUsersFetchAll bool

	ListUsersQuiet bool
)

func NewListUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Users",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUsers(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUsersQuiet {
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
					if !ListUsersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUsersFetchAll || len(items) == 0 {
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

			if ListUsersFetchAll && pageCount > 1 && !ListUsersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUsersBackupDir, "user", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUsersQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUsersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUsersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUsersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUsersQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUsersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListUsersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUsersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUsersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListUsersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUsersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUsersCmd := NewListUsersCmd()
	UserCmd.AddCommand(ListUsersCmd)
}

var (
	UpdateUseruserId string

	UpdateUserdata string

	UpdateUserQuiet bool
)

func NewUpdateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: "Update a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.UpdateUser(apiClient.GetConfig().Context, UpdateUseruserId)

			if UpdateUserdata != "" {
				req = req.Data(UpdateUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UpdateUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateUserCmd := NewUpdateUserCmd()
	UserCmd.AddCommand(UpdateUserCmd)
}

var (
	GetUseruserId string

	GetUserBackupDir string

	GetUserQuiet bool
)

func NewGetUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GetUser(apiClient.GetConfig().Context, GetUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetUserQuiet {
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
				dirPath := filepath.Join(GetUserBackupDir, "user", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetUseruserId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetUserQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the User to a file")

	cmd.Flags().StringVarP(&GetUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetUserCmd := NewGetUserCmd()
	UserCmd.AddCommand(GetUserCmd)
}

var (
	ReplaceUseruserId string

	ReplaceUserdata string

	ReplaceUserQuiet bool
)

func NewReplaceUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ReplaceUser(apiClient.GetConfig().Context, ReplaceUseruserId)

			if ReplaceUserdata != "" {
				req = req.Data(ReplaceUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ReplaceUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceUserCmd := NewReplaceUserCmd()
	UserCmd.AddCommand(ReplaceUserCmd)
}

var (
	DeleteUseruserId string

	DeleteUserQuiet bool
)

func NewDeleteUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.DeleteUser(apiClient.GetConfig().Context, DeleteUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&DeleteUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteUserCmd := NewDeleteUserCmd()
	UserCmd.AddCommand(DeleteUserCmd)
}

var (
	ListAppLinksuserId string

	ListAppLinksBackupDir string

	ListAppLinksLimit    int32
	ListAppLinksPage     string
	ListAppLinksFetchAll bool

	ListAppLinksQuiet bool
)

func NewListAppLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAppLinks",
		Long: "List all Assigned Application Links",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListAppLinks(apiClient.GetConfig().Context, ListAppLinksuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAppLinksQuiet {
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
					if !ListAppLinksQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAppLinksFetchAll || len(items) == 0 {
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

			if ListAppLinksFetchAll && pageCount > 1 && !ListAppLinksQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAppLinksBackupDir, "user", "listAppLinks")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAppLinksQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAppLinksQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAppLinksQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAppLinksQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAppLinksQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAppLinksQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAppLinksuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListAppLinksLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAppLinksPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAppLinksFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListAppLinksBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAppLinksQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAppLinksCmd := NewListAppLinksCmd()
	UserCmd.AddCommand(ListAppLinksCmd)
}

var (
	ListUserBlocksuserId string

	ListUserBlocksBackupDir string

	ListUserBlocksLimit    int32
	ListUserBlocksPage     string
	ListUserBlocksFetchAll bool

	ListUserBlocksQuiet bool
)

func NewListUserBlocksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBlocks",
		Long: "List all User Blocks",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserBlocks(apiClient.GetConfig().Context, ListUserBlocksuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUserBlocksQuiet {
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
					if !ListUserBlocksQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUserBlocksFetchAll || len(items) == 0 {
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

			if ListUserBlocksFetchAll && pageCount > 1 && !ListUserBlocksQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUserBlocksBackupDir, "user", "listBlocks")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUserBlocksQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUserBlocksQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUserBlocksQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUserBlocksQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUserBlocksQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUserBlocksQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListUserBlocksuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListUserBlocksLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUserBlocksPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUserBlocksFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListUserBlocksBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUserBlocksQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUserBlocksCmd := NewListUserBlocksCmd()
	UserCmd.AddCommand(ListUserBlocksCmd)
}

var (
	ListUserClientsuserId string

	ListUserClientsBackupDir string

	ListUserClientsLimit    int32
	ListUserClientsPage     string
	ListUserClientsFetchAll bool

	ListUserClientsQuiet bool
)

func NewListUserClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listClients",
		Long: "List all Clients",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserClients(apiClient.GetConfig().Context, ListUserClientsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUserClientsQuiet {
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
					if !ListUserClientsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUserClientsFetchAll || len(items) == 0 {
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

			if ListUserClientsFetchAll && pageCount > 1 && !ListUserClientsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUserClientsBackupDir, "user", "listClients")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUserClientsQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUserClientsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUserClientsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUserClientsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUserClientsQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUserClientsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListUserClientsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListUserClientsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUserClientsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUserClientsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListUserClientsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUserClientsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUserClientsCmd := NewListUserClientsCmd()
	UserCmd.AddCommand(ListUserClientsCmd)
}

var (
	ListGrantsForUserAndClientuserId string

	ListGrantsForUserAndClientclientId string

	ListGrantsForUserAndClientBackupDir string

	ListGrantsForUserAndClientLimit    int32
	ListGrantsForUserAndClientPage     string
	ListGrantsForUserAndClientFetchAll bool

	ListGrantsForUserAndClientQuiet bool
)

func NewListGrantsForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGrantsForAndClient",
		Long: "List all Grants for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListGrantsForUserAndClient(apiClient.GetConfig().Context, ListGrantsForUserAndClientuserId, ListGrantsForUserAndClientclientId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListGrantsForUserAndClientQuiet {
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
					if !ListGrantsForUserAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListGrantsForUserAndClientFetchAll || len(items) == 0 {
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

			if ListGrantsForUserAndClientFetchAll && pageCount > 1 && !ListGrantsForUserAndClientQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListGrantsForUserAndClientBackupDir, "user", "listGrantsForAndClient")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListGrantsForUserAndClientQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListGrantsForUserAndClientQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListGrantsForUserAndClientQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListGrantsForUserAndClientQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListGrantsForUserAndClientQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListGrantsForUserAndClientQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListGrantsForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListGrantsForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().Int32VarP(&ListGrantsForUserAndClientLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListGrantsForUserAndClientPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListGrantsForUserAndClientFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListGrantsForUserAndClientBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListGrantsForUserAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListGrantsForUserAndClientCmd := NewListGrantsForUserAndClientCmd()
	UserCmd.AddCommand(ListGrantsForUserAndClientCmd)
}

var (
	RevokeGrantsForUserAndClientuserId string

	RevokeGrantsForUserAndClientclientId string

	RevokeGrantsForUserAndClientQuiet bool
)

func NewRevokeGrantsForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeGrantsForAndClient",
		Long: "Revoke all Grants for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeGrantsForUserAndClient(apiClient.GetConfig().Context, RevokeGrantsForUserAndClientuserId, RevokeGrantsForUserAndClientclientId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeGrantsForUserAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeGrantsForUserAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeGrantsForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeGrantsForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().BoolVarP(&RevokeGrantsForUserAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeGrantsForUserAndClientCmd := NewRevokeGrantsForUserAndClientCmd()
	UserCmd.AddCommand(RevokeGrantsForUserAndClientCmd)
}

var (
	ListRefreshTokensForUserAndClientuserId string

	ListRefreshTokensForUserAndClientclientId string

	ListRefreshTokensForUserAndClientBackupDir string

	ListRefreshTokensForUserAndClientLimit    int32
	ListRefreshTokensForUserAndClientPage     string
	ListRefreshTokensForUserAndClientFetchAll bool

	ListRefreshTokensForUserAndClientQuiet bool
)

func NewListRefreshTokensForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listRefreshTokensForAndClient",
		Long: "List all Refresh Tokens for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListRefreshTokensForUserAndClient(apiClient.GetConfig().Context, ListRefreshTokensForUserAndClientuserId, ListRefreshTokensForUserAndClientclientId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRefreshTokensForUserAndClientQuiet {
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
					if !ListRefreshTokensForUserAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRefreshTokensForUserAndClientFetchAll || len(items) == 0 {
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

			if ListRefreshTokensForUserAndClientFetchAll && pageCount > 1 && !ListRefreshTokensForUserAndClientQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRefreshTokensForUserAndClientBackupDir, "user", "listRefreshTokensForAndClient")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRefreshTokensForUserAndClientQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRefreshTokensForUserAndClientQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRefreshTokensForUserAndClientQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRefreshTokensForUserAndClientQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRefreshTokensForUserAndClientQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRefreshTokensForUserAndClientQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListRefreshTokensForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListRefreshTokensForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().Int32VarP(&ListRefreshTokensForUserAndClientLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRefreshTokensForUserAndClientPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRefreshTokensForUserAndClientFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListRefreshTokensForUserAndClientBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRefreshTokensForUserAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRefreshTokensForUserAndClientCmd := NewListRefreshTokensForUserAndClientCmd()
	UserCmd.AddCommand(ListRefreshTokensForUserAndClientCmd)
}

var (
	RevokeTokensForUserAndClientuserId string

	RevokeTokensForUserAndClientclientId string

	RevokeTokensForUserAndClientQuiet bool
)

func NewRevokeTokensForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeTokensForAndClient",
		Long: "Revoke all Refresh Tokens for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeTokensForUserAndClient(apiClient.GetConfig().Context, RevokeTokensForUserAndClientuserId, RevokeTokensForUserAndClientclientId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeTokensForUserAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeTokensForUserAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeTokensForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeTokensForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().BoolVarP(&RevokeTokensForUserAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeTokensForUserAndClientCmd := NewRevokeTokensForUserAndClientCmd()
	UserCmd.AddCommand(RevokeTokensForUserAndClientCmd)
}

var (
	GetRefreshTokenForUserAndClientuserId string

	GetRefreshTokenForUserAndClientclientId string

	GetRefreshTokenForUserAndClienttokenId string

	GetRefreshTokenForUserAndClientBackupDir string

	GetRefreshTokenForUserAndClientQuiet bool
)

func NewGetRefreshTokenForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getRefreshTokenForAndClient",
		Long: "Retrieve a Refresh Token for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GetRefreshTokenForUserAndClient(apiClient.GetConfig().Context, GetRefreshTokenForUserAndClientuserId, GetRefreshTokenForUserAndClientclientId, GetRefreshTokenForUserAndClienttokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRefreshTokenForUserAndClientQuiet {
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
				dirPath := filepath.Join(GetRefreshTokenForUserAndClientBackupDir, "user", "getRefreshTokenForAndClient")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetRefreshTokenForUserAndClientuserId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRefreshTokenForUserAndClientQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRefreshTokenForUserAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClienttokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the User to a file")

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClientBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRefreshTokenForUserAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRefreshTokenForUserAndClientCmd := NewGetRefreshTokenForUserAndClientCmd()
	UserCmd.AddCommand(GetRefreshTokenForUserAndClientCmd)
}

var (
	RevokeTokenForUserAndClientuserId string

	RevokeTokenForUserAndClientclientId string

	RevokeTokenForUserAndClienttokenId string

	RevokeTokenForUserAndClientQuiet bool
)

func NewRevokeTokenForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeTokenForAndClient",
		Long: "Revoke a Token for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeTokenForUserAndClient(apiClient.GetConfig().Context, RevokeTokenForUserAndClientuserId, RevokeTokenForUserAndClientclientId, RevokeTokenForUserAndClienttokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeTokenForUserAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeTokenForUserAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeTokenForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeTokenForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().StringVarP(&RevokeTokenForUserAndClienttokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	cmd.Flags().BoolVarP(&RevokeTokenForUserAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeTokenForUserAndClientCmd := NewRevokeTokenForUserAndClientCmd()
	UserCmd.AddCommand(RevokeTokenForUserAndClientCmd)
}

var (
	ChangePassworduserId string

	ChangePassworddata string

	ChangePasswordQuiet bool
)

func NewChangePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "changePassword",
		Long: "Change Password",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ChangePassword(apiClient.GetConfig().Context, ChangePassworduserId)

			if ChangePassworddata != "" {
				req = req.Data(ChangePassworddata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ChangePasswordQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ChangePasswordQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ChangePassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ChangePassworddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ChangePasswordQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ChangePasswordCmd := NewChangePasswordCmd()
	UserCmd.AddCommand(ChangePasswordCmd)
}

var (
	ChangeRecoveryQuestionuserId string

	ChangeRecoveryQuestiondata string

	ChangeRecoveryQuestionQuiet bool
)

func NewChangeRecoveryQuestionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "changeRecoveryQuestion",
		Long: "Change Recovery Question",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ChangeRecoveryQuestion(apiClient.GetConfig().Context, ChangeRecoveryQuestionuserId)

			if ChangeRecoveryQuestiondata != "" {
				req = req.Data(ChangeRecoveryQuestiondata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ChangeRecoveryQuestionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ChangeRecoveryQuestionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ChangeRecoveryQuestionuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ChangeRecoveryQuestiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ChangeRecoveryQuestionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ChangeRecoveryQuestionCmd := NewChangeRecoveryQuestionCmd()
	UserCmd.AddCommand(ChangeRecoveryQuestionCmd)
}

var (
	ForgotPassworduserId string

	ForgotPasswordQuiet bool
)

func NewForgotPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "forgotPassword",
		Long: "Initiate Forgot Password",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ForgotPassword(apiClient.GetConfig().Context, ForgotPassworduserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ForgotPasswordQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ForgotPasswordQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ForgotPassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&ForgotPasswordQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ForgotPasswordCmd := NewForgotPasswordCmd()
	UserCmd.AddCommand(ForgotPasswordCmd)
}

var (
	ForgotPasswordSetNewPassworduserId string

	ForgotPasswordSetNewPassworddata string

	ForgotPasswordSetNewPasswordQuiet bool
)

func NewForgotPasswordSetNewPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "forgotPasswordSetNewPassword",
		Long: "Reset Password with Recovery Question",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ForgotPasswordSetNewPassword(apiClient.GetConfig().Context, ForgotPasswordSetNewPassworduserId)

			if ForgotPasswordSetNewPassworddata != "" {
				req = req.Data(ForgotPasswordSetNewPassworddata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ForgotPasswordSetNewPasswordQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ForgotPasswordSetNewPasswordQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ForgotPasswordSetNewPassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ForgotPasswordSetNewPassworddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ForgotPasswordSetNewPasswordQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ForgotPasswordSetNewPasswordCmd := NewForgotPasswordSetNewPasswordCmd()
	UserCmd.AddCommand(ForgotPasswordSetNewPasswordCmd)
}

var (
	ListUserGrantsuserId string

	ListUserGrantsBackupDir string

	ListUserGrantsLimit    int32
	ListUserGrantsPage     string
	ListUserGrantsFetchAll bool

	ListUserGrantsQuiet bool
)

func NewListUserGrantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGrants",
		Long: "List all User Grants",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserGrants(apiClient.GetConfig().Context, ListUserGrantsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUserGrantsQuiet {
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
					if !ListUserGrantsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUserGrantsFetchAll || len(items) == 0 {
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

			if ListUserGrantsFetchAll && pageCount > 1 && !ListUserGrantsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUserGrantsBackupDir, "user", "listGrants")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUserGrantsQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUserGrantsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUserGrantsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUserGrantsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUserGrantsQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUserGrantsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListUserGrantsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListUserGrantsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUserGrantsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUserGrantsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListUserGrantsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUserGrantsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUserGrantsCmd := NewListUserGrantsCmd()
	UserCmd.AddCommand(ListUserGrantsCmd)
}

var (
	RevokeUserGrantsuserId string

	RevokeUserGrantsQuiet bool
)

func NewRevokeUserGrantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeGrants",
		Long: "Revoke all User Grants",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeUserGrants(apiClient.GetConfig().Context, RevokeUserGrantsuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeUserGrantsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeUserGrantsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeUserGrantsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&RevokeUserGrantsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeUserGrantsCmd := NewRevokeUserGrantsCmd()
	UserCmd.AddCommand(RevokeUserGrantsCmd)
}

var (
	GetUserGrantuserId string

	GetUserGrantgrantId string

	GetUserGrantBackupDir string

	GetUserGrantQuiet bool
)

func NewGetUserGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getGrant",
		Long: "Retrieve a User Grant",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GetUserGrant(apiClient.GetConfig().Context, GetUserGrantuserId, GetUserGrantgrantId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetUserGrantQuiet {
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
				dirPath := filepath.Join(GetUserGrantBackupDir, "user", "getGrant")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetUserGrantuserId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetUserGrantQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetUserGrantQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetUserGrantuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetUserGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the User to a file")

	cmd.Flags().StringVarP(&GetUserGrantBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetUserGrantQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetUserGrantCmd := NewGetUserGrantCmd()
	UserCmd.AddCommand(GetUserGrantCmd)
}

var (
	RevokeUserGrantuserId string

	RevokeUserGrantgrantId string

	RevokeUserGrantQuiet bool
)

func NewRevokeUserGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeGrant",
		Long: "Revoke a User Grant",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeUserGrant(apiClient.GetConfig().Context, RevokeUserGrantuserId, RevokeUserGrantgrantId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeUserGrantQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeUserGrantQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeUserGrantuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeUserGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	cmd.Flags().BoolVarP(&RevokeUserGrantQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeUserGrantCmd := NewRevokeUserGrantCmd()
	UserCmd.AddCommand(RevokeUserGrantCmd)
}

var (
	ListUserGroupsuserId string

	ListUserGroupsBackupDir string

	ListUserGroupsLimit    int32
	ListUserGroupsPage     string
	ListUserGroupsFetchAll bool

	ListUserGroupsQuiet bool
)

func NewListUserGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGroups",
		Long: "List all Groups",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserGroups(apiClient.GetConfig().Context, ListUserGroupsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUserGroupsQuiet {
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
					if !ListUserGroupsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUserGroupsFetchAll || len(items) == 0 {
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

			if ListUserGroupsFetchAll && pageCount > 1 && !ListUserGroupsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUserGroupsBackupDir, "user", "listGroups")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUserGroupsQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUserGroupsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUserGroupsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUserGroupsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUserGroupsQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUserGroupsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListUserGroupsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListUserGroupsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUserGroupsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUserGroupsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListUserGroupsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUserGroupsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUserGroupsCmd := NewListUserGroupsCmd()
	UserCmd.AddCommand(ListUserGroupsCmd)
}

var (
	ListUserIdentityProvidersuserId string

	ListUserIdentityProvidersBackupDir string

	ListUserIdentityProvidersLimit    int32
	ListUserIdentityProvidersPage     string
	ListUserIdentityProvidersFetchAll bool

	ListUserIdentityProvidersQuiet bool
)

func NewListUserIdentityProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listIdentityProviders",
		Long: "List all Identity Providers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserIdentityProviders(apiClient.GetConfig().Context, ListUserIdentityProvidersuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListUserIdentityProvidersQuiet {
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
					if !ListUserIdentityProvidersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListUserIdentityProvidersFetchAll || len(items) == 0 {
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

			if ListUserIdentityProvidersFetchAll && pageCount > 1 && !ListUserIdentityProvidersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListUserIdentityProvidersBackupDir, "user", "listIdentityProviders")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListUserIdentityProvidersQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListUserIdentityProvidersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListUserIdentityProvidersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListUserIdentityProvidersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListUserIdentityProvidersQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListUserIdentityProvidersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListUserIdentityProvidersuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListUserIdentityProvidersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListUserIdentityProvidersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListUserIdentityProvidersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListUserIdentityProvidersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListUserIdentityProvidersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListUserIdentityProvidersCmd := NewListUserIdentityProvidersCmd()
	UserCmd.AddCommand(ListUserIdentityProvidersCmd)
}

var (
	ActivateUseruserId string

	ActivateUserQuiet bool
)

func NewActivateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ActivateUser(apiClient.GetConfig().Context, ActivateUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&ActivateUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateUserCmd := NewActivateUserCmd()
	UserCmd.AddCommand(ActivateUserCmd)
}

var (
	DeactivateUseruserId string

	DeactivateUserQuiet bool
)

func NewDeactivateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.DeactivateUser(apiClient.GetConfig().Context, DeactivateUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&DeactivateUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateUserCmd := NewDeactivateUserCmd()
	UserCmd.AddCommand(DeactivateUserCmd)
}

var (
	ExpirePassworduserId string

	ExpirePasswordQuiet bool
)

func NewExpirePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "expirePassword",
		Long: "Expire Password",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ExpirePassword(apiClient.GetConfig().Context, ExpirePassworduserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ExpirePasswordQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ExpirePasswordQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ExpirePassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&ExpirePasswordQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ExpirePasswordCmd := NewExpirePasswordCmd()
	UserCmd.AddCommand(ExpirePasswordCmd)
}

var (
	ExpirePasswordAndGetTemporaryPassworduserId string

	ExpirePasswordAndGetTemporaryPasswordQuiet bool
)

func NewExpirePasswordAndGetTemporaryPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "expirePasswordAndGetTemporaryPassword",
		Long: "Expire Password and Set Temporary Password",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ExpirePasswordAndGetTemporaryPassword(apiClient.GetConfig().Context, ExpirePasswordAndGetTemporaryPassworduserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ExpirePasswordAndGetTemporaryPasswordQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ExpirePasswordAndGetTemporaryPasswordQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ExpirePasswordAndGetTemporaryPassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&ExpirePasswordAndGetTemporaryPasswordQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ExpirePasswordAndGetTemporaryPasswordCmd := NewExpirePasswordAndGetTemporaryPasswordCmd()
	UserCmd.AddCommand(ExpirePasswordAndGetTemporaryPasswordCmd)
}

var (
	ReactivateUseruserId string

	ReactivateUserQuiet bool
)

func NewReactivateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "reactivate",
		Long: "Reactivate a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ReactivateUser(apiClient.GetConfig().Context, ReactivateUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReactivateUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReactivateUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReactivateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&ReactivateUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReactivateUserCmd := NewReactivateUserCmd()
	UserCmd.AddCommand(ReactivateUserCmd)
}

var (
	ResetFactorsuserId string

	ResetFactorsQuiet bool
)

func NewResetFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "resetFactors",
		Long: "Reset all Factors",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ResetFactors(apiClient.GetConfig().Context, ResetFactorsuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ResetFactorsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ResetFactorsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ResetFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&ResetFactorsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ResetFactorsCmd := NewResetFactorsCmd()
	UserCmd.AddCommand(ResetFactorsCmd)
}

var (
	GenerateResetPasswordTokenuserId string

	GenerateResetPasswordTokenQuiet bool
)

func NewGenerateResetPasswordTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateResetPasswordToken",
		Long: "Generate a Reset Password Token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GenerateResetPasswordToken(apiClient.GetConfig().Context, GenerateResetPasswordTokenuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GenerateResetPasswordTokenQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GenerateResetPasswordTokenQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GenerateResetPasswordTokenuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&GenerateResetPasswordTokenQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GenerateResetPasswordTokenCmd := NewGenerateResetPasswordTokenCmd()
	UserCmd.AddCommand(GenerateResetPasswordTokenCmd)
}

var (
	SuspendUseruserId string

	SuspendUserQuiet bool
)

func NewSuspendUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "suspend",
		Long: "Suspend a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.SuspendUser(apiClient.GetConfig().Context, SuspendUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !SuspendUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !SuspendUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&SuspendUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&SuspendUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	SuspendUserCmd := NewSuspendUserCmd()
	UserCmd.AddCommand(SuspendUserCmd)
}

var (
	UnlockUseruserId string

	UnlockUserQuiet bool
)

func NewUnlockUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unlock",
		Long: "Unlock a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.UnlockUser(apiClient.GetConfig().Context, UnlockUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnlockUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnlockUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnlockUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&UnlockUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnlockUserCmd := NewUnlockUserCmd()
	UserCmd.AddCommand(UnlockUserCmd)
}

var (
	UnsuspendUseruserId string

	UnsuspendUserQuiet bool
)

func NewUnsuspendUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unsuspend",
		Long: "Unsuspend a User",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.UnsuspendUser(apiClient.GetConfig().Context, UnsuspendUseruserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnsuspendUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnsuspendUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnsuspendUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&UnsuspendUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnsuspendUserCmd := NewUnsuspendUserCmd()
	UserCmd.AddCommand(UnsuspendUserCmd)
}

var (
	SetLinkedObjectForUseruserId string

	SetLinkedObjectForUserprimaryRelationshipName string

	SetLinkedObjectForUserprimaryUserId string

	SetLinkedObjectForUserQuiet bool
)

func NewSetLinkedObjectForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "setLinkedObjectFor",
		Long: "Create a Linked Object for two Users",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.SetLinkedObjectForUser(apiClient.GetConfig().Context, SetLinkedObjectForUseruserId, SetLinkedObjectForUserprimaryRelationshipName, SetLinkedObjectForUserprimaryUserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !SetLinkedObjectForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !SetLinkedObjectForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&SetLinkedObjectForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&SetLinkedObjectForUserprimaryRelationshipName, "primaryRelationshipName", "", "", "")
	cmd.MarkFlagRequired("primaryRelationshipName")

	cmd.Flags().StringVarP(&SetLinkedObjectForUserprimaryUserId, "primaryUserId", "", "", "")
	cmd.MarkFlagRequired("primaryUserId")

	cmd.Flags().BoolVarP(&SetLinkedObjectForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	SetLinkedObjectForUserCmd := NewSetLinkedObjectForUserCmd()
	UserCmd.AddCommand(SetLinkedObjectForUserCmd)
}

var (
	ListLinkedObjectsForUseruserId string

	ListLinkedObjectsForUserrelationshipName string

	ListLinkedObjectsForUserBackupDir string

	ListLinkedObjectsForUserLimit    int32
	ListLinkedObjectsForUserPage     string
	ListLinkedObjectsForUserFetchAll bool

	ListLinkedObjectsForUserQuiet bool
)

func NewListLinkedObjectsForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listLinkedObjectsFor",
		Long: "List all Linked Objects",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListLinkedObjectsForUser(apiClient.GetConfig().Context, ListLinkedObjectsForUseruserId, ListLinkedObjectsForUserrelationshipName)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListLinkedObjectsForUserQuiet {
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
					if !ListLinkedObjectsForUserQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListLinkedObjectsForUserFetchAll || len(items) == 0 {
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

			if ListLinkedObjectsForUserFetchAll && pageCount > 1 && !ListLinkedObjectsForUserQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListLinkedObjectsForUserBackupDir, "user", "listLinkedObjectsFor")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListLinkedObjectsForUserQuiet {
					fmt.Printf("Backing up Users to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListLinkedObjectsForUserQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListLinkedObjectsForUserQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListLinkedObjectsForUserQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListLinkedObjectsForUserQuiet {
					fmt.Printf("Successfully backed up %d/%d Users\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListLinkedObjectsForUserQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListLinkedObjectsForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListLinkedObjectsForUserrelationshipName, "relationshipName", "", "", "")
	cmd.MarkFlagRequired("relationshipName")

	cmd.Flags().Int32VarP(&ListLinkedObjectsForUserLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListLinkedObjectsForUserPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListLinkedObjectsForUserFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Users to a directory")

	cmd.Flags().StringVarP(&ListLinkedObjectsForUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListLinkedObjectsForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListLinkedObjectsForUserCmd := NewListLinkedObjectsForUserCmd()
	UserCmd.AddCommand(ListLinkedObjectsForUserCmd)
}

var (
	DeleteLinkedObjectForUseruserId string

	DeleteLinkedObjectForUserrelationshipName string

	DeleteLinkedObjectForUserQuiet bool
)

func NewDeleteLinkedObjectForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteLinkedObjectFor",
		Long: "Delete a Linked Object",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.DeleteLinkedObjectForUser(apiClient.GetConfig().Context, DeleteLinkedObjectForUseruserId, DeleteLinkedObjectForUserrelationshipName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteLinkedObjectForUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteLinkedObjectForUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteLinkedObjectForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&DeleteLinkedObjectForUserrelationshipName, "relationshipName", "", "", "")
	cmd.MarkFlagRequired("relationshipName")

	cmd.Flags().BoolVarP(&DeleteLinkedObjectForUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteLinkedObjectForUserCmd := NewDeleteLinkedObjectForUserCmd()
	UserCmd.AddCommand(DeleteLinkedObjectForUserCmd)
}

var (
	RevokeUserSessionsuserId string

	RevokeUserSessionsQuiet bool
)

func NewRevokeUserSessionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeSessions",
		Long: "Revoke all User Sessions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeUserSessions(apiClient.GetConfig().Context, RevokeUserSessionsuserId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeUserSessionsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeUserSessionsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeUserSessionsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().BoolVarP(&RevokeUserSessionsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeUserSessionsCmd := NewRevokeUserSessionsCmd()
	UserCmd.AddCommand(RevokeUserSessionsCmd)
}
