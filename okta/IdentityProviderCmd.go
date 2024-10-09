package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var IdentityProviderCmd = &cobra.Command{
	Use:  "identityProvider",
	Long: "Manage IdentityProviderAPI",
}

func init() {
	rootCmd.AddCommand(IdentityProviderCmd)
}

var CreateIdentityProviderdata string

func NewCreateIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.CreateIdentityProvider(apiClient.GetConfig().Context)

			if CreateIdentityProviderdata != "" {
				req = req.Data(CreateIdentityProviderdata)
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

	cmd.Flags().StringVarP(&CreateIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateIdentityProviderCmd := NewCreateIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(CreateIdentityProviderCmd)
}

func NewListIdentityProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviders(apiClient.GetConfig().Context)

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
	ListIdentityProvidersCmd := NewListIdentityProvidersCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProvidersCmd)
}

var CreateIdentityProviderKeydata string

func NewCreateIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createKey",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.CreateIdentityProviderKey(apiClient.GetConfig().Context)

			if CreateIdentityProviderKeydata != "" {
				req = req.Data(CreateIdentityProviderKeydata)
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

	cmd.Flags().StringVarP(&CreateIdentityProviderKeydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateIdentityProviderKeyCmd := NewCreateIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(CreateIdentityProviderKeyCmd)
}

func NewListIdentityProviderKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listKeys",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviderKeys(apiClient.GetConfig().Context)

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
	ListIdentityProviderKeysCmd := NewListIdentityProviderKeysCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProviderKeysCmd)
}

var GetIdentityProviderKeyidpKeyId string

func NewGetIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getKey",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProviderKey(apiClient.GetConfig().Context, GetIdentityProviderKeyidpKeyId)

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

	cmd.Flags().StringVarP(&GetIdentityProviderKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	return cmd
}

func init() {
	GetIdentityProviderKeyCmd := NewGetIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderKeyCmd)
}

var DeleteIdentityProviderKeyidpKeyId string

func NewDeleteIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteKey",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.DeleteIdentityProviderKey(apiClient.GetConfig().Context, DeleteIdentityProviderKeyidpKeyId)

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

	cmd.Flags().StringVarP(&DeleteIdentityProviderKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	return cmd
}

func init() {
	DeleteIdentityProviderKeyCmd := NewDeleteIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(DeleteIdentityProviderKeyCmd)
}

var GetIdentityProvideridpId string

func NewGetIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProvider(apiClient.GetConfig().Context, GetIdentityProvideridpId)

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

	cmd.Flags().StringVarP(&GetIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	GetIdentityProviderCmd := NewGetIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderCmd)
}

var (
	ReplaceIdentityProvideridpId string

	ReplaceIdentityProviderdata string
)

func NewReplaceIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ReplaceIdentityProvider(apiClient.GetConfig().Context, ReplaceIdentityProvideridpId)

			if ReplaceIdentityProviderdata != "" {
				req = req.Data(ReplaceIdentityProviderdata)
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

	cmd.Flags().StringVarP(&ReplaceIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&ReplaceIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceIdentityProviderCmd := NewReplaceIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(ReplaceIdentityProviderCmd)
}

var DeleteIdentityProvideridpId string

func NewDeleteIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.DeleteIdentityProvider(apiClient.GetConfig().Context, DeleteIdentityProvideridpId)

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

	cmd.Flags().StringVarP(&DeleteIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	DeleteIdentityProviderCmd := NewDeleteIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(DeleteIdentityProviderCmd)
}

var (
	GenerateCsrForIdentityProvideridpId string

	GenerateCsrForIdentityProviderdata string
)

func NewGenerateCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "generateCsrFor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GenerateCsrForIdentityProvider(apiClient.GetConfig().Context, GenerateCsrForIdentityProvideridpId)

			if GenerateCsrForIdentityProviderdata != "" {
				req = req.Data(GenerateCsrForIdentityProviderdata)
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

	cmd.Flags().StringVarP(&GenerateCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GenerateCsrForIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	GenerateCsrForIdentityProviderCmd := NewGenerateCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(GenerateCsrForIdentityProviderCmd)
}

var ListCsrsForIdentityProvideridpId string

func NewListCsrsForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listCsrsFor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListCsrsForIdentityProvider(apiClient.GetConfig().Context, ListCsrsForIdentityProvideridpId)

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

	cmd.Flags().StringVarP(&ListCsrsForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	ListCsrsForIdentityProviderCmd := NewListCsrsForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(ListCsrsForIdentityProviderCmd)
}

var (
	GetCsrForIdentityProvideridpId string

	GetCsrForIdentityProvideridpCsrId string
)

func NewGetCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getCsrFor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetCsrForIdentityProvider(apiClient.GetConfig().Context, GetCsrForIdentityProvideridpId, GetCsrForIdentityProvideridpCsrId)

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

	cmd.Flags().StringVarP(&GetCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GetCsrForIdentityProvideridpCsrId, "idpCsrId", "", "", "")
	cmd.MarkFlagRequired("idpCsrId")

	return cmd
}

func init() {
	GetCsrForIdentityProviderCmd := NewGetCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(GetCsrForIdentityProviderCmd)
}

var (
	RevokeCsrForIdentityProvideridpId string

	RevokeCsrForIdentityProvideridpCsrId string
)

func NewRevokeCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "revokeCsrFor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.RevokeCsrForIdentityProvider(apiClient.GetConfig().Context, RevokeCsrForIdentityProvideridpId, RevokeCsrForIdentityProvideridpCsrId)

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

	cmd.Flags().StringVarP(&RevokeCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&RevokeCsrForIdentityProvideridpCsrId, "idpCsrId", "", "", "")
	cmd.MarkFlagRequired("idpCsrId")

	return cmd
}

func init() {
	RevokeCsrForIdentityProviderCmd := NewRevokeCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(RevokeCsrForIdentityProviderCmd)
}

var (
	PublishCsrForIdentityProvideridpId string

	PublishCsrForIdentityProvideridpCsrId string

	PublishCsrForIdentityProviderdata string
)

func NewPublishCsrForIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "publishCsrFor",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.PublishCsrForIdentityProvider(apiClient.GetConfig().Context, PublishCsrForIdentityProvideridpId, PublishCsrForIdentityProvideridpCsrId)

			if PublishCsrForIdentityProviderdata != "" {
				req = req.Data(PublishCsrForIdentityProviderdata)
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

	cmd.Flags().StringVarP(&PublishCsrForIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&PublishCsrForIdentityProvideridpCsrId, "idpCsrId", "", "", "")
	cmd.MarkFlagRequired("idpCsrId")

	cmd.Flags().StringVarP(&PublishCsrForIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	PublishCsrForIdentityProviderCmd := NewPublishCsrForIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(PublishCsrForIdentityProviderCmd)
}

var ListIdentityProviderSigningKeysidpId string

func NewListIdentityProviderSigningKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listSigningKeys",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviderSigningKeys(apiClient.GetConfig().Context, ListIdentityProviderSigningKeysidpId)

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

	cmd.Flags().StringVarP(&ListIdentityProviderSigningKeysidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	ListIdentityProviderSigningKeysCmd := NewListIdentityProviderSigningKeysCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProviderSigningKeysCmd)
}

var GenerateIdentityProviderSigningKeyidpId string

func NewGenerateIdentityProviderSigningKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "generateSigningKey",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GenerateIdentityProviderSigningKey(apiClient.GetConfig().Context, GenerateIdentityProviderSigningKeyidpId)

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

	cmd.Flags().StringVarP(&GenerateIdentityProviderSigningKeyidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	GenerateIdentityProviderSigningKeyCmd := NewGenerateIdentityProviderSigningKeyCmd()
	IdentityProviderCmd.AddCommand(GenerateIdentityProviderSigningKeyCmd)
}

var (
	GetIdentityProviderSigningKeyidpId string

	GetIdentityProviderSigningKeyidpKeyId string
)

func NewGetIdentityProviderSigningKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getSigningKey",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProviderSigningKey(apiClient.GetConfig().Context, GetIdentityProviderSigningKeyidpId, GetIdentityProviderSigningKeyidpKeyId)

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

	cmd.Flags().StringVarP(&GetIdentityProviderSigningKeyidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GetIdentityProviderSigningKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	return cmd
}

func init() {
	GetIdentityProviderSigningKeyCmd := NewGetIdentityProviderSigningKeyCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderSigningKeyCmd)
}

var (
	CloneIdentityProviderKeyidpId string

	CloneIdentityProviderKeyidpKeyId string
)

func NewCloneIdentityProviderKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cloneKey",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.CloneIdentityProviderKey(apiClient.GetConfig().Context, CloneIdentityProviderKeyidpId, CloneIdentityProviderKeyidpKeyId)

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

	cmd.Flags().StringVarP(&CloneIdentityProviderKeyidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&CloneIdentityProviderKeyidpKeyId, "idpKeyId", "", "", "")
	cmd.MarkFlagRequired("idpKeyId")

	return cmd
}

func init() {
	CloneIdentityProviderKeyCmd := NewCloneIdentityProviderKeyCmd()
	IdentityProviderCmd.AddCommand(CloneIdentityProviderKeyCmd)
}

var ActivateIdentityProvideridpId string

func NewActivateIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ActivateIdentityProvider(apiClient.GetConfig().Context, ActivateIdentityProvideridpId)

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

	cmd.Flags().StringVarP(&ActivateIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	ActivateIdentityProviderCmd := NewActivateIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(ActivateIdentityProviderCmd)
}

var DeactivateIdentityProvideridpId string

func NewDeactivateIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.DeactivateIdentityProvider(apiClient.GetConfig().Context, DeactivateIdentityProvideridpId)

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

	cmd.Flags().StringVarP(&DeactivateIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	DeactivateIdentityProviderCmd := NewDeactivateIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(DeactivateIdentityProviderCmd)
}

var ListIdentityProviderApplicationUsersidpId string

func NewListIdentityProviderApplicationUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listApplicationUsers",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListIdentityProviderApplicationUsers(apiClient.GetConfig().Context, ListIdentityProviderApplicationUsersidpId)

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

	cmd.Flags().StringVarP(&ListIdentityProviderApplicationUsersidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	return cmd
}

func init() {
	ListIdentityProviderApplicationUsersCmd := NewListIdentityProviderApplicationUsersCmd()
	IdentityProviderCmd.AddCommand(ListIdentityProviderApplicationUsersCmd)
}

var (
	LinkUserToIdentityProvideridpId string

	LinkUserToIdentityProvideruserId string

	LinkUserToIdentityProviderdata string
)

func NewLinkUserToIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "linkUserTo",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.LinkUserToIdentityProvider(apiClient.GetConfig().Context, LinkUserToIdentityProvideridpId, LinkUserToIdentityProvideruserId)

			if LinkUserToIdentityProviderdata != "" {
				req = req.Data(LinkUserToIdentityProviderdata)
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

	cmd.Flags().StringVarP(&LinkUserToIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&LinkUserToIdentityProvideruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&LinkUserToIdentityProviderdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	LinkUserToIdentityProviderCmd := NewLinkUserToIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(LinkUserToIdentityProviderCmd)
}

var (
	GetIdentityProviderApplicationUseridpId string

	GetIdentityProviderApplicationUseruserId string
)

func NewGetIdentityProviderApplicationUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getApplicationUser",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.GetIdentityProviderApplicationUser(apiClient.GetConfig().Context, GetIdentityProviderApplicationUseridpId, GetIdentityProviderApplicationUseruserId)

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

	cmd.Flags().StringVarP(&GetIdentityProviderApplicationUseridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&GetIdentityProviderApplicationUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	GetIdentityProviderApplicationUserCmd := NewGetIdentityProviderApplicationUserCmd()
	IdentityProviderCmd.AddCommand(GetIdentityProviderApplicationUserCmd)
}

var (
	UnlinkUserFromIdentityProvideridpId string

	UnlinkUserFromIdentityProvideruserId string
)

func NewUnlinkUserFromIdentityProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "unlinkUserFrom",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.UnlinkUserFromIdentityProvider(apiClient.GetConfig().Context, UnlinkUserFromIdentityProvideridpId, UnlinkUserFromIdentityProvideruserId)

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

	cmd.Flags().StringVarP(&UnlinkUserFromIdentityProvideridpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&UnlinkUserFromIdentityProvideruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	UnlinkUserFromIdentityProviderCmd := NewUnlinkUserFromIdentityProviderCmd()
	IdentityProviderCmd.AddCommand(UnlinkUserFromIdentityProviderCmd)
}

var (
	ListSocialAuthTokensidpId string

	ListSocialAuthTokensuserId string
)

func NewListSocialAuthTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listSocialAuthTokens",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentityProviderAPI.ListSocialAuthTokens(apiClient.GetConfig().Context, ListSocialAuthTokensidpId, ListSocialAuthTokensuserId)

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

	cmd.Flags().StringVarP(&ListSocialAuthTokensidpId, "idpId", "", "", "")
	cmd.MarkFlagRequired("idpId")

	cmd.Flags().StringVarP(&ListSocialAuthTokensuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListSocialAuthTokensCmd := NewListSocialAuthTokensCmd()
	IdentityProviderCmd.AddCommand(ListSocialAuthTokensCmd)
}
