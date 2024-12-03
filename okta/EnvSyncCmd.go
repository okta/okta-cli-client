package okta

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
	"github.com/okta/okta-cli-client/utils"
    "github.com/spf13/cobra"
)


var EnvSyncCmd = &cobra.Command{
	Use:  "envsync",
	Long: "backup and restore okta environments",
}

var UserdataPath string

type UserProfile struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Login       string `json:"login"`
	MobilePhone string `json:"mobilePhone"`
	SecondEmail string `json:"secondEmail"`
}

type UserData struct {
	Profile UserProfile `json:"profile"`
}

func init() {
	rootCmd.AddCommand(EnvSyncCmd)
	// Add the create command to EnvSyncPushCmd
	EnvSyncPushCmd.AddCommand(NewEnvSyncPushUserCmd())
}

func NewEnvSyncPushUserCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use: "pushuser",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read and parse the userdata file
			userData, err := os.ReadFile(UserdataPath)
			if err != nil {
				return fmt.Errorf("error reading user data file: %w", err)
			}

			var data UserData
			if err = json.Unmarshal(userData, &data); err != nil {
				return fmt.Errorf("error parsing user data file: %w", err)
			}

			// Create the API request data
			createData := fmt.Sprintf(
				`{"credentials":{"password":{"value":"Hell4W0rld"}},"profile":{"email":"%s","firstName":"%s","lastName":"%s","login":"%s"}}`,
				data.Profile.Email, data.Profile.FirstName, data.Profile.LastName, data.Profile.Login,
			)

			// Make the API request
			req := apiClient.UserAPI.CreateUser(apiClient.GetConfig().Context)
			req = req.Data(createData)

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
			return nil
		},
	}

	createCmd.Flags().StringVarP(&UserdataPath, "userdata", "u", "", "Path to the userdata file")
	createCmd.MarkFlagRequired("userdata")

	return createCmd
}
