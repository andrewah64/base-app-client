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
	"github.com/andrewah64/base-app-client/internal/common/core/password"
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

	aurId, aurIdErr := form.VIntArray(r, "aur-tnt-inf-aur-id")
	if aurIdErr != nil {
		error.IntSrv(ctx, rw, aurIdErr)
		return
	}

	if len(aurId) > 0 {
		delErr := DelAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId, nil)
		if delErr != nil{
			error.IntSrv(ctx, rw, delErr)
			return
		}

		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::success")

		rw.Header().Set("HX-Trigger", "mod")

		message := ""

		if len(aurId) == 1 {
			message = data.T("web-core-auth-aur-tnt-del-form.message-delete-success-singular", "n", strconv.Itoa(len(aurId)))
		} else {
			message = data.T("web-core-auth-aur-tnt-del-form.message-delete-success-plural"  , "n", strconv.Itoa(len(aurId)))
		}

		notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::end")

	return
}

func params(aurNm string, dbrlId *int64, aurEnabled *bool, lngId *int64, pageNumber int) string {
	v := url.Values{}

	v.Set("aur-tnt-inf-aur-nm" , aurNm)

	switch dbrlId {
		case nil:
			v.Set("aur-tnt-inf-dbrl-id" , "")
		default :
			v.Set("aur-tnt-inf-dbrl-id" , strconv.FormatInt(*dbrlId, 10))
	}

	switch aurEnabled {
		case nil:
			v.Set("aur-tnt-inf-aur-enabled" , "")
		default :
			v.Set("aur-tnt-inf-aur-enabled" , strconv.FormatBool(*aurEnabled))
	}

	switch lngId {
		case nil:
			v.Set("aur-tnt-inf-lng-id" , "")
		default :
			v.Set("aur-tnt-inf-lng-id" , strconv.FormatInt(*lngId, 10))
	}

	v.Set("aur-tnt-inf-page-number" , strconv.Itoa(pageNumber))

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
			opts := make(map[string]any)

			optsInfRs, optsInfRsErr := OptsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if optsInfRsErr != nil {
				error.IntSrv(ctx, rw, optsInfRsErr)
				return
			}

			opts["Search"] = optsInfRs

			if data.HasRole ("role_web_core_auth_aur_tnt_reg") {
				optsRegRs, optsRegRsErr := OptsReg(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
				if optsRegRsErr != nil {
					error.IntSrv(ctx, rw, optsRegRsErr)
					return
				}

				opts["Register"] = optsRegRs
			}

			aurRs, aurRsErr := GetAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, "", nil, nil, nil, offset, resultLimit)
			if aurRsErr != nil {
				error.IntSrv(ctx, rw, aurRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(*optsInfRs)", len(*optsInfRs)),
				slog.Int("len(aurRs)"     , len(aurRs)),
			)

			data.FormOpts  = &opts
			data.ResultSet = &map[string]any{
				"Search"      : &aurRs,
				"PageNumber"  : pageNumber,
				"ResultLimit" : resultLimit,
				"Params"      : params("", nil, nil, nil, pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/aur/tnt/content", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [page load]")

		case "aur-tnt-inf-scr": // infinite scroll
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			aurEnabled := form.PBool  (r, "aur-tnt-inf-aur-enabled")
			aurNm      := form.VText  (r, "aur-tnt-inf-aur-nm")
			dbrlId     := form.PInt64 (r, "aur-tnt-inf-dbrl-id")
			lngId      := form.PInt64 (r, "aur-tnt-inf-lng-id")
			pageNumber := form.VInt   (r, "aur-tnt-inf-page-number")
			offset     := (pageNumber - 1) * resultLimit

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.Any   ("aurEnabled" , aurEnabled),
				slog.String("aurNm"      , aurNm),
				slog.Any   ("dbrlId"     , dbrlId),
				slog.Any   ("lngId"      , lngId),
				slog.Int   ("pageNumber" , pageNumber),
				slog.Int   ("offset"     , offset),
			)

			aurRs, aurRsErr := GetAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, aurEnabled, dbrlId, lngId, offset, resultLimit)
			if aurRsErr != nil {
				error.IntSrv(ctx, rw, aurRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(aurRs)" , len(aurRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &aurRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(aurNm, dbrlId, aurEnabled, lngId, pageNumber + 1),
			}

			rw.Header().Set("HX-Trigger", "inf")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/aur/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [infinite scroll]")

		case "aur-tnt-inf-form": // search
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			aurEnabled := form.PBool  (r, "aur-tnt-inf-aur-enabled")
			aurNm      := form.VText  (r, "aur-tnt-inf-aur-nm")
			dbrlId     := form.PInt64 (r, "aur-tnt-inf-dbrl-id")
			lngId      := form.PInt64 (r, "aur-tnt-inf-lng-id")
			pageNumber := form.VInt   (r, "aur-tnt-inf-page-number")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.Any   ("aurEnabled" , aurEnabled),
				slog.String("aurNm"      , aurNm),
				slog.Any   ("dbrlId"     , dbrlId),
				slog.Any   ("lngId"      , lngId),
				slog.Int   ("pageNumber" , pageNumber),
			)

			aurRs, aurRsErr := GetAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, aurEnabled, dbrlId, lngId, offset, resultLimit)
			if aurRsErr != nil {
				error.IntSrv(ctx, rw, aurRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(aurRs)" , len(aurRs)),
			)

			data.ResultSet = &map[string]any{
				"Search"      : &aurRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(aurNm, dbrlId, aurEnabled, lngId, pageNumber),
			}

			rw.Header().Set("HX-Trigger", "src")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/aur/tnt/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [search]")

	}
}

func Post(rw http.ResponseWriter, r *http.Request){
	const (
		valTmpl = "core/auth/aur/tnt/fragment/val"
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
	aurNm := form.VText  (r, "aur-tnt-reg-aur-nm")
	aurPw := form.VText  (r, "aur-tnt-reg-pw")
	grpId := form.PInt64 (r, "aur-tnt-reg-grp-id")
	lngId := form.PInt64 (r, "aur-tnt-reg-lng-id")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
		slog.String("AurNm", aurNm),
		slog.Any   ("GrpId", grpId),
		slog.Any   ("LngId", lngId),
	)

	v.Check(validator.NotBlank(aurNm), "aur-tnt-reg-aur-nm", data.T("web-core-auth-aur-tnt-reg-form.warning-input-aur-nm-blank"))
	v.Check(validator.NotBlank(aurPw), "aur-tnt-reg-pw"    , data.T("web-core-auth-aur-tnt-reg-form.warning-input-pw-blank"))
	v.Check(validator.NotNil  (grpId), "aur-tnt-reg-grp-id", data.T("web-core-auth-aur-tnt-reg-form.warning-input-grp-blank"))
	v.Check(validator.NotNil  (lngId), "aur-tnt-reg-lng-id", data.T("web-core-auth-aur-tnt-reg-form.warning-input-lng-blank"))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::validate data retrieved from form",
		slog.Int("len(v.Errors)", len(v.Errors)),
	)

	if !v.Valid() {
		html.Fragment(ctx, ssd.Logger, rw, r, valTmpl, http.StatusUnprocessableEntity, &data)
		return
	}

	aurHshPw, pErr := password.Hash(aurPw)
	if pErr != nil {
		error.IntSrv(ctx, rw, pErr)
		return
	}

	exptErrs := []string{
		pgerrcode.UniqueViolation,
	}

	regErr := PostAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, grpId, aurNm, aurHshPw, lngId, data.User.AurNm, exptErrs)
	if regErr != nil {
		var pgErr *pgconn.PgError

		if errors.As(regErr, &pgErr) {
			switch pgErr.Code {
				case pgerrcode.UniqueViolation:
					v.AddError("aur-tnt-reg-credentials-taken", data.T("web-core-auth-aur-tnt-reg-form.warning-input-credentials-taken", "aurNm", aurNm))
				default:
					slog.LogAttrs(ctx, slog.LevelError, "unexpected error",
						slog.String("regErr.Error()" , regErr.Error()),
						slog.String("pgErr.Code"     , pgErr.Code),
					)

					v.AddError("aur-tnt-reg-unexpected", data.T("web-core-auth-aur-tnt-reg-form.warning-input-unexpected-error"))
			}

			data.ResultSet = &map[string]any{"Validator" : &v}

			html.Fragment(ctx, ssd.Logger, rw, r, valTmpl, http.StatusUnprocessableEntity, &data)

			return
		} else {
			error.IntSrv(ctx, rw, regErr)

			return
		}
	}

	rw.Header().Set("HX-Trigger", "mod")

	notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : data.T ("web-core-auth-aur-tnt-reg-form.message-input-success", "aurNm", aurNm)}, data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::end")

	return
}
