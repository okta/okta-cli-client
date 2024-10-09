package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var PushProviderCmd = &cobra.Command{
	Use:  "pushProvider",
	Long: "Manage PushProviderAPI",
}

func init() {
	rootCmd.AddCommand(PushProviderCmd)
}

var CreatePushProviderdata string

func NewCreatePushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.CreatePushProvider(apiClient.GetConfig().Context)

			if CreatePushProviderdata != "" {
				req = req.Data(CreatePushProviderdata)
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
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreatePushProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreatePushProviderCmd := NewCreatePushProviderCmd()
	PushProviderCmd.AddCommand(CreatePushProviderCmd)
}

func NewListPushProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.ListPushProviders(apiClient.GetConfig().Context)

			resp, err := req.Execute()
			if err != nil {
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
	ListPushProvidersCmd := NewListPushProvidersCmd()
	PushProviderCmd.AddCommand(ListPushProvidersCmd)
}

var GetPushProviderpushProviderId string

func NewGetPushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.GetPushProvider(apiClient.GetConfig().Context, GetPushProviderpushProviderId)

			resp, err := req.Execute()
			if err != nil {
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

	cmd.Flags().StringVarP(&GetPushProviderpushProviderId, "pushProviderId", "", "", "")
	cmd.MarkFlagRequired("pushProviderId")

	return cmd
}

func init() {
	GetPushProviderCmd := NewGetPushProviderCmd()
	PushProviderCmd.AddCommand(GetPushProviderCmd)
}

var (
	ReplacePushProviderpushProviderId string

	ReplacePushProviderdata string
)

func NewReplacePushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.ReplacePushProvider(apiClient.GetConfig().Context, ReplacePushProviderpushProviderId)

			if ReplacePushProviderdata != "" {
				req = req.Data(ReplacePushProviderdata)
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
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ReplacePushProviderpushProviderId, "pushProviderId", "", "", "")
	cmd.MarkFlagRequired("pushProviderId")

	cmd.Flags().StringVarP(&ReplacePushProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacePushProviderCmd := NewReplacePushProviderCmd()
	PushProviderCmd.AddCommand(ReplacePushProviderCmd)
}

var DeletePushProviderpushProviderId string

func NewDeletePushProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PushProviderAPI.DeletePushProvider(apiClient.GetConfig().Context, DeletePushProviderpushProviderId)

			resp, err := req.Execute()
			if err != nil {
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

	cmd.Flags().StringVarP(&DeletePushProviderpushProviderId, "pushProviderId", "", "", "")
	cmd.MarkFlagRequired("pushProviderId")

	return cmd
}

func init() {
	DeletePushProviderCmd := NewDeletePushProviderCmd()
	PushProviderCmd.AddCommand(DeletePushProviderCmd)
}
