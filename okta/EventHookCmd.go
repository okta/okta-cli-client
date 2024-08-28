package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var EventHookCmd = &cobra.Command{
	Use:  "eventHook",
	Long: "Manage EventHookAPI",
}

func init() {
	rootCmd.AddCommand(EventHookCmd)
}

var CreateEventHookdata string

func NewCreateEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.CreateEventHook(apiClient.GetConfig().Context)

			if CreateEventHookdata != "" {
				req = req.Data(CreateEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateEventHookCmd := NewCreateEventHookCmd()
	EventHookCmd.AddCommand(CreateEventHookCmd)
}

func NewListEventHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.ListEventHooks(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	return cmd
}

func init() {
	ListEventHooksCmd := NewListEventHooksCmd()
	EventHookCmd.AddCommand(ListEventHooksCmd)
}

var GetEventHookeventHookId string

func NewGetEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.GetEventHook(apiClient.GetConfig().Context, GetEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	return cmd
}

func init() {
	GetEventHookCmd := NewGetEventHookCmd()
	EventHookCmd.AddCommand(GetEventHookCmd)
}

var (
	ReplaceEventHookeventHookId string

	ReplaceEventHookdata string
)

func NewReplaceEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.ReplaceEventHook(apiClient.GetConfig().Context, ReplaceEventHookeventHookId)

			if ReplaceEventHookdata != "" {
				req = req.Data(ReplaceEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplaceEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().StringVarP(&ReplaceEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceEventHookCmd := NewReplaceEventHookCmd()
	EventHookCmd.AddCommand(ReplaceEventHookCmd)
}

var DeleteEventHookeventHookId string

func NewDeleteEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.DeleteEventHook(apiClient.GetConfig().Context, DeleteEventHookeventHookId)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	return cmd
}

func init() {
	DeleteEventHookCmd := NewDeleteEventHookCmd()
	EventHookCmd.AddCommand(DeleteEventHookCmd)
}

var (
	ActivateEventHookeventHookId string

	ActivateEventHookdata string
)

func NewActivateEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.ActivateEventHook(apiClient.GetConfig().Context, ActivateEventHookeventHookId)

			if ActivateEventHookdata != "" {
				req = req.Data(ActivateEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ActivateEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().StringVarP(&ActivateEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivateEventHookCmd := NewActivateEventHookCmd()
	EventHookCmd.AddCommand(ActivateEventHookCmd)
}

var (
	DeactivateEventHookeventHookId string

	DeactivateEventHookdata string
)

func NewDeactivateEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.DeactivateEventHook(apiClient.GetConfig().Context, DeactivateEventHookeventHookId)

			if DeactivateEventHookdata != "" {
				req = req.Data(DeactivateEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeactivateEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().StringVarP(&DeactivateEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	DeactivateEventHookCmd := NewDeactivateEventHookCmd()
	EventHookCmd.AddCommand(DeactivateEventHookCmd)
}

var (
	VerifyEventHookeventHookId string

	VerifyEventHookdata string
)

func NewVerifyEventHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "verify",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.EventHookAPI.VerifyEventHook(apiClient.GetConfig().Context, VerifyEventHookeventHookId)

			if VerifyEventHookdata != "" {
				req = req.Data(VerifyEventHookdata)
			}

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&VerifyEventHookeventHookId, "eventHookId", "", "", "")
	cmd.MarkFlagRequired("eventHookId")

	cmd.Flags().StringVarP(&VerifyEventHookdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	VerifyEventHookCmd := NewVerifyEventHookCmd()
	EventHookCmd.AddCommand(VerifyEventHookCmd)
}
