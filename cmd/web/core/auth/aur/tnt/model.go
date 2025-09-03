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

type Opt struct {
	Key   string
	Id    int
	Value string
}

type Inf struct {
	AurId        int
	AurNm        string
	RolName      string
	AurEnabled   bool
	LngNm        string
	PgNm         string
}

func OptsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) (*map[string][]Opt, error) {
	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ref_inf"
			qry    := fmt.Sprintf("select web_core_auth_aur_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"   , cErr.Error()),
					slog.String("qry"     , qry),
					slog.Int   ("tntId"  , tntId),
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

func OptsReg (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) (*map[string][]Opt, error) {
	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ref_inf"
			qry    := fmt.Sprintf("select web_core_auth_aur_tnt_reg.%v($1, $2)", dbFunc)

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

func DelAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId []int, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_aur_tnt_del.del_aur(@p_tnt_id, @p_aur_id)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_aur_id" : aurId,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Any   ("aurId"     , aurId),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func GetAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, aurEnabled *bool, dbrlId *int64, lngId *int64, offset int, limit int) ([]Inf, error) {
	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aur_inf"
			qry    := fmt.Sprintf("select web_core_auth_aur_tnt_inf.%v($1, $2, $3, $4, $5, $6, $7, $8)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm, aurEnabled, dbrlId, lngId, offset, limit)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"      , cErr.Error()),
					slog.String("qry"        , qry),
					slog.Int   ("tntId"      , tntId),
					slog.String("aurNm"      , aurNm),
					slog.Any   ("aurEnabled" , aurEnabled),
					slog.Any   ("dbrlId"     , dbrlId),
					slog.Any   ("lngId"      , lngId),
					slog.Int   ("offset"     , offset),
					slog.Int   ("limit"      , limit),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PostAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, grpId *int64, aurNm string, aurHshPw string, lngId *int64, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_aur_tnt_reg.reg_aur(@p_tnt_id, @p_grp_id, @p_aur_nm, @p_aur_hsh_pw, @p_lng_id, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_grp_id"     : grpId,
			"p_aur_nm"     : aurNm,
			"p_aur_hsh_pw" : aurHshPw,
			"p_lng_id"     : lngId,
			"p_by"         : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Any   ("grpId"     , grpId),
			slog.String("aurNm"     , aurNm),
			slog.Any   ("lngId"     , lngId),
			slog.Any   ("exptErrs"  , exptErrs),
			slog.String("by"        , by),
		)

		return sprocErr
	}

	return nil
}
