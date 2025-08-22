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

var AuthorizationServerAssocCmd = &cobra.Command{
	Use:  "authorizationServerAssoc",
	Long: "Manage AuthorizationServerAssocAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerAssocCmd)
}

var (
	CreateAssociatedServersauthServerId string

	CreateAssociatedServersdata string

	CreateAssociatedServersRestoreFile string

	CreateAssociatedServersQuiet bool
)

func NewCreateAssociatedServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createAssociatedServers",
		Long: "Create an associated Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateAssociatedServersRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateAssociatedServersRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateAssociatedServersdata = string(processedData)

				if !CreateAssociatedServersQuiet {
					fmt.Println("Restoring AuthorizationServerAssoc from:", CreateAssociatedServersRestoreFile)
				}
			}

			req := apiClient.AuthorizationServerAssocAPI.CreateAssociatedServers(apiClient.GetConfig().Context, CreateAssociatedServersauthServerId)

			if CreateAssociatedServersdata != "" {
				req = req.Data(CreateAssociatedServersdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateAssociatedServersQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateAssociatedServersQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateAssociatedServersauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&CreateAssociatedServersdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateAssociatedServersRestoreFile, "restore-from", "r", "", "Restore AuthorizationServerAssoc from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateAssociatedServersQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateAssociatedServersCmd := NewCreateAssociatedServersCmd()
	AuthorizationServerAssocCmd.AddCommand(CreateAssociatedServersCmd)
}

var (
	ListAssociatedServersByTrustedTypeauthServerId string

	ListAssociatedServersByTrustedTypeBackupDir string

	ListAssociatedServersByTrustedTypeLimit    int32
	ListAssociatedServersByTrustedTypePage     string
	ListAssociatedServersByTrustedTypeFetchAll bool

	ListAssociatedServersByTrustedTypeQuiet bool
)

func NewListAssociatedServersByTrustedTypeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAssociatedServersByTrustedType",
		Long: "List all associated Authorization Servers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAssocAPI.ListAssociatedServersByTrustedType(apiClient.GetConfig().Context, ListAssociatedServersByTrustedTypeauthServerId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAssociatedServersByTrustedTypeQuiet {
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
					if !ListAssociatedServersByTrustedTypeQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAssociatedServersByTrustedTypeFetchAll || len(items) == 0 {
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

			if ListAssociatedServersByTrustedTypeFetchAll && pageCount > 1 && !ListAssociatedServersByTrustedTypeQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAssociatedServersByTrustedTypeBackupDir, "authorizationserverassoc", "listAssociatedServersByTrustedType")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAssociatedServersByTrustedTypeQuiet {
					fmt.Printf("Backing up AuthorizationServerAssocs to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAssociatedServersByTrustedTypeQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAssociatedServersByTrustedTypeQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAssociatedServersByTrustedTypeQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAssociatedServersByTrustedTypeQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerAssocs\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAssociatedServersByTrustedTypeQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAssociatedServersByTrustedTypeauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().Int32VarP(&ListAssociatedServersByTrustedTypeLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAssociatedServersByTrustedTypePage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAssociatedServersByTrustedTypeFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerAssocs to a directory")

	cmd.Flags().StringVarP(&ListAssociatedServersByTrustedTypeBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAssociatedServersByTrustedTypeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAssociatedServersByTrustedTypeCmd := NewListAssociatedServersByTrustedTypeCmd()
	AuthorizationServerAssocCmd.AddCommand(ListAssociatedServersByTrustedTypeCmd)
}

var (
	DeleteAssociatedServerauthServerId string

	DeleteAssociatedServerassociatedServerId string

	DeleteAssociatedServerQuiet bool
)

func NewDeleteAssociatedServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteAssociatedServer",
		Long: "Delete an associated Authorization Server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerAssocAPI.DeleteAssociatedServer(apiClient.GetConfig().Context, DeleteAssociatedServerauthServerId, DeleteAssociatedServerassociatedServerId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteAssociatedServerQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteAssociatedServerQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteAssociatedServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&DeleteAssociatedServerassociatedServerId, "associatedServerId", "", "", "")
	cmd.MarkFlagRequired("associatedServerId")

	cmd.Flags().BoolVarP(&DeleteAssociatedServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteAssociatedServerCmd := NewDeleteAssociatedServerCmd()
	AuthorizationServerAssocCmd.AddCommand(DeleteAssociatedServerCmd)
}
