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

var NetworkZoneCmd = &cobra.Command{
	Use:  "networkZone",
	Long: "Manage NetworkZoneAPI",
}

func init() {
	rootCmd.AddCommand(NetworkZoneCmd)
}

var (
	CreateNetworkZonedata string

	CreateNetworkZoneRestoreFile string

	CreateNetworkZoneQuiet bool
)

func NewCreateNetworkZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Network Zone",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateNetworkZoneRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateNetworkZoneRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateNetworkZonedata = string(processedData)

				if !CreateNetworkZoneQuiet {
					fmt.Println("Restoring NetworkZone from:", CreateNetworkZoneRestoreFile)
				}
			}

			req := apiClient.NetworkZoneAPI.CreateNetworkZone(apiClient.GetConfig().Context)

			if CreateNetworkZonedata != "" {
				req = req.Data(CreateNetworkZonedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateNetworkZoneQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateNetworkZoneQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateNetworkZonedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateNetworkZoneRestoreFile, "restore-from", "r", "", "Restore NetworkZone from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateNetworkZoneQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateNetworkZoneCmd := NewCreateNetworkZoneCmd()
	NetworkZoneCmd.AddCommand(CreateNetworkZoneCmd)
}

var (
	ListNetworkZonesBackupDir string

	ListNetworkZonesLimit    int32
	ListNetworkZonesPage     string
	ListNetworkZonesFetchAll bool

	ListNetworkZonesQuiet bool
)

func NewListNetworkZonesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Network Zones",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.NetworkZoneAPI.ListNetworkZones(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListNetworkZonesQuiet {
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
					if !ListNetworkZonesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListNetworkZonesFetchAll || len(items) == 0 {
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

			if ListNetworkZonesFetchAll && pageCount > 1 && !ListNetworkZonesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListNetworkZonesBackupDir, "networkzone", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListNetworkZonesQuiet {
					fmt.Printf("Backing up NetworkZones to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListNetworkZonesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListNetworkZonesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListNetworkZonesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListNetworkZonesQuiet {
					fmt.Printf("Successfully backed up %d/%d NetworkZones\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListNetworkZonesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListNetworkZonesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListNetworkZonesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListNetworkZonesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple NetworkZones to a directory")

	cmd.Flags().StringVarP(&ListNetworkZonesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListNetworkZonesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListNetworkZonesCmd := NewListNetworkZonesCmd()
	NetworkZoneCmd.AddCommand(ListNetworkZonesCmd)
}

var (
	GetNetworkZonezoneId string

	GetNetworkZoneBackupDir string

	GetNetworkZoneQuiet bool
)

func NewGetNetworkZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Network Zone",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.NetworkZoneAPI.GetNetworkZone(apiClient.GetConfig().Context, GetNetworkZonezoneId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetNetworkZoneQuiet {
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
				dirPath := filepath.Join(GetNetworkZoneBackupDir, "networkzone", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetNetworkZonezoneId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetNetworkZoneQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetNetworkZoneQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetNetworkZonezoneId, "zoneId", "", "", "")
	cmd.MarkFlagRequired("zoneId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the NetworkZone to a file")

	cmd.Flags().StringVarP(&GetNetworkZoneBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetNetworkZoneQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetNetworkZoneCmd := NewGetNetworkZoneCmd()
	NetworkZoneCmd.AddCommand(GetNetworkZoneCmd)
}

var (
	ReplaceNetworkZonezoneId string

	ReplaceNetworkZonedata string

	ReplaceNetworkZoneQuiet bool
)

func NewReplaceNetworkZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Network Zone",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.NetworkZoneAPI.ReplaceNetworkZone(apiClient.GetConfig().Context, ReplaceNetworkZonezoneId)

			if ReplaceNetworkZonedata != "" {
				req = req.Data(ReplaceNetworkZonedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceNetworkZoneQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceNetworkZoneQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceNetworkZonezoneId, "zoneId", "", "", "")
	cmd.MarkFlagRequired("zoneId")

	cmd.Flags().StringVarP(&ReplaceNetworkZonedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceNetworkZoneQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceNetworkZoneCmd := NewReplaceNetworkZoneCmd()
	NetworkZoneCmd.AddCommand(ReplaceNetworkZoneCmd)
}

var (
	DeleteNetworkZonezoneId string

	DeleteNetworkZoneQuiet bool
)

func NewDeleteNetworkZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Network Zone",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.NetworkZoneAPI.DeleteNetworkZone(apiClient.GetConfig().Context, DeleteNetworkZonezoneId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteNetworkZoneQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteNetworkZoneQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteNetworkZonezoneId, "zoneId", "", "", "")
	cmd.MarkFlagRequired("zoneId")

	cmd.Flags().BoolVarP(&DeleteNetworkZoneQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteNetworkZoneCmd := NewDeleteNetworkZoneCmd()
	NetworkZoneCmd.AddCommand(DeleteNetworkZoneCmd)
}

var (
	ActivateNetworkZonezoneId string

	ActivateNetworkZoneQuiet bool
)

func NewActivateNetworkZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Network Zone",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.NetworkZoneAPI.ActivateNetworkZone(apiClient.GetConfig().Context, ActivateNetworkZonezoneId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateNetworkZoneQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateNetworkZoneQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateNetworkZonezoneId, "zoneId", "", "", "")
	cmd.MarkFlagRequired("zoneId")

	cmd.Flags().BoolVarP(&ActivateNetworkZoneQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateNetworkZoneCmd := NewActivateNetworkZoneCmd()
	NetworkZoneCmd.AddCommand(ActivateNetworkZoneCmd)
}

var (
	DeactivateNetworkZonezoneId string

	DeactivateNetworkZoneQuiet bool
)

func NewDeactivateNetworkZoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Network Zone",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.NetworkZoneAPI.DeactivateNetworkZone(apiClient.GetConfig().Context, DeactivateNetworkZonezoneId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateNetworkZoneQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateNetworkZoneQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateNetworkZonezoneId, "zoneId", "", "", "")
	cmd.MarkFlagRequired("zoneId")

	cmd.Flags().BoolVarP(&DeactivateNetworkZoneQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateNetworkZoneCmd := NewDeactivateNetworkZoneCmd()
	NetworkZoneCmd.AddCommand(DeactivateNetworkZoneCmd)
}
