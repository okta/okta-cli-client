package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var BehaviorCmd = &cobra.Command{
	Use:   "behavior",
	Long:  "Manage BehaviorAPI",
}

func NewBehaviorCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "behavior",
		Long:  "Manage BehaviorAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(BehaviorCmd)
}

var (
    
    
            CreateBehaviorDetectionRuledata string
        
    
)

func NewCreateBehaviorDetectionRuleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "createDetectionRule",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.CreateBehaviorDetectionRule(apiClient.GetConfig().Context)
            
            
            if CreateBehaviorDetectionRuledata != "" {
                req = req.Data(CreateBehaviorDetectionRuledata)
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

    
    
        cmd.Flags().StringVarP(&CreateBehaviorDetectionRuledata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	CreateBehaviorDetectionRuleCmd := NewCreateBehaviorDetectionRuleCmd()
    BehaviorCmd.AddCommand(CreateBehaviorDetectionRuleCmd)
}

var (
    
    
    
)

func NewListBehaviorDetectionRulesCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listDetectionRules",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.ListBehaviorDetectionRules(apiClient.GetConfig().Context)
            
            
            
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
	ListBehaviorDetectionRulesCmd := NewListBehaviorDetectionRulesCmd()
    BehaviorCmd.AddCommand(ListBehaviorDetectionRulesCmd)
}

var (
    
    
            GetBehaviorDetectionRulebehaviorId string
        
    
)

func NewGetBehaviorDetectionRuleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getDetectionRule",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.GetBehaviorDetectionRule(apiClient.GetConfig().Context, GetBehaviorDetectionRulebehaviorId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
        cmd.MarkFlagRequired("behaviorId")
        
    

	return cmd
}

func init() {
	GetBehaviorDetectionRuleCmd := NewGetBehaviorDetectionRuleCmd()
    BehaviorCmd.AddCommand(GetBehaviorDetectionRuleCmd)
}

var (
    
    
            ReplaceBehaviorDetectionRulebehaviorId string
        
            ReplaceBehaviorDetectionRuledata string
        
    
)

func NewReplaceBehaviorDetectionRuleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replaceDetectionRule",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.ReplaceBehaviorDetectionRule(apiClient.GetConfig().Context, ReplaceBehaviorDetectionRulebehaviorId)
            
            
            if ReplaceBehaviorDetectionRuledata != "" {
                req = req.Data(ReplaceBehaviorDetectionRuledata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
        cmd.MarkFlagRequired("behaviorId")
        
        cmd.Flags().StringVarP(&ReplaceBehaviorDetectionRuledata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceBehaviorDetectionRuleCmd := NewReplaceBehaviorDetectionRuleCmd()
    BehaviorCmd.AddCommand(ReplaceBehaviorDetectionRuleCmd)
}

var (
    
    
            DeleteBehaviorDetectionRulebehaviorId string
        
    
)

func NewDeleteBehaviorDetectionRuleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deleteDetectionRule",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.DeleteBehaviorDetectionRule(apiClient.GetConfig().Context, DeleteBehaviorDetectionRulebehaviorId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&DeleteBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
        cmd.MarkFlagRequired("behaviorId")
        
    

	return cmd
}

func init() {
	DeleteBehaviorDetectionRuleCmd := NewDeleteBehaviorDetectionRuleCmd()
    BehaviorCmd.AddCommand(DeleteBehaviorDetectionRuleCmd)
}

var (
    
    
            ActivateBehaviorDetectionRulebehaviorId string
        
            ActivateBehaviorDetectionRuledata string
        
    
)

func NewActivateBehaviorDetectionRuleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "activateDetectionRule",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.ActivateBehaviorDetectionRule(apiClient.GetConfig().Context, ActivateBehaviorDetectionRulebehaviorId)
            
            
            if ActivateBehaviorDetectionRuledata != "" {
                req = req.Data(ActivateBehaviorDetectionRuledata)
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

    
    
        cmd.Flags().StringVarP(&ActivateBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
        cmd.MarkFlagRequired("behaviorId")
        
        cmd.Flags().StringVarP(&ActivateBehaviorDetectionRuledata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ActivateBehaviorDetectionRuleCmd := NewActivateBehaviorDetectionRuleCmd()
    BehaviorCmd.AddCommand(ActivateBehaviorDetectionRuleCmd)
}

var (
    
    
            DeactivateBehaviorDetectionRulebehaviorId string
        
            DeactivateBehaviorDetectionRuledata string
        
    
)

func NewDeactivateBehaviorDetectionRuleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deactivateDetectionRule",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.BehaviorAPI.DeactivateBehaviorDetectionRule(apiClient.GetConfig().Context, DeactivateBehaviorDetectionRulebehaviorId)
            
            
            if DeactivateBehaviorDetectionRuledata != "" {
                req = req.Data(DeactivateBehaviorDetectionRuledata)
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

    
    
        cmd.Flags().StringVarP(&DeactivateBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
        cmd.MarkFlagRequired("behaviorId")
        
        cmd.Flags().StringVarP(&DeactivateBehaviorDetectionRuledata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	DeactivateBehaviorDetectionRuleCmd := NewDeactivateBehaviorDetectionRuleCmd()
    BehaviorCmd.AddCommand(DeactivateBehaviorDetectionRuleCmd)
}