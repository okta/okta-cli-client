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
	"reflect"
	"strings"
)

// UserFactorTOTP struct for UserFactorTOTP
type UserFactorTOTP struct {
	UserFactor
	Profile *UserFactorTOTPProfile `json:"profile,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _UserFactorTOTP UserFactorTOTP

// NewUserFactorTOTP instantiates a new UserFactorTOTP object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserFactorTOTP() *UserFactorTOTP {
	this := UserFactorTOTP{}
	return &this
}

// NewUserFactorTOTPWithDefaults instantiates a new UserFactorTOTP object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserFactorTOTPWithDefaults() *UserFactorTOTP {
	this := UserFactorTOTP{}
	return &this
}

// GetProfile returns the Profile field value if set, zero value otherwise.
func (o *UserFactorTOTP) GetProfile() UserFactorTOTPProfile {
	if o == nil || o.Profile == nil {
		var ret UserFactorTOTPProfile
		return ret
	}
	return *o.Profile
}

// GetProfileOk returns a tuple with the Profile field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserFactorTOTP) GetProfileOk() (*UserFactorTOTPProfile, bool) {
	if o == nil || o.Profile == nil {
		return nil, false
	}
	return o.Profile, true
}

// HasProfile returns a boolean if a field has been set.
func (o *UserFactorTOTP) HasProfile() bool {
	if o != nil && o.Profile != nil {
		return true
	}

	return false
}

// SetProfile gets a reference to the given UserFactorTOTPProfile and assigns it to the Profile field.
func (o *UserFactorTOTP) SetProfile(v UserFactorTOTPProfile) {
	o.Profile = &v
}

func (o UserFactorTOTP) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	serializedUserFactor, errUserFactor := json.Marshal(o.UserFactor)
	if errUserFactor != nil {
		return []byte{}, errUserFactor
	}
	errUserFactor = json.Unmarshal([]byte(serializedUserFactor), &toSerialize)
	if errUserFactor != nil {
		return []byte{}, errUserFactor
	}
	if o.Profile != nil {
		toSerialize["profile"] = o.Profile
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *UserFactorTOTP) UnmarshalJSON(bytes []byte) (err error) {
	type UserFactorTOTPWithoutEmbeddedStruct struct {
		Profile *UserFactorTOTPProfile `json:"profile,omitempty"`
	}

	varUserFactorTOTPWithoutEmbeddedStruct := UserFactorTOTPWithoutEmbeddedStruct{}

	err = json.Unmarshal(bytes, &varUserFactorTOTPWithoutEmbeddedStruct)
	if err == nil {
		varUserFactorTOTP := _UserFactorTOTP{}
		varUserFactorTOTP.Profile = varUserFactorTOTPWithoutEmbeddedStruct.Profile
		*o = UserFactorTOTP(varUserFactorTOTP)
	} else {
		return err
	}

	varUserFactorTOTP := _UserFactorTOTP{}

	err = json.Unmarshal(bytes, &varUserFactorTOTP)
	if err == nil {
		o.UserFactor = varUserFactorTOTP.UserFactor
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "profile")

		// remove fields from embedded structs
		reflectUserFactor := reflect.ValueOf(o.UserFactor)
		for i := 0; i < reflectUserFactor.Type().NumField(); i++ {
			t := reflectUserFactor.Type().Field(i)

			if jsonTag := t.Tag.Get("json"); jsonTag != "" {
				fieldName := ""
				if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
					fieldName = jsonTag[:commaIdx]
				} else {
					fieldName = jsonTag
				}
				if fieldName != "AdditionalProperties" {
					delete(additionalProperties, fieldName)
				}
			}
		}

		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableUserFactorTOTP struct {
	value *UserFactorTOTP
	isSet bool
}

func (v NullableUserFactorTOTP) Get() *UserFactorTOTP {
	return v.value
}

func (v *NullableUserFactorTOTP) Set(val *UserFactorTOTP) {
	v.value = val
	v.isSet = true
}

func (v NullableUserFactorTOTP) IsSet() bool {
	return v.isSet
}

func (v *NullableUserFactorTOTP) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserFactorTOTP(val *UserFactorTOTP) *NullableUserFactorTOTP {
	return &NullableUserFactorTOTP{value: val, isSet: true}
}

func (v NullableUserFactorTOTP) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserFactorTOTP) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
