package session

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func Identity(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, un string) error {
	idCall := `"` + strings.ReplaceAll(un, `"`, `""`) + `"`

	_, err := conn.Exec(*ctx, fmt.Sprintf("set role %v", idCall))

	if err != nil {
		logger.LogAttrs(*ctx, slog.LevelError, "set identity",
			slog.String("error"  , err.Error()),
			slog.String("un"     , un),
			slog.String("idCall" , idCall),
		)
		return err
	}

	return nil
}
