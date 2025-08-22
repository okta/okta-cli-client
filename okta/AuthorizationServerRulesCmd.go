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

var AuthorizationServerRulesCmd = &cobra.Command{
	Use:  "authorizationServerRules",
	Long: "Manage AuthorizationServerRulesAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerRulesCmd)
}

var (
	CreateAuthorizationServerPolicyRuleauthServerId string

	CreateAuthorizationServerPolicyRulepolicyId string

	CreateAuthorizationServerPolicyRuledata string

	CreateAuthorizationServerPolicyRuleRestoreFile string

	CreateAuthorizationServerPolicyRuleQuiet bool
)

func NewCreateAuthorizationServerPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createAuthorizationServerPolicyRule",
		Long: "Create a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateAuthorizationServerPolicyRuleRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateAuthorizationServerPolicyRuleRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateAuthorizationServerPolicyRuledata = string(processedData)

				if !CreateAuthorizationServerPolicyRuleQuiet {
					fmt.Println("Restoring AuthorizationServerRules from:", CreateAuthorizationServerPolicyRuleRestoreFile)
				}
			}

			req := apiClient.AuthorizationServerRulesAPI.CreateAuthorizationServerPolicyRule(apiClient.GetConfig().Context, CreateAuthorizationServerPolicyRuleauthServerId, CreateAuthorizationServerPolicyRulepolicyId)

			if CreateAuthorizationServerPolicyRuledata != "" {
				req = req.Data(CreateAuthorizationServerPolicyRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateAuthorizationServerPolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateAuthorizationServerPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyRuleauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateAuthorizationServerPolicyRuleRestoreFile, "restore-from", "r", "", "Restore AuthorizationServerRules from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateAuthorizationServerPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateAuthorizationServerPolicyRuleCmd := NewCreateAuthorizationServerPolicyRuleCmd()
	AuthorizationServerRulesCmd.AddCommand(CreateAuthorizationServerPolicyRuleCmd)
}

var (
	ListAuthorizationServerPolicyRulesauthServerId string

	ListAuthorizationServerPolicyRulespolicyId string

	ListAuthorizationServerPolicyRulesBackupDir string

	ListAuthorizationServerPolicyRulesLimit    int32
	ListAuthorizationServerPolicyRulesPage     string
	ListAuthorizationServerPolicyRulesFetchAll bool

	ListAuthorizationServerPolicyRulesQuiet bool
)

func NewListAuthorizationServerPolicyRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAuthorizationServerPolicyRules",
		Long: "List all Policy Rules",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerRulesAPI.ListAuthorizationServerPolicyRules(apiClient.GetConfig().Context, ListAuthorizationServerPolicyRulesauthServerId, ListAuthorizationServerPolicyRulespolicyId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAuthorizationServerPolicyRulesQuiet {
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
					if !ListAuthorizationServerPolicyRulesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAuthorizationServerPolicyRulesFetchAll || len(items) == 0 {
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

			if ListAuthorizationServerPolicyRulesFetchAll && pageCount > 1 && !ListAuthorizationServerPolicyRulesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAuthorizationServerPolicyRulesBackupDir, "authorizationserverrules", "listAuthorizationServerPolicyRules")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAuthorizationServerPolicyRulesQuiet {
					fmt.Printf("Backing up AuthorizationServerRuless to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAuthorizationServerPolicyRulesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAuthorizationServerPolicyRulesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAuthorizationServerPolicyRulesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAuthorizationServerPolicyRulesQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerRuless\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAuthorizationServerPolicyRulesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAuthorizationServerPolicyRulesauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ListAuthorizationServerPolicyRulespolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().Int32VarP(&ListAuthorizationServerPolicyRulesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAuthorizationServerPolicyRulesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAuthorizationServerPolicyRulesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerRuless to a directory")

	cmd.Flags().StringVarP(&ListAuthorizationServerPolicyRulesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAuthorizationServerPolicyRulesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAuthorizationServerPolicyRulesCmd := NewListAuthorizationServerPolicyRulesCmd()
	AuthorizationServerRulesCmd.AddCommand(ListAuthorizationServerPolicyRulesCmd)
}

var (
	GetAuthorizationServerPolicyRuleauthServerId string

	GetAuthorizationServerPolicyRulepolicyId string

	GetAuthorizationServerPolicyRuleruleId string

	GetAuthorizationServerPolicyRuleBackupDir string

	GetAuthorizationServerPolicyRuleQuiet bool
)

func NewGetAuthorizationServerPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getAuthorizationServerPolicyRule",
		Long: "Retrieve a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerRulesAPI.GetAuthorizationServerPolicyRule(apiClient.GetConfig().Context, GetAuthorizationServerPolicyRuleauthServerId, GetAuthorizationServerPolicyRulepolicyId, GetAuthorizationServerPolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetAuthorizationServerPolicyRuleQuiet {
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
				dirPath := filepath.Join(GetAuthorizationServerPolicyRuleBackupDir, "authorizationserverrules", "getAuthorizationServerPolicyRule")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetAuthorizationServerPolicyRuleauthServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetAuthorizationServerPolicyRuleQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetAuthorizationServerPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyRuleauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AuthorizationServerRules to a file")

	cmd.Flags().StringVarP(&GetAuthorizationServerPolicyRuleBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetAuthorizationServerPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetAuthorizationServerPolicyRuleCmd := NewGetAuthorizationServerPolicyRuleCmd()
	AuthorizationServerRulesCmd.AddCommand(GetAuthorizationServerPolicyRuleCmd)
}

var (
	ReplaceAuthorizationServerPolicyRuleauthServerId string

	ReplaceAuthorizationServerPolicyRulepolicyId string

	ReplaceAuthorizationServerPolicyRuleruleId string

	ReplaceAuthorizationServerPolicyRuledata string

	ReplaceAuthorizationServerPolicyRuleQuiet bool
)

func NewReplaceAuthorizationServerPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceAuthorizationServerPolicyRule",
		Long: "Replace a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerRulesAPI.ReplaceAuthorizationServerPolicyRule(apiClient.GetConfig().Context, ReplaceAuthorizationServerPolicyRuleauthServerId, ReplaceAuthorizationServerPolicyRulepolicyId, ReplaceAuthorizationServerPolicyRuleruleId)

			if ReplaceAuthorizationServerPolicyRuledata != "" {
				req = req.Data(ReplaceAuthorizationServerPolicyRuledata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceAuthorizationServerPolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceAuthorizationServerPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicyRuleauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().StringVarP(&ReplaceAuthorizationServerPolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceAuthorizationServerPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceAuthorizationServerPolicyRuleCmd := NewReplaceAuthorizationServerPolicyRuleCmd()
	AuthorizationServerRulesCmd.AddCommand(ReplaceAuthorizationServerPolicyRuleCmd)
}

var (
	DeleteAuthorizationServerPolicyRuleauthServerId string

	DeleteAuthorizationServerPolicyRulepolicyId string

	DeleteAuthorizationServerPolicyRuleruleId string

	DeleteAuthorizationServerPolicyRuleQuiet bool
)

func NewDeleteAuthorizationServerPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteAuthorizationServerPolicyRule",
		Long: "Delete a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerRulesAPI.DeleteAuthorizationServerPolicyRule(apiClient.GetConfig().Context, DeleteAuthorizationServerPolicyRuleauthServerId, DeleteAuthorizationServerPolicyRulepolicyId, DeleteAuthorizationServerPolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteAuthorizationServerPolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteAuthorizationServerPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicyRuleauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeleteAuthorizationServerPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolVarP(&DeleteAuthorizationServerPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteAuthorizationServerPolicyRuleCmd := NewDeleteAuthorizationServerPolicyRuleCmd()
	AuthorizationServerRulesCmd.AddCommand(DeleteAuthorizationServerPolicyRuleCmd)
}

var (
	ActivateAuthorizationServerPolicyRuleauthServerId string

	ActivateAuthorizationServerPolicyRulepolicyId string

	ActivateAuthorizationServerPolicyRuleruleId string

	ActivateAuthorizationServerPolicyRuleQuiet bool
)

func NewActivateAuthorizationServerPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateAuthorizationServerPolicyRule",
		Long: "Activate a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerRulesAPI.ActivateAuthorizationServerPolicyRule(apiClient.GetConfig().Context, ActivateAuthorizationServerPolicyRuleauthServerId, ActivateAuthorizationServerPolicyRulepolicyId, ActivateAuthorizationServerPolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateAuthorizationServerPolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateAuthorizationServerPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicyRuleauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ActivateAuthorizationServerPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolVarP(&ActivateAuthorizationServerPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateAuthorizationServerPolicyRuleCmd := NewActivateAuthorizationServerPolicyRuleCmd()
	AuthorizationServerRulesCmd.AddCommand(ActivateAuthorizationServerPolicyRuleCmd)
}

var (
	DeactivateAuthorizationServerPolicyRuleauthServerId string

	DeactivateAuthorizationServerPolicyRulepolicyId string

	DeactivateAuthorizationServerPolicyRuleruleId string

	DeactivateAuthorizationServerPolicyRuleQuiet bool
)

func NewDeactivateAuthorizationServerPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateAuthorizationServerPolicyRule",
		Long: "Deactivate a Policy Rule",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerRulesAPI.DeactivateAuthorizationServerPolicyRule(apiClient.GetConfig().Context, DeactivateAuthorizationServerPolicyRuleauthServerId, DeactivateAuthorizationServerPolicyRulepolicyId, DeactivateAuthorizationServerPolicyRuleruleId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeactivateAuthorizationServerPolicyRuleQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeactivateAuthorizationServerPolicyRuleQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicyRuleauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeactivateAuthorizationServerPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().BoolVarP(&DeactivateAuthorizationServerPolicyRuleQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeactivateAuthorizationServerPolicyRuleCmd := NewDeactivateAuthorizationServerPolicyRuleCmd()
	AuthorizationServerRulesCmd.AddCommand(DeactivateAuthorizationServerPolicyRuleCmd)
}
