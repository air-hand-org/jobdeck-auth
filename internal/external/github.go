package external

import (
	"context"
	"fmt"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	"github.com/google/go-github/v93/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

type githubResolver struct {
	issuer string
}

var _ app.ExternalIdentityResolver = (*githubResolver)(nil)

func (r *githubResolver) Resolve(ctx context.Context, userAccessToken string) (domain.ExternalIdentity, error) {
	client, err := github.NewClient(
		github.WithHTTPClient(
			oauth2.NewClient(ctx, githubauth.NewPersonalAccessTokenSource(userAccessToken)),
		),
	)
	if err != nil {
		return domain.ExternalIdentity{}, err
	}
	// use empty string to get self
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return domain.ExternalIdentity{}, err
	}
	return domain.ExternalIdentity{
		Issuer:  r.issuer,
		Subject: fmt.Sprintf("%d", user.GetID()),
		Name:    user.GetName(),
	}, nil
}

func newGitHubResolver(issuer string) *githubResolver {
	return &githubResolver{
		issuer: issuer,
	}
}
