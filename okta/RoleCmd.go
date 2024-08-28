package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var RoleCmd = &cobra.Command{
	Use:   "role",
	Long:  "Manage RoleAPI",
}

func NewRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "role",
		Long:  "Manage RoleAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(RoleCmd)
}

var (
    
    
            CreateRoledata string
        
    
)

func NewCreateRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "create",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.CreateRole(apiClient.GetConfig().Context)
            
            
            if CreateRoledata != "" {
                req = req.Data(CreateRoledata)
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

    
    
        cmd.Flags().StringVarP(&CreateRoledata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	CreateRoleCmd := NewCreateRoleCmd()
    RoleCmd.AddCommand(CreateRoleCmd)
}

var (
    
    
    
)

func NewListRolesCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "lists",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.ListRoles(apiClient.GetConfig().Context)
            
            
            
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
	ListRolesCmd := NewListRolesCmd()
    RoleCmd.AddCommand(ListRolesCmd)
}

var (
    
    
            GetRoleroleIdOrLabel string
        
    
)

func NewGetRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "get",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.GetRole(apiClient.GetConfig().Context, GetRoleroleIdOrLabel)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetRoleroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
    

	return cmd
}

func init() {
	GetRoleCmd := NewGetRoleCmd()
    RoleCmd.AddCommand(GetRoleCmd)
}

var (
    
    
            ReplaceRoleroleIdOrLabel string
        
            ReplaceRoledata string
        
    
)

func NewReplaceRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replace",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.ReplaceRole(apiClient.GetConfig().Context, ReplaceRoleroleIdOrLabel)
            
            
            if ReplaceRoledata != "" {
                req = req.Data(ReplaceRoledata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceRoleroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
        cmd.Flags().StringVarP(&ReplaceRoledata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceRoleCmd := NewReplaceRoleCmd()
    RoleCmd.AddCommand(ReplaceRoleCmd)
}

var (
    
    
            DeleteRoleroleIdOrLabel string
        
    
)

func NewDeleteRoleCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "delete",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.DeleteRole(apiClient.GetConfig().Context, DeleteRoleroleIdOrLabel)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&DeleteRoleroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
    

	return cmd
}

func init() {
	DeleteRoleCmd := NewDeleteRoleCmd()
    RoleCmd.AddCommand(DeleteRoleCmd)
}

var (
    
    
            ListRolePermissionsroleIdOrLabel string
        
    
)

func NewListRolePermissionsCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listPermissions",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.ListRolePermissions(apiClient.GetConfig().Context, ListRolePermissionsroleIdOrLabel)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&ListRolePermissionsroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
    

	return cmd
}

func init() {
	ListRolePermissionsCmd := NewListRolePermissionsCmd()
    RoleCmd.AddCommand(ListRolePermissionsCmd)
}

var (
    
    
            CreateRolePermissionroleIdOrLabel string
        
            CreateRolePermissionpermissionType string
        
            CreateRolePermissiondata string
        
    
)

func NewCreateRolePermissionCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "createPermission",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.CreateRolePermission(apiClient.GetConfig().Context, CreateRolePermissionroleIdOrLabel, CreateRolePermissionpermissionType)
            
            
            if CreateRolePermissiondata != "" {
                req = req.Data(CreateRolePermissiondata)
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

    
    
        cmd.Flags().StringVarP(&CreateRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
        cmd.Flags().StringVarP(&CreateRolePermissionpermissionType, "permissionType", "", "", "")
        cmd.MarkFlagRequired("permissionType")
        
        cmd.Flags().StringVarP(&CreateRolePermissiondata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	CreateRolePermissionCmd := NewCreateRolePermissionCmd()
    RoleCmd.AddCommand(CreateRolePermissionCmd)
}

var (
    
    
            GetRolePermissionroleIdOrLabel string
        
            GetRolePermissionpermissionType string
        
    
)

func NewGetRolePermissionCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "getPermission",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.GetRolePermission(apiClient.GetConfig().Context, GetRolePermissionroleIdOrLabel, GetRolePermissionpermissionType)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
        cmd.Flags().StringVarP(&GetRolePermissionpermissionType, "permissionType", "", "", "")
        cmd.MarkFlagRequired("permissionType")
        
    

	return cmd
}

func init() {
	GetRolePermissionCmd := NewGetRolePermissionCmd()
    RoleCmd.AddCommand(GetRolePermissionCmd)
}

var (
    
    
            ReplaceRolePermissionroleIdOrLabel string
        
            ReplaceRolePermissionpermissionType string
        
            ReplaceRolePermissiondata string
        
    
)

func NewReplaceRolePermissionCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replacePermission",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.ReplaceRolePermission(apiClient.GetConfig().Context, ReplaceRolePermissionroleIdOrLabel, ReplaceRolePermissionpermissionType)
            
            
            if ReplaceRolePermissiondata != "" {
                req = req.Data(ReplaceRolePermissiondata)
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

    
    
        cmd.Flags().StringVarP(&ReplaceRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
        cmd.Flags().StringVarP(&ReplaceRolePermissionpermissionType, "permissionType", "", "", "")
        cmd.MarkFlagRequired("permissionType")
        
        cmd.Flags().StringVarP(&ReplaceRolePermissiondata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceRolePermissionCmd := NewReplaceRolePermissionCmd()
    RoleCmd.AddCommand(ReplaceRolePermissionCmd)
}

var (
    
    
            DeleteRolePermissionroleIdOrLabel string
        
            DeleteRolePermissionpermissionType string
        
    
)

func NewDeleteRolePermissionCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deletePermission",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.RoleAPI.DeleteRolePermission(apiClient.GetConfig().Context, DeleteRolePermissionroleIdOrLabel, DeleteRolePermissionpermissionType)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&DeleteRolePermissionroleIdOrLabel, "roleIdOrLabel", "", "", "")
        cmd.MarkFlagRequired("roleIdOrLabel")
        
        cmd.Flags().StringVarP(&DeleteRolePermissionpermissionType, "permissionType", "", "", "")
        cmd.MarkFlagRequired("permissionType")
        
    

	return cmd
}

func init() {
	DeleteRolePermissionCmd := NewDeleteRolePermissionCmd()
    RoleCmd.AddCommand(DeleteRolePermissionCmd)
}