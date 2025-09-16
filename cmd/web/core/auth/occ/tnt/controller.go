package tnt

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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Get(rw http.ResponseWriter, r *http.Request){
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

	p := r.URL.Query()

	if p.Has("ntf") && p.Has("lvl"){
		notification.Show(ctx, slog.Default(), rw, r, p.Get("lvl"), &map[string]string{"Message" : data.T(p.Get("ntf"))} , data)
	}

	ocpInfRs, ocpInfRsErr := GetOcpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if ocpInfRsErr != nil {
		error.IntSrv(ctx, rw, ocpInfRsErr)
		return
	}

	occInfRs, occInfRsErr := GetOccInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if occInfRsErr != nil {
		error.IntSrv(ctx, rw, occInfRsErr)
		return
	}

	data.ResultSet = &map[string]any{
		"Ocp"  : &ocpInfRs,
		"Occ"  : &occInfRs,
	}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/occ/tnt/content", http.StatusOK, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
}

func Patch(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Patch::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::start")

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request data"))
		return
	}

	occId, occIdErr := strconv.Atoi(r.PathValue("id"))
	if occIdErr != nil || occId < 1 {
		http.NotFound(rw, r)
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	occEnabled      := form.VBool (r, fmt.Sprintf("occ-tnt-mod-enabled-%v"       , occId))
	occUrl          := form.VText (r, fmt.Sprintf("occ-tnt-mod-client-url-%v"    , occId))
	occClientId     := form.VText (r, fmt.Sprintf("occ-tnt-mod-client-id-%v"     , occId))
	occClientSecret := form.VText (r, fmt.Sprintf("occ-tnt-mod-client-secret-%v" , occId))
	uts             := form.VTime (r, fmt.Sprintf("occ-tnt-mod-uts-%v"           , occId))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from occ form",
		slog.Int   ("occId"           , occId),
		slog.Bool  ("occEnabled"      , occEnabled),
		slog.String("occUrl"          , occUrl),
		slog.String("occClientId"     , occClientId),
		slog.String("occClientSecret" , occClientSecret),
		slog.Any   ("uts"             , uts),
	)

	exptErrs := []string{
		"OLOCK",
		pgerrcode.UniqueViolation,
	}

	patchErr := PatchOcc(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, occId, occEnabled, occUrl, occClientId, occClientSecret, data.User.AurNm, uts, exptErrs)
	if patchErr != nil{
		var pgErr *pgconn.PgError

		if errors.As(patchErr, &pgErr) {
			switch pgErr.Code {
				case "OLOCK":
					currentUrl := r.Header.Get("HX-Current-URL")

					rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-occ-tnt-mod-form.warning-input-occ-olock-error", "lvl": "error"}}`, currentUrl))

				case pgerrcode.UniqueViolation:
					notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-occ-tnt-mod-form.warning-input-occ-url-taken")}, data)

				default:
					slog.LogAttrs(ctx, slog.LevelError, "unexpected error",
						slog.String("patchErr.Error()" , patchErr.Error()),
						slog.String("pgErr.Code"       , pgErr.Code),
					)

					notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-occ-tnt-mod-form.warning-input-occ-unexpected-error")}, data)
			}
		}

		return
	}

	occUtsInfRs, occUtsInfRsErr := GetOccUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, occId)
	if occUtsInfRsErr != nil {
		error.IntSrv(ctx, rw, occUtsInfRsErr)
		return
	}

	html.HiddenUtsFragment(rw, fmt.Sprintf("occ-tnt-mod-uts-ctr-%v", occId), fmt.Sprintf("occ-tnt-mod-uts-%v", occId), fmt.Sprintf("occ-tnt-mod-uts-%v", occId), occUtsInfRs[0].Uts, data.TFT())

	notification.Show(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-occ-tnt-mod-form.message-input-success")} , data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}
