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
	cm "github.com/andrewah64/base-app-client/internal/common/core/mw"
	   "github.com/andrewah64/base-app-client/internal/common/core/routes"
	wm "github.com/andrewah64/base-app-client/internal/web/core/mw"
)

import (
	"github.com/andrewah64/base-app-client/ui"
)

func Mux(ctx *context.Context, handlers map[string]http.HandlerFunc) (http.Handler) {
	slog.LogAttrs(*ctx, slog.LevelInfo, "load routes")

	mux := http.NewServeMux()

	standard := alice.New(wm.Recover , cm.ResponseHeaders/*, cm.CSRFHandler*/)
	auth     := alice.New(wm.WebAuth)
	unauth   := alice.New(wm.WebUnauth)

	mux.Handle("GET /static/" , http.FileServerFS(ui.Files))

	cache := routes.CacheCopy()

	for _, v := range cache {
		switch v.MiddlewareChain {
			case "web/auth":
				slog.LogAttrs(*ctx, slog.LevelInfo, "register web/auth route",
					slog.String("HTTPRequestMethod", v.HTTPRequestMethod),
					slog.String("EndpointPath"     , v.EndpointPath),
				)

				mux.Handle(fmt.Sprintf("%v %v", v.HTTPRequestMethod, v.EndpointPath), auth.Then(handlers[v.Handler]))
			case "web/unauth":
				slog.LogAttrs(*ctx, slog.LevelInfo, "register web/unauth route",
					slog.String("HTTPRequestMethod", v.HTTPRequestMethod),
					slog.String("EndpointPath"     , v.EndpointPath),
				)

				mux.Handle(fmt.Sprintf("%v %v", v.HTTPRequestMethod, v.EndpointPath), unauth.Then(handlers[v.Handler]))
		}
	}

	return standard.Then(mux)
}

func InitCache(ctx *context.Context, conn *pgxpool.Conn) error {
	slog.LogAttrs(*ctx, slog.LevelInfo, "initialise web application routes cache")

	const (
		dbSchema = "web_core_unauth_rts_web_inf"
		dbFunc   = "rts_inf"
	)

	err := routes.InitCache(ctx, conn, dbSchema, dbFunc)
	if err != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "initialise web application route cache",
			slog.String("error" , err.Error()),
		)

		panic("The web application route cache was not initialised correctly")
	}

	return nil
}
