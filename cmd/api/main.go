package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

import (
	"github.com/andrewah64/base-app-client/internal/api/core/route"
	"github.com/andrewah64/base-app-client/internal/api/core/ui/i18n"
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/log"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/tenant"
)

import (
	"github.com/andrewah64/base-app-client/cmd/api/core/auth/users/register"
	"github.com/andrewah64/base-app-client/cmd/api/core/unauth/health"
)

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

const version = "1.0.0"

func main(){
	/* capture values passed via flags */
	port        := flag.Int   ("port"        , 8080                        , "Port")
	loglvl      := flag.String("loglvl"      , "info"                      , "Level of default logger (debug|info|error)")
	pghost      := flag.String("pghost"      , "localhost"                 , "Host of PostgreSQL")
	pgport      := flag.Int   ("pgport"      , 5432                        , "Port of PostgreSQL")
	pguser      := flag.String("pguser"      , "postgres"                  , "Name of PostgreSQL user")
	pgpw        := flag.String("pgpw"        , "secret"                    , "Password for 'pguser'")
	pgdb        := flag.String("pgdb"        , "gopgtest"                  , "Database name")
	pgsslmode   := flag.String("pgsslmode"   , "disable"                   , "Secure connections to PG with SSL (enable|disable")
	pgcachesize := flag.Int   ("pgcachesize" , 0                           , "Size of the PG statement cache")
	pgapp       := flag.String("pgapp"       , "myapp"                     , "Name of the application")

	flag.Parse()

	/* set up a request Id so we can track startup activity */
	ctx := session.NewContext(context.Background(), &session.CtxData{
		RequestId: uuid.NewString(),
	})

	/* set up the default logger. this will log things that happen at startup and all errors*/
	llvl := slog.LevelInfo

	switch *loglvl {
		case "debug": llvl = slog.LevelDebug
		case "info" : llvl = slog.LevelInfo
		case "error": llvl = slog.LevelError
		default     : panic(fmt.Sprintf("'loglvl' can be (debug|info|error). '%v' is an invalid choice", *loglvl))
	}

	slog.SetDefault(log.Setup(slog.Level(llvl)))

	slog.LogAttrs(ctx, slog.LevelInfo, "start")

	/* get  & validate a connection pool to the postgres DB */
	pool, cpErr := db.ConnPool(&ctx, slog.Default(), pghost, pgport, pgdb, pguser, pgpw, pgsslmode, pgcachesize, pgapp)
	if cpErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "get pool",
			slog.String("error", cpErr.Error()),
		)

		panic(cpErr)
	}

	if pingErr := pool.Ping(ctx); pingErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "ping pool",
			slog.String("error", pingErr.Error()),
		)

		panic(pingErr)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "acquired connection pool")

	defer pool.Close()

	slog.LogAttrs(ctx, slog.LevelInfo, "deferred close of connection pool")

	/* populate our i18n bundle cache */
	i18nCacheErr := i18n.InitCache(ctx, language.English)
	if i18nCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the bundle cache",
			slog.String("error", i18nCacheErr.Error()),
		)

		panic(i18nCacheErr)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "initialised bundle cache")

	/* populate our tenant and route caches*/
	conn, connErr := db.Conn(&ctx, slog.Default(), pool)
	if connErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the tenant cache",
			slog.String("error", connErr.Error()),
		)

		panic(connErr)
	}

	/* tenant cache */
	idErr := session.Identity(&ctx, slog.Default(), conn, "role_all_core_unauth_tnt_all_inf")
	if idErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the tenant cache",
			slog.String("error", idErr.Error()),
		)

		panic(connErr)
	}

	tntCacheErr := tenant.InitCache(&ctx, conn)
	if tntCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the tenant cache",
			slog.String("error", tntCacheErr.Error()),
		)

		defer conn.Release()

		panic(tntCacheErr)
	}

	/*route cache*/
	rtsIdErr := session.Identity(&ctx, slog.Default(), conn, "role_api_core_rts_api_inf")
	if rtsIdErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the route cache",
			slog.String("error", rtsIdErr.Error()),
		)

		panic(connErr)
	}

	rtsCacheErr := route.InitCache(&ctx, conn)
	if rtsCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the route cache",
			slog.String("error", rtsCacheErr.Error()),
		)

		defer conn.Release()

		panic(rtsCacheErr)
	}

	conn.Release()

	/* TLS preferences */
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion      : tls.VersionTLS13,
	}

	/* start the server */
	var (
		handlers = map[string]http.HandlerFunc{
			"api.core.auth.aur.tnt.reg.Register" : register.Register,
			"api.core.unauth.health.Check"          : health.Check,
		}
	)

	server := &http.Server{
		Addr        :	fmt.Sprintf(":%d", *port),
		Handler     :	route.Mux(&ctx, handlers),
		BaseContext :	func(_ net.Listener) context.Context {
					return db.NewContext(
						context.Background(),
						&db.Pool {
							Pool: pool,
						},
					)
				},
		TLSConfig   :	tlsConfig,
		IdleTimeout :	time.Minute,
		ReadTimeout :	5  * time.Second,
		WriteTimeout:	10 * time.Second,
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "start server",
		slog.Int("port", *port),
	)

	/* if the server falls over capture errors from it & quit  */
	srvErr := server.ListenAndServeTLS("cert.pem", "key.pem")

	slog.LogAttrs(ctx, slog.LevelError, "server error",
		slog.String("error", srvErr.Error()),
	)

	os.Exit(1)
}
