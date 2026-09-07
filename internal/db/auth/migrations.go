// Package migrations embeds auth database migrations.
package migrations

import (
	"embed"

	// Enable go migrations
	_ "github.com/air-hand-org/jobdeck-auth/internal/db/auth/migrations"
)

// FS contains the authentication database migrations.
//
//go:embed migrations/*.sql
var FS embed.FS
