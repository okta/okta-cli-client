package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var IdentitySourceCmd = &cobra.Command{
	Use:  "identitySource",
	Long: "Manage IdentitySourceAPI",
}

func init() {
	rootCmd.AddCommand(IdentitySourceCmd)
}

var (
	CreateIdentitySourceSessionidentitySourceId string

	CreateIdentitySourceSessiondata string
)

func NewCreateIdentitySourceSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createSession",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.CreateIdentitySourceSession(apiClient.GetConfig().Context, CreateIdentitySourceSessionidentitySourceId)

			if CreateIdentitySourceSessiondata != "" {
				req = req.Data(CreateIdentitySourceSessiondata)
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
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&CreateIdentitySourceSessionidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&CreateIdentitySourceSessiondata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateIdentitySourceSessionCmd := NewCreateIdentitySourceSessionCmd()
	IdentitySourceCmd.AddCommand(CreateIdentitySourceSessionCmd)
}

var ListIdentitySourceSessionsidentitySourceId string

func NewListIdentitySourceSessionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listSessions",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.ListIdentitySourceSessions(apiClient.GetConfig().Context, ListIdentitySourceSessionsidentitySourceId)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&ListIdentitySourceSessionsidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	return cmd
}

func init() {
	ListIdentitySourceSessionsCmd := NewListIdentitySourceSessionsCmd()
	IdentitySourceCmd.AddCommand(ListIdentitySourceSessionsCmd)
}

var (
	GetIdentitySourceSessionidentitySourceId string

	GetIdentitySourceSessionsessionId string
)

func NewGetIdentitySourceSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getSession",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.GetIdentitySourceSession(apiClient.GetConfig().Context, GetIdentitySourceSessionidentitySourceId, GetIdentitySourceSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetIdentitySourceSessionidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&GetIdentitySourceSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	return cmd
}

func init() {
	GetIdentitySourceSessionCmd := NewGetIdentitySourceSessionCmd()
	IdentitySourceCmd.AddCommand(GetIdentitySourceSessionCmd)
}

var (
	DeleteIdentitySourceSessionidentitySourceId string

	DeleteIdentitySourceSessionsessionId string
)

func NewDeleteIdentitySourceSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteSession",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.DeleteIdentitySourceSession(apiClient.GetConfig().Context, DeleteIdentitySourceSessionidentitySourceId, DeleteIdentitySourceSessionsessionId)

			resp, err := req.Execute()
			if err != nil {
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&DeleteIdentitySourceSessionidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&DeleteIdentitySourceSessionsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	return cmd
}

func init() {
	DeleteIdentitySourceSessionCmd := NewDeleteIdentitySourceSessionCmd()
	IdentitySourceCmd.AddCommand(DeleteIdentitySourceSessionCmd)
}

var (
	UploadIdentitySourceDataForDeleteidentitySourceId string

	UploadIdentitySourceDataForDeletesessionId string

	UploadIdentitySourceDataForDeletedata string
)

func NewUploadIdentitySourceDataForDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "uploadDataForDelete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.UploadIdentitySourceDataForDelete(apiClient.GetConfig().Context, UploadIdentitySourceDataForDeleteidentitySourceId, UploadIdentitySourceDataForDeletesessionId)

			if UploadIdentitySourceDataForDeletedata != "" {
				req = req.Data(UploadIdentitySourceDataForDeletedata)
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
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForDeleteidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForDeletesessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForDeletedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UploadIdentitySourceDataForDeleteCmd := NewUploadIdentitySourceDataForDeleteCmd()
	IdentitySourceCmd.AddCommand(UploadIdentitySourceDataForDeleteCmd)
}

var (
	UploadIdentitySourceDataForUpsertidentitySourceId string

	UploadIdentitySourceDataForUpsertsessionId string

	UploadIdentitySourceDataForUpsertdata string
)

func NewUploadIdentitySourceDataForUpsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "uploadDataForUpsert",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.UploadIdentitySourceDataForUpsert(apiClient.GetConfig().Context, UploadIdentitySourceDataForUpsertidentitySourceId, UploadIdentitySourceDataForUpsertsessionId)

			if UploadIdentitySourceDataForUpsertdata != "" {
				req = req.Data(UploadIdentitySourceDataForUpsertdata)
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
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForUpsertidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForUpsertsessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().StringVarP(&UploadIdentitySourceDataForUpsertdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UploadIdentitySourceDataForUpsertCmd := NewUploadIdentitySourceDataForUpsertCmd()
	IdentitySourceCmd.AddCommand(UploadIdentitySourceDataForUpsertCmd)
}

var (
	StartImportFromIdentitySourceidentitySourceId string

	StartImportFromIdentitySourcesessionId string

	StartImportFromIdentitySourcedata string
)

func NewStartImportFromIdentitySourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "startImportFrom",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.IdentitySourceAPI.StartImportFromIdentitySource(apiClient.GetConfig().Context, StartImportFromIdentitySourceidentitySourceId, StartImportFromIdentitySourcesessionId)

			if StartImportFromIdentitySourcedata != "" {
				req = req.Data(StartImportFromIdentitySourcedata)
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
			cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&StartImportFromIdentitySourceidentitySourceId, "identitySourceId", "", "", "")
	cmd.MarkFlagRequired("identitySourceId")

	cmd.Flags().StringVarP(&StartImportFromIdentitySourcesessionId, "sessionId", "", "", "")
	cmd.MarkFlagRequired("sessionId")

	cmd.Flags().StringVarP(&StartImportFromIdentitySourcedata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	StartImportFromIdentitySourceCmd := NewStartImportFromIdentitySourceCmd()
	IdentitySourceCmd.AddCommand(StartImportFromIdentitySourceCmd)
}
