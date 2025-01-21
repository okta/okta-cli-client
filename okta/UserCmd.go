package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:  "user",
	Long: "Manage UserAPI",
}

func init() {
	rootCmd.AddCommand(UserCmd)
}

var CreateUserdata string

func NewCreateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.CreateUser(apiClient.GetConfig().Context)

			if CreateUserdata != "" {
				req = req.Data(CreateUserdata)
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

	cmd.Flags().StringVarP(&CreateUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateUserCmd := NewCreateUserCmd()
	UserCmd.AddCommand(CreateUserCmd)
}

func NewListUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Users",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUsers(apiClient.GetConfig().Context)

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
	ListUsersCmd := NewListUsersCmd()
	UserCmd.AddCommand(ListUsersCmd)
}

var (
	UpdateUseruserId string

	UpdateUserdata string
)

func NewUpdateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: "Update a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.UpdateUser(apiClient.GetConfig().Context, UpdateUseruserId)

			if UpdateUserdata != "" {
				req = req.Data(UpdateUserdata)
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

	cmd.Flags().StringVarP(&UpdateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&UpdateUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateUserCmd := NewUpdateUserCmd()
	UserCmd.AddCommand(UpdateUserCmd)
}

var GetUseruserId string

func NewGetUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GetUser(apiClient.GetConfig().Context, GetUseruserId)

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

	cmd.Flags().StringVarP(&GetUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	GetUserCmd := NewGetUserCmd()
	UserCmd.AddCommand(GetUserCmd)
}

var (
	ReplaceUseruserId string

	ReplaceUserdata string
)

func NewReplaceUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ReplaceUser(apiClient.GetConfig().Context, ReplaceUseruserId)

			if ReplaceUserdata != "" {
				req = req.Data(ReplaceUserdata)
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

	cmd.Flags().StringVarP(&ReplaceUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ReplaceUserdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceUserCmd := NewReplaceUserCmd()
	UserCmd.AddCommand(ReplaceUserCmd)
}

var DeleteUseruserId string

func NewDeleteUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.DeleteUser(apiClient.GetConfig().Context, DeleteUseruserId)

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

	cmd.Flags().StringVarP(&DeleteUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	DeleteUserCmd := NewDeleteUserCmd()
	UserCmd.AddCommand(DeleteUserCmd)
}

var ListAppLinksuserId string

func NewListAppLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listAppLinks",
		Long: "List all Assigned Application Links",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListAppLinks(apiClient.GetConfig().Context, ListAppLinksuserId)

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

	cmd.Flags().StringVarP(&ListAppLinksuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListAppLinksCmd := NewListAppLinksCmd()
	UserCmd.AddCommand(ListAppLinksCmd)
}

var ListUserBlocksuserId string

func NewListUserBlocksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listBlocks",
		Long: "List all User Blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserBlocks(apiClient.GetConfig().Context, ListUserBlocksuserId)

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

	cmd.Flags().StringVarP(&ListUserBlocksuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListUserBlocksCmd := NewListUserBlocksCmd()
	UserCmd.AddCommand(ListUserBlocksCmd)
}

var ListUserClientsuserId string

func NewListUserClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listClients",
		Long: "List all Clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserClients(apiClient.GetConfig().Context, ListUserClientsuserId)

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

	cmd.Flags().StringVarP(&ListUserClientsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListUserClientsCmd := NewListUserClientsCmd()
	UserCmd.AddCommand(ListUserClientsCmd)
}

var (
	ListGrantsForUserAndClientuserId string

	ListGrantsForUserAndClientclientId string
)

func NewListGrantsForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGrantsForAndClient",
		Long: "List all Grants for a Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListGrantsForUserAndClient(apiClient.GetConfig().Context, ListGrantsForUserAndClientuserId, ListGrantsForUserAndClientclientId)

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

	cmd.Flags().StringVarP(&ListGrantsForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListGrantsForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	return cmd
}

func init() {
	ListGrantsForUserAndClientCmd := NewListGrantsForUserAndClientCmd()
	UserCmd.AddCommand(ListGrantsForUserAndClientCmd)
}

var (
	RevokeGrantsForUserAndClientuserId string

	RevokeGrantsForUserAndClientclientId string
)

func NewRevokeGrantsForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeGrantsForAndClient",
		Long: "Revoke all Grants for a Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeGrantsForUserAndClient(apiClient.GetConfig().Context, RevokeGrantsForUserAndClientuserId, RevokeGrantsForUserAndClientclientId)

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

	cmd.Flags().StringVarP(&RevokeGrantsForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeGrantsForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	return cmd
}

func init() {
	RevokeGrantsForUserAndClientCmd := NewRevokeGrantsForUserAndClientCmd()
	UserCmd.AddCommand(RevokeGrantsForUserAndClientCmd)
}

var (
	ListRefreshTokensForUserAndClientuserId string

	ListRefreshTokensForUserAndClientclientId string
)

func NewListRefreshTokensForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listRefreshTokensForAndClient",
		Long: "List all Refresh Tokens for a Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListRefreshTokensForUserAndClient(apiClient.GetConfig().Context, ListRefreshTokensForUserAndClientuserId, ListRefreshTokensForUserAndClientclientId)

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

	cmd.Flags().StringVarP(&ListRefreshTokensForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListRefreshTokensForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	return cmd
}

func init() {
	ListRefreshTokensForUserAndClientCmd := NewListRefreshTokensForUserAndClientCmd()
	UserCmd.AddCommand(ListRefreshTokensForUserAndClientCmd)
}

var (
	RevokeTokensForUserAndClientuserId string

	RevokeTokensForUserAndClientclientId string
)

func NewRevokeTokensForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeTokensForAndClient",
		Long: "Revoke all Refresh Tokens for a Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeTokensForUserAndClient(apiClient.GetConfig().Context, RevokeTokensForUserAndClientuserId, RevokeTokensForUserAndClientclientId)

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

	cmd.Flags().StringVarP(&RevokeTokensForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeTokensForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	return cmd
}

func init() {
	RevokeTokensForUserAndClientCmd := NewRevokeTokensForUserAndClientCmd()
	UserCmd.AddCommand(RevokeTokensForUserAndClientCmd)
}

var (
	GetRefreshTokenForUserAndClientuserId string

	GetRefreshTokenForUserAndClientclientId string

	GetRefreshTokenForUserAndClienttokenId string
)

func NewGetRefreshTokenForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getRefreshTokenForAndClient",
		Long: "Retrieve a Refresh Token for a Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GetRefreshTokenForUserAndClient(apiClient.GetConfig().Context, GetRefreshTokenForUserAndClientuserId, GetRefreshTokenForUserAndClientclientId, GetRefreshTokenForUserAndClienttokenId)

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

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().StringVarP(&GetRefreshTokenForUserAndClienttokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	return cmd
}

func init() {
	GetRefreshTokenForUserAndClientCmd := NewGetRefreshTokenForUserAndClientCmd()
	UserCmd.AddCommand(GetRefreshTokenForUserAndClientCmd)
}

var (
	RevokeTokenForUserAndClientuserId string

	RevokeTokenForUserAndClientclientId string

	RevokeTokenForUserAndClienttokenId string
)

func NewRevokeTokenForUserAndClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeTokenForAndClient",
		Long: "Revoke a Token for a Client",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeTokenForUserAndClient(apiClient.GetConfig().Context, RevokeTokenForUserAndClientuserId, RevokeTokenForUserAndClientclientId, RevokeTokenForUserAndClienttokenId)

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

	cmd.Flags().StringVarP(&RevokeTokenForUserAndClientuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeTokenForUserAndClientclientId, "clientId", "", "", "")
	cmd.MarkFlagRequired("clientId")

	cmd.Flags().StringVarP(&RevokeTokenForUserAndClienttokenId, "tokenId", "", "", "")
	cmd.MarkFlagRequired("tokenId")

	return cmd
}

func init() {
	RevokeTokenForUserAndClientCmd := NewRevokeTokenForUserAndClientCmd()
	UserCmd.AddCommand(RevokeTokenForUserAndClientCmd)
}

var (
	ChangePassworduserId string

	ChangePassworddata string
)

func NewChangePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "changePassword",
		Long: "Change Password",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ChangePassword(apiClient.GetConfig().Context, ChangePassworduserId)

			if ChangePassworddata != "" {
				req = req.Data(ChangePassworddata)
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

	cmd.Flags().StringVarP(&ChangePassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ChangePassworddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ChangePasswordCmd := NewChangePasswordCmd()
	UserCmd.AddCommand(ChangePasswordCmd)
}

var (
	ChangeRecoveryQuestionuserId string

	ChangeRecoveryQuestiondata string
)

func NewChangeRecoveryQuestionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "changeRecoveryQuestion",
		Long: "Change Recovery Question",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ChangeRecoveryQuestion(apiClient.GetConfig().Context, ChangeRecoveryQuestionuserId)

			if ChangeRecoveryQuestiondata != "" {
				req = req.Data(ChangeRecoveryQuestiondata)
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

	cmd.Flags().StringVarP(&ChangeRecoveryQuestionuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ChangeRecoveryQuestiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ChangeRecoveryQuestionCmd := NewChangeRecoveryQuestionCmd()
	UserCmd.AddCommand(ChangeRecoveryQuestionCmd)
}

var ForgotPassworduserId string

func NewForgotPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "forgotPassword",
		Long: "Initiate Forgot Password",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ForgotPassword(apiClient.GetConfig().Context, ForgotPassworduserId)

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

	cmd.Flags().StringVarP(&ForgotPassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ForgotPasswordCmd := NewForgotPasswordCmd()
	UserCmd.AddCommand(ForgotPasswordCmd)
}

var (
	ForgotPasswordSetNewPassworduserId string

	ForgotPasswordSetNewPassworddata string
)

func NewForgotPasswordSetNewPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "forgotPasswordSetNewPassword",
		Long: "Reset Password with Recovery Question",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ForgotPasswordSetNewPassword(apiClient.GetConfig().Context, ForgotPasswordSetNewPassworduserId)

			if ForgotPasswordSetNewPassworddata != "" {
				req = req.Data(ForgotPasswordSetNewPassworddata)
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

	cmd.Flags().StringVarP(&ForgotPasswordSetNewPassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ForgotPasswordSetNewPassworddata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ForgotPasswordSetNewPasswordCmd := NewForgotPasswordSetNewPasswordCmd()
	UserCmd.AddCommand(ForgotPasswordSetNewPasswordCmd)
}

var ListUserGrantsuserId string

func NewListUserGrantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGrants",
		Long: "List all User Grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserGrants(apiClient.GetConfig().Context, ListUserGrantsuserId)

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

	cmd.Flags().StringVarP(&ListUserGrantsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListUserGrantsCmd := NewListUserGrantsCmd()
	UserCmd.AddCommand(ListUserGrantsCmd)
}

var RevokeUserGrantsuserId string

func NewRevokeUserGrantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeGrants",
		Long: "Revoke all User Grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeUserGrants(apiClient.GetConfig().Context, RevokeUserGrantsuserId)

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

	cmd.Flags().StringVarP(&RevokeUserGrantsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	RevokeUserGrantsCmd := NewRevokeUserGrantsCmd()
	UserCmd.AddCommand(RevokeUserGrantsCmd)
}

var (
	GetUserGrantuserId string

	GetUserGrantgrantId string
)

func NewGetUserGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getGrant",
		Long: "Retrieve a User Grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GetUserGrant(apiClient.GetConfig().Context, GetUserGrantuserId, GetUserGrantgrantId)

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

	cmd.Flags().StringVarP(&GetUserGrantuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&GetUserGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	return cmd
}

func init() {
	GetUserGrantCmd := NewGetUserGrantCmd()
	UserCmd.AddCommand(GetUserGrantCmd)
}

var (
	RevokeUserGrantuserId string

	RevokeUserGrantgrantId string
)

func NewRevokeUserGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeGrant",
		Long: "Revoke a User Grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeUserGrant(apiClient.GetConfig().Context, RevokeUserGrantuserId, RevokeUserGrantgrantId)

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

	cmd.Flags().StringVarP(&RevokeUserGrantuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&RevokeUserGrantgrantId, "grantId", "", "", "")
	cmd.MarkFlagRequired("grantId")

	return cmd
}

func init() {
	RevokeUserGrantCmd := NewRevokeUserGrantCmd()
	UserCmd.AddCommand(RevokeUserGrantCmd)
}

var ListUserGroupsuserId string

func NewListUserGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listGroups",
		Long: "List all Groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserGroups(apiClient.GetConfig().Context, ListUserGroupsuserId)

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

	cmd.Flags().StringVarP(&ListUserGroupsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListUserGroupsCmd := NewListUserGroupsCmd()
	UserCmd.AddCommand(ListUserGroupsCmd)
}

var ListUserIdentityProvidersuserId string

func NewListUserIdentityProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listIdentityProviders",
		Long: "List all Identity Providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListUserIdentityProviders(apiClient.GetConfig().Context, ListUserIdentityProvidersuserId)

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

	cmd.Flags().StringVarP(&ListUserIdentityProvidersuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ListUserIdentityProvidersCmd := NewListUserIdentityProvidersCmd()
	UserCmd.AddCommand(ListUserIdentityProvidersCmd)
}

var ActivateUseruserId string

func NewActivateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ActivateUser(apiClient.GetConfig().Context, ActivateUseruserId)

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

	cmd.Flags().StringVarP(&ActivateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ActivateUserCmd := NewActivateUserCmd()
	UserCmd.AddCommand(ActivateUserCmd)
}

var DeactivateUseruserId string

func NewDeactivateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.DeactivateUser(apiClient.GetConfig().Context, DeactivateUseruserId)

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

	cmd.Flags().StringVarP(&DeactivateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	DeactivateUserCmd := NewDeactivateUserCmd()
	UserCmd.AddCommand(DeactivateUserCmd)
}

var ExpirePassworduserId string

func NewExpirePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "expirePassword",
		Long: "Expire Password",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ExpirePassword(apiClient.GetConfig().Context, ExpirePassworduserId)

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

	cmd.Flags().StringVarP(&ExpirePassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ExpirePasswordCmd := NewExpirePasswordCmd()
	UserCmd.AddCommand(ExpirePasswordCmd)
}

var ExpirePasswordAndGetTemporaryPassworduserId string

func NewExpirePasswordAndGetTemporaryPasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "expirePasswordAndGetTemporaryPassword",
		Long: "Expire Password and Set Temporary Password",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ExpirePasswordAndGetTemporaryPassword(apiClient.GetConfig().Context, ExpirePasswordAndGetTemporaryPassworduserId)

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

	cmd.Flags().StringVarP(&ExpirePasswordAndGetTemporaryPassworduserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ExpirePasswordAndGetTemporaryPasswordCmd := NewExpirePasswordAndGetTemporaryPasswordCmd()
	UserCmd.AddCommand(ExpirePasswordAndGetTemporaryPasswordCmd)
}

var ReactivateUseruserId string

func NewReactivateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "reactivate",
		Long: "Reactivate a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ReactivateUser(apiClient.GetConfig().Context, ReactivateUseruserId)

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

	cmd.Flags().StringVarP(&ReactivateUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ReactivateUserCmd := NewReactivateUserCmd()
	UserCmd.AddCommand(ReactivateUserCmd)
}

var ResetFactorsuserId string

func NewResetFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "resetFactors",
		Long: "Reset all Factors",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ResetFactors(apiClient.GetConfig().Context, ResetFactorsuserId)

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

	cmd.Flags().StringVarP(&ResetFactorsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	ResetFactorsCmd := NewResetFactorsCmd()
	UserCmd.AddCommand(ResetFactorsCmd)
}

var GenerateResetPasswordTokenuserId string

func NewGenerateResetPasswordTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generateResetPasswordToken",
		Long: "Generate a Reset Password Token",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.GenerateResetPasswordToken(apiClient.GetConfig().Context, GenerateResetPasswordTokenuserId)

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

	cmd.Flags().StringVarP(&GenerateResetPasswordTokenuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	GenerateResetPasswordTokenCmd := NewGenerateResetPasswordTokenCmd()
	UserCmd.AddCommand(GenerateResetPasswordTokenCmd)
}

var SuspendUseruserId string

func NewSuspendUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "suspend",
		Long: "Suspend a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.SuspendUser(apiClient.GetConfig().Context, SuspendUseruserId)

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

	cmd.Flags().StringVarP(&SuspendUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	SuspendUserCmd := NewSuspendUserCmd()
	UserCmd.AddCommand(SuspendUserCmd)
}

var UnlockUseruserId string

func NewUnlockUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unlock",
		Long: "Unlock a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.UnlockUser(apiClient.GetConfig().Context, UnlockUseruserId)

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

	cmd.Flags().StringVarP(&UnlockUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	UnlockUserCmd := NewUnlockUserCmd()
	UserCmd.AddCommand(UnlockUserCmd)
}

var UnsuspendUseruserId string

func NewUnsuspendUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unsuspend",
		Long: "Unsuspend a User",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.UnsuspendUser(apiClient.GetConfig().Context, UnsuspendUseruserId)

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

	cmd.Flags().StringVarP(&UnsuspendUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	UnsuspendUserCmd := NewUnsuspendUserCmd()
	UserCmd.AddCommand(UnsuspendUserCmd)
}

var (
	SetLinkedObjectForUseruserId string

	SetLinkedObjectForUserprimaryRelationshipName string

	SetLinkedObjectForUserprimaryUserId string
)

func NewSetLinkedObjectForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "setLinkedObjectFor",
		Long: "Create a Linked Object for two Users",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.SetLinkedObjectForUser(apiClient.GetConfig().Context, SetLinkedObjectForUseruserId, SetLinkedObjectForUserprimaryRelationshipName, SetLinkedObjectForUserprimaryUserId)

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

	cmd.Flags().StringVarP(&SetLinkedObjectForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&SetLinkedObjectForUserprimaryRelationshipName, "primaryRelationshipName", "", "", "")
	cmd.MarkFlagRequired("primaryRelationshipName")

	cmd.Flags().StringVarP(&SetLinkedObjectForUserprimaryUserId, "primaryUserId", "", "", "")
	cmd.MarkFlagRequired("primaryUserId")

	return cmd
}

func init() {
	SetLinkedObjectForUserCmd := NewSetLinkedObjectForUserCmd()
	UserCmd.AddCommand(SetLinkedObjectForUserCmd)
}

var (
	ListLinkedObjectsForUseruserId string

	ListLinkedObjectsForUserrelationshipName string
)

func NewListLinkedObjectsForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listLinkedObjectsFor",
		Long: "List all Linked Objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.ListLinkedObjectsForUser(apiClient.GetConfig().Context, ListLinkedObjectsForUseruserId, ListLinkedObjectsForUserrelationshipName)

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

	cmd.Flags().StringVarP(&ListLinkedObjectsForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&ListLinkedObjectsForUserrelationshipName, "relationshipName", "", "", "")
	cmd.MarkFlagRequired("relationshipName")

	return cmd
}

func init() {
	ListLinkedObjectsForUserCmd := NewListLinkedObjectsForUserCmd()
	UserCmd.AddCommand(ListLinkedObjectsForUserCmd)
}

var (
	DeleteLinkedObjectForUseruserId string

	DeleteLinkedObjectForUserrelationshipName string
)

func NewDeleteLinkedObjectForUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteLinkedObjectFor",
		Long: "Delete a Linked Object",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.DeleteLinkedObjectForUser(apiClient.GetConfig().Context, DeleteLinkedObjectForUseruserId, DeleteLinkedObjectForUserrelationshipName)

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

	cmd.Flags().StringVarP(&DeleteLinkedObjectForUseruserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	cmd.Flags().StringVarP(&DeleteLinkedObjectForUserrelationshipName, "relationshipName", "", "", "")
	cmd.MarkFlagRequired("relationshipName")

	return cmd
}

func init() {
	DeleteLinkedObjectForUserCmd := NewDeleteLinkedObjectForUserCmd()
	UserCmd.AddCommand(DeleteLinkedObjectForUserCmd)
}

var RevokeUserSessionsuserId string

func NewRevokeUserSessionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revokeSessions",
		Long: "Revoke all User Sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.UserAPI.RevokeUserSessions(apiClient.GetConfig().Context, RevokeUserSessionsuserId)

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

	cmd.Flags().StringVarP(&RevokeUserSessionsuserId, "userId", "", "", "")
	cmd.MarkFlagRequired("userId")

	return cmd
}

func init() {
	RevokeUserSessionsCmd := NewRevokeUserSessionsCmd()
	UserCmd.AddCommand(RevokeUserSessionsCmd)
}
