package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationCmd = &cobra.Command{
	Use:  "application",
	Long: "Manage ApplicationAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationCmd)
}

var CreateApplicationdata string

func NewCreateApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.CreateApplication(apiClient.GetConfig().Context)

			if CreateApplicationdata != "" {
				req = req.Data(CreateApplicationdata)
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

	cmd.Flags().StringVarP(&CreateApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateApplicationCmd := NewCreateApplicationCmd()
	ApplicationCmd.AddCommand(CreateApplicationCmd)
}

func NewListApplicationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.ListApplications(apiClient.GetConfig().Context)

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
	ListApplicationsCmd := NewListApplicationsCmd()
	ApplicationCmd.AddCommand(ListApplicationsCmd)
}

var GetApplicationappId string

func NewGetApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.GetApplication(apiClient.GetConfig().Context, GetApplicationappId)

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

	cmd.Flags().StringVarP(&GetApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	GetApplicationCmd := NewGetApplicationCmd()
	ApplicationCmd.AddCommand(GetApplicationCmd)
}

var (
	ReplaceApplicationappId string

	ReplaceApplicationdata string
)

func NewReplaceApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.ReplaceApplication(apiClient.GetConfig().Context, ReplaceApplicationappId)

			if ReplaceApplicationdata != "" {
				req = req.Data(ReplaceApplicationdata)
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

	cmd.Flags().StringVarP(&ReplaceApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&ReplaceApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceApplicationCmd := NewReplaceApplicationCmd()
	ApplicationCmd.AddCommand(ReplaceApplicationCmd)
}

var DeleteApplicationappId string

func NewDeleteApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.DeleteApplication(apiClient.GetConfig().Context, DeleteApplicationappId)

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

	cmd.Flags().StringVarP(&DeleteApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	DeleteApplicationCmd := NewDeleteApplicationCmd()
	ApplicationCmd.AddCommand(DeleteApplicationCmd)
}

var ActivateApplicationappId string

func NewActivateApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.ActivateApplication(apiClient.GetConfig().Context, ActivateApplicationappId)

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

	cmd.Flags().StringVarP(&ActivateApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	ActivateApplicationCmd := NewActivateApplicationCmd()
	ApplicationCmd.AddCommand(ActivateApplicationCmd)
}

var DeactivateApplicationappId string

func NewDeactivateApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationAPI.DeactivateApplication(apiClient.GetConfig().Context, DeactivateApplicationappId)

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

	cmd.Flags().StringVarP(&DeactivateApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	DeactivateApplicationCmd := NewDeactivateApplicationCmd()
	ApplicationCmd.AddCommand(DeactivateApplicationCmd)
}
