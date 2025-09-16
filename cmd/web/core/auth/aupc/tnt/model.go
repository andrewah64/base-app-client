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
			qry := fmt.Sprintf("select web_core_auth_aupc_tnt_inf.%v($1, $2)", dbFunc)

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

type AupcInf struct {
	TntId            int
	AupcAurNmMinLen  int
	AupcAurNmMaxLen  int
	AupcAurPwdMinLen int
	AupcAurPwdMaxLen int
	AupcAurPwdIncSym bool
	AupcAurPwdIncNum bool
	AupcEnabled      bool
	AupcMfaEnabled   bool
	Uts              time.Time
}

func GetAupcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AupcInf, error) {
	rs, rErr := db.DataSet[AupcInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aupc_inf"
			qry    := fmt.Sprintf("select web_core_auth_aupc_tnt_inf.%v($1, $2)", dbFunc)

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

func GetAupcUtsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]UtsInf, error) {
	rs, rErr := db.DataSet[UtsInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aupc_uts_inf"
			qry    := fmt.Sprintf("select web_core_auth_aupc_tnt_inf.%v($1, $2)", dbFunc)

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

func PatchAupc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aupcAurNmMinLen int, aupcAurNmMaxLen int, aupcAurPwdMinLen int, aupcAurPwdMaxLen int, aupcAurPwdIncSym bool, aupcAurPwdIncNum bool, aupcEnabled bool, aupcMfaEnabled bool, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_aupc_tnt_mod.mod_aupc(@p_tnt_id, @p_aupc_aur_nm_min_len, @p_aupc_aur_nm_max_len, @p_aupc_aur_pwd_min_len, @p_aupc_aur_pwd_max_len, @p_aupc_aur_pwd_inc_sym, @p_aupc_aur_pwd_inc_num, @p_aupc_enabled, @p_aupc_mfa_enabled, @p_by, @p_uts)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"               : tntId,
			"p_aupc_aur_nm_min_len"  : aupcAurNmMinLen,
			"p_aupc_aur_nm_max_len"  : aupcAurNmMaxLen,
			"p_aupc_aur_pwd_min_len" : aupcAurPwdMinLen,
			"p_aupc_aur_pwd_max_len" : aupcAurPwdMaxLen,
			"p_aupc_aur_pwd_inc_sym" : aupcAurPwdIncSym,
			"p_aupc_aur_pwd_inc_num" : aupcAurPwdIncNum,
			"p_aupc_enabled"         : aupcEnabled,
			"p_aupc_mfa_enabled"     : aupcMfaEnabled,
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
			slog.Int   ("aupcAurNmMinLen"  , aupcAurNmMinLen),
			slog.Int   ("aupcAurNmMaxLen"  , aupcAurNmMaxLen),
			slog.Int   ("aupcAurPwdMinLen" , aupcAurPwdMinLen),
			slog.Int   ("aupcAurPwdMaxLen" , aupcAurPwdMaxLen),
			slog.Bool  ("aupcAurPwdIncSym" , aupcAurPwdIncSym),
			slog.Bool  ("aupcAurPwdIncNum" , aupcAurPwdIncNum),
			slog.Bool  ("aupcEnabled"      , aupcEnabled),
			slog.Bool  ("aupcMfaEnabled"   , aupcMfaEnabled),
			slog.String("by"               , by),
			slog.Any   ("uts"              , uts),
		)

		return sprocErr
	}

	return nil
}
