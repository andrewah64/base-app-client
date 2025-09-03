package session

import (
	"context"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type CtxData struct {
	RequestId    string
	TntId     int
	Conn        *pgxpool.Conn
	Logger      *slog.Logger
}

type key int

var ctxDataKey key

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, ctxData *CtxData) context.Context {
	return context.WithValue(ctx, ctxDataKey, ctxData)
}

// FromContext returns the value stored in ctx, if any.
func FromContext(ctx context.Context) (*CtxData, bool) {
	p, ok := ctx.Value(ctxDataKey).(*CtxData)
	return p, ok
}
