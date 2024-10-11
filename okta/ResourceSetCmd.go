package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ResourceSetCmd = &cobra.Command{
	Use:  "resourceSet",
	Long: "Manage ResourceSetAPI",
}

func init() {
	rootCmd.AddCommand(ResourceSetCmd)
}

var CreateResourceSetdata string

func NewCreateResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.CreateResourceSet(apiClient.GetConfig().Context)

			if CreateResourceSetdata != "" {
				req = req.Data(CreateResourceSetdata)
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

	cmd.Flags().StringVarP(&CreateResourceSetdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateResourceSetCmd := NewCreateResourceSetCmd()
	ResourceSetCmd.AddCommand(CreateResourceSetCmd)
}

func NewListResourceSetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListResourceSets(apiClient.GetConfig().Context)

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
	ListResourceSetsCmd := NewListResourceSetsCmd()
	ResourceSetCmd.AddCommand(ListResourceSetsCmd)
}

var GetResourceSetresourceSetId string

func NewGetResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.GetResourceSet(apiClient.GetConfig().Context, GetResourceSetresourceSetId)

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

	cmd.Flags().StringVarP(&GetResourceSetresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	return cmd
}

func init() {
	GetResourceSetCmd := NewGetResourceSetCmd()
	ResourceSetCmd.AddCommand(GetResourceSetCmd)
}

var (
	ReplaceResourceSetresourceSetId string

	ReplaceResourceSetdata string
)

func NewReplaceResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ReplaceResourceSet(apiClient.GetConfig().Context, ReplaceResourceSetresourceSetId)

			if ReplaceResourceSetdata != "" {
				req = req.Data(ReplaceResourceSetdata)
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

	cmd.Flags().StringVarP(&ReplaceResourceSetresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&ReplaceResourceSetdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceResourceSetCmd := NewReplaceResourceSetCmd()
	ResourceSetCmd.AddCommand(ReplaceResourceSetCmd)
}

var DeleteResourceSetresourceSetId string

func NewDeleteResourceSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.DeleteResourceSet(apiClient.GetConfig().Context, DeleteResourceSetresourceSetId)

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

	cmd.Flags().StringVarP(&DeleteResourceSetresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	return cmd
}

func init() {
	DeleteResourceSetCmd := NewDeleteResourceSetCmd()
	ResourceSetCmd.AddCommand(DeleteResourceSetCmd)
}

var (
	CreateResourceSetBindingresourceSetId string

	CreateResourceSetBindingdata string
)

func NewCreateResourceSetBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.CreateResourceSetBinding(apiClient.GetConfig().Context, CreateResourceSetBindingresourceSetId)

			if CreateResourceSetBindingdata != "" {
				req = req.Data(CreateResourceSetBindingdata)
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

	cmd.Flags().StringVarP(&CreateResourceSetBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&CreateResourceSetBindingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateResourceSetBindingCmd := NewCreateResourceSetBindingCmd()
	ResourceSetCmd.AddCommand(CreateResourceSetBindingCmd)
}

var ListBindingsresourceSetId string

func NewListBindingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listBindings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListBindings(apiClient.GetConfig().Context, ListBindingsresourceSetId)

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

	cmd.Flags().StringVarP(&ListBindingsresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	return cmd
}

func init() {
	ListBindingsCmd := NewListBindingsCmd()
	ResourceSetCmd.AddCommand(ListBindingsCmd)
}

var (
	GetBindingresourceSetId string

	GetBindingroleIdOrLabel string
)

func NewGetBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.GetBinding(apiClient.GetConfig().Context, GetBindingresourceSetId, GetBindingroleIdOrLabel)

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

	cmd.Flags().StringVarP(&GetBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&GetBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	return cmd
}

func init() {
	GetBindingCmd := NewGetBindingCmd()
	ResourceSetCmd.AddCommand(GetBindingCmd)
}

var (
	DeleteBindingresourceSetId string

	DeleteBindingroleIdOrLabel string
)

func NewDeleteBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.DeleteBinding(apiClient.GetConfig().Context, DeleteBindingresourceSetId, DeleteBindingroleIdOrLabel)

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

	cmd.Flags().StringVarP(&DeleteBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&DeleteBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	return cmd
}

func init() {
	DeleteBindingCmd := NewDeleteBindingCmd()
	ResourceSetCmd.AddCommand(DeleteBindingCmd)
}

var (
	ListMembersOfBindingresourceSetId string

	ListMembersOfBindingroleIdOrLabel string
)

func NewListMembersOfBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listMembersOfBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListMembersOfBinding(apiClient.GetConfig().Context, ListMembersOfBindingresourceSetId, ListMembersOfBindingroleIdOrLabel)

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

	cmd.Flags().StringVarP(&ListMembersOfBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&ListMembersOfBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	return cmd
}

func init() {
	ListMembersOfBindingCmd := NewListMembersOfBindingCmd()
	ResourceSetCmd.AddCommand(ListMembersOfBindingCmd)
}

var (
	AddMembersToBindingresourceSetId string

	AddMembersToBindingroleIdOrLabel string

	AddMembersToBindingdata string
)

func NewAddMembersToBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "addMembersToBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.AddMembersToBinding(apiClient.GetConfig().Context, AddMembersToBindingresourceSetId, AddMembersToBindingroleIdOrLabel)

			if AddMembersToBindingdata != "" {
				req = req.Data(AddMembersToBindingdata)
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

	cmd.Flags().StringVarP(&AddMembersToBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&AddMembersToBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&AddMembersToBindingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AddMembersToBindingCmd := NewAddMembersToBindingCmd()
	ResourceSetCmd.AddCommand(AddMembersToBindingCmd)
}

var (
	GetMemberOfBindingresourceSetId string

	GetMemberOfBindingroleIdOrLabel string

	GetMemberOfBindingmemberId string
)

func NewGetMemberOfBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getMemberOfBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.GetMemberOfBinding(apiClient.GetConfig().Context, GetMemberOfBindingresourceSetId, GetMemberOfBindingroleIdOrLabel, GetMemberOfBindingmemberId)

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

	cmd.Flags().StringVarP(&GetMemberOfBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&GetMemberOfBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&GetMemberOfBindingmemberId, "memberId", "", "", "")
	cmd.MarkFlagRequired("memberId")

	return cmd
}

func init() {
	GetMemberOfBindingCmd := NewGetMemberOfBindingCmd()
	ResourceSetCmd.AddCommand(GetMemberOfBindingCmd)
}

var (
	UnassignMemberFromBindingresourceSetId string

	UnassignMemberFromBindingroleIdOrLabel string

	UnassignMemberFromBindingmemberId string
)

func NewUnassignMemberFromBindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "unassignMemberFromBinding",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.UnassignMemberFromBinding(apiClient.GetConfig().Context, UnassignMemberFromBindingresourceSetId, UnassignMemberFromBindingroleIdOrLabel, UnassignMemberFromBindingmemberId)

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

	cmd.Flags().StringVarP(&UnassignMemberFromBindingresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&UnassignMemberFromBindingroleIdOrLabel, "roleIdOrLabel", "", "", "")
	cmd.MarkFlagRequired("roleIdOrLabel")

	cmd.Flags().StringVarP(&UnassignMemberFromBindingmemberId, "memberId", "", "", "")
	cmd.MarkFlagRequired("memberId")

	return cmd
}

func init() {
	UnassignMemberFromBindingCmd := NewUnassignMemberFromBindingCmd()
	ResourceSetCmd.AddCommand(UnassignMemberFromBindingCmd)
}

var ListResourceSetResourcesresourceSetId string

func NewListResourceSetResourcesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listResources",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.ListResourceSetResources(apiClient.GetConfig().Context, ListResourceSetResourcesresourceSetId)

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

	cmd.Flags().StringVarP(&ListResourceSetResourcesresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	return cmd
}

func init() {
	ListResourceSetResourcesCmd := NewListResourceSetResourcesCmd()
	ResourceSetCmd.AddCommand(ListResourceSetResourcesCmd)
}

var (
	AddResourceSetResourceresourceSetId string

	AddResourceSetResourcedata string
)

func NewAddResourceSetResourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "addResource",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.AddResourceSetResource(apiClient.GetConfig().Context, AddResourceSetResourceresourceSetId)

			if AddResourceSetResourcedata != "" {
				req = req.Data(AddResourceSetResourcedata)
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

	cmd.Flags().StringVarP(&AddResourceSetResourceresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&AddResourceSetResourcedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AddResourceSetResourceCmd := NewAddResourceSetResourceCmd()
	ResourceSetCmd.AddCommand(AddResourceSetResourceCmd)
}

var (
	DeleteResourceSetResourceresourceSetId string

	DeleteResourceSetResourceresourceId string
)

func NewDeleteResourceSetResourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteResource",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ResourceSetAPI.DeleteResourceSetResource(apiClient.GetConfig().Context, DeleteResourceSetResourceresourceSetId, DeleteResourceSetResourceresourceId)

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

	cmd.Flags().StringVarP(&DeleteResourceSetResourceresourceSetId, "resourceSetId", "", "", "")
	cmd.MarkFlagRequired("resourceSetId")

	cmd.Flags().StringVarP(&DeleteResourceSetResourceresourceId, "resourceId", "", "", "")
	cmd.MarkFlagRequired("resourceId")

	return cmd
}

func init() {
	DeleteResourceSetResourceCmd := NewDeleteResourceSetResourceCmd()
	ResourceSetCmd.AddCommand(DeleteResourceSetResourceCmd)
}
