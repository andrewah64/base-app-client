package aur

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type Inf struct {
	AaukId      int
	AaukNm      string
	AaukEnabled bool
	NumRoles    int
}

type Opt struct {
	Key         string
	Id          int
	Value       string
}

func OptsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int) (*map[string][]Opt, error) {
	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ref_inf"
			qry    := fmt.Sprintf("select web_core_auth_key_aur_mod.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"   , cErr.Error()),
					slog.String("qry"     , qry),
					slog.Int   ("tntId"   , tntId),
					slog.Int   ("aurId"   , aurId),
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

func DelKey (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukId []int, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_key_aur_mod.del_key(@p_tnt_id, @p_aur_id, @p_aauk_id)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"  : tntId,
			"p_aur_id"  : aurId,
			"p_aauk_id" : aaukId,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("aurId"     , aurId),
			slog.Any   ("aaukId"    , aaukId),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func GetKey (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukNm string, aaukEnabled *bool, dbrlId *int64, offset int, limit int) ([]Inf, error) {
	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "key_inf"
			qry    := fmt.Sprintf("select web_core_auth_key_aur_mod.%v($1, $2, $3, $4, $5, $6, $7, $8)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, aaukNm, aaukEnabled, dbrlId, offset, limit)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"       , cErr.Error()),
					slog.String("qry"         , qry),
					slog.Int   ("tntId"       , tntId),
					slog.Int   ("aurId"       , aurId),
					slog.String("aaukNm"      , aaukNm),
					slog.Any   ("aaukEnabled" , aaukEnabled),
					slog.Any   ("dbrlId"      , dbrlId),
					slog.Int   ("offset"      , offset),
					slog.Int   ("limit"       , limit),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PostKey(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukKey []byte, aaukEnabled bool, aaukNm string, by string, exptErrs []string) error {
	const (
		dbSchema = "web_core_auth_key_aur_mod"
		dbSproc  = "reg_key"
	)

	var (
		sprocCall   = fmt.Sprintf("call %v.%v(@p_tnt_id, @p_aur_id, @p_aauk_key, @p_aauk_enabled, @p_aauk_nm, @p_by)", dbSchema, dbSproc)
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"       : tntId,
			"p_aur_id"       : aurId,
			"p_aauk_key"     : aaukKey,
			"p_aauk_enabled" : aaukEnabled,
			"p_aauk_nm"      : aaukNm,
			"p_by"           : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelError, "call sproc",
			slog.String("error"       , sprocErr.Error()),
			slog.String("sprocCall"   , sprocCall),
		)
		return sprocErr
	}

	return nil
}
