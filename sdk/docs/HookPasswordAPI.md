# \HookPasswordAPI

All URIs are relative to *https://subdomain.okta.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreatePasswordImportInlineHook**](HookPasswordAPI.md#CreatePasswordImportInlineHook) | **Post** /your-endpoint | Create an Okta Password Import Inline Hook



## CreatePasswordImportInlineHook

> PasswordImportResponse CreatePasswordImportInlineHook(ctx).CreatePasswordImportInlineHookRequest(createPasswordImportInlineHookRequest).Execute()

Create an Okta Password Import Inline Hook



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
    createPasswordImportInlineHookRequest := *openapiclient.NewCreatePasswordImportInlineHookRequest() // CreatePasswordImportInlineHookRequest | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.HookPasswordAPI.CreatePasswordImportInlineHook(context.Background()).CreatePasswordImportInlineHookRequest(createPasswordImportInlineHookRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `HookPasswordAPI.CreatePasswordImportInlineHook``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreatePasswordImportInlineHook`: PasswordImportResponse
    fmt.Fprintf(os.Stdout, "Response from `HookPasswordAPI.CreatePasswordImportInlineHook`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreatePasswordImportInlineHookRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createPasswordImportInlineHookRequest** | [**CreatePasswordImportInlineHookRequest**](CreatePasswordImportInlineHookRequest.md) |  | 

### Return type

[**PasswordImportResponse**](PasswordImportResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

