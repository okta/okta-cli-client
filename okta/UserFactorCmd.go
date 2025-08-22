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

var UserFactorCmd = &cobra.Command{
	Use:  "userFactor",
	Long: "Manage UserFactorAPI",
}

func init() {
	rootCmd.AddCommand(UserFactorCmd)
}

var (
	EnrollFactoruserId string

	EnrollFactordata string

	EnrollFactorQuiet bool
)

func NewEnrollFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "enrollFactor",
		Long: "Enroll a Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.EnrollFactor(apiClient.GetConfig().Context, EnrollFactoruserId)

			if EnrollFactordata != "" {
				req = req.Data(EnrollFactordata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !EnrollFactorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !EnrollFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&EnrollFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&EnrollFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&EnrollFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	EnrollFactorCmd := NewEnrollFactorCmd()
	UserFactorCmd.AddCommand(EnrollFactorCmd)
}

var (
	ListFactorsuserId string

	ListFactorsBackupDir string

	ListFactorsLimit    int32
	ListFactorsPage     string
	ListFactorsFetchAll bool

	ListFactorsQuiet bool
)

func NewListFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listFactors",
		Long: "List all enrolled Factors",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ListFactors(apiClient.GetConfig().Context, ListFactorsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListFactorsQuiet {
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
					if !ListFactorsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListFactorsFetchAll || len(items) == 0 {
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

			if ListFactorsFetchAll && pageCount > 1 && !ListFactorsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListFactorsBackupDir, "userfactor", "listFactors")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListFactorsQuiet {
					fmt.Printf("Backing up UserFactors to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListFactorsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListFactorsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListFactorsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListFactorsQuiet {
					fmt.Printf("Successfully backed up %d/%d UserFactors\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListFactorsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListFactorsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListFactorsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListFactorsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple UserFactors to a directory")

	cmd.Flags().StringVarP(&ListFactorsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListFactorsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListFactorsCmd := NewListFactorsCmd()
	UserFactorCmd.AddCommand(ListFactorsCmd)
}

var (
	ListSupportedFactorsuserId string

	ListSupportedFactorsBackupDir string

	ListSupportedFactorsLimit    int32
	ListSupportedFactorsPage     string
	ListSupportedFactorsFetchAll bool

	ListSupportedFactorsQuiet bool
)

func NewListSupportedFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSupportedFactors",
		Long: "List all supported Factors",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ListSupportedFactors(apiClient.GetConfig().Context, ListSupportedFactorsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListSupportedFactorsQuiet {
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
					if !ListSupportedFactorsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListSupportedFactorsFetchAll || len(items) == 0 {
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

			if ListSupportedFactorsFetchAll && pageCount > 1 && !ListSupportedFactorsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListSupportedFactorsBackupDir, "userfactor", "listSupportedFactors")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListSupportedFactorsQuiet {
					fmt.Printf("Backing up UserFactors to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListSupportedFactorsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListSupportedFactorsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListSupportedFactorsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListSupportedFactorsQuiet {
					fmt.Printf("Successfully backed up %d/%d UserFactors\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListSupportedFactorsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListSupportedFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListSupportedFactorsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListSupportedFactorsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListSupportedFactorsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple UserFactors to a directory")

	cmd.Flags().StringVarP(&ListSupportedFactorsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListSupportedFactorsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListSupportedFactorsCmd := NewListSupportedFactorsCmd()
	UserFactorCmd.AddCommand(ListSupportedFactorsCmd)
}

var (
	ListSupportedSecurityQuestionsuserId string

	ListSupportedSecurityQuestionsBackupDir string

	ListSupportedSecurityQuestionsLimit    int32
	ListSupportedSecurityQuestionsPage     string
	ListSupportedSecurityQuestionsFetchAll bool

	ListSupportedSecurityQuestionsQuiet bool
)

func NewListSupportedSecurityQuestionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSupportedSecurityQuestions",
		Long: "List all supported Security Questions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ListSupportedSecurityQuestions(apiClient.GetConfig().Context, ListSupportedSecurityQuestionsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListSupportedSecurityQuestionsQuiet {
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
					if !ListSupportedSecurityQuestionsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListSupportedSecurityQuestionsFetchAll || len(items) == 0 {
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

			if ListSupportedSecurityQuestionsFetchAll && pageCount > 1 && !ListSupportedSecurityQuestionsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListSupportedSecurityQuestionsBackupDir, "userfactor", "listSupportedSecurityQuestions")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListSupportedSecurityQuestionsQuiet {
					fmt.Printf("Backing up UserFactors to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListSupportedSecurityQuestionsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListSupportedSecurityQuestionsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListSupportedSecurityQuestionsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListSupportedSecurityQuestionsQuiet {
					fmt.Printf("Successfully backed up %d/%d UserFactors\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListSupportedSecurityQuestionsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListSupportedSecurityQuestionsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListSupportedSecurityQuestionsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListSupportedSecurityQuestionsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListSupportedSecurityQuestionsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple UserFactors to a directory")

	cmd.Flags().StringVarP(&ListSupportedSecurityQuestionsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListSupportedSecurityQuestionsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListSupportedSecurityQuestionsCmd := NewListSupportedSecurityQuestionsCmd()
	UserFactorCmd.AddCommand(ListSupportedSecurityQuestionsCmd)
}

var (
	GetFactoruserId string

	GetFactorfactorId string

	GetFactorBackupDir string

	GetFactorQuiet bool
)

func NewGetFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getFactor",
		Long: "Retrieve a Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.GetFactor(apiClient.GetConfig().Context, GetFactoruserId, GetFactorfactorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetFactorQuiet {
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
				dirPath := filepath.Join(GetFactorBackupDir, "userfactor", "getFactor")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetFactoruserId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetFactorQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the UserFactor to a file")

	cmd.Flags().StringVarP(&GetFactorBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetFactorCmd := NewGetFactorCmd()
	UserFactorCmd.AddCommand(GetFactorCmd)
}

var (
	UnenrollFactoruserId string

	UnenrollFactorfactorId string

	UnenrollFactorQuiet bool
)

func NewUnenrollFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unenrollFactor",
		Long: "Unenroll a Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.UnenrollFactor(apiClient.GetConfig().Context, UnenrollFactoruserId, UnenrollFactorfactorId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UnenrollFactorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UnenrollFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnenrollFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UnenrollFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().BoolVarP(&UnenrollFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UnenrollFactorCmd := NewUnenrollFactorCmd()
	UserFactorCmd.AddCommand(UnenrollFactorCmd)
}

var (
	ActivateFactoruserId string

	ActivateFactorfactorId string

	ActivateFactordata string

	ActivateFactorQuiet bool
)

func NewActivateFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateFactor",
		Long: "Activate a Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ActivateFactor(apiClient.GetConfig().Context, ActivateFactoruserId, ActivateFactorfactorId)

			if ActivateFactordata != "" {
				req = req.Data(ActivateFactordata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivateFactorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivateFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ActivateFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&ActivateFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ActivateFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivateFactorCmd := NewActivateFactorCmd()
	UserFactorCmd.AddCommand(ActivateFactorCmd)
}

var (
	ResendEnrollFactoruserId string

	ResendEnrollFactorfactorId string

	ResendEnrollFactordata string

	ResendEnrollFactorQuiet bool
)

func NewResendEnrollFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "resendEnrollFactor",
		Long: "Resend a Factor enrollment",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ResendEnrollFactor(apiClient.GetConfig().Context, ResendEnrollFactoruserId, ResendEnrollFactorfactorId)

			if ResendEnrollFactordata != "" {
				req = req.Data(ResendEnrollFactordata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ResendEnrollFactorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ResendEnrollFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ResendEnrollFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ResendEnrollFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&ResendEnrollFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ResendEnrollFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ResendEnrollFactorCmd := NewResendEnrollFactorCmd()
	UserFactorCmd.AddCommand(ResendEnrollFactorCmd)
}

var (
	GetFactorTransactionStatususerId string

	GetFactorTransactionStatusfactorId string

	GetFactorTransactionStatustransactionId string

	GetFactorTransactionStatusBackupDir string

	GetFactorTransactionStatusQuiet bool
)

func NewGetFactorTransactionStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getFactorTransactionStatus",
		Long: "Retrieve a Factor transaction status",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.GetFactorTransactionStatus(apiClient.GetConfig().Context, GetFactorTransactionStatususerId, GetFactorTransactionStatusfactorId, GetFactorTransactionStatustransactionId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetFactorTransactionStatusQuiet {
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
				dirPath := filepath.Join(GetFactorTransactionStatusBackupDir, "userfactor", "getFactorTransactionStatus")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetFactorTransactionStatususerId
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetFactorTransactionStatusQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetFactorTransactionStatusQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetFactorTransactionStatususerId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetFactorTransactionStatusfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&GetFactorTransactionStatustransactionId, "transactionId", "", "", "")
	cmd.MarkFlagRequired("transactionId")

	cmd.Flags().BoolP("backup", "b", false, "Backup the UserFactor to a file")

	cmd.Flags().StringVarP(&GetFactorTransactionStatusBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetFactorTransactionStatusQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetFactorTransactionStatusCmd := NewGetFactorTransactionStatusCmd()
	UserFactorCmd.AddCommand(GetFactorTransactionStatusCmd)
}

var (
	VerifyFactoruserId string

	VerifyFactorfactorId string

	VerifyFactordata string

	VerifyFactorQuiet bool
)

func NewVerifyFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verifyFactor",
		Long: "Verify a Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.VerifyFactor(apiClient.GetConfig().Context, VerifyFactoruserId, VerifyFactorfactorId)

			if VerifyFactordata != "" {
				req = req.Data(VerifyFactordata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !VerifyFactorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !VerifyFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&VerifyFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&VerifyFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&VerifyFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&VerifyFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	VerifyFactorCmd := NewVerifyFactorCmd()
	UserFactorCmd.AddCommand(VerifyFactorCmd)
}
