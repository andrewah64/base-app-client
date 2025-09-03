package id

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

type Inf struct {
	AaukId      int
	AaukNm      string
	AaukEnabled bool
	NumRoles    int
}

type Mod struct {
	AaukId      int
	AaukNm      string
	AaukEnabled bool
	NumRoles    int
	Uts         time.Time
}

const (
	dbSchema = "web_core_auth_key_aur_mod"
)

func GetRowKeyInf(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukId int) ([]Inf, error) {
	const (
		dbFunc = "row_key_inf"
	)

	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, aaukId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"  , cErr.Error()),
					slog.String("qry"    , qry),
					slog.Int   ("tntId"  , tntId),
					slog.Int   ("aurId"  , aurId),
					slog.Int   ("aaukId" , aaukId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetRowKeyMod(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukId int) ([]Mod, error) {
	const (
		dbFunc = "row_key_mod"
	)

	rs, rErr := db.DataSet[Mod](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId, aaukId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"  , cErr.Error()),
					slog.String("qry"    , qry),
					slog.Int   ("tntId"  , tntId),
					slog.Int   ("aurId"  , aurId),
					slog.Int   ("aaukId" , aaukId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func PatchKey (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukId int, aaukNm string, aaukEnabled bool, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = fmt.Sprintf("call %v.row_mod_key(@p_tnt_id, @p_aur_id, @p_aauk_id, @p_aauk_nm, @p_aauk_enabled, @p_by, @p_uts)", dbSchema)
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"       : tntId,
			"p_aur_id"       : aurId,
			"p_aauk_id"      : aaukId,
			"p_aauk_nm"      : aaukNm,
			"p_aauk_enabled" : aaukEnabled,
			"p_by"           : by,
			"p_uts"          : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"   , sprocCall),
			slog.String("error"       , sprocErr.Error()),
			slog.Int   ("tntId"       , tntId),
			slog.Int   ("aurId"       , aurId),
			slog.Int   ("aaukId"      , aaukId),
			slog.String("aaukNm"      , aaukNm),
			slog.Bool  ("aaukEnabled" , aaukEnabled),
			slog.String("by"          , by),
			slog.Any   ("uts"         , uts),
			slog.Any   ("exptErrs"    , exptErrs),
		)

		return sprocErr
	}

	return nil
}
