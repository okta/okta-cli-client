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

// ResourceSetResourcePatchRequest struct for ResourceSetResourcePatchRequest
type ResourceSetResourcePatchRequest struct {
	Additions            []string `json:"additions,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ResourceSetResourcePatchRequest ResourceSetResourcePatchRequest

// NewResourceSetResourcePatchRequest instantiates a new ResourceSetResourcePatchRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewResourceSetResourcePatchRequest() *ResourceSetResourcePatchRequest {
	this := ResourceSetResourcePatchRequest{}
	return &this
}

// NewResourceSetResourcePatchRequestWithDefaults instantiates a new ResourceSetResourcePatchRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewResourceSetResourcePatchRequestWithDefaults() *ResourceSetResourcePatchRequest {
	this := ResourceSetResourcePatchRequest{}
	return &this
}

// GetAdditions returns the Additions field value if set, zero value otherwise.
func (o *ResourceSetResourcePatchRequest) GetAdditions() []string {
	if o == nil || o.Additions == nil {
		var ret []string
		return ret
	}
	return o.Additions
}

// GetAdditionsOk returns a tuple with the Additions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResourceSetResourcePatchRequest) GetAdditionsOk() ([]string, bool) {
	if o == nil || o.Additions == nil {
		return nil, false
	}
	return o.Additions, true
}

// HasAdditions returns a boolean if a field has been set.
func (o *ResourceSetResourcePatchRequest) HasAdditions() bool {
	if o != nil && o.Additions != nil {
		return true
	}

	return false
}

// SetAdditions gets a reference to the given []string and assigns it to the Additions field.
func (o *ResourceSetResourcePatchRequest) SetAdditions(v []string) {
	o.Additions = v
}

func (o ResourceSetResourcePatchRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Additions != nil {
		toSerialize["additions"] = o.Additions
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *ResourceSetResourcePatchRequest) UnmarshalJSON(bytes []byte) (err error) {
	varResourceSetResourcePatchRequest := _ResourceSetResourcePatchRequest{}

	err = json.Unmarshal(bytes, &varResourceSetResourcePatchRequest)
	if err == nil {
		*o = ResourceSetResourcePatchRequest(varResourceSetResourcePatchRequest)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "additions")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableResourceSetResourcePatchRequest struct {
	value *ResourceSetResourcePatchRequest
	isSet bool
}

func (v NullableResourceSetResourcePatchRequest) Get() *ResourceSetResourcePatchRequest {
	return v.value
}

func (v *NullableResourceSetResourcePatchRequest) Set(val *ResourceSetResourcePatchRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableResourceSetResourcePatchRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableResourceSetResourcePatchRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResourceSetResourcePatchRequest(val *ResourceSetResourcePatchRequest) *NullableResourceSetResourcePatchRequest {
	return &NullableResourceSetResourcePatchRequest{value: val, isSet: true}
}

func (v NullableResourceSetResourcePatchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResourceSetResourcePatchRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
