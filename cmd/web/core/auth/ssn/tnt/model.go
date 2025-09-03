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

type Result struct {
	AurNm      string
	RolName    string
	WauhsSsnTk string
	Cts        time.Time
	WauhsExpTs time.Time
}

func DelSsn (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, wauhsSsnTk []string, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_ssn_tnt_del.del_ssn(@p_tnt_id, @p_wauhs_ssn_tk, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"       : tntId,
			"p_wauhs_ssn_tk" : wauhsSsnTk,
			"p_by"           : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"  , sprocCall),
			slog.String("error"      , sprocErr.Error()),
			slog.Int   ("tntId"      , tntId),
			slog.Any   ("wauhsSsnTk" , wauhsSsnTk),
			slog.String("by"         , by),
			slog.Any   ("exptErrs"   , exptErrs),
		)

		return sprocErr
	}

	return nil
}

func GetSsn(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurNm string, offset int, limit int) ([]Result, error) {
	rs, rErr := db.DataSet[Result](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			dbFunc := "ssn_inf"
			qry    := fmt.Sprintf("select web_core_auth_ssn_tnt_inf.%v($1, $2, $3, $4, $5)", dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurNm, offset, limit)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"   , cErr.Error()),
					slog.String("qry"     , qry),
					slog.Int   ("tntId"   , tntId),
					slog.String("aurNm"   , aurNm),
					slog.Int   ("offset"  , offset),
					slog.Int   ("limit"   , limit),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}
