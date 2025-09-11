package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

import (
	"github.com/andrewah64/base-app-client/internal/api/core/route"
	"github.com/andrewah64/base-app-client/internal/api/core/ui/i18n"
	"github.com/andrewah64/base-app-client/internal/common/core/db"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/startup"
)

import (
	"github.com/andrewah64/base-app-client/cmd/api/core/auth/users/register"
	"github.com/andrewah64/base-app-client/cmd/api/core/unauth/health"
)

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

func main(){
	rtp := startup.GetRuntimeParams()

	ctx := session.NewContext(context.Background(), &session.CtxData{
		RequestId: uuid.NewString(),
	})

	startup.SetupDefaultLogger(*rtp.LogLvl)

	pool := startup.SetupPGConnectionPool(ctx, rtp)

	defer pool.Close()

	i18nCacheErr := i18n.InitCache(ctx, language.English)
	if i18nCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the bundle cache",
			slog.String("error", i18nCacheErr.Error()),
		)

		panic(i18nCacheErr)
	}

	conn, connErr := db.Conn(&ctx, slog.Default(), pool)
	if connErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the tenant cache",
			slog.String("error", connErr.Error()),
		)

		panic(connErr)
	}

	defer conn.Release()

	startup.SetupTenantCache(ctx, conn)

	rtsIdErr := session.Identity(&ctx, slog.Default(), conn, "role_api_core_rts_api_inf")
	if rtsIdErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the route cache",
			slog.String("error", rtsIdErr.Error()),
		)

		panic(connErr)
	}

	rtsCacheErr := route.InitCache(&ctx, conn)
	if rtsCacheErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "initialise the route cache",
			slog.String("error", rtsCacheErr.Error()),
		)

		panic(rtsCacheErr)
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion      : tls.VersionTLS13,
	}

	var (
		handlers = map[string]http.HandlerFunc{
			"api.core.auth.aur.tnt.reg.Register" : register.Register,
			"api.core.unauth.health.Check"       : health.Check,
		}
	)

	server := &http.Server{
		Addr        :	fmt.Sprintf(":%d", *rtp.HttpPort),
		Handler     :	route.Mux(&ctx, handlers),
		BaseContext :	func(_ net.Listener) context.Context {
					return db.NewContext(
						context.Background(),
						&db.Pool {
							Pool: pool,
						},
					)
				},
		TLSConfig   :	tlsConfig,
		IdleTimeout :	time.Minute,
		ReadTimeout :	5  * time.Second,
		WriteTimeout:	10 * time.Second,
	}

	srvErr := server.ListenAndServeTLS("cert.pem", "key.pem")

	slog.LogAttrs(ctx, slog.LevelError, "server error",
		slog.String("error", srvErr.Error()),
	)

	os.Exit(1)
}
