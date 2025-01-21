package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApiServiceIntegrationsCmd = &cobra.Command{
	Use:  "apiServiceIntegrations",
	Long: "Manage ApiServiceIntegrationsAPI",
}

func init() {
	rootCmd.AddCommand(ApiServiceIntegrationsCmd)
}

var CreateApiServiceIntegrationInstancedata string

func NewCreateApiServiceIntegrationInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createApiServiceIntegrationInstance",
		Long: "Create an API Service Integration instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.CreateApiServiceIntegrationInstance(apiClient.GetConfig().Context)

			if CreateApiServiceIntegrationInstancedata != "" {
				req = req.Data(CreateApiServiceIntegrationInstancedata)
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

	cmd.Flags().StringVarP(&CreateApiServiceIntegrationInstancedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateApiServiceIntegrationInstanceCmd := NewCreateApiServiceIntegrationInstanceCmd()
	ApiServiceIntegrationsCmd.AddCommand(CreateApiServiceIntegrationInstanceCmd)
}

func NewListApiServiceIntegrationInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApiServiceIntegrationInstances",
		Long: "List all API Service Integration instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.ListApiServiceIntegrationInstances(apiClient.GetConfig().Context)

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
	ListApiServiceIntegrationInstancesCmd := NewListApiServiceIntegrationInstancesCmd()
	ApiServiceIntegrationsCmd.AddCommand(ListApiServiceIntegrationInstancesCmd)
}

var GetApiServiceIntegrationInstanceapiServiceId string

func NewGetApiServiceIntegrationInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApiServiceIntegrationInstance",
		Long: "Retrieve an API Service Integration instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.GetApiServiceIntegrationInstance(apiClient.GetConfig().Context, GetApiServiceIntegrationInstanceapiServiceId)

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

	cmd.Flags().StringVarP(&GetApiServiceIntegrationInstanceapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	return cmd
}

func init() {
	GetApiServiceIntegrationInstanceCmd := NewGetApiServiceIntegrationInstanceCmd()
	ApiServiceIntegrationsCmd.AddCommand(GetApiServiceIntegrationInstanceCmd)
}

var DeleteApiServiceIntegrationInstanceapiServiceId string

func NewDeleteApiServiceIntegrationInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteApiServiceIntegrationInstance",
		Long: "Delete an API Service Integration instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.DeleteApiServiceIntegrationInstance(apiClient.GetConfig().Context, DeleteApiServiceIntegrationInstanceapiServiceId)

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

	cmd.Flags().StringVarP(&DeleteApiServiceIntegrationInstanceapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	return cmd
}

func init() {
	DeleteApiServiceIntegrationInstanceCmd := NewDeleteApiServiceIntegrationInstanceCmd()
	ApiServiceIntegrationsCmd.AddCommand(DeleteApiServiceIntegrationInstanceCmd)
}

var CreateApiServiceIntegrationInstanceSecretapiServiceId string

func NewCreateApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createApiServiceIntegrationInstanceSecret",
		Long: "Create an API Service Integration instance Secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.CreateApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, CreateApiServiceIntegrationInstanceSecretapiServiceId)

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

	cmd.Flags().StringVarP(&CreateApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	return cmd
}

func init() {
	CreateApiServiceIntegrationInstanceSecretCmd := NewCreateApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(CreateApiServiceIntegrationInstanceSecretCmd)
}

var ListApiServiceIntegrationInstanceSecretsapiServiceId string

func NewListApiServiceIntegrationInstanceSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listApiServiceIntegrationInstanceSecrets",
		Long: "List all API Service Integration instance Secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.ListApiServiceIntegrationInstanceSecrets(apiClient.GetConfig().Context, ListApiServiceIntegrationInstanceSecretsapiServiceId)

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

	cmd.Flags().StringVarP(&ListApiServiceIntegrationInstanceSecretsapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	return cmd
}

func init() {
	ListApiServiceIntegrationInstanceSecretsCmd := NewListApiServiceIntegrationInstanceSecretsCmd()
	ApiServiceIntegrationsCmd.AddCommand(ListApiServiceIntegrationInstanceSecretsCmd)
}

var (
	DeleteApiServiceIntegrationInstanceSecretapiServiceId string

	DeleteApiServiceIntegrationInstanceSecretsecretId string
)

func NewDeleteApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteApiServiceIntegrationInstanceSecret",
		Long: "Delete an API Service Integration instance Secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.DeleteApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, DeleteApiServiceIntegrationInstanceSecretapiServiceId, DeleteApiServiceIntegrationInstanceSecretsecretId)

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

	cmd.Flags().StringVarP(&DeleteApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&DeleteApiServiceIntegrationInstanceSecretsecretId, "secretId", "", "", "")
	cmd.MarkFlagRequired("secretId")

	return cmd
}

func init() {
	DeleteApiServiceIntegrationInstanceSecretCmd := NewDeleteApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(DeleteApiServiceIntegrationInstanceSecretCmd)
}

var (
	ActivateApiServiceIntegrationInstanceSecretapiServiceId string

	ActivateApiServiceIntegrationInstanceSecretsecretId string
)

func NewActivateApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateApiServiceIntegrationInstanceSecret",
		Long: "Activate an API Service Integration instance Secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.ActivateApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, ActivateApiServiceIntegrationInstanceSecretapiServiceId, ActivateApiServiceIntegrationInstanceSecretsecretId)

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

	cmd.Flags().StringVarP(&ActivateApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&ActivateApiServiceIntegrationInstanceSecretsecretId, "secretId", "", "", "")
	cmd.MarkFlagRequired("secretId")

	return cmd
}

func init() {
	ActivateApiServiceIntegrationInstanceSecretCmd := NewActivateApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(ActivateApiServiceIntegrationInstanceSecretCmd)
}

var (
	DeactivateApiServiceIntegrationInstanceSecretapiServiceId string

	DeactivateApiServiceIntegrationInstanceSecretsecretId string
)

func NewDeactivateApiServiceIntegrationInstanceSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateApiServiceIntegrationInstanceSecret",
		Long: "Deactivate an API Service Integration instance Secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApiServiceIntegrationsAPI.DeactivateApiServiceIntegrationInstanceSecret(apiClient.GetConfig().Context, DeactivateApiServiceIntegrationInstanceSecretapiServiceId, DeactivateApiServiceIntegrationInstanceSecretsecretId)

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

	cmd.Flags().StringVarP(&DeactivateApiServiceIntegrationInstanceSecretapiServiceId, "apiServiceId", "", "", "")
	cmd.MarkFlagRequired("apiServiceId")

	cmd.Flags().StringVarP(&DeactivateApiServiceIntegrationInstanceSecretsecretId, "secretId", "", "", "")
	cmd.MarkFlagRequired("secretId")

	return cmd
}

func init() {
	DeactivateApiServiceIntegrationInstanceSecretCmd := NewDeactivateApiServiceIntegrationInstanceSecretCmd()
	ApiServiceIntegrationsCmd.AddCommand(DeactivateApiServiceIntegrationInstanceSecretCmd)
}
