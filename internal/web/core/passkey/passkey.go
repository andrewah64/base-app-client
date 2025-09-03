package passkey

import (
	"context"
	"fmt"
	"log/slog"
)

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/db"
)

var (
	cache map[int]*webauthn.WebAuthn = make(map[int]*webauthn.WebAuthn)
)

func InitCache(ctx *context.Context, conn *pgxpool.Conn) error {
	slog.LogAttrs(*ctx, slog.LevelInfo, "initialise passkey cache")

	const (
		dbSchema = "all_core_unauth_tnt_all_inf"
		dbFunc   = "tnt_inf"
	)

	type tenant struct {
		TntId     int
		TntPrtc   string
		TntFqdn   string
		TntPort   int
		TntOrigin string
	}

	rs, rsErr := db.DataSet[tenant](ctx, slog.Default(), conn, func(ctx *context.Context, tx *pgx.Tx)(string, string, *pgx.Rows, error){
		qry := fmt.Sprintf("select %v.%v($1)", dbSchema, dbFunc)

		call, err := (*tx).Query(*ctx, qry, dbFunc)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "get tenant data",
				slog.String("error", err.Error()),
			)

			return qry, dbFunc, nil, fmt.Errorf("get tenant data: %w", err)
		}

		return qry, dbFunc, &call, nil
	})

	if rsErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "get tenant info",
			slog.String("error", rsErr.Error()),
		)
		return rsErr
	}

	for _, v := range rs {
		wa, err := webauthn.New(
			&webauthn.Config{
				RPDisplayName : v.TntFqdn,
				RPID          : v.TntFqdn,
				RPOrigins     : []string{v.TntOrigin},
			},
		)
		if err != nil {
			slog.LogAttrs(*ctx, slog.LevelError, "initialise webauthn",
				slog.String("error", err.Error()),
			)

			panic("failed to initiatise webauthn")
		}
		cache[v.TntId] = wa
	}

	slog.LogAttrs(*ctx, slog.LevelInfo, "initialised webauthn",
		slog.Any("cache", cache),
	)

	return nil
}

func WebAuthn (ctx *context.Context, logger *slog.Logger, tntId int) *webauthn.WebAuthn {
	logger.LogAttrs(*ctx, slog.LevelDebug, "get webauthn",
		slog.Int("tntId", tntId),
	)

	if webauthn, ok := cache[tntId]; ok {
		return webauthn
	} else {
		slog.LogAttrs(*ctx, slog.LevelError, "webauthn not found",
			slog.Int("tntId", tntId),
		)
		panic("webauthn not found")
	}
}
