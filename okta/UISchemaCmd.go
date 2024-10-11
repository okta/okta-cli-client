package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var UISchemaCmd = &cobra.Command{
	Use:  "uISchema",
	Long: "Manage UISchemaAPI",
}

func init() {
	rootCmd.AddCommand(UISchemaCmd)
}

var CreateUISchemadata string

func NewCreateUISchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.CreateUISchema(apiClient.GetConfig().Context)

			if CreateUISchemadata != "" {
				req = req.Data(CreateUISchemadata)
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

	cmd.Flags().StringVarP(&CreateUISchemadata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateUISchemaCmd := NewCreateUISchemaCmd()
	UISchemaCmd.AddCommand(CreateUISchemaCmd)
}

func NewListUISchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.ListUISchemas(apiClient.GetConfig().Context)

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
	ListUISchemasCmd := NewListUISchemasCmd()
	UISchemaCmd.AddCommand(ListUISchemasCmd)
}

var GetUISchemaid string

func NewGetUISchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.GetUISchema(apiClient.GetConfig().Context, GetUISchemaid)

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

	cmd.Flags().StringVarP(&GetUISchemaid, "id", "", "", "")
	cmd.MarkFlagRequired("id")

	return cmd
}

func init() {
	GetUISchemaCmd := NewGetUISchemaCmd()
	UISchemaCmd.AddCommand(GetUISchemaCmd)
}

var (
	ReplaceUISchemasid string

	ReplaceUISchemasdata string
)

func NewReplaceUISchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaces",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.ReplaceUISchemas(apiClient.GetConfig().Context, ReplaceUISchemasid)

			if ReplaceUISchemasdata != "" {
				req = req.Data(ReplaceUISchemasdata)
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

	cmd.Flags().StringVarP(&ReplaceUISchemasid, "id", "", "", "")
	cmd.MarkFlagRequired("id")

	cmd.Flags().StringVarP(&ReplaceUISchemasdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceUISchemasCmd := NewReplaceUISchemasCmd()
	UISchemaCmd.AddCommand(ReplaceUISchemasCmd)
}

var DeleteUISchemasid string

func NewDeleteUISchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deletes",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UISchemaAPI.DeleteUISchemas(apiClient.GetConfig().Context, DeleteUISchemasid)

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

	cmd.Flags().StringVarP(&DeleteUISchemasid, "id", "", "", "")
	cmd.MarkFlagRequired("id")

	return cmd
}

func init() {
	DeleteUISchemasCmd := NewDeleteUISchemasCmd()
	UISchemaCmd.AddCommand(DeleteUISchemasCmd)
}
