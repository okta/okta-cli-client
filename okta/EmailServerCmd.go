package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var EmailServerCmd = &cobra.Command{
	Use:  "emailServer",
	Long: "Manage EmailServerAPI",
}

func init() {
	rootCmd.AddCommand(EmailServerCmd)
}

var CreateEmailServerdata string

func NewCreateEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a custom SMTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.CreateEmailServer(apiClient.GetConfig().Context)

			if CreateEmailServerdata != "" {
				req = req.Data(CreateEmailServerdata)
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

	cmd.Flags().StringVarP(&CreateEmailServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateEmailServerCmd := NewCreateEmailServerCmd()
	EmailServerCmd.AddCommand(CreateEmailServerCmd)
}

func NewListEmailServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all enrolled SMTP servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.ListEmailServers(apiClient.GetConfig().Context)

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
	ListEmailServersCmd := NewListEmailServersCmd()
	EmailServerCmd.AddCommand(ListEmailServersCmd)
}

var GetEmailServeremailServerId string

func NewGetEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve an SMTP Server configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.GetEmailServer(apiClient.GetConfig().Context, GetEmailServeremailServerId)

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

	cmd.Flags().StringVarP(&GetEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	return cmd
}

func init() {
	GetEmailServerCmd := NewGetEmailServerCmd()
	EmailServerCmd.AddCommand(GetEmailServerCmd)
}

var DeleteEmailServeremailServerId string

func NewDeleteEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete an SMTP Server configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.DeleteEmailServer(apiClient.GetConfig().Context, DeleteEmailServeremailServerId)

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

	cmd.Flags().StringVarP(&DeleteEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	return cmd
}

func init() {
	DeleteEmailServerCmd := NewDeleteEmailServerCmd()
	EmailServerCmd.AddCommand(DeleteEmailServerCmd)
}

var (
	UpdateEmailServeremailServerId string

	UpdateEmailServerdata string
)

func NewUpdateEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: "Update an SMTP Server configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.UpdateEmailServer(apiClient.GetConfig().Context, UpdateEmailServeremailServerId)

			if UpdateEmailServerdata != "" {
				req = req.Data(UpdateEmailServerdata)
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

	cmd.Flags().StringVarP(&UpdateEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	cmd.Flags().StringVarP(&UpdateEmailServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateEmailServerCmd := NewUpdateEmailServerCmd()
	EmailServerCmd.AddCommand(UpdateEmailServerCmd)
}

var (
	TestEmailServeremailServerId string

	TestEmailServerdata string
)

func NewTestEmailServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "test",
		Long: "Test an SMTP Server configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailServerAPI.TestEmailServer(apiClient.GetConfig().Context, TestEmailServeremailServerId)

			if TestEmailServerdata != "" {
				req = req.Data(TestEmailServerdata)
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

	cmd.Flags().StringVarP(&TestEmailServeremailServerId, "emailServerId", "", "", "")
	cmd.MarkFlagRequired("emailServerId")

	cmd.Flags().StringVarP(&TestEmailServerdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	TestEmailServerCmd := NewTestEmailServerCmd()
	EmailServerCmd.AddCommand(TestEmailServerCmd)
}
