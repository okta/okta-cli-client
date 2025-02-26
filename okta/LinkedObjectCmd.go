package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var LinkedObjectCmd = &cobra.Command{
	Use:  "linkedObject",
	Long: "Manage LinkedObjectAPI",
}

func init() {
	rootCmd.AddCommand(LinkedObjectCmd)
}

var CreateLinkedObjectDefinitiondata string

func NewCreateLinkedObjectDefinitionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createDefinition",
		Long: "Create a Linked Object Definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LinkedObjectAPI.CreateLinkedObjectDefinition(apiClient.GetConfig().Context)

			if CreateLinkedObjectDefinitiondata != "" {
				req = req.Data(CreateLinkedObjectDefinitiondata)
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

	cmd.Flags().StringVarP(&CreateLinkedObjectDefinitiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateLinkedObjectDefinitionCmd := NewCreateLinkedObjectDefinitionCmd()
	LinkedObjectCmd.AddCommand(CreateLinkedObjectDefinitionCmd)
}

func NewListLinkedObjectDefinitionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listDefinitions",
		Long: "List all Linked Object Definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LinkedObjectAPI.ListLinkedObjectDefinitions(apiClient.GetConfig().Context)

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
	ListLinkedObjectDefinitionsCmd := NewListLinkedObjectDefinitionsCmd()
	LinkedObjectCmd.AddCommand(ListLinkedObjectDefinitionsCmd)
}

var GetLinkedObjectDefinitionlinkedObjectName string

func NewGetLinkedObjectDefinitionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDefinition",
		Long: "Retrieve a Linked Object Definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LinkedObjectAPI.GetLinkedObjectDefinition(apiClient.GetConfig().Context, GetLinkedObjectDefinitionlinkedObjectName)

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

	cmd.Flags().StringVarP(&GetLinkedObjectDefinitionlinkedObjectName, "linkedObjectName", "", "", "")
	cmd.MarkFlagRequired("linkedObjectName")

	return cmd
}

func init() {
	GetLinkedObjectDefinitionCmd := NewGetLinkedObjectDefinitionCmd()
	LinkedObjectCmd.AddCommand(GetLinkedObjectDefinitionCmd)
}

var DeleteLinkedObjectDefinitionlinkedObjectName string

func NewDeleteLinkedObjectDefinitionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteDefinition",
		Long: "Delete a Linked Object Definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LinkedObjectAPI.DeleteLinkedObjectDefinition(apiClient.GetConfig().Context, DeleteLinkedObjectDefinitionlinkedObjectName)

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

	cmd.Flags().StringVarP(&DeleteLinkedObjectDefinitionlinkedObjectName, "linkedObjectName", "", "", "")
	cmd.MarkFlagRequired("linkedObjectName")

	return cmd
}

func init() {
	DeleteLinkedObjectDefinitionCmd := NewDeleteLinkedObjectDefinitionCmd()
	LinkedObjectCmd.AddCommand(DeleteLinkedObjectDefinitionCmd)
}
