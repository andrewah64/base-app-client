package routes

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type Route struct {
	EndpointPath        string
	Handler             string
	MiddlewareChain     string
	HTTPRequestMethod   string
	Role              []*string
}

var (
	cache map[string]*Route = make(map[string]*Route)
)

func Add(ctx *context.Context, logger *slog.Logger, key string, route *Route){
	cache[key] = route
}

func CacheCopy () map[string]*Route {
	cp := make(map[string]*Route)

	for k, v := range cache {
		cp[k] = v
	}

	return cp
}

func Count () int {
	return len(cache)
}

func EndpointRoute(ctx *context.Context, logger *slog.Logger, endpoint string) (*Route, error) {
	logger.LogAttrs(*ctx, slog.LevelDebug, "get route",
		slog.String("endpoint", endpoint),
	)

	if route, ok := cache[endpoint]; ok {
		return route, nil
	} else {
		return nil, fmt.Errorf("endpoint '%v' not found", endpoint)
	}
}

func InitCache(ctx *context.Context, conn *pgxpool.Conn, dbSchema string, dbFunc string) error {
	slog.LogAttrs(*ctx, slog.LevelInfo, "initialise tenant cache")

	rs, rsErr := db.DataSet[Route](ctx, slog.Default(), conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "get route data",
				slog.String("error", err.Error()),
			)

			return qry, dbFunc, nil, fmt.Errorf("call database function: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	if rsErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "get route info",
			slog.String("error", rsErr.Error()),
		)
		return rsErr
	}

	for _, v := range rs {
		Add(ctx, slog.Default(), Key(v.HTTPRequestMethod, v.EndpointPath), &v)
	}

	if Count() == 0 || Count() != len(rs) {
		slog.LogAttrs(*ctx, slog.LevelError, "route cache is empty",
			slog.Int("Count()" , Count()),
			slog.Int("len(rs)" , len(rs)),
		)

		panic("The route cache was not initialised correctly")
	}

	return nil
}

func Key (hrm string, epp string) string{
	return fmt.Sprintf("%v/%v", hrm, epp)
}
