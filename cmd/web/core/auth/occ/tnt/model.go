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

type UtsInf struct {
	Uts time.Time
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
			qry    := fmt.Sprintf("select web_core_auth_occ_tnt_inf.%v($1, $2)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_occ_tnt_inf.%v($1, $2, $3)", dbFunc)

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
			qry    := fmt.Sprintf("select web_core_auth_occ_tnt_inf.%v($1, $2)", dbFunc)

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

func PatchOcc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, occId int, occEnabled bool, occUrl string, occClientId string, occClientSecret string, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_occ_tnt_mod.mod_occ(@p_tnt_id, @p_occ_id, @p_occ_enabled, @p_occ_url, @p_occ_client_id, @p_occ_client_secret, @p_by, @p_uts)"
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
