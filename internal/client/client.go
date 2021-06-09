package client

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/micahlmartin/terraform-provider-harness/internal/envvar"
	"github.com/micahlmartin/terraform-provider-harness/internal/httphelpers"
)

type ApiClient struct {
	HTTPClient *http.Client
	Endpoint   string
	UserAgent  string
	APIKey     string
	ApiToken   string
	AccountId  string
}

func New() *ApiClient {
	return &ApiClient{
		UserAgent: "micahlmartin-harness-go-sdk-0.0.1",
		Endpoint:  DefaultApiUrl,
		AccountId: os.Getenv(envvar.HarnessAccountId),
		APIKey:    os.Getenv(envvar.HarnessApiKey),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Creates a new unauthenticated HTTP request
func (client *ApiClient) NewRequest(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", client.Endpoint, path), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set(httphelpers.HeaderUserAgent, client.UserAgent)
	return req, err
}

// Creates an HTTP request using the bearer Token for authentication
func (client *ApiClient) NewAuthorizedRequestWithBearerToken(path string) (*http.Request, error) {
	req, err := client.NewRequest(path)

	if err != nil {
		return nil, err
	}

	req.Header.Set(httphelpers.HeaderAuthorization, fmt.Sprintf("Bearer %s", client.ApiToken))
	return req, nil
}

// Creates an HTTP request using an API key for authentication
func (client *ApiClient) NewAuthorizedRequestWithApiKey(path string) (*http.Request, error) {
	req, err := client.NewRequest(path)

	if err != nil {
		return nil, err
	}

	req.Header.Set(httphelpers.HeaderApiKey, client.APIKey)
	return req, nil
}
