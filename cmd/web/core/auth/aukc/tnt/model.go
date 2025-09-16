package tnt

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

type Opt struct {
	Key   string
	Id    int
	Value string
}

type UtsInf struct {
	Uts time.Time
}

func Opts (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) (*map[string][]Opt, error) {
	const (
		dbFunc   = "ref_inf"
	)

	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_aukc_tnt_inf.%v($1, $2)", dbFunc)

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

type AukcInf struct {
	TntId            int
	AukcAurNmMinLen  int
	AukcAurNmMaxLen  int
	AukcEnabled      bool
	PkaId            int
	PktId            int
	PdcId            int
	PuvRegId         int
	PuvAtnId         int
	Uts              time.Time
}

func GetAukcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AukcInf, error) {
	rs, rErr := db.DataSet[AukcInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aukc_inf"
			qry    := fmt.Sprintf("select web_core_auth_aukc_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetAukcUtsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]UtsInf, error) {
	rs, rErr := db.DataSet[UtsInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aukc_uts_inf"
			qry    := fmt.Sprintf("select web_core_auth_aukc_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type PahInf struct {
	PkhId    int
	PkhNm    string
	Assigned bool
}

func GetPahInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]PahInf, error) {
	rs, rErr := db.DataSet[PahInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "pah_inf"
			qry    := fmt.Sprintf("select web_core_auth_aukc_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type PkgInf struct {
	PkgId    int
	PkgNm    string
	Assigned bool
}

func GetPkgInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]PkgInf, error) {
	rs, rErr := db.DataSet[PkgInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "pkg_inf"
			qry    := fmt.Sprintf("select web_core_auth_aukc_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type PrhInf struct {
	PkhId    int
	PkhNm    string
	Assigned bool
}

func GetPrhInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]PrhInf, error) {
	rs, rErr := db.DataSet[PrhInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "prh_inf"
			qry    := fmt.Sprintf("select web_core_auth_aukc_tnt_inf.%v($1, $2)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchAukc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aukcAurNmMinLen int, aukcAurNmMaxLen int, aukcEnabled bool, pkaId int, pktId int, pdcId int, puvRegId int, puvAtnId int, pkgId []int, prhId []int, pahId []int, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_aukc_tnt_mod.mod_aukc(@p_tnt_id, @p_aukc_aur_nm_min_len, @p_aukc_aur_nm_max_len, @p_aukc_enabled, @p_pka_id, @p_pkt_id, @p_pdc_id, @p_puv_reg_id, @p_puv_atn_id, @p_pkg_id, @p_prh_id, @p_pah_id, @p_by, @p_uts)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"               : tntId,
			"p_aukc_aur_nm_min_len"  : aukcAurNmMinLen,
			"p_aukc_aur_nm_max_len"  : aukcAurNmMaxLen,
			"p_aukc_enabled"         : aukcEnabled,
			"p_pka_id"               : pkaId,
			"p_pkt_id"               : pktId,
			"p_pdc_id"               : pdcId,
			"p_puv_reg_id"           : puvRegId,
			"p_puv_atn_id"           : puvAtnId,
			"p_pkg_id"               : pkgId,
			"p_prh_id"               : prhId,
			"p_pah_id"               : pahId,
			"p_by"                   : by,
			"p_uts"                  : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"        , sprocCall),
			slog.String("error"            , sprocErr.Error()),
			slog.Int   ("tntId"            , tntId),
			slog.Int   ("aukcAurNmMinLen"  , aukcAurNmMinLen),
			slog.Int   ("aukcAurNmMaxLen"  , aukcAurNmMaxLen),
			slog.Bool  ("aukcEnabled"      , aukcEnabled),
			slog.Int   ("pkaId"            , pkaId),
			slog.Int   ("pktId"            , pktId),
			slog.Int   ("pdcId"            , pdcId),
			slog.Any   ("pkgId"            , pkgId),
			slog.Any   ("prhId"            , prhId),
			slog.Any   ("pahId"            , pahId),
			slog.String("by"               , by),
		)

		return sprocErr
	}

	return nil
}
