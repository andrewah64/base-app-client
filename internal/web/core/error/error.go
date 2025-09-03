package error

import (
	"context"
	"log/slog"
	"net/http"
)

func IntSrv(ctx context.Context, rw http.ResponseWriter, err error){
	slog.LogAttrs(ctx, slog.LevelError, "Unexpected error",
		slog.String("error" , err.Error()),
	)

	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
