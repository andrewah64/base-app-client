package route

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justinas/alice"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/routes"
	"github.com/andrewah64/base-app-client/internal/api/core/mw"
)

func Mux(ctx *context.Context, handlers map[string]http.HandlerFunc) (http.Handler) {
	slog.LogAttrs(*ctx, slog.LevelInfo, "load routes")

	mux    := http.NewServeMux()
	always := alice.New(mw.Recover, mw.Authorise)

	cache := routes.CacheCopy()

	for _, v := range cache {
		mux.Handle(fmt.Sprintf("%v %v", v.HTTPRequestMethod, v.EndpointPath), handlers[v.Handler])
	}

	return always.Then(mux)
}

func InitCache(ctx *context.Context, conn *pgxpool.Conn) error {
	slog.LogAttrs(*ctx, slog.LevelInfo, "initialise web application routes cache")

	const (
		dbSchema = "api_core_rts_api_inf"
		dbFunc   = "rts_inf"
	)

	err := routes.InitCache(ctx, conn, dbSchema, dbFunc)
	if err != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "initialise api route cache",
			slog.String("error" , err.Error()),
		)

		panic("The web application route cache was not initialised correctly")
	}

	return nil
}
