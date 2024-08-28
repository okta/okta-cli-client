# \PrivilegedResourceAPI

All URIs are relative to *https://subdomain.okta.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ClaimPrivilegedResource**](PrivilegedResourceAPI.md#ClaimPrivilegedResource) | **Post** /api/v1/privileged-resource/{id}/claim | Claim a privileged resource for management
[**CreatePrivilegedResource**](PrivilegedResourceAPI.md#CreatePrivilegedResource) | **Post** /api/v1/privileged-resource | Create a privileged resource
[**DeletePrivilegedResource**](PrivilegedResourceAPI.md#DeletePrivilegedResource) | **Delete** /api/v1/privileged-resource/{id} | Delete a privileged resource
[**GetPrivilegedResource**](PrivilegedResourceAPI.md#GetPrivilegedResource) | **Get** /api/v1/privileged-resource/{id} | Retrieve a privileged resource
[**ReplacePrivilegedResource**](PrivilegedResourceAPI.md#ReplacePrivilegedResource) | **Put** /api/v1/privileged-resource/{id} | Replace a privileged resource



## ClaimPrivilegedResource

> CreatePrivilegedResourceRequest ClaimPrivilegedResource(ctx, id).Execute()

Claim a privileged resource for management



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
    id := "id_example" // string | ID of an existing privileged resource

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PrivilegedResourceAPI.ClaimPrivilegedResource(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PrivilegedResourceAPI.ClaimPrivilegedResource``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ClaimPrivilegedResource`: CreatePrivilegedResourceRequest
    fmt.Fprintf(os.Stdout, "Response from `PrivilegedResourceAPI.ClaimPrivilegedResource`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of an existing privileged resource | 

### Other Parameters

Other parameters are passed through a pointer to a apiClaimPrivilegedResourceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CreatePrivilegedResourceRequest**](CreatePrivilegedResourceRequest.md)

### Authorization

[apiToken](../README.md#apiToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreatePrivilegedResource

> CreatePrivilegedResourceRequest CreatePrivilegedResource(ctx).Body(body).Execute()

Create a privileged resource



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
    body := openapiclient.createPrivilegedResource_request{PrivilegedResourceAccountApp: openapiclient.NewPrivilegedResourceAccountApp("ContainerId_example")} // CreatePrivilegedResourceRequest | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PrivilegedResourceAPI.CreatePrivilegedResource(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PrivilegedResourceAPI.CreatePrivilegedResource``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreatePrivilegedResource`: CreatePrivilegedResourceRequest
    fmt.Fprintf(os.Stdout, "Response from `PrivilegedResourceAPI.CreatePrivilegedResource`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreatePrivilegedResourceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CreatePrivilegedResourceRequest**](CreatePrivilegedResourceRequest.md) |  | 

### Return type

[**CreatePrivilegedResourceRequest**](CreatePrivilegedResourceRequest.md)

### Authorization

[apiToken](../README.md#apiToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeletePrivilegedResource

> CreatePrivilegedResourceRequest DeletePrivilegedResource(ctx, id).Execute()

Delete a privileged resource



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
    id := "id_example" // string | ID of an existing privileged resource

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PrivilegedResourceAPI.DeletePrivilegedResource(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PrivilegedResourceAPI.DeletePrivilegedResource``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeletePrivilegedResource`: CreatePrivilegedResourceRequest
    fmt.Fprintf(os.Stdout, "Response from `PrivilegedResourceAPI.DeletePrivilegedResource`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of an existing privileged resource | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeletePrivilegedResourceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CreatePrivilegedResourceRequest**](CreatePrivilegedResourceRequest.md)

### Authorization

[apiToken](../README.md#apiToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPrivilegedResource

> CreatePrivilegedResourceRequest GetPrivilegedResource(ctx, id).Execute()

Retrieve a privileged resource



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
    id := "id_example" // string | ID of an existing privileged resource

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PrivilegedResourceAPI.GetPrivilegedResource(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PrivilegedResourceAPI.GetPrivilegedResource``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetPrivilegedResource`: CreatePrivilegedResourceRequest
    fmt.Fprintf(os.Stdout, "Response from `PrivilegedResourceAPI.GetPrivilegedResource`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of an existing privileged resource | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPrivilegedResourceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CreatePrivilegedResourceRequest**](CreatePrivilegedResourceRequest.md)

### Authorization

[apiToken](../README.md#apiToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplacePrivilegedResource

> CreatePrivilegedResourceRequest ReplacePrivilegedResource(ctx, id).Body(body).Execute()

Replace a privileged resource



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
    id := "id_example" // string | ID of an existing privileged resource
    body := *openapiclient.NewPrivilegedResourceCredentials() // PrivilegedResourceCredentials | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PrivilegedResourceAPI.ReplacePrivilegedResource(context.Background(), id).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PrivilegedResourceAPI.ReplacePrivilegedResource``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplacePrivilegedResource`: CreatePrivilegedResourceRequest
    fmt.Fprintf(os.Stdout, "Response from `PrivilegedResourceAPI.ReplacePrivilegedResource`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of an existing privileged resource | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplacePrivilegedResourceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**PrivilegedResourceCredentials**](PrivilegedResourceCredentials.md) |  | 

### Return type

[**CreatePrivilegedResourceRequest**](CreatePrivilegedResourceRequest.md)

### Authorization

[apiToken](../README.md#apiToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

