package repository_test

import (
	"context"
	"testing"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	"github.com/air-hand-org/jobdeck-auth/internal/migrator"
	"github.com/air-hand-org/jobdeck-auth/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func testDBSetup(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	container, err := postgres.Run(
		ctx,
		"postgres:18.4-bookworm",
		postgres.WithDatabase("jobdeck"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, container)

	uri, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	require.NoError(t, migrator.Up(ctx, uri))

	pool, err := pgxpool.New(ctx, uri)
	require.NoError(t, err)
	t.Cleanup(func() {
		pool.Close()
	})
	return pool
}

func TestAuthRepositoryIntegration_GetOrCreateActorExternalUser(t *testing.T) {
	pool := testDBSetup(t, t.Context())

	store := repository.NewStore(pool)

	identity := domain.ExternalIdentity{
		Issuer: "foobar",
		Subject: "1111",
	}

	// first signup
	var createdActor domain.Actor
	err := store.WithinAuthTx(t.Context(), func(ctx context.Context, auths app.AuthRepository) error {
		actor, err := auths.GetOrCreateActorByExternalIdentity(ctx, identity)
		require.NoError(t, err)
		assert.Equal(t, "USER", actor.Type)
		assert.True(t, actor.Active)
		createdActor = actor
		return err
	})
	require.NoError(t, err)

	// second signin
	err = store.WithinAuthTx(t.Context(), func(ctx context.Context, auths app.AuthRepository) error {
		actor, err := auths.GetOrCreateActorByExternalIdentity(ctx, identity)
		require.NoError(t, err)
		assert.Equal(t, createdActor, actor)
		return err
	})
	require.NoError(t, err)
}
