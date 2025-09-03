package tnt

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type Grp struct {
	GrpId      int
	GrpNm      string
	GrpCanEdt  bool
}

type Rol struct {
	DbrlId     int
	DbrlDs     string
	DbrlHas    bool
}

func GetGrp(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, grpId int) ([]Grp, error) {
	const (
		dbFunc = "grp_inf"
	)

	rs, rErr := db.DataSet[Grp](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_rol_grp_tnt_inf.%v($1, $2, $3, $4)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, grpId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("aurId" , aurId),
					slog.Int   ("grpId" , grpId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetRol(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, grpId int) ([]Rol, error) {
	const (
		dbFunc = "rol_inf"
	)

	rs, rErr := db.DataSet[Rol](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_rol_grp_tnt_inf.%v($1, $2, $3, $4)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, grpId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("aurId" , aurId),
					slog.Int   ("grpId" , grpId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchRol (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId int, dbrlId []int, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_rol_grp_tnt_mod.mod_rol(@p_tnt_id, @p_grp_id, @p_dbrl_id, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"  : tntId,
			"p_grp_id"  : grpId,
			"p_dbrl_id" : dbrlId,
			"p_by"      : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("grpId"     , grpId),
			slog.Any   ("dbrlId"    , dbrlId),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
