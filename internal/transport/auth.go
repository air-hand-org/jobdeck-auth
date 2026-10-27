// Package transport exposes auth application services over HTTP and RPC.
package transport

import (
	"context"

	"github.com/air-hand-org/jobdeck-auth/internal/app"

	authv1 "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1"
	"github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1/authv1connect"

	"connectrpc.com/connect/v2"
)

// AuthServiceHandler exposes the auth use case over Connect, gRPC, and gRPC-Web.
type AuthServiceHandler struct {
	usecase app.AuthUsecase
	authv1connect.UnimplementedAuthServiceHandler
}

var _ authv1connect.AuthServiceHandler = (*AuthServiceHandler)(nil)

// AuthenticateLocalClient authenticates a local client and returns an access token.
func (s *AuthServiceHandler) AuthenticateLocalClient(ctx context.Context, req *authv1.AuthenticateLocalClientRequest) (*authv1.AuthenticateLocalClientResponse, error) {
	token, err := s.usecase.Authenticate(ctx, req.GetClientId(), req.GetClientSecret(), req.GetAudience())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, "failed to authenticate").WithCause(err)
	}
	return &authv1.AuthenticateLocalClientResponse{
		Token: &authv1.Token{
			AccessToken:      token.Token,
			TokenType:        token.TokenType,
			ExpiresInSeconds: token.ExpiresInSec,
		},
	}, nil
}

// NewAuthServiceHandler creates an auth RPC handler.
func NewAuthServiceHandler(usecase *app.AuthUsecase) *AuthServiceHandler {
	return &AuthServiceHandler{
		usecase: *usecase,
	}
}
