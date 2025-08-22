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

var ResourceSetCmd = &cobra.Command{
	Use:  "resourceSet",
	Long: "Manage ResourceSetAPI",
}

func init() {
	rootCmd.AddCommand(ResourceSetCmd)
}

var (
	CreateResourceSetdata string

	CreateResourceSetRestoreFile string

	CreateResourceSetQuiet bool
)

func NewCreateResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateResourceSetRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateResourceSetRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateResourceSetdata = string(processedData)

				if !CreateResourceSetQuiet {
					fmt.Println("Restoring ResourceSet from:", CreateResourceSetRestoreFile)
				}
			}

			req := apiClient.ResourceSetAPI.CreateResourceSet(apiClient.GetConfig().Context)

			if CreateResourceSetdata != "" {
				req = req.Data(CreateResourceSetdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateResourceSetQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateResourceSetQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateResourceSetdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateResourceSetRestoreFile, "restore-from", "r", "", "Restore ResourceSet from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateResourceSetQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateResourceSetCmd := NewCreateResourceSetCmd()
	ResourceSetCmd.AddCommand(CreateResourceSetCmd)
}

var (
	ListResourceSetsBackupDir string

	ListResourceSetsLimit    int32
	ListResourceSetsPage     string
	ListResourceSetsFetchAll bool

	ListResourceSetsQuiet bool
)

func NewListResourceSetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Resource Sets",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListResourceSets(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListResourceSetsQuiet {
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
					if !ListResourceSetsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListResourceSetsFetchAll || len(items) == 0 {
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

			if ListResourceSetsFetchAll && pageCount > 1 && !ListResourceSetsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListResourceSetsBackupDir, "resourceset", "lists")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListResourceSetsQuiet {
					fmt.Printf("Backing up ResourceSets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListResourceSetsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListResourceSetsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListResourceSetsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListResourceSetsQuiet {
					fmt.Printf("Successfully backed up %d/%d ResourceSets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListResourceSetsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListResourceSetsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListResourceSetsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListResourceSetsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ResourceSets to a directory")

	cmd.Flags().StringVarP(&ListResourceSetsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListResourceSetsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListResourceSetsCmd := NewListResourceSetsCmd()
	ResourceSetCmd.AddCommand(ListResourceSetsCmd)
}

var (
	GetResourceSetresourceSetId string

	GetResourceSetBackupDir string

	GetResourceSetQuiet bool
)

func NewGetResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.GetResourceSet(apiClient.GetConfig().Context, GetResourceSetresourceSetId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetResourceSetQuiet {
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
				dirPath := filepath.Join(GetResourceSetBackupDir, "resourceset", "get")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetResourceSetresourceSetId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetResourceSetQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetResourceSetQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetResourceSetresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ResourceSet to a file")

	cmd.Flags().StringVarP(&GetResourceSetBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetResourceSetQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetResourceSetCmd := NewGetResourceSetCmd()
	ResourceSetCmd.AddCommand(GetResourceSetCmd)
}

var (
	ReplaceResourceSetresourceSetId string

	ReplaceResourceSetdata string

	ReplaceResourceSetQuiet bool
)

func NewReplaceResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ReplaceResourceSet(apiClient.GetConfig().Context, ReplaceResourceSetresourceSetId)

			if ReplaceResourceSetdata != "" {
				req = req.Data(ReplaceResourceSetdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceResourceSetQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceResourceSetQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceResourceSetresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&ReplaceResourceSetdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceResourceSetQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceResourceSetCmd := NewReplaceResourceSetCmd()
	ResourceSetCmd.AddCommand(ReplaceResourceSetCmd)
}

var (
	DeleteResourceSetresourceSetId string

	DeleteResourceSetQuiet bool
)

func NewDeleteResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.DeleteResourceSet(apiClient.GetConfig().Context, DeleteResourceSetresourceSetId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteResourceSetQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteResourceSetQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteResourceSetresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().BoolVarP(&DeleteResourceSetQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteResourceSetCmd := NewDeleteResourceSetCmd()
	ResourceSetCmd.AddCommand(DeleteResourceSetCmd)
}

var (
	CreateResourceSetBindingresourceSetId string

	CreateResourceSetBindingdata string

	CreateResourceSetBindingRestoreFile string

	CreateResourceSetBindingQuiet bool
)

func NewCreateResourceSetBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createBinding",
		Long: "Create a Resource Set Binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateResourceSetBindingRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateResourceSetBindingRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateResourceSetBindingdata = string(processedData)

				if !CreateResourceSetBindingQuiet {
					fmt.Println("Restoring ResourceSet from:", CreateResourceSetBindingRestoreFile)
				}
			}

			req := apiClient.ResourceSetAPI.CreateResourceSetBinding(apiClient.GetConfig().Context, CreateResourceSetBindingresourceSetId)

			if CreateResourceSetBindingdata != "" {
				req = req.Data(CreateResourceSetBindingdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateResourceSetBindingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateResourceSetBindingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateResourceSetBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&CreateResourceSetBindingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateResourceSetBindingRestoreFile, "restore-from", "r", "", "Restore ResourceSet from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateResourceSetBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateResourceSetBindingCmd := NewCreateResourceSetBindingCmd()
	ResourceSetCmd.AddCommand(CreateResourceSetBindingCmd)
}

var (
	ListBindingsresourceSetId string

	ListBindingsBackupDir string

	ListBindingsLimit    int32
	ListBindingsPage     string
	ListBindingsFetchAll bool

	ListBindingsQuiet bool
)

func NewListBindingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBindings",
		Long: "List all Bindings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListBindings(apiClient.GetConfig().Context, ListBindingsresourceSetId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListBindingsQuiet {
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
					if !ListBindingsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListBindingsFetchAll || len(items) == 0 {
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

			if ListBindingsFetchAll && pageCount > 1 && !ListBindingsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListBindingsBackupDir, "resourceset", "listBindings")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListBindingsQuiet {
					fmt.Printf("Backing up ResourceSets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListBindingsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListBindingsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListBindingsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListBindingsQuiet {
					fmt.Printf("Successfully backed up %d/%d ResourceSets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListBindingsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListBindingsresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().Int32VarP(&ListBindingsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListBindingsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListBindingsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ResourceSets to a directory")

	cmd.Flags().StringVarP(&ListBindingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListBindingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListBindingsCmd := NewListBindingsCmd()
	ResourceSetCmd.AddCommand(ListBindingsCmd)
}

var (
	GetBindingresourceSetId string

	GetBindingroleIdOrLabel string

	GetBindingBackupDir string

	GetBindingQuiet bool
)

func NewGetBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getBinding",
		Long: "Retrieve a Binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.GetBinding(apiClient.GetConfig().Context, GetBindingresourceSetId, GetBindingroleIdOrLabel)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetBindingQuiet {
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
				dirPath := filepath.Join(GetBindingBackupDir, "resourceset", "getBinding")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetBindingresourceSetId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetBindingQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetBindingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&GetBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ResourceSet to a file")

	cmd.Flags().StringVarP(&GetBindingBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetBindingCmd := NewGetBindingCmd()
	ResourceSetCmd.AddCommand(GetBindingCmd)
}

var (
	DeleteBindingresourceSetId string

	DeleteBindingroleIdOrLabel string

	DeleteBindingQuiet bool
)

func NewDeleteBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBinding",
		Long: "Delete a Binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.DeleteBinding(apiClient.GetConfig().Context, DeleteBindingresourceSetId, DeleteBindingroleIdOrLabel)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteBindingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteBindingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&DeleteBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().BoolVarP(&DeleteBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteBindingCmd := NewDeleteBindingCmd()
	ResourceSetCmd.AddCommand(DeleteBindingCmd)
}

var (
	ListMembersOfBindingresourceSetId string

	ListMembersOfBindingroleIdOrLabel string

	ListMembersOfBindingBackupDir string

	ListMembersOfBindingLimit    int32
	ListMembersOfBindingPage     string
	ListMembersOfBindingFetchAll bool

	ListMembersOfBindingQuiet bool
)

func NewListMembersOfBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listMembersOfBinding",
		Long: "List all Members of a binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListMembersOfBinding(apiClient.GetConfig().Context, ListMembersOfBindingresourceSetId, ListMembersOfBindingroleIdOrLabel)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListMembersOfBindingQuiet {
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
					if !ListMembersOfBindingQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListMembersOfBindingFetchAll || len(items) == 0 {
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

			if ListMembersOfBindingFetchAll && pageCount > 1 && !ListMembersOfBindingQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListMembersOfBindingBackupDir, "resourceset", "listMembersOfBinding")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListMembersOfBindingQuiet {
					fmt.Printf("Backing up ResourceSets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListMembersOfBindingQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListMembersOfBindingQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListMembersOfBindingQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListMembersOfBindingQuiet {
					fmt.Printf("Successfully backed up %d/%d ResourceSets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListMembersOfBindingQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListMembersOfBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&ListMembersOfBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().Int32VarP(&ListMembersOfBindingLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListMembersOfBindingPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListMembersOfBindingFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ResourceSets to a directory")

	cmd.Flags().StringVarP(&ListMembersOfBindingBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListMembersOfBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListMembersOfBindingCmd := NewListMembersOfBindingCmd()
	ResourceSetCmd.AddCommand(ListMembersOfBindingCmd)
}

var (
	AddMembersToBindingresourceSetId string

	AddMembersToBindingroleIdOrLabel string

	AddMembersToBindingdata string

	AddMembersToBindingQuiet bool
)

func NewAddMembersToBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "addMembersToBinding",
		Long: "Add more Members to a binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.AddMembersToBinding(apiClient.GetConfig().Context, AddMembersToBindingresourceSetId, AddMembersToBindingroleIdOrLabel)

			if AddMembersToBindingdata != "" {
				req = req.Data(AddMembersToBindingdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AddMembersToBindingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AddMembersToBindingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AddMembersToBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&AddMembersToBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&AddMembersToBindingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AddMembersToBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AddMembersToBindingCmd := NewAddMembersToBindingCmd()
	ResourceSetCmd.AddCommand(AddMembersToBindingCmd)
}

var (
	GetMemberOfBindingresourceSetId string

	GetMemberOfBindingroleIdOrLabel string

	GetMemberOfBindingmemberId string

	GetMemberOfBindingBackupDir string

	GetMemberOfBindingQuiet bool
)

func NewGetMemberOfBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getMemberOfBinding",
		Long: "Retrieve a Member of a binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.GetMemberOfBinding(apiClient.GetConfig().Context, GetMemberOfBindingresourceSetId, GetMemberOfBindingroleIdOrLabel, GetMemberOfBindingmemberId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetMemberOfBindingQuiet {
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
				dirPath := filepath.Join(GetMemberOfBindingBackupDir, "resourceset", "getMemberOfBinding")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetMemberOfBindingresourceSetId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetMemberOfBindingQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetMemberOfBindingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetMemberOfBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&GetMemberOfBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&GetMemberOfBindingmemberId, "memberId", "", "", "")
	cmd.MarkFlagRequired("memberId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ResourceSet to a file")

	cmd.Flags().StringVarP(&GetMemberOfBindingBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetMemberOfBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetMemberOfBindingCmd := NewGetMemberOfBindingCmd()
	ResourceSetCmd.AddCommand(GetMemberOfBindingCmd)
}

var (
	UnassignMemberFromBindingresourceSetId string

	UnassignMemberFromBindingroleIdOrLabel string

	UnassignMemberFromBindingmemberId string

	UnassignMemberFromBindingQuiet bool
)

func NewUnassignMemberFromBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignMemberFromBinding",
		Long: "Unassign a Member from a binding",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.UnassignMemberFromBinding(apiClient.GetConfig().Context, UnassignMemberFromBindingresourceSetId, UnassignMemberFromBindingroleIdOrLabel, UnassignMemberFromBindingmemberId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnassignMemberFromBindingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnassignMemberFromBindingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignMemberFromBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&UnassignMemberFromBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&UnassignMemberFromBindingmemberId, "memberId", "", "", "")
	cmd.MarkFlagRequired("memberId")

	cmd.Flags().BoolVarP(&UnassignMemberFromBindingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnassignMemberFromBindingCmd := NewUnassignMemberFromBindingCmd()
	ResourceSetCmd.AddCommand(UnassignMemberFromBindingCmd)
}

var (
	ListResourceSetResourcesresourceSetId string

	ListResourceSetResourcesBackupDir string

	ListResourceSetResourcesLimit    int32
	ListResourceSetResourcesPage     string
	ListResourceSetResourcesFetchAll bool

	ListResourceSetResourcesQuiet bool
)

func NewListResourceSetResourcesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listResources",
		Long: "List all Resources of a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListResourceSetResources(apiClient.GetConfig().Context, ListResourceSetResourcesresourceSetId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListResourceSetResourcesQuiet {
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
					if !ListResourceSetResourcesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListResourceSetResourcesFetchAll || len(items) == 0 {
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

			if ListResourceSetResourcesFetchAll && pageCount > 1 && !ListResourceSetResourcesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListResourceSetResourcesBackupDir, "resourceset", "listResources")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListResourceSetResourcesQuiet {
					fmt.Printf("Backing up ResourceSets to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListResourceSetResourcesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListResourceSetResourcesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListResourceSetResourcesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListResourceSetResourcesQuiet {
					fmt.Printf("Successfully backed up %d/%d ResourceSets\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListResourceSetResourcesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListResourceSetResourcesresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().Int32VarP(&ListResourceSetResourcesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListResourceSetResourcesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListResourceSetResourcesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ResourceSets to a directory")

	cmd.Flags().StringVarP(&ListResourceSetResourcesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListResourceSetResourcesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListResourceSetResourcesCmd := NewListResourceSetResourcesCmd()
	ResourceSetCmd.AddCommand(ListResourceSetResourcesCmd)
}

var (
	AddResourceSetResourceresourceSetId string

	AddResourceSetResourcedata string

	AddResourceSetResourceQuiet bool
)

func NewAddResourceSetResourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "addResource",
		Long: "Add more Resource to a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.AddResourceSetResource(apiClient.GetConfig().Context, AddResourceSetResourceresourceSetId)

			if AddResourceSetResourcedata != "" {
				req = req.Data(AddResourceSetResourcedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AddResourceSetResourceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AddResourceSetResourceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AddResourceSetResourceresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&AddResourceSetResourcedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AddResourceSetResourceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AddResourceSetResourceCmd := NewAddResourceSetResourceCmd()
	ResourceSetCmd.AddCommand(AddResourceSetResourceCmd)
}

var (
	DeleteResourceSetResourceresourceSetId string

	DeleteResourceSetResourceresourceId string

	DeleteResourceSetResourceQuiet bool
)

func NewDeleteResourceSetResourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteResource",
		Long: "Delete a Resource from a Resource Set",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.DeleteResourceSetResource(apiClient.GetConfig().Context, DeleteResourceSetResourceresourceSetId, DeleteResourceSetResourceresourceId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteResourceSetResourceQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteResourceSetResourceQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteResourceSetResourceresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&DeleteResourceSetResourceresourceId, "resourceId", "", "", "")
	cmd.MarkFlagRequired("resourceId")

	cmd.Flags().BoolVarP(&DeleteResourceSetResourceQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteResourceSetResourceCmd := NewDeleteResourceSetResourceCmd()
	ResourceSetCmd.AddCommand(DeleteResourceSetResourceCmd)
}
