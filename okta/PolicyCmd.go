package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var PolicyCmd = &cobra.Command{
	Use:  "policy",
	Long: "Manage PolicyAPI",
}

func init() {
	rootCmd.AddCommand(PolicyCmd)
}

var CreatePolicydata string

func NewCreatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.CreatePolicy(apiClient.GetConfig().Context)

			if CreatePolicydata != "" {
				req = req.Data(CreatePolicydata)
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

	cmd.Flags().StringVarP(&CreatePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreatePolicyCmd := NewCreatePolicyCmd()
	PolicyCmd.AddCommand(CreatePolicyCmd)
}

func NewListPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listPolicies",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicies(apiClient.GetConfig().Context)

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

	return cmd
}

func init() {
	ListPoliciesCmd := NewListPoliciesCmd()
	PolicyCmd.AddCommand(ListPoliciesCmd)
}

var CreatePolicySimulationdata string

func NewCreatePolicySimulationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createSimulation",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.CreatePolicySimulation(apiClient.GetConfig().Context)

			if CreatePolicySimulationdata != "" {
				req = req.Data(CreatePolicySimulationdata)
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

	cmd.Flags().StringVarP(&CreatePolicySimulationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreatePolicySimulationCmd := NewCreatePolicySimulationCmd()
	PolicyCmd.AddCommand(CreatePolicySimulationCmd)
}

var GetPolicypolicyId string

func NewGetPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.GetPolicy(apiClient.GetConfig().Context, GetPolicypolicyId)

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

	cmd.Flags().StringVarP(&GetPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	GetPolicyCmd := NewGetPolicyCmd()
	PolicyCmd.AddCommand(GetPolicyCmd)
}

var (
	ReplacePolicypolicyId string

	ReplacePolicydata string
)

func NewReplacePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replace",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ReplacePolicy(apiClient.GetConfig().Context, ReplacePolicypolicyId)

			if ReplacePolicydata != "" {
				req = req.Data(ReplacePolicydata)
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

	cmd.Flags().StringVarP(&ReplacePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplacePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacePolicyCmd := NewReplacePolicyCmd()
	PolicyCmd.AddCommand(ReplacePolicyCmd)
}

var DeletePolicypolicyId string

func NewDeletePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeletePolicy(apiClient.GetConfig().Context, DeletePolicypolicyId)

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

	cmd.Flags().StringVarP(&DeletePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	DeletePolicyCmd := NewDeletePolicyCmd()
	PolicyCmd.AddCommand(DeletePolicyCmd)
}

var ListPolicyAppspolicyId string

func NewListPolicyAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listApps",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicyApps(apiClient.GetConfig().Context, ListPolicyAppspolicyId)

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

	cmd.Flags().StringVarP(&ListPolicyAppspolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	ListPolicyAppsCmd := NewListPolicyAppsCmd()
	PolicyCmd.AddCommand(ListPolicyAppsCmd)
}

var (
	ClonePolicypolicyId string

	ClonePolicydata string
)

func NewClonePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "clone",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ClonePolicy(apiClient.GetConfig().Context, ClonePolicypolicyId)

			if ClonePolicydata != "" {
				req = req.Data(ClonePolicydata)
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

	cmd.Flags().StringVarP(&ClonePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ClonePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ClonePolicyCmd := NewClonePolicyCmd()
	PolicyCmd.AddCommand(ClonePolicyCmd)
}

var (
	ActivatePolicypolicyId string

	ActivatePolicydata string
)

func NewActivatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ActivatePolicy(apiClient.GetConfig().Context, ActivatePolicypolicyId)

			if ActivatePolicydata != "" {
				req = req.Data(ActivatePolicydata)
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

	cmd.Flags().StringVarP(&ActivatePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ActivatePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivatePolicyCmd := NewActivatePolicyCmd()
	PolicyCmd.AddCommand(ActivatePolicyCmd)
}

var (
	DeactivatePolicypolicyId string

	DeactivatePolicydata string
)

func NewDeactivatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivate",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeactivatePolicy(apiClient.GetConfig().Context, DeactivatePolicypolicyId)

			if DeactivatePolicydata != "" {
				req = req.Data(DeactivatePolicydata)
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

	cmd.Flags().StringVarP(&DeactivatePolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeactivatePolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	DeactivatePolicyCmd := NewDeactivatePolicyCmd()
	PolicyCmd.AddCommand(DeactivatePolicyCmd)
}

var (
	MapResourceToPolicypolicyId string

	MapResourceToPolicydata string
)

func NewMapResourceToPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "mapResourceTo",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.MapResourceToPolicy(apiClient.GetConfig().Context, MapResourceToPolicypolicyId)

			if MapResourceToPolicydata != "" {
				req = req.Data(MapResourceToPolicydata)
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

	cmd.Flags().StringVarP(&MapResourceToPolicypolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&MapResourceToPolicydata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	MapResourceToPolicyCmd := NewMapResourceToPolicyCmd()
	PolicyCmd.AddCommand(MapResourceToPolicyCmd)
}

var ListPolicyMappingspolicyId string

func NewListPolicyMappingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listMappings",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicyMappings(apiClient.GetConfig().Context, ListPolicyMappingspolicyId)

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

	cmd.Flags().StringVarP(&ListPolicyMappingspolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	ListPolicyMappingsCmd := NewListPolicyMappingsCmd()
	PolicyCmd.AddCommand(ListPolicyMappingsCmd)
}

var (
	GetPolicyMappingpolicyId string

	GetPolicyMappingmappingId string
)

func NewGetPolicyMappingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getMapping",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.GetPolicyMapping(apiClient.GetConfig().Context, GetPolicyMappingpolicyId, GetPolicyMappingmappingId)

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

	cmd.Flags().StringVarP(&GetPolicyMappingpolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&GetPolicyMappingmappingId, "mappingId", "", "", "")
	cmd.MarkFlagRequired("mappingId")

	return cmd
}

func init() {
	GetPolicyMappingCmd := NewGetPolicyMappingCmd()
	PolicyCmd.AddCommand(GetPolicyMappingCmd)
}

var (
	DeletePolicyResourceMappingpolicyId string

	DeletePolicyResourceMappingmappingId string
)

func NewDeletePolicyResourceMappingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteResourceMapping",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeletePolicyResourceMapping(apiClient.GetConfig().Context, DeletePolicyResourceMappingpolicyId, DeletePolicyResourceMappingmappingId)

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

	cmd.Flags().StringVarP(&DeletePolicyResourceMappingpolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeletePolicyResourceMappingmappingId, "mappingId", "", "", "")
	cmd.MarkFlagRequired("mappingId")

	return cmd
}

func init() {
	DeletePolicyResourceMappingCmd := NewDeletePolicyResourceMappingCmd()
	PolicyCmd.AddCommand(DeletePolicyResourceMappingCmd)
}

var (
	CreatePolicyRulepolicyId string

	CreatePolicyRuledata string
)

func NewCreatePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "createRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.CreatePolicyRule(apiClient.GetConfig().Context, CreatePolicyRulepolicyId)

			if CreatePolicyRuledata != "" {
				req = req.Data(CreatePolicyRuledata)
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

	cmd.Flags().StringVarP(&CreatePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&CreatePolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	CreatePolicyRuleCmd := NewCreatePolicyRuleCmd()
	PolicyCmd.AddCommand(CreatePolicyRuleCmd)
}

var ListPolicyRulespolicyId string

func NewListPolicyRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listRules",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ListPolicyRules(apiClient.GetConfig().Context, ListPolicyRulespolicyId)

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

	cmd.Flags().StringVarP(&ListPolicyRulespolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	return cmd
}

func init() {
	ListPolicyRulesCmd := NewListPolicyRulesCmd()
	PolicyCmd.AddCommand(ListPolicyRulesCmd)
}

var (
	GetPolicyRulepolicyId string

	GetPolicyRuleruleId string
)

func NewGetPolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.GetPolicyRule(apiClient.GetConfig().Context, GetPolicyRulepolicyId, GetPolicyRuleruleId)

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

	cmd.Flags().StringVarP(&GetPolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&GetPolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	return cmd
}

func init() {
	GetPolicyRuleCmd := NewGetPolicyRuleCmd()
	PolicyCmd.AddCommand(GetPolicyRuleCmd)
}

var (
	ReplacePolicyRulepolicyId string

	ReplacePolicyRuleruleId string

	ReplacePolicyRuledata string
)

func NewReplacePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "replaceRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ReplacePolicyRule(apiClient.GetConfig().Context, ReplacePolicyRulepolicyId, ReplacePolicyRuleruleId)

			if ReplacePolicyRuledata != "" {
				req = req.Data(ReplacePolicyRuledata)
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

	cmd.Flags().StringVarP(&ReplacePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ReplacePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().StringVarP(&ReplacePolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ReplacePolicyRuleCmd := NewReplacePolicyRuleCmd()
	PolicyCmd.AddCommand(ReplacePolicyRuleCmd)
}

var (
	DeletePolicyRulepolicyId string

	DeletePolicyRuleruleId string
)

func NewDeletePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deleteRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeletePolicyRule(apiClient.GetConfig().Context, DeletePolicyRulepolicyId, DeletePolicyRuleruleId)

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

	cmd.Flags().StringVarP(&DeletePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeletePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	return cmd
}

func init() {
	DeletePolicyRuleCmd := NewDeletePolicyRuleCmd()
	PolicyCmd.AddCommand(DeletePolicyRuleCmd)
}

var (
	ActivatePolicyRulepolicyId string

	ActivatePolicyRuleruleId string

	ActivatePolicyRuledata string
)

func NewActivatePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "activateRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.ActivatePolicyRule(apiClient.GetConfig().Context, ActivatePolicyRulepolicyId, ActivatePolicyRuleruleId)

			if ActivatePolicyRuledata != "" {
				req = req.Data(ActivatePolicyRuledata)
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

	cmd.Flags().StringVarP(&ActivatePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&ActivatePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().StringVarP(&ActivatePolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	ActivatePolicyRuleCmd := NewActivatePolicyRuleCmd()
	PolicyCmd.AddCommand(ActivatePolicyRuleCmd)
}

var (
	DeactivatePolicyRulepolicyId string

	DeactivatePolicyRuleruleId string

	DeactivatePolicyRuledata string
)

func NewDeactivatePolicyRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deactivateRule",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.PolicyAPI.DeactivatePolicyRule(apiClient.GetConfig().Context, DeactivatePolicyRulepolicyId, DeactivatePolicyRuleruleId)

			if DeactivatePolicyRuledata != "" {
				req = req.Data(DeactivatePolicyRuledata)
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

	cmd.Flags().StringVarP(&DeactivatePolicyRulepolicyId, "policyId", "", "", "")
	cmd.MarkFlagRequired("policyId")

	cmd.Flags().StringVarP(&DeactivatePolicyRuleruleId, "ruleId", "", "", "")
	cmd.MarkFlagRequired("ruleId")

	cmd.Flags().StringVarP(&DeactivatePolicyRuledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	DeactivatePolicyRuleCmd := NewDeactivatePolicyRuleCmd()
	PolicyCmd.AddCommand(DeactivatePolicyRuleCmd)
}
