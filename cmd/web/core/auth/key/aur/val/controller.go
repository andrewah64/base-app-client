package val

import (
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
	"github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
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

	aaukNm := form.VText (r, "key-aur-reg-aauk-nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get result of validation",
		slog.Int   ("ssd.TntId" , ssd.TntId),
		slog.String("aaukNm"    , aaukNm),
	)

	valRs, valRsErr := GetInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukNm)
	if valRsErr != nil {
		error.IntSrv(ctx, rw, valRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("len(valRs)" , len(valRs)),
	)

	msgs := make([]string, 0)

	if ! valRs[0].AaukNmOk {
		msgs = append(msgs, data.T("web-core-auth-key-aur-reg-form.warning-input-aauk-nm-taken", "aaukNm", aaukNm))
	}

	notification.Vrl(ctx, ssd.Logger, rw, r,
		data.T("web-core-auth-key-aur-page.title"),
		data.T("web-core-auth-key-aur-reg-form.title-warning-singular", "n", strconv.Itoa(len(msgs))),
		data.T("web-core-auth-key-aur-reg-form.title-warning-plural"  , "n", strconv.Itoa(len(msgs))),
		&msgs,
		data,
	)
}
