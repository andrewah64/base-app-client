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

type Result struct {
	AuellId int
	AurNm   string
	EppPt   string
	HrmNm   string
	LvlNm   string
}

func Opts (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) (*map[string][]Opt, error) {
	const (
		dbFunc   = "ref_inf"
	)

	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_log_aur_tnt_inf.%v($1, $2)", dbFunc)

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

func GetLog(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, eppPt string, hrmId *int64, lvlId *int64, offset int, limit int) ([]Result, error) {
	const (
		dbFunc = "log_inf"
	)

	rs, rErr := db.DataSet[Result](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_log_aur_tnt_inf.%v($1, $2, $3, $4, $5, $6, $7, $8)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm, eppPt, hrmId, lvlId, offset, limit)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"  , cErr.Error()),
					slog.String("qry"    , qry),
					slog.Int   ("tntId"  , tntId),
					slog.String("aurNm"  , aurNm),
					slog.String("eppPt"  , eppPt),
					slog.Any   ("hrmId"  , hrmId),
					slog.Any   ("lvlId"  , lvlId),
					slog.Int   ("offset" , offset),
					slog.Int   ("limit"  , limit),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PutLog (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, eppPt string, hrmId *int64, lvlId *int64, tgtLvlId int, by string, exptErrs []string) error {
	var (
		sprocCall   = fmt.Sprintf("call web_core_auth_log_aur_tnt_mod.mod_log(@p_tnt_id, @p_aur_nm, @p_epp_pt, @p_hrm_id, @p_lvl_id, @p_tgt_lvl_id, @p_by)")
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_aur_nm"     : aurNm,
			"p_epp_pt"     : eppPt,
			"p_hrm_id"     : hrmId,
			"p_lvl_id"     : lvlId,
			"p_tgt_lvl_id" : tgtLvlId,
			"p_by"         : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.String("aurNm"     , aurNm),
			slog.String("eppPt"     , eppPt),
			slog.Any   ("hrmId"     , hrmId),
			slog.Any   ("lvlId"     , lvlId),
			slog.Int   ("tgtLvlId"  , tgtLvlId),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
