package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var WebAuthnPreregistrationCmd = &cobra.Command{
	Use:  "webAuthnPreregistration",
	Long: "Manage WebAuthnPreregistrationAPI",
}

func init() {
	rootCmd.AddCommand(WebAuthnPreregistrationCmd)
}

var ActivatePreregistrationEnrollmentdata string

func NewActivatePreregistrationEnrollmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activatePreregistrationEnrollment",
		Long: "Activate a Preregistered WebAuthn Factor",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.ActivatePreregistrationEnrollment(apiClient.GetConfig().Context)

			if ActivatePreregistrationEnrollmentdata != "" {
				req = req.Data(ActivatePreregistrationEnrollmentdata)
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

	cmd.Flags().StringVarP(&ActivatePreregistrationEnrollmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivatePreregistrationEnrollmentCmd := NewActivatePreregistrationEnrollmentCmd()
	WebAuthnPreregistrationCmd.AddCommand(ActivatePreregistrationEnrollmentCmd)
}

var EnrollPreregistrationEnrollmentdata string

func NewEnrollPreregistrationEnrollmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "enrollPreregistrationEnrollment",
		Long: "Enroll a Preregistered WebAuthn Factor",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.EnrollPreregistrationEnrollment(apiClient.GetConfig().Context)

			if EnrollPreregistrationEnrollmentdata != "" {
				req = req.Data(EnrollPreregistrationEnrollmentdata)
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

	cmd.Flags().StringVarP(&EnrollPreregistrationEnrollmentdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	EnrollPreregistrationEnrollmentCmd := NewEnrollPreregistrationEnrollmentCmd()
	WebAuthnPreregistrationCmd.AddCommand(EnrollPreregistrationEnrollmentCmd)
}

var GenerateFulfillmentRequestdata string

func NewGenerateFulfillmentRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateFulfillmentRequest",
		Long: "Generate a Fulfillment Request",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.GenerateFulfillmentRequest(apiClient.GetConfig().Context)

			if GenerateFulfillmentRequestdata != "" {
				req = req.Data(GenerateFulfillmentRequestdata)
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

	cmd.Flags().StringVarP(&GenerateFulfillmentRequestdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	GenerateFulfillmentRequestCmd := NewGenerateFulfillmentRequestCmd()
	WebAuthnPreregistrationCmd.AddCommand(GenerateFulfillmentRequestCmd)
}

var ListWebAuthnPreregistrationFactorsuserId string

func NewListWebAuthnPreregistrationFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listFactors",
		Long: "List all WebAuthn Preregistration Factors",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.ListWebAuthnPreregistrationFactors(apiClient.GetConfig().Context, ListWebAuthnPreregistrationFactorsuserId)

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

	cmd.Flags().StringVarP(&ListWebAuthnPreregistrationFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListWebAuthnPreregistrationFactorsCmd := NewListWebAuthnPreregistrationFactorsCmd()
	WebAuthnPreregistrationCmd.AddCommand(ListWebAuthnPreregistrationFactorsCmd)
}

var (
	DeleteWebAuthnPreregistrationFactoruserId string

	DeleteWebAuthnPreregistrationFactorauthenticatorEnrollmentId string
)

func NewDeleteWebAuthnPreregistrationFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteFactor",
		Long: "Delete a WebAuthn Preregistration Factor",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.WebAuthnPreregistrationAPI.DeleteWebAuthnPreregistrationFactor(apiClient.GetConfig().Context, DeleteWebAuthnPreregistrationFactoruserId, DeleteWebAuthnPreregistrationFactorauthenticatorEnrollmentId)

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

	cmd.Flags().StringVarP(&DeleteWebAuthnPreregistrationFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&DeleteWebAuthnPreregistrationFactorauthenticatorEnrollmentId, "authenticatorEnrollmentId", "", "", "")
	cmd.MarkFlagRequired("authenticatorEnrollmentId")

	return cmd
}

func init() {
	DeleteWebAuthnPreregistrationFactorCmd := NewDeleteWebAuthnPreregistrationFactorCmd()
	WebAuthnPreregistrationCmd.AddCommand(DeleteWebAuthnPreregistrationFactorCmd)
}
