package id

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

import (
	"github.com/jackc/pgx/v5/pgconn"
)

func Get(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::start")

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request data"))
		return
	}

	nm := r.PathValue("nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get: switch path",
		slog.String("nm" , nm),
	)

	switch nm {
		case "idp" :
			idpId, idpIdErr := strconv.Atoi(r.PathValue("id"))
			if idpIdErr != nil || idpId < 1 {
				error.IntSrv(ctx, rw, fmt.Errorf("Get::get idpId"))
				return
			}

			idpRs, idpRsErr := GetRowIdpMod(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpId)
			if idpRsErr != nil {
				error.IntSrv(ctx, rw, idpRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("idpId"      , idpId),
				slog.Int("len(idpRs)" , len(idpRs)),
			)

			data.ResultSet = &map[string]any{"Idp": &idpRs}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/s2c/tnt/fragment/modrow-idp", http.StatusCreated, &data)

			if len(idpRs) == 0 {
				notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-idp-form.warning-input-olock-error")}, data)
			}
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")

	return
}

func Patch(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Patch::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::start")

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Patch::get request data"))
		return
	}

	nm := r.PathValue("nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch: switch path",
		slog.String("nm" , nm),
	)

	switch nm {
		case "idp" :
			idpId, idpIdErr := strconv.Atoi(r.PathValue("id"))
			if idpIdErr != nil || idpId < 1 {
				error.IntSrv(ctx, rw, fmt.Errorf("Get::get idpId"))
				return
			}

			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			idpNm      := form.VText (r, fmt.Sprintf("s2c-tnt-inf-idp-nm-%v"      , idpId))
			idpEnabled := form.VBool (r, fmt.Sprintf("s2c-tnt-inf-idp-enabled-%v" , idpId))
			uts        := form.VTime (r, fmt.Sprintf("s2c-tnt-inf-idp-uts-%v"     , idpId))

			idpValRs, idpValRsErr := GetRowIdpVal(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpId, idpEnabled, idpNm)
			if idpValRsErr != nil {
				error.IntSrv(ctx, rw, idpValRsErr)
				return
			}

			if ! idpValRs[0].IdpEnabledOk {
				Get(rw, r)

				notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-idp-form.warning-input-idp-enabled")}, data)

				return
			}

			if ! idpValRs[0].IdpNmOk {
				Get(rw, r)

				notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-idp-form.warning-input-idp-nm-taken", "idpNm" , idpNm)}, data)

				return
			}

			exptErrs := []string{
				"OLOKU",
				"OLOKD",
			}

			patchErr := PatchIdp (&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpId, idpNm, idpEnabled, data.User.AurNm, uts, exptErrs)
			if patchErr != nil{
				Get(rw, r)

				var pgErr *pgconn.PgError

				if errors.As(patchErr, &pgErr) {
					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch: PatchGrp params",
						slog.Int   ("ssd.TntId"  , ssd.TntId),
						slog.Int   ("idpId"      , idpId),
						slog.String("idpNm"      , idpNm),
						slog.Bool  ("idpEnabled" , idpEnabled),
						slog.String("patchErr"   , patchErr.Error()),
					)

					switch pgErr.Code {
						case "OLOKU":
							notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-idp-form.warning-input-olock-error")}, data)

						case "OLOKD":
							//intentionally empty

						default:
							notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-idp-form.warning-input-unexpected-error")}, data)
					}

					return
				}
			}

			idpRs, idpRsErr := GetRowIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpId)
			if idpRsErr != nil {
				error.IntSrv(ctx, rw, idpRsErr)
				return
			}

			data.ResultSet = &map[string]any{"Idp": &idpRs}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/s2c/tnt/fragment/infrow-idp", http.StatusCreated, &data)

			notification.Toast(ctx, slog.Default(), rw, r, "success" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-idp-form.message-input-success")}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")

	return
}
