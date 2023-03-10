package main

import (
	"context"

	httptransport "github.com/AraiYuno/golang-rest-api/internal/api/http"
	"github.com/AraiYuno/golang-rest-api/internal/config"
	"github.com/AraiYuno/golang-rest-api/internal/core/app"
	"github.com/AraiYuno/golang-rest-api/internal/core/drivers/psql"
	"github.com/AraiYuno/golang-rest-api/internal/core/listeners/http"
	"github.com/AraiYuno/golang-rest-api/internal/core/logging"
	"github.com/AraiYuno/golang-rest-api/internal/events"
	"github.com/AraiYuno/golang-rest-api/internal/users"
	"github.com/AraiYuno/golang-rest-api/internal/users/store"
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	// Load configuration from config/config.yaml which contains details such as DB connection params
	cfg, err := config.Load(ctx)
	if err != nil {
		return nil, err
	}

	// Connect to the postgres DB
	db, err := initDatabase(ctx, cfg, a)
	if err != nil {
		return nil, err
	}

	// Run our migrations which will update the DB or create it if it doesn't exist
	if err := db.MigratePostgres(ctx, "file://migrations"); err != nil {
		return nil, err
	}
	a.OnShutdown(func() {
		// Temp for development so database is cleared on shutdown
		if err := db.RevertMigrations(ctx, "file://migrations"); err != nil {
			logging.From(ctx).Error("failed to revert migrations", zap.Error(err))
		}
	})

	// Instantiate and connect all our classes
	us := store.New(db.GetDB())
	e := events.New()
	u := users.New(us, e)

	httpServer := httptransport.New(u, db.GetDB())

	// Create a HTTP server
	h, err := http.New(httpServer, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	// Start listening for HTTP requests
	return []app.Listener{
		h,
	}, nil
}

func initDatabase(ctx context.Context, cfg *config.Config, a *app.App) (*psql.Driver, error) {
	db := psql.New(cfg.PSQL)

	err := backoff.Retry(func() error {
		return db.Connect(ctx)
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	a.OnShutdown(func() {
		// Shutdown connection when server terminated
		logging.From(ctx).Info("shutting down db connection")
		if err := db.Close(ctx); err != nil {
			logging.From(ctx).Error("failed to close db connection", zap.Error(err))
		}
	})

	return db, nil
}
