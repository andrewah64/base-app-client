package aur

import (
	"fmt"
	"net/http"
	"log/slog"
)

import (
	cs "github.com/andrewah64/base-app-client/internal/common/core/session"
	   "github.com/andrewah64/base-app-client/internal/web/core/error"
	ws "github.com/andrewah64/base-app-client/internal/web/core/session"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
)

func Delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	ssd, ok := cs.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "get session-scoped data",
		slog.Int("ssd.TntId", ssd.TntId),
	)

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("get request data"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "get page-scoped data",
		slog.Int("data.User.AurId", data.User.AurId),
	)

	ssnTkn, err := r.Cookie("session_token")
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelDebug, "get session token",
			slog.String("error", err.Error()),
		)

		rw.Header().Set("HX-Redirect", "/")

		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "end http session",
		slog.Int("data.User.AurId", data.User.AurId),
	)

	ws.End(&ctx, ssd.Logger, ssd.Conn, rw, ssnTkn)

	rw.Header().Set("HX-Redirect", "/")

	return
}
