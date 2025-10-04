package mw

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

import (
	"github.com/andrewah64/base-app-client/internal/api/core/key"
	"github.com/andrewah64/base-app-client/internal/api/core/error"
	ck "github.com/andrewah64/base-app-client/internal/common/core/key"
	"github.com/andrewah64/base-app-client/internal/common/core/log"
	"github.com/andrewah64/base-app-client/internal/common/core/mw/auth"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
)

func Authorise(next http.Handler) http.Handler {
	return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request){
		ctx, ssd, eppPt, hrmNm, _, err := auth.Setup(r)
		if err != nil {
			error.IntSrv(ctx, rw, err)
			return
		}

		defer ssd.Conn.Release()

		rw.Header().Add("Vary", "Authorization")

		authHdr := strings.Split(r.Header.Get("Authorization"), " ")

		slog.LogAttrs(ctx, slog.LevelDebug, "check number of elements in Authorization header",
			slog.Int("len(authHdr)" , len(authHdr)),
		)

		switch len(authHdr) {
			case 2:
				authType  := authHdr[0]
				authToken := authHdr[1]

				slog.LogAttrs(ctx, slog.LevelDebug, "check 'Authorization' header",
					slog.String("authType" , authType),
					slog.String("authToken", authToken),
				)

				if len(authToken) != 26 {
					error.ValErr(ctx, rw, map[string]string{"token" : "must be 26 characters"})
					return
				}

				switch authType {
					case "Key":
						idErr := session.Identity(&ctx, slog.Default(), ssd.Conn, "role_api_core_key_aur_lgn")
						if idErr != nil {
							slog.LogAttrs(ctx, slog.LevelError, "set db connection's user",
								slog.String("error", idErr.Error()),
							)

							panic(idErr)
						}

						rs, rsErr := key.UserInfo(&ctx, slog.Default(), ssd.Conn, ssd.TntId, ck.Hash(authToken), *eppPt, *hrmNm)
						if rsErr != nil{
							error.IntSrv(ctx, rw, rsErr)
							return
						}

						if len(rs) == 1 {
							ssd.Logger = log.Setup(slog.Level(rs[0].LogLevel))

							idErr = session.Identity(&ctx, slog.Default(), ssd.Conn, rs[0].UserRole)
							if idErr != nil {
								error.IntSrv(ctx, rw, idErr)
								return
							}
						} else {
							error.NotAuth(ctx, rw, fmt.Errorf("right header format, no permission"))
							return
						}
					case "Token":
						fmt.Printf("xxxxxxxxxxxxx")
					default:
						error.NotAuth(ctx, rw, fmt.Errorf("wrong header format"))
						return
				}
			default:
				error.NotAuth(ctx, rw, fmt.Errorf("wrong header format"))
				return
		}
		next.ServeHTTP(rw, r)
	})
}
