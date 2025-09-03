package tnt

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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

func Delete (rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Delete::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::start")

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Delete::get request data"))
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr) 
		return
	}

	ssnTkn := form.VTextArray(r, "ssn-tnt-inf-wauhs-ssn-tk")

	if len(ssnTkn) > 0 {
		delErr := DelSsn(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, ssnTkn, data.User.AurNm, nil)
		if delErr != nil{
			error.IntSrv(ctx, rw, delErr)
			return
		}

		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::success")

		rw.Header().Set("HX-Trigger", "mod")

		message := ""

		if len(ssnTkn) == 1 {
			message = data.T("web-core-auth-ssn-tnt-del-form.message-delete-success-singular", "n", strconv.Itoa(len(ssnTkn)))
		} else {
			message = data.T("web-core-auth-ssn-tnt-del-form.message-delete-success-plural"  , "n", strconv.Itoa(len(ssnTkn)))
		}

		notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::end")

	return
}

func params(aurNm string, pageNumber int) string {
	v := url.Values{}

	v.Set("ssn-tnt-inf-aur-nm"     , aurNm)

	v.Set("ssn-tnt-inf-page-number", strconv.Itoa(pageNumber))

	return v.Encode()
}

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

	pageNumber  := 2
	offset      := 0
	resultLimit := 50
	trigger     := r.Header.Get("HX-Trigger")

	switch trigger {
		case "": // page load
			ssnRs, ssnRsErr := GetSsn(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, "", offset, resultLimit)
			if ssnRsErr != nil {
				error.IntSrv(ctx, rw, ssnRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve ssn dataset",
				slog.Int("len(ssnRs)", len(ssnRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &ssnRs,
				"PageNumber"  : pageNumber,
				"ResultLimit" : resultLimit,
				"Params"      : params("", pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/ssn/tnt/content", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [page load]")

		case "ssn-tnt-inf-scr": // infinite scroll
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			aurNm      := form.VText (r, "ssn-tnt-inf-aur-nm")
			pageNumber := form.VInt  (r, "ssn-tnt-inf-page-number")
			offset     := (pageNumber - 1) * resultLimit

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("aurNm"      , aurNm),
				slog.Int   ("pageNumber" , pageNumber),
				slog.Int   ("offset"     , offset),
			)

			ssnRs, ssnRsErr := GetSsn(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, offset, resultLimit)
			if ssnRsErr != nil {
				error.IntSrv(ctx, rw, ssnRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(ssnRs)" , len(ssnRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &ssnRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(aurNm, pageNumber + 1),
			}

			rw.Header().Set("HX-Trigger", "inf")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/ssn/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [infinite scroll]")

		case "ssn-tnt-inf-form":
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			aurNm      := form.VText (r, "ssn-tnt-inf-aur-nm")
			pageNumber := form.VInt  (r, "ssn-tnt-inf-page-number")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("aurNm"      , aurNm),
				slog.Int   ("pageNumber" , pageNumber),
			)

			ssnRs, ssnRsErr := GetSsn(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, offset, resultLimit)
			if ssnRsErr != nil {
				error.IntSrv(ctx, rw, ssnRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(ssnRs)" , len(ssnRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &ssnRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(aurNm, pageNumber + 1),
			}

			rw.Header().Set("HX-Trigger", "src")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/ssn/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [search]")
	}
}
