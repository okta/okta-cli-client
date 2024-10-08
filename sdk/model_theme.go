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

// Theme struct for Theme
type Theme struct {
	BackgroundImage                   *string    `json:"backgroundImage,omitempty"`
	EmailTemplateTouchPointVariant    *string    `json:"emailTemplateTouchPointVariant,omitempty"`
	EndUserDashboardTouchPointVariant *string    `json:"endUserDashboardTouchPointVariant,omitempty"`
	ErrorPageTouchPointVariant        *string    `json:"errorPageTouchPointVariant,omitempty"`
	LoadingPageTouchPointVariant      *string    `json:"loadingPageTouchPointVariant,omitempty"`
	PrimaryColorContrastHex           *string    `json:"primaryColorContrastHex,omitempty"`
	PrimaryColorHex                   *string    `json:"primaryColorHex,omitempty"`
	SecondaryColorContrastHex         *string    `json:"secondaryColorContrastHex,omitempty"`
	SecondaryColorHex                 *string    `json:"secondaryColorHex,omitempty"`
	SignInPageTouchPointVariant       *string    `json:"signInPageTouchPointVariant,omitempty"`
	Links                             *LinksSelf `json:"_links,omitempty"`
	AdditionalProperties              map[string]interface{}
}

type _Theme Theme

// NewTheme instantiates a new Theme object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTheme() *Theme {
	this := Theme{}
	return &this
}

// NewThemeWithDefaults instantiates a new Theme object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewThemeWithDefaults() *Theme {
	this := Theme{}
	return &this
}

// GetBackgroundImage returns the BackgroundImage field value if set, zero value otherwise.
func (o *Theme) GetBackgroundImage() string {
	if o == nil || o.BackgroundImage == nil {
		var ret string
		return ret
	}
	return *o.BackgroundImage
}

// GetBackgroundImageOk returns a tuple with the BackgroundImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetBackgroundImageOk() (*string, bool) {
	if o == nil || o.BackgroundImage == nil {
		return nil, false
	}
	return o.BackgroundImage, true
}

// HasBackgroundImage returns a boolean if a field has been set.
func (o *Theme) HasBackgroundImage() bool {
	if o != nil && o.BackgroundImage != nil {
		return true
	}

	return false
}

// SetBackgroundImage gets a reference to the given string and assigns it to the BackgroundImage field.
func (o *Theme) SetBackgroundImage(v string) {
	o.BackgroundImage = &v
}

// GetEmailTemplateTouchPointVariant returns the EmailTemplateTouchPointVariant field value if set, zero value otherwise.
func (o *Theme) GetEmailTemplateTouchPointVariant() string {
	if o == nil || o.EmailTemplateTouchPointVariant == nil {
		var ret string
		return ret
	}
	return *o.EmailTemplateTouchPointVariant
}

// GetEmailTemplateTouchPointVariantOk returns a tuple with the EmailTemplateTouchPointVariant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetEmailTemplateTouchPointVariantOk() (*string, bool) {
	if o == nil || o.EmailTemplateTouchPointVariant == nil {
		return nil, false
	}
	return o.EmailTemplateTouchPointVariant, true
}

// HasEmailTemplateTouchPointVariant returns a boolean if a field has been set.
func (o *Theme) HasEmailTemplateTouchPointVariant() bool {
	if o != nil && o.EmailTemplateTouchPointVariant != nil {
		return true
	}

	return false
}

// SetEmailTemplateTouchPointVariant gets a reference to the given string and assigns it to the EmailTemplateTouchPointVariant field.
func (o *Theme) SetEmailTemplateTouchPointVariant(v string) {
	o.EmailTemplateTouchPointVariant = &v
}

// GetEndUserDashboardTouchPointVariant returns the EndUserDashboardTouchPointVariant field value if set, zero value otherwise.
func (o *Theme) GetEndUserDashboardTouchPointVariant() string {
	if o == nil || o.EndUserDashboardTouchPointVariant == nil {
		var ret string
		return ret
	}
	return *o.EndUserDashboardTouchPointVariant
}

// GetEndUserDashboardTouchPointVariantOk returns a tuple with the EndUserDashboardTouchPointVariant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetEndUserDashboardTouchPointVariantOk() (*string, bool) {
	if o == nil || o.EndUserDashboardTouchPointVariant == nil {
		return nil, false
	}
	return o.EndUserDashboardTouchPointVariant, true
}

// HasEndUserDashboardTouchPointVariant returns a boolean if a field has been set.
func (o *Theme) HasEndUserDashboardTouchPointVariant() bool {
	if o != nil && o.EndUserDashboardTouchPointVariant != nil {
		return true
	}

	return false
}

// SetEndUserDashboardTouchPointVariant gets a reference to the given string and assigns it to the EndUserDashboardTouchPointVariant field.
func (o *Theme) SetEndUserDashboardTouchPointVariant(v string) {
	o.EndUserDashboardTouchPointVariant = &v
}

// GetErrorPageTouchPointVariant returns the ErrorPageTouchPointVariant field value if set, zero value otherwise.
func (o *Theme) GetErrorPageTouchPointVariant() string {
	if o == nil || o.ErrorPageTouchPointVariant == nil {
		var ret string
		return ret
	}
	return *o.ErrorPageTouchPointVariant
}

// GetErrorPageTouchPointVariantOk returns a tuple with the ErrorPageTouchPointVariant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetErrorPageTouchPointVariantOk() (*string, bool) {
	if o == nil || o.ErrorPageTouchPointVariant == nil {
		return nil, false
	}
	return o.ErrorPageTouchPointVariant, true
}

// HasErrorPageTouchPointVariant returns a boolean if a field has been set.
func (o *Theme) HasErrorPageTouchPointVariant() bool {
	if o != nil && o.ErrorPageTouchPointVariant != nil {
		return true
	}

	return false
}

// SetErrorPageTouchPointVariant gets a reference to the given string and assigns it to the ErrorPageTouchPointVariant field.
func (o *Theme) SetErrorPageTouchPointVariant(v string) {
	o.ErrorPageTouchPointVariant = &v
}

// GetLoadingPageTouchPointVariant returns the LoadingPageTouchPointVariant field value if set, zero value otherwise.
func (o *Theme) GetLoadingPageTouchPointVariant() string {
	if o == nil || o.LoadingPageTouchPointVariant == nil {
		var ret string
		return ret
	}
	return *o.LoadingPageTouchPointVariant
}

// GetLoadingPageTouchPointVariantOk returns a tuple with the LoadingPageTouchPointVariant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetLoadingPageTouchPointVariantOk() (*string, bool) {
	if o == nil || o.LoadingPageTouchPointVariant == nil {
		return nil, false
	}
	return o.LoadingPageTouchPointVariant, true
}

// HasLoadingPageTouchPointVariant returns a boolean if a field has been set.
func (o *Theme) HasLoadingPageTouchPointVariant() bool {
	if o != nil && o.LoadingPageTouchPointVariant != nil {
		return true
	}

	return false
}

// SetLoadingPageTouchPointVariant gets a reference to the given string and assigns it to the LoadingPageTouchPointVariant field.
func (o *Theme) SetLoadingPageTouchPointVariant(v string) {
	o.LoadingPageTouchPointVariant = &v
}

// GetPrimaryColorContrastHex returns the PrimaryColorContrastHex field value if set, zero value otherwise.
func (o *Theme) GetPrimaryColorContrastHex() string {
	if o == nil || o.PrimaryColorContrastHex == nil {
		var ret string
		return ret
	}
	return *o.PrimaryColorContrastHex
}

// GetPrimaryColorContrastHexOk returns a tuple with the PrimaryColorContrastHex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetPrimaryColorContrastHexOk() (*string, bool) {
	if o == nil || o.PrimaryColorContrastHex == nil {
		return nil, false
	}
	return o.PrimaryColorContrastHex, true
}

// HasPrimaryColorContrastHex returns a boolean if a field has been set.
func (o *Theme) HasPrimaryColorContrastHex() bool {
	if o != nil && o.PrimaryColorContrastHex != nil {
		return true
	}

	return false
}

// SetPrimaryColorContrastHex gets a reference to the given string and assigns it to the PrimaryColorContrastHex field.
func (o *Theme) SetPrimaryColorContrastHex(v string) {
	o.PrimaryColorContrastHex = &v
}

// GetPrimaryColorHex returns the PrimaryColorHex field value if set, zero value otherwise.
func (o *Theme) GetPrimaryColorHex() string {
	if o == nil || o.PrimaryColorHex == nil {
		var ret string
		return ret
	}
	return *o.PrimaryColorHex
}

// GetPrimaryColorHexOk returns a tuple with the PrimaryColorHex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetPrimaryColorHexOk() (*string, bool) {
	if o == nil || o.PrimaryColorHex == nil {
		return nil, false
	}
	return o.PrimaryColorHex, true
}

// HasPrimaryColorHex returns a boolean if a field has been set.
func (o *Theme) HasPrimaryColorHex() bool {
	if o != nil && o.PrimaryColorHex != nil {
		return true
	}

	return false
}

// SetPrimaryColorHex gets a reference to the given string and assigns it to the PrimaryColorHex field.
func (o *Theme) SetPrimaryColorHex(v string) {
	o.PrimaryColorHex = &v
}

// GetSecondaryColorContrastHex returns the SecondaryColorContrastHex field value if set, zero value otherwise.
func (o *Theme) GetSecondaryColorContrastHex() string {
	if o == nil || o.SecondaryColorContrastHex == nil {
		var ret string
		return ret
	}
	return *o.SecondaryColorContrastHex
}

// GetSecondaryColorContrastHexOk returns a tuple with the SecondaryColorContrastHex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetSecondaryColorContrastHexOk() (*string, bool) {
	if o == nil || o.SecondaryColorContrastHex == nil {
		return nil, false
	}
	return o.SecondaryColorContrastHex, true
}

// HasSecondaryColorContrastHex returns a boolean if a field has been set.
func (o *Theme) HasSecondaryColorContrastHex() bool {
	if o != nil && o.SecondaryColorContrastHex != nil {
		return true
	}

	return false
}

// SetSecondaryColorContrastHex gets a reference to the given string and assigns it to the SecondaryColorContrastHex field.
func (o *Theme) SetSecondaryColorContrastHex(v string) {
	o.SecondaryColorContrastHex = &v
}

// GetSecondaryColorHex returns the SecondaryColorHex field value if set, zero value otherwise.
func (o *Theme) GetSecondaryColorHex() string {
	if o == nil || o.SecondaryColorHex == nil {
		var ret string
		return ret
	}
	return *o.SecondaryColorHex
}

// GetSecondaryColorHexOk returns a tuple with the SecondaryColorHex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetSecondaryColorHexOk() (*string, bool) {
	if o == nil || o.SecondaryColorHex == nil {
		return nil, false
	}
	return o.SecondaryColorHex, true
}

// HasSecondaryColorHex returns a boolean if a field has been set.
func (o *Theme) HasSecondaryColorHex() bool {
	if o != nil && o.SecondaryColorHex != nil {
		return true
	}

	return false
}

// SetSecondaryColorHex gets a reference to the given string and assigns it to the SecondaryColorHex field.
func (o *Theme) SetSecondaryColorHex(v string) {
	o.SecondaryColorHex = &v
}

// GetSignInPageTouchPointVariant returns the SignInPageTouchPointVariant field value if set, zero value otherwise.
func (o *Theme) GetSignInPageTouchPointVariant() string {
	if o == nil || o.SignInPageTouchPointVariant == nil {
		var ret string
		return ret
	}
	return *o.SignInPageTouchPointVariant
}

// GetSignInPageTouchPointVariantOk returns a tuple with the SignInPageTouchPointVariant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetSignInPageTouchPointVariantOk() (*string, bool) {
	if o == nil || o.SignInPageTouchPointVariant == nil {
		return nil, false
	}
	return o.SignInPageTouchPointVariant, true
}

// HasSignInPageTouchPointVariant returns a boolean if a field has been set.
func (o *Theme) HasSignInPageTouchPointVariant() bool {
	if o != nil && o.SignInPageTouchPointVariant != nil {
		return true
	}

	return false
}

// SetSignInPageTouchPointVariant gets a reference to the given string and assigns it to the SignInPageTouchPointVariant field.
func (o *Theme) SetSignInPageTouchPointVariant(v string) {
	o.SignInPageTouchPointVariant = &v
}

// GetLinks returns the Links field value if set, zero value otherwise.
func (o *Theme) GetLinks() LinksSelf {
	if o == nil || o.Links == nil {
		var ret LinksSelf
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Theme) GetLinksOk() (*LinksSelf, bool) {
	if o == nil || o.Links == nil {
		return nil, false
	}
	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *Theme) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

// SetLinks gets a reference to the given LinksSelf and assigns it to the Links field.
func (o *Theme) SetLinks(v LinksSelf) {
	o.Links = &v
}

func (o Theme) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.BackgroundImage != nil {
		toSerialize["backgroundImage"] = o.BackgroundImage
	}
	if o.EmailTemplateTouchPointVariant != nil {
		toSerialize["emailTemplateTouchPointVariant"] = o.EmailTemplateTouchPointVariant
	}
	if o.EndUserDashboardTouchPointVariant != nil {
		toSerialize["endUserDashboardTouchPointVariant"] = o.EndUserDashboardTouchPointVariant
	}
	if o.ErrorPageTouchPointVariant != nil {
		toSerialize["errorPageTouchPointVariant"] = o.ErrorPageTouchPointVariant
	}
	if o.LoadingPageTouchPointVariant != nil {
		toSerialize["loadingPageTouchPointVariant"] = o.LoadingPageTouchPointVariant
	}
	if o.PrimaryColorContrastHex != nil {
		toSerialize["primaryColorContrastHex"] = o.PrimaryColorContrastHex
	}
	if o.PrimaryColorHex != nil {
		toSerialize["primaryColorHex"] = o.PrimaryColorHex
	}
	if o.SecondaryColorContrastHex != nil {
		toSerialize["secondaryColorContrastHex"] = o.SecondaryColorContrastHex
	}
	if o.SecondaryColorHex != nil {
		toSerialize["secondaryColorHex"] = o.SecondaryColorHex
	}
	if o.SignInPageTouchPointVariant != nil {
		toSerialize["signInPageTouchPointVariant"] = o.SignInPageTouchPointVariant
	}
	if o.Links != nil {
		toSerialize["_links"] = o.Links
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *Theme) UnmarshalJSON(bytes []byte) (err error) {
	varTheme := _Theme{}

	err = json.Unmarshal(bytes, &varTheme)
	if err == nil {
		*o = Theme(varTheme)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "backgroundImage")
		delete(additionalProperties, "emailTemplateTouchPointVariant")
		delete(additionalProperties, "endUserDashboardTouchPointVariant")
		delete(additionalProperties, "errorPageTouchPointVariant")
		delete(additionalProperties, "loadingPageTouchPointVariant")
		delete(additionalProperties, "primaryColorContrastHex")
		delete(additionalProperties, "primaryColorHex")
		delete(additionalProperties, "secondaryColorContrastHex")
		delete(additionalProperties, "secondaryColorHex")
		delete(additionalProperties, "signInPageTouchPointVariant")
		delete(additionalProperties, "_links")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableTheme struct {
	value *Theme
	isSet bool
}

func (v NullableTheme) Get() *Theme {
	return v.value
}

func (v *NullableTheme) Set(val *Theme) {
	v.value = val
	v.isSet = true
}

func (v NullableTheme) IsSet() bool {
	return v.isSet
}

func (v *NullableTheme) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTheme(val *Theme) *NullableTheme {
	return &NullableTheme{value: val, isSet: true}
}

func (v NullableTheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTheme) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
