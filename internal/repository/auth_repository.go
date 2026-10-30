// Package repository provides auth persistence implementations.
package repository

import (
	"context"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	authdb "github.com/air-hand-org/jobdeck-auth/internal/repository/gen"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

type authQuerier interface {
	GetMatchedLocalAuthClient(ctx context.Context, clientID string) (authdb.AuthLocalAuthClientDetail, error)
	UpdateLocalAuthClientSecretHash(ctx context.Context, arg authdb.UpdateLocalAuthClientSecretHashParams) (int64, error)

	CreateActor(ctx context.Context, arg authdb.CreateActorParams) (int64, error)
	GetMatchedActorIdentity(ctx context.Context, arg authdb.GetMatchedActorIdentityParams) (authdb.AuthActor, error)
}

type authRepository struct {
	q authQuerier
}

var _ app.AuthRepository = (*authRepository)(nil)

func (r *authRepository) GetLocalAuthClient(ctx context.Context, clientID string) (domain.LocalAuthClient, error) {
	row, err := r.q.GetMatchedLocalAuthClient(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.LocalAuthClient{}, domain.ErrAuthClientNotFound
		}
		return domain.LocalAuthClient{}, err
	}
	return localAuthClientFromDB(row), nil
}

func (r *authRepository) UpdateLocalAuthClientSecretHash(ctx context.Context, actorID domain.ActorId, secretHash string) error {
	updatedCount, err := r.q.UpdateLocalAuthClientSecretHash(ctx, authdb.UpdateLocalAuthClientSecretHashParams{
		ActorID:    actorID,
		SecretHash: secretHash,
	})
	if err != nil {
		return err
	}
	if updatedCount != 1 {
		return domain.ErrAuthClientNotFound
	}
	return nil
}

func (r *authRepository) GetOrCreateActorByExternalIdentity(ctx context.Context, identity domain.ExternalIdentity) (domain.Actor, error) {
	row, err := r.q.GetMatchedActorIdentity(ctx, authdb.GetMatchedActorIdentityParams{
		Issuer:  identity.Issuer,
		Subject: identity.Subject,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return domain.Actor{}, err
		}
		// Create Actor if not found
		createdCount, err := r.q.CreateActor(ctx, authdb.CreateActorParams{
			Issuer:  identity.Issuer,
			Subject: identity.Subject,
			Name:    identity.Name,
			Type:    "USER",
		})
		if err != nil {
			return domain.Actor{}, err
		}
		if createdCount != 1 {
			return domain.Actor{}, errors.New("failed to create actor")
		}
		row, err := r.q.GetMatchedActorIdentity(ctx, authdb.GetMatchedActorIdentityParams{
			Issuer:  identity.Issuer,
			Subject: identity.Subject,
		})
		if err != nil {
			return domain.Actor{}, err
		}
		return actorFromDB(row), nil
	}
	return actorFromDB(row), nil
}

func newAuthRepository(q authQuerier) *authRepository {
	return &authRepository{q: q}
}
