# PrivilegedResourceAccountOkta

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ResourceId** | **string** | The user ID associated with the Okta privileged resource | 
**Credentials** | Pointer to [**PrivilegedResourceAccountOktaAllOfCredentials**](PrivilegedResourceAccountOktaAllOfCredentials.md) |  | [optional] 
**Links** | Pointer to [**UserLink**](UserLink.md) |  | [optional] 

## Methods

### NewPrivilegedResourceAccountOkta

`func NewPrivilegedResourceAccountOkta(resourceId string, ) *PrivilegedResourceAccountOkta`

NewPrivilegedResourceAccountOkta instantiates a new PrivilegedResourceAccountOkta object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrivilegedResourceAccountOktaWithDefaults

`func NewPrivilegedResourceAccountOktaWithDefaults() *PrivilegedResourceAccountOkta`

NewPrivilegedResourceAccountOktaWithDefaults instantiates a new PrivilegedResourceAccountOkta object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetResourceId

`func (o *PrivilegedResourceAccountOkta) GetResourceId() string`

GetResourceId returns the ResourceId field if non-nil, zero value otherwise.

### GetResourceIdOk

`func (o *PrivilegedResourceAccountOkta) GetResourceIdOk() (*string, bool)`

GetResourceIdOk returns a tuple with the ResourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceId

`func (o *PrivilegedResourceAccountOkta) SetResourceId(v string)`

SetResourceId sets ResourceId field to given value.


### GetCredentials

`func (o *PrivilegedResourceAccountOkta) GetCredentials() PrivilegedResourceAccountOktaAllOfCredentials`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *PrivilegedResourceAccountOkta) GetCredentialsOk() (*PrivilegedResourceAccountOktaAllOfCredentials, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *PrivilegedResourceAccountOkta) SetCredentials(v PrivilegedResourceAccountOktaAllOfCredentials)`

SetCredentials sets Credentials field to given value.

### HasCredentials

`func (o *PrivilegedResourceAccountOkta) HasCredentials() bool`

HasCredentials returns a boolean if a field has been set.

### GetLinks

`func (o *PrivilegedResourceAccountOkta) GetLinks() UserLink`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *PrivilegedResourceAccountOkta) GetLinksOk() (*UserLink, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *PrivilegedResourceAccountOkta) SetLinks(v UserLink)`

SetLinks sets Links field to given value.

### HasLinks

`func (o *PrivilegedResourceAccountOkta) HasLinks() bool`

HasLinks returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


