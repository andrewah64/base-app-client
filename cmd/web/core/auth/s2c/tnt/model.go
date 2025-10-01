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

func DelSpc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcId []int, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_s2c_tnt_mod.del_spc(@p_tnt_id, @p_spc_id)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_spc_id" : spcId,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Any   ("spcId"     , spcId),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}

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
	SpcOk       bool
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

type SpcInf struct {
	SpcId      int
	SpcNm      string
	SpcCnNm    string
	SpcIncTs   time.Time
	SpcExpTs   time.Time
	SpcEnabled bool
}

func GetSpcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcNm string, spcIncTs *time.Time, spcExpTs *time.Time, spcEnabled *bool, offset int, limit int) ([]SpcInf, error) {
	rs, rErr := db.DataSet[SpcInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "spc_inf"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_inf.%v($1, $2, $3, $4, $5, $6, $7, $8)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, spcNm, spcIncTs, spcExpTs, spcEnabled, offset, limit)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"      , cErr.Error()),
					slog.String("qry"        , qry),
					slog.Int   ("tntId"      , tntId),
					slog.String("spcNm"      , spcNm),
					slog.Any   ("spcIncTs"   , spcIncTs),
					slog.Any   ("spcExpTs"   , spcExpTs),
					slog.Any   ("spcEnabled" , spcEnabled),
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

func PostIdp (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpNm string, idpEntityId string, ipcCrt [][]byte, cruNm []string, ipcIncTs []time.Time, ipcExpTs []time.Time, mdeUrl *string, sloUrl []string, sloUrlBnd []string, ssoUrl []string, ssoUrlBnd []string, by string, exptErrs []string) error {
	var (
		sprocCall = "call web_core_auth_s2c_tnt_mod.reg_idp(@p_tnt_id, @p_idp_nm, @p_idp_entity_id, @p_ipc_crt, @p_cru_nm, @p_ipc_inc_ts, @p_ipc_exp_ts, @p_mde_url, @p_slo_url, @p_slo_url_bnd, @p_sso_url, @p_sso_url_bnd, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"        : tntId,
			"p_idp_nm"        : idpNm,
			"p_idp_entity_id" : idpEntityId,
			"p_ipc_crt"       : ipcCrt,
			"p_cru_nm"        : cruNm,
			"p_ipc_inc_ts"    : ipcIncTs,
			"p_ipc_exp_ts"    : ipcExpTs,
			"p_mde_url"       : mdeUrl,
			"p_slo_url"       : sloUrl,
			"p_slo_url_bnd"   : sloUrlBnd,
			"p_sso_url"       : ssoUrl,
			"p_sso_url_bnd"   : ssoUrlBnd,
			"p_by"            : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"   , sprocCall),
			slog.String("error"       , sprocErr.Error()),
			slog.Int   ("tntId"       , tntId),
			slog.String("idpNm"       , idpNm),
			slog.String("idpEntityId" , idpEntityId),
			slog.Any   ("cruNm"       , cruNm),
			slog.Any   ("ipcIncTs"    , ipcIncTs),
			slog.Any   ("ipcExpTs"    , ipcExpTs),
			slog.Any   ("mdeUrl"      , mdeUrl),
			slog.Any   ("sloUrl"      , sloUrl),
			slog.Any   ("sloUrlBnd"   , sloUrlBnd),
			slog.Any   ("ssoUrl"      , ssoUrl),
			slog.Any   ("ssoUrlBnd"   , ssoUrlBnd),
			slog.String("by"          , by),
		)

		return sprocErr
	}

	return nil
}

func PostSpc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcNm string, spcCnNm string, spcOrgNm string, spcEncCrt []byte, spcEncPvk []byte, spcSgnCrt []byte, spcSgnPvk []byte, spcIncTs time.Time, spcExpTs time.Time, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_s2c_tnt_mod.reg_spc(@p_tnt_id, @p_spc_nm, @p_spc_cn_nm, @p_spc_org_nm, @p_spc_enc_crt, @p_spc_enc_pvk, @p_spc_sgn_crt, @p_spc_sgn_pvk, @p_spc_inc_ts, @p_spc_exp_ts, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"      : tntId,
			"p_spc_nm"      : spcNm,
			"p_spc_cn_nm"   : spcCnNm,
			"p_spc_org_nm"  : spcOrgNm,
			"p_spc_enc_crt" : spcEncCrt,
			"p_spc_enc_pvk" : spcEncPvk,
			"p_spc_sgn_crt" : spcSgnCrt,
			"p_spc_sgn_pvk" : spcSgnPvk,
			"p_spc_inc_ts"  : spcIncTs,
			"p_spc_exp_ts"  : spcExpTs,
			"p_by"          : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.String("spcNm"     , spcNm),
			slog.String("spcCnNm"   , spcCnNm),
			slog.String("spcOrgNm"  , spcOrgNm),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
