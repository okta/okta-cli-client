# PrivilegedResourceAccountApp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ContainerId** | **string** | The application ID associated with the privileged account | 
**Credentials** | Pointer to [**PrivilegedResourceCredentials**](PrivilegedResourceCredentials.md) |  | [optional] 
**ContainerDisplayName** | Pointer to **string** | Human-readable name of the container that owns the privileged resource | [optional] [readonly] 
**Links** | Pointer to [**AppLink**](AppLink.md) |  | [optional] 

## Methods

### NewPrivilegedResourceAccountApp

`func NewPrivilegedResourceAccountApp(containerId string, ) *PrivilegedResourceAccountApp`

NewPrivilegedResourceAccountApp instantiates a new PrivilegedResourceAccountApp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrivilegedResourceAccountAppWithDefaults

`func NewPrivilegedResourceAccountAppWithDefaults() *PrivilegedResourceAccountApp`

NewPrivilegedResourceAccountAppWithDefaults instantiates a new PrivilegedResourceAccountApp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContainerId

`func (o *PrivilegedResourceAccountApp) GetContainerId() string`

GetContainerId returns the ContainerId field if non-nil, zero value otherwise.

### GetContainerIdOk

`func (o *PrivilegedResourceAccountApp) GetContainerIdOk() (*string, bool)`

GetContainerIdOk returns a tuple with the ContainerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainerId

`func (o *PrivilegedResourceAccountApp) SetContainerId(v string)`

SetContainerId sets ContainerId field to given value.


### GetCredentials

`func (o *PrivilegedResourceAccountApp) GetCredentials() PrivilegedResourceCredentials`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *PrivilegedResourceAccountApp) GetCredentialsOk() (*PrivilegedResourceCredentials, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *PrivilegedResourceAccountApp) SetCredentials(v PrivilegedResourceCredentials)`

SetCredentials sets Credentials field to given value.

### HasCredentials

`func (o *PrivilegedResourceAccountApp) HasCredentials() bool`

HasCredentials returns a boolean if a field has been set.

### GetContainerDisplayName

`func (o *PrivilegedResourceAccountApp) GetContainerDisplayName() string`

GetContainerDisplayName returns the ContainerDisplayName field if non-nil, zero value otherwise.

### GetContainerDisplayNameOk

`func (o *PrivilegedResourceAccountApp) GetContainerDisplayNameOk() (*string, bool)`

GetContainerDisplayNameOk returns a tuple with the ContainerDisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainerDisplayName

`func (o *PrivilegedResourceAccountApp) SetContainerDisplayName(v string)`

SetContainerDisplayName sets ContainerDisplayName field to given value.

### HasContainerDisplayName

`func (o *PrivilegedResourceAccountApp) HasContainerDisplayName() bool`

HasContainerDisplayName returns a boolean if a field has been set.

### GetLinks

`func (o *PrivilegedResourceAccountApp) GetLinks() AppLink`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *PrivilegedResourceAccountApp) GetLinksOk() (*AppLink, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *PrivilegedResourceAccountApp) SetLinks(v AppLink)`

SetLinks sets Links field to given value.

### HasLinks

`func (o *PrivilegedResourceAccountApp) HasLinks() bool`

HasLinks returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


