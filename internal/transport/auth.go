// Package transport exposes auth application services over HTTP and RPC.
package transport

import (
	"context"
	"errors"

	"github.com/air-hand-org/jobdeck-auth/internal/app"

	authv1 "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1"
	"github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1/authv1connect"

	"connectrpc.com/connect"
)

// AuthServiceHandler exposes the auth use case over Connect, gRPC, and gRPC-Web.
type AuthServiceHandler struct {
	usecase app.AuthUsecase
}

var _ authv1connect.AuthServiceHandler = (*AuthServiceHandler)(nil)

// AuthenticateLocalClient authenticates a local client and returns an access token.
func (s *AuthServiceHandler) AuthenticateLocalClient(ctx context.Context, req *authv1.AuthenticateLocalClientRequest) (*authv1.AuthenticateLocalClientResponse, error) {
	token, err := s.usecase.Authenticate(ctx, req.GetClientId(), req.GetClientSecret(), req.GetAudience())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("failed to authenticate"))
	}
	return &authv1.AuthenticateLocalClientResponse{
		Token: &authv1.Token{
			AccessToken:      token.Token,
			TokenType:        token.TokenType,
			ExpiresInSeconds: token.ExpiresInSec,
		},
	}, nil
}

// ExchangeGoogleIdToken is reserved for a future Google identity integration.
//
//nolint:revive // The protobuf-generated handler interface fixes this RPC method name.
func (s *AuthServiceHandler) ExchangeGoogleIdToken(context.Context, *authv1.ExchangeGoogleIdTokenRequest) (*authv1.ExchangeGoogleIdTokenResponse, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ExchangeGoogleIdToken is not implemented"))
}

// NewAuthServiceHandler creates an auth RPC handler.
func NewAuthServiceHandler(usecase *app.AuthUsecase) *AuthServiceHandler {
	return &AuthServiceHandler{
		usecase: *usecase,
	}
}
