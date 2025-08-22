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

var AuthorizationServerClaimsCmd = &cobra.Command{
	Use:  "authorizationServerClaims",
	Long: "Manage AuthorizationServerClaimsAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerClaimsCmd)
}

var (
	CreateOAuth2ClaimauthServerId string

	CreateOAuth2Claimdata string

	CreateOAuth2ClaimRestoreFile string

	CreateOAuth2ClaimQuiet bool
)

func NewCreateOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createOAuth2Claim",
		Long: "Create a custom token Claim",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateOAuth2ClaimRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateOAuth2ClaimRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateOAuth2Claimdata = string(processedData)

				if !CreateOAuth2ClaimQuiet {
					fmt.Println("Restoring AuthorizationServerClaims from:", CreateOAuth2ClaimRestoreFile)
				}
			}

			req := apiClient.AuthorizationServerClaimsAPI.CreateOAuth2Claim(apiClient.GetConfig().Context, CreateOAuth2ClaimauthServerId)

			if CreateOAuth2Claimdata != "" {
				req = req.Data(CreateOAuth2Claimdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateOAuth2ClaimQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateOAuth2ClaimQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateOAuth2Claimdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateOAuth2ClaimRestoreFile, "restore-from", "r", "", "Restore AuthorizationServerClaims from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateOAuth2ClaimQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateOAuth2ClaimCmd := NewCreateOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(CreateOAuth2ClaimCmd)
}

var (
	ListOAuth2ClaimsauthServerId string

	ListOAuth2ClaimsBackupDir string

	ListOAuth2ClaimsLimit    int32
	ListOAuth2ClaimsPage     string
	ListOAuth2ClaimsFetchAll bool

	ListOAuth2ClaimsQuiet bool
)

func NewListOAuth2ClaimsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listOAuth2Claims",
		Long: "List all custom token Claims",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.ListOAuth2Claims(apiClient.GetConfig().Context, ListOAuth2ClaimsauthServerId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListOAuth2ClaimsQuiet {
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
					if !ListOAuth2ClaimsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListOAuth2ClaimsFetchAll || len(items) == 0 {
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

			if ListOAuth2ClaimsFetchAll && pageCount > 1 && !ListOAuth2ClaimsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListOAuth2ClaimsBackupDir, "authorizationserverclaims", "listOAuth2Claims")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListOAuth2ClaimsQuiet {
					fmt.Printf("Backing up AuthorizationServerClaimss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListOAuth2ClaimsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListOAuth2ClaimsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListOAuth2ClaimsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListOAuth2ClaimsQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerClaimss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListOAuth2ClaimsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListOAuth2ClaimsauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().Int32VarP(&ListOAuth2ClaimsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListOAuth2ClaimsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListOAuth2ClaimsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerClaimss to a directory")

	cmd.Flags().StringVarP(&ListOAuth2ClaimsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListOAuth2ClaimsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListOAuth2ClaimsCmd := NewListOAuth2ClaimsCmd()
	AuthorizationServerClaimsCmd.AddCommand(ListOAuth2ClaimsCmd)
}

var (
	GetOAuth2ClaimauthServerId string

	GetOAuth2ClaimclaimId string

	GetOAuth2ClaimBackupDir string

	GetOAuth2ClaimQuiet bool
)

func NewGetOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOAuth2Claim",
		Long: "Retrieve a custom token Claim",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.GetOAuth2Claim(apiClient.GetConfig().Context, GetOAuth2ClaimauthServerId, GetOAuth2ClaimclaimId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOAuth2ClaimQuiet {
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
				dirPath := filepath.Join(GetOAuth2ClaimBackupDir, "authorizationserverclaims", "getOAuth2Claim")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetOAuth2ClaimauthServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOAuth2ClaimQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOAuth2ClaimQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetOAuth2ClaimclaimId, "claimId", "", "", "")
	cmd.MarkFlagRequired("claimId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AuthorizationServerClaims to a file")

	cmd.Flags().StringVarP(&GetOAuth2ClaimBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOAuth2ClaimQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOAuth2ClaimCmd := NewGetOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(GetOAuth2ClaimCmd)
}

var (
	ReplaceOAuth2ClaimauthServerId string

	ReplaceOAuth2ClaimclaimId string

	ReplaceOAuth2Claimdata string

	ReplaceOAuth2ClaimQuiet bool
)

func NewReplaceOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceOAuth2Claim",
		Long: "Replace a custom token Claim",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.ReplaceOAuth2Claim(apiClient.GetConfig().Context, ReplaceOAuth2ClaimauthServerId, ReplaceOAuth2ClaimclaimId)

			if ReplaceOAuth2Claimdata != "" {
				req = req.Data(ReplaceOAuth2Claimdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceOAuth2ClaimQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceOAuth2ClaimQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ReplaceOAuth2ClaimclaimId, "claimId", "", "", "")
	cmd.MarkFlagRequired("claimId")

	cmd.Flags().StringVarP(&ReplaceOAuth2Claimdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceOAuth2ClaimQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceOAuth2ClaimCmd := NewReplaceOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(ReplaceOAuth2ClaimCmd)
}

var (
	DeleteOAuth2ClaimauthServerId string

	DeleteOAuth2ClaimclaimId string

	DeleteOAuth2ClaimQuiet bool
)

func NewDeleteOAuth2ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteOAuth2Claim",
		Long: "Delete a custom token Claim",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClaimsAPI.DeleteOAuth2Claim(apiClient.GetConfig().Context, DeleteOAuth2ClaimauthServerId, DeleteOAuth2ClaimclaimId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteOAuth2ClaimQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteOAuth2ClaimQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteOAuth2ClaimauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteOAuth2ClaimclaimId, "claimId", "", "", "")
	cmd.MarkFlagRequired("claimId")

	cmd.Flags().BoolVarP(&DeleteOAuth2ClaimQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteOAuth2ClaimCmd := NewDeleteOAuth2ClaimCmd()
	AuthorizationServerClaimsCmd.AddCommand(DeleteOAuth2ClaimCmd)
}
