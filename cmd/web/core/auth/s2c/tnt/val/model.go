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

type IdpInf struct {
	IdpNmOk bool
}

func GetIdpInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, idpNm string) ([]IdpInf, error) {
	rs, rErr := db.DataSet[IdpInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "idp_val"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_mod.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, idpNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("idpNm" , idpNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

type SpcInf struct {
	SpcNmOk bool
}

func GetSpcInf (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, spcNm string) ([]SpcInf, error) {
	rs, rErr := db.DataSet[SpcInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "spc_val"
			qry    := fmt.Sprintf("select web_core_auth_s2c_tnt_mod.%v($1, $2, $3)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, spcNm)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.String("spcNm" , spcNm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}
