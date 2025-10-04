package tnt

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

	grpId, grpIdErr := strconv.Atoi(r.PathValue("id"))
	if grpIdErr != nil || grpId < 1 {
		http.NotFound(rw, r)
		return
	}

	grpRs, grpRsErr := GetGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, grpId)
	if grpRsErr != nil {
		error.IntSrv(ctx, rw, grpRsErr)
		return
	}

	aurRs, aurRsErr := GetAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId)
	if aurRsErr != nil {
		error.IntSrv(ctx, rw, aurRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("grpId"      , grpId),
		slog.Int("len(grpRs)" , len(grpRs)),
		slog.Int("len(aurRs)" , len(aurRs)),
	)

	data.ResultSet = &map[string]any{"Group": &grpRs , "Search": &aurRs}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/aur/grp/tnt/content", http.StatusOK, &data)

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

	aurId, aurIdErr := form.VIntArray(r, "aur-grp-tnt-inf-aur-id")
	if aurIdErr != nil {
		error.IntSrv(ctx, rw, aurIdErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int("grpId" , grpId),
		slog.Any("aurId" , aurId),
	)

	patchErr := PatchAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId, aurId, nil)
	if patchErr != nil{
		error.IntSrv(ctx, rw, patchErr)
		return
	}

	message := ""

	if len(aurId) == 1 {
		message = data.T("web-core-auth-aur-grp-tnt-mod-form.message-input-success-singular", "n", strconv.Itoa(len(aurId)))
	} else {
		message = data.T("web-core-auth-aur-grp-tnt-mod-form.message-input-success-plural", "n", strconv.Itoa(len(aurId)))
	}

	notification.Toast(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")

	return
}
