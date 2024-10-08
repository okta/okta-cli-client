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

// NetworkZoneLocation struct for NetworkZoneLocation
type NetworkZoneLocation struct {
	// Format of the country value: length 2 [ISO-3166-1](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) country code. Do not use continent codes as they are treated as generic codes for undesignated countries.
	Country *string `json:"country,omitempty"`
	// Format of the region value (optional): region code [ISO-3166-2](https://en.wikipedia.org/wiki/ISO_3166-2) appended to country code (`countryCode-regionCode`), or `null` if empty. Do not use continent codes as they are treated as generic codes for undesignated regions.
	Region               *string `json:"region,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _NetworkZoneLocation NetworkZoneLocation

// NewNetworkZoneLocation instantiates a new NetworkZoneLocation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNetworkZoneLocation() *NetworkZoneLocation {
	this := NetworkZoneLocation{}
	return &this
}

// NewNetworkZoneLocationWithDefaults instantiates a new NetworkZoneLocation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNetworkZoneLocationWithDefaults() *NetworkZoneLocation {
	this := NetworkZoneLocation{}
	return &this
}

// GetCountry returns the Country field value if set, zero value otherwise.
func (o *NetworkZoneLocation) GetCountry() string {
	if o == nil || o.Country == nil {
		var ret string
		return ret
	}
	return *o.Country
}

// GetCountryOk returns a tuple with the Country field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkZoneLocation) GetCountryOk() (*string, bool) {
	if o == nil || o.Country == nil {
		return nil, false
	}
	return o.Country, true
}

// HasCountry returns a boolean if a field has been set.
func (o *NetworkZoneLocation) HasCountry() bool {
	if o != nil && o.Country != nil {
		return true
	}

	return false
}

// SetCountry gets a reference to the given string and assigns it to the Country field.
func (o *NetworkZoneLocation) SetCountry(v string) {
	o.Country = &v
}

// GetRegion returns the Region field value if set, zero value otherwise.
func (o *NetworkZoneLocation) GetRegion() string {
	if o == nil || o.Region == nil {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkZoneLocation) GetRegionOk() (*string, bool) {
	if o == nil || o.Region == nil {
		return nil, false
	}
	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *NetworkZoneLocation) HasRegion() bool {
	if o != nil && o.Region != nil {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *NetworkZoneLocation) SetRegion(v string) {
	o.Region = &v
}

func (o NetworkZoneLocation) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Country != nil {
		toSerialize["country"] = o.Country
	}
	if o.Region != nil {
		toSerialize["region"] = o.Region
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *NetworkZoneLocation) UnmarshalJSON(bytes []byte) (err error) {
	varNetworkZoneLocation := _NetworkZoneLocation{}

	err = json.Unmarshal(bytes, &varNetworkZoneLocation)
	if err == nil {
		*o = NetworkZoneLocation(varNetworkZoneLocation)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "country")
		delete(additionalProperties, "region")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableNetworkZoneLocation struct {
	value *NetworkZoneLocation
	isSet bool
}

func (v NullableNetworkZoneLocation) Get() *NetworkZoneLocation {
	return v.value
}

func (v *NullableNetworkZoneLocation) Set(val *NetworkZoneLocation) {
	v.value = val
	v.isSet = true
}

func (v NullableNetworkZoneLocation) IsSet() bool {
	return v.isSet
}

func (v *NullableNetworkZoneLocation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNetworkZoneLocation(val *NetworkZoneLocation) *NullableNetworkZoneLocation {
	return &NullableNetworkZoneLocation{value: val, isSet: true}
}

func (v NullableNetworkZoneLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNetworkZoneLocation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
