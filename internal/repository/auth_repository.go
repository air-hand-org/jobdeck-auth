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

func newAuthRepository(q authQuerier) *authRepository {
	return &authRepository{q: q}
}
