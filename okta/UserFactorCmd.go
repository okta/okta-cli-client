package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var UserFactorCmd = &cobra.Command{
	Use:  "userFactor",
	Long: "Manage UserFactorAPI",
}

func init() {
	rootCmd.AddCommand(UserFactorCmd)
}

var (
	EnrollFactoruserId string

	EnrollFactordata string
)

func NewEnrollFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "enrollFactor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.EnrollFactor(apiClient.GetConfig().Context, EnrollFactoruserId)

			if EnrollFactordata != "" {
				req = req.Data(EnrollFactordata)
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

	cmd.Flags().StringVarP(&EnrollFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&EnrollFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	EnrollFactorCmd := NewEnrollFactorCmd()
	UserFactorCmd.AddCommand(EnrollFactorCmd)
}

var ListFactorsuserId string

func NewListFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listFactors",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ListFactors(apiClient.GetConfig().Context, ListFactorsuserId)

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

	cmd.Flags().StringVarP(&ListFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListFactorsCmd := NewListFactorsCmd()
	UserFactorCmd.AddCommand(ListFactorsCmd)
}

var ListSupportedFactorsuserId string

func NewListSupportedFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listSupportedFactors",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ListSupportedFactors(apiClient.GetConfig().Context, ListSupportedFactorsuserId)

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

	cmd.Flags().StringVarP(&ListSupportedFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListSupportedFactorsCmd := NewListSupportedFactorsCmd()
	UserFactorCmd.AddCommand(ListSupportedFactorsCmd)
}

var ListSupportedSecurityQuestionsuserId string

func NewListSupportedSecurityQuestionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listSupportedSecurityQuestions",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ListSupportedSecurityQuestions(apiClient.GetConfig().Context, ListSupportedSecurityQuestionsuserId)

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

	cmd.Flags().StringVarP(&ListSupportedSecurityQuestionsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListSupportedSecurityQuestionsCmd := NewListSupportedSecurityQuestionsCmd()
	UserFactorCmd.AddCommand(ListSupportedSecurityQuestionsCmd)
}

var (
	GetFactoruserId string

	GetFactorfactorId string
)

func NewGetFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getFactor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.GetFactor(apiClient.GetConfig().Context, GetFactoruserId, GetFactorfactorId)

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

	cmd.Flags().StringVarP(&GetFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	return cmd
}

func init() {
	GetFactorCmd := NewGetFactorCmd()
	UserFactorCmd.AddCommand(GetFactorCmd)
}

var (
	UnenrollFactoruserId string

	UnenrollFactorfactorId string
)

func NewUnenrollFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "unenrollFactor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.UnenrollFactor(apiClient.GetConfig().Context, UnenrollFactoruserId, UnenrollFactorfactorId)

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

	cmd.Flags().StringVarP(&UnenrollFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UnenrollFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	return cmd
}

func init() {
	UnenrollFactorCmd := NewUnenrollFactorCmd()
	UserFactorCmd.AddCommand(UnenrollFactorCmd)
}

var (
	ActivateFactoruserId string

	ActivateFactorfactorId string

	ActivateFactordata string
)

func NewActivateFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activateFactor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ActivateFactor(apiClient.GetConfig().Context, ActivateFactoruserId, ActivateFactorfactorId)

			if ActivateFactordata != "" {
				req = req.Data(ActivateFactordata)
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

	cmd.Flags().StringVarP(&ActivateFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ActivateFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&ActivateFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivateFactorCmd := NewActivateFactorCmd()
	UserFactorCmd.AddCommand(ActivateFactorCmd)
}

var (
	ResendEnrollFactoruserId string

	ResendEnrollFactorfactorId string

	ResendEnrollFactordata string
)

func NewResendEnrollFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "resendEnrollFactor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.ResendEnrollFactor(apiClient.GetConfig().Context, ResendEnrollFactoruserId, ResendEnrollFactorfactorId)

			if ResendEnrollFactordata != "" {
				req = req.Data(ResendEnrollFactordata)
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

	cmd.Flags().StringVarP(&ResendEnrollFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ResendEnrollFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&ResendEnrollFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ResendEnrollFactorCmd := NewResendEnrollFactorCmd()
	UserFactorCmd.AddCommand(ResendEnrollFactorCmd)
}

var (
	GetFactorTransactionStatususerId string

	GetFactorTransactionStatusfactorId string

	GetFactorTransactionStatustransactionId string
)

func NewGetFactorTransactionStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getFactorTransactionStatus",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.GetFactorTransactionStatus(apiClient.GetConfig().Context, GetFactorTransactionStatususerId, GetFactorTransactionStatusfactorId, GetFactorTransactionStatustransactionId)

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

	cmd.Flags().StringVarP(&GetFactorTransactionStatususerId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetFactorTransactionStatusfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&GetFactorTransactionStatustransactionId, "transactionId", "", "", "")
	cmd.MarkFlagRequired("transactionId")

	return cmd
}

func init() {
	GetFactorTransactionStatusCmd := NewGetFactorTransactionStatusCmd()
	UserFactorCmd.AddCommand(GetFactorTransactionStatusCmd)
}

var (
	VerifyFactoruserId string

	VerifyFactorfactorId string

	VerifyFactordata string
)

func NewVerifyFactorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "verifyFactor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserFactorAPI.VerifyFactor(apiClient.GetConfig().Context, VerifyFactoruserId, VerifyFactorfactorId)

			if VerifyFactordata != "" {
				req = req.Data(VerifyFactordata)
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

	cmd.Flags().StringVarP(&VerifyFactoruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&VerifyFactorfactorId, "factorId", "", "", "")
	cmd.MarkFlagRequired("factorId")

	cmd.Flags().StringVarP(&VerifyFactordata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	VerifyFactorCmd := NewVerifyFactorCmd()
	UserFactorCmd.AddCommand(VerifyFactorCmd)
}
