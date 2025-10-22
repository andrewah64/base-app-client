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

func WebUnauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		ctx, ssd, eppPt, hrmNm, origin, err := auth.Setup(r)
		if err != nil {
			error.IntSrv(ctx, rw, err)
			return
		}

		defer ssd.Conn.Release()

		idErr := cs.Identity(&ctx, slog.Default(), ssd.Conn, "role_web_core_unauth_ssn_ep_inf")
		if idErr != nil {
			error.IntSrv(ctx, rw, idErr)
			return
		}

		uasuiRs, uasuiErr := ws.UnauthSessionEndpointInfo(&ctx, slog.Default(), ssd.Conn, ssd.TntId, *eppPt, *hrmNm)
		if uasuiErr != nil{
			error.IntSrv(ctx, rw, uasuiErr)
			return
		}

		ssd.Logger = log.Setup(slog.Level(uasuiRs[0].LvlNb))

		var (
			ssnTkn, _ = r.Cookie("session_token")
		)

		if ssnTkn == nil {
			slog.LogAttrs(ctx, slog.LevelDebug, "no HTTP session is in effect")

			lngCd := r.Header.Get("Accept-Language")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "get language",
				slog.String("accept-language", lngCd),
			)

			next.ServeHTTP(rw, r.WithContext(
				page.NewContext(
					ctx,
					&page.Data{
						CSRFToken : nosurf.Token(r),
						Localiser : i18n.Localiser(ctx, ssd.Logger, lngCd),
					},
				),
			))
		} else {
			idErr := cs.Identity(&ctx, slog.Default(), ssd.Conn, "role_web_core_auth_ssn_aur_inf")
			if idErr != nil {
				error.IntSrv(ctx, rw, idErr)
				return
			}

			asuiRs, asuiErr := ws.AuthSessionUserInfo(&ctx, slog.Default(), ssd.Conn, ssd.TntId, ssnTkn.Value, *eppPt, *hrmNm)
			if asuiErr != nil{
				error.IntSrv(ctx, rw, asuiErr)
				return
			}

			switch len(asuiRs){
				case 1:
					ssd.Logger = log.Setup(slog.Level(asuiRs[0].LvlNb))

					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "found active session",
						slog.Int   ("asuiRs[0].AurId"                   , asuiRs[0].AurId),
						slog.String("asuiRs[0].AurNm"                   , asuiRs[0].AurNm),
						slog.String("asuiRs[0].LngCd"                   , asuiRs[0].LngCd),
						slog.String("asuiRs[0].RolName"                 , asuiRs[0].RolName),
						slog.String(`fmt.Sprintf("%", asuiRs[0].Roles)` , fmt.Sprintf("%", asuiRs[0].Roles)),
						slog.Int   ("asuiRs[0].LvlNb"                   , asuiRs[0].LvlNb),
						slog.String("asuiRs[0].EppPt"                   , asuiRs[0].EppPt),
						slog.String("asuiRs[0].HrmNm"                   , asuiRs[0].HrmNm),
					)

					idErr := cs.Identity(&ctx, ssd.Logger, ssd.Conn, asuiRs[0].RolName)
					if idErr != nil {
						error.IntSrv(ctx, rw, idErr)
						return
					}

					rt, rtErr := routes.EndpointRoute(&ctx, ssd.Logger, routes.Key(asuiRs[0].HrmNm, asuiRs[0].EppPt))
					if rtErr != nil {
						error.IntSrv(ctx, rw, rtErr)
						return
					}

					a := false

					switch {
						case len(rt.Role) == 0:
							a = true
						case len(rt.Role) > 0:
							for _, v := range rt.Role {
								if slices.Contains(asuiRs[0].Roles, *v) {
									a = true
									break;
								}
							}
					}

					if a {
						ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "permit access",
							slog.Bool  ("rt.Role == nil"  , rt.Role == nil),
							slog.Int   ("asuiRs[0].AurId" , asuiRs[0].AurId),
							slog.String("eppPt"           , *eppPt),
							slog.String("hrmNm"           , *hrmNm),
						)

						http.Redirect(rw, r, asuiRs[0].EppPt, http.StatusSeeOther)

						return
					}

					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "deny access",
						slog.Int   ("asuiRs[0].AurId" , asuiRs[0].AurId),
						slog.String("eppPt"           , *eppPt),
						slog.String("hrmNm"           , *hrmNm),
						slog.Any   ("rt.Role"         , rt.Role),
					)

					fallthrough
				case 0:
					slog.LogAttrs(ctx, slog.LevelDebug, "user details not found",
						slog.String("ssnTkn.Value" , ssnTkn.Value),
						slog.String("eppPt"        , *eppPt),
						slog.String("hrmNm"        , *hrmNm),
						slog.String("origin"       , *origin),
					)

					idErr := cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_auth_ssn_aur_end")
					if idErr != nil {
						error.IntSrv(ctx, rw, idErr)
						return
					}

					endErr := ws.End(&ctx, ssd.Logger, ssd.Conn, rw, ssnTkn)
					if endErr != nil {
						error.IntSrv(ctx, rw, idErr)
						return
					}

					rw.Header().Set("HX-Redirect", "/")

					return
				default:
					error.IntSrv(ctx, rw, fmt.Errorf("details of more than one user found"))
			}
		}
	})
}
