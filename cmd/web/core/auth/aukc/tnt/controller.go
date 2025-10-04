package tnt

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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
		notification.Toast(ctx, slog.Default(), rw, r, p.Get("lvl"), &map[string]string{"Message" : data.T(p.Get("ntf"))} , data)
	}

	aukcInfRs, aukcInfRsErr := GetAukcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aukcInfRsErr != nil {
		error.IntSrv(ctx, rw, aukcInfRsErr)
		return
	}

	optsRs, optsRsErr := Opts(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if optsRsErr != nil {
		error.IntSrv(ctx, rw, optsRsErr)
		return
	}

	pahInfRs, pahInfRsErr := GetPahInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if pahInfRsErr != nil {
		error.IntSrv(ctx, rw, pahInfRsErr)
		return
	}

	pkgInfRs, pkgInfRsErr := GetPkgInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if pkgInfRsErr != nil {
		error.IntSrv(ctx, rw, pkgInfRsErr)
		return
	}

	prhInfRs, prhInfRsErr := GetPrhInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if prhInfRsErr != nil {
		error.IntSrv(ctx, rw, prhInfRsErr)
		return
	}

	data.FormOpts  = &map[string]any{
		"Opts" : &optsRs,
	}

	data.ResultSet = &map[string]any{
		"Aukc" : &aukcInfRs,
		"Pah"  : &pahInfRs,
		"Pkg"  : &pkgInfRs,
		"Prh"  : &prhInfRs,
	}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/aukc/tnt/content", http.StatusOK, data)

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

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	aukcAurNmMinLen  := form.VInt  (r, "aukc-tnt-mod-aur-nm-min-len")
	aukcAurNmMaxLen  := form.VInt  (r, "aukc-tnt-mod-aur-nm-max-len")
	aukcEnabled      := form.VBool (r, "aukc-tnt-mod-aur-pky-enabled")
	pkaId            := form.VInt  (r, "aukc-tnt-mod-pka-id")
	pktId            := form.VInt  (r, "aukc-tnt-mod-pkt-id")
	pdcId            := form.VInt  (r, "aukc-tnt-mod-pdc-id")
	puvRegId         := form.VInt  (r, "aukc-tnt-mod-puv-reg-id")
	puvAtnId         := form.VInt  (r, "aukc-tnt-mod-puv-atn-id")
	uts              := form.VTime (r, "aukc-tnt-mod-uts")

	pkgId, pkgIdErr := form.VIntArray  (r, "aukc-tnt-mod-pkg-id")
	if pkgIdErr != nil {
		error.IntSrv(ctx, rw, pkgIdErr)
		return
	}

	pahId, pahIdErr := form.VIntArray (r, "aukc-tnt-mod-pah-id")
	if pahIdErr != nil {
		error.IntSrv(ctx, rw, pahIdErr)
		return
	}

	prhId, prhIdErr := form.VIntArray (r, "aukc-tnt-mod-prh-id")
	if prhIdErr != nil {
		error.IntSrv(ctx, rw, prhIdErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from aukc form",
		slog.Int  ("aukcAurNmMinLen" , aukcAurNmMinLen),
		slog.Int  ("aukcAurNmMaxLen" , aukcAurNmMaxLen),
		slog.Bool ("aukcEnabled"     , aukcEnabled),
		slog.Int  ("pkaId"           , pkaId),
		slog.Int  ("pktId"           , pktId),
		slog.Int  ("pdcId"           , pktId),
		slog.Int  ("puvRegId"        , puvRegId),
		slog.Int  ("puvAtnId"        , puvAtnId),
		slog.Any  ("pkgId"           , pkgId),
		slog.Any  ("pahId"           , pahId),
		slog.Any  ("prhId"           , prhId),
		slog.Any  ("uts"             , uts),
	)

	exptErrs := []string{
		"OLOCK",
	}

	patchErr := PatchAukc(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aukcAurNmMinLen, aukcAurNmMaxLen, aukcEnabled, pkaId, pktId, pdcId, puvRegId, puvAtnId, pkgId, prhId, pahId, data.User.AurNm, uts, exptErrs)
	if patchErr != nil{
		var pgErr *pgconn.PgError

		if errors.As(patchErr, &pgErr) {
			switch pgErr.Code {
				case "OLOCK":
					currentUrl := r.Header.Get("HX-Current-URL")

					rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-aukc-tnt-mod-form.warning-input-aukc-olock-error", "lvl": "error"}}`, currentUrl))

				default:
					slog.LogAttrs(ctx, slog.LevelError, "Patch::unexpected error",
						slog.String("patchErr.Error()" , patchErr.Error()),
						slog.Int   ("aukcAurNmMinLen"  , aukcAurNmMinLen),
						slog.Int   ("aukcAurNmMaxLen"  , aukcAurNmMaxLen),
						slog.Bool  ("aukcEnabled"      , aukcEnabled),
						slog.Int   ("pkaId"            , pkaId),
						slog.Int   ("pktId"            , pktId),
						slog.Int   ("pdcId"            , pdcId),
						slog.Int   ("puvRegId"         , puvRegId),
						slog.Int   ("puvAtnId"         , puvAtnId),
						slog.Any   ("pkgId"            , pkgId),
						slog.Any   ("prhId"            , prhId),
						slog.Any   ("pahId"            , pahId),
						slog.Any   ("uts"              , uts),
					)

					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-aukc-tnt-mod-form.warning-input-aukc-unexpected-error")}, data)
			}
		}

		return
	}

	aukcUtsInfRs, aukcUtsInfRsErr := GetAukcUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aukcUtsInfRsErr != nil {
		error.IntSrv(ctx, rw, aukcUtsInfRsErr)
		return
	}

	html.HiddenUtsFragment(rw, "aukc-tnt-mod-uts-ctr", "aukc-tnt-mod-uts", "aukc-tnt-mod-uts", aukcUtsInfRs[0].Uts, data.TFT())

	notification.Toast(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-aukc-tnt-mod-form.message-input-success")} , data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}
