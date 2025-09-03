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

type Inf struct {
	AupcAurPwdMinLen int
	AupcAurPwdMaxLen int
	AupcAurPwdIncSym bool
	AupcAurPwdIncNum bool
}

func GetInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]Inf, error) {
	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "pwd_inf"
			qry    := fmt.Sprintf("select web_core_auth_pwd_aur_tnt_mod.%v($1, $2)", dbFunc)

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
