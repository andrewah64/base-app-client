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

	aukcInfRs, aukcInfRsErr := GetAukcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aukcInfRsErr != nil {
		error.IntSrv(ctx, rw, aukcInfRsErr)
		return
	}

	aupcInfRs, aupcInfRsErr := GetAupcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aupcInfRsErr != nil {
		error.IntSrv(ctx, rw, aupcInfRsErr)
		return
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
		"Aupc" : &aupcInfRs,
		"Ocp"  : &ocpInfRs,
		"Occ"  : &occInfRs,
		"Pah"  : &pahInfRs,
		"Pkg"  : &pkgInfRs,
		"Prh"  : &prhInfRs,
	}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/atn/tnt/content", http.StatusOK, data)

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

	currentUrl := r.Header.Get("HX-Current-URL")

	switch r.PathValue("nm") {
		case "aupc" :
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			aupcAurNmMinLen  := form.VInt  (r, "atn-tnt-mod-aupc-aur-nm-min-len")
			aupcAurNmMaxLen  := form.VInt  (r, "atn-tnt-mod-aupc-aur-nm-max-len")
			aupcAurPwdMinLen := form.VInt  (r, "atn-tnt-mod-aupc-aur-pwd-min-len")
			aupcAurPwdMaxLen := form.VInt  (r, "atn-tnt-mod-aupc-aur-pwd-max-len")
			aupcAurPwdIncSym := form.VBool (r, "atn-tnt-mod-aupc-aur-pwd-inc-sym")
			aupcAurPwdIncNum := form.VBool (r, "atn-tnt-mod-aupc-aur-pwd-inc-num")
			aupcEnabled      := form.VBool (r, "atn-tnt-mod-aupc-aur-pwd-enabled")
			aupcMfaEnabled   := form.VBool (r, "atn-tnt-mod-aupc-aur-pwd-mfa-enabled")
			uts              := form.VTime (r, "atn-tnt-mod-aupc-uts")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from aupc form",
				slog.Int  ("aupcAurNmMinLen"  , aupcAurNmMinLen),
				slog.Int  ("aupcAurNmMaxLen"  , aupcAurNmMaxLen),
				slog.Int  ("aupcAurPwdMinLen" , aupcAurPwdMinLen),
				slog.Int  ("aupcAurPwdMaxLen" , aupcAurPwdMaxLen),
				slog.Bool ("aupcAurPwdIncSym" , aupcAurPwdIncSym),
				slog.Bool ("aupcAurPwdIncNum" , aupcAurPwdIncNum),
				slog.Bool ("aupcEnabled"      , aupcEnabled),
				slog.Bool ("aupcMfaEnabled"   , aupcMfaEnabled),
				slog.Any  ("uts"              , uts),
			)

			exptErrs := []string{
				"OLOCK",
			}

			patchErr := PatchAupc(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aupcAurNmMinLen, aupcAurNmMaxLen, aupcAurPwdMinLen, aupcAurPwdMaxLen, aupcAurPwdIncSym, aupcAurPwdIncNum, aupcEnabled, aupcMfaEnabled, data.User.AurNm, uts, exptErrs)
			if patchErr != nil{
				var pgErr *pgconn.PgError

				if errors.As(patchErr, &pgErr) {
					switch pgErr.Code {
						case "OLOCK":
							rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-atn-tnt-aupc-tab.warning-input-aupc-olock-error", "lvl": "error"}}`, currentUrl))

						default:
							slog.LogAttrs(ctx, slog.LevelError, "Patch::unexpected error",
								slog.Int  ("aupcAurNmMinLen"  , aupcAurNmMinLen),
								slog.Int  ("aupcAurNmMaxLen"  , aupcAurNmMaxLen),
								slog.Int  ("aupcAurPwdMinLen" , aupcAurPwdMinLen),
								slog.Int  ("aupcAurPwdMaxLen" , aupcAurPwdMaxLen),
								slog.Bool ("aupcAurPwdIncSym" , aupcAurPwdIncSym),
								slog.Bool ("aupcAurPwdIncNum" , aupcAurPwdIncNum),
								slog.Bool ("aupcEnabled"      , aupcEnabled),
								slog.Bool ("aupcMfaEnabled"   , aupcMfaEnabled),
								slog.Any  ("uts"              , uts),
							)

							notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-aupc-tab.warning-input-aupc-unexpected-error")}, data)

					}
				}

				return
			}

			aupcUtsInfRs, aupcUtsInfRsErr := GetAupcUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if aupcUtsInfRsErr != nil {
				error.IntSrv(ctx, rw, aupcUtsInfRsErr)
				return
			}

			html.HiddenUtsFragment(rw, "atn-tnt-mod-aupc-uts-ctr", "atn-tnt-mod-aupc-uts", "atn-tnt-mod-aupc-uts", aupcUtsInfRs[0].Uts, data.TFT())

			notification.Show(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-aupc-tab.message-input-success")} , data)

		case "aukc" :
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			aukcAurNmMinLen  := form.VInt  (r, "atn-tnt-mod-aukc-aur-nm-min-len")
			aukcAurNmMaxLen  := form.VInt  (r, "atn-tnt-mod-aukc-aur-nm-max-len")
			aukcEnabled      := form.VBool (r, "atn-tnt-mod-aukc-aur-pky-enabled")
			pkaId            := form.VInt  (r, "atn-tnt-mod-aukc-pka-id")
			pktId            := form.VInt  (r, "atn-tnt-mod-aukc-pkt-id")
			pdcId            := form.VInt  (r, "atn-tnt-mod-aukc-pdc-id")
			puvRegId         := form.VInt  (r, "atn-tnt-mod-aukc-puv-reg-id")
			puvAtnId         := form.VInt  (r, "atn-tnt-mod-aukc-puv-atn-id")
			uts              := form.VTime (r, "atn-tnt-mod-aukc-uts")

			pkgId, pkgIdErr := form.VIntArray  (r, "atn-tnt-mod-aukc-pkg-id")
			if pkgIdErr != nil {
				error.IntSrv(ctx, rw, pkgIdErr)
				return
			}

			pahId, pahIdErr := form.VIntArray (r, "atn-tnt-mod-aukc-pah-id")
			if pahIdErr != nil {
				error.IntSrv(ctx, rw, pahIdErr)
				return
			}

			prhId, prhIdErr := form.VIntArray (r, "atn-tnt-mod-aukc-prh-id")
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
							rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-atn-tnt-aukc-tab.warning-input-aukc-olock-error", "lvl": "error"}}`, currentUrl))

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

							notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-aukc-tab.warning-input-aukc-unexpected-error")}, data)
					}
				}

				return
			}

			aukcUtsInfRs, aukcUtsInfRsErr := GetAukcUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if aukcUtsInfRsErr != nil {
				error.IntSrv(ctx, rw, aukcUtsInfRsErr)
				return
			}

			html.HiddenUtsFragment(rw, "atn-tnt-mod-aukc-uts-ctr", "atn-tnt-mod-aukc-uts", "atn-tnt-mod-aukc-uts", aukcUtsInfRs[0].Uts, data.TFT())

			notification.Show(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-aukc-tab.message-input-success")} , data)

		case "oidc" :
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

			occEnabled      := form.VBool (r, fmt.Sprintf("atn-tnt-mod-occ-enabled-%v"       , occId))
			occUrl          := form.VText (r, fmt.Sprintf("atn-tnt-mod-occ-url-%v"           , occId))
			occClientId     := form.VText (r, fmt.Sprintf("atn-tnt-mod-occ-client-id-%v"     , occId))
			occClientSecret := form.VText (r, fmt.Sprintf("atn-tnt-mod-occ-client-secret-%v" , occId))
			uts             := form.VTime (r, fmt.Sprintf("atn-tnt-mod-occ-uts-%v"           , occId))

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
							rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-atn-tnt-oidc-tab.warning-input-occ-olock-error", "lvl": "error"}}`, currentUrl))

						case pgerrcode.UniqueViolation:
							notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-oidc-tab.warning-input-occ-url-taken")}, data)

						default:
							slog.LogAttrs(ctx, slog.LevelError, "unexpected error",
								slog.String("patchErr.Error()" , patchErr.Error()),
								slog.String("pgErr.Code"       , pgErr.Code),
							)

							notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-oidc-tab.warning-input-occ-unexpected-error")}, data)
					}
				}

				return
			}

			occUtsInfRs, occUtsInfRsErr := GetOccUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, occId)
			if occUtsInfRsErr != nil {
				error.IntSrv(ctx, rw, occUtsInfRsErr)
				return
			}

			html.HiddenUtsFragment(rw, fmt.Sprintf("atn-tnt-mod-occ-uts-ctr-%v", occId), fmt.Sprintf("atn-tnt-mod-occ-uts-%v", occId), fmt.Sprintf("atn-tnt-mod-occ-uts-%v", occId), occUtsInfRs[0].Uts, data.TFT())

			notification.Show(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-oidc-tab.message-input-success")} , data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}

