# \ResourceSelectorsAPI

All URIs are relative to *https://subdomain.okta.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateResourceSelector**](ResourceSelectorsAPI.md#CreateResourceSelector) | **Post** /api/v1/resource-selectors | Create a Resource Selector
[**DeleteResourceSelector**](ResourceSelectorsAPI.md#DeleteResourceSelector) | **Delete** /api/v1/resource-selectors/{resourceSelectorId} | Delete a Resource Selector
[**GetResourceSelector**](ResourceSelectorsAPI.md#GetResourceSelector) | **Get** /api/v1/resource-selectors/{resourceSelectorId} | Retrieve a Resource Selector
[**ListResourceSelectors**](ResourceSelectorsAPI.md#ListResourceSelectors) | **Get** /api/v1/resource-selectors | List all Resource Selectors
[**UpdateResourceSelector**](ResourceSelectorsAPI.md#UpdateResourceSelector) | **Patch** /api/v1/resource-selectors/{resourceSelectorId} | Update a Resource Selector



## CreateResourceSelector

> ResourceSelectorResponseSchema CreateResourceSelector(ctx).Instance(instance).Execute()

Create a Resource Selector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    instance := *openapiclient.NewResourceSelectorCreateRequestSchema() // ResourceSelectorCreateRequestSchema | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ResourceSelectorsAPI.CreateResourceSelector(context.Background()).Instance(instance).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ResourceSelectorsAPI.CreateResourceSelector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateResourceSelector`: ResourceSelectorResponseSchema
    fmt.Fprintf(os.Stdout, "Response from `ResourceSelectorsAPI.CreateResourceSelector`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateResourceSelectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instance** | [**ResourceSelectorCreateRequestSchema**](ResourceSelectorCreateRequestSchema.md) |  | 

### Return type

[**ResourceSelectorResponseSchema**](ResourceSelectorResponseSchema.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteResourceSelector

> DeleteResourceSelector(ctx, resourceSelectorId).Execute()

Delete a Resource Selector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    resourceSelectorId := "rsl1hx31gVEa6x10v0g5" // string | `id` of a Resource Selector

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ResourceSelectorsAPI.DeleteResourceSelector(context.Background(), resourceSelectorId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ResourceSelectorsAPI.DeleteResourceSelector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**resourceSelectorId** | **string** | &#x60;id&#x60; of a Resource Selector | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteResourceSelectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetResourceSelector

> ResourceSelectorResponseSchema GetResourceSelector(ctx, resourceSelectorId).Execute()

Retrieve a Resource Selector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    resourceSelectorId := "rsl1hx31gVEa6x10v0g5" // string | `id` of a Resource Selector

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ResourceSelectorsAPI.GetResourceSelector(context.Background(), resourceSelectorId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ResourceSelectorsAPI.GetResourceSelector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetResourceSelector`: ResourceSelectorResponseSchema
    fmt.Fprintf(os.Stdout, "Response from `ResourceSelectorsAPI.GetResourceSelector`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**resourceSelectorId** | **string** | &#x60;id&#x60; of a Resource Selector | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetResourceSelectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ResourceSelectorResponseSchema**](ResourceSelectorResponseSchema.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListResourceSelectors

> ResourceSelectorsSchema ListResourceSelectors(ctx).After(after).Limit(limit).Execute()

List all Resource Selectors



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    after := "after_example" // string | Specifies the pagination cursor for the next page of Resource Selectors (optional)
    limit := int32(56) // int32 | Specifies the number of results returned. Defaults to `10` (optional) (default to 10)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ResourceSelectorsAPI.ListResourceSelectors(context.Background()).After(after).Limit(limit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ResourceSelectorsAPI.ListResourceSelectors``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListResourceSelectors`: ResourceSelectorsSchema
    fmt.Fprintf(os.Stdout, "Response from `ResourceSelectorsAPI.ListResourceSelectors`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListResourceSelectorsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **after** | **string** | Specifies the pagination cursor for the next page of Resource Selectors | 
 **limit** | **int32** | Specifies the number of results returned. Defaults to &#x60;10&#x60; | [default to 10]

### Return type

[**ResourceSelectorsSchema**](ResourceSelectorsSchema.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateResourceSelector

> ResourceSelectorResponseSchema UpdateResourceSelector(ctx, resourceSelectorId).Instance(instance).Execute()

Update a Resource Selector



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    resourceSelectorId := "rsl1hx31gVEa6x10v0g5" // string | `id` of a Resource Selector
    instance := *openapiclient.NewResourceSelectorPatchRequestSchema() // ResourceSelectorPatchRequestSchema | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ResourceSelectorsAPI.UpdateResourceSelector(context.Background(), resourceSelectorId).Instance(instance).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ResourceSelectorsAPI.UpdateResourceSelector``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateResourceSelector`: ResourceSelectorResponseSchema
    fmt.Fprintf(os.Stdout, "Response from `ResourceSelectorsAPI.UpdateResourceSelector`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**resourceSelectorId** | **string** | &#x60;id&#x60; of a Resource Selector | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateResourceSelectorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **instance** | [**ResourceSelectorPatchRequestSchema**](ResourceSelectorPatchRequestSchema.md) |  | 

### Return type

[**ResourceSelectorResponseSchema**](ResourceSelectorResponseSchema.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

