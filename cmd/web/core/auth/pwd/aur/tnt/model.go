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

type PwdAurInf struct {
	AupcAurPwdMinLen int
	AupcAurPwdMaxLen int
	AupcAurPwdIncSym bool
	AupcAurPwdIncNum bool
	AurId            int
	AurNm            string
}

func GetPwdAurInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int) ([]PwdAurInf, error) {
	rs, rErr := db.DataSet[PwdAurInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "pwd_aur_inf"
			qry    := fmt.Sprintf("select web_core_auth_pwd_aur_tnt_mod.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"       , cErr.Error()),
					slog.String("qry"         , qry),
					slog.Int   ("tntId"       , tntId),
					slog.Int   ("aurId"       , aurId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchPwd (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aurHshPw string, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_pwd_aur_tnt_mod.mod_pwd(@p_tnt_id, @p_aur_id, @p_aur_hsh_pw, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"     : tntId,
			"p_aur_id"     : aurId,
			"p_aur_hsh_pw" : aurHshPw,
			"p_by"         : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("aurId"     , aurId),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
