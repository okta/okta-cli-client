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

var AuthorizationServerPoliciesCmd = &cobra.Command{
	Use:  "authorizationServerPolicies",
	Long: "Manage AuthorizationServerPoliciesAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerPoliciesCmd)
}

var (
	CreateAuthorizationServerPolicyauthServerId string

	CreateAuthorizationServerPolicydata string

	CreateAuthorizationServerPolicyRestoreFile string

	CreateAuthorizationServerPolicyQuiet bool
)

func NewCreateAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createAuthorizationServerPolicy",
		Long: "Create a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateAuthorizationServerPolicyRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateAuthorizationServerPolicyRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateAuthorizationServerPolicydata = string(processedData)

				if !CreateAuthorizationServerPolicyQuiet {
					fmt.Println("Restoring AuthorizationServerPolicies from:", CreateAuthorizationServerPolicyRestoreFile)
				}
			}

			req := apiClient.AuthorizationServerPoliciesAPI.CreateAuthorizationServerPolicy(apiClient.GetConfig().Context, CreateAuthorizationServerPolicyauthServerId)

			if CreateAuthorizationServerPolicydata != "" {
				req = req.Data(CreateAuthorizationServerPolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateAuthorizationServerPolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateAuthorizationServerPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyRestoreFile, "restore-from", "r", "", "Restore AuthorizationServerPolicies from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateAuthorizationServerPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateAuthorizationServerPolicyCmd := NewCreateAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(CreateAuthorizationServerPolicyCmd)
}

var (
	ListAuthorizationServerPoliciesauthServerId string

	ListAuthorizationServerPoliciesBackupDir string

	ListAuthorizationServerPoliciesLimit    int32
	ListAuthorizationServerPoliciesPage     string
	ListAuthorizationServerPoliciesFetchAll bool

	ListAuthorizationServerPoliciesQuiet bool
)

func NewListAuthorizationServerPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List all Policies",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.ListAuthorizationServerPolicies(apiClient.GetConfig().Context, ListAuthorizationServerPoliciesauthServerId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAuthorizationServerPoliciesQuiet {
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
					if !ListAuthorizationServerPoliciesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAuthorizationServerPoliciesFetchAll || len(items) == 0 {
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

			if ListAuthorizationServerPoliciesFetchAll && pageCount > 1 && !ListAuthorizationServerPoliciesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAuthorizationServerPoliciesBackupDir, "authorizationserverpolicies", "list")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAuthorizationServerPoliciesQuiet {
					fmt.Printf("Backing up AuthorizationServerPoliciess to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAuthorizationServerPoliciesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAuthorizationServerPoliciesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAuthorizationServerPoliciesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAuthorizationServerPoliciesQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerPoliciess\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAuthorizationServerPoliciesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAuthorizationServerPoliciesauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().Int32VarP(&ListAuthorizationServerPoliciesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAuthorizationServerPoliciesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAuthorizationServerPoliciesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerPoliciess to a directory")

	cmd.Flags().StringVarP(&ListAuthorizationServerPoliciesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAuthorizationServerPoliciesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAuthorizationServerPoliciesCmd := NewListAuthorizationServerPoliciesCmd()
	AuthorizationServerPoliciesCmd.AddCommand(ListAuthorizationServerPoliciesCmd)
}

var (
	GetAuthorizationServerPolicyauthServerId string

	GetAuthorizationServerPolicypolicyId string

	GetAuthorizationServerPolicyBackupDir string

	GetAuthorizationServerPolicyQuiet bool
)

func NewGetAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getAuthorizationServerPolicy",
		Long: "Retrieve a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.GetAuthorizationServerPolicy(apiClient.GetConfig().Context, GetAuthorizationServerPolicyauthServerId, GetAuthorizationServerPolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAuthorizationServerPolicyQuiet {
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
				dirPath := filepath.Join(GetAuthorizationServerPolicyBackupDir, "authorizationserverpolicies", "getAuthorizationServerPolicy")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAuthorizationServerPolicyauthServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAuthorizationServerPolicyQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAuthorizationServerPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AuthorizationServerPolicies to a file")

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAuthorizationServerPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAuthorizationServerPolicyCmd := NewGetAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(GetAuthorizationServerPolicyCmd)
}

var (
	ReplaceAuthorizationServerPolicyauthServerId string

	ReplaceAuthorizationServerPolicypolicyId string

	ReplaceAuthorizationServerPolicydata string

	ReplaceAuthorizationServerPolicyQuiet bool
)

func NewReplaceAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceAuthorizationServerPolicy",
		Long: "Replace a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.ReplaceAuthorizationServerPolicy(apiClient.GetConfig().Context, ReplaceAuthorizationServerPolicyauthServerId, ReplaceAuthorizationServerPolicypolicyId)

			if ReplaceAuthorizationServerPolicydata != "" {
				req = req.Data(ReplaceAuthorizationServerPolicydata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceAuthorizationServerPolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceAuthorizationServerPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceAuthorizationServerPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceAuthorizationServerPolicyCmd := NewReplaceAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(ReplaceAuthorizationServerPolicyCmd)
}

var (
	DeleteAuthorizationServerPolicyauthServerId string

	DeleteAuthorizationServerPolicypolicyId string

	DeleteAuthorizationServerPolicyQuiet bool
)

func NewDeleteAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteAuthorizationServerPolicy",
		Long: "Delete a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.DeleteAuthorizationServerPolicy(apiClient.GetConfig().Context, DeleteAuthorizationServerPolicyauthServerId, DeleteAuthorizationServerPolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteAuthorizationServerPolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteAuthorizationServerPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&DeleteAuthorizationServerPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteAuthorizationServerPolicyCmd := NewDeleteAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(DeleteAuthorizationServerPolicyCmd)
}

var (
	ActivateAuthorizationServerPolicyauthServerId string

	ActivateAuthorizationServerPolicypolicyId string

	ActivateAuthorizationServerPolicyQuiet bool
)

func NewActivateAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateAuthorizationServerPolicy",
		Long: "Activate a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.ActivateAuthorizationServerPolicy(apiClient.GetConfig().Context, ActivateAuthorizationServerPolicyauthServerId, ActivateAuthorizationServerPolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateAuthorizationServerPolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateAuthorizationServerPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&ActivateAuthorizationServerPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateAuthorizationServerPolicyCmd := NewActivateAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(ActivateAuthorizationServerPolicyCmd)
}

var (
	DeactivateAuthorizationServerPolicyauthServerId string

	DeactivateAuthorizationServerPolicypolicyId string

	DeactivateAuthorizationServerPolicyQuiet bool
)

func NewDeactivateAuthorizationServerPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateAuthorizationServerPolicy",
		Long: "Deactivate a Policy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerPoliciesAPI.DeactivateAuthorizationServerPolicy(apiClient.GetConfig().Context, DeactivateAuthorizationServerPolicyauthServerId, DeactivateAuthorizationServerPolicypolicyId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateAuthorizationServerPolicyQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateAuthorizationServerPolicyQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicyauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().BoolVarP(&DeactivateAuthorizationServerPolicyQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateAuthorizationServerPolicyCmd := NewDeactivateAuthorizationServerPolicyCmd()
	AuthorizationServerPoliciesCmd.AddCommand(DeactivateAuthorizationServerPolicyCmd)
}
