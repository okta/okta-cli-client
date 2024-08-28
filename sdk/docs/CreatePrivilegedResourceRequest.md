# CreatePrivilegedResourceRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ResourceId** | **string** | The user ID associated with the Okta privileged resource | 
**Credentials** | [**PrivilegedResourceAccountOktaAllOfCredentials**](PrivilegedResourceAccountOktaAllOfCredentials.md) |  | 
**Links** | Pointer to **map[string]interface{}** | Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the current status of an application using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-06) specification. This object is used for dynamic discovery of related resources and lifecycle operations. | [optional] [readonly] 
**Created** | Pointer to **time.Time** | Timestamp when the object was created | [optional] [readonly] 
**CredentialLastChanged** | Pointer to **time.Time** | Timestamp when the credential was last changed | [optional] [readonly] 
**CredentialLastSyncState** | Pointer to **string** | Current credential sync status of the privileged resource | [optional] [readonly] 
**Id** | Pointer to **string** | ID of the privileged resource | [optional] [readonly] 
**LastUpdated** | Pointer to **time.Time** | Timestamp when the object was last updated | [optional] [readonly] 
**Profile** | Pointer to **map[string]map[string]interface{}** | Specific profile properties for the privileged account | [optional] [readonly] 
**ResourceType** | Pointer to **string** | The type of the resource | [optional] 
**Status** | Pointer to **string** | Current status of the privileged resource | [optional] [readonly] 

## Methods

### NewCreatePrivilegedResourceRequest

`func NewCreatePrivilegedResourceRequest(resourceId string, credentials PrivilegedResourceAccountOktaAllOfCredentials, ) *CreatePrivilegedResourceRequest`

NewCreatePrivilegedResourceRequest instantiates a new CreatePrivilegedResourceRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePrivilegedResourceRequestWithDefaults

`func NewCreatePrivilegedResourceRequestWithDefaults() *CreatePrivilegedResourceRequest`

NewCreatePrivilegedResourceRequestWithDefaults instantiates a new CreatePrivilegedResourceRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetResourceId

`func (o *CreatePrivilegedResourceRequest) GetResourceId() string`

GetResourceId returns the ResourceId field if non-nil, zero value otherwise.

### GetResourceIdOk

`func (o *CreatePrivilegedResourceRequest) GetResourceIdOk() (*string, bool)`

GetResourceIdOk returns a tuple with the ResourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceId

`func (o *CreatePrivilegedResourceRequest) SetResourceId(v string)`

SetResourceId sets ResourceId field to given value.


### GetCredentials

`func (o *CreatePrivilegedResourceRequest) GetCredentials() PrivilegedResourceAccountOktaAllOfCredentials`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *CreatePrivilegedResourceRequest) GetCredentialsOk() (*PrivilegedResourceAccountOktaAllOfCredentials, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *CreatePrivilegedResourceRequest) SetCredentials(v PrivilegedResourceAccountOktaAllOfCredentials)`

SetCredentials sets Credentials field to given value.


### GetLinks

`func (o *CreatePrivilegedResourceRequest) GetLinks() map[string]interface{}`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *CreatePrivilegedResourceRequest) GetLinksOk() (*map[string]interface{}, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *CreatePrivilegedResourceRequest) SetLinks(v map[string]interface{})`

SetLinks sets Links field to given value.

### HasLinks

`func (o *CreatePrivilegedResourceRequest) HasLinks() bool`

HasLinks returns a boolean if a field has been set.

### GetCreated

`func (o *CreatePrivilegedResourceRequest) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *CreatePrivilegedResourceRequest) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *CreatePrivilegedResourceRequest) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *CreatePrivilegedResourceRequest) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetCredentialLastChanged

`func (o *CreatePrivilegedResourceRequest) GetCredentialLastChanged() time.Time`

GetCredentialLastChanged returns the CredentialLastChanged field if non-nil, zero value otherwise.

### GetCredentialLastChangedOk

`func (o *CreatePrivilegedResourceRequest) GetCredentialLastChangedOk() (*time.Time, bool)`

GetCredentialLastChangedOk returns a tuple with the CredentialLastChanged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentialLastChanged

`func (o *CreatePrivilegedResourceRequest) SetCredentialLastChanged(v time.Time)`

SetCredentialLastChanged sets CredentialLastChanged field to given value.

### HasCredentialLastChanged

`func (o *CreatePrivilegedResourceRequest) HasCredentialLastChanged() bool`

HasCredentialLastChanged returns a boolean if a field has been set.

### GetCredentialLastSyncState

`func (o *CreatePrivilegedResourceRequest) GetCredentialLastSyncState() string`

GetCredentialLastSyncState returns the CredentialLastSyncState field if non-nil, zero value otherwise.

### GetCredentialLastSyncStateOk

`func (o *CreatePrivilegedResourceRequest) GetCredentialLastSyncStateOk() (*string, bool)`

GetCredentialLastSyncStateOk returns a tuple with the CredentialLastSyncState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentialLastSyncState

`func (o *CreatePrivilegedResourceRequest) SetCredentialLastSyncState(v string)`

SetCredentialLastSyncState sets CredentialLastSyncState field to given value.

### HasCredentialLastSyncState

`func (o *CreatePrivilegedResourceRequest) HasCredentialLastSyncState() bool`

HasCredentialLastSyncState returns a boolean if a field has been set.

### GetId

`func (o *CreatePrivilegedResourceRequest) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CreatePrivilegedResourceRequest) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CreatePrivilegedResourceRequest) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CreatePrivilegedResourceRequest) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLastUpdated

`func (o *CreatePrivilegedResourceRequest) GetLastUpdated() time.Time`

GetLastUpdated returns the LastUpdated field if non-nil, zero value otherwise.

### GetLastUpdatedOk

`func (o *CreatePrivilegedResourceRequest) GetLastUpdatedOk() (*time.Time, bool)`

GetLastUpdatedOk returns a tuple with the LastUpdated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUpdated

`func (o *CreatePrivilegedResourceRequest) SetLastUpdated(v time.Time)`

SetLastUpdated sets LastUpdated field to given value.

### HasLastUpdated

`func (o *CreatePrivilegedResourceRequest) HasLastUpdated() bool`

HasLastUpdated returns a boolean if a field has been set.

### GetProfile

`func (o *CreatePrivilegedResourceRequest) GetProfile() map[string]map[string]interface{}`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *CreatePrivilegedResourceRequest) GetProfileOk() (*map[string]map[string]interface{}, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *CreatePrivilegedResourceRequest) SetProfile(v map[string]map[string]interface{})`

SetProfile sets Profile field to given value.

### HasProfile

`func (o *CreatePrivilegedResourceRequest) HasProfile() bool`

HasProfile returns a boolean if a field has been set.

### GetResourceType

`func (o *CreatePrivilegedResourceRequest) GetResourceType() string`

GetResourceType returns the ResourceType field if non-nil, zero value otherwise.

### GetResourceTypeOk

`func (o *CreatePrivilegedResourceRequest) GetResourceTypeOk() (*string, bool)`

GetResourceTypeOk returns a tuple with the ResourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceType

`func (o *CreatePrivilegedResourceRequest) SetResourceType(v string)`

SetResourceType sets ResourceType field to given value.

### HasResourceType

`func (o *CreatePrivilegedResourceRequest) HasResourceType() bool`

HasResourceType returns a boolean if a field has been set.

### GetStatus

`func (o *CreatePrivilegedResourceRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CreatePrivilegedResourceRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CreatePrivilegedResourceRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CreatePrivilegedResourceRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


