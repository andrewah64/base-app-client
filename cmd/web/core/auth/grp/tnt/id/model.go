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

type Inf struct {
	GrpId     int
	GrpNm     string
	GrpCanDel bool
	GrpCanEdt bool
	NumRoles  int
	NumUsers  int
}

type Mod struct {
	GrpId     int
	GrpNm     string
	GrpCanDel bool
	GrpCanEdt bool
	NumRoles  int
	NumUsers  int
	Uts       time.Time
}

const (
	dbSchema = "web_core_auth_grp_tnt_mod"
)

func GetRowGrpInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId int) ([]Inf, error) {
	const (
		dbFunc = "row_grp_inf"
	)

	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, grpId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("grpId" , grpId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetRowGrpMod (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId int) ([]Mod, error) {
	const (
		dbFunc = "row_grp_mod"
	)

	rs, rErr := db.DataSet[Mod](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, grpId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("grpId" , grpId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchGrp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId int, grpNm string, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = fmt.Sprintf("call %v.row_mod_grp(@p_tnt_id, @p_grp_id, @p_grp_nm, @p_by, @p_uts)", dbSchema)
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_grp_id" : grpId,
			"p_grp_nm" : grpNm,
			"p_by"     : by,
			"p_uts"    : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"  , sprocCall),
			slog.String("error"      , sprocErr.Error()),
			slog.Int   ("tntId"      , tntId),
			slog.Int   ("grpId"      , grpId),
			slog.String("grpNm"      , grpNm),
			slog.String("by"         , by),
			slog.Any   ("uts"        , uts),
			slog.Any   ("exptErrs"   , exptErrs),
		)

		return sprocErr
	}

	return nil
}
