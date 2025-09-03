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
	AurId  int
	AurNm  string
	AurHas bool
}

type Grp struct {
	GrpId      int
	GrpNm      string
	GrpCanEdt  bool
}

func GetAur(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId int) ([]Aur, error) {
	const (
		dbFunc = "aur_inf"
	)

	rs, rErr := db.DataSet[Aur](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_aur_grp_tnt_inf.%v($1, $2, $3)", dbFunc)

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

func GetGrp(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, grpId int) ([]Grp, error) {
	const (
		dbFunc = "grp_inf"
	)

	rs, rErr := db.DataSet[Grp](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_aur_grp_tnt_inf.%v($1, $2, $3, $4)", dbFunc)

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

func PatchAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId int, aurId []int, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_aur_grp_tnt_mod.mod_aur(@p_tnt_id, @p_grp_id, @p_aur_id)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_grp_id" : grpId,
			"p_aur_id" : aurId,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("grpId"     , grpId),
			slog.Any   ("aurId"     , aurId),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
