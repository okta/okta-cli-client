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

// PasswordDictionary struct for PasswordDictionary
type PasswordDictionary struct {
	Common               *PasswordDictionaryCommon `json:"common,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PasswordDictionary PasswordDictionary

// NewPasswordDictionary instantiates a new PasswordDictionary object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPasswordDictionary() *PasswordDictionary {
	this := PasswordDictionary{}
	return &this
}

// NewPasswordDictionaryWithDefaults instantiates a new PasswordDictionary object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPasswordDictionaryWithDefaults() *PasswordDictionary {
	this := PasswordDictionary{}
	return &this
}

// GetCommon returns the Common field value if set, zero value otherwise.
func (o *PasswordDictionary) GetCommon() PasswordDictionaryCommon {
	if o == nil || o.Common == nil {
		var ret PasswordDictionaryCommon
		return ret
	}
	return *o.Common
}

// GetCommonOk returns a tuple with the Common field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PasswordDictionary) GetCommonOk() (*PasswordDictionaryCommon, bool) {
	if o == nil || o.Common == nil {
		return nil, false
	}
	return o.Common, true
}

// HasCommon returns a boolean if a field has been set.
func (o *PasswordDictionary) HasCommon() bool {
	if o != nil && o.Common != nil {
		return true
	}

	return false
}

// SetCommon gets a reference to the given PasswordDictionaryCommon and assigns it to the Common field.
func (o *PasswordDictionary) SetCommon(v PasswordDictionaryCommon) {
	o.Common = &v
}

func (o PasswordDictionary) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Common != nil {
		toSerialize["common"] = o.Common
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PasswordDictionary) UnmarshalJSON(bytes []byte) (err error) {
	varPasswordDictionary := _PasswordDictionary{}

	err = json.Unmarshal(bytes, &varPasswordDictionary)
	if err == nil {
		*o = PasswordDictionary(varPasswordDictionary)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "common")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullablePasswordDictionary struct {
	value *PasswordDictionary
	isSet bool
}

func (v NullablePasswordDictionary) Get() *PasswordDictionary {
	return v.value
}

func (v *NullablePasswordDictionary) Set(val *PasswordDictionary) {
	v.value = val
	v.isSet = true
}

func (v NullablePasswordDictionary) IsSet() bool {
	return v.isSet
}

func (v *NullablePasswordDictionary) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePasswordDictionary(val *PasswordDictionary) *NullablePasswordDictionary {
	return &NullablePasswordDictionary{value: val, isSet: true}
}

func (v NullablePasswordDictionary) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePasswordDictionary) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
