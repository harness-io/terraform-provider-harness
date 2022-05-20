package internal

import (
	"context"

	"github.com/harness/harness-go-sdk/harness/cd"
	"github.com/harness/harness-go-sdk/harness/nextgen"
)

type Session struct {
	AccountId string
	Endpoint  string
	CDClient  *cd.ApiClient
	PLClient  *nextgen.APIClient
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
