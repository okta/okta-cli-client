package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var LogStreamCmd = &cobra.Command{
	Use:  "logStream",
	Long: "Manage LogStreamAPI",
}

func init() {
	rootCmd.AddCommand(LogStreamCmd)
}

var CreateLogStreamdata string

func NewCreateLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a Log Stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.CreateLogStream(apiClient.GetConfig().Context)

			if CreateLogStreamdata != "" {
				req = req.Data(CreateLogStreamdata)
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

	cmd.Flags().StringVarP(&CreateLogStreamdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateLogStreamCmd := NewCreateLogStreamCmd()
	LogStreamCmd.AddCommand(CreateLogStreamCmd)
}

func NewListLogStreamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Log Streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.ListLogStreams(apiClient.GetConfig().Context)

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
	ListLogStreamsCmd := NewListLogStreamsCmd()
	LogStreamCmd.AddCommand(ListLogStreamsCmd)
}

var GetLogStreamlogStreamId string

func NewGetLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Log Stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.GetLogStream(apiClient.GetConfig().Context, GetLogStreamlogStreamId)

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

	cmd.Flags().StringVarP(&GetLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	return cmd
}

func init() {
	GetLogStreamCmd := NewGetLogStreamCmd()
	LogStreamCmd.AddCommand(GetLogStreamCmd)
}

var (
	ReplaceLogStreamlogStreamId string

	ReplaceLogStreamdata string
)

func NewReplaceLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replace",
		Long: "Replace a Log Stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.ReplaceLogStream(apiClient.GetConfig().Context, ReplaceLogStreamlogStreamId)

			if ReplaceLogStreamdata != "" {
				req = req.Data(ReplaceLogStreamdata)
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

	cmd.Flags().StringVarP(&ReplaceLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	cmd.Flags().StringVarP(&ReplaceLogStreamdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceLogStreamCmd := NewReplaceLogStreamCmd()
	LogStreamCmd.AddCommand(ReplaceLogStreamCmd)
}

var DeleteLogStreamlogStreamId string

func NewDeleteLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete a Log Stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.DeleteLogStream(apiClient.GetConfig().Context, DeleteLogStreamlogStreamId)

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

	cmd.Flags().StringVarP(&DeleteLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	return cmd
}

func init() {
	DeleteLogStreamCmd := NewDeleteLogStreamCmd()
	LogStreamCmd.AddCommand(DeleteLogStreamCmd)
}

var ActivateLogStreamlogStreamId string

func NewActivateLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activate",
		Long: "Activate a Log Stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.ActivateLogStream(apiClient.GetConfig().Context, ActivateLogStreamlogStreamId)

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

	cmd.Flags().StringVarP(&ActivateLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	return cmd
}

func init() {
	ActivateLogStreamCmd := NewActivateLogStreamCmd()
	LogStreamCmd.AddCommand(ActivateLogStreamCmd)
}

var DeactivateLogStreamlogStreamId string

func NewDeactivateLogStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivate",
		Long: "Deactivate a Log Stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.LogStreamAPI.DeactivateLogStream(apiClient.GetConfig().Context, DeactivateLogStreamlogStreamId)

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

	cmd.Flags().StringVarP(&DeactivateLogStreamlogStreamId, "logStreamId", "", "", "")
	cmd.MarkFlagRequired("logStreamId")

	return cmd
}

func init() {
	DeactivateLogStreamCmd := NewDeactivateLogStreamCmd()
	LogStreamCmd.AddCommand(DeactivateLogStreamCmd)
}
