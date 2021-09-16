package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/gofiber/fiber/v2"
	// "github.com/go-ozzo/ozzo-routing/v2/content"
	// "github.com/go-ozzo/ozzo-routing/v2/cors"
	
	_ "github.com/go-sql-driver/mysql"
	// "github.com/tvitcom/fusion-framework/internal/album"
	// "github.com/tvitcom/fusion-framework/internal/auth"
	"github.com/tvitcom/fusion-framework/internal/config"
	// "github.com/tvitcom/fusion-framework/internal/errors"
	"github.com/tvitcom/fusion-framework/internal/healthcheck"
	// "github.com/tvitcom/fusion-framework/pkg/accesslog"
	"github.com/tvitcom/fusion-framework/pkg/dbcontext"
	logz "github.com/tvitcom/fusion-framework/pkg/log"
	"os"
	"time"
	"os/signal"
	"syscall"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var configFile = flag.String("config", "./configs/dev.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := logz.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*configFile, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	db, err := dbx.MustOpen(cfg.DBType, cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	router := fiber.New(fiber.Config{
		BodyLimit: (4 * 1024 * 1024),
		ReadTimeout: 6 * time.Second,
		WriteTimeout: 6 * time.Second,
        Prefork:       (cfg.AppMode == "prod"),
        CaseSensitive: false,
        StrictRouting: true,
        ServerHeader:  "fusion-server",
    })

	buildHandler(router, logger, dbcontext.New(db), cfg)

	go func() {
		if err := router.Listen(cfg.HttpEntrypoint); err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}()
	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	_ = <-c // This blocks the main thread until an interrupt is received

	_ = router.Shutdown()
	
	// Your cleanup tasks go here
	if err := db.Close(); err != nil {
		logger.Error(err)
	}
	logger.Infof("server %v was successful shutdown.")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
// func buildHandler(app *fiber.App, logger logz.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
func buildHandler(router *fiber.App, logger logz.Logger, db *dbcontext.DB, cfg *config.Config) {

	// router := routing.New()

	// router.Use(
	// 	accesslog.Handler(logger),
	// 	errors.Handler(logger),
	// 	content.TypeNegotiator(content.JSON),
	// 	cors.Handler(cors.AllowAll),
	// )

	router.Static("/assets", "./web/assets", fiber.Static{
	  Compress:      true,
	  ByteRange:     true,
	  Browse:        true,
	  Index:         "index.html",
	  CacheDuration: 120 * time.Minute,
	  MaxAge:        3600,
	})

	healthcheck.RegisterHandlers(router, Version)

	// rg := router.Group("/v1")

	// authHandler := auth.Handler(cfg.JWTSigningKey)

	// album.RegisterHandlers(rg.Group(""),
	// 	album.NewAgregator(album.NewRepository(db, logger), logger),
	// 	authHandler, logger,
	// )

	// auth.RegisterHandlers(rg.Group(""),
	// 	auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
	// 	logger,
	// )

	// return router

// api := r.Group("/api", handlerMW1)  // /api
// v1 := api.Group("/v1", handlerMW2)   // /api/v1
// v1.Get("/list", handler3)          // /api/v1/list
// v1.Get("/user", handler4)          // /api/v1/user

// v2 := api.Group("/v2", handlerMW3)   // /api/v2
// v2.Get("/list", handlerList)          // /api/v2/list
// v2.Get("/user", handlerUser)          // /api/v2/user
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger logz.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger logz.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
