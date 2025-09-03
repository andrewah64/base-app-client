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

type Opt struct {
	Key   string
	Id    int
	Value string
}

type Inf struct {
	AurId      int
	AurNm      string
	RolName    string
	AurEnabled bool
	LngNm      string
	PgNm       string
}

type Mod struct {
	AurId      int
	AurNm      string
	RolName    string
	AurEnabled bool
	LngId      int
	PgId       int
	Uts        time.Time
}

const (
	dbSchema = "web_core_auth_aur_tnt_mod"
)

func GetRowAurInf(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int) ([]Inf, error) {
	const (
		dbFunc = "row_aur_inf"
	)

	rs, rErr := db.DataSet[Inf](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("aurId" , aurId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func GetRowAurMod(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int) ([]Mod, error) {
	const (
		dbFunc = "row_aur_mod"
	)

	rs, rErr := db.DataSet[Mod](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("aurId" , aurId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return rs, rErr
}

func Opts (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int) (*map[string][]Opt, error) {
	const (
		dbFunc = "row_ref_inf"
	)

	rs, rErr := db.DataSet[Opt](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tntId, aurId)
			if cErr != nil {
				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("error" , cErr.Error()),
					slog.String("qry"   , qry),
					slog.Int   ("tntId" , tntId),
					slog.Int   ("aurId" , aurId),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		},
	)

	if rErr != nil {
		return nil, fmt.Errorf("get Opts dataset: %w", rErr)
	}

	idValMap := make(map[string][]Opt)

	for _, v := range rs {
		idValMap[v.Key] = append(idValMap[v.Key], Opt{Id: v.Id, Value: v.Value})
	}

	return &idValMap, nil
}

func PatchAur (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tntId int, aurId int, aurNm string, aurEnabled bool, lngId int, pgId int, by string, uts time.Time, exptErrs []string) error {
	var (
		sprocCall   = fmt.Sprintf("call %v.row_mod_aur(@p_tnt_id, @p_aur_id, @p_aur_nm, @p_aur_enabled, @p_lng_id, @p_pg_id, @p_by, @p_uts)", dbSchema)
		sprocParams = pgx.NamedArgs{
			"p_tnt_id"      : tntId,
			"p_aur_id"      : aurId,
			"p_aur_nm"      : aurNm,
			"p_aur_enabled" : aurEnabled,
			"p_lng_id"      : lngId,
			"p_pg_id"       : pgId,
			"p_by"          : by,
			"p_uts"         : uts,
		}
	)

	sprocErr := db.Sproc(ctx, logger, conn, sprocCall, sprocParams, exptErrs)
	if sprocErr != nil {
		logger.LogAttrs(*ctx, slog.LevelDebug, "call sproc",
			slog.String("sprocCall"  , sprocCall),
			slog.String("error"      , sprocErr.Error()),
			slog.Int   ("tntId"      , tntId),
			slog.Int   ("aurId"      , aurId),
			slog.String("aurNm"      , aurNm),
			slog.Bool  ("aurEnabled" , aurEnabled),
			slog.Int   ("lngId"      , lngId),
			slog.Int   ("pgId"       , pgId),
			slog.String("by"         , by),
			slog.Any   ("uts"        , uts),
			slog.Any   ("exptErrs"   , exptErrs),
		)

		return sprocErr
	}

	return nil
}
