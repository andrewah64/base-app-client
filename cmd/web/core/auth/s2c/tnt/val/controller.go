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

	switch r.PathValue("nm") {
		case "mde" :
			idpNm := form.VText (r, "s2c-tnt-reg-mde-idp-nm")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get result of validation",
				slog.Int   ("ssd.TntId" , ssd.TntId),
				slog.String("idpNm"     , idpNm),
			)

			valRs, valRsErr := GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm)
			if valRsErr != nil {
				error.IntSrv(ctx, rw, valRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(valRs)" , len(valRs)),
			)

			msgs := make([]string, 0)

			if ! valRs[0].IdpNmOk {
				msgs = append(msgs, data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-idp-nm-taken", "idpNm", idpNm))
			}


			notification.Vrl(ctx, ssd.Logger, rw, r,
				data.T("web-core-auth-s2c-tnt-page.title-edit"),
				data.T("web-core-auth-s2c-tnt-reg-mde-form.title-warning-singular", "n", strconv.Itoa(len(valRs))),
				data.T("web-core-auth-s2c-tnt-reg-mde-form.title-warning-plural"  , "n", strconv.Itoa(len(valRs))),
				&msgs,
				data,
			)

		case "xml" :
			idpNm := form.VText (r, "s2c-tnt-reg-xml-idp-nm")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get result of validation",
				slog.Int   ("ssd.TntId" , ssd.TntId),
				slog.String("idpNm"     , idpNm),
			)

			valRs, valRsErr := GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm)
			if valRsErr != nil {
				error.IntSrv(ctx, rw, valRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(valRs)" , len(valRs)),
			)

			msgs := make([]string, 0)

			if ! valRs[0].IdpNmOk {
				msgs = append(msgs, data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-idp-nm-taken", "idpNm", idpNm))
			}

			notification.Vrl(ctx, ssd.Logger, rw, r,
				data.T("web-core-auth-s2c-tnt-page.title-edit"),
				data.T("web-core-auth-s2c-tnt-reg-xml-form.title-warning-singular", "n", strconv.Itoa(len(valRs))),
				data.T("web-core-auth-s2c-tnt-reg-xml-form.title-warning-plural"  , "n", strconv.Itoa(len(valRs))),
				&msgs,
				data,
			)

		case "spc" :
			spcNm := form.VText (r, "s2c-tnt-reg-spc-nm")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get result of validation",
				slog.Int   ("ssd.TntId" , ssd.TntId),
				slog.String("spcNm"     , spcNm),
			)

			valRs, valRsErr := GetSpcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, spcNm)
			if valRsErr != nil {
				error.IntSrv(ctx, rw, valRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(valRs)" , len(valRs)),
			)

			msgs := make([]string, 0)

			if !valRs[0].SpcNmOk {
				msgs = append(msgs, data.T("web-core-auth-s2c-tnt-reg-spc-form.warning-input-spc-nm-taken", "spcNm", spcNm))
			}

			notification.Vrl(ctx, ssd.Logger, rw, r,
				data.T("web-core-auth-s2c-tnt-page.title-edit"),
				data.T("web-core-auth-s2c-tnt-reg-spc-form.title-warning-singular", "n", strconv.Itoa(len(valRs))),
				data.T("web-core-auth-s2c-tnt-reg-spc-form.title-warning-plural"  , "n", strconv.Itoa(len(valRs))),
				&msgs,
				data,
			)
	}
}
