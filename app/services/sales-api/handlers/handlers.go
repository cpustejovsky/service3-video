//Package handlers contains routes
package handlers

import (
	"expvar"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/cpustejovsky/service3-video/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/cpustejovsky/service3-video/app/services/sales-api/handlers/v1/testgrp"
	"github.com/cpustejovsky/service3-video/business/sys/auth"
	"github.com/cpustejovsky/service3-video/business/web/mid"
	"github.com/cpustejovsky/service3-video/foundation/web"
	"go.uber.org/zap"
)

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
		DB:    db,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	DB       *sqlx.DB
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		//Panics should always be the closest onion ring
		mid.Panics(),
	)

	v1(app, cfg)

	return app
}

func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"
	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, mid.Authenticate(cfg.Auth), mid.Authorize("ADMIN"))
}
