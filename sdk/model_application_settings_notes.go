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

// ApplicationSettingsNotes struct for ApplicationSettingsNotes
type ApplicationSettingsNotes struct {
	Admin                *string `json:"admin,omitempty"`
	Enduser              *string `json:"enduser,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ApplicationSettingsNotes ApplicationSettingsNotes

// NewApplicationSettingsNotes instantiates a new ApplicationSettingsNotes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApplicationSettingsNotes() *ApplicationSettingsNotes {
	this := ApplicationSettingsNotes{}
	return &this
}

// NewApplicationSettingsNotesWithDefaults instantiates a new ApplicationSettingsNotes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApplicationSettingsNotesWithDefaults() *ApplicationSettingsNotes {
	this := ApplicationSettingsNotes{}
	return &this
}

// GetAdmin returns the Admin field value if set, zero value otherwise.
func (o *ApplicationSettingsNotes) GetAdmin() string {
	if o == nil || o.Admin == nil {
		var ret string
		return ret
	}
	return *o.Admin
}

// GetAdminOk returns a tuple with the Admin field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationSettingsNotes) GetAdminOk() (*string, bool) {
	if o == nil || o.Admin == nil {
		return nil, false
	}
	return o.Admin, true
}

// HasAdmin returns a boolean if a field has been set.
func (o *ApplicationSettingsNotes) HasAdmin() bool {
	if o != nil && o.Admin != nil {
		return true
	}

	return false
}

// SetAdmin gets a reference to the given string and assigns it to the Admin field.
func (o *ApplicationSettingsNotes) SetAdmin(v string) {
	o.Admin = &v
}

// GetEnduser returns the Enduser field value if set, zero value otherwise.
func (o *ApplicationSettingsNotes) GetEnduser() string {
	if o == nil || o.Enduser == nil {
		var ret string
		return ret
	}
	return *o.Enduser
}

// GetEnduserOk returns a tuple with the Enduser field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApplicationSettingsNotes) GetEnduserOk() (*string, bool) {
	if o == nil || o.Enduser == nil {
		return nil, false
	}
	return o.Enduser, true
}

// HasEnduser returns a boolean if a field has been set.
func (o *ApplicationSettingsNotes) HasEnduser() bool {
	if o != nil && o.Enduser != nil {
		return true
	}

	return false
}

// SetEnduser gets a reference to the given string and assigns it to the Enduser field.
func (o *ApplicationSettingsNotes) SetEnduser(v string) {
	o.Enduser = &v
}

func (o ApplicationSettingsNotes) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Admin != nil {
		toSerialize["admin"] = o.Admin
	}
	if o.Enduser != nil {
		toSerialize["enduser"] = o.Enduser
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *ApplicationSettingsNotes) UnmarshalJSON(bytes []byte) (err error) {
	varApplicationSettingsNotes := _ApplicationSettingsNotes{}

	err = json.Unmarshal(bytes, &varApplicationSettingsNotes)
	if err == nil {
		*o = ApplicationSettingsNotes(varApplicationSettingsNotes)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "admin")
		delete(additionalProperties, "enduser")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableApplicationSettingsNotes struct {
	value *ApplicationSettingsNotes
	isSet bool
}

func (v NullableApplicationSettingsNotes) Get() *ApplicationSettingsNotes {
	return v.value
}

func (v *NullableApplicationSettingsNotes) Set(val *ApplicationSettingsNotes) {
	v.value = val
	v.isSet = true
}

func (v NullableApplicationSettingsNotes) IsSet() bool {
	return v.isSet
}

func (v *NullableApplicationSettingsNotes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApplicationSettingsNotes(val *ApplicationSettingsNotes) *NullableApplicationSettingsNotes {
	return &NullableApplicationSettingsNotes{value: val, isSet: true}
}

func (v NullableApplicationSettingsNotes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApplicationSettingsNotes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
