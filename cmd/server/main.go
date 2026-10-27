// Package main runs the auth HTTP and RPC server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/air-hand-org/jobdeck-auth/internal/app"
	"github.com/air-hand-org/jobdeck-auth/internal/repository"
	"github.com/air-hand-org/jobdeck-auth/internal/secret"
	"github.com/air-hand-org/jobdeck-auth/internal/telemetry"
	"github.com/air-hand-org/jobdeck-auth/internal/token"
	"github.com/air-hand-org/jobdeck-auth/internal/transport"

	jwksserver "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/openapi/v1/jwks/server"
	"github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/proto/v1/authv1connect"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/grpchealth/v2"
	"connectrpc.com/grpcreflect/v2"
	"connectrpc.com/validate"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runHealthCheck(ctx context.Context, checker *grpchealth.StaticChecker, pool *pgxpool.Pool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	check := func() {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		if err := pool.Ping(pingCtx); err != nil {
			checker.SetStatus("", grpchealth.StatusNotServing)
			return
		}

		checker.SetStatus("", grpchealth.StatusServing)
	}

	check()

	for {
		select {
		case <-ctx.Done():
			checker.SetStatus("", grpchealth.StatusNotServing)
			return
		case <-ticker.C:
			check()
		}
	}
}

func main() {
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	shutdownTelemetry, err := telemetry.Init(signalCtx, "auth")
	if err != nil {
		log.Fatalf("initialize telemetry: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := shutdownTelemetry(ctx); err != nil {
			log.Printf("shutdown telemetry: %v", err)
		}
	}()

	connStr := os.Getenv("POSTGRESDB_URI")
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("failed to parse conn string %v", err)
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	pool, err := pgxpool.NewWithConfig(signalCtx, cfg)
	if err != nil {
		log.Fatalf("failed to connect postgres %v", err)
	}
	defer pool.Close()

	store := repository.NewStore(pool)

	jwtKeyID := os.Getenv("AUTH_JWT_KEY_ID")
	issuer, err := token.NewJWTIssuer(os.Getenv("AUTH_JWT_PRIVATE_KEY"), jwtKeyID)
	if err != nil {
		log.Fatalf("failed to new jwt isser %v", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to new proto validate %v", err)
	}

	jwksController := jwksserver.NewDefaultAPIController(transport.NewJWKSService(issuer.PublicKey(), jwtKeyID))
	mux := http.NewServeMux()
	mux.Handle("/", jwksserver.NewRouter(jwksController))

	authService := transport.NewAuthServiceHandler(
		app.NewAuthUsecase(
			store, issuer, secret.NewPHCBuilderFactory(),
		),
	)
	server := connect.NewServer(
		validate.NewServerInterceptor(validate.WithValidator(validator)),
	)
	authv1connect.RegisterAuthServiceHandler(server, authService)

	statusServer := connect.NewServer()
	healthcheck := grpchealth.NewStaticChecker()
	healthcheck.SetStatus("", grpchealth.StatusNotServing)
	grpchealth.Register(statusServer, healthcheck)
	if strings.EqualFold(os.Getenv("GRPC_REFLECTION_ENABLED"), "true") {
		grpcreflect.Register(
			statusServer,
			grpcreflect.WithNamer(grpcreflect.NamerFunc(func() []string {
				return []string{
					authv1connect.AuthServiceName,
					grpchealth.HealthV1ServiceName,
					grpcreflect.ReflectV1ServiceName,
					grpcreflect.ReflectV1AlphaServiceName,
				}
			})),
		)
	}
	connecthttp.Mount(mux, server)
	connecthttp.Mount(mux, statusServer)

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           otelhttp.NewHandler(mux, "auth"),
		Protocols:         protocols,
		ReadHeaderTimeout: 5 * time.Second,
	}

	g, ctx := errgroup.WithContext(signalCtx)

	g.Go(func() error {
		runHealthCheck(ctx, healthcheck, pool)
		return nil
	})

	g.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("%+v", err)
	}
}
