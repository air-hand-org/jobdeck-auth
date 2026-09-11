// Package main runs the auth gRPC server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
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

	authv1 "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/go/v1"
	jwksserver "github.com/air-hand-org/api-contracts/packages/jobdeck/auth/http/go/v1/jwks/server"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runHealthCheck(ctx context.Context, h *health.Server, pool *pgxpool.Pool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	check := func() {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		if err := pool.Ping(pingCtx); err != nil {
			h.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
			return
		}

		h.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
	}

	check()

	for {
		select {
		case <-ctx.Done():
			h.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
			return
		case <-ticker.C:
			check()
		}
	}
}

func main() {
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9091"
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
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

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
		),
	)
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s, healthcheck)
	authv1.RegisterAuthServiceServer(s, transport.NewAuthServiceServer(
		app.NewAuthUsecase(
			store, issuer, secret.NewPHCBuilderFactory(),
		),
	))
	if strings.EqualFold(os.Getenv("GRPC_REFLECTION_ENABLED"), "true") {
		reflection.Register(s)
	}
	jwksController := jwksserver.NewDefaultAPIController(transport.NewJWKSService(issuer.PublicKey(), jwtKeyID))
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", httpPort),
		Handler:           otelhttp.NewHandler(jwksserver.NewRouter(jwksController), "jwks"),
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
		s.GracefulStop()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	})

	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			return err
		}
		if err := s.Serve(lis); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("%+v", err)
	}
}
