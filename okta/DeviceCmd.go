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

var DeviceCmd = &cobra.Command{
	Use:  "device",
	Long: "Manage DeviceAPI",
}

func init() {
	rootCmd.AddCommand(DeviceCmd)
}

var (
	ListDevicesBackupDir string

	ListDevicesLimit    int32
	ListDevicesPage     string
	ListDevicesFetchAll bool

	ListDevicesQuiet bool
)

func NewListDevicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Devices",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.ListDevices(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListDevicesQuiet {
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
					if !ListDevicesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListDevicesFetchAll || len(items) == 0 {
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

			if ListDevicesFetchAll && pageCount > 1 && !ListDevicesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListDevicesBackupDir, "device", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListDevicesQuiet {
					fmt.Printf("Backing up Devices to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListDevicesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListDevicesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListDevicesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListDevicesQuiet {
					fmt.Printf("Successfully backed up %d/%d Devices\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListDevicesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListDevicesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListDevicesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListDevicesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Devices to a directory")

	cmd.Flags().StringVarP(&ListDevicesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListDevicesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListDevicesCmd := NewListDevicesCmd()
	DeviceCmd.AddCommand(ListDevicesCmd)
}

var (
	GetDevicedeviceId string

	GetDeviceBackupDir string

	GetDeviceQuiet bool
)

func NewGetDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.GetDevice(apiClient.GetConfig().Context, GetDevicedeviceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetDeviceQuiet {
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
				dirPath := filepath.Join(GetDeviceBackupDir, "device", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetDevicedeviceId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetDeviceQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetDeviceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetDevicedeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Device to a file")

	cmd.Flags().StringVarP(&GetDeviceBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetDeviceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetDeviceCmd := NewGetDeviceCmd()
	DeviceCmd.AddCommand(GetDeviceCmd)
}

var (
	DeleteDevicedeviceId string

	DeleteDeviceQuiet bool
)

func NewDeleteDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.DeleteDevice(apiClient.GetConfig().Context, DeleteDevicedeviceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteDeviceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteDeviceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteDevicedeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().BoolVarP(&DeleteDeviceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteDeviceCmd := NewDeleteDeviceCmd()
	DeviceCmd.AddCommand(DeleteDeviceCmd)
}

var (
	ActivateDevicedeviceId string

	ActivateDeviceQuiet bool
)

func NewActivateDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.ActivateDevice(apiClient.GetConfig().Context, ActivateDevicedeviceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateDeviceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateDeviceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateDevicedeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().BoolVarP(&ActivateDeviceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateDeviceCmd := NewActivateDeviceCmd()
	DeviceCmd.AddCommand(ActivateDeviceCmd)
}

var (
	DeactivateDevicedeviceId string

	DeactivateDeviceQuiet bool
)

func NewDeactivateDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.DeactivateDevice(apiClient.GetConfig().Context, DeactivateDevicedeviceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateDeviceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateDeviceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateDevicedeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().BoolVarP(&DeactivateDeviceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateDeviceCmd := NewDeactivateDeviceCmd()
	DeviceCmd.AddCommand(DeactivateDeviceCmd)
}

var (
	SuspendDevicedeviceId string

	SuspendDeviceQuiet bool
)

func NewSuspendDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "suspend",
		Long: "Suspend a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.SuspendDevice(apiClient.GetConfig().Context, SuspendDevicedeviceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !SuspendDeviceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !SuspendDeviceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&SuspendDevicedeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().BoolVarP(&SuspendDeviceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	SuspendDeviceCmd := NewSuspendDeviceCmd()
	DeviceCmd.AddCommand(SuspendDeviceCmd)
}

var (
	UnsuspendDevicedeviceId string

	UnsuspendDeviceQuiet bool
)

func NewUnsuspendDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unsuspend",
		Long: "Unsuspend a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.UnsuspendDevice(apiClient.GetConfig().Context, UnsuspendDevicedeviceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnsuspendDeviceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnsuspendDeviceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnsuspendDevicedeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().BoolVarP(&UnsuspendDeviceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnsuspendDeviceCmd := NewUnsuspendDeviceCmd()
	DeviceCmd.AddCommand(UnsuspendDeviceCmd)
}

var (
	ListDeviceUsersdeviceId string

	ListDeviceUsersBackupDir string

	ListDeviceUsersLimit    int32
	ListDeviceUsersPage     string
	ListDeviceUsersFetchAll bool

	ListDeviceUsersQuiet bool
)

func NewListDeviceUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listUsers",
		Long: "List all Users for a Device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAPI.ListDeviceUsers(apiClient.GetConfig().Context, ListDeviceUsersdeviceId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListDeviceUsersQuiet {
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
					if !ListDeviceUsersQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListDeviceUsersFetchAll || len(items) == 0 {
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

			if ListDeviceUsersFetchAll && pageCount > 1 && !ListDeviceUsersQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListDeviceUsersBackupDir, "device", "listUsers")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListDeviceUsersQuiet {
					fmt.Printf("Backing up Devices to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListDeviceUsersQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListDeviceUsersQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListDeviceUsersQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListDeviceUsersQuiet {
					fmt.Printf("Successfully backed up %d/%d Devices\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListDeviceUsersQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListDeviceUsersdeviceId, "deviceId", "", "", "")
	cmd.MarkFlagRequired("deviceId")

	cmd.Flags().Int32VarP(&ListDeviceUsersLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListDeviceUsersPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListDeviceUsersFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Devices to a directory")

	cmd.Flags().StringVarP(&ListDeviceUsersBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListDeviceUsersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListDeviceUsersCmd := NewListDeviceUsersCmd()
	DeviceCmd.AddCommand(ListDeviceUsersCmd)
}
