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

// SourceLinksAllOfSchema struct for SourceLinksAllOfSchema
type SourceLinksAllOfSchema struct {
	Hints *HrefHints `json:"hints,omitempty"`
	// Link URI
	Href string `json:"href"`
	// Link name
	Name *string `json:"name,omitempty"`
	// The media type of the link. If omitted, it is implicitly `application/json`.
	Type                 *string `json:"type,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _SourceLinksAllOfSchema SourceLinksAllOfSchema

// NewSourceLinksAllOfSchema instantiates a new SourceLinksAllOfSchema object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSourceLinksAllOfSchema(href string) *SourceLinksAllOfSchema {
	this := SourceLinksAllOfSchema{}
	this.Href = href
	return &this
}

// NewSourceLinksAllOfSchemaWithDefaults instantiates a new SourceLinksAllOfSchema object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSourceLinksAllOfSchemaWithDefaults() *SourceLinksAllOfSchema {
	this := SourceLinksAllOfSchema{}
	return &this
}

// GetHints returns the Hints field value if set, zero value otherwise.
func (o *SourceLinksAllOfSchema) GetHints() HrefHints {
	if o == nil || o.Hints == nil {
		var ret HrefHints
		return ret
	}
	return *o.Hints
}

// GetHintsOk returns a tuple with the Hints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SourceLinksAllOfSchema) GetHintsOk() (*HrefHints, bool) {
	if o == nil || o.Hints == nil {
		return nil, false
	}
	return o.Hints, true
}

// HasHints returns a boolean if a field has been set.
func (o *SourceLinksAllOfSchema) HasHints() bool {
	if o != nil && o.Hints != nil {
		return true
	}

	return false
}

// SetHints gets a reference to the given HrefHints and assigns it to the Hints field.
func (o *SourceLinksAllOfSchema) SetHints(v HrefHints) {
	o.Hints = &v
}

// GetHref returns the Href field value
func (o *SourceLinksAllOfSchema) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *SourceLinksAllOfSchema) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *SourceLinksAllOfSchema) SetHref(v string) {
	o.Href = v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *SourceLinksAllOfSchema) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SourceLinksAllOfSchema) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *SourceLinksAllOfSchema) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *SourceLinksAllOfSchema) SetName(v string) {
	o.Name = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *SourceLinksAllOfSchema) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SourceLinksAllOfSchema) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *SourceLinksAllOfSchema) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *SourceLinksAllOfSchema) SetType(v string) {
	o.Type = &v
}

func (o SourceLinksAllOfSchema) MarshalJSON() ([]byte, error) {
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

func (o *SourceLinksAllOfSchema) UnmarshalJSON(bytes []byte) (err error) {
	varSourceLinksAllOfSchema := _SourceLinksAllOfSchema{}

	err = json.Unmarshal(bytes, &varSourceLinksAllOfSchema)
	if err == nil {
		*o = SourceLinksAllOfSchema(varSourceLinksAllOfSchema)
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

type NullableSourceLinksAllOfSchema struct {
	value *SourceLinksAllOfSchema
	isSet bool
}

func (v NullableSourceLinksAllOfSchema) Get() *SourceLinksAllOfSchema {
	return v.value
}

func (v *NullableSourceLinksAllOfSchema) Set(val *SourceLinksAllOfSchema) {
	v.value = val
	v.isSet = true
}

func (v NullableSourceLinksAllOfSchema) IsSet() bool {
	return v.isSet
}

func (v *NullableSourceLinksAllOfSchema) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSourceLinksAllOfSchema(val *SourceLinksAllOfSchema) *NullableSourceLinksAllOfSchema {
	return &NullableSourceLinksAllOfSchema{value: val, isSet: true}
}

func (v NullableSourceLinksAllOfSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSourceLinksAllOfSchema) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
