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

var CustomizationCmd = &cobra.Command{
	Use:  "customization",
	Long: "Manage CustomizationAPI",
}

func init() {
	rootCmd.AddCommand(CustomizationCmd)
}

var (
	CreateBranddata string

	CreateBrandRestoreFile string

	CreateBrandQuiet bool
)

func NewCreateBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createBrand",
		Long: "Create a Brand",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateBrandRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateBrandRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateBranddata = string(processedData)

				if !CreateBrandQuiet {
					fmt.Println("Restoring Customization from:", CreateBrandRestoreFile)
				}
			}

			req := apiClient.CustomizationAPI.CreateBrand(apiClient.GetConfig().Context)

			if CreateBranddata != "" {
				req = req.Data(CreateBranddata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateBrandQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateBrandQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateBranddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateBrandRestoreFile, "restore-from", "r", "", "Restore Customization from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateBrandQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateBrandCmd := NewCreateBrandCmd()
	CustomizationCmd.AddCommand(CreateBrandCmd)
}

var (
	ListBrandsBackupDir string

	ListBrandsLimit    int32
	ListBrandsPage     string
	ListBrandsFetchAll bool

	ListBrandsQuiet bool
)

func NewListBrandsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBrands",
		Long: "List all Brands",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListBrands(apiClient.GetConfig().Context)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListBrandsQuiet {
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
					if !ListBrandsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListBrandsFetchAll || len(items) == 0 {
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

			if ListBrandsFetchAll && pageCount > 1 && !ListBrandsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListBrandsBackupDir, "customization", "listBrands")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListBrandsQuiet {
					fmt.Printf("Backing up Customizations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListBrandsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListBrandsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListBrandsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListBrandsQuiet {
					fmt.Printf("Successfully backed up %d/%d Customizations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListBrandsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().Int32VarP(&ListBrandsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListBrandsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListBrandsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Customizations to a directory")

	cmd.Flags().StringVarP(&ListBrandsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListBrandsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListBrandsCmd := NewListBrandsCmd()
	CustomizationCmd.AddCommand(ListBrandsCmd)
}

var (
	GetBrandbrandId string

	GetBrandBackupDir string

	GetBrandQuiet bool
)

func NewGetBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getBrand",
		Long: "Retrieve a Brand",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetBrand(apiClient.GetConfig().Context, GetBrandbrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetBrandQuiet {
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
				dirPath := filepath.Join(GetBrandBackupDir, "customization", "getBrand")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetBrandbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetBrandQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetBrandQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetBrandbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetBrandBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetBrandQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetBrandCmd := NewGetBrandCmd()
	CustomizationCmd.AddCommand(GetBrandCmd)
}

var (
	ReplaceBrandbrandId string

	ReplaceBranddata string

	ReplaceBrandQuiet bool
)

func NewReplaceBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceBrand",
		Long: "Replace a Brand",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceBrand(apiClient.GetConfig().Context, ReplaceBrandbrandId)

			if ReplaceBranddata != "" {
				req = req.Data(ReplaceBranddata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceBrandQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceBrandQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceBrandbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceBranddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceBrandQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceBrandCmd := NewReplaceBrandCmd()
	CustomizationCmd.AddCommand(ReplaceBrandCmd)
}

var (
	DeleteBrandbrandId string

	DeleteBrandQuiet bool
)

func NewDeleteBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrand",
		Long: "Delete a brand",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrand(apiClient.GetConfig().Context, DeleteBrandbrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteBrandQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteBrandQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteBrandbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolVarP(&DeleteBrandQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteBrandCmd := NewDeleteBrandCmd()
	CustomizationCmd.AddCommand(DeleteBrandCmd)
}

var (
	ListBrandDomainsbrandId string

	ListBrandDomainsBackupDir string

	ListBrandDomainsLimit    int32
	ListBrandDomainsPage     string
	ListBrandDomainsFetchAll bool

	ListBrandDomainsQuiet bool
)

func NewListBrandDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBrandDomains",
		Long: "List all Domains associated with a Brand",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListBrandDomains(apiClient.GetConfig().Context, ListBrandDomainsbrandId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListBrandDomainsQuiet {
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
					if !ListBrandDomainsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListBrandDomainsFetchAll || len(items) == 0 {
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

			if ListBrandDomainsFetchAll && pageCount > 1 && !ListBrandDomainsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListBrandDomainsBackupDir, "customization", "listBrandDomains")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListBrandDomainsQuiet {
					fmt.Printf("Backing up Customizations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListBrandDomainsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListBrandDomainsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListBrandDomainsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListBrandDomainsQuiet {
					fmt.Printf("Successfully backed up %d/%d Customizations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListBrandDomainsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListBrandDomainsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().Int32VarP(&ListBrandDomainsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListBrandDomainsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListBrandDomainsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Customizations to a directory")

	cmd.Flags().StringVarP(&ListBrandDomainsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListBrandDomainsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListBrandDomainsCmd := NewListBrandDomainsCmd()
	CustomizationCmd.AddCommand(ListBrandDomainsCmd)
}

var (
	GetErrorPagebrandId string

	GetErrorPageBackupDir string

	GetErrorPageQuiet bool
)

func NewGetErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getErrorPage",
		Long: "Retrieve the Error Page Sub-Resources",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetErrorPage(apiClient.GetConfig().Context, GetErrorPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetErrorPageQuiet {
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
				dirPath := filepath.Join(GetErrorPageBackupDir, "customization", "getErrorPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetErrorPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetErrorPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetErrorPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetErrorPageCmd := NewGetErrorPageCmd()
	CustomizationCmd.AddCommand(GetErrorPageCmd)
}

var (
	GetCustomizedErrorPagebrandId string

	GetCustomizedErrorPageBackupDir string

	GetCustomizedErrorPageQuiet bool
)

func NewGetCustomizedErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCustomizedErrorPage",
		Long: "Retrieve the Customized Error Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetCustomizedErrorPage(apiClient.GetConfig().Context, GetCustomizedErrorPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCustomizedErrorPageQuiet {
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
				dirPath := filepath.Join(GetCustomizedErrorPageBackupDir, "customization", "getCustomizedErrorPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCustomizedErrorPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCustomizedErrorPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCustomizedErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCustomizedErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetCustomizedErrorPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCustomizedErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCustomizedErrorPageCmd := NewGetCustomizedErrorPageCmd()
	CustomizationCmd.AddCommand(GetCustomizedErrorPageCmd)
}

var (
	ReplaceCustomizedErrorPagebrandId string

	ReplaceCustomizedErrorPagedata string

	ReplaceCustomizedErrorPageQuiet bool
)

func NewReplaceCustomizedErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceCustomizedErrorPage",
		Long: "Replace the Customized Error Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceCustomizedErrorPage(apiClient.GetConfig().Context, ReplaceCustomizedErrorPagebrandId)

			if ReplaceCustomizedErrorPagedata != "" {
				req = req.Data(ReplaceCustomizedErrorPagedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceCustomizedErrorPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceCustomizedErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceCustomizedErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceCustomizedErrorPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceCustomizedErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceCustomizedErrorPageCmd := NewReplaceCustomizedErrorPageCmd()
	CustomizationCmd.AddCommand(ReplaceCustomizedErrorPageCmd)
}

var (
	DeleteCustomizedErrorPagebrandId string

	DeleteCustomizedErrorPageQuiet bool
)

func NewDeleteCustomizedErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteCustomizedErrorPage",
		Long: "Delete the Customized Error Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteCustomizedErrorPage(apiClient.GetConfig().Context, DeleteCustomizedErrorPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteCustomizedErrorPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteCustomizedErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteCustomizedErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolVarP(&DeleteCustomizedErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteCustomizedErrorPageCmd := NewDeleteCustomizedErrorPageCmd()
	CustomizationCmd.AddCommand(DeleteCustomizedErrorPageCmd)
}

var (
	GetDefaultErrorPagebrandId string

	GetDefaultErrorPageBackupDir string

	GetDefaultErrorPageQuiet bool
)

func NewGetDefaultErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDefaultErrorPage",
		Long: "Retrieve the Default Error Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetDefaultErrorPage(apiClient.GetConfig().Context, GetDefaultErrorPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetDefaultErrorPageQuiet {
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
				dirPath := filepath.Join(GetDefaultErrorPageBackupDir, "customization", "getDefaultErrorPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetDefaultErrorPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetDefaultErrorPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetDefaultErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetDefaultErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetDefaultErrorPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetDefaultErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetDefaultErrorPageCmd := NewGetDefaultErrorPageCmd()
	CustomizationCmd.AddCommand(GetDefaultErrorPageCmd)
}

var (
	GetPreviewErrorPagebrandId string

	GetPreviewErrorPageBackupDir string

	GetPreviewErrorPageQuiet bool
)

func NewGetPreviewErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPreviewErrorPage",
		Long: "Retrieve the Preview Error Page Preview",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetPreviewErrorPage(apiClient.GetConfig().Context, GetPreviewErrorPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPreviewErrorPageQuiet {
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
				dirPath := filepath.Join(GetPreviewErrorPageBackupDir, "customization", "getPreviewErrorPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPreviewErrorPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPreviewErrorPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPreviewErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPreviewErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetPreviewErrorPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPreviewErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPreviewErrorPageCmd := NewGetPreviewErrorPageCmd()
	CustomizationCmd.AddCommand(GetPreviewErrorPageCmd)
}

var (
	ReplacePreviewErrorPagebrandId string

	ReplacePreviewErrorPagedata string

	ReplacePreviewErrorPageQuiet bool
)

func NewReplacePreviewErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePreviewErrorPage",
		Long: "Replace the Preview Error Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplacePreviewErrorPage(apiClient.GetConfig().Context, ReplacePreviewErrorPagebrandId)

			if ReplacePreviewErrorPagedata != "" {
				req = req.Data(ReplacePreviewErrorPagedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacePreviewErrorPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacePreviewErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePreviewErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplacePreviewErrorPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacePreviewErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacePreviewErrorPageCmd := NewReplacePreviewErrorPageCmd()
	CustomizationCmd.AddCommand(ReplacePreviewErrorPageCmd)
}

var (
	DeletePreviewErrorPagebrandId string

	DeletePreviewErrorPageQuiet bool
)

func NewDeletePreviewErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletePreviewErrorPage",
		Long: "Delete the Preview Error Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeletePreviewErrorPage(apiClient.GetConfig().Context, DeletePreviewErrorPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeletePreviewErrorPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeletePreviewErrorPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeletePreviewErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolVarP(&DeletePreviewErrorPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeletePreviewErrorPageCmd := NewDeletePreviewErrorPageCmd()
	CustomizationCmd.AddCommand(DeletePreviewErrorPageCmd)
}

var (
	GetSignInPagebrandId string

	GetSignInPageBackupDir string

	GetSignInPageQuiet bool
)

func NewGetSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSignInPage",
		Long: "Retrieve the Sign-in Page Sub-Resources",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetSignInPage(apiClient.GetConfig().Context, GetSignInPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetSignInPageQuiet {
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
				dirPath := filepath.Join(GetSignInPageBackupDir, "customization", "getSignInPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetSignInPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetSignInPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetSignInPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetSignInPageCmd := NewGetSignInPageCmd()
	CustomizationCmd.AddCommand(GetSignInPageCmd)
}

var (
	GetCustomizedSignInPagebrandId string

	GetCustomizedSignInPageBackupDir string

	GetCustomizedSignInPageQuiet bool
)

func NewGetCustomizedSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCustomizedSignInPage",
		Long: "Retrieve the Customized Sign-in Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetCustomizedSignInPage(apiClient.GetConfig().Context, GetCustomizedSignInPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCustomizedSignInPageQuiet {
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
				dirPath := filepath.Join(GetCustomizedSignInPageBackupDir, "customization", "getCustomizedSignInPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCustomizedSignInPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCustomizedSignInPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCustomizedSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCustomizedSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetCustomizedSignInPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCustomizedSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCustomizedSignInPageCmd := NewGetCustomizedSignInPageCmd()
	CustomizationCmd.AddCommand(GetCustomizedSignInPageCmd)
}

var (
	ReplaceCustomizedSignInPagebrandId string

	ReplaceCustomizedSignInPagedata string

	ReplaceCustomizedSignInPageQuiet bool
)

func NewReplaceCustomizedSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceCustomizedSignInPage",
		Long: "Replace the Customized Sign-in Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceCustomizedSignInPage(apiClient.GetConfig().Context, ReplaceCustomizedSignInPagebrandId)

			if ReplaceCustomizedSignInPagedata != "" {
				req = req.Data(ReplaceCustomizedSignInPagedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceCustomizedSignInPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceCustomizedSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceCustomizedSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceCustomizedSignInPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceCustomizedSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceCustomizedSignInPageCmd := NewReplaceCustomizedSignInPageCmd()
	CustomizationCmd.AddCommand(ReplaceCustomizedSignInPageCmd)
}

var (
	DeleteCustomizedSignInPagebrandId string

	DeleteCustomizedSignInPageQuiet bool
)

func NewDeleteCustomizedSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteCustomizedSignInPage",
		Long: "Delete the Customized Sign-in Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteCustomizedSignInPage(apiClient.GetConfig().Context, DeleteCustomizedSignInPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteCustomizedSignInPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteCustomizedSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteCustomizedSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolVarP(&DeleteCustomizedSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteCustomizedSignInPageCmd := NewDeleteCustomizedSignInPageCmd()
	CustomizationCmd.AddCommand(DeleteCustomizedSignInPageCmd)
}

var (
	GetDefaultSignInPagebrandId string

	GetDefaultSignInPageBackupDir string

	GetDefaultSignInPageQuiet bool
)

func NewGetDefaultSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDefaultSignInPage",
		Long: "Retrieve the Default Sign-in Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetDefaultSignInPage(apiClient.GetConfig().Context, GetDefaultSignInPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetDefaultSignInPageQuiet {
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
				dirPath := filepath.Join(GetDefaultSignInPageBackupDir, "customization", "getDefaultSignInPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetDefaultSignInPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetDefaultSignInPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetDefaultSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetDefaultSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetDefaultSignInPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetDefaultSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetDefaultSignInPageCmd := NewGetDefaultSignInPageCmd()
	CustomizationCmd.AddCommand(GetDefaultSignInPageCmd)
}

var (
	GetPreviewSignInPagebrandId string

	GetPreviewSignInPageBackupDir string

	GetPreviewSignInPageQuiet bool
)

func NewGetPreviewSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPreviewSignInPage",
		Long: "Retrieve the Preview Sign-in Page Preview",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetPreviewSignInPage(apiClient.GetConfig().Context, GetPreviewSignInPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetPreviewSignInPageQuiet {
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
				dirPath := filepath.Join(GetPreviewSignInPageBackupDir, "customization", "getPreviewSignInPage")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetPreviewSignInPagebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetPreviewSignInPageQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetPreviewSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetPreviewSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetPreviewSignInPageBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetPreviewSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetPreviewSignInPageCmd := NewGetPreviewSignInPageCmd()
	CustomizationCmd.AddCommand(GetPreviewSignInPageCmd)
}

var (
	ReplacePreviewSignInPagebrandId string

	ReplacePreviewSignInPagedata string

	ReplacePreviewSignInPageQuiet bool
)

func NewReplacePreviewSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePreviewSignInPage",
		Long: "Replace the Preview Sign-in Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplacePreviewSignInPage(apiClient.GetConfig().Context, ReplacePreviewSignInPagebrandId)

			if ReplacePreviewSignInPagedata != "" {
				req = req.Data(ReplacePreviewSignInPagedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplacePreviewSignInPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplacePreviewSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePreviewSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplacePreviewSignInPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplacePreviewSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplacePreviewSignInPageCmd := NewReplacePreviewSignInPageCmd()
	CustomizationCmd.AddCommand(ReplacePreviewSignInPageCmd)
}

var (
	DeletePreviewSignInPagebrandId string

	DeletePreviewSignInPageQuiet bool
)

func NewDeletePreviewSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletePreviewSignInPage",
		Long: "Delete the Preview Sign-in Page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeletePreviewSignInPage(apiClient.GetConfig().Context, DeletePreviewSignInPagebrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeletePreviewSignInPageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeletePreviewSignInPageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeletePreviewSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolVarP(&DeletePreviewSignInPageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeletePreviewSignInPageCmd := NewDeletePreviewSignInPageCmd()
	CustomizationCmd.AddCommand(DeletePreviewSignInPageCmd)
}

var (
	ListAllSignInWidgetVersionsbrandId string

	ListAllSignInWidgetVersionsBackupDir string

	ListAllSignInWidgetVersionsLimit    int32
	ListAllSignInWidgetVersionsPage     string
	ListAllSignInWidgetVersionsFetchAll bool

	ListAllSignInWidgetVersionsQuiet bool
)

func NewListAllSignInWidgetVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAllSignInWidgetVersions",
		Long: "List all Sign-in Widget Versions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListAllSignInWidgetVersions(apiClient.GetConfig().Context, ListAllSignInWidgetVersionsbrandId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListAllSignInWidgetVersionsQuiet {
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
					if !ListAllSignInWidgetVersionsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListAllSignInWidgetVersionsFetchAll || len(items) == 0 {
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

			if ListAllSignInWidgetVersionsFetchAll && pageCount > 1 && !ListAllSignInWidgetVersionsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListAllSignInWidgetVersionsBackupDir, "customization", "listAllSignInWidgetVersions")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListAllSignInWidgetVersionsQuiet {
					fmt.Printf("Backing up Customizations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListAllSignInWidgetVersionsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListAllSignInWidgetVersionsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListAllSignInWidgetVersionsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListAllSignInWidgetVersionsQuiet {
					fmt.Printf("Successfully backed up %d/%d Customizations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListAllSignInWidgetVersionsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListAllSignInWidgetVersionsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().Int32VarP(&ListAllSignInWidgetVersionsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListAllSignInWidgetVersionsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListAllSignInWidgetVersionsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Customizations to a directory")

	cmd.Flags().StringVarP(&ListAllSignInWidgetVersionsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListAllSignInWidgetVersionsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListAllSignInWidgetVersionsCmd := NewListAllSignInWidgetVersionsCmd()
	CustomizationCmd.AddCommand(ListAllSignInWidgetVersionsCmd)
}

var (
	GetSignOutPageSettingsbrandId string

	GetSignOutPageSettingsBackupDir string

	GetSignOutPageSettingsQuiet bool
)

func NewGetSignOutPageSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSignOutPageSettings",
		Long: "Retrieve the Sign-out Page Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetSignOutPageSettings(apiClient.GetConfig().Context, GetSignOutPageSettingsbrandId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetSignOutPageSettingsQuiet {
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
				dirPath := filepath.Join(GetSignOutPageSettingsBackupDir, "customization", "getSignOutPageSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetSignOutPageSettingsbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetSignOutPageSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetSignOutPageSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetSignOutPageSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetSignOutPageSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetSignOutPageSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetSignOutPageSettingsCmd := NewGetSignOutPageSettingsCmd()
	CustomizationCmd.AddCommand(GetSignOutPageSettingsCmd)
}

var (
	ReplaceSignOutPageSettingsbrandId string

	ReplaceSignOutPageSettingsdata string

	ReplaceSignOutPageSettingsQuiet bool
)

func NewReplaceSignOutPageSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceSignOutPageSettings",
		Long: "Replace the Sign-out Page Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceSignOutPageSettings(apiClient.GetConfig().Context, ReplaceSignOutPageSettingsbrandId)

			if ReplaceSignOutPageSettingsdata != "" {
				req = req.Data(ReplaceSignOutPageSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceSignOutPageSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceSignOutPageSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceSignOutPageSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceSignOutPageSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceSignOutPageSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceSignOutPageSettingsCmd := NewReplaceSignOutPageSettingsCmd()
	CustomizationCmd.AddCommand(ReplaceSignOutPageSettingsCmd)
}

var (
	ListEmailTemplatesbrandId string

	ListEmailTemplatesBackupDir string

	ListEmailTemplatesLimit    int32
	ListEmailTemplatesPage     string
	ListEmailTemplatesFetchAll bool

	ListEmailTemplatesQuiet bool
)

func NewListEmailTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listEmailTemplates",
		Long: "List all Email Templates",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListEmailTemplates(apiClient.GetConfig().Context, ListEmailTemplatesbrandId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListEmailTemplatesQuiet {
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
					if !ListEmailTemplatesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListEmailTemplatesFetchAll || len(items) == 0 {
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

			if ListEmailTemplatesFetchAll && pageCount > 1 && !ListEmailTemplatesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListEmailTemplatesBackupDir, "customization", "listEmailTemplates")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListEmailTemplatesQuiet {
					fmt.Printf("Backing up Customizations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListEmailTemplatesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListEmailTemplatesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListEmailTemplatesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListEmailTemplatesQuiet {
					fmt.Printf("Successfully backed up %d/%d Customizations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListEmailTemplatesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListEmailTemplatesbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().Int32VarP(&ListEmailTemplatesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListEmailTemplatesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListEmailTemplatesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Customizations to a directory")

	cmd.Flags().StringVarP(&ListEmailTemplatesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListEmailTemplatesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListEmailTemplatesCmd := NewListEmailTemplatesCmd()
	CustomizationCmd.AddCommand(ListEmailTemplatesCmd)
}

var (
	GetEmailTemplatebrandId string

	GetEmailTemplatetemplateName string

	GetEmailTemplateBackupDir string

	GetEmailTemplateQuiet bool
)

func NewGetEmailTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailTemplate",
		Long: "Retrieve an Email Template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailTemplate(apiClient.GetConfig().Context, GetEmailTemplatebrandId, GetEmailTemplatetemplateName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailTemplateQuiet {
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
				dirPath := filepath.Join(GetEmailTemplateBackupDir, "customization", "getEmailTemplate")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailTemplatebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailTemplateQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailTemplateQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailTemplatebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailTemplatetemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetEmailTemplateBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailTemplateQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailTemplateCmd := NewGetEmailTemplateCmd()
	CustomizationCmd.AddCommand(GetEmailTemplateCmd)
}

var (
	CreateEmailCustomizationbrandId string

	CreateEmailCustomizationtemplateName string

	CreateEmailCustomizationdata string

	CreateEmailCustomizationRestoreFile string

	CreateEmailCustomizationQuiet bool
)

func NewCreateEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createEmail",
		Long: "Create an Email Customization",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if CreateEmailCustomizationRestoreFile != "" {

				jsonData, err := os.ReadFile(CreateEmailCustomizationRestoreFile)
				if err != nil {
					return fmt.Errorf("failed to read restore file: %w", err)
				}

				processedData, err := utils.PrepareDataForRestore(jsonData)
				if err != nil {
					return fmt.Errorf("failed to process restore data: %w", err)
				}

				CreateEmailCustomizationdata = string(processedData)

				if !CreateEmailCustomizationQuiet {
					fmt.Println("Restoring Customization from:", CreateEmailCustomizationRestoreFile)
				}
			}

			req := apiClient.CustomizationAPI.CreateEmailCustomization(apiClient.GetConfig().Context, CreateEmailCustomizationbrandId, CreateEmailCustomizationtemplateName)

			if CreateEmailCustomizationdata != "" {
				req = req.Data(CreateEmailCustomizationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !CreateEmailCustomizationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !CreateEmailCustomizationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&CreateEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&CreateEmailCustomizationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().StringVarP(&CreateEmailCustomizationRestoreFile, "restore-from", "r", "", "Restore Customization from a JSON backup file")

	cmd.Flags().BoolVarP(&CreateEmailCustomizationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	CreateEmailCustomizationCmd := NewCreateEmailCustomizationCmd()
	CustomizationCmd.AddCommand(CreateEmailCustomizationCmd)
}

var (
	ListEmailCustomizationsbrandId string

	ListEmailCustomizationstemplateName string

	ListEmailCustomizationsBackupDir string

	ListEmailCustomizationsLimit    int32
	ListEmailCustomizationsPage     string
	ListEmailCustomizationsFetchAll bool

	ListEmailCustomizationsQuiet bool
)

func NewListEmailCustomizationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listEmails",
		Long: "List all Email Customizations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListEmailCustomizations(apiClient.GetConfig().Context, ListEmailCustomizationsbrandId, ListEmailCustomizationstemplateName)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListEmailCustomizationsQuiet {
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
					if !ListEmailCustomizationsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListEmailCustomizationsFetchAll || len(items) == 0 {
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

			if ListEmailCustomizationsFetchAll && pageCount > 1 && !ListEmailCustomizationsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListEmailCustomizationsBackupDir, "customization", "listEmails")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListEmailCustomizationsQuiet {
					fmt.Printf("Backing up Customizations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListEmailCustomizationsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListEmailCustomizationsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListEmailCustomizationsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListEmailCustomizationsQuiet {
					fmt.Printf("Successfully backed up %d/%d Customizations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListEmailCustomizationsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListEmailCustomizationsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ListEmailCustomizationstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().Int32VarP(&ListEmailCustomizationsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListEmailCustomizationsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListEmailCustomizationsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Customizations to a directory")

	cmd.Flags().StringVarP(&ListEmailCustomizationsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListEmailCustomizationsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListEmailCustomizationsCmd := NewListEmailCustomizationsCmd()
	CustomizationCmd.AddCommand(ListEmailCustomizationsCmd)
}

var (
	DeleteAllCustomizationsbrandId string

	DeleteAllCustomizationstemplateName string

	DeleteAllCustomizationsQuiet bool
)

func NewDeleteAllCustomizationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteAlls",
		Long: "Delete all Email Customizations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteAllCustomizations(apiClient.GetConfig().Context, DeleteAllCustomizationsbrandId, DeleteAllCustomizationstemplateName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteAllCustomizationsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteAllCustomizationsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteAllCustomizationsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteAllCustomizationstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().BoolVarP(&DeleteAllCustomizationsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteAllCustomizationsCmd := NewDeleteAllCustomizationsCmd()
	CustomizationCmd.AddCommand(DeleteAllCustomizationsCmd)
}

var (
	GetEmailCustomizationbrandId string

	GetEmailCustomizationtemplateName string

	GetEmailCustomizationcustomizationId string

	GetEmailCustomizationBackupDir string

	GetEmailCustomizationQuiet bool
)

func NewGetEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmail",
		Long: "Retrieve an Email Customization",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailCustomization(apiClient.GetConfig().Context, GetEmailCustomizationbrandId, GetEmailCustomizationtemplateName, GetEmailCustomizationcustomizationId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailCustomizationQuiet {
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
				dirPath := filepath.Join(GetEmailCustomizationBackupDir, "customization", "getEmail")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailCustomizationbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailCustomizationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailCustomizationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&GetEmailCustomizationcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetEmailCustomizationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailCustomizationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailCustomizationCmd := NewGetEmailCustomizationCmd()
	CustomizationCmd.AddCommand(GetEmailCustomizationCmd)
}

var (
	ReplaceEmailCustomizationbrandId string

	ReplaceEmailCustomizationtemplateName string

	ReplaceEmailCustomizationcustomizationId string

	ReplaceEmailCustomizationdata string

	ReplaceEmailCustomizationQuiet bool
)

func NewReplaceEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceEmail",
		Long: "Replace an Email Customization",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceEmailCustomization(apiClient.GetConfig().Context, ReplaceEmailCustomizationbrandId, ReplaceEmailCustomizationtemplateName, ReplaceEmailCustomizationcustomizationId)

			if ReplaceEmailCustomizationdata != "" {
				req = req.Data(ReplaceEmailCustomizationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceEmailCustomizationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceEmailCustomizationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceEmailCustomizationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceEmailCustomizationCmd := NewReplaceEmailCustomizationCmd()
	CustomizationCmd.AddCommand(ReplaceEmailCustomizationCmd)
}

var (
	DeleteEmailCustomizationbrandId string

	DeleteEmailCustomizationtemplateName string

	DeleteEmailCustomizationcustomizationId string

	DeleteEmailCustomizationQuiet bool
)

func NewDeleteEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteEmail",
		Long: "Delete an Email Customization",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteEmailCustomization(apiClient.GetConfig().Context, DeleteEmailCustomizationbrandId, DeleteEmailCustomizationtemplateName, DeleteEmailCustomizationcustomizationId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteEmailCustomizationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteEmailCustomizationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&DeleteEmailCustomizationcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	cmd.Flags().BoolVarP(&DeleteEmailCustomizationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteEmailCustomizationCmd := NewDeleteEmailCustomizationCmd()
	CustomizationCmd.AddCommand(DeleteEmailCustomizationCmd)
}

var (
	GetCustomizationPreviewbrandId string

	GetCustomizationPreviewtemplateName string

	GetCustomizationPreviewcustomizationId string

	GetCustomizationPreviewBackupDir string

	GetCustomizationPreviewQuiet bool
)

func NewGetCustomizationPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPreview",
		Long: "Retrieve a Preview of an Email Customization",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetCustomizationPreview(apiClient.GetConfig().Context, GetCustomizationPreviewbrandId, GetCustomizationPreviewtemplateName, GetCustomizationPreviewcustomizationId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetCustomizationPreviewQuiet {
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
				dirPath := filepath.Join(GetCustomizationPreviewBackupDir, "customization", "getPreview")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetCustomizationPreviewbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetCustomizationPreviewQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetCustomizationPreviewQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetCustomizationPreviewbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetCustomizationPreviewtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&GetCustomizationPreviewcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetCustomizationPreviewBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetCustomizationPreviewQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetCustomizationPreviewCmd := NewGetCustomizationPreviewCmd()
	CustomizationCmd.AddCommand(GetCustomizationPreviewCmd)
}

var (
	GetEmailDefaultContentbrandId string

	GetEmailDefaultContenttemplateName string

	GetEmailDefaultContentBackupDir string

	GetEmailDefaultContentQuiet bool
)

func NewGetEmailDefaultContentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailDefaultContent",
		Long: "Retrieve an Email Template Default Content",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailDefaultContent(apiClient.GetConfig().Context, GetEmailDefaultContentbrandId, GetEmailDefaultContenttemplateName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailDefaultContentQuiet {
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
				dirPath := filepath.Join(GetEmailDefaultContentBackupDir, "customization", "getEmailDefaultContent")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailDefaultContentbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailDefaultContentQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailDefaultContentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailDefaultContentbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailDefaultContenttemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetEmailDefaultContentBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailDefaultContentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailDefaultContentCmd := NewGetEmailDefaultContentCmd()
	CustomizationCmd.AddCommand(GetEmailDefaultContentCmd)
}

var (
	GetEmailDefaultPreviewbrandId string

	GetEmailDefaultPreviewtemplateName string

	GetEmailDefaultPreviewBackupDir string

	GetEmailDefaultPreviewQuiet bool
)

func NewGetEmailDefaultPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailDefaultPreview",
		Long: "Retrieve a Preview of the Email Template default content",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailDefaultPreview(apiClient.GetConfig().Context, GetEmailDefaultPreviewbrandId, GetEmailDefaultPreviewtemplateName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailDefaultPreviewQuiet {
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
				dirPath := filepath.Join(GetEmailDefaultPreviewBackupDir, "customization", "getEmailDefaultPreview")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailDefaultPreviewbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailDefaultPreviewQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailDefaultPreviewQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailDefaultPreviewbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailDefaultPreviewtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetEmailDefaultPreviewBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailDefaultPreviewQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailDefaultPreviewCmd := NewGetEmailDefaultPreviewCmd()
	CustomizationCmd.AddCommand(GetEmailDefaultPreviewCmd)
}

var (
	GetEmailSettingsbrandId string

	GetEmailSettingstemplateName string

	GetEmailSettingsBackupDir string

	GetEmailSettingsQuiet bool
)

func NewGetEmailSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailSettings",
		Long: "Retrieve the Email Template Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailSettings(apiClient.GetConfig().Context, GetEmailSettingsbrandId, GetEmailSettingstemplateName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetEmailSettingsQuiet {
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
				dirPath := filepath.Join(GetEmailSettingsBackupDir, "customization", "getEmailSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetEmailSettingsbrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetEmailSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetEmailSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEmailSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailSettingstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetEmailSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetEmailSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetEmailSettingsCmd := NewGetEmailSettingsCmd()
	CustomizationCmd.AddCommand(GetEmailSettingsCmd)
}

var (
	ReplaceEmailSettingsbrandId string

	ReplaceEmailSettingstemplateName string

	ReplaceEmailSettingsdata string

	ReplaceEmailSettingsQuiet bool
)

func NewReplaceEmailSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceEmailSettings",
		Long: "Replace the Email Template Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceEmailSettings(apiClient.GetConfig().Context, ReplaceEmailSettingsbrandId, ReplaceEmailSettingstemplateName)

			if ReplaceEmailSettingsdata != "" {
				req = req.Data(ReplaceEmailSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceEmailSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceEmailSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceEmailSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceEmailSettingstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&ReplaceEmailSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceEmailSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceEmailSettingsCmd := NewReplaceEmailSettingsCmd()
	CustomizationCmd.AddCommand(ReplaceEmailSettingsCmd)
}

var (
	SendTestEmailbrandId string

	SendTestEmailtemplateName string

	SendTestEmailQuiet bool
)

func NewSendTestEmailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "sendTestEmail",
		Long: "Send a Test Email",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.SendTestEmail(apiClient.GetConfig().Context, SendTestEmailbrandId, SendTestEmailtemplateName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !SendTestEmailQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !SendTestEmailQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&SendTestEmailbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&SendTestEmailtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().BoolVarP(&SendTestEmailQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	SendTestEmailCmd := NewSendTestEmailCmd()
	CustomizationCmd.AddCommand(SendTestEmailCmd)
}

var (
	ListBrandThemesbrandId string

	ListBrandThemesBackupDir string

	ListBrandThemesLimit    int32
	ListBrandThemesPage     string
	ListBrandThemesFetchAll bool

	ListBrandThemesQuiet bool
)

func NewListBrandThemesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBrandThemes",
		Long: "List all Themes",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListBrandThemes(apiClient.GetConfig().Context, ListBrandThemesbrandId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListBrandThemesQuiet {
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
					if !ListBrandThemesQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListBrandThemesFetchAll || len(items) == 0 {
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

			if ListBrandThemesFetchAll && pageCount > 1 && !ListBrandThemesQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListBrandThemesBackupDir, "customization", "listBrandThemes")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListBrandThemesQuiet {
					fmt.Printf("Backing up Customizations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListBrandThemesQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListBrandThemesQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListBrandThemesQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListBrandThemesQuiet {
					fmt.Printf("Successfully backed up %d/%d Customizations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListBrandThemesQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListBrandThemesbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().Int32VarP(&ListBrandThemesLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListBrandThemesPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListBrandThemesFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple Customizations to a directory")

	cmd.Flags().StringVarP(&ListBrandThemesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListBrandThemesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListBrandThemesCmd := NewListBrandThemesCmd()
	CustomizationCmd.AddCommand(ListBrandThemesCmd)
}

var (
	GetBrandThemebrandId string

	GetBrandThemethemeId string

	GetBrandThemeBackupDir string

	GetBrandThemeQuiet bool
)

func NewGetBrandThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getBrandTheme",
		Long: "Retrieve a Theme",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetBrandTheme(apiClient.GetConfig().Context, GetBrandThemebrandId, GetBrandThemethemeId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetBrandThemeQuiet {
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
				dirPath := filepath.Join(GetBrandThemeBackupDir, "customization", "getBrandTheme")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetBrandThemebrandId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetBrandThemeQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetBrandThemeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetBrandThemebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetBrandThemethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the Customization to a file")

	cmd.Flags().StringVarP(&GetBrandThemeBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetBrandThemeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetBrandThemeCmd := NewGetBrandThemeCmd()
	CustomizationCmd.AddCommand(GetBrandThemeCmd)
}

var (
	ReplaceBrandThemebrandId string

	ReplaceBrandThemethemeId string

	ReplaceBrandThemedata string

	ReplaceBrandThemeQuiet bool
)

func NewReplaceBrandThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceBrandTheme",
		Long: "Replace a Theme",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceBrandTheme(apiClient.GetConfig().Context, ReplaceBrandThemebrandId, ReplaceBrandThemethemeId)

			if ReplaceBrandThemedata != "" {
				req = req.Data(ReplaceBrandThemedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceBrandThemeQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceBrandThemeQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceBrandThemebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceBrandThemethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&ReplaceBrandThemedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceBrandThemeQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceBrandThemeCmd := NewReplaceBrandThemeCmd()
	CustomizationCmd.AddCommand(ReplaceBrandThemeCmd)
}

var (
	UploadBrandThemeBackgroundImagebrandId string

	UploadBrandThemeBackgroundImagethemeId string

	UploadBrandThemeBackgroundImagedata string

	UploadBrandThemeBackgroundImageQuiet bool
)

func NewUploadBrandThemeBackgroundImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadBrandThemeBackgroundImage",
		Long: "Upload the Background Image",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.UploadBrandThemeBackgroundImage(apiClient.GetConfig().Context, UploadBrandThemeBackgroundImagebrandId, UploadBrandThemeBackgroundImagethemeId)

			if UploadBrandThemeBackgroundImagedata != "" {
				req = req.Data(UploadBrandThemeBackgroundImagedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadBrandThemeBackgroundImageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadBrandThemeBackgroundImageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadBrandThemeBackgroundImagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&UploadBrandThemeBackgroundImagethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&UploadBrandThemeBackgroundImagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadBrandThemeBackgroundImageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadBrandThemeBackgroundImageCmd := NewUploadBrandThemeBackgroundImageCmd()
	CustomizationCmd.AddCommand(UploadBrandThemeBackgroundImageCmd)
}

var (
	DeleteBrandThemeBackgroundImagebrandId string

	DeleteBrandThemeBackgroundImagethemeId string

	DeleteBrandThemeBackgroundImageQuiet bool
)

func NewDeleteBrandThemeBackgroundImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrandThemeBackgroundImage",
		Long: "Delete the Background Image",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrandThemeBackgroundImage(apiClient.GetConfig().Context, DeleteBrandThemeBackgroundImagebrandId, DeleteBrandThemeBackgroundImagethemeId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteBrandThemeBackgroundImageQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteBrandThemeBackgroundImageQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteBrandThemeBackgroundImagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteBrandThemeBackgroundImagethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().BoolVarP(&DeleteBrandThemeBackgroundImageQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteBrandThemeBackgroundImageCmd := NewDeleteBrandThemeBackgroundImageCmd()
	CustomizationCmd.AddCommand(DeleteBrandThemeBackgroundImageCmd)
}

var (
	UploadBrandThemeFaviconbrandId string

	UploadBrandThemeFaviconthemeId string

	UploadBrandThemeFavicondata string

	UploadBrandThemeFaviconQuiet bool
)

func NewUploadBrandThemeFaviconCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadBrandThemeFavicon",
		Long: "Upload the Favicon",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.UploadBrandThemeFavicon(apiClient.GetConfig().Context, UploadBrandThemeFaviconbrandId, UploadBrandThemeFaviconthemeId)

			if UploadBrandThemeFavicondata != "" {
				req = req.Data(UploadBrandThemeFavicondata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadBrandThemeFaviconQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadBrandThemeFaviconQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadBrandThemeFaviconbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&UploadBrandThemeFaviconthemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&UploadBrandThemeFavicondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadBrandThemeFaviconQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadBrandThemeFaviconCmd := NewUploadBrandThemeFaviconCmd()
	CustomizationCmd.AddCommand(UploadBrandThemeFaviconCmd)
}

var (
	DeleteBrandThemeFaviconbrandId string

	DeleteBrandThemeFaviconthemeId string

	DeleteBrandThemeFaviconQuiet bool
)

func NewDeleteBrandThemeFaviconCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrandThemeFavicon",
		Long: "Delete the Favicon",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrandThemeFavicon(apiClient.GetConfig().Context, DeleteBrandThemeFaviconbrandId, DeleteBrandThemeFaviconthemeId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteBrandThemeFaviconQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteBrandThemeFaviconQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteBrandThemeFaviconbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteBrandThemeFaviconthemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().BoolVarP(&DeleteBrandThemeFaviconQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteBrandThemeFaviconCmd := NewDeleteBrandThemeFaviconCmd()
	CustomizationCmd.AddCommand(DeleteBrandThemeFaviconCmd)
}

var (
	UploadBrandThemeLogobrandId string

	UploadBrandThemeLogothemeId string

	UploadBrandThemeLogodata string

	UploadBrandThemeLogoQuiet bool
)

func NewUploadBrandThemeLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadBrandThemeLogo",
		Long: "Upload the Logo",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.UploadBrandThemeLogo(apiClient.GetConfig().Context, UploadBrandThemeLogobrandId, UploadBrandThemeLogothemeId)

			if UploadBrandThemeLogodata != "" {
				req = req.Data(UploadBrandThemeLogodata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadBrandThemeLogoQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadBrandThemeLogoQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadBrandThemeLogobrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&UploadBrandThemeLogothemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&UploadBrandThemeLogodata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadBrandThemeLogoQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadBrandThemeLogoCmd := NewUploadBrandThemeLogoCmd()
	CustomizationCmd.AddCommand(UploadBrandThemeLogoCmd)
}

var (
	DeleteBrandThemeLogobrandId string

	DeleteBrandThemeLogothemeId string

	DeleteBrandThemeLogoQuiet bool
)

func NewDeleteBrandThemeLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrandThemeLogo",
		Long: "Delete the Logo",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrandThemeLogo(apiClient.GetConfig().Context, DeleteBrandThemeLogobrandId, DeleteBrandThemeLogothemeId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteBrandThemeLogoQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteBrandThemeLogoQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteBrandThemeLogobrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteBrandThemeLogothemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().BoolVarP(&DeleteBrandThemeLogoQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteBrandThemeLogoCmd := NewDeleteBrandThemeLogoCmd()
	CustomizationCmd.AddCommand(DeleteBrandThemeLogoCmd)
}
