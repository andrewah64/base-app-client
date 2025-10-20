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
			qry    := fmt.Sprintf("select web_core_unauth_saml2_acs_inf.%v($1, $2)", dbFunc)

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
