package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationFeaturesCmd = &cobra.Command{
	Use:  "applicationFeatures",
	Long: "Manage ApplicationFeaturesAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationFeaturesCmd)
}

var ListFeaturesForApplicationappId string

func NewListFeaturesForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "listFeaturesForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationFeaturesAPI.ListFeaturesForApplication(apiClient.GetConfig().Context, ListFeaturesForApplicationappId)

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

	cmd.Flags().StringVarP(&ListFeaturesForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	return cmd
}

func init() {
	ListFeaturesForApplicationCmd := NewListFeaturesForApplicationCmd()
	ApplicationFeaturesCmd.AddCommand(ListFeaturesForApplicationCmd)
}

var (
	GetFeatureForApplicationappId string

	GetFeatureForApplicationfeatureName string
)

func NewGetFeatureForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "getFeatureForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationFeaturesAPI.GetFeatureForApplication(apiClient.GetConfig().Context, GetFeatureForApplicationappId, GetFeatureForApplicationfeatureName)

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

	cmd.Flags().StringVarP(&GetFeatureForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetFeatureForApplicationfeatureName, "featureName", "", "", "")
	cmd.MarkFlagRequired("featureName")

	return cmd
}

func init() {
	GetFeatureForApplicationCmd := NewGetFeatureForApplicationCmd()
	ApplicationFeaturesCmd.AddCommand(GetFeatureForApplicationCmd)
}

var (
	UpdateFeatureForApplicationappId string

	UpdateFeatureForApplicationfeatureName string

	UpdateFeatureForApplicationdata string
)

func NewUpdateFeatureForApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "updateFeatureForApplication",

		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationFeaturesAPI.UpdateFeatureForApplication(apiClient.GetConfig().Context, UpdateFeatureForApplicationappId, UpdateFeatureForApplicationfeatureName)

			if UpdateFeatureForApplicationdata != "" {
				req = req.Data(UpdateFeatureForApplicationdata)
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

	cmd.Flags().StringVarP(&UpdateFeatureForApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UpdateFeatureForApplicationfeatureName, "featureName", "", "", "")
	cmd.MarkFlagRequired("featureName")

	cmd.Flags().StringVarP(&UpdateFeatureForApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateFeatureForApplicationCmd := NewUpdateFeatureForApplicationCmd()
	ApplicationFeaturesCmd.AddCommand(UpdateFeatureForApplicationCmd)
}
