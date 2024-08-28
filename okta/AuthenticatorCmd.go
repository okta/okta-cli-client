package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var AuthenticatorCmd = &cobra.Command{
	Use:   "authenticator",
	Long:  "Manage AuthenticatorAPI",
}

func NewAuthenticatorCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "authenticator",
		Long:  "Manage AuthenticatorAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(AuthenticatorCmd)
}

var (
    
    
    
)

func NewGetWellKnownAppAuthenticatorConfigurationCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getWellKnownAppConfiguration",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.GetWellKnownAppAuthenticatorConfiguration(apiClient.GetConfig().Context)
            
            
            
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
	GetWellKnownAppAuthenticatorConfigurationCmd := NewGetWellKnownAppAuthenticatorConfigurationCmd()
    AuthenticatorCmd.AddCommand(GetWellKnownAppAuthenticatorConfigurationCmd)
}

var (
    
    
            CreateAuthenticatordata string
        
    
)

func NewCreateAuthenticatorCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "create",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.CreateAuthenticator(apiClient.GetConfig().Context)
            
            
            if CreateAuthenticatordata != "" {
                req = req.Data(CreateAuthenticatordata)
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

    
    
        cmd.Flags().StringVarP(&CreateAuthenticatordata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	CreateAuthenticatorCmd := NewCreateAuthenticatorCmd()
    AuthenticatorCmd.AddCommand(CreateAuthenticatorCmd)
}

var (
    
    
    
)

func NewListAuthenticatorsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "lists",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.ListAuthenticators(apiClient.GetConfig().Context)
            
            
            
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
	ListAuthenticatorsCmd := NewListAuthenticatorsCmd()
    AuthenticatorCmd.AddCommand(ListAuthenticatorsCmd)
}

var (
    
    
            GetAuthenticatorauthenticatorId string
        
    
)

func NewGetAuthenticatorCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "get",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.GetAuthenticator(apiClient.GetConfig().Context, GetAuthenticatorauthenticatorId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
    

	return cmd
}

func init() {
	GetAuthenticatorCmd := NewGetAuthenticatorCmd()
    AuthenticatorCmd.AddCommand(GetAuthenticatorCmd)
}

var (
    
    
            ReplaceAuthenticatorauthenticatorId string
        
            ReplaceAuthenticatordata string
        
    
)

func NewReplaceAuthenticatorCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replace",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.ReplaceAuthenticator(apiClient.GetConfig().Context, ReplaceAuthenticatorauthenticatorId)
            
            
            if ReplaceAuthenticatordata != "" {
                req = req.Data(ReplaceAuthenticatordata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&ReplaceAuthenticatordata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceAuthenticatorCmd := NewReplaceAuthenticatorCmd()
    AuthenticatorCmd.AddCommand(ReplaceAuthenticatorCmd)
}

var (
    
    
            ActivateAuthenticatorauthenticatorId string
        
            ActivateAuthenticatordata string
        
    
)

func NewActivateAuthenticatorCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "activate",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.ActivateAuthenticator(apiClient.GetConfig().Context, ActivateAuthenticatorauthenticatorId)
            
            
            if ActivateAuthenticatordata != "" {
                req = req.Data(ActivateAuthenticatordata)
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

    
    
        cmd.Flags().StringVarP(&ActivateAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&ActivateAuthenticatordata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ActivateAuthenticatorCmd := NewActivateAuthenticatorCmd()
    AuthenticatorCmd.AddCommand(ActivateAuthenticatorCmd)
}

var (
    
    
            DeactivateAuthenticatorauthenticatorId string
        
            DeactivateAuthenticatordata string
        
    
)

func NewDeactivateAuthenticatorCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deactivate",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.DeactivateAuthenticator(apiClient.GetConfig().Context, DeactivateAuthenticatorauthenticatorId)
            
            
            if DeactivateAuthenticatordata != "" {
                req = req.Data(DeactivateAuthenticatordata)
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

    
    
        cmd.Flags().StringVarP(&DeactivateAuthenticatorauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&DeactivateAuthenticatordata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	DeactivateAuthenticatorCmd := NewDeactivateAuthenticatorCmd()
    AuthenticatorCmd.AddCommand(DeactivateAuthenticatorCmd)
}

var (
    
    
            ListAuthenticatorMethodsauthenticatorId string
        
    
)

func NewListAuthenticatorMethodsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listMethods",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.ListAuthenticatorMethods(apiClient.GetConfig().Context, ListAuthenticatorMethodsauthenticatorId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&ListAuthenticatorMethodsauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
    

	return cmd
}

func init() {
	ListAuthenticatorMethodsCmd := NewListAuthenticatorMethodsCmd()
    AuthenticatorCmd.AddCommand(ListAuthenticatorMethodsCmd)
}

var (
    
    
            GetAuthenticatorMethodauthenticatorId string
        
            GetAuthenticatorMethodmethodType string
        
    
)

func NewGetAuthenticatorMethodCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getMethod",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.GetAuthenticatorMethod(apiClient.GetConfig().Context, GetAuthenticatorMethodauthenticatorId, GetAuthenticatorMethodmethodType)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&GetAuthenticatorMethodmethodType, "methodType", "", "", "")
        cmd.MarkFlagRequired("methodType")
        
    

	return cmd
}

func init() {
	GetAuthenticatorMethodCmd := NewGetAuthenticatorMethodCmd()
    AuthenticatorCmd.AddCommand(GetAuthenticatorMethodCmd)
}

var (
    
    
            ReplaceAuthenticatorMethodauthenticatorId string
        
            ReplaceAuthenticatorMethodmethodType string
        
            ReplaceAuthenticatorMethoddata string
        
    
)

func NewReplaceAuthenticatorMethodCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replaceMethod",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.ReplaceAuthenticatorMethod(apiClient.GetConfig().Context, ReplaceAuthenticatorMethodauthenticatorId, ReplaceAuthenticatorMethodmethodType)
            
            
            if ReplaceAuthenticatorMethoddata != "" {
                req = req.Data(ReplaceAuthenticatorMethoddata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&ReplaceAuthenticatorMethodmethodType, "methodType", "", "", "")
        cmd.MarkFlagRequired("methodType")
        
        cmd.Flags().StringVarP(&ReplaceAuthenticatorMethoddata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceAuthenticatorMethodCmd := NewReplaceAuthenticatorMethodCmd()
    AuthenticatorCmd.AddCommand(ReplaceAuthenticatorMethodCmd)
}

var (
    
    
            ActivateAuthenticatorMethodauthenticatorId string
        
            ActivateAuthenticatorMethodmethodType string
        
            ActivateAuthenticatorMethoddata string
        
    
)

func NewActivateAuthenticatorMethodCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "activateMethod",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.ActivateAuthenticatorMethod(apiClient.GetConfig().Context, ActivateAuthenticatorMethodauthenticatorId, ActivateAuthenticatorMethodmethodType)
            
            
            if ActivateAuthenticatorMethoddata != "" {
                req = req.Data(ActivateAuthenticatorMethoddata)
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

    
    
        cmd.Flags().StringVarP(&ActivateAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&ActivateAuthenticatorMethodmethodType, "methodType", "", "", "")
        cmd.MarkFlagRequired("methodType")
        
        cmd.Flags().StringVarP(&ActivateAuthenticatorMethoddata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ActivateAuthenticatorMethodCmd := NewActivateAuthenticatorMethodCmd()
    AuthenticatorCmd.AddCommand(ActivateAuthenticatorMethodCmd)
}

var (
    
    
            DeactivateAuthenticatorMethodauthenticatorId string
        
            DeactivateAuthenticatorMethodmethodType string
        
            DeactivateAuthenticatorMethoddata string
        
    
)

func NewDeactivateAuthenticatorMethodCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deactivateMethod",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.AuthenticatorAPI.DeactivateAuthenticatorMethod(apiClient.GetConfig().Context, DeactivateAuthenticatorMethodauthenticatorId, DeactivateAuthenticatorMethodmethodType)
            
            
            if DeactivateAuthenticatorMethoddata != "" {
                req = req.Data(DeactivateAuthenticatorMethoddata)
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

    
    
        cmd.Flags().StringVarP(&DeactivateAuthenticatorMethodauthenticatorId, "authenticatorId", "", "", "")
        cmd.MarkFlagRequired("authenticatorId")
        
        cmd.Flags().StringVarP(&DeactivateAuthenticatorMethodmethodType, "methodType", "", "", "")
        cmd.MarkFlagRequired("methodType")
        
        cmd.Flags().StringVarP(&DeactivateAuthenticatorMethoddata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	DeactivateAuthenticatorMethodCmd := NewDeactivateAuthenticatorMethodCmd()
    AuthenticatorCmd.AddCommand(DeactivateAuthenticatorMethodCmd)
}