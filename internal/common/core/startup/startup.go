package startup

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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

import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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
	PgCred      *string
	AwsProfile  *string
	AwsSecretNm *string
}

func GetRuntimeParams () *RuntimeParams {
	httpPort    := flag.Int   ("port"        , 8081        , "Port")
	logLvl      := flag.String("loglvl"      , "info"      , "Level of default logger (debug|info|error)")
	pgHost      := flag.String("pghost"      , "localhost" , "Host of PostgreSQL")
	pgPort      := flag.Int   ("pgport"      , 5432        , "Port of PostgreSQL")
	pgUser      := flag.String("pguser"      , "postgres"  , "Name of PostgreSQL user")
	pgPw        := flag.String("pgpw"        , ""          , "Password for 'pguser'")
	pgDb        := flag.String("pgdb"        , "base-app"  , "Database name")
	pgSslMode   := flag.String("pgsslmode"   , "disable"   , "Secure connections to PG with SSL (enable|disable")
	pgCacheSize := flag.Int   ("pgcachesize" , 0           , "Size of the PG statement cache")
	pgApp       := flag.String("pgapp"       , "myapp"     , "Name of the application")
	pgCred      := flag.String("pgcred"      , "systemd"   , "PostgreSQL password retrieval method (systemd)")
	awsProfile  := flag.String("awsprofile"  , ""          , "AWS profile used to retrieve pgpw from secret's manager")
	awsSecretNm := flag.String("awssecretnm" , ""          , "Name of AWS secret")

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
		PgCred      : pgCred,
		AwsProfile  : awsProfile,
		AwsSecretNm : awsSecretNm,
	}

	flag.Parse()

	provided := make(map[string]bool)

	flag.Visit(func(f *flag.Flag){
		provided[f.Name] = true
	})

	if ! provided["pgcred"] {
		panic("pgcred must be supplied and can be (password-plain|password-systemd)")
	} else {
		passwordPlain             := "password-plain"
		passwordSystemd           := "password-systemd"
		passwordAWSSecretsManager := "password-aws-secrets-manager"

		if ! ( *p.PgCred == passwordPlain || *p.PgCred == passwordSystemd || *p.PgCred == passwordAWSSecretsManager) {
			panic("pgcred must be supplied and can be (password-plain|password-systemd|password-aws-secrets-manager)")
		}

		if *p.PgCred == passwordPlain && ! provided["pgpw"] {
			panic("pgpw must be supplied")
		}

		if *p.PgCred == passwordAWSSecretsManager && ! provided["awsprofile"] && ! provided["awssecretnm"] {
			panic("awsprofile and awssecretnm must be provided")
		}

		switch *p.PgCred {
			case passwordSystemd:
				if provided["pgpw"] {
					panic("pgpw must not be supplied when pgcred is password-systemd")
				}

				credPath := os.Getenv("CREDENTIALS_DIRECTORY")

				_, credPathErr := os.Open(credPath)
				if credPathErr != nil {
					panic("Failed to locate the systemd credentials folder")
				}

				pgPw, pgPwErr := os.ReadFile(filepath.Join(credPath,"postgres-password"))
				if pgPwErr != nil {
					panic("Failed to retrieve postgres password from systemd")
				}

				*p.PgPw = strings.TrimSpace(string(pgPw))
			case passwordAWSSecretsManager:
				if provided["pgpw"] {
					panic("pgpw must not be supplied when pgcred is password-aws-secrets-manager")
				}

				cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(*awsProfile))
				if err != nil {
					panic("Failed to retrieve profile from the aws config file")
				}

				input := &secretsmanager.GetSecretValueInput{
					SecretId:     aws.String(*awsSecretNm),
					VersionStage: aws.String("AWSCURRENT"),
				}

				pgpw, err := secretsmanager.NewFromConfig(cfg).GetSecretValue(context.TODO(), input)
				if err != nil {
					panic("Failed to retrieve password from AWS Secrets Manager")
				}

				*p.PgPw = *pgpw.SecretString
		}
	}

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
