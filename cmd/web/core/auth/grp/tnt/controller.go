package tnt

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/validator"
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

	grpId, grpIdErr := form.VIntArray(r, "grp-tnt-inf-grp-id")
	if grpIdErr != nil {
		error.IntSrv(ctx, rw, grpIdErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::grpId",
		slog.Any("grpId", grpId),
	)

	if len(grpId) > 0 {
		delErr := DelGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId, nil)
		if delErr != nil{
			error.IntSrv(ctx, rw, delErr)
			return
		}

		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::success")

		rw.Header().Set("HX-Trigger", "mod")

		message := ""

		if len(grpId) == 1 {
			message = data.T("web-core-auth-grp-tnt-del-form.message-delete-success-singular", "n", strconv.Itoa(len(grpId)))
		} else {
			message = data.T("web-core-auth-grp-tnt-del-form.message-delete-success-plural"  , "n", strconv.Itoa(len(grpId)))
		}

		notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::end")

	return
}

func params(grpNm string, aurNm string, dbrlId *int64, pageNumber int) string {
	v := url.Values{}

	v.Set("grp-tnt-inf-grp-nm" , grpNm)

	v.Set("grp-tnt-inf-aur-nm" , aurNm)

	switch dbrlId {
		case nil:
			v.Set("grp-tnt-inf-dbrl-id" , "")
		default :
			v.Set("grp-tnt-inf-dbrl-id" , strconv.FormatInt(*dbrlId, 10))
	}

	v.Set("grp-tnt-inf-page-number" , strconv.Itoa(pageNumber))

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

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
		slog.Int   ("pageNumber"  , pageNumber),
		slog.Int   ("offset"      , offset),
		slog.Int   ("resultLimit" , resultLimit),
		slog.String("trigger"     , trigger),
	)

	switch trigger {
		case "": // page load
			optsInfRs, optsInfRsErr := OptsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if optsInfRsErr != nil {
				error.IntSrv(ctx, rw, optsInfRsErr)
				return
			}

			grpRs, grpRsErr := GetGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, "", "", nil, offset, resultLimit)
			if grpRsErr != nil {
				error.IntSrv(ctx, rw, grpRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(*optsInfRs)" , len(*optsInfRs)),
				slog.Int("len(grpRs)"      , len(grpRs)),
			)

			data.FormOpts  = &map[string]any{"Search": &optsInfRs}
			data.ResultSet = &map[string]any{
				"Search"      : &grpRs,
				"PageNumber"  : pageNumber,
				"ResultLimit" : resultLimit,
				"Params"      : params("", "", nil, pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/grp/tnt/content", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [page load]")

		case "grp-tnt-inf-scr": // infinite scroll
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			grpNm      := form.VText  (r, "grp-tnt-inf-grp-nm")
			aurNm      := form.VText  (r, "grp-tnt-inf-aur-nm")
			dbrlId     := form.PInt64 (r, "grp-tnt-inf-dbrl-id")
			pageNumber := form.VInt   (r, "grp-tnt-inf-page-number")
			offset     := (pageNumber - 1) * resultLimit

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("grpNm"      , grpNm),
				slog.String("aurNm"      , aurNm),
				slog.Any   ("dbrlId"     , dbrlId),
				slog.Int   ("pageNumber" , pageNumber),
				slog.Int   ("offset"     , offset),
			)

			grpRs, grpRsErr := GetGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, grpNm, aurNm, dbrlId, offset, resultLimit)
			if grpRsErr != nil {
				error.IntSrv(ctx, rw, grpRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(grpRs)" , len(grpRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &grpRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(grpNm, aurNm, dbrlId, pageNumber + 1),
			}

			rw.Header().Set("HX-Trigger", "inf")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/grp/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [infinite scroll]")

		case "grp-tnt-inf-form":
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			grpNm      := form.VText  (r, "grp-tnt-inf-grp-nm")
			aurNm      := form.VText  (r, "grp-tnt-inf-aur-nm")
			dbrlId     := form.PInt64 (r, "grp-tnt-inf-dbrl-id")
			pageNumber := form.VInt   (r, "grp-tnt-inf-page-number")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("grpNm"      , grpNm),
				slog.String("aurNm"      , aurNm),
				slog.Any   ("dbrlId"     , dbrlId),
				slog.Int   ("pageNumber" , pageNumber),
			)

			grpRs, grpRsErr := GetGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, grpNm, aurNm, dbrlId, offset, resultLimit)
			if grpRsErr != nil {
				error.IntSrv(ctx, rw, grpRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(grpRs)" , len(grpRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &grpRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(grpNm, aurNm, dbrlId, pageNumber),
			}

			rw.Header().Set("HX-Trigger", "src")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/grp/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [search]")
	}

	return
}

func Post (rw http.ResponseWriter, r *http.Request) {
	const (
		valTmpl = "core/auth/grp/tnt/fragment/val"
	)

	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Post::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::start")

	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Post::get request data"))
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	v     := validator.New()
	grpNm := form.VText (r, "grp-tnt-reg-grp-nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
		slog.String("grpNm", grpNm),
	)

	v.Check(validator.NotBlank(grpNm), "grp-tnt-reg-grp-nm-blank", data.T("web-core-auth-grp-tnt-reg-form.warning-input-grp-nm-blank"))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::validate data retrieved from form",
		slog.Int("len(v.Errors)", len(v.Errors)),
	)

	if !v.Valid() {
		data.ResultSet = &map[string]any{"Validator" : &v}

		html.Fragment(ctx, ssd.Logger, rw, r, valTmpl, http.StatusUnprocessableEntity, &data)

		return
	}

	exptErrs := []string{
		pgerrcode.UniqueViolation,
	}

	regErr := PostGrp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpNm, data.User.AurNm, exptErrs)
	if regErr != nil {
		var pgErr *pgconn.PgError

		if errors.As(regErr, &pgErr) {
			switch pgErr.Code {
				case pgerrcode.UniqueViolation:
					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "group already exists",
						slog.String("grpNm"  , grpNm),
						slog.String("regErr" , regErr.Error()),
					)

					v.AddError("grp-tnt-reg-grp-nm-taken", data.T("web-core-auth-grp-tnt-reg-form.warning-input-grp-nm-taken", "grpNm", grpNm))
				default:
					slog.LogAttrs(ctx, slog.LevelError, "unexpected error",
						slog.String("regErr"    , regErr.Error()),
						slog.String("pgErr.code", pgErr.Code),
					)

					v.AddError("grp-tnt-reg-unexpected", data.T("web-core-auth-grp-tnt-reg-form.warning-input-unexpected-error"))
			}

			data.ResultSet = &map[string]any{"Validator" : &v}

			html.Fragment(ctx, ssd.Logger, rw, r, valTmpl, http.StatusUnprocessableEntity, &data)

			return
		} else {
			error.IntSrv(ctx, rw, regErr)
			return
		}
	}

	data.ResultSet = &map[string]any{"GrpNm": &grpNm}

	rw.Header().Set("HX-Trigger", "mod")
	
	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/grp/tnt/fragment/res", http.StatusCreated, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::end")

	return
}
