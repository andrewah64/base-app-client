package aur

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/key"
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

	aaukId, aaukIdErr := form.VIntArray(r, "key-aur-mod-aauk-id")
	if aaukIdErr != nil {
		error.IntSrv(ctx, rw, aaukIdErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::get selected keys",
		slog.Int("len(aaukId)", len(aaukId)),
	)

	if len(aaukId) > 0 {
		delErr := DelKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukId, nil)
		if delErr != nil{
			error.IntSrv(ctx, rw, delErr)
			return
		}

		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::success")

		rw.Header().Set("HX-Trigger", "mod")

		message := ""

		if len(aaukId) == 1 {
			message = data.T("web-core-auth-key-aur-del-form.message-delete-success-singular", "n", strconv.Itoa(len(aaukId)))
		} else {
			message = data.T("web-core-auth-key-aur-del-form.message-delete-success-plural"  , "n", strconv.Itoa(len(aaukId)))
		}

		notification.Show(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::end")

	return
}

func params(aaukNm string, dbrlId *int64, aaukEnabled *bool, pageNumber int) string {
	v := url.Values{}

	v.Set("key-aur-inf-aauk-nm" , aaukNm)

	switch dbrlId {
		case nil:
			v.Set("key-aur-inf-dbrl-id" , "")
		default :
			v.Set("key-aur-inf-dbrl-id" , strconv.FormatInt(*dbrlId, 10))
	}

	switch aaukEnabled {
		case nil:
			v.Set("key-aur-inf-aauk-enabled" , "")
		default :
			v.Set("key-aur-inf-aauk-enabled" , strconv.FormatBool(*aaukEnabled))
	}

	v.Set("key-aur-inf-page-number" , strconv.Itoa(pageNumber))

	return v.Encode()
}

func Get(rw http.ResponseWriter, r *http.Request){
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
			optsInfRs, optsInfRsErr := OptsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId)
			if optsInfRsErr != nil {
				error.IntSrv(ctx, rw, optsInfRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve opts dataset",
				slog.Int("len(*optsInfRs)", len(*optsInfRs)),
			)

			keyRs, keyRsErr := GetKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, "", nil, nil, offset, resultLimit)
			if keyRsErr != nil {
				error.IntSrv(ctx, rw, keyRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve key dataset",
				slog.Int("len(keyRs)", len(keyRs)),
			)

			data.FormOpts  = &map[string]any{"Search" : &optsInfRs}
			data.ResultSet = &map[string]any{
				"Search"      : &keyRs,
				"PageNumber"  : pageNumber,
				"ResultLimit" : resultLimit,
				"Params"      : params("", nil, nil, pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/key/aur/content", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [page load]")

		case "key-aur-mod-scr": // infinite scroll
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr) 
				return
			}

			aaukNm      := form.VText  (r, "key-aur-inf-aauk-nm")
			dbrlId      := form.PInt64 (r, "key-aur-inf-dbrl-id")
			aaukEnabled := form.PBool  (r, "key-aur-inf-aauk-enabled")
			pageNumber  := form.VInt   (r, "key-aur-inf-page-number")
			offset      := (pageNumber - 1) * resultLimit

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("aaukNm"     , aaukNm),
				slog.Any   ("dbrlId"     , dbrlId),
				slog.Any   ("aaukEnabled", aaukEnabled),
				slog.Int   ("pageNumber" , pageNumber),
				slog.Int   ("offset"     , offset),
			)

			keyRs, keyRsErr := GetKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukNm, aaukEnabled, dbrlId, offset, resultLimit)
			if keyRsErr != nil {
				error.IntSrv(ctx, rw, keyRsErr)
				return
			}

			data.ResultSet = &map[string]any{
				"Search"      : &keyRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(aaukNm, dbrlId, aaukEnabled, pageNumber + 1),
			}

			rw.Header().Set("HX-Trigger", "inf")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/key/aur/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [infinite scroll]")

		case "key-aur-inf-form": // search
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr) 
				return
			}

			aaukNm      := form.VText  (r, "key-aur-inf-aauk-nm")
			dbrlId      := form.PInt64 (r, "key-aur-inf-dbrl-id")
			aaukEnabled := form.PBool  (r, "key-aur-inf-aauk-enabled")
			pageNumber  := form.VInt   (r, "key-aur-inf-page-number")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("aaukNm"     , aaukNm),
				slog.Any   ("dbrlId"     , dbrlId),
				slog.Any   ("aaukEnabled", aaukEnabled),
				slog.Int   ("pageNumber ", pageNumber),
			)

			keyRs, keyRsErr := GetKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, aaukNm, aaukEnabled, dbrlId, offset, resultLimit)
			if keyRsErr != nil {
				error.IntSrv(ctx, rw, keyRsErr)
				return
			}

			data.ResultSet = &map[string]any{
				"Search"      : &keyRs,
				"ResultLimit" : resultLimit,
				"Params"      : params(aaukNm, dbrlId, aaukEnabled, pageNumber),
			}

			rw.Header().Set("HX-Trigger", "src")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/key/aur/template/res", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [search]")
	}

	return
}

func Post(rw http.ResponseWriter, r *http.Request){
	const (
		valTmpl = "core/auth/key/aur/fragment/val"
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

	v      := validator.New()
	aaukNm := form.VText(r, "key-aur-reg-aauk-nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
		slog.String("aaukNm", aaukNm),
	)

	v.Check(validator.NotBlank(aaukNm), "key-aur-reg-aauk-nm-blank", data.T("web-core-auth-key-aur-reg-form.warning-input-aauk-nm-blank"))

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::validate data retrieved from form",
		slog.Int("len(v.Errors)", len(v.Errors)),
	)

	if !v.Valid() {
		data.ResultSet = &map[string]any{"Validator" : &v}

		html.Fragment(ctx, ssd.Logger, rw, r, valTmpl, http.StatusUnprocessableEntity, &data)

		return
	}

	aaukKey, aaukKeyErr := key.Key(16)
	if aaukKeyErr != nil {
		error.IntSrv(ctx, rw, aaukKeyErr)
		return
	}

	exptErrs := []string{
		pgerrcode.UniqueViolation,
	}

	regErr := PostKey(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, data.User.AurId, key.Hash(aaukKey), true, aaukNm, data.User.AurNm, exptErrs)
	if regErr != nil {
		var pgErr *pgconn.PgError

		if errors.As(regErr, &pgErr) {
			switch pgErr.Code {
				case pgerrcode.UniqueViolation:
					v.AddError("key-aur-reg-aauk-nm-taken", data.T("web-core-auth-key-aur-reg-form.warning-input-aauk-nm-taken", "aaukNm", aaukNm))
				default:
					v.AddError("key-aur-reg-unexpected"   , data.T("web-core-auth-key-aur-tnt-reg-form.warning-input-unexpected-error"))
			}

			data.ResultSet = &map[string]any{"Validator" : &v}

			html.Fragment(ctx, ssd.Logger, rw, r, valTmpl, http.StatusUnprocessableEntity, &data)

			return
		} else {
			error.IntSrv(ctx, rw, regErr)
			return
		}
	}

	data.ResultSet = &map[string]any{"Key": &aaukKey}

	rw.Header().Set("HX-Trigger", "mod")

	html.Fragment(ctx, ssd.Logger, rw, r, "core/auth/key/aur/fragment/res", http.StatusCreated, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::end")

	return
}
