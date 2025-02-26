/*
Okta Admin Management

Allows customers to easily access the Okta Management APIs

API version: 5.1.0
Contact: devex-public@okta.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sdk

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HookPasswordAPI interface {

	/*
			CreatePasswordImportInlineHook Create an Okta Password Import Inline Hook

			Creates an Okta password import inline hook request. This is an automated request from Okta to your third-party service endpoint.
		The outbound call from Okta to your external service includes the following JSON body.

		The objects that you return in the JSON payload of your response to this Okta request are an array of one or more objects,
		which specify the Okta commands to execute.

		>**Note:** The size of your response payload must be less than 256 KB.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@return ApiCreatePasswordImportInlineHookRequest
	*/
	CreatePasswordImportInlineHook(ctx context.Context) ApiCreatePasswordImportInlineHookRequest

	// CreatePasswordImportInlineHookExecute executes the request
	//  @return PasswordImportResponse
	// TODU
	CreatePasswordImportInlineHookExecute(r ApiCreatePasswordImportInlineHookRequest) (*APIResponse, error)
}

// HookPasswordAPIService HookPasswordAPI service
type HookPasswordAPIService service

type ApiCreatePasswordImportInlineHookRequest struct {
	ctx                                   context.Context
	ApiService                            HookPasswordAPI
	createPasswordImportInlineHookRequest *CreatePasswordImportInlineHookRequest
	// TODU
	data       interface{}
	retryCount int32
}

func (r ApiCreatePasswordImportInlineHookRequest) CreatePasswordImportInlineHookRequest(createPasswordImportInlineHookRequest CreatePasswordImportInlineHookRequest) ApiCreatePasswordImportInlineHookRequest {
	r.createPasswordImportInlineHookRequest = &createPasswordImportInlineHookRequest
	return r
}

// TODU
func (r ApiCreatePasswordImportInlineHookRequest) Data(data interface{}) ApiCreatePasswordImportInlineHookRequest {
	r.data = data
	return r
}

// TODU
func (r ApiCreatePasswordImportInlineHookRequest) Execute() (*APIResponse, error) {
	return r.ApiService.CreatePasswordImportInlineHookExecute(r)
}

/*
CreatePasswordImportInlineHook Create an Okta Password Import Inline Hook

Creates an Okta password import inline hook request. This is an automated request from Okta to your third-party service endpoint.
The outbound call from Okta to your external service includes the following JSON body.

The objects that you return in the JSON payload of your response to this Okta request are an array of one or more objects,
which specify the Okta commands to execute.

>**Note:** The size of your response payload must be less than 256 KB.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreatePasswordImportInlineHookRequest
*/
// TODU

func (a *HookPasswordAPIService) CreatePasswordImportInlineHook(ctx context.Context) ApiCreatePasswordImportInlineHookRequest {
	return ApiCreatePasswordImportInlineHookRequest{
		ApiService: a,
		ctx:        ctx,
		retryCount: 0,
	}
}

// Execute executes the request
//  @return PasswordImportResponse

func (a *HookPasswordAPIService) CreatePasswordImportInlineHookExecute(r ApiCreatePasswordImportInlineHookRequest) (*APIResponse, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   interface{}
		formFiles          []formFile
		// TODU
		localVarHTTPResponse *http.Response
		localAPIResponse     *APIResponse
		err                  error
	)

	if a.client.cfg.Okta.Client.RequestTimeout > 0 {
		localctx, cancel := context.WithTimeout(r.ctx, time.Second*time.Duration(a.client.cfg.Okta.Client.RequestTimeout))
		r.ctx = localctx
		defer cancel()
	}
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "HookPasswordAPIService.CreatePasswordImportInlineHook")
	if err != nil {
		// TODU
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/your-endpoint"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// TODU
	// body params
	// localVarPostBody = r.createPasswordImportInlineHookRequest
	localVarPostBody = r.data
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["apiToken"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.PrepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		// TODU
		return nil, err
	}
	localVarHTTPResponse, err = a.client.Do(r.ctx, req)
	if err != nil {
		localAPIResponse = newAPIResponse(localVarHTTPResponse, a.client, nil)
		// TODU
		return localAPIResponse, &GenericOpenAPIError{error: err.Error()}
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		localAPIResponse = newAPIResponse(localVarHTTPResponse, a.client, nil)
		// TODU
		return localAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				localAPIResponse = newAPIResponse(localVarHTTPResponse, a.client, nil)
				// TODU
				return localAPIResponse, newErr
			}
			newErr.model = v
		}
		localAPIResponse = newAPIResponse(localVarHTTPResponse, a.client, nil)
		// TODU
		return localAPIResponse, newErr
	}

	localAPIResponse = newAPIResponse(localVarHTTPResponse, a.client, nil)
	return localAPIResponse, nil
}
