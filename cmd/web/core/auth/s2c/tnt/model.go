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

func Opts (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) (*map[string][]Opt, error) {
	const (
		dbFunc = "ref_inf"
	)

	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select web_core_auth_s2c_tnt_inf.%v($1, $2)", dbFunc)

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

type UtsInf struct {
	Uts time.Time
}

func GetS2cUtsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]UtsInf, error) {
	rs, rErr := db.DataSet[UtsInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "s2c_uts_inf"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_inf.%v($1, $2)", dbFunc)

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

func GetS2gUtsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]UtsInf, error) {
	rs, rErr := db.DataSet[UtsInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "s2g_uts_inf"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_inf.%v($1, $2)", dbFunc)

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

type S2cInf struct {
	S2cEntityId string
	S2cEnabled  bool
	AumId       int
	EppAcsPt    string
	EppMtdPt    string
	Uts         time.Time
}

func GetS2cInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]S2cInf, error) {
	rs, rErr := db.DataSet[S2cInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "s2c_inf"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_inf.%v($1, $2)", dbFunc)

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

type S2gInf struct {
	S2gCrtCn  string
	S2gCrtOrg string
	Uts       time.Time
}

func GetS2gInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]S2gInf, error) {
	rs, rErr := db.DataSet[S2gInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "s2g_inf"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_inf.%v($1, $2)", dbFunc)

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

func PatchS2c (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, s2cEnabled bool, s2cEntityId string, aumId int, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_s2c_tnt_mod.mod_s2c(@p_tnt_id, @p_s2c_enabled, @p_s2c_entity_id, @p_aum_id, @p_by, @p_uts)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"        : tntId,
			"p_s2c_enabled"   : s2cEnabled,
			"p_s2c_entity_id" : s2cEntityId,
			"p_aum_id"        : aumId,
			"p_by"            : by,
			"p_uts"           : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"   , sprocCall),
			slog.String("error"       , sprocErr.Error()),
			slog.Int   ("tntId"       , tntId),
			slog.Bool  ("s2cEnabled"  , s2cEnabled),
			slog.String("s2cEntityId" , s2cEntityId),
			slog.Int   ("aumId"       , aumId),
			slog.String("by"          , by),
			slog.Any   ("uts"         , uts),
			slog.Any   ("exptErrs"    , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func PatchS2g (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, s2gCrtCn string, s2gCrtOrg string, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_s2c_tnt_mod.mod_s2g(@p_tnt_id, @p_s2g_crt_cn, @p_s2g_crt_org, @p_by, @p_uts)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"      : tntId,
			"p_s2g_crt_cn"  : s2gCrtCn,
			"p_s2g_crt_org" : s2gCrtOrg,
			"p_by"          : by,
			"p_uts"         : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.String("s2gCrtCn"  , s2gCrtCn),
			slog.String("s2gCrtOrg" , s2gCrtOrg),
			slog.String("by"        , by),
			slog.Any   ("uts"       , uts),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
