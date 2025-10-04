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

	aurId, aurIdErr := strconv.Atoi(r.PathValue("id"))
	if aurIdErr != nil || aurId < 1 {
		http.NotFound(rw, r)
		return
	}

	aurRs, aurRsErr := GetAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId)
	if aurRsErr != nil {
		error.IntSrv(ctx, rw, aurRsErr)
		return
	}

	grpRs, grpRsErr := GetGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aurId)
	if grpRsErr != nil {
		error.IntSrv(ctx, rw, grpRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int("aurId"      , aurId),
		slog.Int("len(aurRs)" , len(aurRs)),
		slog.Int("len(grpRs)" , len(grpRs)),
	)

	data.ResultSet = &map[string]any{"Aur": &aurRs , "Search": &grpRs}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/grp/aur/tnt/content", http.StatusOK, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
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

	tgtAurId, tgtAurIdErr := strconv.Atoi(r.PathValue("id"))
	if tgtAurIdErr != nil || tgtAurId < 1 {
		http.NotFound(rw, r)
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	grpId, grpIdErr := form.VIntArray(r, "grp-aur-tnt-inf-grp-id")
	if grpIdErr != nil {
		error.IntSrv(ctx, rw, grpIdErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from form",
		slog.Int("tgtAurId" , tgtAurId),
		slog.Any("grpId"    , grpId),
	)

	patchErr := PatchGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, tgtAurId, grpId, data.User.AurNm, nil)
	if patchErr != nil{
		error.IntSrv(ctx, rw, patchErr)
		return
	}

	data.ResultSet = &map[string]any{"GroupCount": len(grpId)}

	numGrps := len(grpId)
	message := ""

	switch numGrps {
		case 1 :
			message = data.T("web-core-auth-grp-aur-tnt-mod-form.message-input-success-singular", "n", fmt.Sprintf("%v", numGrps))
		default :
			message = data.T("web-core-auth-grp-aur-tnt-mod-form.message-input-success-plural"  , "n", fmt.Sprintf("%v", numGrps))
	}

	notification.Toast(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}
