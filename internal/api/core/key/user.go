package key

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

type User struct {
	UserRole string
	LogLevel int
}

func UserInfo(ctx *context.Context, logger *slog.Logger, conn *pgxpool.Conn, tenant int, key []byte, endPointURL string, hrm string) ([]User, error) {
	const (
		dbSchema = "api_core_key_aur_lgn"
		dbFunc   = "aur_inf"
	)

	results, err := db.DataSet[User](ctx, logger, conn,
		func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
			qry := fmt.Sprintf("select %v.%v($1, $2, $3, $4, $5)", dbSchema, dbFunc)

			c, cErr := (*tx).Query(*ctx, qry, dbFunc, tenant, key, endPointURL, hrm)
			if cErr != nil {

				slog.LogAttrs(*ctx, slog.LevelError, "get dataset",
					slog.String("cErr.Error()" , cErr.Error()),
					slog.String("qry"          , qry),
					slog.Int   ("tenant"       , tenant),
					slog.String("endPointURL"  , endPointURL),
					slog.String("hrm"          , hrm),
				)

				return qry, dbFunc, nil, fmt.Errorf("call database function: %w", cErr)
			}

			return qry, dbFunc, &c, nil
		})

	return results, err
}
