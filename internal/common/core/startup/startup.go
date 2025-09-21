package startup

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/log"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/tenant"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type RuntimeParams struct {
	HttpPort    *int
	LogLvl      *string
	PgHost      *string
	PgPort      *int
	PgUser      *string
	PgPw        *string
	PgDb        *string
	PgSslMode   *string
	PgCacheSize *int
	PgApp       *string
}

func GetRuntimeParams () *RuntimeParams {
	httpPort    := flag.Int   ("port"        , 8081        , "Port")
	logLvl      := flag.String("loglvl"      , "info"      , "Level of default logger (debug|info|error)")
	pgHost      := flag.String("pghost"      , "localhost" , "Host of PostgreSQL")
	pgPort      := flag.Int   ("pgport"      , 5432        , "Port of PostgreSQL")
	pgUser      := flag.String("pguser"      , "postgres"  , "Name of PostgreSQL user")
	pgPw        := flag.String("pgpw"        , "secret"    , "Password for 'pguser'")
	pgDb        := flag.String("pgdb"        , "gopgtest"  , "Database name")
	pgSslMode   := flag.String("pgsslmode"   , "disable"   , "Secure connections to PG with SSL (enable|disable")
	pgCacheSize := flag.Int   ("pgcachesize" , 0           , "Size of the PG statement cache")
	pgApp       := flag.String("pgapp"       , "myapp"     , "Name of the application")

	p := &RuntimeParams {
		HttpPort    : httpPort,
		LogLvl      : logLvl,
		PgHost      : pgHost,
		PgPort      : pgPort,
		PgUser      : pgUser,
		PgPw        : pgPw,
		PgDb        : pgDb,
		PgSslMode   : pgSslMode,
		PgCacheSize : pgCacheSize,
		PgApp       : pgApp,
	}

	flag.Parse()

	return p
}

func SetupDefaultLogger (logLvl string) {
	lvl := slog.LevelInfo

	switch logLvl {
		case "debug" :
			lvl = slog.LevelDebug
		case "info"  :
			lvl = slog.LevelInfo
		case "error" :
			lvl = slog.LevelError
		default      :
			panic(fmt.Sprintf("'loglvl' can be (debug|info|error). '%v' is an invalid choice", logLvl))
	}

	slog.SetDefault(log.Setup(slog.Level(lvl)))
}

func SetupPGConnectionPool (ctx context.Context, rtp *RuntimeParams) (*pgxpool.Pool) {
	pool, cpErr := db.ConnPool(&ctx, slog.Default(), rtp.PgHost, rtp.PgPort, rtp.PgDb, rtp.PgUser, rtp.PgPw, rtp.PgSslMode, rtp.PgCacheSize, rtp.PgApp)
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

	return pool
}

func SetupTenantCache (ctx context.Context, conn *pgxpool.Conn) {
	idErr := session.Identity(&ctx, slog.Default(), conn, "role_all_core_unauth_tnt_all_inf")
	if idErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the tenant cache",
			slog.String("error", idErr.Error()),
		)

		panic(idErr)
	}

	tntCacheErr := tenant.InitCache(&ctx, conn)
	if tntCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the tenant cache",
			slog.String("error", tntCacheErr.Error()),
		)

		panic(tntCacheErr)
	}
}
