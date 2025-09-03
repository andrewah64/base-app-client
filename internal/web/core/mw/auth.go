package mware

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
)

import (
	   "github.com/andrewah64/base-app-client/internal/common/core/i18n"
	   "github.com/andrewah64/base-app-client/internal/common/core/log"
	   "github.com/andrewah64/base-app-client/internal/common/core/mw/auth"
	   "github.com/andrewah64/base-app-client/internal/common/core/routes"
	cs "github.com/andrewah64/base-app-client/internal/common/core/session"
	   "github.com/andrewah64/base-app-client/internal/web/core/error"
	ws "github.com/andrewah64/base-app-client/internal/web/core/session"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
)

import (
	"github.com/justinas/nosurf"
)

func WebAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		ctx, ssd, eppPt, hrmNm, origin, err := auth.Setup(r)
		if err != nil {
			error.IntSrv(ctx, rw, err)
			return
		}

		defer ssd.Conn.Release()

		slog.LogAttrs(ctx, slog.LevelDebug, "setup middleware",
			slog.String("eppPt"  , *eppPt),
			slog.String("hrmNm"  , *hrmNm),
			slog.String("origin" , *origin),
		)

		var (
			ssnTkn, _ = r.Cookie("session_token")
		)

		if ssnTkn != nil {

			slog.LogAttrs(ctx, slog.LevelDebug, "validate http session info")

			idErr := cs.Identity(&ctx, slog.Default(), ssd.Conn, "role_web_core_auth_ssn_aur_inf")
			if idErr != nil {
				error.IntSrv(ctx, rw, idErr)
				return
			}

			rs, rsErr := ws.AuthSessionUserInfo(&ctx, slog.Default(), ssd.Conn, ssd.TntId, ssnTkn.Value, *eppPt, *hrmNm)
			if rsErr != nil{
				error.IntSrv(ctx, rw, rsErr)
				return
			}

			switch len(rs){
				case 0:
					slog.LogAttrs(ctx, slog.LevelDebug, "user details not found",
						slog.String("ssnTkn" , ssnTkn.Value),
						slog.String("eppPt"  , *eppPt),
						slog.String("hrmNm"  , *hrmNm),
						slog.String("origin" , *origin),
					)

					idErr := cs.Identity(&ctx, slog.Default(), ssd.Conn, "role_web_core_auth_ssn_aur_end")
					if idErr != nil {
						error.IntSrv(ctx, rw, idErr)
						return
					}

					endErr := ws.End(&ctx, slog.Default(), ssd.Conn, rw, ssnTkn)
					if endErr != nil {
						error.IntSrv(ctx, rw, idErr)
						return
					}

					rw.Header().Set("HX-Redirect", "/")

					return
				case 1:
					ssd.Logger = log.Setup(slog.Level(rs[0].LvlNb))

					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "add user details to context",
						slog.Int   ("rs[0].AurId"                   , rs[0].AurId),
						slog.String("rs[0].AurNm"                   , rs[0].AurNm),
						slog.String("rs[0].LngCd"                   , rs[0].LngCd),
						slog.String("rs[0].RolName"                 , rs[0].RolName),
						slog.String(`fmt.Sprintf("%", rs[0].Roles)` , fmt.Sprintf("%", rs[0].Roles)),
						slog.Int   ("rs[0].LvlNb"                   , rs[0].LvlNb),
						slog.String("rs[0].EppPt"                   , rs[0].EppPt),
						slog.String("rs[0].HrmNm"                   , rs[0].HrmNm),
					)

					idErr := cs.Identity(&ctx, ssd.Logger, ssd.Conn, rs[0].RolName)
					if idErr != nil {
						error.IntSrv(ctx, rw, idErr)
						return
					}

					rt, rtErr := routes.EndpointRoute(&ctx, ssd.Logger, routes.Key(*hrmNm, *eppPt))
					if rtErr != nil {
						error.IntSrv(ctx, rw, rsErr)
						return
					}

					a := false

					switch {
						case len(rt.Role) == 0:
							a = true
						case len(rt.Role) > 0:
							for _, v := range rt.Role {
								if slices.Contains(rs[0].Roles, *v) {
									a = true
									break;
								}
							}
					}

					if a {
						ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "authorise access",
							slog.Bool  ("rt.Role == nil", rt.Role == nil),
							slog.Int   ("UserId"        , rs[0].AurId),
							slog.String("eppPt"         , *eppPt),
							slog.String("hrmNm"         , *hrmNm),
						)

						next.ServeHTTP(rw, r.WithContext(
							page.NewContext(
								cs.NewContext(ctx, ssd),
									&page.Data{
										CSRFToken : nosurf.Token(r),
										User      : &rs[0],
										Localiser : i18n.Localiser(ctx, ssd.Logger, rs[0].LngCd),
									},
								),
							),
						)
					} else {
						ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "deny access",
							slog.Int   ("rs[0].AurId" , rs[0].AurId),
							slog.String("eppPt"       , *eppPt),
							slog.String("hrmNm"       , *hrmNm),
							slog.Any   ("rt.Role"     , rt.Role),
						)

						rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content"}`, rs[0].EppPt))

						return
					}
				default:
					error.IntSrv(ctx, rw, fmt.Errorf("details of more than one user found"))
			}
		} else {
			slog.LogAttrs(ctx, slog.LevelDebug, "no HTTP session is in effect")

			http.Redirect(rw, r, "/", http.StatusSeeOther)

			return
		}
	})
}
