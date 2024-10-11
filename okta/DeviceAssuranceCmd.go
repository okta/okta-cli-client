package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var DeviceAssuranceCmd = &cobra.Command{
	Use:  "deviceAssurance",
	Long: "Manage DeviceAssuranceAPI",
}

func init() {
	rootCmd.AddCommand(DeviceAssuranceCmd)
}

var CreateDeviceAssurancePolicydata string

func NewCreateDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.CreateDeviceAssurancePolicy(apiClient.GetConfig().Context)

			if CreateDeviceAssurancePolicydata != "" {
				req = req.Data(CreateDeviceAssurancePolicydata)
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

	cmd.Flags().StringVarP(&CreateDeviceAssurancePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateDeviceAssurancePolicyCmd := NewCreateDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(CreateDeviceAssurancePolicyCmd)
}

func NewListDeviceAssurancePoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listPolicies",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.ListDeviceAssurancePolicies(apiClient.GetConfig().Context)

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
	ListDeviceAssurancePoliciesCmd := NewListDeviceAssurancePoliciesCmd()
	DeviceAssuranceCmd.AddCommand(ListDeviceAssurancePoliciesCmd)
}

var GetDeviceAssurancePolicydeviceAssuranceId string

func NewGetDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getPolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.GetDeviceAssurancePolicy(apiClient.GetConfig().Context, GetDeviceAssurancePolicydeviceAssuranceId)

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

	cmd.Flags().StringVarP(&GetDeviceAssurancePolicydeviceAssuranceId, "deviceAssuranceId", "", "", "")
	cmd.MarkFlagRequired("deviceAssuranceId")

	return cmd
}

func init() {
	GetDeviceAssurancePolicyCmd := NewGetDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(GetDeviceAssurancePolicyCmd)
}

var (
	ReplaceDeviceAssurancePolicydeviceAssuranceId string

	ReplaceDeviceAssurancePolicydata string
)

func NewReplaceDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replacePolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.ReplaceDeviceAssurancePolicy(apiClient.GetConfig().Context, ReplaceDeviceAssurancePolicydeviceAssuranceId)

			if ReplaceDeviceAssurancePolicydata != "" {
				req = req.Data(ReplaceDeviceAssurancePolicydata)
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

	cmd.Flags().StringVarP(&ReplaceDeviceAssurancePolicydeviceAssuranceId, "deviceAssuranceId", "", "", "")
	cmd.MarkFlagRequired("deviceAssuranceId")

	cmd.Flags().StringVarP(&ReplaceDeviceAssurancePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceDeviceAssurancePolicyCmd := NewReplaceDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(ReplaceDeviceAssurancePolicyCmd)
}

var DeleteDeviceAssurancePolicydeviceAssuranceId string

func NewDeleteDeviceAssurancePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deletePolicy",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.DeviceAssuranceAPI.DeleteDeviceAssurancePolicy(apiClient.GetConfig().Context, DeleteDeviceAssurancePolicydeviceAssuranceId)

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

	cmd.Flags().StringVarP(&DeleteDeviceAssurancePolicydeviceAssuranceId, "deviceAssuranceId", "", "", "")
	cmd.MarkFlagRequired("deviceAssuranceId")

	return cmd
}

func init() {
	DeleteDeviceAssurancePolicyCmd := NewDeleteDeviceAssurancePolicyCmd()
	DeviceAssuranceCmd.AddCommand(DeleteDeviceAssurancePolicyCmd)
}
