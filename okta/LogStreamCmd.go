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

var LogStreamCmd = &cobra.Command{
	Use:  "logStream",
	Long: "Manage LogStreamAPI",
}

func init() {
	rootCmd.AddCommand(LogStreamCmd)
}

var (
	CreateLogStreamdata string

	CreateLogStreamRestoreFile string

	CreateLogStreamQuiet bool
)

func NewCreateLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Log Stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateLogStreamRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateLogStreamRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateLogStreamdata = string(processedData)

				if !CreateLogStreamQuiet {
					fmt.Println("Restoring LogStream from:", CreateLogStreamRestoreFile)
				}
			}

			req := apiClient.LogStreamAPI.CreateLogStream(apiClient.GetConfig().Context)

			if CreateLogStreamdata != "" {
				req = req.Data(CreateLogStreamdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateLogStreamQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateLogStreamQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateLogStreamdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateLogStreamRestoreFile, "restore-from", "r", "", "Restore LogStream from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateLogStreamQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateLogStreamCmd := NewCreateLogStreamCmd()
	LogStreamCmd.AddCommand(CreateLogStreamCmd)
}

var (
	ListLogStreamsBackupDir string

	ListLogStreamsLimit    int32
	ListLogStreamsPage     string
	ListLogStreamsFetchAll bool

	ListLogStreamsQuiet bool
)

func NewListLogStreamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Log Streams",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.ListLogStreams(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListLogStreamsQuiet {
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
					if !ListLogStreamsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListLogStreamsFetchAll || len(items) == 0 {
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

			if ListLogStreamsFetchAll && pageCount > 1 && !ListLogStreamsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListLogStreamsBackupDir, "logstream", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListLogStreamsQuiet {
					fmt.Printf("Backing up LogStreams to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListLogStreamsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListLogStreamsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListLogStreamsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListLogStreamsQuiet {
					fmt.Printf("Successfully backed up %d/%d LogStreams\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListLogStreamsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListLogStreamsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListLogStreamsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListLogStreamsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple LogStreams to a directory")

	cmd.Flags().StringVarP(&ListLogStreamsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListLogStreamsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListLogStreamsCmd := NewListLogStreamsCmd()
	LogStreamCmd.AddCommand(ListLogStreamsCmd)
}

var (
	GetLogStreamlogStreamId string

	GetLogStreamBackupDir string

	GetLogStreamQuiet bool
)

func NewGetLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Log Stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.GetLogStream(apiClient.GetConfig().Context, GetLogStreamlogStreamId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetLogStreamQuiet {
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
				dirPath := filepath.Join(GetLogStreamBackupDir, "logstream", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetLogStreamlogStreamId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetLogStreamQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetLogStreamQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the LogStream to a file")

	cmd.Flags().StringVarP(&GetLogStreamBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetLogStreamQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetLogStreamCmd := NewGetLogStreamCmd()
	LogStreamCmd.AddCommand(GetLogStreamCmd)
}

var (
	ReplaceLogStreamlogStreamId string

	ReplaceLogStreamdata string

	ReplaceLogStreamQuiet bool
)

func NewReplaceLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Log Stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.ReplaceLogStream(apiClient.GetConfig().Context, ReplaceLogStreamlogStreamId)

			if ReplaceLogStreamdata != "" {
				req = req.Data(ReplaceLogStreamdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceLogStreamQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceLogStreamQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	cmd.Flags().StringVarP(&ReplaceLogStreamdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceLogStreamQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceLogStreamCmd := NewReplaceLogStreamCmd()
	LogStreamCmd.AddCommand(ReplaceLogStreamCmd)
}

var (
	DeleteLogStreamlogStreamId string

	DeleteLogStreamQuiet bool
)

func NewDeleteLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Log Stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.DeleteLogStream(apiClient.GetConfig().Context, DeleteLogStreamlogStreamId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteLogStreamQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteLogStreamQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	cmd.Flags().BoolVarP(&DeleteLogStreamQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteLogStreamCmd := NewDeleteLogStreamCmd()
	LogStreamCmd.AddCommand(DeleteLogStreamCmd)
}

var (
	ActivateLogStreamlogStreamId string

	ActivateLogStreamQuiet bool
)

func NewActivateLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Log Stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.ActivateLogStream(apiClient.GetConfig().Context, ActivateLogStreamlogStreamId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateLogStreamQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateLogStreamQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	cmd.Flags().BoolVarP(&ActivateLogStreamQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateLogStreamCmd := NewActivateLogStreamCmd()
	LogStreamCmd.AddCommand(ActivateLogStreamCmd)
}

var (
	DeactivateLogStreamlogStreamId string

	DeactivateLogStreamQuiet bool
)

func NewDeactivateLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Log Stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.DeactivateLogStream(apiClient.GetConfig().Context, DeactivateLogStreamlogStreamId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateLogStreamQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateLogStreamQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	cmd.Flags().BoolVarP(&DeactivateLogStreamQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateLogStreamCmd := NewDeactivateLogStreamCmd()
	LogStreamCmd.AddCommand(DeactivateLogStreamCmd)
}
