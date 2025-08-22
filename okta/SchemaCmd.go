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

var SchemaCmd = &cobra.Command{
	Use:  "schema",
	Long: "Manage SchemaAPI",
}

func init() {
	rootCmd.AddCommand(SchemaCmd)
}

var (
	UpdateApplicationUserProfileappId string

	UpdateApplicationUserProfiledata string

	UpdateApplicationUserProfileQuiet bool
)

func NewUpdateApplicationUserProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateApplicationUserProfile",
		Long: "Update the default Application User Schema for an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.UpdateApplicationUserProfile(apiClient.GetConfig().Context, UpdateApplicationUserProfileappId)

			if UpdateApplicationUserProfiledata != "" {
				req = req.Data(UpdateApplicationUserProfiledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateApplicationUserProfileQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateApplicationUserProfileQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateApplicationUserProfileappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateApplicationUserProfiledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateApplicationUserProfileQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateApplicationUserProfileCmd := NewUpdateApplicationUserProfileCmd()
	SchemaCmd.AddCommand(UpdateApplicationUserProfileCmd)
}

var (
	GetApplicationUserSchemaappId string

	GetApplicationUserSchemaBackupDir string

	GetApplicationUserSchemaQuiet bool
)

func NewGetApplicationUserSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationUser",
		Long: "Retrieve the default Application User Schema for an Application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetApplicationUserSchema(apiClient.GetConfig().Context, GetApplicationUserSchemaappId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetApplicationUserSchemaQuiet {
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
				dirPath := filepath.Join(GetApplicationUserSchemaBackupDir, "schema", "getApplicationUser")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetApplicationUserSchemaappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetApplicationUserSchemaQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetApplicationUserSchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApplicationUserSchemaappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Schema to a file")

	cmd.Flags().StringVarP(&GetApplicationUserSchemaBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetApplicationUserSchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetApplicationUserSchemaCmd := NewGetApplicationUserSchemaCmd()
	SchemaCmd.AddCommand(GetApplicationUserSchemaCmd)
}

var (
	UpdateGroupSchemadata string

	UpdateGroupSchemaQuiet bool
)

func NewUpdateGroupSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateGroup",
		Long: "Update the default Group Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.UpdateGroupSchema(apiClient.GetConfig().Context)

			if UpdateGroupSchemadata != "" {
				req = req.Data(UpdateGroupSchemadata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateGroupSchemaQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateGroupSchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateGroupSchemadata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateGroupSchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateGroupSchemaCmd := NewUpdateGroupSchemaCmd()
	SchemaCmd.AddCommand(UpdateGroupSchemaCmd)
}

var (
	GetGroupSchemaBackupDir string

	GetGroupSchemaQuiet bool
)

func NewGetGroupSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getGroup",
		Long: "Retrieve the default Group Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetGroupSchema(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetGroupSchemaQuiet {
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
				dirPath := filepath.Join(GetGroupSchemaBackupDir, "schema", "getGroup")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "schema.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetGroupSchemaQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetGroupSchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the Schema to a file")

	cmd.Flags().StringVarP(&GetGroupSchemaBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetGroupSchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetGroupSchemaCmd := NewGetGroupSchemaCmd()
	SchemaCmd.AddCommand(GetGroupSchemaCmd)
}

var (
	ListLogStreamSchemasBackupDir string

	ListLogStreamSchemasLimit    int32
	ListLogStreamSchemasPage     string
	ListLogStreamSchemasFetchAll bool

	ListLogStreamSchemasQuiet bool
)

func NewListLogStreamSchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listLogStreams",
		Long: "List the Log Stream Schemas",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.ListLogStreamSchemas(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListLogStreamSchemasQuiet {
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
					if !ListLogStreamSchemasQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListLogStreamSchemasFetchAll || len(items) == 0 {
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

			if ListLogStreamSchemasFetchAll && pageCount > 1 && !ListLogStreamSchemasQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListLogStreamSchemasBackupDir, "schema", "listLogStreams")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListLogStreamSchemasQuiet {
					fmt.Printf("Backing up Schemas to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListLogStreamSchemasQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListLogStreamSchemasQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListLogStreamSchemasQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListLogStreamSchemasQuiet {
					fmt.Printf("Successfully backed up %d/%d Schemas\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListLogStreamSchemasQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListLogStreamSchemasLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListLogStreamSchemasPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListLogStreamSchemasFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Schemas to a directory")

	cmd.Flags().StringVarP(&ListLogStreamSchemasBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListLogStreamSchemasQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListLogStreamSchemasCmd := NewListLogStreamSchemasCmd()
	SchemaCmd.AddCommand(ListLogStreamSchemasCmd)
}

var (
	GetLogStreamSchemalogStreamType string

	GetLogStreamSchemaBackupDir string

	GetLogStreamSchemaQuiet bool
)

func NewGetLogStreamSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getLogStream",
		Long: "Retrieve the Log Stream Schema for the schema type",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetLogStreamSchema(apiClient.GetConfig().Context, GetLogStreamSchemalogStreamType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetLogStreamSchemaQuiet {
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
				dirPath := filepath.Join(GetLogStreamSchemaBackupDir, "schema", "getLogStream")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetLogStreamSchemalogStreamType
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetLogStreamSchemaQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetLogStreamSchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetLogStreamSchemalogStreamType, "logStreamType", "", "", "")
	cmd.MarkFlagRequired("logStreamType")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Schema to a file")

	cmd.Flags().StringVarP(&GetLogStreamSchemaBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetLogStreamSchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetLogStreamSchemaCmd := NewGetLogStreamSchemaCmd()
	SchemaCmd.AddCommand(GetLogStreamSchemaCmd)
}

var (
	UpdateUserProfileschemaId string

	UpdateUserProfiledata string

	UpdateUserProfileQuiet bool
)

func NewUpdateUserProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateUserProfile",
		Long: "Update a User Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.UpdateUserProfile(apiClient.GetConfig().Context, UpdateUserProfileschemaId)

			if UpdateUserProfiledata != "" {
				req = req.Data(UpdateUserProfiledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateUserProfileQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateUserProfileQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateUserProfileschemaId, "schemaId", "", "", "")
	cmd.MarkFlagRequired("schemaId")

	cmd.Flags().StringVarP(&UpdateUserProfiledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateUserProfileQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateUserProfileCmd := NewUpdateUserProfileCmd()
	SchemaCmd.AddCommand(UpdateUserProfileCmd)
}

var (
	GetUserSchemaschemaId string

	GetUserSchemaBackupDir string

	GetUserSchemaQuiet bool
)

func NewGetUserSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getUser",
		Long: "Retrieve a User Schema",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetUserSchema(apiClient.GetConfig().Context, GetUserSchemaschemaId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetUserSchemaQuiet {
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
				dirPath := filepath.Join(GetUserSchemaBackupDir, "schema", "getUser")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetUserSchemaschemaId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetUserSchemaQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetUserSchemaQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetUserSchemaschemaId, "schemaId", "", "", "")
	cmd.MarkFlagRequired("schemaId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Schema to a file")

	cmd.Flags().StringVarP(&GetUserSchemaBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetUserSchemaQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetUserSchemaCmd := NewGetUserSchemaCmd()
	SchemaCmd.AddCommand(GetUserSchemaCmd)
}
