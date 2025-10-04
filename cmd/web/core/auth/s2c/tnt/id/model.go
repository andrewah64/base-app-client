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
	IdpNmOk bool
}

func GetRowIdpVal (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpId int, idpNm string) ([]RowIdpVal, error) {
	const (
		dbFunc = "row_idp_val"
	)

	rs, rErr := db.DataSet[RowIdpVal](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, idpId, idpNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"      , cErr.Error()),
					slog.String("qry"        , qry),
					slog.Int   ("tntId"      , tntId),
					slog.Int   ("idpId"      , idpId),
					slog.String("idpNm"      , idpNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type RowSpcInf struct {
	SpcId      int
	SpcNm      string
	SpcCnNm    string
	SpcIncTs   time.Time
	SpcExpTs   time.Time
	SpcEnabled bool
}

func GetRowSpcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcId int) ([]RowSpcInf, error) {
	const (
		dbFunc = "row_spc_inf"
	)

	rs, rErr := db.DataSet[RowSpcInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, spcId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("spcId" , spcId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type RowSpcMod struct {
	SpcId      int
	SpcNm      string
	SpcCnNm    string
	SpcIncTs   time.Time
	SpcExpTs   time.Time
	SpcEnabled bool
	Uts        time.Time
}

func GetRowSpcMod (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcId int) ([]RowSpcMod, error) {
	const (
		dbFunc = "row_spc_mod"
	)

	rs, rErr := db.DataSet[RowSpcMod](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, spcId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("spcId" , spcId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type RowSpcVal struct {
	SpcNmOk      bool
	SpcEnabledOk bool
	SpcTsOk      bool
}

func GetRowSpcVal (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcId int, spcNm string, spcEnabled bool) ([]RowSpcVal, error) {
	const (
		dbFunc = "row_spc_val"
	)

	rs, rErr := db.DataSet[RowSpcVal](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4, $5)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, spcId, spcNm, spcEnabled)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"      , cErr.Error()),
					slog.String("qry"        , qry),
					slog.Int   ("tntId"      , tntId),
					slog.Int   ("spcId"      , spcId),
					slog.String("spcNm"      , spcNm),
					slog.Bool  ("spcEnabled" , spcEnabled),
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

func PatchSpc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcId int, spcNm string, spcEnabled bool, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = fmt.Sprintf("call %v.row_mod_spc(@p_tnt_id, @p_spc_id, @p_spc_nm, @p_spc_enabled, @p_by, @p_uts)", dbSchema)
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"      : tntId,
			"p_spc_id"      : spcId,
			"p_spc_nm"      : spcNm,
			"p_spc_enabled" : spcEnabled,
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
			slog.Int   ("spcId"      , spcId),
			slog.String("spcNm"      , spcNm),
			slog.Bool  ("spcEnabled" , spcEnabled),
			slog.String("by"         , by),
			slog.Any   ("uts"        , uts),
			slog.Any   ("exptErrs"   , exptErrs),
		)

		return sprocErr
	}

	return nil
}
