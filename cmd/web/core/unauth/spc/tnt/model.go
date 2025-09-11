package tnt

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

type SpcInf struct {
	S2cEntityId string
	S2cAcsUrl   string
	SpcSgnCrt   []byte
	SpcEncCrt   []byte
}

func GetSpc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int) ([]SpcInf, error) {
	rs, rErr := db.DataSet[SpcInf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "spc_inf"
			qry    := fmt.Sprintf("select web_core_unauth_spc_tnt_inf.%v($1, $2)", dbFunc)

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
