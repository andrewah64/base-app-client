package tnt

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/password"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

import (
	"github.com/andrewah64/base-app-client/cmd/web/core/auth/pwd/aur/tnt/val"
)

func Get (rw http.ResponseWriter, r *http.Request) {
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

	aurId, aurIdErr := strconv.Atoi(r.PathValue("id"))
	if aurIdErr != nil || aurId < 1 {
		http.NotFound(rw, r)
		return
	}

	infRs, infRsErr := GetPwdAurInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId)
	if infRsErr != nil {
		error.IntSrv(ctx, rw, infRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("aurId"      , aurId),
		slog.Int("len(infRs)" , len(infRs)),
	)

	data.ResultSet = &map[string]any{"Inf": &infRs}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/pwd/aur/tnt/content", http.StatusOK, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
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

	aurId, aurIdErr := strconv.Atoi(r.PathValue("id"))
	if aurIdErr != nil || aurId < 1 {
		http.NotFound(rw, r)
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	aurPwd := form.VText (r, "pwd-aur-tnt-mod-aur-pwd")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int("aurId" , aurId),
	)

	aurHshPw, aurHshPwErr := password.Hash(aurPwd)
	if aurHshPwErr != nil {
		error.IntSrv(ctx, rw, aurHshPwErr)
		return
	}

	valInfRs, valInfRsErr := val.GetInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if valInfRsErr != nil {
		error.IntSrv(ctx, rw, valInfRsErr)
		return
	}

	lenPass, symPass, numPass := password.Validate(aurPwd, valInfRs[0].AupcAurPwdMinLen, valInfRs[0].AupcAurPwdMaxLen, valInfRs[0].AupcAurPwdIncSym, valInfRs[0].AupcAurPwdIncNum)

	if lenPass && symPass && numPass {
		patchErr := PatchPwd(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId, aurHshPw, data.User.AurNm, nil)
		if patchErr != nil{
			error.IntSrv(ctx, rw, patchErr)
			return
		}

		notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : data.T("web-core-auth-pwd-aur-tnt-mod-form.message-success")}, data)
	} else {
		notification.Show(ctx, ssd.Logger, rw, r, "error"   , &map[string]string{"Message" : data.T("web-core-auth-pwd-aur-tnt-mod-form.message-error")}, data)
	}

	
	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}
