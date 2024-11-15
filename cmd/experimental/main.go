package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/larwef/yggdrasil/internal/server"
)

// Variables injected at compile time.
var (
	appName = "No name provided"
	version = "No version provided"
)

type Config struct {
	LogLvl    slog.Level `envconfig:"LOG_LEVEL" default:"info"`
	LogSource bool       `envconfig:"LOG_SOURCE"`
	LogJSON   bool       `envconfig:"LOG_JSON" default:"true"`

	ServerConfig server.Config
}

func main() {
	var conf Config
	if err := envconfig.Process("EXPERIMENTAL", &conf); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Set the default slog logger in stead of passing it around.
	logHandlerOpts := &slog.HandlerOptions{
		Level:     conf.LogLvl,
		AddSource: conf.LogSource,
	}

	var logHandler slog.Handler
	if conf.LogJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, logHandlerOpts)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, logHandlerOpts)
	}
	logger := slog.
		New(logHandler).
		With(slog.Group("application", "name", appName, "version", version))

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := realMain(ctx, conf, logger)
	done()
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("program finished with error", "error", err)
	} else {
		logger.Info("program finished gracefully")
	}
}

func realMain(ctx context.Context, conf Config, logger *slog.Logger) error {
	logger.Info("starting application")
	srv := server.New(logger.With(slog.String("component", "server")), conf.ServerConfig, routes())
	return srv.ListenAndServeContext(ctx)
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorld())
	return mux
}

// Just a simple example of a handler.
func helloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello world")); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
