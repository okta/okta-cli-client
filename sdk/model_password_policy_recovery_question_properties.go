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

// PasswordPolicyRecoveryQuestionProperties struct for PasswordPolicyRecoveryQuestionProperties
type PasswordPolicyRecoveryQuestionProperties struct {
	Complexity           *PasswordPolicyRecoveryQuestionComplexity `json:"complexity,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PasswordPolicyRecoveryQuestionProperties PasswordPolicyRecoveryQuestionProperties

// NewPasswordPolicyRecoveryQuestionProperties instantiates a new PasswordPolicyRecoveryQuestionProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPasswordPolicyRecoveryQuestionProperties() *PasswordPolicyRecoveryQuestionProperties {
	this := PasswordPolicyRecoveryQuestionProperties{}
	return &this
}

// NewPasswordPolicyRecoveryQuestionPropertiesWithDefaults instantiates a new PasswordPolicyRecoveryQuestionProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPasswordPolicyRecoveryQuestionPropertiesWithDefaults() *PasswordPolicyRecoveryQuestionProperties {
	this := PasswordPolicyRecoveryQuestionProperties{}
	return &this
}

// GetComplexity returns the Complexity field value if set, zero value otherwise.
func (o *PasswordPolicyRecoveryQuestionProperties) GetComplexity() PasswordPolicyRecoveryQuestionComplexity {
	if o == nil || o.Complexity == nil {
		var ret PasswordPolicyRecoveryQuestionComplexity
		return ret
	}
	return *o.Complexity
}

// GetComplexityOk returns a tuple with the Complexity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PasswordPolicyRecoveryQuestionProperties) GetComplexityOk() (*PasswordPolicyRecoveryQuestionComplexity, bool) {
	if o == nil || o.Complexity == nil {
		return nil, false
	}
	return o.Complexity, true
}

// HasComplexity returns a boolean if a field has been set.
func (o *PasswordPolicyRecoveryQuestionProperties) HasComplexity() bool {
	if o != nil && o.Complexity != nil {
		return true
	}

	return false
}

// SetComplexity gets a reference to the given PasswordPolicyRecoveryQuestionComplexity and assigns it to the Complexity field.
func (o *PasswordPolicyRecoveryQuestionProperties) SetComplexity(v PasswordPolicyRecoveryQuestionComplexity) {
	o.Complexity = &v
}

func (o PasswordPolicyRecoveryQuestionProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Complexity != nil {
		toSerialize["complexity"] = o.Complexity
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PasswordPolicyRecoveryQuestionProperties) UnmarshalJSON(bytes []byte) (err error) {
	varPasswordPolicyRecoveryQuestionProperties := _PasswordPolicyRecoveryQuestionProperties{}

	err = json.Unmarshal(bytes, &varPasswordPolicyRecoveryQuestionProperties)
	if err == nil {
		*o = PasswordPolicyRecoveryQuestionProperties(varPasswordPolicyRecoveryQuestionProperties)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "complexity")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullablePasswordPolicyRecoveryQuestionProperties struct {
	value *PasswordPolicyRecoveryQuestionProperties
	isSet bool
}

func (v NullablePasswordPolicyRecoveryQuestionProperties) Get() *PasswordPolicyRecoveryQuestionProperties {
	return v.value
}

func (v *NullablePasswordPolicyRecoveryQuestionProperties) Set(val *PasswordPolicyRecoveryQuestionProperties) {
	v.value = val
	v.isSet = true
}

func (v NullablePasswordPolicyRecoveryQuestionProperties) IsSet() bool {
	return v.isSet
}

func (v *NullablePasswordPolicyRecoveryQuestionProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePasswordPolicyRecoveryQuestionProperties(val *PasswordPolicyRecoveryQuestionProperties) *NullablePasswordPolicyRecoveryQuestionProperties {
	return &NullablePasswordPolicyRecoveryQuestionProperties{value: val, isSet: true}
}

func (v NullablePasswordPolicyRecoveryQuestionProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePasswordPolicyRecoveryQuestionProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
