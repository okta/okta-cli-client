package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var CustomDomainCmd = &cobra.Command{
	Use:  "customDomain",
	Long: "Manage CustomDomainAPI",
}

func init() {
	rootCmd.AddCommand(CustomDomainCmd)
}

var CreateCustomDomaindata string

func NewCreateCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Custom Domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.CreateCustomDomain(apiClient.GetConfig().Context)

			if CreateCustomDomaindata != "" {
				req = req.Data(CreateCustomDomaindata)
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

	cmd.Flags().StringVarP(&CreateCustomDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateCustomDomainCmd := NewCreateCustomDomainCmd()
	CustomDomainCmd.AddCommand(CreateCustomDomainCmd)
}

func NewListCustomDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Custom Domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.ListCustomDomains(apiClient.GetConfig().Context)

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
	ListCustomDomainsCmd := NewListCustomDomainsCmd()
	CustomDomainCmd.AddCommand(ListCustomDomainsCmd)
}

var GetCustomDomaindomainId string

func NewGetCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Custom Domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.GetCustomDomain(apiClient.GetConfig().Context, GetCustomDomaindomainId)

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

	cmd.Flags().StringVarP(&GetCustomDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	return cmd
}

func init() {
	GetCustomDomainCmd := NewGetCustomDomainCmd()
	CustomDomainCmd.AddCommand(GetCustomDomainCmd)
}

var (
	ReplaceCustomDomaindomainId string

	ReplaceCustomDomaindata string
)

func NewReplaceCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Custom Domain's Brand",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.ReplaceCustomDomain(apiClient.GetConfig().Context, ReplaceCustomDomaindomainId)

			if ReplaceCustomDomaindata != "" {
				req = req.Data(ReplaceCustomDomaindata)
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

	cmd.Flags().StringVarP(&ReplaceCustomDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().StringVarP(&ReplaceCustomDomaindata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceCustomDomainCmd := NewReplaceCustomDomainCmd()
	CustomDomainCmd.AddCommand(ReplaceCustomDomainCmd)
}

var DeleteCustomDomaindomainId string

func NewDeleteCustomDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Custom Domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.DeleteCustomDomain(apiClient.GetConfig().Context, DeleteCustomDomaindomainId)

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

	cmd.Flags().StringVarP(&DeleteCustomDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	return cmd
}

func init() {
	DeleteCustomDomainCmd := NewDeleteCustomDomainCmd()
	CustomDomainCmd.AddCommand(DeleteCustomDomainCmd)
}

var (
	UpsertCertificatedomainId string

	UpsertCertificatedata string
)

func NewUpsertCertificateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "upsertCertificate",
		Long: "Upsert the Custom Domain's Certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.UpsertCertificate(apiClient.GetConfig().Context, UpsertCertificatedomainId)

			if UpsertCertificatedata != "" {
				req = req.Data(UpsertCertificatedata)
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

	cmd.Flags().StringVarP(&UpsertCertificatedomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	cmd.Flags().StringVarP(&UpsertCertificatedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpsertCertificateCmd := NewUpsertCertificateCmd()
	CustomDomainCmd.AddCommand(UpsertCertificateCmd)
}

var VerifyDomaindomainId string

func NewVerifyDomainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "verifyDomain",
		Long: "Verify a Custom Domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomDomainAPI.VerifyDomain(apiClient.GetConfig().Context, VerifyDomaindomainId)

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

	cmd.Flags().StringVarP(&VerifyDomaindomainId, "domainId", "", "", "")
	cmd.MarkFlagRequired("domainId")

	return cmd
}

func init() {
	VerifyDomainCmd := NewVerifyDomainCmd()
	CustomDomainCmd.AddCommand(VerifyDomainCmd)
}
