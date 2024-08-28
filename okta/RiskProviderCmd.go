package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var RiskProviderCmd = &cobra.Command{
	Use:   "riskProvider",
	Long:  "Manage RiskProviderAPI",
}

func NewRiskProviderCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "riskProvider",
		Long:  "Manage RiskProviderAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(RiskProviderCmd)
}

var (
    
    
            CreateRiskProviderdata string
        
    
)

func NewCreateRiskProviderCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "create",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RiskProviderAPI.CreateRiskProvider(apiClient.GetConfig().Context)
            
            
            if CreateRiskProviderdata != "" {
                req = req.Data(CreateRiskProviderdata)
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

    
    
        cmd.Flags().StringVarP(&CreateRiskProviderdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	CreateRiskProviderCmd := NewCreateRiskProviderCmd()
    RiskProviderCmd.AddCommand(CreateRiskProviderCmd)
}

var (
    
    
    
)

func NewListRiskProvidersCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "lists",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RiskProviderAPI.ListRiskProviders(apiClient.GetConfig().Context)
            
            
            
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
	ListRiskProvidersCmd := NewListRiskProvidersCmd()
    RiskProviderCmd.AddCommand(ListRiskProvidersCmd)
}

var (
    
    
            GetRiskProviderriskProviderId string
        
    
)

func NewGetRiskProviderCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "get",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RiskProviderAPI.GetRiskProvider(apiClient.GetConfig().Context, GetRiskProviderriskProviderId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetRiskProviderriskProviderId, "riskProviderId", "", "", "")
        cmd.MarkFlagRequired("riskProviderId")
        
    

	return cmd
}

func init() {
	GetRiskProviderCmd := NewGetRiskProviderCmd()
    RiskProviderCmd.AddCommand(GetRiskProviderCmd)
}

var (
    
    
            ReplaceRiskProviderriskProviderId string
        
            ReplaceRiskProviderdata string
        
    
)

func NewReplaceRiskProviderCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replace",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RiskProviderAPI.ReplaceRiskProvider(apiClient.GetConfig().Context, ReplaceRiskProviderriskProviderId)
            
            
            if ReplaceRiskProviderdata != "" {
                req = req.Data(ReplaceRiskProviderdata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceRiskProviderriskProviderId, "riskProviderId", "", "", "")
        cmd.MarkFlagRequired("riskProviderId")
        
        cmd.Flags().StringVarP(&ReplaceRiskProviderdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceRiskProviderCmd := NewReplaceRiskProviderCmd()
    RiskProviderCmd.AddCommand(ReplaceRiskProviderCmd)
}

var (
    
    
            DeleteRiskProviderriskProviderId string
        
    
)

func NewDeleteRiskProviderCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "delete",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RiskProviderAPI.DeleteRiskProvider(apiClient.GetConfig().Context, DeleteRiskProviderriskProviderId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&DeleteRiskProviderriskProviderId, "riskProviderId", "", "", "")
        cmd.MarkFlagRequired("riskProviderId")
        
    

	return cmd
}

func init() {
	DeleteRiskProviderCmd := NewDeleteRiskProviderCmd()
    RiskProviderCmd.AddCommand(DeleteRiskProviderCmd)
}