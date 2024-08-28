package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var RoleAssignmentCmd = &cobra.Command{
	Use:   "roleAssignment",
	Long:  "Manage RoleAssignmentAPI",
}

func NewRoleAssignmentCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "roleAssignment",
		Long:  "Manage RoleAssignmentAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(RoleAssignmentCmd)
}

var (
    
    
            AssignRoleToGroupgroupId string
        
            AssignRoleToGroupdata string
        
    
)

func NewAssignRoleToGroupCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "assignRoleToGroup",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.AssignRoleToGroup(apiClient.GetConfig().Context, AssignRoleToGroupgroupId)
            
            
            if AssignRoleToGroupdata != "" {
                req = req.Data(AssignRoleToGroupdata)
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

    
    
        cmd.Flags().StringVarP(&AssignRoleToGroupgroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
        cmd.Flags().StringVarP(&AssignRoleToGroupdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	AssignRoleToGroupCmd := NewAssignRoleToGroupCmd()
    RoleAssignmentCmd.AddCommand(AssignRoleToGroupCmd)
}

var (
    
    
            ListGroupAssignedRolesgroupId string
        
    
)

func NewListGroupAssignedRolesCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listGroupAssignedRoles",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.ListGroupAssignedRoles(apiClient.GetConfig().Context, ListGroupAssignedRolesgroupId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&ListGroupAssignedRolesgroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
    

	return cmd
}

func init() {
	ListGroupAssignedRolesCmd := NewListGroupAssignedRolesCmd()
    RoleAssignmentCmd.AddCommand(ListGroupAssignedRolesCmd)
}

var (
    
    
            GetGroupAssignedRolegroupId string
        
            GetGroupAssignedRoleroleId string
        
    
)

func NewGetGroupAssignedRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getGroupAssignedRole",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.GetGroupAssignedRole(apiClient.GetConfig().Context, GetGroupAssignedRolegroupId, GetGroupAssignedRoleroleId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetGroupAssignedRolegroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
        cmd.Flags().StringVarP(&GetGroupAssignedRoleroleId, "roleId", "", "", "")
        cmd.MarkFlagRequired("roleId")
        
    

	return cmd
}

func init() {
	GetGroupAssignedRoleCmd := NewGetGroupAssignedRoleCmd()
    RoleAssignmentCmd.AddCommand(GetGroupAssignedRoleCmd)
}

var (
    
    
            UnassignRoleFromGroupgroupId string
        
            UnassignRoleFromGrouproleId string
        
    
)

func NewUnassignRoleFromGroupCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "unassignRoleFromGroup",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.UnassignRoleFromGroup(apiClient.GetConfig().Context, UnassignRoleFromGroupgroupId, UnassignRoleFromGrouproleId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&UnassignRoleFromGroupgroupId, "groupId", "", "", "")
        cmd.MarkFlagRequired("groupId")
        
        cmd.Flags().StringVarP(&UnassignRoleFromGrouproleId, "roleId", "", "", "")
        cmd.MarkFlagRequired("roleId")
        
    

	return cmd
}

func init() {
	UnassignRoleFromGroupCmd := NewUnassignRoleFromGroupCmd()
    RoleAssignmentCmd.AddCommand(UnassignRoleFromGroupCmd)
}

var (
    
    
    
)

func NewListUsersWithRoleAssignmentsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listUsersWiths",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.ListUsersWithRoleAssignments(apiClient.GetConfig().Context)
            
            
            
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
	ListUsersWithRoleAssignmentsCmd := NewListUsersWithRoleAssignmentsCmd()
    RoleAssignmentCmd.AddCommand(ListUsersWithRoleAssignmentsCmd)
}

var (
    
    
            AssignRoleToUseruserId string
        
            AssignRoleToUserdata string
        
    
)

func NewAssignRoleToUserCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "assignRoleToUser",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.AssignRoleToUser(apiClient.GetConfig().Context, AssignRoleToUseruserId)
            
            
            if AssignRoleToUserdata != "" {
                req = req.Data(AssignRoleToUserdata)
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

    
    
        cmd.Flags().StringVarP(&AssignRoleToUseruserId, "userId", "", "", "")
        cmd.MarkFlagRequired("userId")
        
        cmd.Flags().StringVarP(&AssignRoleToUserdata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	AssignRoleToUserCmd := NewAssignRoleToUserCmd()
    RoleAssignmentCmd.AddCommand(AssignRoleToUserCmd)
}

var (
    
    
            ListAssignedRolesForUseruserId string
        
    
)

func NewListAssignedRolesForUserCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listAssignedRolesForUser",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.ListAssignedRolesForUser(apiClient.GetConfig().Context, ListAssignedRolesForUseruserId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&ListAssignedRolesForUseruserId, "userId", "", "", "")
        cmd.MarkFlagRequired("userId")
        
    

	return cmd
}

func init() {
	ListAssignedRolesForUserCmd := NewListAssignedRolesForUserCmd()
    RoleAssignmentCmd.AddCommand(ListAssignedRolesForUserCmd)
}

var (
    
    
            GetUserAssignedRoleuserId string
        
            GetUserAssignedRoleroleId string
        
    
)

func NewGetUserAssignedRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getUserAssignedRole",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.GetUserAssignedRole(apiClient.GetConfig().Context, GetUserAssignedRoleuserId, GetUserAssignedRoleroleId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetUserAssignedRoleuserId, "userId", "", "", "")
        cmd.MarkFlagRequired("userId")
        
        cmd.Flags().StringVarP(&GetUserAssignedRoleroleId, "roleId", "", "", "")
        cmd.MarkFlagRequired("roleId")
        
    

	return cmd
}

func init() {
	GetUserAssignedRoleCmd := NewGetUserAssignedRoleCmd()
    RoleAssignmentCmd.AddCommand(GetUserAssignedRoleCmd)
}

var (
    
    
            UnassignRoleFromUseruserId string
        
            UnassignRoleFromUserroleId string
        
    
)

func NewUnassignRoleFromUserCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "unassignRoleFromUser",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAssignmentAPI.UnassignRoleFromUser(apiClient.GetConfig().Context, UnassignRoleFromUseruserId, UnassignRoleFromUserroleId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&UnassignRoleFromUseruserId, "userId", "", "", "")
        cmd.MarkFlagRequired("userId")
        
        cmd.Flags().StringVarP(&UnassignRoleFromUserroleId, "roleId", "", "", "")
        cmd.MarkFlagRequired("roleId")
        
    

	return cmd
}

func init() {
	UnassignRoleFromUserCmd := NewUnassignRoleFromUserCmd()
    RoleAssignmentCmd.AddCommand(UnassignRoleFromUserCmd)
}