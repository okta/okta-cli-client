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

var PolicyCmd = &cobra.Command{
	Use:  "policy",
	Long: "Manage PolicyAPI",
}

func init() {
	rootCmd.AddCommand(PolicyCmd)
}

var (
	CreatePolicydata string

	CreatePolicyRestoreFile string

	CreatePolicyQuiet bool
)

func NewCreatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreatePolicyRestoreFile != "" {

				jsonData, err := os.ReadFile(CreatePolicyRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreatePolicydata = string(processedData)

				if !CreatePolicyQuiet {
					fmt.Println("Restoring Policy from:", CreatePolicyRestoreFile)
				}
			}

			req := apiClient.PolicyAPI.CreatePolicy(apiClient.GetConfig().Context)

			if CreatePolicydata != "" {
				req = req.Data(CreatePolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreatePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreatePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreatePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreatePolicyRestoreFile, "restore-from", "r", "", "Restore Policy from a JSON backup file")

	cmd.Flags().BoolVarP(&CreatePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreatePolicyCmd := NewCreatePolicyCmd()
	PolicyCmd.AddCommand(CreatePolicyCmd)
}

var (
	ListPoliciesBackupDir string

	ListPoliciesLimit    int32
	ListPoliciesPage     string
	ListPoliciesFetchAll bool

	ListPoliciesQuiet bool
)

func NewListPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listPolicies",
		Long: "List all Policies",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicies(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListPoliciesQuiet {
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
					if !ListPoliciesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListPoliciesFetchAll || len(items) == 0 {
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

			if ListPoliciesFetchAll && pageCount > 1 && !ListPoliciesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListPoliciesBackupDir, "policy", "listPolicies")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListPoliciesQuiet {
					fmt.Printf("Backing up Policys to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListPoliciesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListPoliciesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListPoliciesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListPoliciesQuiet {
					fmt.Printf("Successfully backed up %d/%d Policys\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListPoliciesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListPoliciesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListPoliciesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListPoliciesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Policys to a directory")

	cmd.Flags().StringVarP(&ListPoliciesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListPoliciesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListPoliciesCmd := NewListPoliciesCmd()
	PolicyCmd.AddCommand(ListPoliciesCmd)
}

var (
	CreatePolicySimulationdata string

	CreatePolicySimulationRestoreFile string

	CreatePolicySimulationQuiet bool
)

func NewCreatePolicySimulationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createSimulation",
		Long: "Create a Policy Simulation",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreatePolicySimulationRestoreFile != "" {

				jsonData, err := os.ReadFile(CreatePolicySimulationRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreatePolicySimulationdata = string(processedData)

				if !CreatePolicySimulationQuiet {
					fmt.Println("Restoring Policy from:", CreatePolicySimulationRestoreFile)
				}
			}

			req := apiClient.PolicyAPI.CreatePolicySimulation(apiClient.GetConfig().Context)

			if CreatePolicySimulationdata != "" {
				req = req.Data(CreatePolicySimulationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreatePolicySimulationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreatePolicySimulationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreatePolicySimulationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreatePolicySimulationRestoreFile, "restore-from", "r", "", "Restore Policy from a JSON backup file")

	cmd.Flags().BoolVarP(&CreatePolicySimulationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreatePolicySimulationCmd := NewCreatePolicySimulationCmd()
	PolicyCmd.AddCommand(CreatePolicySimulationCmd)
}

var (
	GetPolicypolicyId string

	GetPolicyBackupDir string

	GetPolicyQuiet bool
)

func NewGetPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.GetPolicy(apiClient.GetConfig().Context, GetPolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPolicyQuiet {
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
				dirPath := filepath.Join(GetPolicyBackupDir, "policy", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPolicypolicyId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPolicyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Policy to a file")

	cmd.Flags().StringVarP(&GetPolicyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPolicyCmd := NewGetPolicyCmd()
	PolicyCmd.AddCommand(GetPolicyCmd)
}

var (
	ReplacePolicypolicyId string

	ReplacePolicydata string

	ReplacePolicyQuiet bool
)

func NewReplacePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ReplacePolicy(apiClient.GetConfig().Context, ReplacePolicypolicyId)

			if ReplacePolicydata != "" {
				req = req.Data(ReplacePolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplacePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacePolicyCmd := NewReplacePolicyCmd()
	PolicyCmd.AddCommand(ReplacePolicyCmd)
}

var (
	DeletePolicypolicyId string

	DeletePolicyQuiet bool
)

func NewDeletePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeletePolicy(apiClient.GetConfig().Context, DeletePolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeletePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeletePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeletePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&DeletePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeletePolicyCmd := NewDeletePolicyCmd()
	PolicyCmd.AddCommand(DeletePolicyCmd)
}

var (
	ListPolicyAppspolicyId string

	ListPolicyAppsBackupDir string

	ListPolicyAppsLimit    int32
	ListPolicyAppsPage     string
	ListPolicyAppsFetchAll bool

	ListPolicyAppsQuiet bool
)

func NewListPolicyAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApps",
		Long: "List all Applications mapped to a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicyApps(apiClient.GetConfig().Context, ListPolicyAppspolicyId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListPolicyAppsQuiet {
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
					if !ListPolicyAppsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListPolicyAppsFetchAll || len(items) == 0 {
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

			if ListPolicyAppsFetchAll && pageCount > 1 && !ListPolicyAppsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListPolicyAppsBackupDir, "policy", "listApps")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListPolicyAppsQuiet {
					fmt.Printf("Backing up Policys to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListPolicyAppsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListPolicyAppsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListPolicyAppsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListPolicyAppsQuiet {
					fmt.Printf("Successfully backed up %d/%d Policys\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListPolicyAppsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListPolicyAppspolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().Int32VarP(&ListPolicyAppsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListPolicyAppsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListPolicyAppsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Policys to a directory")

	cmd.Flags().StringVarP(&ListPolicyAppsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListPolicyAppsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListPolicyAppsCmd := NewListPolicyAppsCmd()
	PolicyCmd.AddCommand(ListPolicyAppsCmd)
}

var (
	ClonePolicypolicyId string

	ClonePolicyQuiet bool
)

func NewClonePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "clone",
		Long: "Clone an existing Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ClonePolicy(apiClient.GetConfig().Context, ClonePolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ClonePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ClonePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ClonePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&ClonePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ClonePolicyCmd := NewClonePolicyCmd()
	PolicyCmd.AddCommand(ClonePolicyCmd)
}

var (
	ActivatePolicypolicyId string

	ActivatePolicyQuiet bool
)

func NewActivatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ActivatePolicy(apiClient.GetConfig().Context, ActivatePolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivatePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivatePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivatePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&ActivatePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivatePolicyCmd := NewActivatePolicyCmd()
	PolicyCmd.AddCommand(ActivatePolicyCmd)
}

var (
	DeactivatePolicypolicyId string

	DeactivatePolicyQuiet bool
)

func NewDeactivatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeactivatePolicy(apiClient.GetConfig().Context, DeactivatePolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivatePolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivatePolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivatePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&DeactivatePolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivatePolicyCmd := NewDeactivatePolicyCmd()
	PolicyCmd.AddCommand(DeactivatePolicyCmd)
}

var (
	MapResourceToPolicypolicyId string

	MapResourceToPolicydata string

	MapResourceToPolicyQuiet bool
)

func NewMapResourceToPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mapResourceTo",
		Long: "Map a resource to a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.MapResourceToPolicy(apiClient.GetConfig().Context, MapResourceToPolicypolicyId)

			if MapResourceToPolicydata != "" {
				req = req.Data(MapResourceToPolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !MapResourceToPolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !MapResourceToPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&MapResourceToPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&MapResourceToPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&MapResourceToPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	MapResourceToPolicyCmd := NewMapResourceToPolicyCmd()
	PolicyCmd.AddCommand(MapResourceToPolicyCmd)
}

var (
	ListPolicyMappingspolicyId string

	ListPolicyMappingsBackupDir string

	ListPolicyMappingsLimit    int32
	ListPolicyMappingsPage     string
	ListPolicyMappingsFetchAll bool

	ListPolicyMappingsQuiet bool
)

func NewListPolicyMappingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listMappings",
		Long: "List all resources mapped to a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicyMappings(apiClient.GetConfig().Context, ListPolicyMappingspolicyId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListPolicyMappingsQuiet {
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
					if !ListPolicyMappingsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListPolicyMappingsFetchAll || len(items) == 0 {
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

			if ListPolicyMappingsFetchAll && pageCount > 1 && !ListPolicyMappingsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListPolicyMappingsBackupDir, "policy", "listMappings")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListPolicyMappingsQuiet {
					fmt.Printf("Backing up Policys to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListPolicyMappingsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListPolicyMappingsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListPolicyMappingsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListPolicyMappingsQuiet {
					fmt.Printf("Successfully backed up %d/%d Policys\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListPolicyMappingsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListPolicyMappingspolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().Int32VarP(&ListPolicyMappingsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListPolicyMappingsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListPolicyMappingsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Policys to a directory")

	cmd.Flags().StringVarP(&ListPolicyMappingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListPolicyMappingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListPolicyMappingsCmd := NewListPolicyMappingsCmd()
	PolicyCmd.AddCommand(ListPolicyMappingsCmd)
}

var (
	GetPolicyMappingpolicyId string

	GetPolicyMappingmappingId string

	GetPolicyMappingBackupDir string

	GetPolicyMappingQuiet bool
)

func NewGetPolicyMappingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getMapping",
		Long: "Retrieve a policy resource Mapping",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.GetPolicyMapping(apiClient.GetConfig().Context, GetPolicyMappingpolicyId, GetPolicyMappingmappingId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPolicyMappingQuiet {
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
				dirPath := filepath.Join(GetPolicyMappingBackupDir, "policy", "getMapping")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPolicyMappingpolicyId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPolicyMappingQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPolicyMappingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPolicyMappingpolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&GetPolicyMappingmappingId, "mappingId", "", "", "")
	cmd.MarkFlagRequired("mappingId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Policy to a file")

	cmd.Flags().StringVarP(&GetPolicyMappingBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPolicyMappingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPolicyMappingCmd := NewGetPolicyMappingCmd()
	PolicyCmd.AddCommand(GetPolicyMappingCmd)
}

var (
	DeletePolicyResourceMappingpolicyId string

	DeletePolicyResourceMappingmappingId string

	DeletePolicyResourceMappingQuiet bool
)

func NewDeletePolicyResourceMappingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteResourceMapping",
		Long: "Delete a policy resource Mapping",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeletePolicyResourceMapping(apiClient.GetConfig().Context, DeletePolicyResourceMappingpolicyId, DeletePolicyResourceMappingmappingId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeletePolicyResourceMappingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeletePolicyResourceMappingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeletePolicyResourceMappingpolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeletePolicyResourceMappingmappingId, "mappingId", "", "", "")
	cmd.MarkFlagRequired("mappingId")

	cmd.Flags().BoolVarP(&DeletePolicyResourceMappingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeletePolicyResourceMappingCmd := NewDeletePolicyResourceMappingCmd()
	PolicyCmd.AddCommand(DeletePolicyResourceMappingCmd)
}

var (
	CreatePolicyRulepolicyId string

	CreatePolicyRuledata string

	CreatePolicyRuleRestoreFile string

	CreatePolicyRuleQuiet bool
)

func NewCreatePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createRule",
		Long: "Create a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreatePolicyRuleRestoreFile != "" {

				jsonData, err := os.ReadFile(CreatePolicyRuleRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreatePolicyRuledata = string(processedData)

				if !CreatePolicyRuleQuiet {
					fmt.Println("Restoring Policy from:", CreatePolicyRuleRestoreFile)
				}
			}

			req := apiClient.PolicyAPI.CreatePolicyRule(apiClient.GetConfig().Context, CreatePolicyRulepolicyId)

			if CreatePolicyRuledata != "" {
				req = req.Data(CreatePolicyRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreatePolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreatePolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreatePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&CreatePolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreatePolicyRuleRestoreFile, "restore-from", "r", "", "Restore Policy from a JSON backup file")

	cmd.Flags().BoolVarP(&CreatePolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreatePolicyRuleCmd := NewCreatePolicyRuleCmd()
	PolicyCmd.AddCommand(CreatePolicyRuleCmd)
}

var (
	ListPolicyRulespolicyId string

	ListPolicyRulesBackupDir string

	ListPolicyRulesLimit    int32
	ListPolicyRulesPage     string
	ListPolicyRulesFetchAll bool

	ListPolicyRulesQuiet bool
)

func NewListPolicyRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listRules",
		Long: "List all Policy Rules",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicyRules(apiClient.GetConfig().Context, ListPolicyRulespolicyId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListPolicyRulesQuiet {
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
					if !ListPolicyRulesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListPolicyRulesFetchAll || len(items) == 0 {
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

			if ListPolicyRulesFetchAll && pageCount > 1 && !ListPolicyRulesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListPolicyRulesBackupDir, "policy", "listRules")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListPolicyRulesQuiet {
					fmt.Printf("Backing up Policys to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListPolicyRulesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListPolicyRulesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListPolicyRulesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListPolicyRulesQuiet {
					fmt.Printf("Successfully backed up %d/%d Policys\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListPolicyRulesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListPolicyRulespolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().Int32VarP(&ListPolicyRulesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListPolicyRulesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListPolicyRulesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Policys to a directory")

	cmd.Flags().StringVarP(&ListPolicyRulesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListPolicyRulesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListPolicyRulesCmd := NewListPolicyRulesCmd()
	PolicyCmd.AddCommand(ListPolicyRulesCmd)
}

var (
	GetPolicyRulepolicyId string

	GetPolicyRuleruleId string

	GetPolicyRuleBackupDir string

	GetPolicyRuleQuiet bool
)

func NewGetPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getRule",
		Long: "Retrieve a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.GetPolicyRule(apiClient.GetConfig().Context, GetPolicyRulepolicyId, GetPolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPolicyRuleQuiet {
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
				dirPath := filepath.Join(GetPolicyRuleBackupDir, "policy", "getRule")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPolicyRulepolicyId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPolicyRuleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&GetPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Policy to a file")

	cmd.Flags().StringVarP(&GetPolicyRuleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPolicyRuleCmd := NewGetPolicyRuleCmd()
	PolicyCmd.AddCommand(GetPolicyRuleCmd)
}

var (
	ReplacePolicyRulepolicyId string

	ReplacePolicyRuleruleId string

	ReplacePolicyRuledata string

	ReplacePolicyRuleQuiet bool
)

func NewReplacePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceRule",
		Long: "Replace a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ReplacePolicyRule(apiClient.GetConfig().Context, ReplacePolicyRulepolicyId, ReplacePolicyRuleruleId)

			if ReplacePolicyRuledata != "" {
				req = req.Data(ReplacePolicyRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacePolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacePolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplacePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().StringVarP(&ReplacePolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacePolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacePolicyRuleCmd := NewReplacePolicyRuleCmd()
	PolicyCmd.AddCommand(ReplacePolicyRuleCmd)
}

var (
	DeletePolicyRulepolicyId string

	DeletePolicyRuleruleId string

	DeletePolicyRuleQuiet bool
)

func NewDeletePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteRule",
		Long: "Delete a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeletePolicyRule(apiClient.GetConfig().Context, DeletePolicyRulepolicyId, DeletePolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeletePolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeletePolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeletePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeletePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolVarP(&DeletePolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeletePolicyRuleCmd := NewDeletePolicyRuleCmd()
	PolicyCmd.AddCommand(DeletePolicyRuleCmd)
}

var (
	ActivatePolicyRulepolicyId string

	ActivatePolicyRuleruleId string

	ActivatePolicyRuleQuiet bool
)

func NewActivatePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateRule",
		Long: "Activate a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ActivatePolicyRule(apiClient.GetConfig().Context, ActivatePolicyRulepolicyId, ActivatePolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivatePolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivatePolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivatePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ActivatePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolVarP(&ActivatePolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivatePolicyRuleCmd := NewActivatePolicyRuleCmd()
	PolicyCmd.AddCommand(ActivatePolicyRuleCmd)
}

var (
	DeactivatePolicyRulepolicyId string

	DeactivatePolicyRuleruleId string

	DeactivatePolicyRuleQuiet bool
)

func NewDeactivatePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateRule",
		Long: "Deactivate a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeactivatePolicyRule(apiClient.GetConfig().Context, DeactivatePolicyRulepolicyId, DeactivatePolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivatePolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivatePolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivatePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeactivatePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolVarP(&DeactivatePolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivatePolicyRuleCmd := NewDeactivatePolicyRuleCmd()
	PolicyCmd.AddCommand(DeactivatePolicyRuleCmd)
}
