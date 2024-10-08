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

// PolicyContextDevice struct for PolicyContextDevice
type PolicyContextDevice struct {
	// The platform of the device, for example, IOS.
	Platform *string `json:"platform,omitempty"`
	// If the device is registered
	Registered *bool `json:"registered,omitempty"`
	// If the device is managed
	Managed              *bool `json:"managed,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PolicyContextDevice PolicyContextDevice

// NewPolicyContextDevice instantiates a new PolicyContextDevice object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPolicyContextDevice() *PolicyContextDevice {
	this := PolicyContextDevice{}
	return &this
}

// NewPolicyContextDeviceWithDefaults instantiates a new PolicyContextDevice object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPolicyContextDeviceWithDefaults() *PolicyContextDevice {
	this := PolicyContextDevice{}
	return &this
}

// GetPlatform returns the Platform field value if set, zero value otherwise.
func (o *PolicyContextDevice) GetPlatform() string {
	if o == nil || o.Platform == nil {
		var ret string
		return ret
	}
	return *o.Platform
}

// GetPlatformOk returns a tuple with the Platform field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyContextDevice) GetPlatformOk() (*string, bool) {
	if o == nil || o.Platform == nil {
		return nil, false
	}
	return o.Platform, true
}

// HasPlatform returns a boolean if a field has been set.
func (o *PolicyContextDevice) HasPlatform() bool {
	if o != nil && o.Platform != nil {
		return true
	}

	return false
}

// SetPlatform gets a reference to the given string and assigns it to the Platform field.
func (o *PolicyContextDevice) SetPlatform(v string) {
	o.Platform = &v
}

// GetRegistered returns the Registered field value if set, zero value otherwise.
func (o *PolicyContextDevice) GetRegistered() bool {
	if o == nil || o.Registered == nil {
		var ret bool
		return ret
	}
	return *o.Registered
}

// GetRegisteredOk returns a tuple with the Registered field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyContextDevice) GetRegisteredOk() (*bool, bool) {
	if o == nil || o.Registered == nil {
		return nil, false
	}
	return o.Registered, true
}

// HasRegistered returns a boolean if a field has been set.
func (o *PolicyContextDevice) HasRegistered() bool {
	if o != nil && o.Registered != nil {
		return true
	}

	return false
}

// SetRegistered gets a reference to the given bool and assigns it to the Registered field.
func (o *PolicyContextDevice) SetRegistered(v bool) {
	o.Registered = &v
}

// GetManaged returns the Managed field value if set, zero value otherwise.
func (o *PolicyContextDevice) GetManaged() bool {
	if o == nil || o.Managed == nil {
		var ret bool
		return ret
	}
	return *o.Managed
}

// GetManagedOk returns a tuple with the Managed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PolicyContextDevice) GetManagedOk() (*bool, bool) {
	if o == nil || o.Managed == nil {
		return nil, false
	}
	return o.Managed, true
}

// HasManaged returns a boolean if a field has been set.
func (o *PolicyContextDevice) HasManaged() bool {
	if o != nil && o.Managed != nil {
		return true
	}

	return false
}

// SetManaged gets a reference to the given bool and assigns it to the Managed field.
func (o *PolicyContextDevice) SetManaged(v bool) {
	o.Managed = &v
}

func (o PolicyContextDevice) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Platform != nil {
		toSerialize["platform"] = o.Platform
	}
	if o.Registered != nil {
		toSerialize["registered"] = o.Registered
	}
	if o.Managed != nil {
		toSerialize["managed"] = o.Managed
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PolicyContextDevice) UnmarshalJSON(bytes []byte) (err error) {
	varPolicyContextDevice := _PolicyContextDevice{}

	err = json.Unmarshal(bytes, &varPolicyContextDevice)
	if err == nil {
		*o = PolicyContextDevice(varPolicyContextDevice)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "platform")
		delete(additionalProperties, "registered")
		delete(additionalProperties, "managed")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullablePolicyContextDevice struct {
	value *PolicyContextDevice
	isSet bool
}

func (v NullablePolicyContextDevice) Get() *PolicyContextDevice {
	return v.value
}

func (v *NullablePolicyContextDevice) Set(val *PolicyContextDevice) {
	v.value = val
	v.isSet = true
}

func (v NullablePolicyContextDevice) IsSet() bool {
	return v.isSet
}

func (v *NullablePolicyContextDevice) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePolicyContextDevice(val *PolicyContextDevice) *NullablePolicyContextDevice {
	return &NullablePolicyContextDevice{value: val, isSet: true}
}

func (v NullablePolicyContextDevice) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePolicyContextDevice) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
