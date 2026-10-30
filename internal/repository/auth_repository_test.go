package repository

import (
	"context"
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	authdb "github.com/air-hand-org/jobdeck-auth/internal/repository/gen"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type authQuerierStub struct {
	updatedCount int64
}

func (s authQuerierStub) GetMatchedLocalAuthClient(context.Context, string) (authdb.AuthLocalAuthClientDetail, error) {
	return authdb.AuthLocalAuthClientDetail{}, nil
}

func (s authQuerierStub) UpdateLocalAuthClientSecretHash(context.Context, authdb.UpdateLocalAuthClientSecretHashParams) (int64, error) {
	return s.updatedCount, nil
}

func (s authQuerierStub) GetMatchedActorIdentity(context.Context, authdb.GetMatchedActorIdentityParams) (authdb.AuthActor, error) {
	return authdb.AuthActor{}, domain.ErrActorIdentityNotFound
}

func (s authQuerierStub) CreateActor(context.Context, authdb.CreateActorParams) (int64, error) {
	return 0, nil
}


func TestAuthRepositoryUpdateLocalAuthClientSecretHashRequiresOneRow(t *testing.T) {
	t.Parallel()

	actorID := uuid.New()
	for _, count := range []int64{0, 2} {
		t.Run("unexpected affected rows", func(t *testing.T) {
			repository := newAuthRepository(authQuerierStub{updatedCount: count})
			err := repository.UpdateLocalAuthClientSecretHash(context.Background(), actorID, "hash")
			require.ErrorIs(t, err, domain.ErrAuthClientNotFound)
		})
	}

	repository := newAuthRepository(authQuerierStub{updatedCount: 1})
	err := repository.UpdateLocalAuthClientSecretHash(context.Background(), actorID, "hash")
	require.NoError(t, err)
}
