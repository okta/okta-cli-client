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

var DeviceAssuranceCmd = &cobra.Command{
	Use:  "deviceAssurance",
	Long: "Manage DeviceAssuranceAPI",
}

func init() {
	rootCmd.AddCommand(DeviceAssuranceCmd)
}

var (
	CreateDeviceAssurancePolicydata string

	CreateDeviceAssurancePolicyRestoreFile string

	CreateDeviceAssurancePolicyQuiet bool
)

func NewCreateDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createPolicy",
		Long: "Create a Device Assurance Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateDeviceAssurancePolicyRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateDeviceAssurancePolicyRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateDeviceAssurancePolicydata = string(processedData)

				if !CreateDeviceAssurancePolicyQuiet {
					fmt.Println("Restoring DeviceAssurance from:", CreateDeviceAssurancePolicyRestoreFile)
				}
			}

			req := apiClient.DeviceAssuranceAPI.CreateDeviceAssurancePolicy(apiClient.GetConfig().Context)

			if CreateDeviceAssurancePolicydata != "" {
				req = req.Data(CreateDeviceAssurancePolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateDeviceAssurancePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateDeviceAssurancePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateDeviceAssurancePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateDeviceAssurancePolicyRestoreFile, "restore-from", "r", "", "Restore DeviceAssurance from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateDeviceAssurancePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateDeviceAssurancePolicyCmd := NewCreateDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(CreateDeviceAssurancePolicyCmd)
}

var (
	ListDeviceAssurancePoliciesBackupDir string

	ListDeviceAssurancePoliciesLimit    int32
	ListDeviceAssurancePoliciesPage     string
	ListDeviceAssurancePoliciesFetchAll bool

	ListDeviceAssurancePoliciesQuiet bool
)

func NewListDeviceAssurancePoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listPolicies",
		Long: "List all Device Assurance Policies",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.ListDeviceAssurancePolicies(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListDeviceAssurancePoliciesQuiet {
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
					if !ListDeviceAssurancePoliciesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListDeviceAssurancePoliciesFetchAll || len(items) == 0 {
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

			if ListDeviceAssurancePoliciesFetchAll && pageCount > 1 && !ListDeviceAssurancePoliciesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListDeviceAssurancePoliciesBackupDir, "deviceassurance", "listPolicies")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListDeviceAssurancePoliciesQuiet {
					fmt.Printf("Backing up DeviceAssurances to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListDeviceAssurancePoliciesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListDeviceAssurancePoliciesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListDeviceAssurancePoliciesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListDeviceAssurancePoliciesQuiet {
					fmt.Printf("Successfully backed up %d/%d DeviceAssurances\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListDeviceAssurancePoliciesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListDeviceAssurancePoliciesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListDeviceAssurancePoliciesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListDeviceAssurancePoliciesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple DeviceAssurances to a directory")

	cmd.Flags().StringVarP(&ListDeviceAssurancePoliciesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListDeviceAssurancePoliciesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListDeviceAssurancePoliciesCmd := NewListDeviceAssurancePoliciesCmd()
	DeviceAssuranceCmd.AddCommand(ListDeviceAssurancePoliciesCmd)
}

var (
	GetDeviceAssurancePolicydeviceAssuranceId string

	GetDeviceAssurancePolicyBackupDir string

	GetDeviceAssurancePolicyQuiet bool
)

func NewGetDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPolicy",
		Long: "Retrieve a Device Assurance Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.GetDeviceAssurancePolicy(apiClient.GetConfig().Context, GetDeviceAssurancePolicydeviceAssuranceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetDeviceAssurancePolicyQuiet {
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
				dirPath := filepath.Join(GetDeviceAssurancePolicyBackupDir, "deviceassurance", "getPolicy")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetDeviceAssurancePolicydeviceAssuranceId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetDeviceAssurancePolicyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetDeviceAssurancePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetDeviceAssurancePolicydeviceAssuranceId, "deviceAssuranceId", "", "", "")
	cmd.MarkFlagRequired("deviceAssuranceId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the DeviceAssurance to a file")

	cmd.Flags().StringVarP(&GetDeviceAssurancePolicyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetDeviceAssurancePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetDeviceAssurancePolicyCmd := NewGetDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(GetDeviceAssurancePolicyCmd)
}

var (
	ReplaceDeviceAssurancePolicydeviceAssuranceId string

	ReplaceDeviceAssurancePolicydata string

	ReplaceDeviceAssurancePolicyQuiet bool
)

func NewReplaceDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePolicy",
		Long: "Replace a Device Assurance Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.ReplaceDeviceAssurancePolicy(apiClient.GetConfig().Context, ReplaceDeviceAssurancePolicydeviceAssuranceId)

			if ReplaceDeviceAssurancePolicydata != "" {
				req = req.Data(ReplaceDeviceAssurancePolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceDeviceAssurancePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceDeviceAssurancePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceDeviceAssurancePolicydeviceAssuranceId, "deviceAssuranceId", "", "", "")
	cmd.MarkFlagRequired("deviceAssuranceId")

	cmd.Flags().StringVarP(&ReplaceDeviceAssurancePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceDeviceAssurancePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceDeviceAssurancePolicyCmd := NewReplaceDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(ReplaceDeviceAssurancePolicyCmd)
}

var (
	DeleteDeviceAssurancePolicydeviceAssuranceId string

	DeleteDeviceAssurancePolicyQuiet bool
)

func NewDeleteDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletePolicy",
		Long: "Delete a Device Assurance Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.DeleteDeviceAssurancePolicy(apiClient.GetConfig().Context, DeleteDeviceAssurancePolicydeviceAssuranceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteDeviceAssurancePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteDeviceAssurancePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteDeviceAssurancePolicydeviceAssuranceId, "deviceAssuranceId", "", "", "")
	cmd.MarkFlagRequired("deviceAssuranceId")

	cmd.Flags().BoolVarP(&DeleteDeviceAssurancePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteDeviceAssurancePolicyCmd := NewDeleteDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(DeleteDeviceAssurancePolicyCmd)
}
