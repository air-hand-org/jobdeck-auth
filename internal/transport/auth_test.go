package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	"github.com/air-hand-org/jobdeck-auth/internal/secret"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	authv1 "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1"
	"github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1/authv1connect"
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

func TestAuthServiceHandlerProtocols(t *testing.T) {
	t.Parallel()

	factory := secret.NewPHCBuilderFactory()
	usecase := app.NewAuthUsecase(
		authStoreStub{repository: authRepositoryStub{client: domain.LocalAuthClient{
			ClientId:   "worker",
			SecretHash: factory.NewBuilder([]byte("0123456789abcdef")).PHC([]byte("secret")).String(),
			Type:       "SERVICE",
			Active:     true,
		}}},
		tokenIssuerStub{},
		factory,
	)

	path, handler := authv1connect.NewAuthServiceHandler(NewAuthServiceHandler(usecase),
		connect.WithInterceptors(validate.NewInterceptor()))
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	request := &authv1.AuthenticateLocalClientRequest{
		ClientId: "worker", ClientSecret: "secret", Audience: "jobdeck-job-api",
	}
	connectClient := authv1connect.NewAuthServiceClient(server.Client(), server.URL)
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

	_, err = authv1.NewAuthServiceClient(connection).ExchangeGoogleIdToken(context.Background(), &authv1.ExchangeGoogleIdTokenRequest{
		IdToken: "some-token", Audience: "jobdeck-job-api",
	})
	require.Equal(t, codes.Unimplemented, status.Code(err))
}
