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

// BrowserPluginApplication struct for BrowserPluginApplication
type BrowserPluginApplication struct {
	Application
	Credentials *SchemeApplicationCredentials `json:"credentials,omitempty"`
	// Unique key for the app definition
	Name                 *string                 `json:"name,omitempty"`
	Settings             *SwaApplicationSettings `json:"settings,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _BrowserPluginApplication BrowserPluginApplication

// NewBrowserPluginApplication instantiates a new BrowserPluginApplication object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBrowserPluginApplication() *BrowserPluginApplication {
	this := BrowserPluginApplication{}
	return &this
}

// NewBrowserPluginApplicationWithDefaults instantiates a new BrowserPluginApplication object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBrowserPluginApplicationWithDefaults() *BrowserPluginApplication {
	this := BrowserPluginApplication{}
	return &this
}

// GetCredentials returns the Credentials field value if set, zero value otherwise.
func (o *BrowserPluginApplication) GetCredentials() SchemeApplicationCredentials {
	if o == nil || o.Credentials == nil {
		var ret SchemeApplicationCredentials
		return ret
	}
	return *o.Credentials
}

// GetCredentialsOk returns a tuple with the Credentials field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BrowserPluginApplication) GetCredentialsOk() (*SchemeApplicationCredentials, bool) {
	if o == nil || o.Credentials == nil {
		return nil, false
	}
	return o.Credentials, true
}

// HasCredentials returns a boolean if a field has been set.
func (o *BrowserPluginApplication) HasCredentials() bool {
	if o != nil && o.Credentials != nil {
		return true
	}

	return false
}

// SetCredentials gets a reference to the given SchemeApplicationCredentials and assigns it to the Credentials field.
func (o *BrowserPluginApplication) SetCredentials(v SchemeApplicationCredentials) {
	o.Credentials = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *BrowserPluginApplication) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BrowserPluginApplication) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *BrowserPluginApplication) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *BrowserPluginApplication) SetName(v string) {
	o.Name = &v
}

// GetSettings returns the Settings field value if set, zero value otherwise.
func (o *BrowserPluginApplication) GetSettings() SwaApplicationSettings {
	if o == nil || o.Settings == nil {
		var ret SwaApplicationSettings
		return ret
	}
	return *o.Settings
}

// GetSettingsOk returns a tuple with the Settings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BrowserPluginApplication) GetSettingsOk() (*SwaApplicationSettings, bool) {
	if o == nil || o.Settings == nil {
		return nil, false
	}
	return o.Settings, true
}

// HasSettings returns a boolean if a field has been set.
func (o *BrowserPluginApplication) HasSettings() bool {
	if o != nil && o.Settings != nil {
		return true
	}

	return false
}

// SetSettings gets a reference to the given SwaApplicationSettings and assigns it to the Settings field.
func (o *BrowserPluginApplication) SetSettings(v SwaApplicationSettings) {
	o.Settings = &v
}

func (o BrowserPluginApplication) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	serializedApplication, errApplication := json.Marshal(o.Application)
	if errApplication != nil {
		return []byte{}, errApplication
	}
	errApplication = json.Unmarshal([]byte(serializedApplication), &toSerialize)
	if errApplication != nil {
		return []byte{}, errApplication
	}
	if o.Credentials != nil {
		toSerialize["credentials"] = o.Credentials
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Settings != nil {
		toSerialize["settings"] = o.Settings
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *BrowserPluginApplication) UnmarshalJSON(bytes []byte) (err error) {
	type BrowserPluginApplicationWithoutEmbeddedStruct struct {
		Credentials *SchemeApplicationCredentials `json:"credentials,omitempty"`
		// Unique key for the app definition
		Name     *string                 `json:"name,omitempty"`
		Settings *SwaApplicationSettings `json:"settings,omitempty"`
	}

	varBrowserPluginApplicationWithoutEmbeddedStruct := BrowserPluginApplicationWithoutEmbeddedStruct{}

	err = json.Unmarshal(bytes, &varBrowserPluginApplicationWithoutEmbeddedStruct)
	if err == nil {
		varBrowserPluginApplication := _BrowserPluginApplication{}
		varBrowserPluginApplication.Credentials = varBrowserPluginApplicationWithoutEmbeddedStruct.Credentials
		varBrowserPluginApplication.Name = varBrowserPluginApplicationWithoutEmbeddedStruct.Name
		varBrowserPluginApplication.Settings = varBrowserPluginApplicationWithoutEmbeddedStruct.Settings
		*o = BrowserPluginApplication(varBrowserPluginApplication)
	} else {
		return err
	}

	varBrowserPluginApplication := _BrowserPluginApplication{}

	err = json.Unmarshal(bytes, &varBrowserPluginApplication)
	if err == nil {
		o.Application = varBrowserPluginApplication.Application
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "credentials")
		delete(additionalProperties, "name")
		delete(additionalProperties, "settings")

		// remove fields from embedded structs
		reflectApplication := reflect.ValueOf(o.Application)
		for i := 0; i < reflectApplication.Type().NumField(); i++ {
			t := reflectApplication.Type().Field(i)

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

type NullableBrowserPluginApplication struct {
	value *BrowserPluginApplication
	isSet bool
}

func (v NullableBrowserPluginApplication) Get() *BrowserPluginApplication {
	return v.value
}

func (v *NullableBrowserPluginApplication) Set(val *BrowserPluginApplication) {
	v.value = val
	v.isSet = true
}

func (v NullableBrowserPluginApplication) IsSet() bool {
	return v.isSet
}

func (v *NullableBrowserPluginApplication) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBrowserPluginApplication(val *BrowserPluginApplication) *NullableBrowserPluginApplication {
	return &NullableBrowserPluginApplication{value: val, isSet: true}
}

func (v NullableBrowserPluginApplication) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBrowserPluginApplication) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
