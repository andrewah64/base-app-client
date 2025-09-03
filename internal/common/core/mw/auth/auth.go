package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
        "github.com/andrewah64/base-app-client/internal/common/core/tenant"
)

import (
	"github.com/google/uuid"
)

func Setup(r *http.Request) (context.Context, *session.CtxData, *string, *string, *string, error){
	ssd := &session.CtxData{
		RequestId: uuid.NewString(),
	}

	ctx := session.NewContext(r.Context(), ssd)

	connPool, poolOk := db.FromContext(ctx)
	if ! poolOk {
		return nil, nil, nil, nil, nil, fmt.Errorf("could not acquire connection pool")
	}

	conn, connErr := db.Conn(&ctx, slog.Default(), connPool.Pool)
	if connErr != nil {
		return nil, nil, nil, nil, nil, connErr
	}

	ssd.Conn = conn

	var (
		epp    = strings.Split(r.Pattern, " ")[1]
		origin = tenant.Origin(r)
		hrm    = r.Method
	)

	ssd.TntId = tenant.Tenant(&ctx, slog.Default(), origin)

	slog.LogAttrs(ctx, slog.LevelDebug, "setup Auth middleware",
		slog.String("epp"   , epp),
		slog.String("hrm"   , hrm),
		slog.String("origin", origin),
	)

	return ctx, ssd, &epp, &hrm, &origin, nil
}
