# LinksAppAndUserGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hints** | Pointer to [**HrefHints**](HrefHints.md) |  | [optional] 
**Href** | **string** | Link URI | 
**Name** | Pointer to **string** | Link name | [optional] 
**Type** | Pointer to **string** | The media type of the link. If omitted, it is implicitly &#x60;application/json&#x60;. | [optional] 

## Methods

### NewLinksAppAndUserGroup

`func NewLinksAppAndUserGroup(href string, ) *LinksAppAndUserGroup`

NewLinksAppAndUserGroup instantiates a new LinksAppAndUserGroup object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLinksAppAndUserGroupWithDefaults

`func NewLinksAppAndUserGroupWithDefaults() *LinksAppAndUserGroup`

NewLinksAppAndUserGroupWithDefaults instantiates a new LinksAppAndUserGroup object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHints

`func (o *LinksAppAndUserGroup) GetHints() HrefHints`

GetHints returns the Hints field if non-nil, zero value otherwise.

### GetHintsOk

`func (o *LinksAppAndUserGroup) GetHintsOk() (*HrefHints, bool)`

GetHintsOk returns a tuple with the Hints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHints

`func (o *LinksAppAndUserGroup) SetHints(v HrefHints)`

SetHints sets Hints field to given value.

### HasHints

`func (o *LinksAppAndUserGroup) HasHints() bool`

HasHints returns a boolean if a field has been set.

### GetHref

`func (o *LinksAppAndUserGroup) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *LinksAppAndUserGroup) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *LinksAppAndUserGroup) SetHref(v string)`

SetHref sets Href field to given value.


### GetName

`func (o *LinksAppAndUserGroup) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *LinksAppAndUserGroup) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *LinksAppAndUserGroup) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *LinksAppAndUserGroup) HasName() bool`

HasName returns a boolean if a field has been set.

### GetType

`func (o *LinksAppAndUserGroup) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *LinksAppAndUserGroup) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *LinksAppAndUserGroup) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *LinksAppAndUserGroup) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


