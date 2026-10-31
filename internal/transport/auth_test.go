package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	"github.com/air-hand-org/jobdeck-auth/internal/external"
	"github.com/air-hand-org/jobdeck-auth/internal/secret"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/validate"
	authv1 "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1"
	"github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1/authv1connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type authRepositoryStub struct {
	client domain.LocalAuthClient
}

func (s authRepositoryStub) GetLocalAuthClient(context.Context, string) (domain.LocalAuthClient, error) {
	return s.client, nil
}

func (authRepositoryStub) UpdateLocalAuthClientSecretHash(context.Context, domain.ActorId, string) error {
	return nil
}

func (authRepositoryStub) GetOrCreateActorByExternalIdentity(context.Context, domain.ExternalIdentity) (domain.Actor, error) {
	return domain.Actor{}, errors.New("not found")
}

type authStoreStub struct {
	repository app.AuthRepository
}

func (s authStoreStub) Auths() app.AuthRepository { return s.repository }

func (s authStoreStub) WithinAuthTx(ctx context.Context, fn func(context.Context, app.AuthRepository) error) error {
	return fn(ctx, s.repository)
}

type tokenIssuerStub struct{}

func (tokenIssuerStub) Token(domain.ActorId, string, string, time.Time) (app.AccessToken, error) {
	return app.AccessToken{Token: "token", TokenType: "Bearer", ExpiresInSec: 60}, nil
}

//go:generate moq -out gen_auth_store_moq_test.go -pkg transport ../app AuthStore
//go:generate moq -out gen_auth_respository_moq_test.go -pkg transport ../app AuthRepository
//go:generate moq -out gen_token_issuer_moq_test.go -pkg transport ../app TokenIssuer
//go:generate moq -out gen_phc_builder_factory_moq_test.go -pkg transport ../app PHCBuilderFactory
//go:generate moq -out gen_external_identity_resolver_moq_test.go -pkg transport ../app ExternalIdentityResolver
//go:generate moq -out gen_external_identity_resolver_registry_moq_test.go -pkg transport ../app ExternalIdentityResolverRegistry

func TestAuthServiceHandlerProtocols(t *testing.T) {
	t.Parallel()

	factory := secret.NewPHCBuilderFactory()
	usecase := app.NewAuthUsecase(
		authStoreStub{repository: authRepositoryStub{
			client: domain.LocalAuthClient{
				ClientId:   "worker",
				SecretHash: factory.NewBuilder([]byte("0123456789abcdef")).PHC([]byte("secret")).String(),
				Type:       "SERVICE",
				Active:     true,
			},
		}},
		tokenIssuerStub{},
		factory,
		&ExternalIdentityResolverRegistryMock{
			GetFunc: func(app.ExternalProvider) (app.ExternalIdentityResolver, error) {
				return &ExternalIdentityResolverMock{
					ResolveFunc: func(context.Context, string) (domain.ExternalIdentity, error) {
						return domain.ExternalIdentity{}, errors.New("failed to resolve")
					},
				}, nil
			},
		},
	)

	connectServer := connect.NewServer(validate.NewServerInterceptor())
	authv1connect.RegisterAuthServiceHandler(connectServer, NewAuthServiceHandler(usecase))
	mux := http.NewServeMux()
	connecthttp.Mount(mux, connectServer)
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	request := &authv1.AuthenticateLocalClientRequest{
		ClientId: "worker", ClientSecret: "secret", Audience: "jobdeck-job-api",
	}
	connectClient := authv1connect.NewAuthServiceClient(connect.NewClient(connecthttp.NewTransport(server.Client(), server.URL)))
	response, err := connectClient.AuthenticateLocalClient(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, "token", response.GetToken().GetAccessToken())

	connection, err := grpc.NewClient(server.Listener.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(server.Client().Transport.(*http.Transport).TLSClientConfig)))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, connection.Close()) })

	response, err = authv1.NewAuthServiceClient(connection).AuthenticateLocalClient(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, "token", response.GetToken().GetAccessToken())

	_, err = authv1.NewAuthServiceClient(connection).AuthenticateLocalClient(context.Background(), &authv1.AuthenticateLocalClientRequest{
		ClientId: "worker", ClientSecret: "wrong", Audience: "jobdeck-job-api",
	})
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	_, err = authv1.NewAuthServiceClient(connection).ExchangeExternalToken(context.Background(), &authv1.ExchangeExternalTokenRequest{
		Provider: authv1.ExternalTokenProvider_EXTERNAL_TOKEN_PROVIDER_GITHUB,
		Token: "this.is.external.identity.token",
		Audience: "jobdeck-ui",
	})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthServiceHandler_ExchangeExternalTokenGitHub(t *testing.T) {
	t.Parallel()

	expectedExternalToken := "this.is.external.provider.token"
	expectedExternalID := "12235950"
	expectedActorID := uuid.Max
	expectedToken := fmt.Sprintf("token-%v", expectedActorID)
	repository := &AuthRepositoryMock {
		GetOrCreateActorByExternalIdentityFunc: func(ctx context.Context, actualIdentity domain.ExternalIdentity) (domain.Actor, error) {
			require.Equal(t, string(external.ExternalProviderGitHub), actualIdentity.Issuer)
			require.Equal(t, expectedExternalID, actualIdentity.Subject)
			return domain.Actor{
				ActorID: expectedActorID,
				Type: "USER",
				Active: true,
			}, nil
		},
	}
	usecase := app.NewAuthUsecase(
		&AuthStoreMock{
			AuthsFunc: func() app.AuthRepository {
				return repository
			},
			WithinAuthTxFunc: func(ctx context.Context, fn func(context.Context, app.AuthRepository) error) error {
				return fn(ctx, repository)
			},
		},
		&TokenIssuerMock{
			TokenFunc: func (actualActorID domain.ActorId, actualActorType string, actualAudience string, actualNow time.Time) (app.AccessToken, error) {
				require.Equal(t, expectedActorID, actualActorID)
				return app.AccessToken{Token: expectedToken, TokenType: "Bearer", ExpiresInSec: 60}, nil
			},
		},
		&PHCBuilderFactoryMock{},
		&ExternalIdentityResolverRegistryMock{
			GetFunc: func(actualProvider app.ExternalProvider) (app.ExternalIdentityResolver, error) {
				require.Equal(t, external.ExternalProviderGitHub, actualProvider)
				return &ExternalIdentityResolverMock{
					ResolveFunc: func(ctx context.Context, actualExternalToken string) (domain.ExternalIdentity, error) {
						require.Equal(t, expectedExternalToken, actualExternalToken)
						return domain.ExternalIdentity{
							Issuer: string(actualProvider),
							Subject: expectedExternalID,
						}, nil
					},
				}, nil
			},
		},
	)

	connectServer := connect.NewServer(validate.NewServerInterceptor())
	authv1connect.RegisterAuthServiceHandler(connectServer, NewAuthServiceHandler(usecase))
	mux := http.NewServeMux()
	connecthttp.Mount(mux, connectServer)
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	request := &authv1.ExchangeExternalTokenRequest{
		Provider: authv1.ExternalTokenProvider_EXTERNAL_TOKEN_PROVIDER_GITHUB,
		Token: expectedExternalToken,
		Audience: "jobdeck-ui",
	}
	connectClient := authv1connect.NewAuthServiceClient(connect.NewClient(connecthttp.NewTransport(server.Client(), server.URL)))
	response, err := connectClient.ExchangeExternalToken(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, expectedToken, response.GetToken().GetAccessToken())

	connection, err := grpc.NewClient(server.Listener.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(server.Client().Transport.(*http.Transport).TLSClientConfig)))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, connection.Close()) })

	response, err = authv1.NewAuthServiceClient(connection).ExchangeExternalToken(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, expectedToken, response.GetToken().GetAccessToken())
}
