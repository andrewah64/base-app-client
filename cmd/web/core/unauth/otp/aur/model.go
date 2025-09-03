package otp

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type OtpInf struct {
	OtpSecret string
}

func GetOtpInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, otpId string) ([]OtpInf, error) {
	const (
		dbSchema = "web_core_unauth_otp_aur_mod"
		dbFunc   = "otp_inf"
	)

	results, err := db.DataSet[OtpInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, otpId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetOtpInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
				slog.Int   ("aurId" , aurId),
				slog.String("otpId" , otpId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetOtpInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

type OtpAurInf struct {
	AurId     int
	AurNm     string
	OtpSecret string
}

func GetOtpAurInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, otpId string) ([]OtpAurInf, error) {
	const (
		dbSchema = "web_core_unauth_otp_aur_inf"
		dbFunc   = "otp_aur_inf"
	)

	results, err := db.DataSet[OtpAurInf](ctx, logger, conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc, tntId, otpId)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "GetOtpAurInf::get dataset",
				slog.String("error" , err.Error()),
				slog.String("qry"   , qry),
				slog.Int   ("tntId" , tntId),
				slog.String("otpId" , otpId),
			)

			return qry, dbFunc, nil, fmt.Errorf("GetOtpAurInf::call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	return results, err
}

func PostOtp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, otpId string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_otp_aur_mod.mod_otp(@p_tnt_id, @p_aur_id, @p_otp_id)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_aur_id"     : aurId,
			"p_otp_id"     : otpId,
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
