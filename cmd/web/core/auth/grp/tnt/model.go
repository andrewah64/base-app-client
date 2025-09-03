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

type Inf struct {
	GrpId     int
	GrpNm     string
	GrpCanDel bool
	GrpCanEdt bool
	NumRoles  int
	NumUsers  int
}

type Opt struct {
	Key   string
	Id    int
	Value string
}

func DelGrp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId []int, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_grp_tnt_del.del_grp(@p_tnt_id, @p_grp_id)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_grp_id" : grpId,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Any   ("grpId"     , grpId),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func GetGrp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, grpNm string, aurNm string, dbrlId *int64, offset int, limit int) ([]Inf, error) {
	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "grp_inf"
			qry    := fmt.Sprintf("select web_core_auth_grp_tnt_inf.%v($1, $2, $3, $4, $5, $6, $7, $8)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, grpNm, aurNm, dbrlId, offset, limit)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("cErr.Error()" , cErr.Error()),
					slog.String("qry"          , qry),
					slog.Int   ("tntId"        , tntId),
					slog.Int   ("aurId"        , aurId),
					slog.String("grpNm"        , grpNm),
					slog.String("aurNm"        , aurNm),
					slog.Any   ("dbrlId"       , dbrlId),
					slog.Int   ("offset"       , offset),
					slog.Int   ("limit"        , limit),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func OptsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) (*map[string][]Opt, error) {
	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ref_inf"
			qry    := fmt.Sprintf("select web_core_auth_grp_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"   , cErr.Error()),
					slog.String("qry"     , qry),
					slog.Int   ("tntId"   , tntId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		},
	)

	if rErr != nil {
		return nil, fmt.Errorf("get Opts dataset: %w", rErr)
	}

	idValMap := make(map[string][]Opt)

	for _, v := range rs {
		idValMap[v.Key] = append(idValMap[v.Key], Opt{Id: v.Id, Value: v.Value})
	}

	return &idValMap, nil
}

func PostGrp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpNm string, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_grp_tnt_reg.reg_grp(@p_tnt_id, @p_grp_nm, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_grp_nm" : grpNm,
			"p_by"     : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Any   ("grpNm"     , grpNm),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
