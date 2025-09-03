package session

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/token"
)

func Begin(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, rw http.ResponseWriter, userId int, expiry time.Time) error {
	ssnTkn, stErr := token.Token(32)
	if (stErr != nil) {
		return stErr
	}

	http.SetCookie(rw, &http.Cookie{
		Name    : "session_token",
		Value   : ssnTkn,
		Expires : expiry,
		HttpOnly: true,
		Secure  : true,
		Path    : "/",
		SameSite: http.SameSiteLaxMode,
	})

	const (
		dbSchema = "web_core_unauth_ssn_aur_reg"
		dbSproc  = "reg_ssn"
	)

	var (
		sprocCall   = fmt.Sprintf("call %v.%v(@p_aur_id, @p_wauhs_ssn_tk, @p_wauhs_exp_ts)", dbSchema, dbSproc)
		sprocParams = pgx.NamedArgs{
			"p_aur_id"       : userId,
			"p_wauhs_ssn_tk" : ssnTkn,
			"p_wauhs_exp_ts" : expiry,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, nil)
	if sprocErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "call sproc",
			slog.String("error"       , sprocErr.Error()),
			slog.String("sprocCall"   , sprocCall),
			slog.Int   ("userid"      , userId),
			slog.String("sessionToken", ssnTkn),
			slog.Time  ("cookieExpiry", expiry),
		)
		return sprocErr
	}

	return nil
}

func End(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, rw http.ResponseWriter, ssnTkn *http.Cookie) error {
	expiry := time.Now()

	http.SetCookie(rw, &http.Cookie{
		Name    : "session_token",
		Value   : "",
		Expires : expiry,
		HttpOnly: true,
		Secure  : true,
		Path    : "/",
		SameSite: http.SameSiteLaxMode,
	})

	const (
		dbSchema = "web_core_auth_ssn_aur_end"
		dbSproc  = "end_ssn"
	)

	var (
		sprocCall   = fmt.Sprintf("call %v.%v(@p_wauhs_ssn_tk)", dbSchema, dbSproc)
		sprocParams = pgx.NamedArgs{
			"p_wauhs_ssn_tk" : ssnTkn.Value,
		}
	)

	err := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, nil)
	if err != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "call sproc",
			slog.String("sprocCall", sprocCall),
			slog.String("ssnTkn"   , ssnTkn.Value),
			slog.Time  ("expiry"   , expiry),
		)
		return err
	}

	return nil
}
