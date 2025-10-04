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

func params(eppPt string, hrmId *int64, lvlId *int64, pageNumber int) string {
	v := url.Values{}

	v.Set("log-ep-tnt-inf-epp-pt" , eppPt)

	switch hrmId {
		case nil:
			v.Set("log-ep-tnt-inf-hrm-id" , "")
		default :
			v.Set("log-ep-tnt-inf-hrm-id" , strconv.FormatInt(*hrmId, 10))
	}

	switch lvlId {
		case nil:
			v.Set("log-ep-tnt-inf-lvl-id" , "")
		default :
			v.Set("log-ep-tnt-inf-lvl-id" , strconv.FormatInt(*lvlId, 10))
	}

	v.Set("log-ep-tnt-inf-page-number" , strconv.Itoa(pageNumber))

	return v.Encode()
}

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

	pageNumber  := 2
	offset      := 0
	resultLimit := 50
	trigger     := r.Header.Get("HX-Trigger")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int   ("pageNumber"  , pageNumber),
		slog.Int   ("offset"      , offset),
		slog.Int   ("resultLimit" , resultLimit),
		slog.String("trigger"     , trigger),
	)

	switch trigger {
		case "": // page load
			optsRs, optsRsErr := Opts(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if optsRsErr != nil {
				error.IntSrv(ctx, rw, optsRsErr)
				return
			}

			logRs, logRsErr := GetLog(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, "", nil, nil, offset, resultLimit)
			if logRsErr != nil {
				error.IntSrv(ctx, rw, logRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(*optsRs)", len(*optsRs)),
				slog.Int("len(logRs)"  , len(logRs)),
			)

			data.FormOpts  = &map[string]any{"Search": &optsRs}
			data.ResultSet = &map[string]any{
				"Search"      : &logRs,
				"PageNumber"  : pageNumber,
				"ResultLimit" : resultLimit,
				"Params"      : params("", nil, nil, pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/log/ep/tnt/content", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [page load]")
		case "log-ep-tnt-inf-scr": // infinite scroll
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			eppPt      := form.VText (r, "log-ep-tnt-inf-epp-pt")
			hrmId      := form.PInt64(r, "log-ep-tnt-inf-hrm-id")
			lvlId      := form.PInt64(r, "log-ep-tnt-inf-lvl-id")
			pageNumber := form.VInt  (r, "log-ep-tnt-inf-page-number")
			offset     := (pageNumber - 1) * resultLimit

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("eppPt"      , eppPt),
				slog.Any   ("hrmId"      , hrmId),
				slog.Any   ("lvlId"      , lvlId),
				slog.Any   ("pageNumber" , pageNumber),
			)

			logRs, logRsErr := GetLog(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, eppPt, hrmId, lvlId, offset, resultLimit)
			if logRsErr != nil {
				error.IntSrv(ctx, rw, logRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(logRs)" , len(logRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &logRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(eppPt, hrmId, lvlId, pageNumber + 1),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/log/ep/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [infinite scroll]")
		case "log-ep-tnt-inf-form": // search
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			eppPt      := form.VText (r, "log-ep-tnt-inf-epp-pt")
			hrmId      := form.PInt64(r, "log-ep-tnt-inf-hrm-id")
			lvlId      := form.PInt64(r, "log-ep-tnt-inf-lvl-id")
			pageNumber := form.VInt  (r, "log-ep-tnt-inf-page-number")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("eppPt"      , eppPt),
				slog.Any   ("hrmId"      , hrmId),
				slog.Any   ("lvlId"      , lvlId),
				slog.Any   ("pageNumber" , pageNumber),
			)

			logRs, logRsErr := GetLog(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, eppPt, hrmId, lvlId, offset, resultLimit)
			if logRsErr != nil {
				error.IntSrv(ctx, rw, logRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(logRs)" , len(logRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &logRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(eppPt, hrmId, lvlId, pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/log/ep/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [default]")
	}
}

func Put(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Put::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Put::start")

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Put::get request data"))
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr) 
		return
	}

	eppPt    := form.VText (r, "log-ep-tnt-inf-epp-pt")
	hrmId    := form.PInt64(r, "log-ep-tnt-inf-hrm-id")
	lvlId    := form.PInt64(r, "log-ep-tnt-inf-lvl-id")
	tgtLvlId := form.VInt  (r, "log-ep-tnt-mod-tgt-lvl-id")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Put::get data from form",
		slog.String("eppPt"    , eppPt),
		slog.Any   ("hrmId"    , hrmId),
		slog.Any   ("lvlId"    , lvlId),
		slog.Any   ("tgtLvlId" , tgtLvlId),
	)

	exptErrs := []string{}

	putErr := PutLog(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, eppPt, hrmId, lvlId, tgtLvlId, data.User.AurNm, exptErrs)
	if putErr != nil {
		error.IntSrv(ctx, rw, putErr) 
		return
	}

	rw.Header().Set("HX-Trigger", "mod-log")

	notification.Toast(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : data.T ("web-core-auth-log-ep-tnt-mod-bulk-form.message-input-success")}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Put::end")

	return
}
