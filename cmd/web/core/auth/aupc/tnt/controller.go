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
		notification.Show(ctx, slog.Default(), rw, r, p.Get("lvl"), &map[string]string{"Message" : data.T(p.Get("ntf"))} , data)
	}

	aupcInfRs, aupcInfRsErr := GetAupcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aupcInfRsErr != nil {
		error.IntSrv(ctx, rw, aupcInfRsErr)
		return
	}

	data.ResultSet = &map[string]any{
		"Aupc" : &aupcInfRs,
	}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/aupc/tnt/content", http.StatusOK, data)

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

	aupcAurNmMinLen  := form.VInt  (r, "aupc-tnt-mod-aur-nm-min-len")
	aupcAurNmMaxLen  := form.VInt  (r, "aupc-tnt-mod-aur-nm-max-len")
	aupcAurPwdMinLen := form.VInt  (r, "aupc-tnt-mod-aur-pwd-min-len")
	aupcAurPwdMaxLen := form.VInt  (r, "aupc-tnt-mod-aur-pwd-max-len")
	aupcAurPwdIncSym := form.VBool (r, "aupc-tnt-mod-aur-pwd-inc-sym")
	aupcAurPwdIncNum := form.VBool (r, "aupc-tnt-mod-aur-pwd-inc-num")
	aupcEnabled      := form.VBool (r, "aupc-tnt-mod-aur-pwd-enabled")
	aupcMfaEnabled   := form.VBool (r, "aupc-tnt-mod-aur-pwd-mfa-enabled")
	uts              := form.VTime (r, "aupc-tnt-mod-uts")

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
					currentUrl := r.Header.Get("HX-Current-URL")

					rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-aupc-tnt-mod-form.warning-input-aupc-olock-error", "lvl": "error"}}`, currentUrl))

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

	html.HiddenUtsFragment(rw, "aupc-tnt-mod-uts-ctr", "aupc-tnt-mod-uts", "aupc-tnt-mod-uts", aupcUtsInfRs[0].Uts, data.TFT())

	notification.Show(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-atn-tnt-aupc-tab.message-input-success")} , data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}

