package external

import (
	"fmt"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
)

const ExternalProviderGitHub app.ExternalProvider = "github"

type externalResolverRegistry struct {
	resolvers map[app.ExternalProvider]app.ExternalIdentityResolver
}

var _ app.ExternalIdentityResolverRegistry = (*externalResolverRegistry)(nil)

func (r *externalResolverRegistry) Get(provider app.ExternalProvider) (app.ExternalIdentityResolver, error) {
	if resolver, ok := r.resolvers[provider]; ok {
		return resolver, nil
	}
	return nil, fmt.Errorf("Provider [%s] is not supported", provider)
}

func NewExternalIdentityResolverRegistry() *externalResolverRegistry {
	return &externalResolverRegistry{
		resolvers: map[app.ExternalProvider]app.ExternalIdentityResolver{
			ExternalProviderGitHub: newGitHubResolver(string(ExternalProviderGitHub)),
		},
	}
}
