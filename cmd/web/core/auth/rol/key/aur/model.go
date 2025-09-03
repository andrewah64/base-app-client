package aur

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

type Key struct {
	AaukId   int
	AaukNm   string
}

type Rol struct {
	DbrlId   int
	DbrlDs   string
	Assigned bool
}

func GetKey(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aaukId int) ([]Key, error) {
	const (
		dbSchema = "web_core_auth_key_aur_mod"
		dbFunc   = "key_inf"
	)

	rs, rErr := db.DataSet[Key](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aaukId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error"  , cErr.Error()),
					slog.String("qry"    , qry),
					slog.Int   ("tntId"  , tntId),
					slog.Int   ("aaukId" , aaukId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetRol(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukId int) ([]Rol, error) {
	const (
		dbSchema = "web_core_auth_key_aur_mod"
		dbFunc   = "rol_inf"
	)

	rs, rErr := db.DataSet[Rol](ctx, logger, conn,
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

func PatchRol (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aaukId int, dbrlId []int, by string, exptErrs []string) error {
	var (
		sprocCall   = "call web_core_auth_key_aur_mod.mod_rol(@p_tnt_id, @p_aur_id, @p_aauk_id, @p_dbrl_id, @p_by)"
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"  : tntId,
			"p_aur_id"  : aurId,
			"p_aauk_id" : aaukId,
			"p_dbrl_id" : dbrlId,
			"p_by"      : by,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall" , sprocCall),
			slog.String("error"     , sprocErr.Error()),
			slog.Int   ("tntId"     , tntId),
			slog.Int   ("aurId"     , aurId),
			slog.Int   ("aaukId"    , aaukId),
			slog.Any   ("dbrlId"    , dbrlId),
			slog.String("by"        , by),
			slog.Any   ("exptErrs"  , exptErrs),
		)

		return sprocErr
	}

	return nil
}
