package tnt

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type AcsInf struct {
	SsoUrl      string
	SsoBndNm    string
	IdpEntityId string
	AcsEppPt    string
	S2cEntityId string
	IpcCrt      [][]byte
}

func GetAcsInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]AcsInf, error) {
	rs, rErr := db.DataSet[AcsInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "acs_inf"
			qry    := fmt.Sprintf("select web_core_unauth_saml2_acs_mod.%v($1, $2)", dbFunc)

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

type AurInf struct {
	AurId int
	SsnDn time.Duration
	EppPt string
}

func GetAurInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurEa string) ([]AurInf, error) {
	rs, rErr := db.DataSet[AurInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "aur_inf"
			qry    := fmt.Sprintf("select web_core_unauth_saml2_acs_mod.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurEa)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("aurEa" , aurEa),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})
	
	return rs, rErr
}

func RegAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurEa string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_unauth_saml2_acs_mod.reg_aur(@p_tnt_id, @p_aur_ea)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id" : tntId,
			"p_aur_ea" : aurEa,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.String("aurEa"     , aurEa),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
