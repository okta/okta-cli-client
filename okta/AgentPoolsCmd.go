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

var AgentPoolsCmd = &cobra.Command{
	Use:  "agentPools",
	Long: "Manage AgentPoolsAPI",
}

func init() {
	rootCmd.AddCommand(AgentPoolsCmd)
}

var (
	ListAgentPoolsBackupDir string

	ListAgentPoolsLimit    int32
	ListAgentPoolsPage     string
	ListAgentPoolsFetchAll bool

	ListAgentPoolsQuiet bool
)

func NewListAgentPoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List all Agent Pools",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.ListAgentPools(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAgentPoolsQuiet {
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
					if !ListAgentPoolsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAgentPoolsFetchAll || len(items) == 0 {
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

			if ListAgentPoolsFetchAll && pageCount > 1 && !ListAgentPoolsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAgentPoolsBackupDir, "agentpools", "list")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAgentPoolsQuiet {
					fmt.Printf("Backing up AgentPoolss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAgentPoolsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAgentPoolsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAgentPoolsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAgentPoolsQuiet {
					fmt.Printf("Successfully backed up %d/%d AgentPoolss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAgentPoolsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListAgentPoolsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAgentPoolsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAgentPoolsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AgentPoolss to a directory")

	cmd.Flags().StringVarP(&ListAgentPoolsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAgentPoolsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAgentPoolsCmd := NewListAgentPoolsCmd()
	AgentPoolsCmd.AddCommand(ListAgentPoolsCmd)
}

var (
	CreateAgentPoolsUpdatepoolId string

	CreateAgentPoolsUpdatedata string

	CreateAgentPoolsUpdateRestoreFile string

	CreateAgentPoolsUpdateQuiet bool
)

func NewCreateAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createUpdate",
		Long: "Create an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateAgentPoolsUpdateRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateAgentPoolsUpdateRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateAgentPoolsUpdatedata = string(processedData)

				if !CreateAgentPoolsUpdateQuiet {
					fmt.Println("Restoring AgentPools from:", CreateAgentPoolsUpdateRestoreFile)
				}
			}

			req := apiClient.AgentPoolsAPI.CreateAgentPoolsUpdate(apiClient.GetConfig().Context, CreateAgentPoolsUpdatepoolId)

			if CreateAgentPoolsUpdatedata != "" {
				req = req.Data(CreateAgentPoolsUpdatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&CreateAgentPoolsUpdatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateAgentPoolsUpdateRestoreFile, "restore-from", "r", "", "Restore AgentPools from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateAgentPoolsUpdateCmd := NewCreateAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(CreateAgentPoolsUpdateCmd)
}

var (
	ListAgentPoolsUpdatespoolId string

	ListAgentPoolsUpdatesBackupDir string

	ListAgentPoolsUpdatesLimit    int32
	ListAgentPoolsUpdatesPage     string
	ListAgentPoolsUpdatesFetchAll bool

	ListAgentPoolsUpdatesQuiet bool
)

func NewListAgentPoolsUpdatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listUpdates",
		Long: "List all Agent Pool updates",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.ListAgentPoolsUpdates(apiClient.GetConfig().Context, ListAgentPoolsUpdatespoolId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAgentPoolsUpdatesQuiet {
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
					if !ListAgentPoolsUpdatesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAgentPoolsUpdatesFetchAll || len(items) == 0 {
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

			if ListAgentPoolsUpdatesFetchAll && pageCount > 1 && !ListAgentPoolsUpdatesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAgentPoolsUpdatesBackupDir, "agentpools", "listUpdates")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAgentPoolsUpdatesQuiet {
					fmt.Printf("Backing up AgentPoolss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAgentPoolsUpdatesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAgentPoolsUpdatesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAgentPoolsUpdatesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAgentPoolsUpdatesQuiet {
					fmt.Printf("Successfully backed up %d/%d AgentPoolss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAgentPoolsUpdatesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAgentPoolsUpdatespoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().Int32VarP(&ListAgentPoolsUpdatesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAgentPoolsUpdatesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAgentPoolsUpdatesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AgentPoolss to a directory")

	cmd.Flags().StringVarP(&ListAgentPoolsUpdatesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAgentPoolsUpdatesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAgentPoolsUpdatesCmd := NewListAgentPoolsUpdatesCmd()
	AgentPoolsCmd.AddCommand(ListAgentPoolsUpdatesCmd)
}

var (
	UpdateAgentPoolsUpdateSettingspoolId string

	UpdateAgentPoolsUpdateSettingsdata string

	UpdateAgentPoolsUpdateSettingsQuiet bool
)

func NewUpdateAgentPoolsUpdateSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateUpdateSettings",
		Long: "Update an Agent Pool update settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.UpdateAgentPoolsUpdateSettings(apiClient.GetConfig().Context, UpdateAgentPoolsUpdateSettingspoolId)

			if UpdateAgentPoolsUpdateSettingsdata != "" {
				req = req.Data(UpdateAgentPoolsUpdateSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateAgentPoolsUpdateSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateAgentPoolsUpdateSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateAgentPoolsUpdateSettingspoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&UpdateAgentPoolsUpdateSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateAgentPoolsUpdateSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateAgentPoolsUpdateSettingsCmd := NewUpdateAgentPoolsUpdateSettingsCmd()
	AgentPoolsCmd.AddCommand(UpdateAgentPoolsUpdateSettingsCmd)
}

var (
	GetAgentPoolsUpdateSettingspoolId string

	GetAgentPoolsUpdateSettingsBackupDir string

	GetAgentPoolsUpdateSettingsQuiet bool
)

func NewGetAgentPoolsUpdateSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getUpdateSettings",
		Long: "Retrieve an Agent Pool update's settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.GetAgentPoolsUpdateSettings(apiClient.GetConfig().Context, GetAgentPoolsUpdateSettingspoolId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAgentPoolsUpdateSettingsQuiet {
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
				dirPath := filepath.Join(GetAgentPoolsUpdateSettingsBackupDir, "agentpools", "getUpdateSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAgentPoolsUpdateSettingspoolId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAgentPoolsUpdateSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAgentPoolsUpdateSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAgentPoolsUpdateSettingspoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AgentPools to a file")

	cmd.Flags().StringVarP(&GetAgentPoolsUpdateSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAgentPoolsUpdateSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAgentPoolsUpdateSettingsCmd := NewGetAgentPoolsUpdateSettingsCmd()
	AgentPoolsCmd.AddCommand(GetAgentPoolsUpdateSettingsCmd)
}

var (
	UpdateAgentPoolsUpdatepoolId string

	UpdateAgentPoolsUpdateupdateId string

	UpdateAgentPoolsUpdatedata string

	UpdateAgentPoolsUpdateQuiet bool
)

func NewUpdateAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateUpdate",
		Long: "Update an Agent Pool update by id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.UpdateAgentPoolsUpdate(apiClient.GetConfig().Context, UpdateAgentPoolsUpdatepoolId, UpdateAgentPoolsUpdateupdateId)

			if UpdateAgentPoolsUpdatedata != "" {
				req = req.Data(UpdateAgentPoolsUpdatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&UpdateAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().StringVarP(&UpdateAgentPoolsUpdatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateAgentPoolsUpdateCmd := NewUpdateAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(UpdateAgentPoolsUpdateCmd)
}

var (
	GetAgentPoolsUpdateInstancepoolId string

	GetAgentPoolsUpdateInstanceupdateId string

	GetAgentPoolsUpdateInstanceBackupDir string

	GetAgentPoolsUpdateInstanceQuiet bool
)

func NewGetAgentPoolsUpdateInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getUpdateInstance",
		Long: "Retrieve an Agent Pool update by id",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.GetAgentPoolsUpdateInstance(apiClient.GetConfig().Context, GetAgentPoolsUpdateInstancepoolId, GetAgentPoolsUpdateInstanceupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAgentPoolsUpdateInstanceQuiet {
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
				dirPath := filepath.Join(GetAgentPoolsUpdateInstanceBackupDir, "agentpools", "getUpdateInstance")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAgentPoolsUpdateInstancepoolId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAgentPoolsUpdateInstanceQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAgentPoolsUpdateInstanceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAgentPoolsUpdateInstancepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&GetAgentPoolsUpdateInstanceupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AgentPools to a file")

	cmd.Flags().StringVarP(&GetAgentPoolsUpdateInstanceBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAgentPoolsUpdateInstanceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAgentPoolsUpdateInstanceCmd := NewGetAgentPoolsUpdateInstanceCmd()
	AgentPoolsCmd.AddCommand(GetAgentPoolsUpdateInstanceCmd)
}

var (
	DeleteAgentPoolsUpdatepoolId string

	DeleteAgentPoolsUpdateupdateId string

	DeleteAgentPoolsUpdateQuiet bool
)

func NewDeleteAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteUpdate",
		Long: "Delete an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.DeleteAgentPoolsUpdate(apiClient.GetConfig().Context, DeleteAgentPoolsUpdatepoolId, DeleteAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&DeleteAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&DeleteAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteAgentPoolsUpdateCmd := NewDeleteAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(DeleteAgentPoolsUpdateCmd)
}

var (
	ActivateAgentPoolsUpdatepoolId string

	ActivateAgentPoolsUpdateupdateId string

	ActivateAgentPoolsUpdateQuiet bool
)

func NewActivateAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateUpdate",
		Long: "Activate an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.ActivateAgentPoolsUpdate(apiClient.GetConfig().Context, ActivateAgentPoolsUpdatepoolId, ActivateAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&ActivateAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&ActivateAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateAgentPoolsUpdateCmd := NewActivateAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(ActivateAgentPoolsUpdateCmd)
}

var (
	DeactivateAgentPoolsUpdatepoolId string

	DeactivateAgentPoolsUpdateupdateId string

	DeactivateAgentPoolsUpdateQuiet bool
)

func NewDeactivateAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateUpdate",
		Long: "Deactivate an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.DeactivateAgentPoolsUpdate(apiClient.GetConfig().Context, DeactivateAgentPoolsUpdatepoolId, DeactivateAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&DeactivateAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&DeactivateAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateAgentPoolsUpdateCmd := NewDeactivateAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(DeactivateAgentPoolsUpdateCmd)
}

var (
	PauseAgentPoolsUpdatepoolId string

	PauseAgentPoolsUpdateupdateId string

	PauseAgentPoolsUpdateQuiet bool
)

func NewPauseAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "pauseUpdate",
		Long: "Pause an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.PauseAgentPoolsUpdate(apiClient.GetConfig().Context, PauseAgentPoolsUpdatepoolId, PauseAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !PauseAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !PauseAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&PauseAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&PauseAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&PauseAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	PauseAgentPoolsUpdateCmd := NewPauseAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(PauseAgentPoolsUpdateCmd)
}

var (
	ResumeAgentPoolsUpdatepoolId string

	ResumeAgentPoolsUpdateupdateId string

	ResumeAgentPoolsUpdateQuiet bool
)

func NewResumeAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "resumeUpdate",
		Long: "Resume an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.ResumeAgentPoolsUpdate(apiClient.GetConfig().Context, ResumeAgentPoolsUpdatepoolId, ResumeAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ResumeAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ResumeAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ResumeAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&ResumeAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&ResumeAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ResumeAgentPoolsUpdateCmd := NewResumeAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(ResumeAgentPoolsUpdateCmd)
}

var (
	RetryAgentPoolsUpdatepoolId string

	RetryAgentPoolsUpdateupdateId string

	RetryAgentPoolsUpdateQuiet bool
)

func NewRetryAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "retryUpdate",
		Long: "Retry an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.RetryAgentPoolsUpdate(apiClient.GetConfig().Context, RetryAgentPoolsUpdatepoolId, RetryAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RetryAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RetryAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RetryAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&RetryAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&RetryAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RetryAgentPoolsUpdateCmd := NewRetryAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(RetryAgentPoolsUpdateCmd)
}

var (
	StopAgentPoolsUpdatepoolId string

	StopAgentPoolsUpdateupdateId string

	StopAgentPoolsUpdateQuiet bool
)

func NewStopAgentPoolsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "stopUpdate",
		Long: "Stop an Agent Pool update",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AgentPoolsAPI.StopAgentPoolsUpdate(apiClient.GetConfig().Context, StopAgentPoolsUpdatepoolId, StopAgentPoolsUpdateupdateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !StopAgentPoolsUpdateQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !StopAgentPoolsUpdateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&StopAgentPoolsUpdatepoolId, "poolId", "", "", "")
	cmd.MarkFlagRequired("poolId")

	cmd.Flags().StringVarP(&StopAgentPoolsUpdateupdateId, "updateId", "", "", "")
	cmd.MarkFlagRequired("updateId")

	cmd.Flags().BoolVarP(&StopAgentPoolsUpdateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	StopAgentPoolsUpdateCmd := NewStopAgentPoolsUpdateCmd()
	AgentPoolsCmd.AddCommand(StopAgentPoolsUpdateCmd)
}
