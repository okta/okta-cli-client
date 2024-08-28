package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var GroupOwnerCmd = &cobra.Command{
	Use:  "groupOwner",
	Long: "Manage GroupOwnerAPI",
}

func init() {
	rootCmd.AddCommand(GroupOwnerCmd)
}

var (
	AssignGroupOwnergroupId string

	AssignGroupOwnerdata string
)

func NewAssignGroupOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "assign",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupOwnerAPI.AssignGroupOwner(apiClient.GetConfig().Context, AssignGroupOwnergroupId)

			if AssignGroupOwnerdata != "" {
				req = req.Data(AssignGroupOwnerdata)
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

	cmd.Flags().StringVarP(&AssignGroupOwnergroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignGroupOwnerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AssignGroupOwnerCmd := NewAssignGroupOwnerCmd()
	GroupOwnerCmd.AddCommand(AssignGroupOwnerCmd)
}

var ListGroupOwnersgroupId string

func NewListGroupOwnersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupOwnerAPI.ListGroupOwners(apiClient.GetConfig().Context, ListGroupOwnersgroupId)

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

	cmd.Flags().StringVarP(&ListGroupOwnersgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	ListGroupOwnersCmd := NewListGroupOwnersCmd()
	GroupOwnerCmd.AddCommand(ListGroupOwnersCmd)
}

var (
	DeleteGroupOwnergroupId string

	DeleteGroupOwnerownerId string
)

func NewDeleteGroupOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.GroupOwnerAPI.DeleteGroupOwner(apiClient.GetConfig().Context, DeleteGroupOwnergroupId, DeleteGroupOwnerownerId)

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

	cmd.Flags().StringVarP(&DeleteGroupOwnergroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&DeleteGroupOwnerownerId, "ownerId", "", "", "")
	cmd.MarkFlagRequired("ownerId")

	return cmd
}

func init() {
	DeleteGroupOwnerCmd := NewDeleteGroupOwnerCmd()
	GroupOwnerCmd.AddCommand(DeleteGroupOwnerCmd)
}
