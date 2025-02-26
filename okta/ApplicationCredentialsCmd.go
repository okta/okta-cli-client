package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationCredentialsCmd = &cobra.Command{
	Use:  "applicationCredentials",
	Long: "Manage ApplicationCredentialsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationCredentialsCmd)
}

var (
	GenerateCsrForApplicationappId string

	GenerateCsrForApplicationdata string
)

func NewGenerateCsrForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateCsrForApplication",
		Long: "Generate a Certificate Signing Request",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GenerateCsrForApplication(apiClient.GetConfig().Context, GenerateCsrForApplicationappId)

			if GenerateCsrForApplicationdata != "" {
				req = req.Data(GenerateCsrForApplicationdata)
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

	cmd.Flags().StringVarP(&GenerateCsrForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GenerateCsrForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	GenerateCsrForApplicationCmd := NewGenerateCsrForApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(GenerateCsrForApplicationCmd)
}

var ListCsrsForApplicationappId string

func NewListCsrsForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listCsrsForApplication",
		Long: "List all Certificate Signing Requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.ListCsrsForApplication(apiClient.GetConfig().Context, ListCsrsForApplicationappId)

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

	cmd.Flags().StringVarP(&ListCsrsForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	ListCsrsForApplicationCmd := NewListCsrsForApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(ListCsrsForApplicationCmd)
}

var (
	GetCsrForApplicationappId string

	GetCsrForApplicationcsrId string
)

func NewGetCsrForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCsrForApplication",
		Long: "Retrieve a Certificate Signing Request",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GetCsrForApplication(apiClient.GetConfig().Context, GetCsrForApplicationappId, GetCsrForApplicationcsrId)

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

	cmd.Flags().StringVarP(&GetCsrForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetCsrForApplicationcsrId, "csrId", "", "", "")
	cmd.MarkFlagRequired("csrId")

	return cmd
}

func init() {
	GetCsrForApplicationCmd := NewGetCsrForApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(GetCsrForApplicationCmd)
}

var (
	RevokeCsrFromApplicationappId string

	RevokeCsrFromApplicationcsrId string
)

func NewRevokeCsrFromApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeCsrFromApplication",
		Long: "Revoke a Certificate Signing Request",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.RevokeCsrFromApplication(apiClient.GetConfig().Context, RevokeCsrFromApplicationappId, RevokeCsrFromApplicationcsrId)

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

	cmd.Flags().StringVarP(&RevokeCsrFromApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&RevokeCsrFromApplicationcsrId, "csrId", "", "", "")
	cmd.MarkFlagRequired("csrId")

	return cmd
}

func init() {
	RevokeCsrFromApplicationCmd := NewRevokeCsrFromApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(RevokeCsrFromApplicationCmd)
}

var (
	PublishCsrFromApplicationappId string

	PublishCsrFromApplicationcsrId string

	PublishCsrFromApplicationdata string
)

func NewPublishCsrFromApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "publishCsrFromApplication",
		Long: "Publish a Certificate Signing Request",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.PublishCsrFromApplication(apiClient.GetConfig().Context, PublishCsrFromApplicationappId, PublishCsrFromApplicationcsrId)

			if PublishCsrFromApplicationdata != "" {
				req = req.Data(PublishCsrFromApplicationdata)
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

	cmd.Flags().StringVarP(&PublishCsrFromApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&PublishCsrFromApplicationcsrId, "csrId", "", "", "")
	cmd.MarkFlagRequired("csrId")

	cmd.Flags().StringVarP(&PublishCsrFromApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	PublishCsrFromApplicationCmd := NewPublishCsrFromApplicationCmd()
	ApplicationCredentialsCmd.AddCommand(PublishCsrFromApplicationCmd)
}

var ListApplicationKeysappId string

func NewListApplicationKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApplicationKeys",
		Long: "List all Key Credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.ListApplicationKeys(apiClient.GetConfig().Context, ListApplicationKeysappId)

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

	cmd.Flags().StringVarP(&ListApplicationKeysappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	ListApplicationKeysCmd := NewListApplicationKeysCmd()
	ApplicationCredentialsCmd.AddCommand(ListApplicationKeysCmd)
}

var GenerateApplicationKeyappId string

func NewGenerateApplicationKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateApplicationKey",
		Long: "Generate a Key Credential",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GenerateApplicationKey(apiClient.GetConfig().Context, GenerateApplicationKeyappId)

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

	cmd.Flags().StringVarP(&GenerateApplicationKeyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	GenerateApplicationKeyCmd := NewGenerateApplicationKeyCmd()
	ApplicationCredentialsCmd.AddCommand(GenerateApplicationKeyCmd)
}

var (
	GetApplicationKeyappId string

	GetApplicationKeykeyId string
)

func NewGetApplicationKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationKey",
		Long: "Retrieve a Key Credential",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.GetApplicationKey(apiClient.GetConfig().Context, GetApplicationKeyappId, GetApplicationKeykeyId)

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

	cmd.Flags().StringVarP(&GetApplicationKeyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetApplicationKeykeyId, "keyId", "", "", "")
	cmd.MarkFlagRequired("keyId")

	return cmd
}

func init() {
	GetApplicationKeyCmd := NewGetApplicationKeyCmd()
	ApplicationCredentialsCmd.AddCommand(GetApplicationKeyCmd)
}

var (
	CloneApplicationKeyappId string

	CloneApplicationKeykeyId string
)

func NewCloneApplicationKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "cloneApplicationKey",
		Long: "Clone a Key Credential",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationCredentialsAPI.CloneApplicationKey(apiClient.GetConfig().Context, CloneApplicationKeyappId, CloneApplicationKeykeyId)

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

	cmd.Flags().StringVarP(&CloneApplicationKeyappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&CloneApplicationKeykeyId, "keyId", "", "", "")
	cmd.MarkFlagRequired("keyId")

	return cmd
}

func init() {
	CloneApplicationKeyCmd := NewCloneApplicationKeyCmd()
	ApplicationCredentialsCmd.AddCommand(CloneApplicationKeyCmd)
}
