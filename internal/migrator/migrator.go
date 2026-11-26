// Package migrator applies auth database migrations.
package migrator

import (
	"context"
	"database/sql"

	migrations "github.com/air-hand-org/jobdeck-auth/internal/db/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib" // Register the pgx database driver.
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func configure() error {
	goose.SetBaseFS(migrations.FS)
	goose.SetTableName("public.goose_db_version")
	return goose.SetDialect("postgres")
}

// Up applies all pending auth migrations.
func Up(ctx context.Context, connectionString string) error {
	return withDB(ctx, connectionString, func(db *sql.DB) error {
		return up(ctx, db)
	})
}

// DownTo0 rolls back all auth migrations.
func DownTo0(ctx context.Context, connectionString string) error {
	return withDB(ctx, connectionString, func(db *sql.DB) error {
		return downTo0(ctx, db)
	})
}

func withDB(ctx context.Context, connectionString string, fn func(*sql.DB) error) error {
	// FIXME: implement to other services
	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return err
	}
//	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
//		_, err := conn.Exec(ctx, "SET ROLE migrator;")
//		return err
//	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() { pool.Close() }()

	db := stdlib.OpenDBFromPool(pool)
	defer func() { _ = db.Close() }()

	return fn(db)
}

func up(ctx context.Context, db *sql.DB) error {
	if err := configure(); err != nil {
		return err
	}
	return goose.UpContext(ctx, db, migrationsDir)
}

func downTo0(ctx context.Context, db *sql.DB) error {
	if err := configure(); err != nil {
		return err
	}
	return goose.DownToContext(ctx, db, migrationsDir, 0)
}
