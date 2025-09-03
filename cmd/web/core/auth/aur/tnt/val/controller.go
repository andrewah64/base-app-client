package val

import (
	"fmt"
	"log/slog"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/validator"
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

	v     := validator.New()
	aurNm := form.VText (r, "aur-tnt-reg-aur-nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get result of validation",
		slog.Int   ("ssd.TntId" , ssd.TntId),
		slog.String("aurNm"     , aurNm),
	)

	valRs, valRsErr := GetInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
	if valRsErr != nil {
		error.IntSrv(ctx, rw, valRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("len(valRs)" , len(valRs)),
	)

	switch len(valRs) {
		case 1:
			if ! valRs[0].AurNmOk {
				v.AddError("aur-tnt-reg-credentials", data.T("web-core-auth-aur-tnt-reg-form.warning-input-credentials-taken", "aurNm", aurNm))
			}
		default:
			v.AddError("aur-tnt-reg-unexpected", data.T("web-core-auth-aur-tnt-reg-form.warning-input-unexpected-error"))
	}

	data.ResultSet = &map[string]any{"Validator" : &v}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/aur/tnt/fragment/val", http.StatusUnprocessableEntity, &data)
	return
}
