package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	v5backOff "github.com/cenkalti/backoff/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	DefaultCliAPIURL           = "https://start.okta.dev"
	DefaultRegistrationBaseURL = "https://okta-devok12.okta.com"
	DefaultRegistrationID      = "reg405abrRAkn0TRf5d6"
	DefaultConfigPath          = ".okta/okta.yaml"
)

type OrgRequest struct {
	UserProfile OrgUserProfile `json:"userProfile,omitempty"`
}

type OrgUserProfile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Country   string `json:"country"`
	OktaOie   bool   `json:"okta_oie"`
}

type OrgResponse struct {
	ActivationToken      string `json:"activationToken"`
	DeveloperOrgCliToken string `json:"developerOrgCliToken"`
}

type VerifyResponse struct {
	Status   string `json:"status"`
	APIToken string `json:"apiToken"`
	OrgURL   string `json:"orgUrl"`
}

type Registration struct {
	FirstName string
	LastName  string
	Email     string
	Company   string
	Country   string
	Oie       bool
}

type Config struct {
	Okta struct {
		Client struct {
			OrgURL string `yaml:"orgUrl"`
			Token  string `yaml:"token"`
		} `yaml:"client"`
	} `yaml:"okta"`
}

var (
	firstName = Flag{
		Name:       "First Name",
		LongForm:   "first-name",
		Help:       "First name used when registering a new Okta account.",
		IsRequired: true,
	}
	lastName = Flag{
		Name:       "Last Name",
		LongForm:   "last-name",
		Help:       "Last name used when registering a new Okta account.",
		IsRequired: true,
	}
	emailAddress = Flag{
		Name:       "Email",
		LongForm:   "email",
		Help:       "Email used when registering a new Okta account.",
		IsRequired: true,
	}
	country = Flag{
		Name:       "Country",
		LongForm:   "country",
		Help:       "Country of residence",
		IsRequired: true,
	}
	company = Flag{
		Name:       "Company",
		LongForm:   "company",
		Help:       "Company/organization used when registering a new Okta account.",
		IsRequired: true,
	}
	oie = Flag{
		Name:       "OIE",
		LongForm:   "oie",
		Help:       "Create Okta account with OIE enabled.",
		IsRequired: false,
	}
)

func NewRegisterCmd() *cobra.Command {
	inputs := Registration{}
	cmd := &cobra.Command{
		Use:   "register",
		Args:  cobra.NoArgs,
		Short: "Sign up for a new Okta account",
		Long: "Sign up for a new Okta account.\n\n" +
			"To register interactively, use `okta-cli-client register` with no arguments.\n\n" +
			"To register non-interactively, supply at least the first name, and type through the flags.",
		Example: `okta-cli-client register
okta-cli-client register --first-name firstName
okta-cli-client register --first-name firstName --last-name lastName
okta-cli-client register --first-name firstName --last-name lastName --email email 
okta-cli-client register --first-name firstName --last-name lastName --email email --country country
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := firstName.Ask(cmd, &inputs.FirstName, nil); err != nil {
				return err
			}
			if err := lastName.Ask(cmd, &inputs.LastName, nil); err != nil {
				return err
			}
			if err := emailAddress.Ask(cmd, &inputs.Email, nil); err != nil {
				return err
			}
			if err := country.Ask(cmd, &inputs.Country, nil); err != nil {
				return err
			}
			if err := backedConfigFile(); err != nil {
				return err
			}
			orgResp, err := createNewOrg(inputs)
			if err != nil {
				return err
			}
			fmt.Println("Check your email")
			var status *string
			var resp *VerifyResponse
			operation := func() (*string, error) {
				resp, status, err = verifyOrg(orgResp.DeveloperOrgCliToken)
				if status == nil || *status != "ACTIVE" {
					return nil, v5backOff.RetryAfter(20)
				}
				return status, err
			}
			result, err := v5backOff.Retry(context.TODO(), operation, v5backOff.WithMaxTries(12))
			if err != nil {
				return err
			}
			if result != nil && *result == "ACTIVE" {
				err = writeConfigFile(resp)
				if err != nil {
					return err
				}
			} else {
				return errors.New("Failed to verify new Okta Organization")
			}
			fmt.Println("New Okta Account created!")
			fmt.Printf("Your Okta Domain: %v", resp.OrgURL)
			return nil
		},
	}
	firstName.RegisterString(cmd, &inputs.FirstName, "")
	lastName.RegisterString(cmd, &inputs.LastName, "")
	emailAddress.RegisterString(cmd, &inputs.Email, "")
	company.RegisterString(cmd, &inputs.Company, "")
	country.RegisterString(cmd, &inputs.Country, "")
	oie.RegisterBool(cmd, &inputs.Oie, true)
	return cmd
}

func init() {
	rootCmd.AddCommand(NewRegisterCmd())
}

func createNewOrg(inputs Registration) (*OrgResponse, error) {
	fmt.Println("Creating new Okta Organization, this may take a minute:")
	reqBody := OrgRequest{
		UserProfile: OrgUserProfile{
			FirstName: inputs.FirstName,
			LastName:  inputs.LastName,
			Email:     inputs.Email,
			Company:   inputs.Company,
			Country:   inputs.Country,
			OktaOie:   inputs.Oie,
		},
	}
	reqHeader := map[string]string{}
	req, err := apiClient.PrepareRequest(apiClient.GetConfig().Context, fmt.Sprintf("%v/api/v1/registration/%v/register", DefaultRegistrationBaseURL, DefaultRegistrationID), http.MethodPost, reqBody, reqHeader, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Okta Organization. You can register manually by going to https://developer.okta.com/signup. Error: %v", err)
	}
	resp, err := apiClient.Do(apiClient.GetConfig().Context, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Okta Organization. You can register manually by going to https://developer.okta.com/signup. Error: %v", err)
	}
	respBody := OrgResponse{}
	if resp.StatusCode != http.StatusOK || json.NewDecoder(resp.Body).Decode(&respBody) != nil {
		return nil, fmt.Errorf("Failed to create Okta Organization. You can register manually by going to https://developer.okta.com/signup. Error: %v", err)
	}
	fmt.Println("An account activation email has been sent to you.")
	return &respBody, nil
}

func verifyOrg(identifier string) (*VerifyResponse, *string, error) {
	var status string
	respBody := VerifyResponse{}
	reqHeader := map[string]string{}
	req, err := apiClient.PrepareRequest(apiClient.GetConfig().Context, fmt.Sprintf("%v/api/internal/v1/developer/redeem/%v", DefaultRegistrationBaseURL, identifier), http.MethodGet, nil, reqHeader, nil, nil, nil)
	if err != nil {
		return &respBody, &status, err
	}
	resp, err := apiClient.Do(apiClient.GetConfig().Context, req)
	if err != nil {
		return &respBody, &status, err
	}
	if resp.StatusCode != http.StatusOK || json.NewDecoder(resp.Body).Decode(&respBody) != nil {
		return &respBody, &status, errors.New("error when reading body")
	}
	status = respBody.Status
	return &respBody, &status, nil
}

func backedConfigFile() error {
	configPath, err := getOktaConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	newFilePath := fmt.Sprintf("%v.%v", configPath, time.Now().Unix())
	srcFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	if err != nil {
		return err
	}
	fmt.Printf("Configuration file backed: %v", newFilePath)
	return nil
}

func writeConfigFile(resp *VerifyResponse) error {
	fmt.Println("Check your email")
	var config Config
	config.Okta.Client.OrgURL = resp.OrgURL
	config.Okta.Client.Token = resp.APIToken
	configPath, err := getOktaConfigPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(config)
}

func getOktaConfigPath() (string, error) {
	usr, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v/%v", usr, DefaultConfigPath), nil
}
