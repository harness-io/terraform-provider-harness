package internal

import (
	"context"
	"github.com/harness/harness-go-sdk/harness/cd"
	"github.com/harness/harness-go-sdk/harness/chaos"
	"github.com/harness/harness-go-sdk/harness/code"
	"github.com/harness/harness-go-sdk/harness/dbops"
	"github.com/harness/harness-go-sdk/harness/har"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/harness-go-sdk/harness/policymgmt"
	openapi_client_nextgen "github.com/harness/harness-openapi-go-client/nextgen"
)

type Session struct {
	AccountId   string
	Endpoint    string
	CDClient    *cd.ApiClient
	PLClient    *nextgen.APIClient
	DBOpsClient *dbops.APIClient
	Client      *openapi_client_nextgen.APIClient
	CodeClient  *code.APIClient
	ChaosClient *chaos.APIClient
	HARClient  	*har.APIClient
}

func (s *Session) GetPlatformClient() (*nextgen.APIClient, context.Context) {
	return s.GetPlatformClientWithContext(nil)
}

func (s *Session) GetPlatformClientWithContext(ctx context.Context) (*nextgen.APIClient, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	return s.PLClient.WithAuthContext(ctx)
}

func (s *Session) GetClientWithContext(ctx context.Context) (*openapi_client_nextgen.APIClient, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	return s.Client.WithAuthContext(ctx)
}

func (s *Session) GetPolicyManagementClient() *policymgmt.APIClient {
	c := policymgmt.NewAPIClient(
		policymgmt.NewConfiguration(),
	)
	c.ChangeBasePath(s.Endpoint + "/pm")
	return c
}

func (s *Session) GetDBOpsClientWithContext(ctx context.Context) (*dbops.APIClient, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	return s.DBOpsClient.WithAuthContext(ctx)
}

func (s *Session) GetCodeClientWithContext(ctx context.Context) (*code.APIClient, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	return s.CodeClient.WithAuthContext(ctx)
}

func (s *Session) GetChaosClientWithContext(ctx context.Context) (*chaos.APIClient, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	return s.ChaosClient.WithAuthContext(ctx)
}

func (s *Session) GetHarClientWithContext(ctx context.Context) (*har.APIClient, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	return s.HARClient.WithAuthContext(ctx)
}
