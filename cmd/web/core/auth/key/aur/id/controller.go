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

	aaukId, aaukIdErr := strconv.Atoi(r.PathValue("id"))
	if aaukIdErr != nil || aaukId < 1 {
		http.NotFound(rw, r)
		return
	}

	keyRs, keyRsErr := GetRowKeyMod(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukId)
	if keyRsErr != nil {
		error.IntSrv(ctx, rw, keyRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("aaukId"     , aaukId),
		slog.Int("len(keyRs)" , len(keyRs)),
	)

	data.ResultSet = &map[string]any{"Key": &keyRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/key/aur/fragment/modrow", http.StatusCreated, &data)

	if len(keyRs) == 0 {
		notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-key-aur-mod-form.warning-input-aauk-olock-error")}, data)
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

	aaukId, aaukIdErr := strconv.Atoi(r.PathValue("id"))
	if aaukIdErr != nil || aaukId < 1 {
		http.NotFound(rw, r)
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	aaukNm      := form.VText(r, fmt.Sprintf("key-aur-mod-aauk-nm-%v"      , aaukId))
	aaukEnabled := form.VBool(r, fmt.Sprintf("key-aur-mod-aauk-enabled-%v" , aaukId))
	uts         := form.VTime(r, fmt.Sprintf("key-aur-mod-uts-%v"          , aaukId))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int   ("aaukId"      , aaukId),
		slog.String("aaukNm"      , aaukNm),
		slog.Bool  ("aaukEnabled" , aaukEnabled),
		slog.Any   ("uts"         , uts),
	)

	exptErrs := []string{
		"OLOKU",
		"OLOKD",
		pgerrcode.UniqueViolation,
		pgerrcode.CheckViolation,
	}

	patchErr := PatchKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukId, aaukNm, aaukEnabled, data.User.AurNm, uts, exptErrs)
	if patchErr != nil{
		Get(rw, r)

		var pgErr *pgconn.PgError

		if errors.As(patchErr, &pgErr) {
			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch: PatchKey params",
				slog.Int   ("ssd.TntId"   , ssd.TntId),
				slog.Int   ("aaukId"      , aaukId),
				slog.String("aaukNm"      , aaukNm),
				slog.Bool  ("aaukEnabled" , aaukEnabled),
				slog.String("patchErr"    , patchErr.Error()),
			)

			switch pgErr.Code {
				case "OLOKU":
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-key-aur-mod-form.warning-input-aauk-olock-error")}, data)

				case "OLOKD":
					// intentionally empty

				case pgerrcode.CheckViolation:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-key-aur-mod-form.warning-input-aauk-nm-blank")}, data)

				case pgerrcode.UniqueViolation:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-key-aur-mod-form.warning-input-aauk-nm-taken", "aaukNm", aaukNm)}, data)

				default:
					notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-key-aur-mod-form.warning-input-unexpected-error")}, data)

			}

			return
		}
	}

	keyRs, keyRsErr := GetRowKeyInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukId)
	if keyRsErr != nil {
		error.IntSrv(ctx, rw, keyRsErr)
		return
	}

	data.ResultSet = &map[string]any{"Key": &keyRs}

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/key/aur/fragment/infrow", http.StatusCreated, &data)

	notification.Toast(ctx, ssd.Logger, rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-key-aur-mod-form.message-input-success")}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")

	return
}
