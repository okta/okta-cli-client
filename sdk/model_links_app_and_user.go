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

// LinksAppAndUser Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification. This object is used for dynamic discovery of resources related to the Application User.
type LinksAppAndUser struct {
	App                  *HrefObjectAppLink    `json:"app,omitempty"`
	Group                *LinksAppAndUserGroup `json:"group,omitempty"`
	User                 *HrefObjectUserLink   `json:"user,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _LinksAppAndUser LinksAppAndUser

// NewLinksAppAndUser instantiates a new LinksAppAndUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinksAppAndUser() *LinksAppAndUser {
	this := LinksAppAndUser{}
	return &this
}

// NewLinksAppAndUserWithDefaults instantiates a new LinksAppAndUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinksAppAndUserWithDefaults() *LinksAppAndUser {
	this := LinksAppAndUser{}
	return &this
}

// GetApp returns the App field value if set, zero value otherwise.
func (o *LinksAppAndUser) GetApp() HrefObjectAppLink {
	if o == nil || o.App == nil {
		var ret HrefObjectAppLink
		return ret
	}
	return *o.App
}

// GetAppOk returns a tuple with the App field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinksAppAndUser) GetAppOk() (*HrefObjectAppLink, bool) {
	if o == nil || o.App == nil {
		return nil, false
	}
	return o.App, true
}

// HasApp returns a boolean if a field has been set.
func (o *LinksAppAndUser) HasApp() bool {
	if o != nil && o.App != nil {
		return true
	}

	return false
}

// SetApp gets a reference to the given HrefObjectAppLink and assigns it to the App field.
func (o *LinksAppAndUser) SetApp(v HrefObjectAppLink) {
	o.App = &v
}

// GetGroup returns the Group field value if set, zero value otherwise.
func (o *LinksAppAndUser) GetGroup() LinksAppAndUserGroup {
	if o == nil || o.Group == nil {
		var ret LinksAppAndUserGroup
		return ret
	}
	return *o.Group
}

// GetGroupOk returns a tuple with the Group field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinksAppAndUser) GetGroupOk() (*LinksAppAndUserGroup, bool) {
	if o == nil || o.Group == nil {
		return nil, false
	}
	return o.Group, true
}

// HasGroup returns a boolean if a field has been set.
func (o *LinksAppAndUser) HasGroup() bool {
	if o != nil && o.Group != nil {
		return true
	}

	return false
}

// SetGroup gets a reference to the given LinksAppAndUserGroup and assigns it to the Group field.
func (o *LinksAppAndUser) SetGroup(v LinksAppAndUserGroup) {
	o.Group = &v
}

// GetUser returns the User field value if set, zero value otherwise.
func (o *LinksAppAndUser) GetUser() HrefObjectUserLink {
	if o == nil || o.User == nil {
		var ret HrefObjectUserLink
		return ret
	}
	return *o.User
}

// GetUserOk returns a tuple with the User field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinksAppAndUser) GetUserOk() (*HrefObjectUserLink, bool) {
	if o == nil || o.User == nil {
		return nil, false
	}
	return o.User, true
}

// HasUser returns a boolean if a field has been set.
func (o *LinksAppAndUser) HasUser() bool {
	if o != nil && o.User != nil {
		return true
	}

	return false
}

// SetUser gets a reference to the given HrefObjectUserLink and assigns it to the User field.
func (o *LinksAppAndUser) SetUser(v HrefObjectUserLink) {
	o.User = &v
}

func (o LinksAppAndUser) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.App != nil {
		toSerialize["app"] = o.App
	}
	if o.Group != nil {
		toSerialize["group"] = o.Group
	}
	if o.User != nil {
		toSerialize["user"] = o.User
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *LinksAppAndUser) UnmarshalJSON(bytes []byte) (err error) {
	varLinksAppAndUser := _LinksAppAndUser{}

	err = json.Unmarshal(bytes, &varLinksAppAndUser)
	if err == nil {
		*o = LinksAppAndUser(varLinksAppAndUser)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "app")
		delete(additionalProperties, "group")
		delete(additionalProperties, "user")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableLinksAppAndUser struct {
	value *LinksAppAndUser
	isSet bool
}

func (v NullableLinksAppAndUser) Get() *LinksAppAndUser {
	return v.value
}

func (v *NullableLinksAppAndUser) Set(val *LinksAppAndUser) {
	v.value = val
	v.isSet = true
}

func (v NullableLinksAppAndUser) IsSet() bool {
	return v.isSet
}

func (v *NullableLinksAppAndUser) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinksAppAndUser(val *LinksAppAndUser) *NullableLinksAppAndUser {
	return &NullableLinksAppAndUser{value: val, isSet: true}
}

func (v NullableLinksAppAndUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinksAppAndUser) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
