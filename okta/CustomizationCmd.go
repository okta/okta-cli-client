package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var CustomizationCmd = &cobra.Command{
	Use:  "customization",
	Long: "Manage CustomizationAPI",
}

func init() {
	rootCmd.AddCommand(CustomizationCmd)
}

var CreateBranddata string

func NewCreateBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createBrand",
		Long: "Create a Brand",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.CreateBrand(apiClient.GetConfig().Context)

			if CreateBranddata != "" {
				req = req.Data(CreateBranddata)
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

	cmd.Flags().StringVarP(&CreateBranddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateBrandCmd := NewCreateBrandCmd()
	CustomizationCmd.AddCommand(CreateBrandCmd)
}

func NewListBrandsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBrands",
		Long: "List all Brands",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListBrands(apiClient.GetConfig().Context)

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
	ListBrandsCmd := NewListBrandsCmd()
	CustomizationCmd.AddCommand(ListBrandsCmd)
}

var GetBrandbrandId string

func NewGetBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getBrand",
		Long: "Retrieve a Brand",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetBrand(apiClient.GetConfig().Context, GetBrandbrandId)

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

	cmd.Flags().StringVarP(&GetBrandbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetBrandCmd := NewGetBrandCmd()
	CustomizationCmd.AddCommand(GetBrandCmd)
}

var (
	ReplaceBrandbrandId string

	ReplaceBranddata string
)

func NewReplaceBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceBrand",
		Long: "Replace a Brand",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceBrand(apiClient.GetConfig().Context, ReplaceBrandbrandId)

			if ReplaceBranddata != "" {
				req = req.Data(ReplaceBranddata)
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

	cmd.Flags().StringVarP(&ReplaceBrandbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceBranddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceBrandCmd := NewReplaceBrandCmd()
	CustomizationCmd.AddCommand(ReplaceBrandCmd)
}

var DeleteBrandbrandId string

func NewDeleteBrandCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrand",
		Long: "Delete a brand",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrand(apiClient.GetConfig().Context, DeleteBrandbrandId)

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

	cmd.Flags().StringVarP(&DeleteBrandbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	DeleteBrandCmd := NewDeleteBrandCmd()
	CustomizationCmd.AddCommand(DeleteBrandCmd)
}

var ListBrandDomainsbrandId string

func NewListBrandDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBrandDomains",
		Long: "List all Domains associated with a Brand",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListBrandDomains(apiClient.GetConfig().Context, ListBrandDomainsbrandId)

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

	cmd.Flags().StringVarP(&ListBrandDomainsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	ListBrandDomainsCmd := NewListBrandDomainsCmd()
	CustomizationCmd.AddCommand(ListBrandDomainsCmd)
}

var GetErrorPagebrandId string

func NewGetErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getErrorPage",
		Long: "Retrieve the Error Page Sub-Resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetErrorPage(apiClient.GetConfig().Context, GetErrorPagebrandId)

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

	cmd.Flags().StringVarP(&GetErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetErrorPageCmd := NewGetErrorPageCmd()
	CustomizationCmd.AddCommand(GetErrorPageCmd)
}

var GetCustomizedErrorPagebrandId string

func NewGetCustomizedErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCustomizedErrorPage",
		Long: "Retrieve the Customized Error Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetCustomizedErrorPage(apiClient.GetConfig().Context, GetCustomizedErrorPagebrandId)

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

	cmd.Flags().StringVarP(&GetCustomizedErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetCustomizedErrorPageCmd := NewGetCustomizedErrorPageCmd()
	CustomizationCmd.AddCommand(GetCustomizedErrorPageCmd)
}

var (
	ReplaceCustomizedErrorPagebrandId string

	ReplaceCustomizedErrorPagedata string
)

func NewReplaceCustomizedErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceCustomizedErrorPage",
		Long: "Replace the Customized Error Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceCustomizedErrorPage(apiClient.GetConfig().Context, ReplaceCustomizedErrorPagebrandId)

			if ReplaceCustomizedErrorPagedata != "" {
				req = req.Data(ReplaceCustomizedErrorPagedata)
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

	cmd.Flags().StringVarP(&ReplaceCustomizedErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceCustomizedErrorPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceCustomizedErrorPageCmd := NewReplaceCustomizedErrorPageCmd()
	CustomizationCmd.AddCommand(ReplaceCustomizedErrorPageCmd)
}

var DeleteCustomizedErrorPagebrandId string

func NewDeleteCustomizedErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteCustomizedErrorPage",
		Long: "Delete the Customized Error Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteCustomizedErrorPage(apiClient.GetConfig().Context, DeleteCustomizedErrorPagebrandId)

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

	cmd.Flags().StringVarP(&DeleteCustomizedErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	DeleteCustomizedErrorPageCmd := NewDeleteCustomizedErrorPageCmd()
	CustomizationCmd.AddCommand(DeleteCustomizedErrorPageCmd)
}

var GetDefaultErrorPagebrandId string

func NewGetDefaultErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDefaultErrorPage",
		Long: "Retrieve the Default Error Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetDefaultErrorPage(apiClient.GetConfig().Context, GetDefaultErrorPagebrandId)

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

	cmd.Flags().StringVarP(&GetDefaultErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetDefaultErrorPageCmd := NewGetDefaultErrorPageCmd()
	CustomizationCmd.AddCommand(GetDefaultErrorPageCmd)
}

var GetPreviewErrorPagebrandId string

func NewGetPreviewErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPreviewErrorPage",
		Long: "Retrieve the Preview Error Page Preview",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetPreviewErrorPage(apiClient.GetConfig().Context, GetPreviewErrorPagebrandId)

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

	cmd.Flags().StringVarP(&GetPreviewErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetPreviewErrorPageCmd := NewGetPreviewErrorPageCmd()
	CustomizationCmd.AddCommand(GetPreviewErrorPageCmd)
}

var (
	ReplacePreviewErrorPagebrandId string

	ReplacePreviewErrorPagedata string
)

func NewReplacePreviewErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePreviewErrorPage",
		Long: "Replace the Preview Error Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplacePreviewErrorPage(apiClient.GetConfig().Context, ReplacePreviewErrorPagebrandId)

			if ReplacePreviewErrorPagedata != "" {
				req = req.Data(ReplacePreviewErrorPagedata)
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

	cmd.Flags().StringVarP(&ReplacePreviewErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplacePreviewErrorPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacePreviewErrorPageCmd := NewReplacePreviewErrorPageCmd()
	CustomizationCmd.AddCommand(ReplacePreviewErrorPageCmd)
}

var DeletePreviewErrorPagebrandId string

func NewDeletePreviewErrorPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletePreviewErrorPage",
		Long: "Delete the Preview Error Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeletePreviewErrorPage(apiClient.GetConfig().Context, DeletePreviewErrorPagebrandId)

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

	cmd.Flags().StringVarP(&DeletePreviewErrorPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	DeletePreviewErrorPageCmd := NewDeletePreviewErrorPageCmd()
	CustomizationCmd.AddCommand(DeletePreviewErrorPageCmd)
}

var GetSignInPagebrandId string

func NewGetSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSignInPage",
		Long: "Retrieve the Sign-in Page Sub-Resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetSignInPage(apiClient.GetConfig().Context, GetSignInPagebrandId)

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

	cmd.Flags().StringVarP(&GetSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetSignInPageCmd := NewGetSignInPageCmd()
	CustomizationCmd.AddCommand(GetSignInPageCmd)
}

var GetCustomizedSignInPagebrandId string

func NewGetCustomizedSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCustomizedSignInPage",
		Long: "Retrieve the Customized Sign-in Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetCustomizedSignInPage(apiClient.GetConfig().Context, GetCustomizedSignInPagebrandId)

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

	cmd.Flags().StringVarP(&GetCustomizedSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetCustomizedSignInPageCmd := NewGetCustomizedSignInPageCmd()
	CustomizationCmd.AddCommand(GetCustomizedSignInPageCmd)
}

var (
	ReplaceCustomizedSignInPagebrandId string

	ReplaceCustomizedSignInPagedata string
)

func NewReplaceCustomizedSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceCustomizedSignInPage",
		Long: "Replace the Customized Sign-in Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceCustomizedSignInPage(apiClient.GetConfig().Context, ReplaceCustomizedSignInPagebrandId)

			if ReplaceCustomizedSignInPagedata != "" {
				req = req.Data(ReplaceCustomizedSignInPagedata)
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

	cmd.Flags().StringVarP(&ReplaceCustomizedSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceCustomizedSignInPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceCustomizedSignInPageCmd := NewReplaceCustomizedSignInPageCmd()
	CustomizationCmd.AddCommand(ReplaceCustomizedSignInPageCmd)
}

var DeleteCustomizedSignInPagebrandId string

func NewDeleteCustomizedSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteCustomizedSignInPage",
		Long: "Delete the Customized Sign-in Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteCustomizedSignInPage(apiClient.GetConfig().Context, DeleteCustomizedSignInPagebrandId)

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

	cmd.Flags().StringVarP(&DeleteCustomizedSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	DeleteCustomizedSignInPageCmd := NewDeleteCustomizedSignInPageCmd()
	CustomizationCmd.AddCommand(DeleteCustomizedSignInPageCmd)
}

var GetDefaultSignInPagebrandId string

func NewGetDefaultSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDefaultSignInPage",
		Long: "Retrieve the Default Sign-in Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetDefaultSignInPage(apiClient.GetConfig().Context, GetDefaultSignInPagebrandId)

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

	cmd.Flags().StringVarP(&GetDefaultSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetDefaultSignInPageCmd := NewGetDefaultSignInPageCmd()
	CustomizationCmd.AddCommand(GetDefaultSignInPageCmd)
}

var GetPreviewSignInPagebrandId string

func NewGetPreviewSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPreviewSignInPage",
		Long: "Retrieve the Preview Sign-in Page Preview",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetPreviewSignInPage(apiClient.GetConfig().Context, GetPreviewSignInPagebrandId)

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

	cmd.Flags().StringVarP(&GetPreviewSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetPreviewSignInPageCmd := NewGetPreviewSignInPageCmd()
	CustomizationCmd.AddCommand(GetPreviewSignInPageCmd)
}

var (
	ReplacePreviewSignInPagebrandId string

	ReplacePreviewSignInPagedata string
)

func NewReplacePreviewSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replacePreviewSignInPage",
		Long: "Replace the Preview Sign-in Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplacePreviewSignInPage(apiClient.GetConfig().Context, ReplacePreviewSignInPagebrandId)

			if ReplacePreviewSignInPagedata != "" {
				req = req.Data(ReplacePreviewSignInPagedata)
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

	cmd.Flags().StringVarP(&ReplacePreviewSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplacePreviewSignInPagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacePreviewSignInPageCmd := NewReplacePreviewSignInPageCmd()
	CustomizationCmd.AddCommand(ReplacePreviewSignInPageCmd)
}

var DeletePreviewSignInPagebrandId string

func NewDeletePreviewSignInPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deletePreviewSignInPage",
		Long: "Delete the Preview Sign-in Page",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeletePreviewSignInPage(apiClient.GetConfig().Context, DeletePreviewSignInPagebrandId)

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

	cmd.Flags().StringVarP(&DeletePreviewSignInPagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	DeletePreviewSignInPageCmd := NewDeletePreviewSignInPageCmd()
	CustomizationCmd.AddCommand(DeletePreviewSignInPageCmd)
}

var ListAllSignInWidgetVersionsbrandId string

func NewListAllSignInWidgetVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAllSignInWidgetVersions",
		Long: "List all Sign-in Widget Versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListAllSignInWidgetVersions(apiClient.GetConfig().Context, ListAllSignInWidgetVersionsbrandId)

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

	cmd.Flags().StringVarP(&ListAllSignInWidgetVersionsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	ListAllSignInWidgetVersionsCmd := NewListAllSignInWidgetVersionsCmd()
	CustomizationCmd.AddCommand(ListAllSignInWidgetVersionsCmd)
}

var GetSignOutPageSettingsbrandId string

func NewGetSignOutPageSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getSignOutPageSettings",
		Long: "Retrieve the Sign-out Page Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetSignOutPageSettings(apiClient.GetConfig().Context, GetSignOutPageSettingsbrandId)

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

	cmd.Flags().StringVarP(&GetSignOutPageSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	GetSignOutPageSettingsCmd := NewGetSignOutPageSettingsCmd()
	CustomizationCmd.AddCommand(GetSignOutPageSettingsCmd)
}

var (
	ReplaceSignOutPageSettingsbrandId string

	ReplaceSignOutPageSettingsdata string
)

func NewReplaceSignOutPageSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceSignOutPageSettings",
		Long: "Replace the Sign-out Page Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceSignOutPageSettings(apiClient.GetConfig().Context, ReplaceSignOutPageSettingsbrandId)

			if ReplaceSignOutPageSettingsdata != "" {
				req = req.Data(ReplaceSignOutPageSettingsdata)
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

	cmd.Flags().StringVarP(&ReplaceSignOutPageSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceSignOutPageSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceSignOutPageSettingsCmd := NewReplaceSignOutPageSettingsCmd()
	CustomizationCmd.AddCommand(ReplaceSignOutPageSettingsCmd)
}

var ListEmailTemplatesbrandId string

func NewListEmailTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listEmailTemplates",
		Long: "List all Email Templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListEmailTemplates(apiClient.GetConfig().Context, ListEmailTemplatesbrandId)

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

	cmd.Flags().StringVarP(&ListEmailTemplatesbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	ListEmailTemplatesCmd := NewListEmailTemplatesCmd()
	CustomizationCmd.AddCommand(ListEmailTemplatesCmd)
}

var (
	GetEmailTemplatebrandId string

	GetEmailTemplatetemplateName string
)

func NewGetEmailTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailTemplate",
		Long: "Retrieve an Email Template",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailTemplate(apiClient.GetConfig().Context, GetEmailTemplatebrandId, GetEmailTemplatetemplateName)

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

	cmd.Flags().StringVarP(&GetEmailTemplatebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailTemplatetemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	GetEmailTemplateCmd := NewGetEmailTemplateCmd()
	CustomizationCmd.AddCommand(GetEmailTemplateCmd)
}

var (
	CreateEmailCustomizationbrandId string

	CreateEmailCustomizationtemplateName string

	CreateEmailCustomizationdata string
)

func NewCreateEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createEmail",
		Long: "Create an Email Customization",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.CreateEmailCustomization(apiClient.GetConfig().Context, CreateEmailCustomizationbrandId, CreateEmailCustomizationtemplateName)

			if CreateEmailCustomizationdata != "" {
				req = req.Data(CreateEmailCustomizationdata)
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

	cmd.Flags().StringVarP(&CreateEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&CreateEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&CreateEmailCustomizationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateEmailCustomizationCmd := NewCreateEmailCustomizationCmd()
	CustomizationCmd.AddCommand(CreateEmailCustomizationCmd)
}

var (
	ListEmailCustomizationsbrandId string

	ListEmailCustomizationstemplateName string
)

func NewListEmailCustomizationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listEmails",
		Long: "List all Email Customizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListEmailCustomizations(apiClient.GetConfig().Context, ListEmailCustomizationsbrandId, ListEmailCustomizationstemplateName)

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

	cmd.Flags().StringVarP(&ListEmailCustomizationsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ListEmailCustomizationstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	ListEmailCustomizationsCmd := NewListEmailCustomizationsCmd()
	CustomizationCmd.AddCommand(ListEmailCustomizationsCmd)
}

var (
	DeleteAllCustomizationsbrandId string

	DeleteAllCustomizationstemplateName string
)

func NewDeleteAllCustomizationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteAlls",
		Long: "Delete all Email Customizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteAllCustomizations(apiClient.GetConfig().Context, DeleteAllCustomizationsbrandId, DeleteAllCustomizationstemplateName)

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

	cmd.Flags().StringVarP(&DeleteAllCustomizationsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteAllCustomizationstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	DeleteAllCustomizationsCmd := NewDeleteAllCustomizationsCmd()
	CustomizationCmd.AddCommand(DeleteAllCustomizationsCmd)
}

var (
	GetEmailCustomizationbrandId string

	GetEmailCustomizationtemplateName string

	GetEmailCustomizationcustomizationId string
)

func NewGetEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmail",
		Long: "Retrieve an Email Customization",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailCustomization(apiClient.GetConfig().Context, GetEmailCustomizationbrandId, GetEmailCustomizationtemplateName, GetEmailCustomizationcustomizationId)

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

	cmd.Flags().StringVarP(&GetEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&GetEmailCustomizationcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	return cmd
}

func init() {
	GetEmailCustomizationCmd := NewGetEmailCustomizationCmd()
	CustomizationCmd.AddCommand(GetEmailCustomizationCmd)
}

var (
	ReplaceEmailCustomizationbrandId string

	ReplaceEmailCustomizationtemplateName string

	ReplaceEmailCustomizationcustomizationId string

	ReplaceEmailCustomizationdata string
)

func NewReplaceEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceEmail",
		Long: "Replace an Email Customization",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceEmailCustomization(apiClient.GetConfig().Context, ReplaceEmailCustomizationbrandId, ReplaceEmailCustomizationtemplateName, ReplaceEmailCustomizationcustomizationId)

			if ReplaceEmailCustomizationdata != "" {
				req = req.Data(ReplaceEmailCustomizationdata)
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

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	cmd.Flags().StringVarP(&ReplaceEmailCustomizationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceEmailCustomizationCmd := NewReplaceEmailCustomizationCmd()
	CustomizationCmd.AddCommand(ReplaceEmailCustomizationCmd)
}

var (
	DeleteEmailCustomizationbrandId string

	DeleteEmailCustomizationtemplateName string

	DeleteEmailCustomizationcustomizationId string
)

func NewDeleteEmailCustomizationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteEmail",
		Long: "Delete an Email Customization",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteEmailCustomization(apiClient.GetConfig().Context, DeleteEmailCustomizationbrandId, DeleteEmailCustomizationtemplateName, DeleteEmailCustomizationcustomizationId)

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

	cmd.Flags().StringVarP(&DeleteEmailCustomizationbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteEmailCustomizationtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&DeleteEmailCustomizationcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	return cmd
}

func init() {
	DeleteEmailCustomizationCmd := NewDeleteEmailCustomizationCmd()
	CustomizationCmd.AddCommand(DeleteEmailCustomizationCmd)
}

var (
	GetCustomizationPreviewbrandId string

	GetCustomizationPreviewtemplateName string

	GetCustomizationPreviewcustomizationId string
)

func NewGetCustomizationPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getPreview",
		Long: "Retrieve a Preview of an Email Customization",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetCustomizationPreview(apiClient.GetConfig().Context, GetCustomizationPreviewbrandId, GetCustomizationPreviewtemplateName, GetCustomizationPreviewcustomizationId)

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

	cmd.Flags().StringVarP(&GetCustomizationPreviewbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetCustomizationPreviewtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&GetCustomizationPreviewcustomizationId, "customizationId", "", "", "")
	cmd.MarkFlagRequired("customizationId")

	return cmd
}

func init() {
	GetCustomizationPreviewCmd := NewGetCustomizationPreviewCmd()
	CustomizationCmd.AddCommand(GetCustomizationPreviewCmd)
}

var (
	GetEmailDefaultContentbrandId string

	GetEmailDefaultContenttemplateName string
)

func NewGetEmailDefaultContentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailDefaultContent",
		Long: "Retrieve an Email Template Default Content",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailDefaultContent(apiClient.GetConfig().Context, GetEmailDefaultContentbrandId, GetEmailDefaultContenttemplateName)

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

	cmd.Flags().StringVarP(&GetEmailDefaultContentbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailDefaultContenttemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	GetEmailDefaultContentCmd := NewGetEmailDefaultContentCmd()
	CustomizationCmd.AddCommand(GetEmailDefaultContentCmd)
}

var (
	GetEmailDefaultPreviewbrandId string

	GetEmailDefaultPreviewtemplateName string
)

func NewGetEmailDefaultPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailDefaultPreview",
		Long: "Retrieve a Preview of the Email Template default content",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailDefaultPreview(apiClient.GetConfig().Context, GetEmailDefaultPreviewbrandId, GetEmailDefaultPreviewtemplateName)

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

	cmd.Flags().StringVarP(&GetEmailDefaultPreviewbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailDefaultPreviewtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	GetEmailDefaultPreviewCmd := NewGetEmailDefaultPreviewCmd()
	CustomizationCmd.AddCommand(GetEmailDefaultPreviewCmd)
}

var (
	GetEmailSettingsbrandId string

	GetEmailSettingstemplateName string
)

func NewGetEmailSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getEmailSettings",
		Long: "Retrieve the Email Template Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetEmailSettings(apiClient.GetConfig().Context, GetEmailSettingsbrandId, GetEmailSettingstemplateName)

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

	cmd.Flags().StringVarP(&GetEmailSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetEmailSettingstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	GetEmailSettingsCmd := NewGetEmailSettingsCmd()
	CustomizationCmd.AddCommand(GetEmailSettingsCmd)
}

var (
	ReplaceEmailSettingsbrandId string

	ReplaceEmailSettingstemplateName string

	ReplaceEmailSettingsdata string
)

func NewReplaceEmailSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceEmailSettings",
		Long: "Replace the Email Template Settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceEmailSettings(apiClient.GetConfig().Context, ReplaceEmailSettingsbrandId, ReplaceEmailSettingstemplateName)

			if ReplaceEmailSettingsdata != "" {
				req = req.Data(ReplaceEmailSettingsdata)
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

	cmd.Flags().StringVarP(&ReplaceEmailSettingsbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceEmailSettingstemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	cmd.Flags().StringVarP(&ReplaceEmailSettingsdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceEmailSettingsCmd := NewReplaceEmailSettingsCmd()
	CustomizationCmd.AddCommand(ReplaceEmailSettingsCmd)
}

var (
	SendTestEmailbrandId string

	SendTestEmailtemplateName string
)

func NewSendTestEmailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "sendTestEmail",
		Long: "Send a Test Email",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.SendTestEmail(apiClient.GetConfig().Context, SendTestEmailbrandId, SendTestEmailtemplateName)

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

	cmd.Flags().StringVarP(&SendTestEmailbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&SendTestEmailtemplateName, "templateName", "", "", "")
	cmd.MarkFlagRequired("templateName")

	return cmd
}

func init() {
	SendTestEmailCmd := NewSendTestEmailCmd()
	CustomizationCmd.AddCommand(SendTestEmailCmd)
}

var ListBrandThemesbrandId string

func NewListBrandThemesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBrandThemes",
		Long: "List all Themes",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ListBrandThemes(apiClient.GetConfig().Context, ListBrandThemesbrandId)

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

	cmd.Flags().StringVarP(&ListBrandThemesbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	return cmd
}

func init() {
	ListBrandThemesCmd := NewListBrandThemesCmd()
	CustomizationCmd.AddCommand(ListBrandThemesCmd)
}

var (
	GetBrandThemebrandId string

	GetBrandThemethemeId string
)

func NewGetBrandThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getBrandTheme",
		Long: "Retrieve a Theme",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.GetBrandTheme(apiClient.GetConfig().Context, GetBrandThemebrandId, GetBrandThemethemeId)

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

	cmd.Flags().StringVarP(&GetBrandThemebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&GetBrandThemethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	return cmd
}

func init() {
	GetBrandThemeCmd := NewGetBrandThemeCmd()
	CustomizationCmd.AddCommand(GetBrandThemeCmd)
}

var (
	ReplaceBrandThemebrandId string

	ReplaceBrandThemethemeId string

	ReplaceBrandThemedata string
)

func NewReplaceBrandThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceBrandTheme",
		Long: "Replace a Theme",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.ReplaceBrandTheme(apiClient.GetConfig().Context, ReplaceBrandThemebrandId, ReplaceBrandThemethemeId)

			if ReplaceBrandThemedata != "" {
				req = req.Data(ReplaceBrandThemedata)
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

	cmd.Flags().StringVarP(&ReplaceBrandThemebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&ReplaceBrandThemethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&ReplaceBrandThemedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceBrandThemeCmd := NewReplaceBrandThemeCmd()
	CustomizationCmd.AddCommand(ReplaceBrandThemeCmd)
}

var (
	UploadBrandThemeBackgroundImagebrandId string

	UploadBrandThemeBackgroundImagethemeId string

	UploadBrandThemeBackgroundImagedata string
)

func NewUploadBrandThemeBackgroundImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadBrandThemeBackgroundImage",
		Long: "Upload the Background Image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.UploadBrandThemeBackgroundImage(apiClient.GetConfig().Context, UploadBrandThemeBackgroundImagebrandId, UploadBrandThemeBackgroundImagethemeId)

			if UploadBrandThemeBackgroundImagedata != "" {
				req = req.Data(UploadBrandThemeBackgroundImagedata)
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

	cmd.Flags().StringVarP(&UploadBrandThemeBackgroundImagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&UploadBrandThemeBackgroundImagethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&UploadBrandThemeBackgroundImagedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UploadBrandThemeBackgroundImageCmd := NewUploadBrandThemeBackgroundImageCmd()
	CustomizationCmd.AddCommand(UploadBrandThemeBackgroundImageCmd)
}

var (
	DeleteBrandThemeBackgroundImagebrandId string

	DeleteBrandThemeBackgroundImagethemeId string
)

func NewDeleteBrandThemeBackgroundImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrandThemeBackgroundImage",
		Long: "Delete the Background Image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrandThemeBackgroundImage(apiClient.GetConfig().Context, DeleteBrandThemeBackgroundImagebrandId, DeleteBrandThemeBackgroundImagethemeId)

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

	cmd.Flags().StringVarP(&DeleteBrandThemeBackgroundImagebrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteBrandThemeBackgroundImagethemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	return cmd
}

func init() {
	DeleteBrandThemeBackgroundImageCmd := NewDeleteBrandThemeBackgroundImageCmd()
	CustomizationCmd.AddCommand(DeleteBrandThemeBackgroundImageCmd)
}

var (
	UploadBrandThemeFaviconbrandId string

	UploadBrandThemeFaviconthemeId string

	UploadBrandThemeFavicondata string
)

func NewUploadBrandThemeFaviconCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadBrandThemeFavicon",
		Long: "Upload the Favicon",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.UploadBrandThemeFavicon(apiClient.GetConfig().Context, UploadBrandThemeFaviconbrandId, UploadBrandThemeFaviconthemeId)

			if UploadBrandThemeFavicondata != "" {
				req = req.Data(UploadBrandThemeFavicondata)
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

	cmd.Flags().StringVarP(&UploadBrandThemeFaviconbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&UploadBrandThemeFaviconthemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&UploadBrandThemeFavicondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UploadBrandThemeFaviconCmd := NewUploadBrandThemeFaviconCmd()
	CustomizationCmd.AddCommand(UploadBrandThemeFaviconCmd)
}

var (
	DeleteBrandThemeFaviconbrandId string

	DeleteBrandThemeFaviconthemeId string
)

func NewDeleteBrandThemeFaviconCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrandThemeFavicon",
		Long: "Delete the Favicon",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrandThemeFavicon(apiClient.GetConfig().Context, DeleteBrandThemeFaviconbrandId, DeleteBrandThemeFaviconthemeId)

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

	cmd.Flags().StringVarP(&DeleteBrandThemeFaviconbrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteBrandThemeFaviconthemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	return cmd
}

func init() {
	DeleteBrandThemeFaviconCmd := NewDeleteBrandThemeFaviconCmd()
	CustomizationCmd.AddCommand(DeleteBrandThemeFaviconCmd)
}

var (
	UploadBrandThemeLogobrandId string

	UploadBrandThemeLogothemeId string

	UploadBrandThemeLogodata string
)

func NewUploadBrandThemeLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "uploadBrandThemeLogo",
		Long: "Upload the Logo",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.UploadBrandThemeLogo(apiClient.GetConfig().Context, UploadBrandThemeLogobrandId, UploadBrandThemeLogothemeId)

			if UploadBrandThemeLogodata != "" {
				req = req.Data(UploadBrandThemeLogodata)
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

	cmd.Flags().StringVarP(&UploadBrandThemeLogobrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&UploadBrandThemeLogothemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	cmd.Flags().StringVarP(&UploadBrandThemeLogodata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UploadBrandThemeLogoCmd := NewUploadBrandThemeLogoCmd()
	CustomizationCmd.AddCommand(UploadBrandThemeLogoCmd)
}

var (
	DeleteBrandThemeLogobrandId string

	DeleteBrandThemeLogothemeId string
)

func NewDeleteBrandThemeLogoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteBrandThemeLogo",
		Long: "Delete the Logo",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.CustomizationAPI.DeleteBrandThemeLogo(apiClient.GetConfig().Context, DeleteBrandThemeLogobrandId, DeleteBrandThemeLogothemeId)

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

	cmd.Flags().StringVarP(&DeleteBrandThemeLogobrandId, "brandId", "", "", "")
	cmd.MarkFlagRequired("brandId")

	cmd.Flags().StringVarP(&DeleteBrandThemeLogothemeId, "themeId", "", "", "")
	cmd.MarkFlagRequired("themeId")

	return cmd
}

func init() {
	DeleteBrandThemeLogoCmd := NewDeleteBrandThemeLogoCmd()
	CustomizationCmd.AddCommand(DeleteBrandThemeLogoCmd)
}
