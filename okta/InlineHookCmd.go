package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var InlineHookCmd = &cobra.Command{
	Use:  "inlineHook",
	Long: "Manage InlineHookAPI",
}

func init() {
	rootCmd.AddCommand(InlineHookCmd)
}

var CreateInlineHookdata string

func NewCreateInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.CreateInlineHook(apiClient.GetConfig().Context)

			if CreateInlineHookdata != "" {
				req = req.Data(CreateInlineHookdata)
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

	cmd.Flags().StringVarP(&CreateInlineHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateInlineHookCmd := NewCreateInlineHookCmd()
	InlineHookCmd.AddCommand(CreateInlineHookCmd)
}

func NewListInlineHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ListInlineHooks(apiClient.GetConfig().Context)

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
	ListInlineHooksCmd := NewListInlineHooksCmd()
	InlineHookCmd.AddCommand(ListInlineHooksCmd)
}

var GetInlineHookinlineHookId string

func NewGetInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.GetInlineHook(apiClient.GetConfig().Context, GetInlineHookinlineHookId)

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

	cmd.Flags().StringVarP(&GetInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	return cmd
}

func init() {
	GetInlineHookCmd := NewGetInlineHookCmd()
	InlineHookCmd.AddCommand(GetInlineHookCmd)
}

var (
	ReplaceInlineHookinlineHookId string

	ReplaceInlineHookdata string
)

func NewReplaceInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ReplaceInlineHook(apiClient.GetConfig().Context, ReplaceInlineHookinlineHookId)

			if ReplaceInlineHookdata != "" {
				req = req.Data(ReplaceInlineHookdata)
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

	cmd.Flags().StringVarP(&ReplaceInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().StringVarP(&ReplaceInlineHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceInlineHookCmd := NewReplaceInlineHookCmd()
	InlineHookCmd.AddCommand(ReplaceInlineHookCmd)
}

var DeleteInlineHookinlineHookId string

func NewDeleteInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.DeleteInlineHook(apiClient.GetConfig().Context, DeleteInlineHookinlineHookId)

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

	cmd.Flags().StringVarP(&DeleteInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	return cmd
}

func init() {
	DeleteInlineHookCmd := NewDeleteInlineHookCmd()
	InlineHookCmd.AddCommand(DeleteInlineHookCmd)
}

var (
	ExecuteInlineHookinlineHookId string

	ExecuteInlineHookdata string
)

func NewExecuteInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "execute",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ExecuteInlineHook(apiClient.GetConfig().Context, ExecuteInlineHookinlineHookId)

			if ExecuteInlineHookdata != "" {
				req = req.Data(ExecuteInlineHookdata)
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

	cmd.Flags().StringVarP(&ExecuteInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	cmd.Flags().StringVarP(&ExecuteInlineHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ExecuteInlineHookCmd := NewExecuteInlineHookCmd()
	InlineHookCmd.AddCommand(ExecuteInlineHookCmd)
}

var ActivateInlineHookinlineHookId string

func NewActivateInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.ActivateInlineHook(apiClient.GetConfig().Context, ActivateInlineHookinlineHookId)

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

	cmd.Flags().StringVarP(&ActivateInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	return cmd
}

func init() {
	ActivateInlineHookCmd := NewActivateInlineHookCmd()
	InlineHookCmd.AddCommand(ActivateInlineHookCmd)
}

var DeactivateInlineHookinlineHookId string

func NewDeactivateInlineHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.InlineHookAPI.DeactivateInlineHook(apiClient.GetConfig().Context, DeactivateInlineHookinlineHookId)

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

	cmd.Flags().StringVarP(&DeactivateInlineHookinlineHookId, "inlineHookId", "", "", "")
	cmd.MarkFlagRequired("inlineHookId")

	return cmd
}

func init() {
	DeactivateInlineHookCmd := NewDeactivateInlineHookCmd()
	InlineHookCmd.AddCommand(DeactivateInlineHookCmd)
}
