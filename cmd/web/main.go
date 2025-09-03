package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"net"
	"net/http"
	"time"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/log"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/tenant"
	"github.com/andrewah64/base-app-client/internal/web/core/passkey"
	"github.com/andrewah64/base-app-client/internal/web/core/route"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/i18n"
)

import (
	authhome         "github.com/andrewah64/base-app-client/cmd/web/core/auth/home"
	authatntnt       "github.com/andrewah64/base-app-client/cmd/web/core/auth/atn/tnt"
	authaurgrptnt    "github.com/andrewah64/base-app-client/cmd/web/core/auth/aur/grp/tnt"
	authgrpaurtnt    "github.com/andrewah64/base-app-client/cmd/web/core/auth/grp/aur/tnt"
	authgrptnt       "github.com/andrewah64/base-app-client/cmd/web/core/auth/grp/tnt"
	authgrptntid     "github.com/andrewah64/base-app-client/cmd/web/core/auth/grp/tnt/id"
	authgrptntval    "github.com/andrewah64/base-app-client/cmd/web/core/auth/grp/tnt/val"
	authkeyaur       "github.com/andrewah64/base-app-client/cmd/web/core/auth/key/aur"
	authkeyaurid     "github.com/andrewah64/base-app-client/cmd/web/core/auth/key/aur/id"
	authkeyaurval    "github.com/andrewah64/base-app-client/cmd/web/core/auth/key/aur/val"
	authlogaurtnt    "github.com/andrewah64/base-app-client/cmd/web/core/auth/log/aur/tnt"
	authlogaurtntid  "github.com/andrewah64/base-app-client/cmd/web/core/auth/log/aur/tnt/id"
	authlogeptnt     "github.com/andrewah64/base-app-client/cmd/web/core/auth/log/ep/tnt"
	authlogeptntid   "github.com/andrewah64/base-app-client/cmd/web/core/auth/log/ep/tnt/id"
	authpwdaurtnt    "github.com/andrewah64/base-app-client/cmd/web/core/auth/pwd/aur/tnt"
	authpwdaurtntval "github.com/andrewah64/base-app-client/cmd/web/core/auth/pwd/aur/tnt/val"
	authrolgrptnt    "github.com/andrewah64/base-app-client/cmd/web/core/auth/rol/grp/tnt"
	authrolkeyaur    "github.com/andrewah64/base-app-client/cmd/web/core/auth/rol/key/aur"
	authssntnt       "github.com/andrewah64/base-app-client/cmd/web/core/auth/ssn/tnt"
	authssnaur       "github.com/andrewah64/base-app-client/cmd/web/core/auth/ssn/aur"
	authaurtnt       "github.com/andrewah64/base-app-client/cmd/web/core/auth/aur/tnt"
	authaurtntid     "github.com/andrewah64/base-app-client/cmd/web/core/auth/aur/tnt/id"
	authaurtntval    "github.com/andrewah64/base-app-client/cmd/web/core/auth/aur/tnt/val"
	oidc             "github.com/andrewah64/base-app-client/cmd/web/core/oidc"
	unauthaurtnt     "github.com/andrewah64/base-app-client/cmd/web/core/unauth/aur/tnt"
	unauthaurtntval  "github.com/andrewah64/base-app-client/cmd/web/core/unauth/aur/tnt/val"
	unauthotpaur     "github.com/andrewah64/base-app-client/cmd/web/core/unauth/otp/aur"
	unauthotpssnaur  "github.com/andrewah64/base-app-client/cmd/web/core/unauth/otp/ssn/aur"
	unauthssnaur     "github.com/andrewah64/base-app-client/cmd/web/core/unauth/ssn/aur"
)

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

func main() {
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

	/* populate our template cache */
	html.InitCache(ctx)

	slog.LogAttrs(ctx, slog.LevelInfo, "initialised template cache")

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

	/* tenant and passkey cache */
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

	pkeyCacheErr := passkey.InitCache(&ctx, conn)
	if pkeyCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the passkey cache",
			slog.String("error", pkeyCacheErr.Error()),
		)

		defer conn.Release()

		panic(pkeyCacheErr)
	}

	/*route cache*/
	rtsIdErr := session.Identity(&ctx, slog.Default(), conn, "role_web_core_unauth_rts_web_inf")
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

	/*end db-related initialisation*/
	conn.Release()

	/* TLS preferences */
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion      : tls.VersionTLS13,
	}

	/* start the server */

	var (
		handlers = map[string]http.HandlerFunc{
			"web.core.auth.atn.tnt.Get"         : authatntnt.Get,
			"web.core.auth.atn.tnt.Patch"       : authatntnt.Patch,
			"web.core.auth.aur.grp.tnt.Get"     : authaurgrptnt.Get,
			"web.core.auth.aur.grp.tnt.Patch"   : authaurgrptnt.Patch,
			"web.core.auth.aur.tnt.Delete"      : authaurtnt.Delete,
			"web.core.auth.aur.tnt.Get"         : authaurtnt.Get,
			"web.core.auth.aur.tnt.Post"        : authaurtnt.Post,
			"web.core.auth.aur.tnt.id.Get"      : authaurtntid.Get,
			"web.core.auth.aur.tnt.id.Patch"    : authaurtntid.Patch,
			"web.core.auth.aur.tnt.val.Get"     : authaurtntval.Get,
			"web.core.auth.grp.aur.tnt.Get"     : authgrpaurtnt.Get,
			"web.core.auth.grp.aur.tnt.Patch"   : authgrpaurtnt.Patch,
			"web.core.auth.grp.tnt.Delete"      : authgrptnt.Delete,
			"web.core.auth.grp.tnt.Get"         : authgrptnt.Get,
			"web.core.auth.grp.tnt.Post"        : authgrptnt.Post,
			"web.core.auth.grp.tnt.id.Get"      : authgrptntid.Get,
			"web.core.auth.grp.tnt.id.Patch"    : authgrptntid.Patch,
			"web.core.auth.grp.tnt.val.Get"     : authgrptntval.Get,
			"web.core.auth.home.Index"          : authhome.Get,
			"web.core.auth.log.aur.tnt.Get"     : authlogaurtnt.Get,
			"web.core.auth.log.aur.tnt.Put"     : authlogaurtnt.Put,
			"web.core.auth.log.aur.tnt.id.Get"  : authlogaurtntid.Get,
			"web.core.auth.log.aur.tnt.id.Patch": authlogaurtntid.Patch,
			"web.core.auth.log.ep.tnt.Get"      : authlogeptnt.Get,
			"web.core.auth.log.ep.tnt.Put"      : authlogeptnt.Put,
			"web.core.auth.log.ep.tnt.id.Get"   : authlogeptntid.Get,
			"web.core.auth.log.ep.tnt.id.Patch" : authlogeptntid.Patch,
			"web.core.auth.key.aur.Get"         : authkeyaur.Get,
			"web.core.auth.key.aur.Post"        : authkeyaur.Post,
			"web.core.auth.key.aur.Delete"      : authkeyaur.Delete,
			"web.core.auth.key.aur.id.Get"      : authkeyaurid.Get,
			"web.core.auth.key.aur.id.Patch"    : authkeyaurid.Patch,
			"web.core.auth.key.aur.val.Get"     : authkeyaurval.Get,
			"web.core.auth.pwd.aur.tnt.Get"     : authpwdaurtnt.Get,
			"web.core.auth.pwd.aur.tnt.Patch"   : authpwdaurtnt.Patch,
			"web.core.auth.pwd.aur.tnt.val.Get" : authpwdaurtntval.Get,
			"web.core.auth.rol.grp.tnt.Get"     : authrolgrptnt.Get,
			"web.core.auth.rol.grp.tnt.Patch"   : authrolgrptnt.Patch,
			"web.core.auth.rol.key.aur.Get"     : authrolkeyaur.Get,
			"web.core.auth.rol.key.aur.Patch"   : authrolkeyaur.Patch,
			"web.core.auth.ssn.tnt.Get"         : authssntnt.Get,
			"web.core.auth.ssn.tnt.Delete"      : authssntnt.Delete,
			"web.core.auth.ssn.aur.Delete"      : authssnaur.Delete,
			"web.core.oidc.Call"                : oidc.Call,
			"web.core.oidc.Callback"            : oidc.Callback,
			"web.core.unauth.aur.tnt.Get"       : unauthaurtnt.Get,
			"web.core.unauth.aur.tnt.Post"      : unauthaurtnt.Post,
			"web.core.unauth.aur.tnt.val.Get"   : unauthaurtntval.Get,
			"web.core.unauth.otp.aur.Get"       : unauthotpaur.Get,
			"web.core.unauth.otp.aur.Post"      : unauthotpaur.Post,
			"web.core.unauth.otp.ssn.aur.Get"   : unauthotpssnaur.Get,
			"web.core.unauth.otp.ssn.aur.Post"  : unauthotpssnaur.Post,
			"web.core.unauth.ssn.aur.Get"       : unauthssnaur.Get,
			"web.core.unauth.ssn.aur.Post"      : unauthssnaur.Post,
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
