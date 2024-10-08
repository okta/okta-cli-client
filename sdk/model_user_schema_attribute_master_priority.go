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

// UserSchemaAttributeMasterPriority struct for UserSchemaAttributeMasterPriority
type UserSchemaAttributeMasterPriority struct {
	Type                 *string `json:"type,omitempty"`
	Value                *string `json:"value,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _UserSchemaAttributeMasterPriority UserSchemaAttributeMasterPriority

// NewUserSchemaAttributeMasterPriority instantiates a new UserSchemaAttributeMasterPriority object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserSchemaAttributeMasterPriority() *UserSchemaAttributeMasterPriority {
	this := UserSchemaAttributeMasterPriority{}
	return &this
}

// NewUserSchemaAttributeMasterPriorityWithDefaults instantiates a new UserSchemaAttributeMasterPriority object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserSchemaAttributeMasterPriorityWithDefaults() *UserSchemaAttributeMasterPriority {
	this := UserSchemaAttributeMasterPriority{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *UserSchemaAttributeMasterPriority) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserSchemaAttributeMasterPriority) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *UserSchemaAttributeMasterPriority) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *UserSchemaAttributeMasterPriority) SetType(v string) {
	o.Type = &v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *UserSchemaAttributeMasterPriority) GetValue() string {
	if o == nil || o.Value == nil {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserSchemaAttributeMasterPriority) GetValueOk() (*string, bool) {
	if o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *UserSchemaAttributeMasterPriority) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *UserSchemaAttributeMasterPriority) SetValue(v string) {
	o.Value = &v
}

func (o UserSchemaAttributeMasterPriority) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.Value != nil {
		toSerialize["value"] = o.Value
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *UserSchemaAttributeMasterPriority) UnmarshalJSON(bytes []byte) (err error) {
	varUserSchemaAttributeMasterPriority := _UserSchemaAttributeMasterPriority{}

	err = json.Unmarshal(bytes, &varUserSchemaAttributeMasterPriority)
	if err == nil {
		*o = UserSchemaAttributeMasterPriority(varUserSchemaAttributeMasterPriority)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "type")
		delete(additionalProperties, "value")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableUserSchemaAttributeMasterPriority struct {
	value *UserSchemaAttributeMasterPriority
	isSet bool
}

func (v NullableUserSchemaAttributeMasterPriority) Get() *UserSchemaAttributeMasterPriority {
	return v.value
}

func (v *NullableUserSchemaAttributeMasterPriority) Set(val *UserSchemaAttributeMasterPriority) {
	v.value = val
	v.isSet = true
}

func (v NullableUserSchemaAttributeMasterPriority) IsSet() bool {
	return v.isSet
}

func (v *NullableUserSchemaAttributeMasterPriority) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserSchemaAttributeMasterPriority(val *UserSchemaAttributeMasterPriority) *NullableUserSchemaAttributeMasterPriority {
	return &NullableUserSchemaAttributeMasterPriority{value: val, isSet: true}
}

func (v NullableUserSchemaAttributeMasterPriority) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserSchemaAttributeMasterPriority) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
