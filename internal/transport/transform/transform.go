package transform

import (
	"fmt"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/external"

	authv1 "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1"
)

func ExternalProviderFromProto(value authv1.ExternalTokenProvider) (app.ExternalProvider, error) {
	switch value {
	case authv1.ExternalTokenProvider_EXTERNAL_TOKEN_PROVIDER_GITHUB:
		return external.ExternalProviderGitHub, nil
	default:
		var zero app.ExternalProvider
		return zero, fmt.Errorf("unsupported ExternalTokenProvider: %v", value)
	}
}
