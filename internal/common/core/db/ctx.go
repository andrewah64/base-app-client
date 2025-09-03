package db

import (
	"context"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	Pool *pgxpool.Pool
}

type key int

var poolKey key

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, pool *Pool) context.Context {
	slog.LogAttrs(ctx, slog.LevelDebug, "pool added to context")

	return context.WithValue(ctx, poolKey, pool)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*Pool, bool) {
	slog.LogAttrs(ctx, slog.LevelDebug, "get pool from context")

	p, ok := ctx.Value(poolKey).(*Pool)

	if ok {
		slog.LogAttrs(ctx, slog.LevelDebug, "pool found in context")
	} else {
		slog.LogAttrs(ctx, slog.LevelError, "pool not found in context")
	}

	return p, ok
}
