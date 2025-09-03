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
			qry := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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

type OccInf struct {
	OccId           int
	OcpNm           string
	OccEnabled      bool
	OccClientId     string
	OccClientSecret string
	OccCbUrl        string
	OccUrl          string
	Uts             time.Time
}

func GetOccInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]OccInf, error) {
	rs, rErr := db.DataSet[OccInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "occ_inf"
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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

func GetOccUtsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, occId int) ([]UtsInf, error) {
	rs, rErr := db.DataSet[UtsInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "occ_uts_inf"
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, occId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("occId" , occId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type OcpInf struct {
	OcpNm string
}

func GetOcpInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]OcpInf, error) {
	rs, rErr := db.DataSet[OcpInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ocp_inf"
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_atn_tnt_inf.%v($1, $2)", dbFunc)

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
		sprocCall   = "call web_core_auth_atn_tnt_mod.mod_aukc(@p_tnt_id, @p_aukc_aur_nm_min_len, @p_aukc_aur_nm_max_len, @p_aukc_enabled, @p_pka_id, @p_pkt_id, @p_pdc_id, @p_puv_reg_id, @p_puv_atn_id, @p_pkg_id, @p_prh_id, @p_pah_id, @p_by, @p_uts)"
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

func PatchAupc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aupcAurNmMinLen int, aupcAurNmMaxLen int, aupcAurPwdMinLen int, aupcAurPwdMaxLen int, aupcAurPwdIncSym bool, aupcAurPwdIncNum bool, aupcEnabled bool, aupcMfaEnabled bool, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_atn_tnt_mod.mod_aupc(@p_tnt_id, @p_aupc_aur_nm_min_len, @p_aupc_aur_nm_max_len, @p_aupc_aur_pwd_min_len, @p_aupc_aur_pwd_max_len, @p_aupc_aur_pwd_inc_sym, @p_aupc_aur_pwd_inc_num, @p_aupc_enabled, @p_aupc_mfa_enabled, @p_by, @p_uts)"
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

func PatchOcc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, occId int, occEnabled bool, occUrl string, occClientId string, occClientSecret string, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_atn_tnt_mod.mod_occ(@p_tnt_id, @p_occ_id, @p_occ_enabled, @p_occ_url, @p_occ_client_id, @p_occ_client_secret, @p_by, @p_uts)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"            : tntId,
			"p_occ_id"            : occId,
			"p_occ_enabled"       : occEnabled,
			"p_occ_url"           : occUrl,
			"p_occ_client_id"     : occClientId,
			"p_occ_client_secret" : occClientSecret,
			"p_by"                : by,
			"p_uts"               : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"   , sprocCall),
			slog.String("error"       , sprocErr.Error()),
			slog.Int   ("tntId"       , tntId),
			slog.Int   ("occId"       , occId),
			slog.Bool  ("occEnabled"  , occEnabled),
			slog.String("occUrl"      , occUrl),
			slog.String("occClientId" , occClientId),
			slog.String("by"          , by),
			slog.Any   ("uts"         , uts),
		)

		return sprocErr
	}

	return nil
}
