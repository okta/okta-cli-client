# \YourOinIntegrationsAPI

All URIs are relative to *https://subdomain.okta.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSubmission**](YourOinIntegrationsAPI.md#CreateSubmission) | **Post** /integrations/api/v1/submissions | Create an OIN Integration
[**GetSubmissionByOperationId**](YourOinIntegrationsAPI.md#GetSubmissionByOperationId) | **Get** /integrations/api/v1/submissions/{submissionId} | Retrieve an OIN Integration
[**GetSubmissionTestInfo**](YourOinIntegrationsAPI.md#GetSubmissionTestInfo) | **Get** /integrations/api/v1/submissions/{submissionId}/testing | Retrieve an OIN Integration Testing Information
[**ListSubmissions**](YourOinIntegrationsAPI.md#ListSubmissions) | **Get** /integrations/api/v1/submissions | List all OIN Integrations
[**ReplaceSubmission**](YourOinIntegrationsAPI.md#ReplaceSubmission) | **Put** /integrations/api/v1/submissions/{submissionId} | Replace an OIN Integration
[**SubmitSubmission**](YourOinIntegrationsAPI.md#SubmitSubmission) | **Post** /integrations/api/v1/submissions/{submissionId}/submit | Submit an OIN Integration
[**UploadSubmissionLogo**](YourOinIntegrationsAPI.md#UploadSubmissionLogo) | **Post** /integrations/api/v1/submissions/logo | Upload an OIN Integration logo
[**UpsertSubmissionTestInfo**](YourOinIntegrationsAPI.md#UpsertSubmissionTestInfo) | **Put** /integrations/api/v1/submissions/{submissionId}/testing | Upsert an OIN Integration Testing Information



## CreateSubmission

> SubmissionResponse CreateSubmission(ctx).SubmissionRequest(submissionRequest).Execute()

Create an OIN Integration



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
    submissionRequest := *openapiclient.NewSubmissionRequest("Your one source for in-season strawberry deals. Okta's Strawberry Central integration allow users to securely access those sweet deals.", "https://acme.okta.com/bc/image/fileStoreRecord?id=fs03xxd3KmkDBwJU80g4", "Strawberry Central") // SubmissionRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.YourOinIntegrationsAPI.CreateSubmission(context.Background()).SubmissionRequest(submissionRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.CreateSubmission``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateSubmission`: SubmissionResponse
    fmt.Fprintf(os.Stdout, "Response from `YourOinIntegrationsAPI.CreateSubmission`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateSubmissionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **submissionRequest** | [**SubmissionRequest**](SubmissionRequest.md) |  | 

### Return type

[**SubmissionResponse**](SubmissionResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSubmissionByOperationId

> SubmissionResponse GetSubmissionByOperationId(ctx, submissionId).Execute()

Retrieve an OIN Integration



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
    submissionId := "acme_submissionapp_1" // string | OIN Integration ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.YourOinIntegrationsAPI.GetSubmissionByOperationId(context.Background(), submissionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.GetSubmissionByOperationId``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetSubmissionByOperationId`: SubmissionResponse
    fmt.Fprintf(os.Stdout, "Response from `YourOinIntegrationsAPI.GetSubmissionByOperationId`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**submissionId** | **string** | OIN Integration ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSubmissionByOperationIdRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**SubmissionResponse**](SubmissionResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSubmissionTestInfo

> TestInfo GetSubmissionTestInfo(ctx, submissionId).Execute()

Retrieve an OIN Integration Testing Information



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
    submissionId := "acme_submissionapp_1" // string | OIN Integration ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.YourOinIntegrationsAPI.GetSubmissionTestInfo(context.Background(), submissionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.GetSubmissionTestInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetSubmissionTestInfo`: TestInfo
    fmt.Fprintf(os.Stdout, "Response from `YourOinIntegrationsAPI.GetSubmissionTestInfo`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**submissionId** | **string** | OIN Integration ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSubmissionTestInfoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**TestInfo**](TestInfo.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSubmissions

> []SubmissionResponse ListSubmissions(ctx).Limit(limit).After(after).Execute()

List all OIN Integrations



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
    limit := int32(56) // int32 | A limit on the number of objects to return (optional) (default to 20)
    after := "after_example" // string | Specify the pagination cursor (OIN Integration instance `id`) for the next page of OIN Integrations. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.YourOinIntegrationsAPI.ListSubmissions(context.Background()).Limit(limit).After(after).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.ListSubmissions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListSubmissions`: []SubmissionResponse
    fmt.Fprintf(os.Stdout, "Response from `YourOinIntegrationsAPI.ListSubmissions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListSubmissionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | A limit on the number of objects to return | [default to 20]
 **after** | **string** | Specify the pagination cursor (OIN Integration instance &#x60;id&#x60;) for the next page of OIN Integrations. | 

### Return type

[**[]SubmissionResponse**](SubmissionResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReplaceSubmission

> SubmissionResponse ReplaceSubmission(ctx, submissionId).SubmissionRequest(submissionRequest).Execute()

Replace an OIN Integration



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
    submissionId := "acme_submissionapp_1" // string | OIN Integration ID
    submissionRequest := *openapiclient.NewSubmissionRequest("Your one source for in-season strawberry deals. Okta's Strawberry Central integration allow users to securely access those sweet deals.", "https://acme.okta.com/bc/image/fileStoreRecord?id=fs03xxd3KmkDBwJU80g4", "Strawberry Central") // SubmissionRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.YourOinIntegrationsAPI.ReplaceSubmission(context.Background(), submissionId).SubmissionRequest(submissionRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.ReplaceSubmission``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ReplaceSubmission`: SubmissionResponse
    fmt.Fprintf(os.Stdout, "Response from `YourOinIntegrationsAPI.ReplaceSubmission`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**submissionId** | **string** | OIN Integration ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiReplaceSubmissionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **submissionRequest** | [**SubmissionRequest**](SubmissionRequest.md) |  | 

### Return type

[**SubmissionResponse**](SubmissionResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SubmitSubmission

> SubmitSubmission(ctx, submissionId).Execute()

Submit an OIN Integration



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
    submissionId := "acme_submissionapp_1" // string | OIN Integration ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.YourOinIntegrationsAPI.SubmitSubmission(context.Background(), submissionId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.SubmitSubmission``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**submissionId** | **string** | OIN Integration ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiSubmitSubmissionRequest struct via the builder pattern


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


## UploadSubmissionLogo

> UploadSubmissionLogo(ctx).File(file).Execute()

Upload an OIN Integration logo



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
    file := os.NewFile(1234, "some_file") // *os.File |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.YourOinIntegrationsAPI.UploadSubmissionLogo(context.Background()).File(file).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.UploadSubmissionLogo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUploadSubmissionLogoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | ***os.File** |  | 

### Return type

 (empty response body)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpsertSubmissionTestInfo

> TestInfo UpsertSubmissionTestInfo(ctx, submissionId).TestInfo(testInfo).Execute()

Upsert an OIN Integration Testing Information



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
    submissionId := "acme_submissionapp_1" // string | OIN Integration ID
    testInfo := *openapiclient.NewTestInfo("strawberry.support@example.com") // TestInfo |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.YourOinIntegrationsAPI.UpsertSubmissionTestInfo(context.Background(), submissionId).TestInfo(testInfo).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `YourOinIntegrationsAPI.UpsertSubmissionTestInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpsertSubmissionTestInfo`: TestInfo
    fmt.Fprintf(os.Stdout, "Response from `YourOinIntegrationsAPI.UpsertSubmissionTestInfo`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**submissionId** | **string** | OIN Integration ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpsertSubmissionTestInfoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **testInfo** | [**TestInfo**](TestInfo.md) |  | 

### Return type

[**TestInfo**](TestInfo.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

