package aur

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
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

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

	aaukId, aaukIdErr := strconv.Atoi(r.PathValue("id"))
	if aaukIdErr != nil || aaukId < 1 {
		http.NotFound(rw, r)
		return
	}

	keyRs, keyRsErr := GetKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aaukId)
	if keyRsErr != nil {
		error.IntSrv(ctx, rw, keyRsErr)
		return
	}

	dbrlRs, dbrlRsErr := GetRol(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukId)
	if dbrlRsErr != nil {
		error.IntSrv(ctx, rw, dbrlRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("aaukId"      , aaukId),
		slog.Int("len(keyRs)"  , len(keyRs)),
		slog.Int("len(dbrlRs)" , len(dbrlRs)),
	)

	data.ResultSet = &map[string]any{"Key": keyRs , "Search": &dbrlRs}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/rol/key/aur/content", http.StatusOK, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")

	return
}

func Patch(rw http.ResponseWriter, r *http.Request){
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

	dbrlId, dbrlIdErr := form.VIntArray(r, "rol-aur-key-inf-dbrl-id")
	if dbrlIdErr != nil {
		error.IntSrv(ctx, rw, dbrlIdErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int("aaukId", aaukId),
		slog.Any("dbrlId", dbrlId),
	)

	patchErr := PatchRol(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukId, dbrlId, data.User.AurNm, nil)
	if patchErr != nil{
		error.IntSrv(ctx, rw, patchErr)
		return
	}

	message := ""

	if len(dbrlId) == 1 {
		message = data.T("web-core-auth-rol-key-aur-mod-form.message-input-success-singular", "n", strconv.Itoa(len(dbrlId)))
	} else {
		message = data.T("web-core-auth-rol-key-aur-mod-form.message-input-success-plural", "n", strconv.Itoa(len(dbrlId)))
	}

	notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")

	return
}
