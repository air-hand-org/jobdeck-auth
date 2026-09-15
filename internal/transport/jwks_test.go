package transport

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	jwksserver "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/openapi/v1/jwks/server"
	"github.com/stretchr/testify/require"
)

func TestJWKSServiceServesSigningPublicKey(t *testing.T) {
	t.Parallel()

	publicKey, _, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	service := NewJWKSService(publicKey, "auth-key-2026-01")
	controller := jwksserver.NewDefaultAPIController(service)
	server := httptest.NewServer(jwksserver.NewRouter(controller))
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL + "/.well-known/jwks.json")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, response.Body.Close()) })
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "public, max-age=300", response.Header.Get("Cache-Control"))

	var body jwksserver.JsonWebKeySet
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Equal(t, []jwksserver.Ed25519JsonWebKey{{
		Alg: "EdDSA",
		Crv: "Ed25519",
		Kid: "auth-key-2026-01",
		Kty: "OKP",
		Use: "sig",
		X:   base64.RawURLEncoding.EncodeToString(publicKey),
	}}, body.Keys)
}
