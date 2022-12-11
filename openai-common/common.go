// Package openaicommon contains common API structures and helper functions, not specific to an endpoint or model.
package openaicommon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	auth "github.com/TannerKvarfordt/gopenai/authentication"
)

const (
	// The version of the API currently implemented by this library.
	APIVersion = "v1"

	// The basis of all API endpoints.
	BaseURL = "https://api.openai.com/" + APIVersion + "/"
)

func setRequestHeaders(req *http.Request, organizationID *string) {
	if req == nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(auth.AuthHeaderKey, auth.AuthHeaderPrefix+auth.APIKey())

	if organizationID != nil {
		req.Header.Set(auth.OrgHeaderKey, *organizationID)
	} else if len(auth.DefaultOrganizationID()) != 0 {
		req.Header.Set(auth.OrgHeaderKey, auth.DefaultOrganizationID())
	}
}

// Send a request to the given OpenAI endpoint, and store the response in the provided response object.
// The method parameter should be an HTTP method, such as GET or POST.
// The organizationID parameter is optional. If provided, it will be included in the request header.
// If not provided, the authorization.DefaultOrganizationID will be used, if it is set.
func MakeRequest[RequestT any, ResponseT any](request *RequestT, response *ResponseT, endpoint, method string, organizationID *string) error {
	if request == nil {
		return errors.New("nil request provided")
	}
	if response == nil {
		return errors.New("nil response destination provided")
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	if req == nil {
		return errors.New("nil request created")
	}

	setRequestHeaders(req, organizationID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("nil response received")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if respBody == nil {
		return errors.New("unable to parse response body")
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return err
	}

	return nil
}
