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

// NetworkZoneLinks struct for NetworkZoneLinks
type NetworkZoneLinks struct {
	Self                 *HrefObjectSelfLink `json:"self,omitempty"`
	Deactivate           *HrefObject         `json:"deactivate,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _NetworkZoneLinks NetworkZoneLinks

// NewNetworkZoneLinks instantiates a new NetworkZoneLinks object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNetworkZoneLinks() *NetworkZoneLinks {
	this := NetworkZoneLinks{}
	return &this
}

// NewNetworkZoneLinksWithDefaults instantiates a new NetworkZoneLinks object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNetworkZoneLinksWithDefaults() *NetworkZoneLinks {
	this := NetworkZoneLinks{}
	return &this
}

// GetSelf returns the Self field value if set, zero value otherwise.
func (o *NetworkZoneLinks) GetSelf() HrefObjectSelfLink {
	if o == nil || o.Self == nil {
		var ret HrefObjectSelfLink
		return ret
	}
	return *o.Self
}

// GetSelfOk returns a tuple with the Self field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkZoneLinks) GetSelfOk() (*HrefObjectSelfLink, bool) {
	if o == nil || o.Self == nil {
		return nil, false
	}
	return o.Self, true
}

// HasSelf returns a boolean if a field has been set.
func (o *NetworkZoneLinks) HasSelf() bool {
	if o != nil && o.Self != nil {
		return true
	}

	return false
}

// SetSelf gets a reference to the given HrefObjectSelfLink and assigns it to the Self field.
func (o *NetworkZoneLinks) SetSelf(v HrefObjectSelfLink) {
	o.Self = &v
}

// GetDeactivate returns the Deactivate field value if set, zero value otherwise.
func (o *NetworkZoneLinks) GetDeactivate() HrefObject {
	if o == nil || o.Deactivate == nil {
		var ret HrefObject
		return ret
	}
	return *o.Deactivate
}

// GetDeactivateOk returns a tuple with the Deactivate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkZoneLinks) GetDeactivateOk() (*HrefObject, bool) {
	if o == nil || o.Deactivate == nil {
		return nil, false
	}
	return o.Deactivate, true
}

// HasDeactivate returns a boolean if a field has been set.
func (o *NetworkZoneLinks) HasDeactivate() bool {
	if o != nil && o.Deactivate != nil {
		return true
	}

	return false
}

// SetDeactivate gets a reference to the given HrefObject and assigns it to the Deactivate field.
func (o *NetworkZoneLinks) SetDeactivate(v HrefObject) {
	o.Deactivate = &v
}

func (o NetworkZoneLinks) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Self != nil {
		toSerialize["self"] = o.Self
	}
	if o.Deactivate != nil {
		toSerialize["deactivate"] = o.Deactivate
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *NetworkZoneLinks) UnmarshalJSON(bytes []byte) (err error) {
	varNetworkZoneLinks := _NetworkZoneLinks{}

	err = json.Unmarshal(bytes, &varNetworkZoneLinks)
	if err == nil {
		*o = NetworkZoneLinks(varNetworkZoneLinks)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "self")
		delete(additionalProperties, "deactivate")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableNetworkZoneLinks struct {
	value *NetworkZoneLinks
	isSet bool
}

func (v NullableNetworkZoneLinks) Get() *NetworkZoneLinks {
	return v.value
}

func (v *NullableNetworkZoneLinks) Set(val *NetworkZoneLinks) {
	v.value = val
	v.isSet = true
}

func (v NullableNetworkZoneLinks) IsSet() bool {
	return v.isSet
}

func (v *NullableNetworkZoneLinks) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNetworkZoneLinks(val *NetworkZoneLinks) *NullableNetworkZoneLinks {
	return &NullableNetworkZoneLinks{value: val, isSet: true}
}

func (v NullableNetworkZoneLinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNetworkZoneLinks) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
