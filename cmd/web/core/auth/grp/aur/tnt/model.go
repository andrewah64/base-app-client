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

type Aur struct {
	AurId int
	AurNm string
}

type Grp struct {
	GrpId     int
	GrpNm     string
	TgtAurHas bool
	CurAurEdt bool
}

func GetAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int) ([]Aur, error) {
	rs, rErr := db.DataSet[Aur](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aur_inf"
			qry    := fmt.Sprintf("select web_core_auth_grp_aur_tnt_mod.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"       , cErr.Error()),
					slog.String("qry"         , qry),
					slog.Int   ("tntId"       , tntId),
					slog.Int   ("aurId"       , aurId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetGrp(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, curAurId int, tgtAurId int) ([]Grp, error) {
	const (
		dbSchema = "web_core_auth_grp_aur_tnt_mod"
		dbFunc   = "grp_inf"
	)

	rs, rErr := db.DataSet[Grp](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, curAurId, tgtAurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"    , cErr.Error()),
					slog.String("qry"      , qry),
					slog.Int   ("tntId"    , tntId),
					slog.Int   ("curAurId" , curAurId),
					slog.Int   ("tgtAurId" , tgtAurId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchGrp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, curAurId int, tgtAurId int, grpId []int, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_grp_aur_tnt_mod.mod_grp(@p_tnt_id, @p_cur_aur_id, @p_tgt_aur_id, @p_grp_id, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_cur_aur_id" : curAurId,
			"p_tgt_aur_id" : tgtAurId,
			"p_grp_id"     : grpId,
			"p_by"         : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("curAurId"  , curAurId),
			slog.Int   ("tgtAurId"  , tgtAurId),
			slog.Any   ("grpId"     , grpId),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
