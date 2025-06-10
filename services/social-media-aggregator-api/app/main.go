package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Torkel-Aannestad/coop-backend/services/social-media-aggregator-api/internal/database"
	"github.com/Torkel-Aannestad/coop-backend/services/social-media-aggregator-api/sql/migrations"
	"github.com/joho/godotenv"
)

type Config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config Config
	logger *slog.Logger
	models *database.Models
}

func main() {
	var cfg Config

	godotenv.Load()
	dsn := os.Getenv("DSN_DB_LOCAL_DEV")

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	//db
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "dsn for PG instance")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	loggerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))

	db, err := database.OpenDB(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Error("failed to open db connection")
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	app := &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
