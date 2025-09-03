package val

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

type AurEaInf struct {
	AurEaAvbPass bool
}

func GetAurEaInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurEa string) ([]AurEaInf, error) {
	rs, rErr := db.DataSet[AurEaInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ea_inf"
			qry    := fmt.Sprintf("select web_core_unauth_aur_tnt_reg.%v($1, $2, $3)", dbFunc)

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

type AurNmInf struct {
	AurNmLenPass bool
	AurNmAvbPass bool
}

func GetAurNmInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string) ([]AurNmInf, error) {
	rs, rErr := db.DataSet[AurNmInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "nm_inf"
			qry    := fmt.Sprintf("select web_core_unauth_aur_tnt_reg.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("aurNm" , aurNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type PwdInf struct {
	AurPwdMinLen int
	AurPwdMaxLen int
	AurPwdIncSym bool
	AurPwdIncNum bool
}

func GetPwdInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]PwdInf, error) {
	rs, rErr := db.DataSet[PwdInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "pwd_inf"
			qry    := fmt.Sprintf("select web_core_unauth_aur_tnt_reg.%v($1, $2)", dbFunc)

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
