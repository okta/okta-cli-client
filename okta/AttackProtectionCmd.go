package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AttackProtectionCmd = &cobra.Command{
	Use:   "attackProtection",
	Long:  "Manage AttackProtectionAPI",
}

func NewAttackProtectionCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "attackProtection",
		Long:  "Manage AttackProtectionAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(AttackProtectionCmd)
}

var (
    
    
    
)

func NewGetAuthenticatorSettingsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getAuthenticatorSettings",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AttackProtectionAPI.GetAuthenticatorSettings(apiClient.GetConfig().Context)
            
            
            
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
	GetAuthenticatorSettingsCmd := NewGetAuthenticatorSettingsCmd()
    AttackProtectionCmd.AddCommand(GetAuthenticatorSettingsCmd)
}

var (
    
    
            ReplaceAuthenticatorSettingsdata string
        
    
)

func NewReplaceAuthenticatorSettingsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replaceAuthenticatorSettings",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AttackProtectionAPI.ReplaceAuthenticatorSettings(apiClient.GetConfig().Context)
            
            
            if ReplaceAuthenticatorSettingsdata != "" {
                req = req.Data(ReplaceAuthenticatorSettingsdata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceAuthenticatorSettingsdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceAuthenticatorSettingsCmd := NewReplaceAuthenticatorSettingsCmd()
    AttackProtectionCmd.AddCommand(ReplaceAuthenticatorSettingsCmd)
}

var (
    
    
    
)

func NewGetUserLockoutSettingsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getUserLockoutSettings",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AttackProtectionAPI.GetUserLockoutSettings(apiClient.GetConfig().Context)
            
            
            
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
	GetUserLockoutSettingsCmd := NewGetUserLockoutSettingsCmd()
    AttackProtectionCmd.AddCommand(GetUserLockoutSettingsCmd)
}

var (
    
    
            ReplaceUserLockoutSettingsdata string
        
    
)

func NewReplaceUserLockoutSettingsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replaceUserLockoutSettings",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AttackProtectionAPI.ReplaceUserLockoutSettings(apiClient.GetConfig().Context)
            
            
            if ReplaceUserLockoutSettingsdata != "" {
                req = req.Data(ReplaceUserLockoutSettingsdata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceUserLockoutSettingsdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceUserLockoutSettingsCmd := NewReplaceUserLockoutSettingsCmd()
    AttackProtectionCmd.AddCommand(ReplaceUserLockoutSettingsCmd)
}