package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var GroupCmd = &cobra.Command{
	Use:  "group",
	Long: "Manage GroupAPI",
}

func init() {
	rootCmd.AddCommand(GroupCmd)
}

var CreateGroupdata string

func NewCreateGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.CreateGroup(apiClient.GetConfig().Context)

			if CreateGroupdata != "" {
				req = req.Data(CreateGroupdata)
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

	cmd.Flags().StringVarP(&CreateGroupdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateGroupCmd := NewCreateGroupCmd()
	GroupCmd.AddCommand(CreateGroupCmd)
}

func NewListGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListGroups(apiClient.GetConfig().Context)

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
	ListGroupsCmd := NewListGroupsCmd()
	GroupCmd.AddCommand(ListGroupsCmd)
}

var CreateGroupRuledata string

func NewCreateGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.CreateGroupRule(apiClient.GetConfig().Context)

			if CreateGroupRuledata != "" {
				req = req.Data(CreateGroupRuledata)
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

	cmd.Flags().StringVarP(&CreateGroupRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateGroupRuleCmd := NewCreateGroupRuleCmd()
	GroupCmd.AddCommand(CreateGroupRuleCmd)
}

func NewListGroupRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listRules",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListGroupRules(apiClient.GetConfig().Context)

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
	ListGroupRulesCmd := NewListGroupRulesCmd()
	GroupCmd.AddCommand(ListGroupRulesCmd)
}

var GetGroupRulegroupRuleId string

func NewGetGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.GetGroupRule(apiClient.GetConfig().Context, GetGroupRulegroupRuleId)

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

	cmd.Flags().StringVarP(&GetGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	return cmd
}

func init() {
	GetGroupRuleCmd := NewGetGroupRuleCmd()
	GroupCmd.AddCommand(GetGroupRuleCmd)
}

var (
	ReplaceGroupRulegroupRuleId string

	ReplaceGroupRuledata string
)

func NewReplaceGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ReplaceGroupRule(apiClient.GetConfig().Context, ReplaceGroupRulegroupRuleId)

			if ReplaceGroupRuledata != "" {
				req = req.Data(ReplaceGroupRuledata)
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

	cmd.Flags().StringVarP(&ReplaceGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().StringVarP(&ReplaceGroupRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceGroupRuleCmd := NewReplaceGroupRuleCmd()
	GroupCmd.AddCommand(ReplaceGroupRuleCmd)
}

var DeleteGroupRulegroupRuleId string

func NewDeleteGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.DeleteGroupRule(apiClient.GetConfig().Context, DeleteGroupRulegroupRuleId)

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

	cmd.Flags().StringVarP(&DeleteGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	return cmd
}

func init() {
	DeleteGroupRuleCmd := NewDeleteGroupRuleCmd()
	GroupCmd.AddCommand(DeleteGroupRuleCmd)
}

var (
	ActivateGroupRulegroupRuleId string

	ActivateGroupRuledata string
)

func NewActivateGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activateRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ActivateGroupRule(apiClient.GetConfig().Context, ActivateGroupRulegroupRuleId)

			if ActivateGroupRuledata != "" {
				req = req.Data(ActivateGroupRuledata)
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

	cmd.Flags().StringVarP(&ActivateGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().StringVarP(&ActivateGroupRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivateGroupRuleCmd := NewActivateGroupRuleCmd()
	GroupCmd.AddCommand(ActivateGroupRuleCmd)
}

var (
	DeactivateGroupRulegroupRuleId string

	DeactivateGroupRuledata string
)

func NewDeactivateGroupRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivateRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.DeactivateGroupRule(apiClient.GetConfig().Context, DeactivateGroupRulegroupRuleId)

			if DeactivateGroupRuledata != "" {
				req = req.Data(DeactivateGroupRuledata)
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

	cmd.Flags().StringVarP(&DeactivateGroupRulegroupRuleId, "groupRuleId", "", "", "")
	cmd.MarkFlagRequired("groupRuleId")

	cmd.Flags().StringVarP(&DeactivateGroupRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	DeactivateGroupRuleCmd := NewDeactivateGroupRuleCmd()
	GroupCmd.AddCommand(DeactivateGroupRuleCmd)
}

var GetGroupgroupId string

func NewGetGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.GetGroup(apiClient.GetConfig().Context, GetGroupgroupId)

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

	cmd.Flags().StringVarP(&GetGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	GetGroupCmd := NewGetGroupCmd()
	GroupCmd.AddCommand(GetGroupCmd)
}

var (
	ReplaceGroupgroupId string

	ReplaceGroupdata string
)

func NewReplaceGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ReplaceGroup(apiClient.GetConfig().Context, ReplaceGroupgroupId)

			if ReplaceGroupdata != "" {
				req = req.Data(ReplaceGroupdata)
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

	cmd.Flags().StringVarP(&ReplaceGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&ReplaceGroupdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceGroupCmd := NewReplaceGroupCmd()
	GroupCmd.AddCommand(ReplaceGroupCmd)
}

var DeleteGroupgroupId string

func NewDeleteGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.DeleteGroup(apiClient.GetConfig().Context, DeleteGroupgroupId)

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

	cmd.Flags().StringVarP(&DeleteGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	DeleteGroupCmd := NewDeleteGroupCmd()
	GroupCmd.AddCommand(DeleteGroupCmd)
}

var ListAssignedApplicationsForGroupgroupId string

func NewListAssignedApplicationsForGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listAssignedApplicationsFor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListAssignedApplicationsForGroup(apiClient.GetConfig().Context, ListAssignedApplicationsForGroupgroupId)

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

	cmd.Flags().StringVarP(&ListAssignedApplicationsForGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	ListAssignedApplicationsForGroupCmd := NewListAssignedApplicationsForGroupCmd()
	GroupCmd.AddCommand(ListAssignedApplicationsForGroupCmd)
}

var ListGroupUsersgroupId string

func NewListGroupUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listUsers",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.ListGroupUsers(apiClient.GetConfig().Context, ListGroupUsersgroupId)

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

	cmd.Flags().StringVarP(&ListGroupUsersgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	ListGroupUsersCmd := NewListGroupUsersCmd()
	GroupCmd.AddCommand(ListGroupUsersCmd)
}

var (
	AssignUserToGroupgroupId string

	AssignUserToGroupuserId string

	AssignUserToGroupdata string
)

func NewAssignUserToGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "assignUserTo",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.AssignUserToGroup(apiClient.GetConfig().Context, AssignUserToGroupgroupId, AssignUserToGroupuserId)

			if AssignUserToGroupdata != "" {
				req = req.Data(AssignUserToGroupdata)
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

	cmd.Flags().StringVarP(&AssignUserToGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignUserToGroupuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&AssignUserToGroupdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AssignUserToGroupCmd := NewAssignUserToGroupCmd()
	GroupCmd.AddCommand(AssignUserToGroupCmd)
}

var (
	UnassignUserFromGroupgroupId string

	UnassignUserFromGroupuserId string
)

func NewUnassignUserFromGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "unassignUserFrom",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupAPI.UnassignUserFromGroup(apiClient.GetConfig().Context, UnassignUserFromGroupgroupId, UnassignUserFromGroupuserId)

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

	cmd.Flags().StringVarP(&UnassignUserFromGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&UnassignUserFromGroupuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	UnassignUserFromGroupCmd := NewUnassignUserFromGroupCmd()
	GroupCmd.AddCommand(UnassignUserFromGroupCmd)
}
