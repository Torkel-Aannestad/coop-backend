package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	port             int
	env              string
	pullingFrequency time.Duration
}

type application struct {
	config    Config
	logger    *slog.Logger
	mu        sync.Mutex
	seenPosts map[string]struct{}
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 5000, "Mastodon service port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.DurationVar(&cfg.pullingFrequency, "pulling-frequency", 30*time.Second, "pulling frequency from mastodon")

	flag.Parse()

	loggerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOptions))

	app := &application{
		config:    cfg,
		logger:    logger,
		seenPosts: make(map[string]struct{}),
	}

	err := app.ChecKDependentService()
	if err != nil {
		app.logger.Error("api healthcheck error", "error", err.Error())
		app.logger.Error("unable to reach dependent service, exiting...")
		os.Exit(1)
	}

	// Run mastodon puller in background
	go app.mastodonBackgroundJob()

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
