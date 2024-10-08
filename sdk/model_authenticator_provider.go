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

// AuthenticatorProvider struct for AuthenticatorProvider
type AuthenticatorProvider struct {
	Configuration        *AuthenticatorProviderConfiguration `json:"configuration,omitempty"`
	Type                 *string                             `json:"type,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _AuthenticatorProvider AuthenticatorProvider

// NewAuthenticatorProvider instantiates a new AuthenticatorProvider object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthenticatorProvider() *AuthenticatorProvider {
	this := AuthenticatorProvider{}
	return &this
}

// NewAuthenticatorProviderWithDefaults instantiates a new AuthenticatorProvider object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthenticatorProviderWithDefaults() *AuthenticatorProvider {
	this := AuthenticatorProvider{}
	return &this
}

// GetConfiguration returns the Configuration field value if set, zero value otherwise.
func (o *AuthenticatorProvider) GetConfiguration() AuthenticatorProviderConfiguration {
	if o == nil || o.Configuration == nil {
		var ret AuthenticatorProviderConfiguration
		return ret
	}
	return *o.Configuration
}

// GetConfigurationOk returns a tuple with the Configuration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProvider) GetConfigurationOk() (*AuthenticatorProviderConfiguration, bool) {
	if o == nil || o.Configuration == nil {
		return nil, false
	}
	return o.Configuration, true
}

// HasConfiguration returns a boolean if a field has been set.
func (o *AuthenticatorProvider) HasConfiguration() bool {
	if o != nil && o.Configuration != nil {
		return true
	}

	return false
}

// SetConfiguration gets a reference to the given AuthenticatorProviderConfiguration and assigns it to the Configuration field.
func (o *AuthenticatorProvider) SetConfiguration(v AuthenticatorProviderConfiguration) {
	o.Configuration = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *AuthenticatorProvider) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProvider) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *AuthenticatorProvider) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *AuthenticatorProvider) SetType(v string) {
	o.Type = &v
}

func (o AuthenticatorProvider) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Configuration != nil {
		toSerialize["configuration"] = o.Configuration
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *AuthenticatorProvider) UnmarshalJSON(bytes []byte) (err error) {
	varAuthenticatorProvider := _AuthenticatorProvider{}

	err = json.Unmarshal(bytes, &varAuthenticatorProvider)
	if err == nil {
		*o = AuthenticatorProvider(varAuthenticatorProvider)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "configuration")
		delete(additionalProperties, "type")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableAuthenticatorProvider struct {
	value *AuthenticatorProvider
	isSet bool
}

func (v NullableAuthenticatorProvider) Get() *AuthenticatorProvider {
	return v.value
}

func (v *NullableAuthenticatorProvider) Set(val *AuthenticatorProvider) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthenticatorProvider) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthenticatorProvider) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthenticatorProvider(val *AuthenticatorProvider) *NullableAuthenticatorProvider {
	return &NullableAuthenticatorProvider{value: val, isSet: true}
}

func (v NullableAuthenticatorProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthenticatorProvider) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
