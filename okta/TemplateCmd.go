package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var TemplateCmd = &cobra.Command{
	Use:  "template",
	Long: "Manage TemplateAPI",
}

func init() {
	rootCmd.AddCommand(TemplateCmd)
}

var CreateSmsTemplatedata string

func NewCreateSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createSms",
		Long: "Create an SMS Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.CreateSmsTemplate(apiClient.GetConfig().Context)

			if CreateSmsTemplatedata != "" {
				req = req.Data(CreateSmsTemplatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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

	cmd.Flags().StringVarP(&CreateSmsTemplatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateSmsTemplateCmd := NewCreateSmsTemplateCmd()
	TemplateCmd.AddCommand(CreateSmsTemplateCmd)
}

func NewListSmsTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listSmss",
		Long: "List all SMS Templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.ListSmsTemplates(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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
	ListSmsTemplatesCmd := NewListSmsTemplatesCmd()
	TemplateCmd.AddCommand(ListSmsTemplatesCmd)
}

var (
	UpdateSmsTemplatetemplateId string

	UpdateSmsTemplatedata string
)

func NewUpdateSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateSms",
		Long: "Update an SMS Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.UpdateSmsTemplate(apiClient.GetConfig().Context, UpdateSmsTemplatetemplateId)

			if UpdateSmsTemplatedata != "" {
				req = req.Data(UpdateSmsTemplatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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

	cmd.Flags().StringVarP(&UpdateSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	cmd.Flags().StringVarP(&UpdateSmsTemplatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateSmsTemplateCmd := NewUpdateSmsTemplateCmd()
	TemplateCmd.AddCommand(UpdateSmsTemplateCmd)
}

var GetSmsTemplatetemplateId string

func NewGetSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSms",
		Long: "Retrieve an SMS Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.GetSmsTemplate(apiClient.GetConfig().Context, GetSmsTemplatetemplateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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

	cmd.Flags().StringVarP(&GetSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	return cmd
}

func init() {
	GetSmsTemplateCmd := NewGetSmsTemplateCmd()
	TemplateCmd.AddCommand(GetSmsTemplateCmd)
}

var (
	ReplaceSmsTemplatetemplateId string

	ReplaceSmsTemplatedata string
)

func NewReplaceSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceSms",
		Long: "Replace an SMS Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.ReplaceSmsTemplate(apiClient.GetConfig().Context, ReplaceSmsTemplatetemplateId)

			if ReplaceSmsTemplatedata != "" {
				req = req.Data(ReplaceSmsTemplatedata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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

	cmd.Flags().StringVarP(&ReplaceSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	cmd.Flags().StringVarP(&ReplaceSmsTemplatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceSmsTemplateCmd := NewReplaceSmsTemplateCmd()
	TemplateCmd.AddCommand(ReplaceSmsTemplateCmd)
}

var DeleteSmsTemplatetemplateId string

func NewDeleteSmsTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteSms",
		Long: "Delete an SMS Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.TemplateAPI.DeleteSmsTemplate(apiClient.GetConfig().Context, DeleteSmsTemplatetemplateId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
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

	cmd.Flags().StringVarP(&DeleteSmsTemplatetemplateId, "templateId", "", "", "")
	cmd.MarkFlagRequired("templateId")

	return cmd
}

func init() {
	DeleteSmsTemplateCmd := NewDeleteSmsTemplateCmd()
	TemplateCmd.AddCommand(DeleteSmsTemplateCmd)
}
