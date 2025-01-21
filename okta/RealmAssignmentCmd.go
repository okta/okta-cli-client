package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var RealmAssignmentCmd = &cobra.Command{
	Use:  "realmAssignment",
	Long: "Manage RealmAssignmentAPI",
}

func init() {
	rootCmd.AddCommand(RealmAssignmentCmd)
}

var CreateRealmAssignmentdata string

func NewCreateRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.CreateRealmAssignment(apiClient.GetConfig().Context)

			if CreateRealmAssignmentdata != "" {
				req = req.Data(CreateRealmAssignmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateRealmAssignmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateRealmAssignmentCmd := NewCreateRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(CreateRealmAssignmentCmd)
}

func NewListRealmAssignmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Realm Assignments",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ListRealmAssignments(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	ListRealmAssignmentsCmd := NewListRealmAssignmentsCmd()
	RealmAssignmentCmd.AddCommand(ListRealmAssignmentsCmd)
}

var ExecuteRealmAssignmentdata string

func NewExecuteRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "execute",
		Long: "Execute a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ExecuteRealmAssignment(apiClient.GetConfig().Context)

			if ExecuteRealmAssignmentdata != "" {
				req = req.Data(ExecuteRealmAssignmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ExecuteRealmAssignmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ExecuteRealmAssignmentCmd := NewExecuteRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(ExecuteRealmAssignmentCmd)
}

func NewListRealmAssignmentOperationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listOperations",
		Long: "List all Realm Assignment operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ListRealmAssignmentOperations(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	ListRealmAssignmentOperationsCmd := NewListRealmAssignmentOperationsCmd()
	RealmAssignmentCmd.AddCommand(ListRealmAssignmentOperationsCmd)
}

var GetRealmAssignmentassignmentId string

func NewGetRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.GetRealmAssignment(apiClient.GetConfig().Context, GetRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	return cmd
}

func init() {
	GetRealmAssignmentCmd := NewGetRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(GetRealmAssignmentCmd)
}

var (
	ReplaceRealmAssignmentassignmentId string

	ReplaceRealmAssignmentdata string
)

func NewReplaceRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ReplaceRealmAssignment(apiClient.GetConfig().Context, ReplaceRealmAssignmentassignmentId)

			if ReplaceRealmAssignmentdata != "" {
				req = req.Data(ReplaceRealmAssignmentdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	cmd.Flags().StringVarP(&ReplaceRealmAssignmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceRealmAssignmentCmd := NewReplaceRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(ReplaceRealmAssignmentCmd)
}

var DeleteRealmAssignmentassignmentId string

func NewDeleteRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.DeleteRealmAssignment(apiClient.GetConfig().Context, DeleteRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	return cmd
}

func init() {
	DeleteRealmAssignmentCmd := NewDeleteRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(DeleteRealmAssignmentCmd)
}

var ActivateRealmAssignmentassignmentId string

func NewActivateRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.ActivateRealmAssignment(apiClient.GetConfig().Context, ActivateRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	return cmd
}

func init() {
	ActivateRealmAssignmentCmd := NewActivateRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(ActivateRealmAssignmentCmd)
}

var DeactivateRealmAssignmentassignmentId string

func NewDeactivateRealmAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Realm Assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.RealmAssignmentAPI.DeactivateRealmAssignment(apiClient.GetConfig().Context, DeactivateRealmAssignmentassignmentId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateRealmAssignmentassignmentId, "assignmentId", "", "", "")
	cmd.MarkFlagRequired("assignmentId")

	return cmd
}

func init() {
	DeactivateRealmAssignmentCmd := NewDeactivateRealmAssignmentCmd()
	RealmAssignmentCmd.AddCommand(DeactivateRealmAssignmentCmd)
}
