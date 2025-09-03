package otp

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type AurInf struct {
	AurSsnDn  time.Duration
	EppPt     string
	OtpSecret string
}

func GetAurInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, nncNonce string) ([]AurInf, error) {
	const (
		dbSchema = "web_core_unauth_otp_ssn_aur_mod"
		dbFunc   = "aur_inf"
	)

	results, err := db.DataSet[AurInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, nncNonce)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetAurInf::get dataset",
				slog.String("error"    , err.Error()),
				slog.String("qry"      , qry),
				slog.Int   ("tntId"    , tntId),
				slog.Int   ("aurId"    , aurId),
				slog.String("nncNonce" , nncNonce),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetAurInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type NncInf struct {
	AurId      int
	AurNm      string
	NncEnabled bool
	OtpEnabled bool
}

func GetNncInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, nncNonce string) ([]NncInf, error) {
	const (
		dbSchema = "web_core_unauth_otp_ssn_aur_mod"
		dbFunc   = "nnc_inf"
	)

	results, err := db.DataSet[NncInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, nncNonce)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetNncInf::get dataset",
				slog.String("error"    , err.Error()),
				slog.String("qry"      , qry),
				slog.Int   ("tntId"    , tntId),
				slog.String("nncNonce" , nncNonce),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetNncInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

func PostOtp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, otpId string, otpSecret string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_otp_ssn_aur_mod.reg_otp(@p_tnt_id, @p_aur_id, @p_otp_id, @p_otp_secret)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_aur_id"     : aurId,
			"p_otp_id"     : otpId,
			"p_otp_secret" : otpSecret,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("aurId"     , aurId),
			slog.String("otpId"     , otpId),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
