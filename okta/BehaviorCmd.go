package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var BehaviorCmd = &cobra.Command{
	Use:  "behavior",
	Long: "Manage BehaviorAPI",
}

func init() {
	rootCmd.AddCommand(BehaviorCmd)
}

var CreateBehaviorDetectionRuledata string

func NewCreateBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "createDetectionRule",
		Long: "Create a Behavior Detection Rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.CreateBehaviorDetectionRule(apiClient.GetConfig().Context)

			if CreateBehaviorDetectionRuledata != "" {
				req = req.Data(CreateBehaviorDetectionRuledata)
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

	cmd.Flags().StringVarP(&CreateBehaviorDetectionRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreateBehaviorDetectionRuleCmd := NewCreateBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(CreateBehaviorDetectionRuleCmd)
}

func NewListBehaviorDetectionRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "listDetectionRules",
		Long: "List all Behavior Detection Rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.ListBehaviorDetectionRules(apiClient.GetConfig().Context)

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
	ListBehaviorDetectionRulesCmd := NewListBehaviorDetectionRulesCmd()
	BehaviorCmd.AddCommand(ListBehaviorDetectionRulesCmd)
}

var GetBehaviorDetectionRulebehaviorId string

func NewGetBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getDetectionRule",
		Long: "Retrieve a Behavior Detection Rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.GetBehaviorDetectionRule(apiClient.GetConfig().Context, GetBehaviorDetectionRulebehaviorId)

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

	cmd.Flags().StringVarP(&GetBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	return cmd
}

func init() {
	GetBehaviorDetectionRuleCmd := NewGetBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(GetBehaviorDetectionRuleCmd)
}

var (
	ReplaceBehaviorDetectionRulebehaviorId string

	ReplaceBehaviorDetectionRuledata string
)

func NewReplaceBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "replaceDetectionRule",
		Long: "Replace a Behavior Detection Rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.ReplaceBehaviorDetectionRule(apiClient.GetConfig().Context, ReplaceBehaviorDetectionRulebehaviorId)

			if ReplaceBehaviorDetectionRuledata != "" {
				req = req.Data(ReplaceBehaviorDetectionRuledata)
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

	cmd.Flags().StringVarP(&ReplaceBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	cmd.Flags().StringVarP(&ReplaceBehaviorDetectionRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplaceBehaviorDetectionRuleCmd := NewReplaceBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(ReplaceBehaviorDetectionRuleCmd)
}

var DeleteBehaviorDetectionRulebehaviorId string

func NewDeleteBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deleteDetectionRule",
		Long: "Delete a Behavior Detection Rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.DeleteBehaviorDetectionRule(apiClient.GetConfig().Context, DeleteBehaviorDetectionRulebehaviorId)

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

	cmd.Flags().StringVarP(&DeleteBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	return cmd
}

func init() {
	DeleteBehaviorDetectionRuleCmd := NewDeleteBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(DeleteBehaviorDetectionRuleCmd)
}

var ActivateBehaviorDetectionRulebehaviorId string

func NewActivateBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "activateDetectionRule",
		Long: "Activate a Behavior Detection Rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.ActivateBehaviorDetectionRule(apiClient.GetConfig().Context, ActivateBehaviorDetectionRulebehaviorId)

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

	cmd.Flags().StringVarP(&ActivateBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	return cmd
}

func init() {
	ActivateBehaviorDetectionRuleCmd := NewActivateBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(ActivateBehaviorDetectionRuleCmd)
}

var DeactivateBehaviorDetectionRulebehaviorId string

func NewDeactivateBehaviorDetectionRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deactivateDetectionRule",
		Long: "Deactivate a Behavior Detection Rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.BehaviorAPI.DeactivateBehaviorDetectionRule(apiClient.GetConfig().Context, DeactivateBehaviorDetectionRulebehaviorId)

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

	cmd.Flags().StringVarP(&DeactivateBehaviorDetectionRulebehaviorId, "behaviorId", "", "", "")
	cmd.MarkFlagRequired("behaviorId")

	return cmd
}

func init() {
	DeactivateBehaviorDetectionRuleCmd := NewDeactivateBehaviorDetectionRuleCmd()
	BehaviorCmd.AddCommand(DeactivateBehaviorDetectionRuleCmd)
}
