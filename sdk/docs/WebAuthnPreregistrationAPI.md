# \WebAuthnPreregistrationAPI

All URIs are relative to *https://subdomain.okta.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivatePreregistrationEnrollment**](WebAuthnPreregistrationAPI.md#ActivatePreregistrationEnrollment) | **Post** /webauthn-registration/api/v1/activate | Activate a Preregistered WebAuthn Factor
[**DeleteWebAuthnPreregistrationFactor**](WebAuthnPreregistrationAPI.md#DeleteWebAuthnPreregistrationFactor) | **Delete** /webauthn-registration/api/v1/users/{userId}/enrollments/{authenticatorEnrollmentId} | Delete a WebAuthn Preregistration Factor
[**EnrollPreregistrationEnrollment**](WebAuthnPreregistrationAPI.md#EnrollPreregistrationEnrollment) | **Post** /webauthn-registration/api/v1/enroll | Enroll a Preregistered WebAuthn Factor
[**GenerateFulfillmentRequest**](WebAuthnPreregistrationAPI.md#GenerateFulfillmentRequest) | **Post** /webauthn-registration/api/v1/initiate-fulfillment-request | Generate a Fulfillment Request
[**ListWebAuthnPreregistrationFactors**](WebAuthnPreregistrationAPI.md#ListWebAuthnPreregistrationFactors) | **Get** /webauthn-registration/api/v1/users/{userId}/enrollments | List all WebAuthn Preregistration Factors
[**SendPin**](WebAuthnPreregistrationAPI.md#SendPin) | **Post** /webauthn-registration/api/v1/send-pin | Send a PIN to user



## ActivatePreregistrationEnrollment

> EnrollmentActivationResponse ActivatePreregistrationEnrollment(ctx).Body(body).Execute()

Activate a Preregistered WebAuthn Factor



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
    body := *openapiclient.NewEnrollmentActivationRequest() // EnrollmentActivationRequest | Enrollment Activation Request (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebAuthnPreregistrationAPI.ActivatePreregistrationEnrollment(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnPreregistrationAPI.ActivatePreregistrationEnrollment``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ActivatePreregistrationEnrollment`: EnrollmentActivationResponse
    fmt.Fprintf(os.Stdout, "Response from `WebAuthnPreregistrationAPI.ActivatePreregistrationEnrollment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiActivatePreregistrationEnrollmentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**EnrollmentActivationRequest**](EnrollmentActivationRequest.md) | Enrollment Activation Request | 

### Return type

[**EnrollmentActivationResponse**](EnrollmentActivationResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteWebAuthnPreregistrationFactor

> DeleteWebAuthnPreregistrationFactor(ctx, userId, authenticatorEnrollmentId).Execute()

Delete a WebAuthn Preregistration Factor



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
    userId := "userId_example" // string | ID of an existing Okta user
    authenticatorEnrollmentId := "authenticatorEnrollmentId_example" // string | ID for a WebAuthn Preregistration Factor in Okta

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.WebAuthnPreregistrationAPI.DeleteWebAuthnPreregistrationFactor(context.Background(), userId, authenticatorEnrollmentId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnPreregistrationAPI.DeleteWebAuthnPreregistrationFactor``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | ID of an existing Okta user | 
**authenticatorEnrollmentId** | **string** | ID for a WebAuthn Preregistration Factor in Okta | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteWebAuthnPreregistrationFactorRequest struct via the builder pattern


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


## EnrollPreregistrationEnrollment

> EnrollmentInitializationResponse EnrollPreregistrationEnrollment(ctx).Body(body).Execute()

Enroll a Preregistered WebAuthn Factor



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
    body := *openapiclient.NewEnrollmentInitializationRequest() // EnrollmentInitializationRequest | Enrollment Initialization Request (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebAuthnPreregistrationAPI.EnrollPreregistrationEnrollment(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnPreregistrationAPI.EnrollPreregistrationEnrollment``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `EnrollPreregistrationEnrollment`: EnrollmentInitializationResponse
    fmt.Fprintf(os.Stdout, "Response from `WebAuthnPreregistrationAPI.EnrollPreregistrationEnrollment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEnrollPreregistrationEnrollmentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**EnrollmentInitializationRequest**](EnrollmentInitializationRequest.md) | Enrollment Initialization Request | 

### Return type

[**EnrollmentInitializationResponse**](EnrollmentInitializationResponse.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GenerateFulfillmentRequest

> GenerateFulfillmentRequest(ctx).Body(body).Execute()

Generate a Fulfillment Request



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
    body := *openapiclient.NewFulfillmentRequest() // FulfillmentRequest | Fulfillment Request (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.WebAuthnPreregistrationAPI.GenerateFulfillmentRequest(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnPreregistrationAPI.GenerateFulfillmentRequest``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGenerateFulfillmentRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**FulfillmentRequest**](FulfillmentRequest.md) | Fulfillment Request | 

### Return type

 (empty response body)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWebAuthnPreregistrationFactors

> []WebAuthnPreregistrationFactor ListWebAuthnPreregistrationFactors(ctx, userId).Execute()

List all WebAuthn Preregistration Factors



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
    userId := "userId_example" // string | ID of an existing Okta user

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.WebAuthnPreregistrationAPI.ListWebAuthnPreregistrationFactors(context.Background(), userId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnPreregistrationAPI.ListWebAuthnPreregistrationFactors``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListWebAuthnPreregistrationFactors`: []WebAuthnPreregistrationFactor
    fmt.Fprintf(os.Stdout, "Response from `WebAuthnPreregistrationAPI.ListWebAuthnPreregistrationFactors`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | ID of an existing Okta user | 

### Other Parameters

Other parameters are passed through a pointer to a apiListWebAuthnPreregistrationFactorsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]WebAuthnPreregistrationFactor**](WebAuthnPreregistrationFactor.md)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SendPin

> SendPin(ctx).Body(body).Execute()

Send a PIN to user



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
    body := *openapiclient.NewPinRequest() // PinRequest | Send PIN Request (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.WebAuthnPreregistrationAPI.SendPin(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `WebAuthnPreregistrationAPI.SendPin``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSendPinRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**PinRequest**](PinRequest.md) | Send PIN Request | 

### Return type

 (empty response body)

### Authorization

[apiToken](../README.md#apiToken), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

