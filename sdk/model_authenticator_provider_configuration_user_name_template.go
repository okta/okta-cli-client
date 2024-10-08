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

// AuthenticatorProviderConfigurationUserNameTemplate struct for AuthenticatorProviderConfigurationUserNameTemplate
type AuthenticatorProviderConfigurationUserNameTemplate struct {
	Template             *string `json:"template,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _AuthenticatorProviderConfigurationUserNameTemplate AuthenticatorProviderConfigurationUserNameTemplate

// NewAuthenticatorProviderConfigurationUserNameTemplate instantiates a new AuthenticatorProviderConfigurationUserNameTemplate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthenticatorProviderConfigurationUserNameTemplate() *AuthenticatorProviderConfigurationUserNameTemplate {
	this := AuthenticatorProviderConfigurationUserNameTemplate{}
	return &this
}

// NewAuthenticatorProviderConfigurationUserNameTemplateWithDefaults instantiates a new AuthenticatorProviderConfigurationUserNameTemplate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthenticatorProviderConfigurationUserNameTemplateWithDefaults() *AuthenticatorProviderConfigurationUserNameTemplate {
	this := AuthenticatorProviderConfigurationUserNameTemplate{}
	return &this
}

// GetTemplate returns the Template field value if set, zero value otherwise.
func (o *AuthenticatorProviderConfigurationUserNameTemplate) GetTemplate() string {
	if o == nil || o.Template == nil {
		var ret string
		return ret
	}
	return *o.Template
}

// GetTemplateOk returns a tuple with the Template field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatorProviderConfigurationUserNameTemplate) GetTemplateOk() (*string, bool) {
	if o == nil || o.Template == nil {
		return nil, false
	}
	return o.Template, true
}

// HasTemplate returns a boolean if a field has been set.
func (o *AuthenticatorProviderConfigurationUserNameTemplate) HasTemplate() bool {
	if o != nil && o.Template != nil {
		return true
	}

	return false
}

// SetTemplate gets a reference to the given string and assigns it to the Template field.
func (o *AuthenticatorProviderConfigurationUserNameTemplate) SetTemplate(v string) {
	o.Template = &v
}

func (o AuthenticatorProviderConfigurationUserNameTemplate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Template != nil {
		toSerialize["template"] = o.Template
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *AuthenticatorProviderConfigurationUserNameTemplate) UnmarshalJSON(bytes []byte) (err error) {
	varAuthenticatorProviderConfigurationUserNameTemplate := _AuthenticatorProviderConfigurationUserNameTemplate{}

	err = json.Unmarshal(bytes, &varAuthenticatorProviderConfigurationUserNameTemplate)
	if err == nil {
		*o = AuthenticatorProviderConfigurationUserNameTemplate(varAuthenticatorProviderConfigurationUserNameTemplate)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "template")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableAuthenticatorProviderConfigurationUserNameTemplate struct {
	value *AuthenticatorProviderConfigurationUserNameTemplate
	isSet bool
}

func (v NullableAuthenticatorProviderConfigurationUserNameTemplate) Get() *AuthenticatorProviderConfigurationUserNameTemplate {
	return v.value
}

func (v *NullableAuthenticatorProviderConfigurationUserNameTemplate) Set(val *AuthenticatorProviderConfigurationUserNameTemplate) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthenticatorProviderConfigurationUserNameTemplate) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthenticatorProviderConfigurationUserNameTemplate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthenticatorProviderConfigurationUserNameTemplate(val *AuthenticatorProviderConfigurationUserNameTemplate) *NullableAuthenticatorProviderConfigurationUserNameTemplate {
	return &NullableAuthenticatorProviderConfigurationUserNameTemplate{value: val, isSet: true}
}

func (v NullableAuthenticatorProviderConfigurationUserNameTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthenticatorProviderConfigurationUserNameTemplate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
