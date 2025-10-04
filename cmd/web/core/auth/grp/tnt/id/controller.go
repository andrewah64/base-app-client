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

	grpId, grpIdErr := strconv.Atoi(r.PathValue("id"))
	if grpIdErr != nil || grpId < 1 {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get grpId"))
		return
	}

	grpRs, grpRsErr := GetRowGrpMod(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId)
	if grpRsErr != nil {
		error.IntSrv(ctx, rw, grpRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("grpId"      , grpId),
		slog.Int("len(grpRs)" , len(grpRs)),
	)

	data.ResultSet = &map[string]any{"Group": &grpRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/grp/tnt/fragment/modrow", http.StatusCreated, &data)

	if len(grpRs) == 0 {
		notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-grp-tnt-mod-form.warning-input-grp-olock-error")}, data)
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

	grpId, grpIdErr := strconv.Atoi(r.PathValue("id"))
	if grpIdErr != nil || grpId < 1 {
		http.NotFound(rw, r)
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	grpNm := form.VText (r, fmt.Sprintf("grp-tnt-mod-grp-nm-%v", grpId))
	uts   := form.VTime (r, fmt.Sprintf("grp-tnt-mod-uts-%v"   , grpId))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int   ("grpId" , grpId),
		slog.String("grpNm" , grpNm),
		slog.Any   ("Any"   , uts),
	)

	exptErrs := []string{
		"OLOKU",
		"OLOKD",
		pgerrcode.CheckViolation,
		pgerrcode.UniqueViolation,
	}

	patchErr := PatchGrp (&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId, grpNm, data.User.AurNm, uts, exptErrs)
	if patchErr != nil{
		Get(rw, r)

		var pgErr *pgconn.PgError

		if errors.As(patchErr, &pgErr) {
			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch: PatchGrp params",
				slog.Int   ("ssd.TntId" , ssd.TntId),
				slog.Int   ("grpId"     , grpId),
				slog.String("grpNm"     , grpNm),
				slog.String("patchErr"  , patchErr.Error()),
			)

			switch pgErr.Code {
				case "OLOKU":
					notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-grp-tnt-mod-form.warning-input-grp-olock-error")}, data)

				case "OLOKD":
					//intentionally empty

				case pgerrcode.CheckViolation:
					notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-grp-tnt-mod-form.warning-input-grp-nm-blank")}, data)

				case pgerrcode.UniqueViolation:
					notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-grp-tnt-mod-form.warning-input-grp-nm-taken", "grpNm" , grpNm)}, data)

				default:
					notification.Show(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-grp-tnt-mod-form.warning-input-unexpected-error")}, data)

			}

			return
		}
	}

	grpRs, grpRsErr := GetRowGrpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId)
	if grpRsErr != nil {
		error.IntSrv(ctx, rw, grpRsErr)
		return
	}

	data.ResultSet = &map[string]any{"Group": &grpRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/grp/tnt/fragment/infrow", http.StatusCreated, &data)

	notification.Show(ctx, slog.Default(), rw, r, "success" , &map[string]string{"Message" : data.T("web-core-auth-grp-tnt-mod-form.message-input-success")}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")

	return
}
