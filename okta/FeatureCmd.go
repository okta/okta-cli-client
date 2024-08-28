package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var FeatureCmd = &cobra.Command{
	Use:  "feature",
	Long: "Manage FeatureAPI",
}

func init() {
	rootCmd.AddCommand(FeatureCmd)
}

func NewListFeaturesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lists",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.FeatureAPI.ListFeatures(apiClient.GetConfig().Context)

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
	ListFeaturesCmd := NewListFeaturesCmd()
	FeatureCmd.AddCommand(ListFeaturesCmd)
}

var GetFeaturefeatureId string

func NewGetFeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.FeatureAPI.GetFeature(apiClient.GetConfig().Context, GetFeaturefeatureId)

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

	cmd.Flags().StringVarP(&GetFeaturefeatureId, "featureId", "", "", "")
	cmd.MarkFlagRequired("featureId")

	return cmd
}

func init() {
	GetFeatureCmd := NewGetFeatureCmd()
	FeatureCmd.AddCommand(GetFeatureCmd)
}

var ListFeatureDependenciesfeatureId string

func NewListFeatureDependenciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listDependencies",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.FeatureAPI.ListFeatureDependencies(apiClient.GetConfig().Context, ListFeatureDependenciesfeatureId)

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

	cmd.Flags().StringVarP(&ListFeatureDependenciesfeatureId, "featureId", "", "", "")
	cmd.MarkFlagRequired("featureId")

	return cmd
}

func init() {
	ListFeatureDependenciesCmd := NewListFeatureDependenciesCmd()
	FeatureCmd.AddCommand(ListFeatureDependenciesCmd)
}

var ListFeatureDependentsfeatureId string

func NewListFeatureDependentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listDependents",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.FeatureAPI.ListFeatureDependents(apiClient.GetConfig().Context, ListFeatureDependentsfeatureId)

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

	cmd.Flags().StringVarP(&ListFeatureDependentsfeatureId, "featureId", "", "", "")
	cmd.MarkFlagRequired("featureId")

	return cmd
}

func init() {
	ListFeatureDependentsCmd := NewListFeatureDependentsCmd()
	FeatureCmd.AddCommand(ListFeatureDependentsCmd)
}

var (
	UpdateFeatureLifecyclefeatureId string

	UpdateFeatureLifecyclelifecycle string

	UpdateFeatureLifecycledata string
)

func NewUpdateFeatureLifecycleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateLifecycle",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.FeatureAPI.UpdateFeatureLifecycle(apiClient.GetConfig().Context, UpdateFeatureLifecyclefeatureId, UpdateFeatureLifecyclelifecycle)

			if UpdateFeatureLifecycledata != "" {
				req = req.Data(UpdateFeatureLifecycledata)
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

	cmd.Flags().StringVarP(&UpdateFeatureLifecyclefeatureId, "featureId", "", "", "")
	cmd.MarkFlagRequired("featureId")

	cmd.Flags().StringVarP(&UpdateFeatureLifecyclelifecycle, "lifecycle", "", "", "")
	cmd.MarkFlagRequired("lifecycle")

	cmd.Flags().StringVarP(&UpdateFeatureLifecycledata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateFeatureLifecycleCmd := NewUpdateFeatureLifecycleCmd()
	FeatureCmd.AddCommand(UpdateFeatureLifecycleCmd)
}
