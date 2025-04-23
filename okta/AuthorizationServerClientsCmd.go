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

var AuthorizationServerClientsCmd = &cobra.Command{
	Use:  "authorizationServerClients",
	Long: "Manage AuthorizationServerClientsAPI",
}

func init() {
	rootCmd.AddCommand(AuthorizationServerClientsCmd)
}

var (
	ListOAuth2ClientsForAuthorizationServerauthServerId string

	ListOAuth2ClientsForAuthorizationServerBackupDir string

	ListOAuth2ClientsForAuthorizationServerLimit    int32
	ListOAuth2ClientsForAuthorizationServerPage     string
	ListOAuth2ClientsForAuthorizationServerFetchAll bool

	ListOAuth2ClientsForAuthorizationServerQuiet bool
)

func NewListOAuth2ClientsForAuthorizationServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listOAuth2ClientsForAuthorizationServer",
		Long: "List all Client resources for an authorization server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClientsAPI.ListOAuth2ClientsForAuthorizationServer(apiClient.GetConfig().Context, ListOAuth2ClientsForAuthorizationServerauthServerId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListOAuth2ClientsForAuthorizationServerQuiet {
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
					if !ListOAuth2ClientsForAuthorizationServerQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListOAuth2ClientsForAuthorizationServerFetchAll || len(items) == 0 {
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

			if ListOAuth2ClientsForAuthorizationServerFetchAll && pageCount > 1 && !ListOAuth2ClientsForAuthorizationServerQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListOAuth2ClientsForAuthorizationServerBackupDir, "authorizationserverclients", "listOAuth2ClientsForAuthorizationServer")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListOAuth2ClientsForAuthorizationServerQuiet {
					fmt.Printf("Backing up AuthorizationServerClientss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListOAuth2ClientsForAuthorizationServerQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListOAuth2ClientsForAuthorizationServerQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListOAuth2ClientsForAuthorizationServerQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListOAuth2ClientsForAuthorizationServerQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerClientss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListOAuth2ClientsForAuthorizationServerQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListOAuth2ClientsForAuthorizationServerauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().Int32VarP(&ListOAuth2ClientsForAuthorizationServerLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListOAuth2ClientsForAuthorizationServerPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListOAuth2ClientsForAuthorizationServerFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerClientss to a directory")

	cmd.Flags().StringVarP(&ListOAuth2ClientsForAuthorizationServerBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListOAuth2ClientsForAuthorizationServerQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListOAuth2ClientsForAuthorizationServerCmd := NewListOAuth2ClientsForAuthorizationServerCmd()
	AuthorizationServerClientsCmd.AddCommand(ListOAuth2ClientsForAuthorizationServerCmd)
}

var (
	ListRefreshTokensForAuthorizationServerAndClientauthServerId string

	ListRefreshTokensForAuthorizationServerAndClientclientId string

	ListRefreshTokensForAuthorizationServerAndClientBackupDir string

	ListRefreshTokensForAuthorizationServerAndClientLimit    int32
	ListRefreshTokensForAuthorizationServerAndClientPage     string
	ListRefreshTokensForAuthorizationServerAndClientFetchAll bool

	ListRefreshTokensForAuthorizationServerAndClientQuiet bool
)

func NewListRefreshTokensForAuthorizationServerAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listRefreshTokensForAuthorizationServerAndClient",
		Long: "List all refresh tokens for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClientsAPI.ListRefreshTokensForAuthorizationServerAndClient(apiClient.GetConfig().Context, ListRefreshTokensForAuthorizationServerAndClientauthServerId, ListRefreshTokensForAuthorizationServerAndClientclientId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListRefreshTokensForAuthorizationServerAndClientQuiet {
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
					if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListRefreshTokensForAuthorizationServerAndClientFetchAll || len(items) == 0 {
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

			if ListRefreshTokensForAuthorizationServerAndClientFetchAll && pageCount > 1 && !ListRefreshTokensForAuthorizationServerAndClientQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListRefreshTokensForAuthorizationServerAndClientBackupDir, "authorizationserverclients", "listRefreshTokensForAuthorizationServerAndClient")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
					fmt.Printf("Backing up AuthorizationServerClientss to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
					fmt.Printf("Successfully backed up %d/%d AuthorizationServerClientss\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListRefreshTokensForAuthorizationServerAndClientQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListRefreshTokensForAuthorizationServerAndClientauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&ListRefreshTokensForAuthorizationServerAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().Int32VarP(&ListRefreshTokensForAuthorizationServerAndClientLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListRefreshTokensForAuthorizationServerAndClientPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListRefreshTokensForAuthorizationServerAndClientFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple AuthorizationServerClientss to a directory")

	cmd.Flags().StringVarP(&ListRefreshTokensForAuthorizationServerAndClientBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListRefreshTokensForAuthorizationServerAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListRefreshTokensForAuthorizationServerAndClientCmd := NewListRefreshTokensForAuthorizationServerAndClientCmd()
	AuthorizationServerClientsCmd.AddCommand(ListRefreshTokensForAuthorizationServerAndClientCmd)
}

var (
	RevokeRefreshTokensForAuthorizationServerAndClientauthServerId string

	RevokeRefreshTokensForAuthorizationServerAndClientclientId string

	RevokeRefreshTokensForAuthorizationServerAndClientQuiet bool
)

func NewRevokeRefreshTokensForAuthorizationServerAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeRefreshTokensForAuthorizationServerAndClient",
		Long: "Revoke all refresh tokens for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClientsAPI.RevokeRefreshTokensForAuthorizationServerAndClient(apiClient.GetConfig().Context, RevokeRefreshTokensForAuthorizationServerAndClientauthServerId, RevokeRefreshTokensForAuthorizationServerAndClientclientId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeRefreshTokensForAuthorizationServerAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeRefreshTokensForAuthorizationServerAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeRefreshTokensForAuthorizationServerAndClientauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&RevokeRefreshTokensForAuthorizationServerAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().BoolVarP(&RevokeRefreshTokensForAuthorizationServerAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeRefreshTokensForAuthorizationServerAndClientCmd := NewRevokeRefreshTokensForAuthorizationServerAndClientCmd()
	AuthorizationServerClientsCmd.AddCommand(RevokeRefreshTokensForAuthorizationServerAndClientCmd)
}

var (
	GetRefreshTokenForAuthorizationServerAndClientauthServerId string

	GetRefreshTokenForAuthorizationServerAndClientclientId string

	GetRefreshTokenForAuthorizationServerAndClienttokenId string

	GetRefreshTokenForAuthorizationServerAndClientBackupDir string

	GetRefreshTokenForAuthorizationServerAndClientQuiet bool
)

func NewGetRefreshTokenForAuthorizationServerAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getRefreshTokenForAuthorizationServerAndClient",
		Long: "Retrieve a refresh token for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClientsAPI.GetRefreshTokenForAuthorizationServerAndClient(apiClient.GetConfig().Context, GetRefreshTokenForAuthorizationServerAndClientauthServerId, GetRefreshTokenForAuthorizationServerAndClientclientId, GetRefreshTokenForAuthorizationServerAndClienttokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetRefreshTokenForAuthorizationServerAndClientQuiet {
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
				dirPath := filepath.Join(GetRefreshTokenForAuthorizationServerAndClientBackupDir, "authorizationserverclients", "getRefreshTokenForAuthorizationServerAndClient")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetRefreshTokenForAuthorizationServerAndClientauthServerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetRefreshTokenForAuthorizationServerAndClientQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetRefreshTokenForAuthorizationServerAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRefreshTokenForAuthorizationServerAndClientauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&GetRefreshTokenForAuthorizationServerAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().StringVarP(&GetRefreshTokenForAuthorizationServerAndClienttokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the AuthorizationServerClients to a file")

	cmd.Flags().StringVarP(&GetRefreshTokenForAuthorizationServerAndClientBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetRefreshTokenForAuthorizationServerAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetRefreshTokenForAuthorizationServerAndClientCmd := NewGetRefreshTokenForAuthorizationServerAndClientCmd()
	AuthorizationServerClientsCmd.AddCommand(GetRefreshTokenForAuthorizationServerAndClientCmd)
}

var (
	RevokeRefreshTokenForAuthorizationServerAndClientauthServerId string

	RevokeRefreshTokenForAuthorizationServerAndClientclientId string

	RevokeRefreshTokenForAuthorizationServerAndClienttokenId string

	RevokeRefreshTokenForAuthorizationServerAndClientQuiet bool
)

func NewRevokeRefreshTokenForAuthorizationServerAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeRefreshTokenForAuthorizationServerAndClient",
		Long: "Revoke a refresh token for a Client",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.AuthorizationServerClientsAPI.RevokeRefreshTokenForAuthorizationServerAndClient(apiClient.GetConfig().Context, RevokeRefreshTokenForAuthorizationServerAndClientauthServerId, RevokeRefreshTokenForAuthorizationServerAndClientclientId, RevokeRefreshTokenForAuthorizationServerAndClienttokenId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeRefreshTokenForAuthorizationServerAndClientQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeRefreshTokenForAuthorizationServerAndClientQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeRefreshTokenForAuthorizationServerAndClientauthServerId, "authServerId", "", "", "")
	cmd.MarkFlagRequired("authServerId")

	cmd.Flags().StringVarP(&RevokeRefreshTokenForAuthorizationServerAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().StringVarP(&RevokeRefreshTokenForAuthorizationServerAndClienttokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	cmd.Flags().BoolVarP(&RevokeRefreshTokenForAuthorizationServerAndClientQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeRefreshTokenForAuthorizationServerAndClientCmd := NewRevokeRefreshTokenForAuthorizationServerAndClientCmd()
	AuthorizationServerClientsCmd.AddCommand(RevokeRefreshTokenForAuthorizationServerAndClientCmd)
}
