package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/go-ozzo/ozzo-dbx"
	//- "github.com/go-ozzo/ozzo-routing/v2"
	//- "github.com/go-ozzo/ozzo-routing/v2/content"
	//- "github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/julienschmidt/httprouter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tvitcom/fusion-framework/internal/album"
	// "github.com/tvitcom/fusion-framework/internal/auth"
	"github.com/tvitcom/fusion-framework/internal/config"
	// "github.com/tvitcom/fusion-framework/internal/errors"
	// "github.com/tvitcom/fusion-framework/internal/healthcheck"
	// "github.com/tvitcom/fusion-framework/pkg/accesslog"
	"github.com/tvitcom/fusion-framework/pkg/dbcontext"
	"github.com/tvitcom/fusion-framework/pkg/log"
	"net/http"
	"fmt"
	"os"
	"time"
)

// Version indicates the current version of the application.
var Version = "1.0.0"
var defaultConfigFile = "./configs/dev.yml"
var configFile = flag.String("config", defaultConfigFile, "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

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

	// build HTTP server
	logger.Infof("server %v is running at %v", Version, cfg.HttpEntrypoint)
	
	r := httprouter.New()

	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Access-Control-Request-Method") != "" {
	    // Set CORS headers
	    header := w.Header()
	    header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
	    header.Set("Access-Control-Allow-Origin", "*")
	}

	// Adjust status code to 204
	w.WriteHeader(http.StatusNoContent)
	})

	// start server on https port
	// hs := http.Server{
	// 	Handler: registerRoutes(logger, dbcontext.New(db), cfg),
	// 	TLSConfig: &tls.Config{
	// 		NextProtos: []string{"h2", "http/1.1"},
	// 	},
	// }
	// if err := hs.ListenAndServe(cfg.HttpEntrypoint, nil); err != nil && err != http.ErrServerClosed {
	// 	logger.Error(err)
	// 	os.Exit(-1)
	// }
	if err := http.ListenAndServe(cfg.HttpEntrypoint, registerRoutes(r, logger, dbcontext.New(db), cfg)); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

// registerRoutes sets up the HTTP routing and builds an HTTP handler.
func registerRoutes(router *httprouter.Router, logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {

	router.GET("/", getIndex)
	
	// router.Use(
	// 	accesslog.Handler(logger),
	// 	errors.Handler(logger),
	// 	content.TypeNegotiator(content.JSON),
	// 	cors.Handler(cors.AllowAll),
	// )

	// healthcheck.RegisterHandlers(router, Version)

	// rg := router.Group("/v1")

	// authHandler := auth.Handler(cfg.JWTSigningKey)

	// album.RegisterHandlers(rg.Group(""),
	// 	album.NewService(album.NewRepository(db, logger), logger),
	// 	authHandler, 
	// 	logger,
	// )

	album.RegisterHandlers(
		router, // Router
		album.NewAgregator(album.NewRepository(db, logger), logger), // Agregator
		//authHandler, // Auth handler: JWT-based authentication middleware
		logger, // Logger
	)

	// auth.RegisterHandlers(rg.Group(""),
	// 	auth.NewAgregator(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
	// 	logger,
	// )

	return router
}

func getIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Welcome to Fusion index route!\n")
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
