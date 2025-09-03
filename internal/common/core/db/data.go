package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DataSet[T any](ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, dataset func(*context.Context, *pgx.Tx) (string, string, *pgx.Rows, error)) ([]T, error) {
	tx, txErr := conn.Begin(*ctx)
	if txErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "open transaction",
			slog.String("error", txErr.Error()),
		)

		return nil, fmt.Errorf("start transaction: %w", txErr)
	}

	defer tx.Commit(*ctx)

	qry, refcursorName, functionCall, refErr := dataset(ctx, &tx)
	if refErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "call function",
			slog.String("error", refErr.Error()),
		)

		return nil, fmt.Errorf("call database function: %w", refErr)
	}

	(*functionCall).Close()

	logger.LogAttrs(*ctx, slog.LevelDebug, "close function call")

	refcursorQuery := fmt.Sprintf("fetch all in %v", refcursorName)

	logger.LogAttrs(*ctx, slog.LevelDebug, "setup refcursor query",
		slog.String("refcursorQuery", refcursorQuery),
	)

	rows, qryErr := tx.Query(*ctx, refcursorQuery)
	if qryErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
			slog.String("error", qryErr.Error()),
		)

		return nil, fmt.Errorf("get dataset: %w", qryErr)
	}

	data, rowErr := pgx.CollectRows(rows, pgx.RowToStructByPos[T])
	if rowErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "collect results into array",
			slog.String("error"          , rowErr.Error()),
			slog.String("T type"         , fmt.Sprintf("%T", *new(T))),
			slog.String("refcursorQuery" , refcursorQuery),
			slog.String("qry"            , qry),
		)

		return nil, fmt.Errorf("collect dataset into struct: %w", rowErr)
	}

	return data, nil
}

func Sproc (ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, sprocCall string, args pgx.NamedArgs, exptErrs []string)(error){
	tx, txErr := conn.Begin(*ctx)
	if txErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "open transaction",
			slog.String("error", txErr.Error()),
		)

		return fmt.Errorf("start transaction: %w", txErr)
	}

	logger.LogAttrs(*ctx, slog.LevelDebug, "begin transaction",
		slog.String("sprocCall", sprocCall),
	)

	defer tx.Commit(*ctx)

	_, sprocErr := tx.Exec(*ctx, sprocCall, args)
	if sprocErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(sprocErr, &pgErr) && exptErrs != nil && slices.Contains(exptErrs, (*pgErr).Code) {
			logger.LogAttrs(*ctx, slog.LevelDebug, "execute sproc",
				slog.String("error"    , sprocErr.Error()),
				slog.String("sprocCall", sprocCall),
			)
		} else {
			slog.LogAttrs(*ctx, slog.LevelError, "execute sproc",
				slog.String("error"    , sprocErr.Error()),
				slog.String("sprocCall", sprocCall),
			)
		}

		return sprocErr
	}

	return nil
}
