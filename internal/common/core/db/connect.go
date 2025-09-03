package db

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnPool(ctx *context.Context, logger *slog.Logger, host *string, port *int, db *string, user *string, pw *string, sslmode *string, cachesize *int, app *string) (*pgxpool.Pool, error) {
	var (
		cs = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v&statement_cache_capacity=%v&application_name=%v", *user, *pw, *host, *port, *db, *sslmode, *cachesize, *app)
	)

	config, err := pgxpool.ParseConfig(cs)

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	connPool, err := pgxpool.NewWithConfig(*ctx, config)
	if err != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "get new pool", slog.String("error", err.Error()))

		return nil, fmt.Errorf("new pool: %w", err)
	}

	logger.LogAttrs(*ctx, slog.LevelInfo, "acquire connection pool")

	return connPool, nil
}

func Conn(ctx *context.Context, logger *slog.Logger, connPool *pgxpool.Pool) (*pgxpool.Conn, error) {
	conn, err := connPool.Acquire(*ctx)
	if err != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "get new conn", slog.String("error", err.Error()))

		return nil, fmt.Errorf("new conn: %w", err)
	}

	logger.LogAttrs(*ctx, slog.LevelDebug, "acquire connection")

	return conn, nil
}
