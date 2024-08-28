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
	"time"
)

// Permission struct for Permission
type Permission struct {
	// Conditions for further restricting a permission
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	// Timestamp when the role was created
	Created *time.Time `json:"created,omitempty"`
	// The permission type
	Label *string `json:"label,omitempty"`
	// Timestamp when the role was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links *PermissionLinks `json:"_links,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Permission Permission

// NewPermission instantiates a new Permission object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPermission() *Permission {
	this := Permission{}
	return &this
}

// NewPermissionWithDefaults instantiates a new Permission object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPermissionWithDefaults() *Permission {
	this := Permission{}
	return &this
}

// GetConditions returns the Conditions field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Permission) GetConditions() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}
	return o.Conditions
}

// GetConditionsOk returns a tuple with the Conditions field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Permission) GetConditionsOk() (map[string]interface{}, bool) {
	if o == nil || o.Conditions == nil {
		return nil, false
	}
	return o.Conditions, true
}

// HasConditions returns a boolean if a field has been set.
func (o *Permission) HasConditions() bool {
	if o != nil && o.Conditions != nil {
		return true
	}

	return false
}

// SetConditions gets a reference to the given map[string]interface{} and assigns it to the Conditions field.
func (o *Permission) SetConditions(v map[string]interface{}) {
	o.Conditions = v
}

// GetCreated returns the Created field value if set, zero value otherwise.
func (o *Permission) GetCreated() time.Time {
	if o == nil || o.Created == nil {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Permission) GetCreatedOk() (*time.Time, bool) {
	if o == nil || o.Created == nil {
		return nil, false
	}
	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *Permission) HasCreated() bool {
	if o != nil && o.Created != nil {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *Permission) SetCreated(v time.Time) {
	o.Created = &v
}

// GetLabel returns the Label field value if set, zero value otherwise.
func (o *Permission) GetLabel() string {
	if o == nil || o.Label == nil {
		var ret string
		return ret
	}
	return *o.Label
}

// GetLabelOk returns a tuple with the Label field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Permission) GetLabelOk() (*string, bool) {
	if o == nil || o.Label == nil {
		return nil, false
	}
	return o.Label, true
}

// HasLabel returns a boolean if a field has been set.
func (o *Permission) HasLabel() bool {
	if o != nil && o.Label != nil {
		return true
	}

	return false
}

// SetLabel gets a reference to the given string and assigns it to the Label field.
func (o *Permission) SetLabel(v string) {
	o.Label = &v
}

// GetLastUpdated returns the LastUpdated field value if set, zero value otherwise.
func (o *Permission) GetLastUpdated() time.Time {
	if o == nil || o.LastUpdated == nil {
		var ret time.Time
		return ret
	}
	return *o.LastUpdated
}

// GetLastUpdatedOk returns a tuple with the LastUpdated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Permission) GetLastUpdatedOk() (*time.Time, bool) {
	if o == nil || o.LastUpdated == nil {
		return nil, false
	}
	return o.LastUpdated, true
}

// HasLastUpdated returns a boolean if a field has been set.
func (o *Permission) HasLastUpdated() bool {
	if o != nil && o.LastUpdated != nil {
		return true
	}

	return false
}

// SetLastUpdated gets a reference to the given time.Time and assigns it to the LastUpdated field.
func (o *Permission) SetLastUpdated(v time.Time) {
	o.LastUpdated = &v
}

// GetLinks returns the Links field value if set, zero value otherwise.
func (o *Permission) GetLinks() PermissionLinks {
	if o == nil || o.Links == nil {
		var ret PermissionLinks
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Permission) GetLinksOk() (*PermissionLinks, bool) {
	if o == nil || o.Links == nil {
		return nil, false
	}
	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *Permission) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

// SetLinks gets a reference to the given PermissionLinks and assigns it to the Links field.
func (o *Permission) SetLinks(v PermissionLinks) {
	o.Links = &v
}

func (o Permission) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Conditions != nil {
		toSerialize["conditions"] = o.Conditions
	}
	if o.Created != nil {
		toSerialize["created"] = o.Created
	}
	if o.Label != nil {
		toSerialize["label"] = o.Label
	}
	if o.LastUpdated != nil {
		toSerialize["lastUpdated"] = o.LastUpdated
	}
	if o.Links != nil {
		toSerialize["_links"] = o.Links
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *Permission) UnmarshalJSON(bytes []byte) (err error) {
	varPermission := _Permission{}

	err = json.Unmarshal(bytes, &varPermission)
	if err == nil {
		*o = Permission(varPermission)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "conditions")
		delete(additionalProperties, "created")
		delete(additionalProperties, "label")
		delete(additionalProperties, "lastUpdated")
		delete(additionalProperties, "_links")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullablePermission struct {
	value *Permission
	isSet bool
}

func (v NullablePermission) Get() *Permission {
	return v.value
}

func (v *NullablePermission) Set(val *Permission) {
	v.value = val
	v.isSet = true
}

func (v NullablePermission) IsSet() bool {
	return v.isSet
}

func (v *NullablePermission) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePermission(val *Permission) *NullablePermission {
	return &NullablePermission{value: val, isSet: true}
}

func (v NullablePermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePermission) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

