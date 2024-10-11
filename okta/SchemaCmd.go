package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var SchemaCmd = &cobra.Command{
	Use:  "schema",
	Long: "Manage SchemaAPI",
}

func init() {
	rootCmd.AddCommand(SchemaCmd)
}

var (
	UpdateApplicationUserProfileappId string

	UpdateApplicationUserProfiledata string
)

func NewUpdateApplicationUserProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateApplicationUserProfile",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.UpdateApplicationUserProfile(apiClient.GetConfig().Context, UpdateApplicationUserProfileappId)

			if UpdateApplicationUserProfiledata != "" {
				req = req.Data(UpdateApplicationUserProfiledata)
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

	cmd.Flags().StringVarP(&UpdateApplicationUserProfileappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateApplicationUserProfiledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateApplicationUserProfileCmd := NewUpdateApplicationUserProfileCmd()
	SchemaCmd.AddCommand(UpdateApplicationUserProfileCmd)
}

var GetApplicationUserSchemaappId string

func NewGetApplicationUserSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getApplicationUser",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetApplicationUserSchema(apiClient.GetConfig().Context, GetApplicationUserSchemaappId)

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

	cmd.Flags().StringVarP(&GetApplicationUserSchemaappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	GetApplicationUserSchemaCmd := NewGetApplicationUserSchemaCmd()
	SchemaCmd.AddCommand(GetApplicationUserSchemaCmd)
}

var UpdateGroupSchemadata string

func NewUpdateGroupSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateGroup",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.UpdateGroupSchema(apiClient.GetConfig().Context)

			if UpdateGroupSchemadata != "" {
				req = req.Data(UpdateGroupSchemadata)
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

	cmd.Flags().StringVarP(&UpdateGroupSchemadata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateGroupSchemaCmd := NewUpdateGroupSchemaCmd()
	SchemaCmd.AddCommand(UpdateGroupSchemaCmd)
}

func NewGetGroupSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getGroup",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetGroupSchema(apiClient.GetConfig().Context)

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
	GetGroupSchemaCmd := NewGetGroupSchemaCmd()
	SchemaCmd.AddCommand(GetGroupSchemaCmd)
}

func NewListLogStreamSchemasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listLogStreams",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.ListLogStreamSchemas(apiClient.GetConfig().Context)

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
	ListLogStreamSchemasCmd := NewListLogStreamSchemasCmd()
	SchemaCmd.AddCommand(ListLogStreamSchemasCmd)
}

var GetLogStreamSchemalogStreamType string

func NewGetLogStreamSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getLogStream",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetLogStreamSchema(apiClient.GetConfig().Context, GetLogStreamSchemalogStreamType)

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

	cmd.Flags().StringVarP(&GetLogStreamSchemalogStreamType, "logStreamType", "", "", "")
	cmd.MarkFlagRequired("logStreamType")

	return cmd
}

func init() {
	GetLogStreamSchemaCmd := NewGetLogStreamSchemaCmd()
	SchemaCmd.AddCommand(GetLogStreamSchemaCmd)
}

var (
	UpdateUserProfileschemaId string

	UpdateUserProfiledata string
)

func NewUpdateUserProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateUserProfile",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.UpdateUserProfile(apiClient.GetConfig().Context, UpdateUserProfileschemaId)

			if UpdateUserProfiledata != "" {
				req = req.Data(UpdateUserProfiledata)
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

	cmd.Flags().StringVarP(&UpdateUserProfileschemaId, "schemaId", "", "", "")
	cmd.MarkFlagRequired("schemaId")

	cmd.Flags().StringVarP(&UpdateUserProfiledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateUserProfileCmd := NewUpdateUserProfileCmd()
	SchemaCmd.AddCommand(UpdateUserProfileCmd)
}

var GetUserSchemaschemaId string

func NewGetUserSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getUser",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SchemaAPI.GetUserSchema(apiClient.GetConfig().Context, GetUserSchemaschemaId)

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

	cmd.Flags().StringVarP(&GetUserSchemaschemaId, "schemaId", "", "", "")
	cmd.MarkFlagRequired("schemaId")

	return cmd
}

func init() {
	GetUserSchemaCmd := NewGetUserSchemaCmd()
	SchemaCmd.AddCommand(GetUserSchemaCmd)
}
