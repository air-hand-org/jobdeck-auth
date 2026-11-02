package external

import (
	"testing"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestExternalProviderResolverRegistryGet(t *testing.T) {
	t.Parallel()

	tests := []struct{
		name string
		provider app.ExternalProvider
	}{
		{
			name: "github",
			provider: ExternalProviderGitHub,
		},
	}

	registry := NewExternalIdentityResolverRegistry()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := registry.Get(tt.provider)
			assert.NoError(t, err)
		})
	}
}
