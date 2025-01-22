package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var SessionCmd = &cobra.Command{
	Use:  "session",
	Long: "Manage SessionAPI",
}

func init() {
	rootCmd.AddCommand(SessionCmd)
}

var CreateSessiondata string

func NewCreateSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Session with session token",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.CreateSession(apiClient.GetConfig().Context)

			if CreateSessiondata != "" {
				req = req.Data(CreateSessiondata)
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

	cmd.Flags().StringVarP(&CreateSessiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateSessionCmd := NewCreateSessionCmd()
	SessionCmd.AddCommand(CreateSessionCmd)
}

func NewGetCurrentSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getCurrent",
		Long: "Retrieve the current Session",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.GetCurrentSession(apiClient.GetConfig().Context)

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
	GetCurrentSessionCmd := NewGetCurrentSessionCmd()
	SessionCmd.AddCommand(GetCurrentSessionCmd)
}

func NewCloseCurrentSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "closeCurrent",
		Long: "Close the current Session",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.CloseCurrentSession(apiClient.GetConfig().Context)

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
	CloseCurrentSessionCmd := NewCloseCurrentSessionCmd()
	SessionCmd.AddCommand(CloseCurrentSessionCmd)
}

func NewRefreshCurrentSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "refreshCurrent",
		Long: "Refresh the current Session",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.RefreshCurrentSession(apiClient.GetConfig().Context)

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
	RefreshCurrentSessionCmd := NewRefreshCurrentSessionCmd()
	SessionCmd.AddCommand(RefreshCurrentSessionCmd)
}

var GetSessionsessionId string

func NewGetSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Session",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.GetSession(apiClient.GetConfig().Context, GetSessionsessionId)

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

	cmd.Flags().StringVarP(&GetSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	return cmd
}

func init() {
	GetSessionCmd := NewGetSessionCmd()
	SessionCmd.AddCommand(GetSessionCmd)
}

var RevokeSessionsessionId string

func NewRevokeSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revoke",
		Long: "Revoke a Session",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.RevokeSession(apiClient.GetConfig().Context, RevokeSessionsessionId)

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

	cmd.Flags().StringVarP(&RevokeSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	return cmd
}

func init() {
	RevokeSessionCmd := NewRevokeSessionCmd()
	SessionCmd.AddCommand(RevokeSessionCmd)
}

var RefreshSessionsessionId string

func NewRefreshSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "refresh",
		Long: "Refresh a Session",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SessionAPI.RefreshSession(apiClient.GetConfig().Context, RefreshSessionsessionId)

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

	cmd.Flags().StringVarP(&RefreshSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	return cmd
}

func init() {
	RefreshSessionCmd := NewRefreshSessionCmd()
	SessionCmd.AddCommand(RefreshSessionCmd)
}
