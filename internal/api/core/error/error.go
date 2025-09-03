package error

import (
	"context"
	"log/slog"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/api/core/json"
)

func manage(ctx context.Context, rw http.ResponseWriter, status int, err any) {
	var (
		msg string
	)

	switch e := err.(type) {
		case error:
			msg = e.Error()
		case map[string]string :
			msg = "???"
		default:
			msg = "???"
	}

	env := json.Envelope{
		"error" : msg,
	}

	jsErr := json.Write(&ctx, slog.Default(), rw, status, env, nil)
	if jsErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, http.StatusText(status),
			slog.String("error" , jsErr.Error()),
		)
		rw.WriteHeader(status)
	}
}

func BadReq(ctx context.Context, rw http.ResponseWriter, err error) {
	manage(ctx, rw, http.StatusBadRequest, err)
}

func IntSrv(ctx context.Context, rw http.ResponseWriter, err error) {
	manage(ctx, rw, http.StatusInternalServerError, err)
}

func NotAuth(ctx context.Context, rw http.ResponseWriter, err error) {
	manage(ctx, rw, http.StatusUnauthorized, err)
}

func ValErr(ctx context.Context, rw http.ResponseWriter, errors map[string]string) {
	manage(ctx, rw, http.StatusUnprocessableEntity, errors)
}
