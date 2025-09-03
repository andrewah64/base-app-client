package tenant

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

var (
	cache map[string]int = make(map[string]int)
)

func InitCache(ctx *context.Context, conn *pgxpool.Conn) error {
	slog.LogAttrs(*ctx, slog.LevelInfo, "initialise tenant cache")

	const (
		dbSchema = "all_core_unauth_tnt_all_inf"
		dbFunc   = "tnt_inf"
	)

	type tenant struct {
		TntId     int
		TntPrtc   string
		TntFqdn   string
		TntPort   int
		TntOrigin string
	}

	rs, rsErr := db.DataSet[tenant](ctx, slog.Default(), conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "get tenant data",
				slog.String("error", err.Error()),
			)

			return qry, dbFunc, nil, fmt.Errorf("call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	if rsErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "get tenant info",
			slog.String("error", rsErr.Error()),
		)
		return rsErr
	}

	for _, v := range rs {
		cache[v.TntOrigin] = v.TntId
	}

	if len(cache) == 0 {
		slog.LogAttrs(*ctx, slog.LevelError, "tenant cache is empty")

		panic("The tenant cache is empty")
	}

	return nil
}

func Origin (r *http.Request) string {
	return fmt.Sprintf("https://%v", r.Host)
}

func Tenant(ctx *context.Context, logger *slog.Logger, origin string) int {
	logger.LogAttrs(*ctx, slog.LevelDebug, "get tenant",
		slog.String("origin", origin),
	)

	if tntId, ok := cache[origin]; ok {
		return tntId
	} else {
		slog.LogAttrs(*ctx, slog.LevelError, "tenant not found",
			slog.String("origin", origin),
		)
		panic("tenant not found")
	}
}
