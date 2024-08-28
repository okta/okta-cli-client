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

// AuthServerLinksAllOfScopes struct for AuthServerLinksAllOfScopes
type AuthServerLinksAllOfScopes struct {
	Hints *HrefHints `json:"hints,omitempty"`
	// Link URI
	Href string `json:"href"`
	// Link name
	Name *string `json:"name,omitempty"`
	// The media type of the link. If omitted, it is implicitly `application/json`.
	Type                 *string `json:"type,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _AuthServerLinksAllOfScopes AuthServerLinksAllOfScopes

// NewAuthServerLinksAllOfScopes instantiates a new AuthServerLinksAllOfScopes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthServerLinksAllOfScopes(href string) *AuthServerLinksAllOfScopes {
	this := AuthServerLinksAllOfScopes{}
	this.Href = href
	return &this
}

// NewAuthServerLinksAllOfScopesWithDefaults instantiates a new AuthServerLinksAllOfScopes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthServerLinksAllOfScopesWithDefaults() *AuthServerLinksAllOfScopes {
	this := AuthServerLinksAllOfScopes{}
	return &this
}

// GetHints returns the Hints field value if set, zero value otherwise.
func (o *AuthServerLinksAllOfScopes) GetHints() HrefHints {
	if o == nil || o.Hints == nil {
		var ret HrefHints
		return ret
	}
	return *o.Hints
}

// GetHintsOk returns a tuple with the Hints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthServerLinksAllOfScopes) GetHintsOk() (*HrefHints, bool) {
	if o == nil || o.Hints == nil {
		return nil, false
	}
	return o.Hints, true
}

// HasHints returns a boolean if a field has been set.
func (o *AuthServerLinksAllOfScopes) HasHints() bool {
	if o != nil && o.Hints != nil {
		return true
	}

	return false
}

// SetHints gets a reference to the given HrefHints and assigns it to the Hints field.
func (o *AuthServerLinksAllOfScopes) SetHints(v HrefHints) {
	o.Hints = &v
}

// GetHref returns the Href field value
func (o *AuthServerLinksAllOfScopes) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *AuthServerLinksAllOfScopes) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *AuthServerLinksAllOfScopes) SetHref(v string) {
	o.Href = v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *AuthServerLinksAllOfScopes) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthServerLinksAllOfScopes) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *AuthServerLinksAllOfScopes) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *AuthServerLinksAllOfScopes) SetName(v string) {
	o.Name = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *AuthServerLinksAllOfScopes) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthServerLinksAllOfScopes) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *AuthServerLinksAllOfScopes) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *AuthServerLinksAllOfScopes) SetType(v string) {
	o.Type = &v
}

func (o AuthServerLinksAllOfScopes) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Hints != nil {
		toSerialize["hints"] = o.Hints
	}
	if true {
		toSerialize["href"] = o.Href
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *AuthServerLinksAllOfScopes) UnmarshalJSON(bytes []byte) (err error) {
	varAuthServerLinksAllOfScopes := _AuthServerLinksAllOfScopes{}

	err = json.Unmarshal(bytes, &varAuthServerLinksAllOfScopes)
	if err == nil {
		*o = AuthServerLinksAllOfScopes(varAuthServerLinksAllOfScopes)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "hints")
		delete(additionalProperties, "href")
		delete(additionalProperties, "name")
		delete(additionalProperties, "type")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableAuthServerLinksAllOfScopes struct {
	value *AuthServerLinksAllOfScopes
	isSet bool
}

func (v NullableAuthServerLinksAllOfScopes) Get() *AuthServerLinksAllOfScopes {
	return v.value
}

func (v *NullableAuthServerLinksAllOfScopes) Set(val *AuthServerLinksAllOfScopes) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthServerLinksAllOfScopes) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthServerLinksAllOfScopes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthServerLinksAllOfScopes(val *AuthServerLinksAllOfScopes) *NullableAuthServerLinksAllOfScopes {
	return &NullableAuthServerLinksAllOfScopes{value: val, isSet: true}
}

func (v NullableAuthServerLinksAllOfScopes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthServerLinksAllOfScopes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
