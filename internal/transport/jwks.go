// Package transport exposes Auth service transports.
package transport

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"net/http"

	jwksserver "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/openapi/v1/jwks/server"
)

// JWKSService serves the public key used to verify Auth-issued JWTs.
type JWKSService struct {
	publicKey ed25519.PublicKey
	keyID     string
}

var _ jwksserver.DefaultAPIServicer = (*JWKSService)(nil)

// NewJWKSService creates a JWKS service for a single Ed25519 signing key.
func NewJWKSService(publicKey ed25519.PublicKey, keyID string) *JWKSService {
	return &JWKSService{publicKey: publicKey, keyID: keyID}
}

// GetJSONWebKeySet returns the signing key in the JSON Web Key Set format.
func (s JWKSService) GetJSONWebKeySet(context.Context) (jwksserver.ImplResponse, error) {
	return jwksserver.ImplResponse{
		Code: http.StatusOK,
		Headers: map[string][]string{
			"Cache-Control": {"public, max-age=300"},
		},
		Body: jwksserver.JsonWebKeySet{Keys: []jwksserver.Ed25519JsonWebKey{{
			Alg: "EdDSA",
			Crv: "Ed25519",
			Kid: s.keyID,
			Kty: "OKP",
			Use: "sig",
			X:   base64.RawURLEncoding.EncodeToString(s.publicKey),
		}}},
	}, nil
}
