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

var IdentitySourceCmd = &cobra.Command{
	Use:  "identitySource",
	Long: "Manage IdentitySourceAPI",
}

func init() {
	rootCmd.AddCommand(IdentitySourceCmd)
}

var (
	CreateIdentitySourceSessionidentitySourceId string

	CreateIdentitySourceSessionRestoreFile string

	CreateIdentitySourceSessionQuiet bool
)

func NewCreateIdentitySourceSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createSession",
		Long: "Create an Identity Source Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateIdentitySourceSessionRestoreFile != "" {

				fmt.Println("Warning: This Create operation doesn't support data input. Cannot restore from file.")
				return fmt.Errorf("restore from file not supported for this operation that doesn't accept data input")

				if !CreateIdentitySourceSessionQuiet {
					fmt.Println("Restoring IdentitySource from:", CreateIdentitySourceSessionRestoreFile)
				}
			}

			req := apiClient.IdentitySourceAPI.CreateIdentitySourceSession(apiClient.GetConfig().Context, CreateIdentitySourceSessionidentitySourceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateIdentitySourceSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateIdentitySourceSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateIdentitySourceSessionidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&CreateIdentitySourceSessionRestoreFile, "restore-from", "r", "", "Restore IdentitySource from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateIdentitySourceSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateIdentitySourceSessionCmd := NewCreateIdentitySourceSessionCmd()
	IdentitySourceCmd.AddCommand(CreateIdentitySourceSessionCmd)
}

var (
	ListIdentitySourceSessionsidentitySourceId string

	ListIdentitySourceSessionsBackupDir string

	ListIdentitySourceSessionsLimit    int32
	ListIdentitySourceSessionsPage     string
	ListIdentitySourceSessionsFetchAll bool

	ListIdentitySourceSessionsQuiet bool
)

func NewListIdentitySourceSessionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSessions",
		Long: "List all Identity Source Sessions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.ListIdentitySourceSessions(apiClient.GetConfig().Context, ListIdentitySourceSessionsidentitySourceId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListIdentitySourceSessionsQuiet {
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
					if !ListIdentitySourceSessionsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListIdentitySourceSessionsFetchAll || len(items) == 0 {
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

			if ListIdentitySourceSessionsFetchAll && pageCount > 1 && !ListIdentitySourceSessionsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListIdentitySourceSessionsBackupDir, "identitysource", "listSessions")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListIdentitySourceSessionsQuiet {
					fmt.Printf("Backing up IdentitySources to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListIdentitySourceSessionsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListIdentitySourceSessionsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListIdentitySourceSessionsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListIdentitySourceSessionsQuiet {
					fmt.Printf("Successfully backed up %d/%d IdentitySources\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListIdentitySourceSessionsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListIdentitySourceSessionsidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().Int32VarP(&ListIdentitySourceSessionsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListIdentitySourceSessionsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListIdentitySourceSessionsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple IdentitySources to a directory")

	cmd.Flags().StringVarP(&ListIdentitySourceSessionsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListIdentitySourceSessionsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListIdentitySourceSessionsCmd := NewListIdentitySourceSessionsCmd()
	IdentitySourceCmd.AddCommand(ListIdentitySourceSessionsCmd)
}

var (
	GetIdentitySourceSessionidentitySourceId string

	GetIdentitySourceSessionsessionId string

	GetIdentitySourceSessionBackupDir string

	GetIdentitySourceSessionQuiet bool
)

func NewGetIdentitySourceSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSession",
		Long: "Retrieve an Identity Source Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.GetIdentitySourceSession(apiClient.GetConfig().Context, GetIdentitySourceSessionidentitySourceId, GetIdentitySourceSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetIdentitySourceSessionQuiet {
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
				dirPath := filepath.Join(GetIdentitySourceSessionBackupDir, "identitysource", "getSession")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetIdentitySourceSessionidentitySourceId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetIdentitySourceSessionQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetIdentitySourceSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetIdentitySourceSessionidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&GetIdentitySourceSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the IdentitySource to a file")

	cmd.Flags().StringVarP(&GetIdentitySourceSessionBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetIdentitySourceSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetIdentitySourceSessionCmd := NewGetIdentitySourceSessionCmd()
	IdentitySourceCmd.AddCommand(GetIdentitySourceSessionCmd)
}

var (
	DeleteIdentitySourceSessionidentitySourceId string

	DeleteIdentitySourceSessionsessionId string

	DeleteIdentitySourceSessionQuiet bool
)

func NewDeleteIdentitySourceSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteSession",
		Long: "Delete an Identity Source Session",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.DeleteIdentitySourceSession(apiClient.GetConfig().Context, DeleteIdentitySourceSessionidentitySourceId, DeleteIdentitySourceSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteIdentitySourceSessionQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteIdentitySourceSessionQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteIdentitySourceSessionidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&DeleteIdentitySourceSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().BoolVarP(&DeleteIdentitySourceSessionQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteIdentitySourceSessionCmd := NewDeleteIdentitySourceSessionCmd()
	IdentitySourceCmd.AddCommand(DeleteIdentitySourceSessionCmd)
}

var (
	UploadIdentitySourceDataForDeleteidentitySourceId string

	UploadIdentitySourceDataForDeletesessionId string

	UploadIdentitySourceDataForDeletedata string

	UploadIdentitySourceDataForDeleteQuiet bool
)

func NewUploadIdentitySourceDataForDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadDataForDelete",
		Long: "Upload the data to be deleted in Okta",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.UploadIdentitySourceDataForDelete(apiClient.GetConfig().Context, UploadIdentitySourceDataForDeleteidentitySourceId, UploadIdentitySourceDataForDeletesessionId)

			if UploadIdentitySourceDataForDeletedata != "" {
				req = req.Data(UploadIdentitySourceDataForDeletedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadIdentitySourceDataForDeleteQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadIdentitySourceDataForDeleteQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForDeleteidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForDeletesessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForDeletedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadIdentitySourceDataForDeleteQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadIdentitySourceDataForDeleteCmd := NewUploadIdentitySourceDataForDeleteCmd()
	IdentitySourceCmd.AddCommand(UploadIdentitySourceDataForDeleteCmd)
}

var (
	UploadIdentitySourceDataForUpsertidentitySourceId string

	UploadIdentitySourceDataForUpsertsessionId string

	UploadIdentitySourceDataForUpsertdata string

	UploadIdentitySourceDataForUpsertQuiet bool
)

func NewUploadIdentitySourceDataForUpsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadDataForUpsert",
		Long: "Upload the data to be upserted in Okta",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.UploadIdentitySourceDataForUpsert(apiClient.GetConfig().Context, UploadIdentitySourceDataForUpsertidentitySourceId, UploadIdentitySourceDataForUpsertsessionId)

			if UploadIdentitySourceDataForUpsertdata != "" {
				req = req.Data(UploadIdentitySourceDataForUpsertdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadIdentitySourceDataForUpsertQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadIdentitySourceDataForUpsertQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForUpsertidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForUpsertsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForUpsertdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadIdentitySourceDataForUpsertQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadIdentitySourceDataForUpsertCmd := NewUploadIdentitySourceDataForUpsertCmd()
	IdentitySourceCmd.AddCommand(UploadIdentitySourceDataForUpsertCmd)
}

var (
	StartImportFromIdentitySourceidentitySourceId string

	StartImportFromIdentitySourcesessionId string

	StartImportFromIdentitySourceQuiet bool
)

func NewStartImportFromIdentitySourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "startImportFrom",
		Long: "Start the import from the Identity Source",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.StartImportFromIdentitySource(apiClient.GetConfig().Context, StartImportFromIdentitySourceidentitySourceId, StartImportFromIdentitySourcesessionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !StartImportFromIdentitySourceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !StartImportFromIdentitySourceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&StartImportFromIdentitySourceidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&StartImportFromIdentitySourcesessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().BoolVarP(&StartImportFromIdentitySourceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	StartImportFromIdentitySourceCmd := NewStartImportFromIdentitySourceCmd()
	IdentitySourceCmd.AddCommand(StartImportFromIdentitySourceCmd)
}
