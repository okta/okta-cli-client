package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ThreatInsightCmd = &cobra.Command{
	Use:   "threatInsight",
	Long:  "Manage ThreatInsightAPI",
}

func NewThreatInsightCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "threatInsight",
		Long:  "Manage ThreatInsightAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(ThreatInsightCmd)
}

var (
    
    
            UpdateConfigurationdata string
        
    
)

func NewUpdateConfigurationCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "updateConfiguration",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.ThreatInsightAPI.UpdateConfiguration(apiClient.GetConfig().Context)
            
            
            if UpdateConfigurationdata != "" {
                req = req.Data(UpdateConfigurationdata)
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

    
    
        cmd.Flags().StringVarP(&UpdateConfigurationdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	UpdateConfigurationCmd := NewUpdateConfigurationCmd()
    ThreatInsightCmd.AddCommand(UpdateConfigurationCmd)
}

var (
    
    
    
)

func NewGetCurrentConfigurationCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getCurrentConfiguration",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.ThreatInsightAPI.GetCurrentConfiguration(apiClient.GetConfig().Context)
            
            
            
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
	GetCurrentConfigurationCmd := NewGetCurrentConfigurationCmd()
    ThreatInsightCmd.AddCommand(GetCurrentConfigurationCmd)
}