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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Get(rw http.ResponseWriter, r *http.Request) {
	const (
		dataKey = "Search"
	)

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

	optsRs, optsRsErr := Opts(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId)
	if optsRsErr != nil {
		error.IntSrv(ctx, rw, optsRsErr)
		return
	}

	aurRs, aurRsErr := GetRowAurMod(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId)
	if aurRsErr != nil {
		error.IntSrv(ctx, rw, aurRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("aurId"       , aurId),
		slog.Int("len(*optsRs)", len(*optsRs)),
		slog.Int("len(aurRs)"  , len(aurRs)),
	)

	data.FormOpts  = &map[string]any{"Search": &optsRs}
	data.ResultSet = &map[string]any{"Search": &aurRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/aur/tnt/fragment/modrow", http.StatusCreated, &data)

	if len(aurRs) == 0 {
		notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-aur-tnt-mod-form.warning-input-log-olock-error")}, data)
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

	aurNm      := form.VText(r, fmt.Sprintf("aur-tnt-mod-aur-nm-%v"  , aurId))
	aurEnabled := form.VBool(r, fmt.Sprintf("aur-tnt-mod-enabled-%v" , aurId))
	lngId      := form.VInt (r, fmt.Sprintf("aur-tnt-mod-language-%v", aurId))
	pgId       := form.VInt (r, fmt.Sprintf("aur-tnt-mod-page-%v"    , aurId))
	uts        := form.VTime(r, fmt.Sprintf("aur-tnt-mod-uts-%v"     , aurId))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int   ("aurId"      , aurId),
		slog.String("aurNm"      , aurNm),
		slog.Bool  ("aurEnabled" , aurEnabled),
		slog.Int   ("lngId"      , lngId),
		slog.Int   ("pgId"       , pgId),
		slog.Any   ("uts"        , uts),
	)

	exptErrs := []string{
		"OLOKU",
		"OLOKD",
		pgerrcode.CheckViolation,
		pgerrcode.UniqueViolation,
	}

	patchErr := PatchAur (&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId, aurNm, aurEnabled, lngId, pgId, data.User.AurNm, uts, exptErrs)
	if patchErr != nil{
		Get(rw, r)

		var pgErr *pgconn.PgError

		if errors.As(patchErr, &pgErr) {
			switch pgErr.Code {
				case "OLOKU":
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-aur-tnt-mod-form.warning-input-log-olock-error")}, data)

				case "OLOKD":
					// intentionally blank

				case pgerrcode.CheckViolation:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-aur-tnt-mod-form.warning-input-aur-nm-blank")}, data)

				case pgerrcode.UniqueViolation:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-aur-tnt-mod-form.warning-input-aur-nm-taken", "aurNm", aurNm)}, data)

				default:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-aur-tnt-mod-form.warning-input-unexpected-error")}, data)

			}
		}

		return
	}

	aurRs, aurRsErr := GetRowAurInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId)
	if aurRsErr != nil {
		error.IntSrv(ctx, rw, aurRsErr)
		return
	}

	data.ResultSet = &map[string]any{"Search": &aurRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/aur/tnt/fragment/infrow", http.StatusCreated, &data)

	notification.Toast(ctx, ssd.Logger, rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-aur-tnt-mod-form.message-input-success")}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}
