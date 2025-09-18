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

	optsRs, optsRsErr := Opts(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if optsRsErr != nil {
		error.IntSrv(ctx, rw, optsRsErr)
		return
	}

	s2cInfRs, s2cInfRsErr := GetS2cInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if s2cInfRsErr != nil {
		error.IntSrv(ctx, rw, s2cInfRsErr)
		return
	}

	data.FormOpts  = &map[string]any{
		"Opts" : &optsRs,
	}

	data.ResultSet = &map[string]any{
		"S2c" : &s2cInfRs,
	}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/s2c/tnt/content", http.StatusOK, data)

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
		case "gen" :
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			s2cEnabled  := form.VBool (r, "s2c-tnt-mod-gen-enabled")
			s2cEntityId := form.VText (r, "s2c-tnt-mod-gen-entity-id")
			aumId       := form.VInt  (r, "s2c-tnt-mod-gen-aum-id")
			uts         := form.VTime (r, "s2c-tnt-mod-gen-uts")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from aupc form",
				slog.Bool  ("s2cEnabled" , s2cEnabled),
				slog.String("s2cEntityId", s2cEntityId),
				slog.Int   ("aumId"      , aumId),
				slog.Any   ("uts"        , uts),
			)

			exptErrs := []string{
				"OLOCK",
			}

			patchErr := PatchS2c(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, s2cEnabled, s2cEntityId, aumId, data.User.AurNm, uts, exptErrs)
			if patchErr != nil {
				var pgErr *pgconn.PgError

				if errors.As(patchErr, &pgErr) {
					switch pgErr.Code {
						case "OLOCK":
							rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-s2c-tnt-mod-gen-form.warning-input-olock-error", "lvl": "error"}}`, currentUrl))

						default:
							slog.LogAttrs(ctx, slog.LevelError, "Patch::unexpected error",
								slog.String("s2cEntityId", s2cEntityId),
								slog.Bool  ("s2cEnabled" , s2cEnabled),
								slog.Int   ("aumId"      , aumId),
								slog.Any   ("uts"        , uts),
							)

							notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-gen-form.warning-input-unexpected-error")}, data)
					}
				}
			}

			s2cUtsInfRs, s2cUtsInfRsErr := GetS2cUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if s2cUtsInfRsErr != nil {
				error.IntSrv(ctx, rw, s2cUtsInfRsErr)
				return
			}

			html.HiddenUtsFragment(rw, "s2c-tnt-mod-gen-uts-ctr", "s2c-tnt-mod-gen-uts", "s2c-tnt-mod-gen-uts", s2cUtsInfRs[0].Uts, data.TFT())

			notification.Show(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-gen-form.message-input-success")} , data)
	/*
		case "cdf":
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			s2cCrtCn := form.VText (r, "s2c-tnt_mod-cdf-crt-cn")
			s2cCrtDn := form.VText (r, "s2c-tnt_mod-cdf-crt-dn")
			uts      := form.VTime (r, "s2c-tnt-mod-cdf-uts")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from aupc form",
				slog.Bool  ("s2cEnabled" , s2cEnabled),
				slog.String("s2cEntityId", s2cEntityId),
				slog.Int   ("aumId"      , aumId),
				slog.Any   ("uts"        , uts),
			)
	*/
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}

