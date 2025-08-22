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

var WebAuthnPreregistrationCmd = &cobra.Command{
	Use:  "webAuthnPreregistration",
	Long: "Manage WebAuthnPreregistrationAPI",
}

func init() {
	rootCmd.AddCommand(WebAuthnPreregistrationCmd)
}

var (
	ActivatePreregistrationEnrollmentdata string

	ActivatePreregistrationEnrollmentQuiet bool
)

func NewActivatePreregistrationEnrollmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activatePreregistrationEnrollment",
		Long: "Activate a Preregistered WebAuthn Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.ActivatePreregistrationEnrollment(apiClient.GetConfig().Context)

			if ActivatePreregistrationEnrollmentdata != "" {
				req = req.Data(ActivatePreregistrationEnrollmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ActivatePreregistrationEnrollmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ActivatePreregistrationEnrollmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivatePreregistrationEnrollmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ActivatePreregistrationEnrollmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ActivatePreregistrationEnrollmentCmd := NewActivatePreregistrationEnrollmentCmd()
	WebAuthnPreregistrationCmd.AddCommand(ActivatePreregistrationEnrollmentCmd)
}

var (
	EnrollPreregistrationEnrollmentdata string

	EnrollPreregistrationEnrollmentQuiet bool
)

func NewEnrollPreregistrationEnrollmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "enrollPreregistrationEnrollment",
		Long: "Enroll a Preregistered WebAuthn Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.EnrollPreregistrationEnrollment(apiClient.GetConfig().Context)

			if EnrollPreregistrationEnrollmentdata != "" {
				req = req.Data(EnrollPreregistrationEnrollmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !EnrollPreregistrationEnrollmentQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !EnrollPreregistrationEnrollmentQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&EnrollPreregistrationEnrollmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&EnrollPreregistrationEnrollmentQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	EnrollPreregistrationEnrollmentCmd := NewEnrollPreregistrationEnrollmentCmd()
	WebAuthnPreregistrationCmd.AddCommand(EnrollPreregistrationEnrollmentCmd)
}

var (
	GenerateFulfillmentRequestdata string

	GenerateFulfillmentRequestQuiet bool
)

func NewGenerateFulfillmentRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateFulfillmentRequest",
		Long: "Generate a Fulfillment Request",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.GenerateFulfillmentRequest(apiClient.GetConfig().Context)

			if GenerateFulfillmentRequestdata != "" {
				req = req.Data(GenerateFulfillmentRequestdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GenerateFulfillmentRequestQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GenerateFulfillmentRequestQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GenerateFulfillmentRequestdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&GenerateFulfillmentRequestQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GenerateFulfillmentRequestCmd := NewGenerateFulfillmentRequestCmd()
	WebAuthnPreregistrationCmd.AddCommand(GenerateFulfillmentRequestCmd)
}

var (
	ListWebAuthnPreregistrationFactorsuserId string

	ListWebAuthnPreregistrationFactorsBackupDir string

	ListWebAuthnPreregistrationFactorsLimit    int32
	ListWebAuthnPreregistrationFactorsPage     string
	ListWebAuthnPreregistrationFactorsFetchAll bool

	ListWebAuthnPreregistrationFactorsQuiet bool
)

func NewListWebAuthnPreregistrationFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listFactors",
		Long: "List all WebAuthn Preregistration Factors",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.ListWebAuthnPreregistrationFactors(apiClient.GetConfig().Context, ListWebAuthnPreregistrationFactorsuserId)

			var allItems []map[string]interface{}
			var pageCount int

			for {
				resp, err := req.Execute()
				if err != nil {
					if resp != nil && resp.Body != nil {
						d, err := io.ReadAll(resp.Body)
						if err == nil && !ListWebAuthnPreregistrationFactorsQuiet {
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
					if !ListWebAuthnPreregistrationFactorsQuiet {
						utils.PrettyPrintByte(d)
					}
					return nil
				}

				allItems = append(allItems, items...)
				pageCount++

				if !ListWebAuthnPreregistrationFactorsFetchAll || len(items) == 0 {
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

			if ListWebAuthnPreregistrationFactorsFetchAll && pageCount > 1 && !ListWebAuthnPreregistrationFactorsQuiet {
				fmt.Printf("Retrieved %d items across %d pages\n", len(allItems), pageCount)
			}

			combinedJSON, err := json.Marshal(allItems)
			if err != nil {
				return fmt.Errorf("error combining results: %w", err)
			}

			if cmd.Flags().Changed("batch-backup") {
				dirPath := filepath.Join(ListWebAuthnPreregistrationFactorsBackupDir, "webauthnpreregistration", "listFactors")

				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				if !ListWebAuthnPreregistrationFactorsQuiet {
					fmt.Printf("Backing up WebAuthnPreregistrations to %s\n", dirPath)
				}

				success := 0
				for _, item := range allItems {
					id, ok := item["id"].(string)
					if !ok {
						if !ListWebAuthnPreregistrationFactorsQuiet {
							fmt.Println("Warning: item missing ID field, skipping")
						}
						continue
					}

					itemJSON, err := json.MarshalIndent(item, "", "  ")
					if err != nil {
						if !ListWebAuthnPreregistrationFactorsQuiet {
							fmt.Printf("Error marshaling item %s: %v\n", id, err)
						}
						continue
					}

					filePath := filepath.Join(dirPath, id+".json")
					if err := os.WriteFile(filePath, itemJSON, 0o644); err != nil {
						if !ListWebAuthnPreregistrationFactorsQuiet {
							fmt.Printf("Error writing file for %s: %v\n", id, err)
						}
						continue
					}

					success++
				}

				if !ListWebAuthnPreregistrationFactorsQuiet {
					fmt.Printf("Successfully backed up %d/%d WebAuthnPreregistrations\n", success, len(allItems))
				}
				return nil
			} else {
				if !ListWebAuthnPreregistrationFactorsQuiet {
					return utils.PrettyPrintByte(combinedJSON)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&ListWebAuthnPreregistrationFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().Int32VarP(&ListWebAuthnPreregistrationFactorsLimit, "limit", "l", 0, "Maximum number of items to return per page")
	cmd.Flags().StringVarP(&ListWebAuthnPreregistrationFactorsPage, "page", "p", "", "Page to fetch (if supported)")
	cmd.Flags().BoolVarP(&ListWebAuthnPreregistrationFactorsFetchAll, "all", "", false, "Fetch all items by following pagination automatically")
	cmd.Flags().BoolP("batch-backup", "b", false, "Backup multiple WebAuthnPreregistrations to a directory")

	cmd.Flags().StringVarP(&ListWebAuthnPreregistrationFactorsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&ListWebAuthnPreregistrationFactorsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ListWebAuthnPreregistrationFactorsCmd := NewListWebAuthnPreregistrationFactorsCmd()
	WebAuthnPreregistrationCmd.AddCommand(ListWebAuthnPreregistrationFactorsCmd)
}

var (
	DeleteWebAuthnPreregistrationFactoruserId string

	DeleteWebAuthnPreregistrationFactorauthenticatorEnrollmentId string

	DeleteWebAuthnPreregistrationFactorQuiet bool
)

func NewDeleteWebAuthnPreregistrationFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteFactor",
		Long: "Delete a WebAuthn Preregistration Factor",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.DeleteWebAuthnPreregistrationFactor(apiClient.GetConfig().Context, DeleteWebAuthnPreregistrationFactoruserId, DeleteWebAuthnPreregistrationFactorauthenticatorEnrollmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !DeleteWebAuthnPreregistrationFactorQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !DeleteWebAuthnPreregistrationFactorQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteWebAuthnPreregistrationFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&DeleteWebAuthnPreregistrationFactorauthenticatorEnrollmentId, "authenticatorEnrollmentId", "", "", "")
	cmd.MarkFlagRequired("authenticatorEnrollmentId")

	cmd.Flags().BoolVarP(&DeleteWebAuthnPreregistrationFactorQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	DeleteWebAuthnPreregistrationFactorCmd := NewDeleteWebAuthnPreregistrationFactorCmd()
	WebAuthnPreregistrationCmd.AddCommand(DeleteWebAuthnPreregistrationFactorCmd)
}
