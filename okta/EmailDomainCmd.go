package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var EmailDomainCmd = &cobra.Command{
	Use:  "emailDomain",
	Long: "Manage EmailDomainAPI",
}

func init() {
	rootCmd.AddCommand(EmailDomainCmd)
}

var CreateEmailDomaindata string

func NewCreateEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.CreateEmailDomain(apiClient.GetConfig().Context)

			if CreateEmailDomaindata != "" {
				req = req.Data(CreateEmailDomaindata)
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

	cmd.Flags().StringVarP(&CreateEmailDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateEmailDomainCmd := NewCreateEmailDomainCmd()
	EmailDomainCmd.AddCommand(CreateEmailDomainCmd)
}

func NewListEmailDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.ListEmailDomains(apiClient.GetConfig().Context)

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
	ListEmailDomainsCmd := NewListEmailDomainsCmd()
	EmailDomainCmd.AddCommand(ListEmailDomainsCmd)
}

var GetEmailDomainemailDomainId string

func NewGetEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.GetEmailDomain(apiClient.GetConfig().Context, GetEmailDomainemailDomainId)

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

	cmd.Flags().StringVarP(&GetEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	return cmd
}

func init() {
	GetEmailDomainCmd := NewGetEmailDomainCmd()
	EmailDomainCmd.AddCommand(GetEmailDomainCmd)
}

var (
	ReplaceEmailDomainemailDomainId string

	ReplaceEmailDomaindata string
)

func NewReplaceEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.ReplaceEmailDomain(apiClient.GetConfig().Context, ReplaceEmailDomainemailDomainId)

			if ReplaceEmailDomaindata != "" {
				req = req.Data(ReplaceEmailDomaindata)
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

	cmd.Flags().StringVarP(&ReplaceEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	cmd.Flags().StringVarP(&ReplaceEmailDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceEmailDomainCmd := NewReplaceEmailDomainCmd()
	EmailDomainCmd.AddCommand(ReplaceEmailDomainCmd)
}

var DeleteEmailDomainemailDomainId string

func NewDeleteEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.DeleteEmailDomain(apiClient.GetConfig().Context, DeleteEmailDomainemailDomainId)

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

	cmd.Flags().StringVarP(&DeleteEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	return cmd
}

func init() {
	DeleteEmailDomainCmd := NewDeleteEmailDomainCmd()
	EmailDomainCmd.AddCommand(DeleteEmailDomainCmd)
}

var VerifyEmailDomainemailDomainId string

func NewVerifyEmailDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "verify",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EmailDomainAPI.VerifyEmailDomain(apiClient.GetConfig().Context, VerifyEmailDomainemailDomainId)

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

	cmd.Flags().StringVarP(&VerifyEmailDomainemailDomainId, "emailDomainId", "", "", "")
	cmd.MarkFlagRequired("emailDomainId")

	return cmd
}

func init() {
	VerifyEmailDomainCmd := NewVerifyEmailDomainCmd()
	EmailDomainCmd.AddCommand(VerifyEmailDomainCmd)
}
