/*
Okta Admin Management

Allows customers to easily access the Okta Management APIs

API version: 5.1.0
Contact: devex-public@okta.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sdk

import (
	"encoding/json"
)

// AuthenticatorProviderConfiguration struct for AuthenticatorProviderConfiguration
type AuthenticatorProviderConfiguration struct {
	AuthPort             *int32                                              `json:"authPort,omitempty"`
	HostName             *string                                             `json:"hostName,omitempty"`
	InstanceId           *string                                             `json:"instanceId,omitempty"`
	SharedSecret         *string                                             `json:"sharedSecret,omitempty"`
	UserNameTemplate     *AuthenticatorProviderConfigurationUserNameTemplate `json:"userNameTemplate,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _AuthenticatorProviderConfiguration AuthenticatorProviderConfiguration

// NewAuthenticatorProviderConfiguration instantiates a new AuthenticatorProviderConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthenticatorProviderConfiguration() *AuthenticatorProviderConfiguration {
	this := AuthenticatorProviderConfiguration{}
	return &this
}

// NewAuthenticatorProviderConfigurationWithDefaults instantiates a new AuthenticatorProviderConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthenticatorProviderConfigurationWithDefaults() *AuthenticatorProviderConfiguration {
	this := AuthenticatorProviderConfiguration{}
	return &this
}

// GetAuthPort returns the AuthPort field value if set, zero value otherwise.
func (o *AuthenticatorProviderConfiguration) GetAuthPort() int32 {
	if o == nil || o.AuthPort == nil {
		var ret int32
		return ret
	}
	return *o.AuthPort
}

// GetAuthPortOk returns a tuple with the AuthPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProviderConfiguration) GetAuthPortOk() (*int32, bool) {
	if o == nil || o.AuthPort == nil {
		return nil, false
	}
	return o.AuthPort, true
}

// HasAuthPort returns a boolean if a field has been set.
func (o *AuthenticatorProviderConfiguration) HasAuthPort() bool {
	if o != nil && o.AuthPort != nil {
		return true
	}

	return false
}

// SetAuthPort gets a reference to the given int32 and assigns it to the AuthPort field.
func (o *AuthenticatorProviderConfiguration) SetAuthPort(v int32) {
	o.AuthPort = &v
}

// GetHostName returns the HostName field value if set, zero value otherwise.
func (o *AuthenticatorProviderConfiguration) GetHostName() string {
	if o == nil || o.HostName == nil {
		var ret string
		return ret
	}
	return *o.HostName
}

// GetHostNameOk returns a tuple with the HostName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProviderConfiguration) GetHostNameOk() (*string, bool) {
	if o == nil || o.HostName == nil {
		return nil, false
	}
	return o.HostName, true
}

// HasHostName returns a boolean if a field has been set.
func (o *AuthenticatorProviderConfiguration) HasHostName() bool {
	if o != nil && o.HostName != nil {
		return true
	}

	return false
}

// SetHostName gets a reference to the given string and assigns it to the HostName field.
func (o *AuthenticatorProviderConfiguration) SetHostName(v string) {
	o.HostName = &v
}

// GetInstanceId returns the InstanceId field value if set, zero value otherwise.
func (o *AuthenticatorProviderConfiguration) GetInstanceId() string {
	if o == nil || o.InstanceId == nil {
		var ret string
		return ret
	}
	return *o.InstanceId
}

// GetInstanceIdOk returns a tuple with the InstanceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProviderConfiguration) GetInstanceIdOk() (*string, bool) {
	if o == nil || o.InstanceId == nil {
		return nil, false
	}
	return o.InstanceId, true
}

// HasInstanceId returns a boolean if a field has been set.
func (o *AuthenticatorProviderConfiguration) HasInstanceId() bool {
	if o != nil && o.InstanceId != nil {
		return true
	}

	return false
}

// SetInstanceId gets a reference to the given string and assigns it to the InstanceId field.
func (o *AuthenticatorProviderConfiguration) SetInstanceId(v string) {
	o.InstanceId = &v
}

// GetSharedSecret returns the SharedSecret field value if set, zero value otherwise.
func (o *AuthenticatorProviderConfiguration) GetSharedSecret() string {
	if o == nil || o.SharedSecret == nil {
		var ret string
		return ret
	}
	return *o.SharedSecret
}

// GetSharedSecretOk returns a tuple with the SharedSecret field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProviderConfiguration) GetSharedSecretOk() (*string, bool) {
	if o == nil || o.SharedSecret == nil {
		return nil, false
	}
	return o.SharedSecret, true
}

// HasSharedSecret returns a boolean if a field has been set.
func (o *AuthenticatorProviderConfiguration) HasSharedSecret() bool {
	if o != nil && o.SharedSecret != nil {
		return true
	}

	return false
}

// SetSharedSecret gets a reference to the given string and assigns it to the SharedSecret field.
func (o *AuthenticatorProviderConfiguration) SetSharedSecret(v string) {
	o.SharedSecret = &v
}

// GetUserNameTemplate returns the UserNameTemplate field value if set, zero value otherwise.
func (o *AuthenticatorProviderConfiguration) GetUserNameTemplate() AuthenticatorProviderConfigurationUserNameTemplate {
	if o == nil || o.UserNameTemplate == nil {
		var ret AuthenticatorProviderConfigurationUserNameTemplate
		return ret
	}
	return *o.UserNameTemplate
}

// GetUserNameTemplateOk returns a tuple with the UserNameTemplate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProviderConfiguration) GetUserNameTemplateOk() (*AuthenticatorProviderConfigurationUserNameTemplate, bool) {
	if o == nil || o.UserNameTemplate == nil {
		return nil, false
	}
	return o.UserNameTemplate, true
}

// HasUserNameTemplate returns a boolean if a field has been set.
func (o *AuthenticatorProviderConfiguration) HasUserNameTemplate() bool {
	if o != nil && o.UserNameTemplate != nil {
		return true
	}

	return false
}

// SetUserNameTemplate gets a reference to the given AuthenticatorProviderConfigurationUserNameTemplate and assigns it to the UserNameTemplate field.
func (o *AuthenticatorProviderConfiguration) SetUserNameTemplate(v AuthenticatorProviderConfigurationUserNameTemplate) {
	o.UserNameTemplate = &v
}

func (o AuthenticatorProviderConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AuthPort != nil {
		toSerialize["authPort"] = o.AuthPort
	}
	if o.HostName != nil {
		toSerialize["hostName"] = o.HostName
	}
	if o.InstanceId != nil {
		toSerialize["instanceId"] = o.InstanceId
	}
	if o.SharedSecret != nil {
		toSerialize["sharedSecret"] = o.SharedSecret
	}
	if o.UserNameTemplate != nil {
		toSerialize["userNameTemplate"] = o.UserNameTemplate
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *AuthenticatorProviderConfiguration) UnmarshalJSON(bytes []byte) (err error) {
	varAuthenticatorProviderConfiguration := _AuthenticatorProviderConfiguration{}

	err = json.Unmarshal(bytes, &varAuthenticatorProviderConfiguration)
	if err == nil {
		*o = AuthenticatorProviderConfiguration(varAuthenticatorProviderConfiguration)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "authPort")
		delete(additionalProperties, "hostName")
		delete(additionalProperties, "instanceId")
		delete(additionalProperties, "sharedSecret")
		delete(additionalProperties, "userNameTemplate")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableAuthenticatorProviderConfiguration struct {
	value *AuthenticatorProviderConfiguration
	isSet bool
}

func (v NullableAuthenticatorProviderConfiguration) Get() *AuthenticatorProviderConfiguration {
	return v.value
}

func (v *NullableAuthenticatorProviderConfiguration) Set(val *AuthenticatorProviderConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthenticatorProviderConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthenticatorProviderConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthenticatorProviderConfiguration(val *AuthenticatorProviderConfiguration) *NullableAuthenticatorProviderConfiguration {
	return &NullableAuthenticatorProviderConfiguration{value: val, isSet: true}
}

func (v NullableAuthenticatorProviderConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthenticatorProviderConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
