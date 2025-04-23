package okta

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var OrgSettingCmd = &cobra.Command{
	Use:  "orgSetting",
	Long: "Manage OrgSettingAPI",
}

func init() {
	rootCmd.AddCommand(OrgSettingCmd)
}

var (
	GetWellknownOrgMetadataBackupDir string

	GetWellknownOrgMetadataQuiet bool
)

func NewGetWellknownOrgMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getWellknownOrgMetadata",
		Long: "Retrieve the Well-Known Org Metadata",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetWellknownOrgMetadata(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetWellknownOrgMetadataQuiet {
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
				dirPath := filepath.Join(GetWellknownOrgMetadataBackupDir, "orgsetting", "getWellknownOrgMetadata")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetWellknownOrgMetadataQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetWellknownOrgMetadataQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetWellknownOrgMetadataBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetWellknownOrgMetadataQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetWellknownOrgMetadataCmd := NewGetWellknownOrgMetadataCmd()
	OrgSettingCmd.AddCommand(GetWellknownOrgMetadataCmd)
}

var (
	UpdateOrgSettingsdata string

	UpdateOrgSettingsQuiet bool
)

func NewUpdateOrgSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updates",
		Long: "Update the Org Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateOrgSettings(apiClient.GetConfig().Context)

			if UpdateOrgSettingsdata != "" {
				req = req.Data(UpdateOrgSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateOrgSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateOrgSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateOrgSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UpdateOrgSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateOrgSettingsCmd := NewUpdateOrgSettingsCmd()
	OrgSettingCmd.AddCommand(UpdateOrgSettingsCmd)
}

var (
	GetOrgSettingsBackupDir string

	GetOrgSettingsQuiet bool
)

func NewGetOrgSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "gets",
		Long: "Retrieve the Org Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOrgSettingsQuiet {
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
				dirPath := filepath.Join(GetOrgSettingsBackupDir, "orgsetting", "gets")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOrgSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOrgSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetOrgSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOrgSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOrgSettingsCmd := NewGetOrgSettingsCmd()
	OrgSettingCmd.AddCommand(GetOrgSettingsCmd)
}

var (
	ReplaceOrgSettingsdata string

	ReplaceOrgSettingsQuiet bool
)

func NewReplaceOrgSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaces",
		Long: "Replace the Org Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.ReplaceOrgSettings(apiClient.GetConfig().Context)

			if ReplaceOrgSettingsdata != "" {
				req = req.Data(ReplaceOrgSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceOrgSettingsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceOrgSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceOrgSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceOrgSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceOrgSettingsCmd := NewReplaceOrgSettingsCmd()
	OrgSettingCmd.AddCommand(ReplaceOrgSettingsCmd)
}

var (
	GetOrgContactTypesBackupDir string

	GetOrgContactTypesQuiet bool
)

func NewGetOrgContactTypesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOrgContactTypes",
		Long: "Retrieve the Org Contact Types",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgContactTypes(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOrgContactTypesQuiet {
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
				dirPath := filepath.Join(GetOrgContactTypesBackupDir, "orgsetting", "getOrgContactTypes")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOrgContactTypesQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOrgContactTypesQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetOrgContactTypesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOrgContactTypesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOrgContactTypesCmd := NewGetOrgContactTypesCmd()
	OrgSettingCmd.AddCommand(GetOrgContactTypesCmd)
}

var (
	GetOrgContactUsercontactType string

	GetOrgContactUserBackupDir string

	GetOrgContactUserQuiet bool
)

func NewGetOrgContactUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOrgContactUser",
		Long: "Retrieve the User of the Contact Type",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgContactUser(apiClient.GetConfig().Context, GetOrgContactUsercontactType)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOrgContactUserQuiet {
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
				dirPath := filepath.Join(GetOrgContactUserBackupDir, "orgsetting", "getOrgContactUser")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				idParam := GetOrgContactUsercontactType
				fileName := fmt.Sprintf("%s.json", idParam)

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOrgContactUserQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOrgContactUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetOrgContactUsercontactType, "contactType", "", "", "")
	cmd.MarkFlagRequired("contactType")

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetOrgContactUserBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOrgContactUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOrgContactUserCmd := NewGetOrgContactUserCmd()
	OrgSettingCmd.AddCommand(GetOrgContactUserCmd)
}

var (
	ReplaceOrgContactUsercontactType string

	ReplaceOrgContactUserdata string

	ReplaceOrgContactUserQuiet bool
)

func NewReplaceOrgContactUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceOrgContactUser",
		Long: "Replace the User of the Contact Type",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.ReplaceOrgContactUser(apiClient.GetConfig().Context, ReplaceOrgContactUsercontactType)

			if ReplaceOrgContactUserdata != "" {
				req = req.Data(ReplaceOrgContactUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ReplaceOrgContactUserQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ReplaceOrgContactUserQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceOrgContactUsercontactType, "contactType", "", "", "")
	cmd.MarkFlagRequired("contactType")

	cmd.Flags().StringVarP(&ReplaceOrgContactUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&ReplaceOrgContactUserQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ReplaceOrgContactUserCmd := NewReplaceOrgContactUserCmd()
	OrgSettingCmd.AddCommand(ReplaceOrgContactUserCmd)
}

var (
	BulkRemoveEmailAddressBouncesdata string

	BulkRemoveEmailAddressBouncesQuiet bool
)

func NewBulkRemoveEmailAddressBouncesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "bulkRemoveEmailAddressBounces",
		Long: "Remove Emails from Email Provider Bounce List",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.BulkRemoveEmailAddressBounces(apiClient.GetConfig().Context)

			if BulkRemoveEmailAddressBouncesdata != "" {
				req = req.Data(BulkRemoveEmailAddressBouncesdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !BulkRemoveEmailAddressBouncesQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !BulkRemoveEmailAddressBouncesQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&BulkRemoveEmailAddressBouncesdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&BulkRemoveEmailAddressBouncesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	BulkRemoveEmailAddressBouncesCmd := NewBulkRemoveEmailAddressBouncesCmd()
	OrgSettingCmd.AddCommand(BulkRemoveEmailAddressBouncesCmd)
}

var (
	UploadOrgLogodata string

	UploadOrgLogoQuiet bool
)

func NewUploadOrgLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadOrgLogo",
		Long: "Upload the Org Logo",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UploadOrgLogo(apiClient.GetConfig().Context)

			if UploadOrgLogodata != "" {
				req = req.Data(UploadOrgLogodata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UploadOrgLogoQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UploadOrgLogoQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadOrgLogodata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&UploadOrgLogoQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UploadOrgLogoCmd := NewUploadOrgLogoCmd()
	OrgSettingCmd.AddCommand(UploadOrgLogoCmd)
}

var UpdateThirdPartyAdminSettingQuiet bool

func NewUpdateThirdPartyAdminSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateThirdPartyAdminSetting",
		Long: "Update the Org Third-Party Admin setting",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateThirdPartyAdminSetting(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateThirdPartyAdminSettingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateThirdPartyAdminSettingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&UpdateThirdPartyAdminSettingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateThirdPartyAdminSettingCmd := NewUpdateThirdPartyAdminSettingCmd()
	OrgSettingCmd.AddCommand(UpdateThirdPartyAdminSettingCmd)
}

var (
	GetThirdPartyAdminSettingBackupDir string

	GetThirdPartyAdminSettingQuiet bool
)

func NewGetThirdPartyAdminSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getThirdPartyAdminSetting",
		Long: "Retrieve the Org Third-Party Admin setting",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetThirdPartyAdminSetting(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetThirdPartyAdminSettingQuiet {
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
				dirPath := filepath.Join(GetThirdPartyAdminSettingBackupDir, "orgsetting", "getThirdPartyAdminSetting")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetThirdPartyAdminSettingQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetThirdPartyAdminSettingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetThirdPartyAdminSettingBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetThirdPartyAdminSettingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetThirdPartyAdminSettingCmd := NewGetThirdPartyAdminSettingCmd()
	OrgSettingCmd.AddCommand(GetThirdPartyAdminSettingCmd)
}

var (
	GetOrgPreferencesBackupDir string

	GetOrgPreferencesQuiet bool
)

func NewGetOrgPreferencesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOrgPreferences",
		Long: "Retrieve the Org Preferences",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgPreferences(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOrgPreferencesQuiet {
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
				dirPath := filepath.Join(GetOrgPreferencesBackupDir, "orgsetting", "getOrgPreferences")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOrgPreferencesQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOrgPreferencesQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetOrgPreferencesBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOrgPreferencesQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOrgPreferencesCmd := NewGetOrgPreferencesCmd()
	OrgSettingCmd.AddCommand(GetOrgPreferencesCmd)
}

var UpdateOrgHideOktaUIFooterQuiet bool

func NewUpdateOrgHideOktaUIFooterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateOrgHideOktaUIFooter",
		Long: "Update the Preference to Hide the Okta Dashboard Footer",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateOrgHideOktaUIFooter(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateOrgHideOktaUIFooterQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateOrgHideOktaUIFooterQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&UpdateOrgHideOktaUIFooterQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateOrgHideOktaUIFooterCmd := NewUpdateOrgHideOktaUIFooterCmd()
	OrgSettingCmd.AddCommand(UpdateOrgHideOktaUIFooterCmd)
}

var UpdateOrgShowOktaUIFooterQuiet bool

func NewUpdateOrgShowOktaUIFooterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateOrgShowOktaUIFooter",
		Long: "Update the Preference to Show the Okta Dashboard Footer",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateOrgShowOktaUIFooter(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !UpdateOrgShowOktaUIFooterQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !UpdateOrgShowOktaUIFooterQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&UpdateOrgShowOktaUIFooterQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	UpdateOrgShowOktaUIFooterCmd := NewUpdateOrgShowOktaUIFooterCmd()
	OrgSettingCmd.AddCommand(UpdateOrgShowOktaUIFooterCmd)
}

var (
	GetOktaCommunicationSettingsBackupDir string

	GetOktaCommunicationSettingsQuiet bool
)

func NewGetOktaCommunicationSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOktaCommunicationSettings",
		Long: "Retrieve the Okta Communication Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOktaCommunicationSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOktaCommunicationSettingsQuiet {
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
				dirPath := filepath.Join(GetOktaCommunicationSettingsBackupDir, "orgsetting", "getOktaCommunicationSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOktaCommunicationSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOktaCommunicationSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetOktaCommunicationSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOktaCommunicationSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOktaCommunicationSettingsCmd := NewGetOktaCommunicationSettingsCmd()
	OrgSettingCmd.AddCommand(GetOktaCommunicationSettingsCmd)
}

var OptInUsersToOktaCommunicationEmailsQuiet bool

func NewOptInUsersToOktaCommunicationEmailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "optInUsersToOktaCommunicationEmails",
		Long: "Opt in all Users to Okta Communication emails",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.OptInUsersToOktaCommunicationEmails(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !OptInUsersToOktaCommunicationEmailsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !OptInUsersToOktaCommunicationEmailsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&OptInUsersToOktaCommunicationEmailsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	OptInUsersToOktaCommunicationEmailsCmd := NewOptInUsersToOktaCommunicationEmailsCmd()
	OrgSettingCmd.AddCommand(OptInUsersToOktaCommunicationEmailsCmd)
}

var OptOutUsersFromOktaCommunicationEmailsQuiet bool

func NewOptOutUsersFromOktaCommunicationEmailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "optOutUsersFromOktaCommunicationEmails",
		Long: "Opt out all Users from Okta Communication emails",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.OptOutUsersFromOktaCommunicationEmails(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !OptOutUsersFromOktaCommunicationEmailsQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !OptOutUsersFromOktaCommunicationEmailsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&OptOutUsersFromOktaCommunicationEmailsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	OptOutUsersFromOktaCommunicationEmailsCmd := NewOptOutUsersFromOktaCommunicationEmailsCmd()
	OrgSettingCmd.AddCommand(OptOutUsersFromOktaCommunicationEmailsCmd)
}

var (
	GetOrgOktaSupportSettingsBackupDir string

	GetOrgOktaSupportSettingsQuiet bool
)

func NewGetOrgOktaSupportSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getOrgOktaSupportSettings",
		Long: "Retrieve the Okta Support Settings",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgOktaSupportSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetOrgOktaSupportSettingsQuiet {
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
				dirPath := filepath.Join(GetOrgOktaSupportSettingsBackupDir, "orgsetting", "getOrgOktaSupportSettings")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetOrgOktaSupportSettingsQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetOrgOktaSupportSettingsQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetOrgOktaSupportSettingsBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetOrgOktaSupportSettingsQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetOrgOktaSupportSettingsCmd := NewGetOrgOktaSupportSettingsCmd()
	OrgSettingCmd.AddCommand(GetOrgOktaSupportSettingsCmd)
}

var ExtendOktaSupportQuiet bool

func NewExtendOktaSupportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "extendOktaSupport",
		Long: "Extend Okta Support Access",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.ExtendOktaSupport(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !ExtendOktaSupportQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !ExtendOktaSupportQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&ExtendOktaSupportQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	ExtendOktaSupportCmd := NewExtendOktaSupportCmd()
	OrgSettingCmd.AddCommand(ExtendOktaSupportCmd)
}

var GrantOktaSupportQuiet bool

func NewGrantOktaSupportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "grantOktaSupport",
		Long: "Grant Okta Support Access to your Org",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GrantOktaSupport(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GrantOktaSupportQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !GrantOktaSupportQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&GrantOktaSupportQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GrantOktaSupportCmd := NewGrantOktaSupportCmd()
	OrgSettingCmd.AddCommand(GrantOktaSupportCmd)
}

var RevokeOktaSupportQuiet bool

func NewRevokeOktaSupportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeOktaSupport",
		Long: "Revoke Okta Support Access",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.RevokeOktaSupport(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !RevokeOktaSupportQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !RevokeOktaSupportQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&RevokeOktaSupportQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	RevokeOktaSupportCmd := NewRevokeOktaSupportCmd()
	OrgSettingCmd.AddCommand(RevokeOktaSupportCmd)
}

var (
	GetClientPrivilegesSettingBackupDir string

	GetClientPrivilegesSettingQuiet bool
)

func NewGetClientPrivilegesSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getClientPrivilegesSetting",
		Long: "Retrieve the Org settings to assign the Super Admin role",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isBackupRequested := cmd.Flags().Changed("backup") || cmd.Flags().Changed("batch-backup")
			if isBackupRequested && !cmd.Flags().Changed("backup-dir") {
				return fmt.Errorf("--backup-dir is required when using backup functionality")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetClientPrivilegesSetting(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !GetClientPrivilegesSettingQuiet {
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
				dirPath := filepath.Join(GetClientPrivilegesSettingBackupDir, "orgsetting", "getClientPrivilegesSetting")
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					return fmt.Errorf("failed to create backup directory: %w", err)
				}

				fileName := "orgsetting.json"

				filePath := filepath.Join(dirPath, fileName)

				if err := os.WriteFile(filePath, d, 0o644); err != nil {
					return fmt.Errorf("failed to write backup file: %w", err)
				}

				if !GetClientPrivilegesSettingQuiet {
					fmt.Printf("Backup completed successfully to %s\n", filePath)
				}
				return nil
			}

			if !GetClientPrivilegesSettingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("backup", "b", false, "Backup the OrgSetting to a file")

	cmd.Flags().StringVarP(&GetClientPrivilegesSettingBackupDir, "backup-dir", "d", "", "Directory to save backups")

	cmd.Flags().BoolVarP(&GetClientPrivilegesSettingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	GetClientPrivilegesSettingCmd := NewGetClientPrivilegesSettingCmd()
	OrgSettingCmd.AddCommand(GetClientPrivilegesSettingCmd)
}

var (
	AssignClientPrivilegesSettingdata string

	AssignClientPrivilegesSettingQuiet bool
)

func NewAssignClientPrivilegesSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignClientPrivilegesSetting",
		Long: "Assign the Super Admin role to a public client app",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.AssignClientPrivilegesSetting(apiClient.GetConfig().Context)

			if AssignClientPrivilegesSettingdata != "" {
				req = req.Data(AssignClientPrivilegesSettingdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil && !AssignClientPrivilegesSettingQuiet {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if !AssignClientPrivilegesSettingQuiet {
				utils.PrettyPrintByte(d)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignClientPrivilegesSettingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	cmd.Flags().BoolVarP(&AssignClientPrivilegesSettingQuiet, "quiet", "q", false, "Suppress normal output")

	return cmd
}

func init() {
	AssignClientPrivilegesSettingCmd := NewAssignClientPrivilegesSettingCmd()
	OrgSettingCmd.AddCommand(AssignClientPrivilegesSettingCmd)
}
