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

var ApplicationFeaturesCmd = &cobra.Command{
	Use:  "applicationFeatures",
	Long: "Manage ApplicationFeaturesAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationFeaturesCmd)
}

var (
	ListFeaturesForApplicationappId string

	ListFeaturesForApplicationBackupDir string

	ListFeaturesForApplicationLimit    int32
	ListFeaturesForApplicationPage     string
	ListFeaturesForApplicationFetchAll bool

	ListFeaturesForApplicationQuiet bool
)

func NewListFeaturesForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listFeaturesForApplication",
		Long: "List all Features",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationFeaturesAPI.ListFeaturesForApplication(apiClient.GetConfig().Context, ListFeaturesForApplicationappId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListFeaturesForApplicationQuiet {
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
					if !ListFeaturesForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListFeaturesForApplicationFetchAll || len(items) == 0 {
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

			if ListFeaturesForApplicationFetchAll && pageCount > 1 && !ListFeaturesForApplicationQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListFeaturesForApplicationBackupDir, "applicationfeatures", "listFeaturesForApplication")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListFeaturesForApplicationQuiet {
					fmt.Printf("Backing up ApplicationFeaturess to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListFeaturesForApplicationQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListFeaturesForApplicationQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListFeaturesForApplicationQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListFeaturesForApplicationQuiet {
					fmt.Printf("Successfully backed up %d/%d ApplicationFeaturess\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListFeaturesForApplicationQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListFeaturesForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().Int32VarP(&ListFeaturesForApplicationLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListFeaturesForApplicationPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListFeaturesForApplicationFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple ApplicationFeaturess to a directory")

	cmd.Flags().StringVarP(&ListFeaturesForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListFeaturesForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListFeaturesForApplicationCmd := NewListFeaturesForApplicationCmd()
	ApplicationFeaturesCmd.AddCommand(ListFeaturesForApplicationCmd)
}

var (
	GetFeatureForApplicationappId string

	GetFeatureForApplicationfeatureName string

	GetFeatureForApplicationBackupDir string

	GetFeatureForApplicationQuiet bool
)

func NewGetFeatureForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getFeatureForApplication",
		Long: "Retrieve a Feature",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationFeaturesAPI.GetFeatureForApplication(apiClient.GetConfig().Context, GetFeatureForApplicationappId, GetFeatureForApplicationfeatureName)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetFeatureForApplicationQuiet {
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
				dirPath := filepath.Join(GetFeatureForApplicationBackupDir, "applicationfeatures", "getFeatureForApplication")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetFeatureForApplicationappId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetFeatureForApplicationQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetFeatureForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetFeatureForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetFeatureForApplicationfeatureName, "featureName", "", "", "")
	cmd.MarkFlagRequired("featureName")

	cmd.Flags().BoolP("backup", "b", false, "Backup the ApplicationFeatures to a file")

	cmd.Flags().StringVarP(&GetFeatureForApplicationBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetFeatureForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetFeatureForApplicationCmd := NewGetFeatureForApplicationCmd()
	ApplicationFeaturesCmd.AddCommand(GetFeatureForApplicationCmd)
}

var (
	UpdateFeatureForApplicationappId string

	UpdateFeatureForApplicationfeatureName string

	UpdateFeatureForApplicationdata string

	UpdateFeatureForApplicationQuiet bool
)

func NewUpdateFeatureForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateFeatureForApplication",
		Long: "Update a Feature",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationFeaturesAPI.UpdateFeatureForApplication(apiClient.GetConfig().Context, UpdateFeatureForApplicationappId, UpdateFeatureForApplicationfeatureName)

			if UpdateFeatureForApplicationdata != "" {
				req = req.Data(UpdateFeatureForApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateFeatureForApplicationQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateFeatureForApplicationQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateFeatureForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateFeatureForApplicationfeatureName, "featureName", "", "", "")
	cmd.MarkFlagRequired("featureName")

	cmd.Flags().StringVarP(&UpdateFeatureForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateFeatureForApplicationQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateFeatureForApplicationCmd := NewUpdateFeatureForApplicationCmd()
	ApplicationFeaturesCmd.AddCommand(UpdateFeatureForApplicationCmd)
}
