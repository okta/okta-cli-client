package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var TemplateCmd = &cobra.Command{
	Use:   "template",
	Long:  "Manage TemplateAPI",
}

func NewTemplateCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "template",
		Long:  "Manage TemplateAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(TemplateCmd)
}

var (
    
    
            CreateSmsTemplatedata string
        
    
)

func NewCreateSmsTemplateCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "createSms",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.TemplateAPI.CreateSmsTemplate(apiClient.GetConfig().Context)
            
            
            if CreateSmsTemplatedata != "" {
                req = req.Data(CreateSmsTemplatedata)
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
            cmd.Println(string(d))
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

var (
    
    
    
)

func NewListSmsTemplatesCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listSmss",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.TemplateAPI.ListSmsTemplates(apiClient.GetConfig().Context)
            
            
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
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
	    Use:   "updateSms",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.TemplateAPI.UpdateSmsTemplate(apiClient.GetConfig().Context, UpdateSmsTemplatetemplateId)
            
            
            if UpdateSmsTemplatedata != "" {
                req = req.Data(UpdateSmsTemplatedata)
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
            cmd.Println(string(d))
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

var (
    
    
            GetSmsTemplatetemplateId string
        
    
)

func NewGetSmsTemplateCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getSms",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.TemplateAPI.GetSmsTemplate(apiClient.GetConfig().Context, GetSmsTemplatetemplateId)
            
            
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
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
	    Use:   "replaceSms",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.TemplateAPI.ReplaceSmsTemplate(apiClient.GetConfig().Context, ReplaceSmsTemplatetemplateId)
            
            
            if ReplaceSmsTemplatedata != "" {
                req = req.Data(ReplaceSmsTemplatedata)
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
            cmd.Println(string(d))
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

var (
    
    
            DeleteSmsTemplatetemplateId string
        
    
)

func NewDeleteSmsTemplateCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deleteSms",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.TemplateAPI.DeleteSmsTemplate(apiClient.GetConfig().Context, DeleteSmsTemplatetemplateId)
            
            
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
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