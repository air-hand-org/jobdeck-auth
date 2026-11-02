package repository

import (
	"github.com/air-hand-org/jobdeck-auth/internal/domain"
	authdb "github.com/air-hand-org/jobdeck-auth/internal/repository/gen"
)

func localAuthClientFromDB(input authdb.AuthLocalAuthClientDetail) domain.LocalAuthClient {
	return domain.LocalAuthClient{
		ActorId:    input.ActorID,
		ClientId:   input.ClientID,
		SecretHash: input.SecretHash,
		Name:       input.Name,
		Type:       input.Type,
		Active:     input.Active,
	}
}

func actorFromDB(input authdb.AuthActor) domain.Actor {
	return domain.Actor{
		ActorID: input.ActorID,
		Name:    input.Name,
		Type:    input.Type,
		Active:  input.Active,
	}
}
