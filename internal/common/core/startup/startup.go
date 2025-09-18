package startup

import (
	"context"
	"crypto/x509"
	"flag"
	"fmt"
	"log/slog"
	"time"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/log"
	"github.com/andrewah64/base-app-client/internal/common/core/saml2"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/tenant"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GenSAML2ServiceProviderCerts (ctx *context.Context, conn *pgxpool.Conn) (error){
	s2cErr := session.Identity(ctx, slog.Default(), conn, "role_all_core_unauth_spc_all_reg")
	if s2cErr != nil {
		panic(s2cErr)
	}

	const (
		dbSchema = "all_core_unauth_spc_all_reg"
		dbFunc   = "s2g_inf"
	)

	type s2g struct {
		TntId       int
		S2gCrtCnNm  string
		S2gCrtDn    time.Duration
		S2gCrtOrgNm string
	}

	s2gRs, s2gRsErr := db.DataSet[s2g](ctx, slog.Default(), conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "get tenant data",
				slog.String("error", err.Error()),
			)

			return qry, dbFunc, nil, fmt.Errorf("call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	if s2gRsErr != nil {
		return s2gRsErr
	}

	sprocCall := fmt.Sprintf("call %v.reg_spc(@p_tnt_id, @p_spc_cn_nm, @p_spc_org_nm, @p_spc_enc_crt, @p_spc_enc_pvk, @p_spc_sgn_crt, @p_spc_sgn_pvk, @p_spc_exp_ts, @p_spc_enabled)", dbSchema)

	encKeyUsage := x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment
	sgnKeyUsage := x509.KeyUsageDigitalSignature
	now         := time.Now()
	spcFromTs   := now.Add(-time.Hour)

	for _, v := range s2gRs {
		spcExpTs := now.Add(v.S2gCrtDn);

		spcEncCrt, spcEncPvk, spcEncCrtErr := saml2.GenCert(v.S2gCrtCnNm, []string{v.S2gCrtOrgNm}, encKeyUsage, spcFromTs, spcExpTs)
		if spcEncCrtErr != nil {
			panic(spcEncCrtErr)
		}

		spcSgnCrt, spcSgnPvk, spcSgnCrtErr := saml2.GenCert(v.S2gCrtCnNm, []string{v.S2gCrtOrgNm}, sgnKeyUsage, spcFromTs, spcExpTs)
		if spcSgnCrtErr != nil {
			panic(spcSgnCrtErr)
		}

		var (
			sprocParams = pgx.NamedArgs{
				"p_tnt_id"      : v.TntId,
				"p_spc_cn_nm"   : v.S2gCrtCnNm,
				"p_spc_org_nm"  : v.S2gCrtOrgNm,
				"p_spc_enc_crt" : spcEncCrt,
				"p_spc_enc_pvk" : spcEncPvk,
				"p_spc_sgn_crt" : spcSgnCrt,
				"p_spc_sgn_pvk" : spcSgnPvk,
				"p_spc_exp_ts"  : spcExpTs,
				"p_spc_enabled" : true,
			}
		)

		sprocErr := db.Sproc(ctx, slog.Default(), conn, sprocCall, sprocParams, nil)
		if sprocErr != nil {
			panic(sprocErr)
		}
	}

	return nil
}

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
