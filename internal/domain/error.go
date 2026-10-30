package domain

import "errors"

// ErrAuthClientNotFound indicates that an authentication client was not found.
var ErrAuthClientNotFound = errors.New("auth client is not found")
var ErrActorIdentityNotFound = errors.New("actor identity is not found")
var ErrActorDeactivated = errors.New("actor is deactivated")
