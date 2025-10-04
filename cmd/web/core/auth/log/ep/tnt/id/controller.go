package id

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
	"github.com/jackc/pgx/v5/pgconn"
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

	ellId, ellIdErr := strconv.Atoi(r.PathValue("id"))
	if ellIdErr != nil || ellId < 1 {
		http.NotFound(rw, r)
		return
	}

	optsRs, optsRsErr := Opts(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if optsRsErr != nil {
		error.IntSrv(ctx, rw, optsRsErr)
		return
	}

	logRs, logRsErr := GetRowLogMod(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, ellId)
	if logRsErr != nil {
		error.IntSrv(ctx, rw, logRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("ellId"       , ellId),
		slog.Int("len(*optsRs)", len(*optsRs)),
		slog.Int("len(logRs)"  , len(logRs)),
	)

	data.FormOpts  = &map[string]any{"Search": &optsRs}
	data.ResultSet = &map[string]any{"Search": &logRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/log/ep/tnt/fragment/modrow", http.StatusCreated, &data)

	if len(logRs) == 0 {
		notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-log-ep-tnt-mod-row-form.warning-input-log-olock-error")}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")

	return
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

	ellId, ellIdErr := strconv.Atoi(r.PathValue("id"))
	if ellIdErr != nil || ellId < 1 {
		http.NotFound(rw, r)
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	lvlId := form.VInt  (r, fmt.Sprintf("log-ep-tnt-mod-lvl-id-%v", ellId))
	uts   := form.VTime (r, fmt.Sprintf("log-ep-tnt-mod-uts-%v"   , ellId))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int("ellId" , ellId),
		slog.Any("lvlId" , lvlId),
		slog.Any("uts"   , uts),
	)

	exptErrs := []string{
		"OLOKU",
		"OLOKD",
	}

	patchErr := PatchLog(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, ellId, lvlId, data.User.AurNm, uts, exptErrs)
	if patchErr != nil{
		Get(rw, r)

		var pgErr *pgconn.PgError

		if errors.As(patchErr, &pgErr) {
			switch pgErr.Code {
				case "OLOKU":
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-log-ep-tnt-mod-row-form.warning-input-log-olock-error")}, data)

				case "OLOKD":
					// intentionally empty

				default:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-log-ep-tnt-mod-row-form.warning-input-unexpected-error")}, data)

			}
		}

		return
	}

	logRs, logRsErr := GetRowLogInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, ellId)
	if logRsErr != nil {
		error.IntSrv(ctx, rw, logRsErr)
		return
	}

	data.ResultSet = &map[string]any{"Search": &logRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/log/ep/tnt/fragment/infrow", http.StatusCreated, &data)

	notification.Toast(ctx, ssd.Logger, rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-log-ep-tnt-mod-row-form.message-input-success")}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}
