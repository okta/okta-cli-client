package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationUsersCmd = &cobra.Command{
	Use:  "applicationUsers",
	Long: "Manage ApplicationUsersAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationUsersCmd)
}

var (
	AssignUserToApplicationappId string

	AssignUserToApplicationdata string
)

func NewAssignUserToApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignUserToApplication",
		Long: "Assign an Application User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.AssignUserToApplication(apiClient.GetConfig().Context, AssignUserToApplicationappId)

			if AssignUserToApplicationdata != "" {
				req = req.Data(AssignUserToApplicationdata)
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

	cmd.Flags().StringVarP(&AssignUserToApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&AssignUserToApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AssignUserToApplicationCmd := NewAssignUserToApplicationCmd()
	ApplicationUsersCmd.AddCommand(AssignUserToApplicationCmd)
}

var ListApplicationUsersappId string

func NewListApplicationUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List all Application Users",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.ListApplicationUsers(apiClient.GetConfig().Context, ListApplicationUsersappId)

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

	cmd.Flags().StringVarP(&ListApplicationUsersappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	ListApplicationUsersCmd := NewListApplicationUsersCmd()
	ApplicationUsersCmd.AddCommand(ListApplicationUsersCmd)
}

var (
	UpdateApplicationUserappId string

	UpdateApplicationUseruserId string

	UpdateApplicationUserdata string
)

func NewUpdateApplicationUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "updateApplicationUser",
		Long: "Update an Application User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.UpdateApplicationUser(apiClient.GetConfig().Context, UpdateApplicationUserappId, UpdateApplicationUseruserId)

			if UpdateApplicationUserdata != "" {
				req = req.Data(UpdateApplicationUserdata)
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

	cmd.Flags().StringVarP(&UpdateApplicationUserappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateApplicationUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UpdateApplicationUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateApplicationUserCmd := NewUpdateApplicationUserCmd()
	ApplicationUsersCmd.AddCommand(UpdateApplicationUserCmd)
}

var (
	GetApplicationUserappId string

	GetApplicationUseruserId string
)

func NewGetApplicationUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationUser",
		Long: "Retrieve an Application User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.GetApplicationUser(apiClient.GetConfig().Context, GetApplicationUserappId, GetApplicationUseruserId)

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

	cmd.Flags().StringVarP(&GetApplicationUserappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetApplicationUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	GetApplicationUserCmd := NewGetApplicationUserCmd()
	ApplicationUsersCmd.AddCommand(GetApplicationUserCmd)
}

var (
	UnassignUserFromApplicationappId string

	UnassignUserFromApplicationuserId string
)

func NewUnassignUserFromApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignUserFromApplication",
		Long: "Unassign an Application User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationUsersAPI.UnassignUserFromApplication(apiClient.GetConfig().Context, UnassignUserFromApplicationappId, UnassignUserFromApplicationuserId)

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

	cmd.Flags().StringVarP(&UnassignUserFromApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UnassignUserFromApplicationuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	UnassignUserFromApplicationCmd := NewUnassignUserFromApplicationCmd()
	ApplicationUsersCmd.AddCommand(UnassignUserFromApplicationCmd)
}
