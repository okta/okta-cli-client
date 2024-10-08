package okta

import (
	"io"

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

func NewGetWellknownOrgMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getWellknownOrgMetadata",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetWellknownOrgMetadata(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetWellknownOrgMetadataCmd := NewGetWellknownOrgMetadataCmd()
	OrgSettingCmd.AddCommand(GetWellknownOrgMetadataCmd)
}

var UpdateOrgSettingsdata string

func NewUpdateOrgSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updates",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateOrgSettings(apiClient.GetConfig().Context)

			if UpdateOrgSettingsdata != "" {
				req = req.Data(UpdateOrgSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateOrgSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateOrgSettingsCmd := NewUpdateOrgSettingsCmd()
	OrgSettingCmd.AddCommand(UpdateOrgSettingsCmd)
}

func NewGetOrgSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gets",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetOrgSettingsCmd := NewGetOrgSettingsCmd()
	OrgSettingCmd.AddCommand(GetOrgSettingsCmd)
}

var ReplaceOrgSettingsdata string

func NewReplaceOrgSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaces",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.ReplaceOrgSettings(apiClient.GetConfig().Context)

			if ReplaceOrgSettingsdata != "" {
				req = req.Data(ReplaceOrgSettingsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceOrgSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceOrgSettingsCmd := NewReplaceOrgSettingsCmd()
	OrgSettingCmd.AddCommand(ReplaceOrgSettingsCmd)
}

func NewGetOrgContactTypesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOrgContactTypes",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgContactTypes(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetOrgContactTypesCmd := NewGetOrgContactTypesCmd()
	OrgSettingCmd.AddCommand(GetOrgContactTypesCmd)
}

var GetOrgContactUsercontactType string

func NewGetOrgContactUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOrgContactUser",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgContactUser(apiClient.GetConfig().Context, GetOrgContactUsercontactType)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetOrgContactUsercontactType, "contactType", "", "", "")
	cmd.MarkFlagRequired("contactType")

	return cmd
}

func init() {
	GetOrgContactUserCmd := NewGetOrgContactUserCmd()
	OrgSettingCmd.AddCommand(GetOrgContactUserCmd)
}

var (
	ReplaceOrgContactUsercontactType string

	ReplaceOrgContactUserdata string
)

func NewReplaceOrgContactUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceOrgContactUser",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.ReplaceOrgContactUser(apiClient.GetConfig().Context, ReplaceOrgContactUsercontactType)

			if ReplaceOrgContactUserdata != "" {
				req = req.Data(ReplaceOrgContactUserdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceOrgContactUsercontactType, "contactType", "", "", "")
	cmd.MarkFlagRequired("contactType")

	cmd.Flags().StringVarP(&ReplaceOrgContactUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceOrgContactUserCmd := NewReplaceOrgContactUserCmd()
	OrgSettingCmd.AddCommand(ReplaceOrgContactUserCmd)
}

var BulkRemoveEmailAddressBouncesdata string

func NewBulkRemoveEmailAddressBouncesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bulkRemoveEmailAddressBounces",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.BulkRemoveEmailAddressBounces(apiClient.GetConfig().Context)

			if BulkRemoveEmailAddressBouncesdata != "" {
				req = req.Data(BulkRemoveEmailAddressBouncesdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&BulkRemoveEmailAddressBouncesdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	BulkRemoveEmailAddressBouncesCmd := NewBulkRemoveEmailAddressBouncesCmd()
	OrgSettingCmd.AddCommand(BulkRemoveEmailAddressBouncesCmd)
}

var UploadOrgLogodata string

func NewUploadOrgLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "uploadOrgLogo",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UploadOrgLogo(apiClient.GetConfig().Context)

			if UploadOrgLogodata != "" {
				req = req.Data(UploadOrgLogodata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadOrgLogodata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UploadOrgLogoCmd := NewUploadOrgLogoCmd()
	OrgSettingCmd.AddCommand(UploadOrgLogoCmd)
}

var UpdateThirdPartyAdminSettingdata string

func NewUpdateThirdPartyAdminSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateThirdPartyAdminSetting",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateThirdPartyAdminSetting(apiClient.GetConfig().Context)

			if UpdateThirdPartyAdminSettingdata != "" {
				req = req.Data(UpdateThirdPartyAdminSettingdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateThirdPartyAdminSettingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateThirdPartyAdminSettingCmd := NewUpdateThirdPartyAdminSettingCmd()
	OrgSettingCmd.AddCommand(UpdateThirdPartyAdminSettingCmd)
}

func NewGetThirdPartyAdminSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getThirdPartyAdminSetting",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetThirdPartyAdminSetting(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetThirdPartyAdminSettingCmd := NewGetThirdPartyAdminSettingCmd()
	OrgSettingCmd.AddCommand(GetThirdPartyAdminSettingCmd)
}

func NewGetOrgPreferencesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOrgPreferences",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgPreferences(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetOrgPreferencesCmd := NewGetOrgPreferencesCmd()
	OrgSettingCmd.AddCommand(GetOrgPreferencesCmd)
}

var UpdateOrgHideOktaUIFooterdata string

func NewUpdateOrgHideOktaUIFooterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateOrgHideOktaUIFooter",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateOrgHideOktaUIFooter(apiClient.GetConfig().Context)

			if UpdateOrgHideOktaUIFooterdata != "" {
				req = req.Data(UpdateOrgHideOktaUIFooterdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateOrgHideOktaUIFooterdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateOrgHideOktaUIFooterCmd := NewUpdateOrgHideOktaUIFooterCmd()
	OrgSettingCmd.AddCommand(UpdateOrgHideOktaUIFooterCmd)
}

var UpdateOrgShowOktaUIFooterdata string

func NewUpdateOrgShowOktaUIFooterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateOrgShowOktaUIFooter",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.UpdateOrgShowOktaUIFooter(apiClient.GetConfig().Context)

			if UpdateOrgShowOktaUIFooterdata != "" {
				req = req.Data(UpdateOrgShowOktaUIFooterdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UpdateOrgShowOktaUIFooterdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateOrgShowOktaUIFooterCmd := NewUpdateOrgShowOktaUIFooterCmd()
	OrgSettingCmd.AddCommand(UpdateOrgShowOktaUIFooterCmd)
}

func NewGetOktaCommunicationSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOktaCommunicationSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOktaCommunicationSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetOktaCommunicationSettingsCmd := NewGetOktaCommunicationSettingsCmd()
	OrgSettingCmd.AddCommand(GetOktaCommunicationSettingsCmd)
}

var OptInUsersToOktaCommunicationEmailsdata string

func NewOptInUsersToOktaCommunicationEmailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "optInUsersToOktaCommunicationEmails",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.OptInUsersToOktaCommunicationEmails(apiClient.GetConfig().Context)

			if OptInUsersToOktaCommunicationEmailsdata != "" {
				req = req.Data(OptInUsersToOktaCommunicationEmailsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&OptInUsersToOktaCommunicationEmailsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	OptInUsersToOktaCommunicationEmailsCmd := NewOptInUsersToOktaCommunicationEmailsCmd()
	OrgSettingCmd.AddCommand(OptInUsersToOktaCommunicationEmailsCmd)
}

var OptOutUsersFromOktaCommunicationEmailsdata string

func NewOptOutUsersFromOktaCommunicationEmailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "optOutUsersFromOktaCommunicationEmails",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.OptOutUsersFromOktaCommunicationEmails(apiClient.GetConfig().Context)

			if OptOutUsersFromOktaCommunicationEmailsdata != "" {
				req = req.Data(OptOutUsersFromOktaCommunicationEmailsdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&OptOutUsersFromOktaCommunicationEmailsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	OptOutUsersFromOktaCommunicationEmailsCmd := NewOptOutUsersFromOktaCommunicationEmailsCmd()
	OrgSettingCmd.AddCommand(OptOutUsersFromOktaCommunicationEmailsCmd)
}

func NewGetOrgOktaSupportSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOrgOktaSupportSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetOrgOktaSupportSettings(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetOrgOktaSupportSettingsCmd := NewGetOrgOktaSupportSettingsCmd()
	OrgSettingCmd.AddCommand(GetOrgOktaSupportSettingsCmd)
}

var ExtendOktaSupportdata string

func NewExtendOktaSupportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "extendOktaSupport",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.ExtendOktaSupport(apiClient.GetConfig().Context)

			if ExtendOktaSupportdata != "" {
				req = req.Data(ExtendOktaSupportdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ExtendOktaSupportdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ExtendOktaSupportCmd := NewExtendOktaSupportCmd()
	OrgSettingCmd.AddCommand(ExtendOktaSupportCmd)
}

var GrantOktaSupportdata string

func NewGrantOktaSupportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "grantOktaSupport",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GrantOktaSupport(apiClient.GetConfig().Context)

			if GrantOktaSupportdata != "" {
				req = req.Data(GrantOktaSupportdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GrantOktaSupportdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	GrantOktaSupportCmd := NewGrantOktaSupportCmd()
	OrgSettingCmd.AddCommand(GrantOktaSupportCmd)
}

var RevokeOktaSupportdata string

func NewRevokeOktaSupportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "revokeOktaSupport",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.RevokeOktaSupport(apiClient.GetConfig().Context)

			if RevokeOktaSupportdata != "" {
				req = req.Data(RevokeOktaSupportdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&RevokeOktaSupportdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	RevokeOktaSupportCmd := NewRevokeOktaSupportCmd()
	OrgSettingCmd.AddCommand(RevokeOktaSupportCmd)
}

func NewGetClientPrivilegesSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getClientPrivilegesSetting",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.GetClientPrivilegesSetting(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	GetClientPrivilegesSettingCmd := NewGetClientPrivilegesSettingCmd()
	OrgSettingCmd.AddCommand(GetClientPrivilegesSettingCmd)
}

var AssignClientPrivilegesSettingdata string

func NewAssignClientPrivilegesSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "assignClientPrivilegesSetting",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.OrgSettingAPI.AssignClientPrivilegesSetting(apiClient.GetConfig().Context)

			if AssignClientPrivilegesSettingdata != "" {
				req = req.Data(AssignClientPrivilegesSettingdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignClientPrivilegesSettingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AssignClientPrivilegesSettingCmd := NewAssignClientPrivilegesSettingCmd()
	OrgSettingCmd.AddCommand(AssignClientPrivilegesSettingCmd)
}
