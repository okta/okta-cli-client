package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ProfileMappingCmd = &cobra.Command{
	Use:  "profileMapping",
	Long: "Manage ProfileMappingAPI",
}

func init() {
	rootCmd.AddCommand(ProfileMappingCmd)
}

func NewListProfileMappingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all Profile Mappings",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ProfileMappingAPI.ListProfileMappings(apiClient.GetConfig().Context)

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
	ListProfileMappingsCmd := NewListProfileMappingsCmd()
	ProfileMappingCmd.AddCommand(ListProfileMappingsCmd)
}

var (
	UpdateProfileMappingmappingId string

	UpdateProfileMappingdata string
)

func NewUpdateProfileMappingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update",
		Long: "Update a Profile Mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ProfileMappingAPI.UpdateProfileMapping(apiClient.GetConfig().Context, UpdateProfileMappingmappingId)

			if UpdateProfileMappingdata != "" {
				req = req.Data(UpdateProfileMappingdata)
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

	cmd.Flags().StringVarP(&UpdateProfileMappingmappingId, "mappingId", "", "", "")
	cmd.MarkFlagRequired("mappingId")

	cmd.Flags().StringVarP(&UpdateProfileMappingdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	UpdateProfileMappingCmd := NewUpdateProfileMappingCmd()
	ProfileMappingCmd.AddCommand(UpdateProfileMappingCmd)
}

var GetProfileMappingmappingId string

func NewGetProfileMappingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a Profile Mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ProfileMappingAPI.GetProfileMapping(apiClient.GetConfig().Context, GetProfileMappingmappingId)

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

	cmd.Flags().StringVarP(&GetProfileMappingmappingId, "mappingId", "", "", "")
	cmd.MarkFlagRequired("mappingId")

	return cmd
}

func init() {
	GetProfileMappingCmd := NewGetProfileMappingCmd()
	ProfileMappingCmd.AddCommand(GetProfileMappingCmd)
}
