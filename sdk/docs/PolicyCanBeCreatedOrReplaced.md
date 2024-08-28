# PolicyCanBeCreatedOrReplaced

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Created** | Pointer to **time.Time** | Timestamp when the Policy was created | [optional] [readonly] 
**Description** | Pointer to **string** | Policy description | [optional] 
**Id** | Pointer to **string** | Policy ID | [optional] [readonly] 
**LastUpdated** | Pointer to **time.Time** | Timestamp when the Policy was last updated | [optional] [readonly] 
**Name** | Pointer to **string** | Policy name | [optional] 
**Priority** | Pointer to **int32** | Specifies the order in which this Policy is evaluated in relation to the other policies | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**System** | Pointer to **bool** | Specifies whether Okta created the Policy | [optional] 
**Type** | Pointer to **string** | All Okta orgs contain only one IdP Discovery Policy with an immutable default Rule routing to your org&#39;s sign-in page. Creating or replacing a policy with &#x60;IDP_DISCOVERY&#x60; type isn&#39;t supported. | [optional] 
**Embedded** | Pointer to **map[string]map[string]interface{}** |  | [optional] [readonly] 
**Links** | Pointer to [**LinksSelf**](LinksSelf.md) |  | [optional] 

## Methods

### NewPolicyCanBeCreatedOrReplaced

`func NewPolicyCanBeCreatedOrReplaced() *PolicyCanBeCreatedOrReplaced`

NewPolicyCanBeCreatedOrReplaced instantiates a new PolicyCanBeCreatedOrReplaced object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPolicyCanBeCreatedOrReplacedWithDefaults

`func NewPolicyCanBeCreatedOrReplacedWithDefaults() *PolicyCanBeCreatedOrReplaced`

NewPolicyCanBeCreatedOrReplacedWithDefaults instantiates a new PolicyCanBeCreatedOrReplaced object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreated

`func (o *PolicyCanBeCreatedOrReplaced) GetCreated() time.Time`

GetCreated returns the Created field if non-nil, zero value otherwise.

### GetCreatedOk

`func (o *PolicyCanBeCreatedOrReplaced) GetCreatedOk() (*time.Time, bool)`

GetCreatedOk returns a tuple with the Created field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreated

`func (o *PolicyCanBeCreatedOrReplaced) SetCreated(v time.Time)`

SetCreated sets Created field to given value.

### HasCreated

`func (o *PolicyCanBeCreatedOrReplaced) HasCreated() bool`

HasCreated returns a boolean if a field has been set.

### GetDescription

`func (o *PolicyCanBeCreatedOrReplaced) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *PolicyCanBeCreatedOrReplaced) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *PolicyCanBeCreatedOrReplaced) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *PolicyCanBeCreatedOrReplaced) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetId

`func (o *PolicyCanBeCreatedOrReplaced) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PolicyCanBeCreatedOrReplaced) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PolicyCanBeCreatedOrReplaced) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *PolicyCanBeCreatedOrReplaced) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLastUpdated

`func (o *PolicyCanBeCreatedOrReplaced) GetLastUpdated() time.Time`

GetLastUpdated returns the LastUpdated field if non-nil, zero value otherwise.

### GetLastUpdatedOk

`func (o *PolicyCanBeCreatedOrReplaced) GetLastUpdatedOk() (*time.Time, bool)`

GetLastUpdatedOk returns a tuple with the LastUpdated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUpdated

`func (o *PolicyCanBeCreatedOrReplaced) SetLastUpdated(v time.Time)`

SetLastUpdated sets LastUpdated field to given value.

### HasLastUpdated

`func (o *PolicyCanBeCreatedOrReplaced) HasLastUpdated() bool`

HasLastUpdated returns a boolean if a field has been set.

### GetName

`func (o *PolicyCanBeCreatedOrReplaced) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PolicyCanBeCreatedOrReplaced) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PolicyCanBeCreatedOrReplaced) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PolicyCanBeCreatedOrReplaced) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPriority

`func (o *PolicyCanBeCreatedOrReplaced) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *PolicyCanBeCreatedOrReplaced) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *PolicyCanBeCreatedOrReplaced) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *PolicyCanBeCreatedOrReplaced) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetStatus

`func (o *PolicyCanBeCreatedOrReplaced) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PolicyCanBeCreatedOrReplaced) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PolicyCanBeCreatedOrReplaced) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *PolicyCanBeCreatedOrReplaced) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetSystem

`func (o *PolicyCanBeCreatedOrReplaced) GetSystem() bool`

GetSystem returns the System field if non-nil, zero value otherwise.

### GetSystemOk

`func (o *PolicyCanBeCreatedOrReplaced) GetSystemOk() (*bool, bool)`

GetSystemOk returns a tuple with the System field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSystem

`func (o *PolicyCanBeCreatedOrReplaced) SetSystem(v bool)`

SetSystem sets System field to given value.

### HasSystem

`func (o *PolicyCanBeCreatedOrReplaced) HasSystem() bool`

HasSystem returns a boolean if a field has been set.

### GetType

`func (o *PolicyCanBeCreatedOrReplaced) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PolicyCanBeCreatedOrReplaced) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PolicyCanBeCreatedOrReplaced) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *PolicyCanBeCreatedOrReplaced) HasType() bool`

HasType returns a boolean if a field has been set.

### GetEmbedded

`func (o *PolicyCanBeCreatedOrReplaced) GetEmbedded() map[string]map[string]interface{}`

GetEmbedded returns the Embedded field if non-nil, zero value otherwise.

### GetEmbeddedOk

`func (o *PolicyCanBeCreatedOrReplaced) GetEmbeddedOk() (*map[string]map[string]interface{}, bool)`

GetEmbeddedOk returns a tuple with the Embedded field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmbedded

`func (o *PolicyCanBeCreatedOrReplaced) SetEmbedded(v map[string]map[string]interface{})`

SetEmbedded sets Embedded field to given value.

### HasEmbedded

`func (o *PolicyCanBeCreatedOrReplaced) HasEmbedded() bool`

HasEmbedded returns a boolean if a field has been set.

### GetLinks

`func (o *PolicyCanBeCreatedOrReplaced) GetLinks() LinksSelf`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *PolicyCanBeCreatedOrReplaced) GetLinksOk() (*LinksSelf, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *PolicyCanBeCreatedOrReplaced) SetLinks(v LinksSelf)`

SetLinks sets Links field to given value.

### HasLinks

`func (o *PolicyCanBeCreatedOrReplaced) HasLinks() bool`

HasLinks returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


