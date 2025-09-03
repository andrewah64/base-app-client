package val

import (
	"fmt"
	"log/slog"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/password"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
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

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	aurPwd := form.VText (r, "pwd-aur-tnt-mod-aur-pwd")

	infRs, infRsErr := GetInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if infRsErr != nil {
		error.IntSrv(ctx, rw, infRsErr)
		return
	}

	lenPass, symPass, numPass := password.Validate(aurPwd, infRs[0].AupcAurPwdMinLen, infRs[0].AupcAurPwdMaxLen, infRs[0].AupcAurPwdIncSym, infRs[0].AupcAurPwdIncNum)

	data.ResultSet = &map[string]any{"LenPass" : lenPass, "SymPass" : symPass, "NumPass" : numPass}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/pwd/aur/tnt/fragment/val", http.StatusOK, &data)

	return
}
