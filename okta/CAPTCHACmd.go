package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var CAPTCHACmd = &cobra.Command{
	Use:  "cAPTCHA",
	Long: "Manage CAPTCHAAPI",
}

func init() {
	rootCmd.AddCommand(CAPTCHACmd)
}

var CreateCaptchaInstancedata string

func NewCreateCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createCaptchaInstance",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.CreateCaptchaInstance(apiClient.GetConfig().Context)

			if CreateCaptchaInstancedata != "" {
				req = req.Data(CreateCaptchaInstancedata)
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

	cmd.Flags().StringVarP(&CreateCaptchaInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateCaptchaInstanceCmd := NewCreateCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(CreateCaptchaInstanceCmd)
}

func NewListCaptchaInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listCaptchaInstances",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.ListCaptchaInstances(apiClient.GetConfig().Context)

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
	ListCaptchaInstancesCmd := NewListCaptchaInstancesCmd()
	CAPTCHACmd.AddCommand(ListCaptchaInstancesCmd)
}

var (
	UpdateCaptchaInstancecaptchaId string

	UpdateCaptchaInstancedata string
)

func NewUpdateCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateCaptchaInstance",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.UpdateCaptchaInstance(apiClient.GetConfig().Context, UpdateCaptchaInstancecaptchaId)

			if UpdateCaptchaInstancedata != "" {
				req = req.Data(UpdateCaptchaInstancedata)
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

	cmd.Flags().StringVarP(&UpdateCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	cmd.Flags().StringVarP(&UpdateCaptchaInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateCaptchaInstanceCmd := NewUpdateCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(UpdateCaptchaInstanceCmd)
}

var GetCaptchaInstancecaptchaId string

func NewGetCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getCaptchaInstance",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.GetCaptchaInstance(apiClient.GetConfig().Context, GetCaptchaInstancecaptchaId)

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

	cmd.Flags().StringVarP(&GetCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	return cmd
}

func init() {
	GetCaptchaInstanceCmd := NewGetCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(GetCaptchaInstanceCmd)
}

var (
	ReplaceCaptchaInstancecaptchaId string

	ReplaceCaptchaInstancedata string
)

func NewReplaceCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceCaptchaInstance",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.ReplaceCaptchaInstance(apiClient.GetConfig().Context, ReplaceCaptchaInstancecaptchaId)

			if ReplaceCaptchaInstancedata != "" {
				req = req.Data(ReplaceCaptchaInstancedata)
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

	cmd.Flags().StringVarP(&ReplaceCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	cmd.Flags().StringVarP(&ReplaceCaptchaInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceCaptchaInstanceCmd := NewReplaceCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(ReplaceCaptchaInstanceCmd)
}

var DeleteCaptchaInstancecaptchaId string

func NewDeleteCaptchaInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteCaptchaInstance",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.DeleteCaptchaInstance(apiClient.GetConfig().Context, DeleteCaptchaInstancecaptchaId)

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

	cmd.Flags().StringVarP(&DeleteCaptchaInstancecaptchaId, "captchaId", "", "", "")
	cmd.MarkFlagRequired("captchaId")

	return cmd
}

func init() {
	DeleteCaptchaInstanceCmd := NewDeleteCaptchaInstanceCmd()
	CAPTCHACmd.AddCommand(DeleteCaptchaInstanceCmd)
}

func NewGetOrgCaptchaSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getOrgCaptchaSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.GetOrgCaptchaSettings(apiClient.GetConfig().Context)

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
	GetOrgCaptchaSettingsCmd := NewGetOrgCaptchaSettingsCmd()
	CAPTCHACmd.AddCommand(GetOrgCaptchaSettingsCmd)
}

var ReplacesOrgCaptchaSettingsdata string

func NewReplacesOrgCaptchaSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replacesOrgCaptchaSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.ReplacesOrgCaptchaSettings(apiClient.GetConfig().Context)

			if ReplacesOrgCaptchaSettingsdata != "" {
				req = req.Data(ReplacesOrgCaptchaSettingsdata)
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

	cmd.Flags().StringVarP(&ReplacesOrgCaptchaSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacesOrgCaptchaSettingsCmd := NewReplacesOrgCaptchaSettingsCmd()
	CAPTCHACmd.AddCommand(ReplacesOrgCaptchaSettingsCmd)
}

func NewDeleteOrgCaptchaSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteOrgCaptchaSettings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CAPTCHAAPI.DeleteOrgCaptchaSettings(apiClient.GetConfig().Context)

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
	DeleteOrgCaptchaSettingsCmd := NewDeleteOrgCaptchaSettingsCmd()
	CAPTCHACmd.AddCommand(DeleteOrgCaptchaSettingsCmd)
}
