package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationGroupsCmd = &cobra.Command{
	Use:   "applicationGroups",
	Long:  "Manage ApplicationGroupsAPI",
}

func NewApplicationGroupsCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "applicationGroups",
		Long:  "Manage ApplicationGroupsAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(ApplicationGroupsCmd)
}

var (
    
    
            ListApplicationGroupAssignmentsappId string
        
    
)

func NewListApplicationGroupAssignmentsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listApplicationGroupAssignments",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.ApplicationGroupsAPI.ListApplicationGroupAssignments(apiClient.GetConfig().Context, ListApplicationGroupAssignmentsappId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&ListApplicationGroupAssignmentsappId, "appId", "", "", "")
        cmd.MarkFlagRequired("appId")
        
    

	return cmd
}

func init() {
	ListApplicationGroupAssignmentsCmd := NewListApplicationGroupAssignmentsCmd()
    ApplicationGroupsCmd.AddCommand(ListApplicationGroupAssignmentsCmd)
}

var (
    
    
            GetApplicationGroupAssignmentappId string
        
            GetApplicationGroupAssignmentgroupId string
        
    
)

func NewGetApplicationGroupAssignmentCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getApplicationGroupAssignment",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.ApplicationGroupsAPI.GetApplicationGroupAssignment(apiClient.GetConfig().Context, GetApplicationGroupAssignmentappId, GetApplicationGroupAssignmentgroupId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetApplicationGroupAssignmentappId, "appId", "", "", "")
        cmd.MarkFlagRequired("appId")
        
        cmd.Flags().StringVarP(&GetApplicationGroupAssignmentgroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
    

	return cmd
}

func init() {
	GetApplicationGroupAssignmentCmd := NewGetApplicationGroupAssignmentCmd()
    ApplicationGroupsCmd.AddCommand(GetApplicationGroupAssignmentCmd)
}

var (
    
    
            AssignGroupToApplicationappId string
        
            AssignGroupToApplicationgroupId string
        
            AssignGroupToApplicationdata string
        
    
)

func NewAssignGroupToApplicationCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "assignGroupToApplication",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.ApplicationGroupsAPI.AssignGroupToApplication(apiClient.GetConfig().Context, AssignGroupToApplicationappId, AssignGroupToApplicationgroupId)
            
            
            if AssignGroupToApplicationdata != "" {
                req = req.Data(AssignGroupToApplicationdata)
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

    
    
        cmd.Flags().StringVarP(&AssignGroupToApplicationappId, "appId", "", "", "")
        cmd.MarkFlagRequired("appId")
        
        cmd.Flags().StringVarP(&AssignGroupToApplicationgroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
        cmd.Flags().StringVarP(&AssignGroupToApplicationdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	AssignGroupToApplicationCmd := NewAssignGroupToApplicationCmd()
    ApplicationGroupsCmd.AddCommand(AssignGroupToApplicationCmd)
}

var (
    
    
            UnassignApplicationFromGroupappId string
        
            UnassignApplicationFromGroupgroupId string
        
    
)

func NewUnassignApplicationFromGroupCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "unassignApplicationFromGroup",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.ApplicationGroupsAPI.UnassignApplicationFromGroup(apiClient.GetConfig().Context, UnassignApplicationFromGroupappId, UnassignApplicationFromGroupgroupId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&UnassignApplicationFromGroupappId, "appId", "", "", "")
        cmd.MarkFlagRequired("appId")
        
        cmd.Flags().StringVarP(&UnassignApplicationFromGroupgroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
    

	return cmd
}

func init() {
	UnassignApplicationFromGroupCmd := NewUnassignApplicationFromGroupCmd()
    ApplicationGroupsCmd.AddCommand(UnassignApplicationFromGroupCmd)
}