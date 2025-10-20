package id

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

const (
	dbSchema = "web_core_auth_s2c_tnt_mod"
)

type RowIdpInf struct {
	IdpId       int
	IdpNm       string
	IdpEntityId string
	IdpEnabled  bool
	NumMde      int
	NumSso      int
	NumSlo      int
}

func GetRowIdpInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpId int) ([]RowIdpInf, error) {
	const (
		dbFunc = "row_idp_inf"
	)

	rs, rErr := db.DataSet[RowIdpInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, idpId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("idpId" , idpId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type RowIdpMod struct {
	IdpId       int
	IdpNm       string
	IdpEntityId string
	IdpEnabled  bool
	NumMde      int
	NumSso      int
	NumSlo      int
	Uts         time.Time
}

func GetRowIdpMod (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpId int) ([]RowIdpMod, error) {
	const (
		dbFunc = "row_idp_mod"
	)

	rs, rErr := db.DataSet[RowIdpMod](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, idpId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("idpId" , idpId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type RowIdpVal struct {
	IdpEnabledOk bool
	IdpNmOk      bool
}

func GetRowIdpVal (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpId int, idpEnabled bool, idpNm string) ([]RowIdpVal, error) {
	const (
		dbFunc = "row_idp_val"
	)

	rs, rErr := db.DataSet[RowIdpVal](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4, $5)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, idpId, idpEnabled, idpNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"      , cErr.Error()),
					slog.String("qry"        , qry),
					slog.Int   ("tntId"      , tntId),
					slog.Int   ("idpId"      , idpId),
					slog.Bool  ("idpEnabled" , idpEnabled),
					slog.String("idpNm"      , idpNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchIdp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpId int, idpNm string, idpEnabled bool, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = fmt.Sprintf("call %v.row_mod_idp(@p_tnt_id, @p_idp_id, @p_idp_nm, @p_idp_enabled, @p_by, @p_uts)", dbSchema)
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"      : tntId,
			"p_idp_id"      : idpId,
			"p_idp_nm"      : idpNm,
			"p_idp_enabled" : idpEnabled,
			"p_by"          : by,
			"p_uts"         : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"  , sprocCall),
			slog.String("error"      , sprocErr.Error()),
			slog.Int   ("tntId"      , tntId),
			slog.Int   ("idpId"      , idpId),
			slog.String("idpNm"      , idpNm),
			slog.Bool  ("idpEnabled" , idpEnabled),
			slog.String("by"         , by),
			slog.Any   ("uts"        , uts),
			slog.Any   ("exptErrs"   , exptErrs),
		)

		return sprocErr
	}

	return nil
}
