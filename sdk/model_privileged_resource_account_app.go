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

// PrivilegedResourceAccountApp struct for PrivilegedResourceAccountApp
type PrivilegedResourceAccountApp struct {
	PrivilegedResource
	// The application ID associated with the privileged account
	ContainerId string `json:"containerId"`
	Credentials *PrivilegedResourceCredentials `json:"credentials,omitempty"`
	// Human-readable name of the container that owns the privileged resource
	ContainerDisplayName *string `json:"containerDisplayName,omitempty"`
	Links *AppLink `json:"_links,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PrivilegedResourceAccountApp PrivilegedResourceAccountApp

// NewPrivilegedResourceAccountApp instantiates a new PrivilegedResourceAccountApp object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivilegedResourceAccountApp(containerId string) *PrivilegedResourceAccountApp {
	this := PrivilegedResourceAccountApp{}
	this.ContainerId = containerId
	return &this
}

// NewPrivilegedResourceAccountAppWithDefaults instantiates a new PrivilegedResourceAccountApp object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivilegedResourceAccountAppWithDefaults() *PrivilegedResourceAccountApp {
	this := PrivilegedResourceAccountApp{}
	return &this
}

// GetContainerId returns the ContainerId field value
func (o *PrivilegedResourceAccountApp) GetContainerId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ContainerId
}

// GetContainerIdOk returns a tuple with the ContainerId field value
// and a boolean to check if the value has been set.
func (o *PrivilegedResourceAccountApp) GetContainerIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ContainerId, true
}

// SetContainerId sets field value
func (o *PrivilegedResourceAccountApp) SetContainerId(v string) {
	o.ContainerId = v
}

// GetCredentials returns the Credentials field value if set, zero value otherwise.
func (o *PrivilegedResourceAccountApp) GetCredentials() PrivilegedResourceCredentials {
	if o == nil || o.Credentials == nil {
		var ret PrivilegedResourceCredentials
		return ret
	}
	return *o.Credentials
}

// GetCredentialsOk returns a tuple with the Credentials field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivilegedResourceAccountApp) GetCredentialsOk() (*PrivilegedResourceCredentials, bool) {
	if o == nil || o.Credentials == nil {
		return nil, false
	}
	return o.Credentials, true
}

// HasCredentials returns a boolean if a field has been set.
func (o *PrivilegedResourceAccountApp) HasCredentials() bool {
	if o != nil && o.Credentials != nil {
		return true
	}

	return false
}

// SetCredentials gets a reference to the given PrivilegedResourceCredentials and assigns it to the Credentials field.
func (o *PrivilegedResourceAccountApp) SetCredentials(v PrivilegedResourceCredentials) {
	o.Credentials = &v
}

// GetContainerDisplayName returns the ContainerDisplayName field value if set, zero value otherwise.
func (o *PrivilegedResourceAccountApp) GetContainerDisplayName() string {
	if o == nil || o.ContainerDisplayName == nil {
		var ret string
		return ret
	}
	return *o.ContainerDisplayName
}

// GetContainerDisplayNameOk returns a tuple with the ContainerDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivilegedResourceAccountApp) GetContainerDisplayNameOk() (*string, bool) {
	if o == nil || o.ContainerDisplayName == nil {
		return nil, false
	}
	return o.ContainerDisplayName, true
}

// HasContainerDisplayName returns a boolean if a field has been set.
func (o *PrivilegedResourceAccountApp) HasContainerDisplayName() bool {
	if o != nil && o.ContainerDisplayName != nil {
		return true
	}

	return false
}

// SetContainerDisplayName gets a reference to the given string and assigns it to the ContainerDisplayName field.
func (o *PrivilegedResourceAccountApp) SetContainerDisplayName(v string) {
	o.ContainerDisplayName = &v
}

// GetLinks returns the Links field value if set, zero value otherwise.
func (o *PrivilegedResourceAccountApp) GetLinks() AppLink {
	if o == nil || o.Links == nil {
		var ret AppLink
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivilegedResourceAccountApp) GetLinksOk() (*AppLink, bool) {
	if o == nil || o.Links == nil {
		return nil, false
	}
	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PrivilegedResourceAccountApp) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

// SetLinks gets a reference to the given AppLink and assigns it to the Links field.
func (o *PrivilegedResourceAccountApp) SetLinks(v AppLink) {
	o.Links = &v
}

func (o PrivilegedResourceAccountApp) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	serializedPrivilegedResource, errPrivilegedResource := json.Marshal(o.PrivilegedResource)
	if errPrivilegedResource != nil {
		return []byte{}, errPrivilegedResource
	}
	errPrivilegedResource = json.Unmarshal([]byte(serializedPrivilegedResource), &toSerialize)
	if errPrivilegedResource != nil {
		return []byte{}, errPrivilegedResource
	}
	if true {
		toSerialize["containerId"] = o.ContainerId
	}
	if o.Credentials != nil {
		toSerialize["credentials"] = o.Credentials
	}
	if o.ContainerDisplayName != nil {
		toSerialize["containerDisplayName"] = o.ContainerDisplayName
	}
	if o.Links != nil {
		toSerialize["_links"] = o.Links
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *PrivilegedResourceAccountApp) UnmarshalJSON(bytes []byte) (err error) {
	type PrivilegedResourceAccountAppWithoutEmbeddedStruct struct {
		// The application ID associated with the privileged account
		ContainerId string `json:"containerId"`
		Credentials *PrivilegedResourceCredentials `json:"credentials,omitempty"`
		// Human-readable name of the container that owns the privileged resource
		ContainerDisplayName *string `json:"containerDisplayName,omitempty"`
		Links *AppLink `json:"_links,omitempty"`
	}

	varPrivilegedResourceAccountAppWithoutEmbeddedStruct := PrivilegedResourceAccountAppWithoutEmbeddedStruct{}

	err = json.Unmarshal(bytes, &varPrivilegedResourceAccountAppWithoutEmbeddedStruct)
	if err == nil {
		varPrivilegedResourceAccountApp := _PrivilegedResourceAccountApp{}
		varPrivilegedResourceAccountApp.ContainerId = varPrivilegedResourceAccountAppWithoutEmbeddedStruct.ContainerId
		varPrivilegedResourceAccountApp.Credentials = varPrivilegedResourceAccountAppWithoutEmbeddedStruct.Credentials
		varPrivilegedResourceAccountApp.ContainerDisplayName = varPrivilegedResourceAccountAppWithoutEmbeddedStruct.ContainerDisplayName
		varPrivilegedResourceAccountApp.Links = varPrivilegedResourceAccountAppWithoutEmbeddedStruct.Links
		*o = PrivilegedResourceAccountApp(varPrivilegedResourceAccountApp)
	} else {
		return err
	}

	varPrivilegedResourceAccountApp := _PrivilegedResourceAccountApp{}

	err = json.Unmarshal(bytes, &varPrivilegedResourceAccountApp)
	if err == nil {
		o.PrivilegedResource = varPrivilegedResourceAccountApp.PrivilegedResource
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "containerId")
		delete(additionalProperties, "credentials")
		delete(additionalProperties, "containerDisplayName")
		delete(additionalProperties, "_links")

		// remove fields from embedded structs
		reflectPrivilegedResource := reflect.ValueOf(o.PrivilegedResource)
		for i := 0; i < reflectPrivilegedResource.Type().NumField(); i++ {
			t := reflectPrivilegedResource.Type().Field(i)

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

type NullablePrivilegedResourceAccountApp struct {
	value *PrivilegedResourceAccountApp
	isSet bool
}

func (v NullablePrivilegedResourceAccountApp) Get() *PrivilegedResourceAccountApp {
	return v.value
}

func (v *NullablePrivilegedResourceAccountApp) Set(val *PrivilegedResourceAccountApp) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivilegedResourceAccountApp) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivilegedResourceAccountApp) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivilegedResourceAccountApp(val *PrivilegedResourceAccountApp) *NullablePrivilegedResourceAccountApp {
	return &NullablePrivilegedResourceAccountApp{value: val, isSet: true}
}

func (v NullablePrivilegedResourceAccountApp) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivilegedResourceAccountApp) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
