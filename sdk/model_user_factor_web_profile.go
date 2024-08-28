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

// UserFactorWebProfile struct for UserFactorWebProfile
type UserFactorWebProfile struct {
	// ID for the Factor credential
	CredentialId *string `json:"credentialId,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _UserFactorWebProfile UserFactorWebProfile

// NewUserFactorWebProfile instantiates a new UserFactorWebProfile object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserFactorWebProfile() *UserFactorWebProfile {
	this := UserFactorWebProfile{}
	return &this
}

// NewUserFactorWebProfileWithDefaults instantiates a new UserFactorWebProfile object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserFactorWebProfileWithDefaults() *UserFactorWebProfile {
	this := UserFactorWebProfile{}
	return &this
}

// GetCredentialId returns the CredentialId field value if set, zero value otherwise.
func (o *UserFactorWebProfile) GetCredentialId() string {
	if o == nil || o.CredentialId == nil {
		var ret string
		return ret
	}
	return *o.CredentialId
}

// GetCredentialIdOk returns a tuple with the CredentialId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserFactorWebProfile) GetCredentialIdOk() (*string, bool) {
	if o == nil || o.CredentialId == nil {
		return nil, false
	}
	return o.CredentialId, true
}

// HasCredentialId returns a boolean if a field has been set.
func (o *UserFactorWebProfile) HasCredentialId() bool {
	if o != nil && o.CredentialId != nil {
		return true
	}

	return false
}

// SetCredentialId gets a reference to the given string and assigns it to the CredentialId field.
func (o *UserFactorWebProfile) SetCredentialId(v string) {
	o.CredentialId = &v
}

func (o UserFactorWebProfile) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CredentialId != nil {
		toSerialize["credentialId"] = o.CredentialId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *UserFactorWebProfile) UnmarshalJSON(bytes []byte) (err error) {
	varUserFactorWebProfile := _UserFactorWebProfile{}

	err = json.Unmarshal(bytes, &varUserFactorWebProfile)
	if err == nil {
		*o = UserFactorWebProfile(varUserFactorWebProfile)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "credentialId")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableUserFactorWebProfile struct {
	value *UserFactorWebProfile
	isSet bool
}

func (v NullableUserFactorWebProfile) Get() *UserFactorWebProfile {
	return v.value
}

func (v *NullableUserFactorWebProfile) Set(val *UserFactorWebProfile) {
	v.value = val
	v.isSet = true
}

func (v NullableUserFactorWebProfile) IsSet() bool {
	return v.isSet
}

func (v *NullableUserFactorWebProfile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserFactorWebProfile(val *UserFactorWebProfile) *NullableUserFactorWebProfile {
	return &NullableUserFactorWebProfile{value: val, isSet: true}
}

func (v NullableUserFactorWebProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserFactorWebProfile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

