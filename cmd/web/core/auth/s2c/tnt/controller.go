package tnt

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	"github.com/andrewah64/base-app-client/cmd/web/core/auth/s2c/tnt/val"
)

import (
	"github.com/jackc/pgx/v5/pgconn"
)

import (
	gosaml2types "github.com/russellhaering/gosaml2/types"
)

func idpParams(idpNm string, idpEntityId string, idpEnabled *bool, pageNumber int) string {
	v := url.Values{}

	v.Set("s2c-tnt-inf-idp-nm"        , idpNm)

	v.Set("s2c-tnt-inf-idp-entity-id" , idpEntityId)

	switch idpEnabled {
		case nil:
			v.Set("s2c-tnt-inf-idp-enabled" , "")
		default :
			v.Set("s2c-tnt-inf-idp-enabled" , strconv.FormatBool(*idpEnabled))
	}

	v.Set("s2c-tnt-inf-idp-page-number" , strconv.Itoa(pageNumber))

	return v.Encode()
}

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

	switch r.PathValue("nm") {
		case "idp" :
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			idpId, idpIdErr := form.VIntArray(r, "s2c-tnt-inf-idp-id")
			if idpIdErr != nil {
				error.IntSrv(ctx, rw, idpIdErr)
				return
			}

			if len(idpId) > 0 {
				delErr := DelIdp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpId, nil)
				if delErr != nil{
					error.IntSrv(ctx, rw, delErr)
					return
				}

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::success")

				rw.Header().Set("HX-Trigger", `{"mod":{"target":"#s2c-tnt-inf-idp-form"}}`)

				message := ""

				if len(idpId) == 1 {
					message = data.T("web-core-auth-s2c-tnt-del-idp-form.message-delete-success-singular", "n", strconv.Itoa(len(idpId)))
				} else {
					message = data.T("web-core-auth-s2c-tnt-del-idp-form.message-delete-success-plural"  , "n", strconv.Itoa(len(idpId)))
				}

				notification.Toast(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : message}, data)
			}
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Delete::end")

	return
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

	p := r.URL.Query()

	if p.Has("ntf") && p.Has("lvl"){
		notification.Toast(ctx, slog.Default(), rw, r, p.Get("lvl"), &map[string]string{"Message" : data.T(p.Get("ntf"))} , data)
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
		case "" : // page load
			optsRs, optsRsErr := Opts(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if optsRsErr != nil {
				error.IntSrv(ctx, rw, optsRsErr)
				return
			}

			idpInfRs, idpInfRsErr := GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, "", "", nil, offset, resultLimit)
			if idpInfRsErr != nil {
				error.IntSrv(ctx, rw, idpInfRsErr)
				return
			}

			s2cInfRs, s2cInfRsErr := GetS2cInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if s2cInfRsErr != nil {
				error.IntSrv(ctx, rw, s2cInfRsErr)
				return
			}

			data.FormOpts  = &map[string]any{
				"Opts" : &optsRs,
			}

			data.ResultSet = &map[string]any{
				"Idp"         : &idpInfRs,
				"S2c"         : &s2cInfRs,
				"PageNumber"  : pageNumber,
				"ResultLimit" : resultLimit,
				"IdpParams"   : idpParams("", "" , nil, pageNumber),
			}

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/s2c/tnt/content", http.StatusOK, data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [page load]")

		case "s2c-tnt-inf-idp-scr" : // idp infinite scroll
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			idpNm       := form.VText (r, "s2c-tnt-inf-idp-nm")
			idpEntityId := form.VText (r, "s2c-tnt-inf-idp-entity-id")
			idpEnabled  := form.PBool (r, "s2c-tnt-inf-idp-enabled")
			pageNumber  := form.VInt  (r, "s2c-tnt-inf-idp-page-number")
			offset      := (pageNumber - 1) * resultLimit

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("idpNm"       , idpNm),
				slog.String("idpEntityId" , idpEntityId),
				slog.Any   ("idpEnabled"  , idpEnabled),
				slog.Int   ("pageNumber"  , pageNumber),
				slog.Int   ("offset"      , offset),
			)

			idpInfRs, idpInfRsErr := GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm, idpEntityId, idpEnabled, offset, resultLimit)
			if idpInfRsErr != nil {
				error.IntSrv(ctx, rw, idpInfRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::retrieve datasets",
				slog.Int("len(idpInfRs)" , len(idpInfRs)),
			)

			data.ResultSet = &map[string]any{
				"Idp"         : &idpInfRs,
				"IdpParams"   : idpParams(idpNm, idpEntityId, idpEnabled, pageNumber + 1),
				"ResultLimit" : resultLimit,
			}

			rw.Header().Set("HX-Trigger", "inf")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/s2c/tnt/template/res-idp", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [infinite scroll]")

		case "s2c-tnt-inf-idp-form" : // idp search
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			idpNm       := form.VText (r, "s2c-tnt-inf-idp-nm")
			idpEntityId := form.VText (r, "s2c-tnt-inf-idp-entity-id")
			idpEnabled  := form.PBool (r, "s2c-tnt-inf-idp-enabled")
			pageNumber  := form.VInt  (r, "s2c-tnt-inf-idp-page-number")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::get data from form",
				slog.String("idpNm"       , idpNm),
				slog.String("idpEntityId" , idpEntityId),
				slog.Any   ("idpEnabled"  , idpEnabled),
				slog.Int   ("pageNumber"  , pageNumber),
			)

			idpInfRs, idpInfRsErr := GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm, idpEntityId, idpEnabled, offset, resultLimit)
			if idpInfRsErr != nil {
				error.IntSrv(ctx, rw, idpInfRsErr)
				return
			}

			data.ResultSet = &map[string]any{
				"Idp"         : &idpInfRs,
				"IdpParams"   : idpParams(idpNm, idpEntityId, idpEnabled, pageNumber),
				"ResultLimit" : resultLimit,
			}

			rw.Header().Set("HX-Trigger", "src")

			html.Tmpl(ctx, ssd.Logger, rw, r, "core/auth/s2c/tnt/template/res-idp", http.StatusOK, &data)

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end [search]")
	}
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
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request data"))
		return
	}

	currentUrl := r.Header.Get("HX-Current-URL")

	switch r.PathValue("nm") {
		case "gen" :
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			s2cEntityId := form.VText (r, "s2c-tnt-mod-gen-entity-id")
			s2cEnabled  := form.VBool (r, "s2c-tnt-mod-gen-enabled")
			aumId       := form.VInt  (r, "s2c-tnt-mod-gen-aum-id")
			uts         := form.VTime (r, "s2c-tnt-mod-gen-uts")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::get data from gen form",
				slog.Bool  ("s2cEnabled" , s2cEnabled),
				slog.String("s2cEntityId", s2cEntityId),
				slog.Int   ("aumId"      , aumId),
				slog.Any   ("uts"        , uts),
			)

			exptErrs := []string{
				"OLOCK",
			}

			patchErr := PatchS2c(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, s2cEnabled, s2cEntityId, aumId, data.User.AurNm, uts, exptErrs)
			if patchErr != nil {
				var pgErr *pgconn.PgError

				if errors.As(patchErr, &pgErr) {
					switch pgErr.Code {
						case "OLOCK":
							rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"%v", "target":"#main", "select":"#content", "swap" : "innerHTML show:window:top", "values":{"ntf": "web-core-auth-s2c-tnt-mod-gen-form.warning-input-olock-error", "lvl": "error"}}`, currentUrl))

						default:
							slog.LogAttrs(ctx, slog.LevelError, "Patch::unexpected error",
								slog.String("s2cEntityId", s2cEntityId),
								slog.Bool  ("s2cEnabled" , s2cEnabled),
								slog.Int   ("aumId"      , aumId),
								slog.Any   ("uts"        , uts),
							)

							notification.Toast(ctx, slog.Default(), rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-gen-form.warning-input-unexpected-error")}, data)

							return
					}
				}
			}

			s2cUtsInfRs, s2cUtsInfRsErr := GetS2cUtsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if s2cUtsInfRsErr != nil {
				error.IntSrv(ctx, rw, s2cUtsInfRsErr)
				return
			}

			html.HiddenUtsFragment(rw, "s2c-tnt-mod-gen-uts-ctr", "s2c-tnt-mod-gen-uts", "s2c-tnt-mod-gen-uts", s2cUtsInfRs[0].Uts, data.TFT())

			notification.Toast(ctx, slog.Default(), rw, r, "success", &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-mod-gen-form.message-input-success")} , data)

	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Patch::end")
}

func Post(rw http.ResponseWriter, r *http.Request){
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

	switch r.PathValue("nm") {
		case "mde":
			pfErr := r.ParseForm()
			if pfErr != nil {
				error.IntSrv(ctx, rw, pfErr)
				return
			}

			idpNm  := form.VText (r, "s2c-tnt-reg-mde-idp-nm")
			mdeUrl := form.VText (r, "s2c-tnt-reg-mde-url")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from mde form",
				slog.String("idpNm"  , idpNm),
				slog.String("mdeUrl" , mdeUrl),
			)

			valRs, valRsErr := val.GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm)
			if valRsErr != nil {
				error.IntSrv(ctx, rw, valRsErr)
				return
			}

			if ! valRs[0].IdpNmOk {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-idp-nm-taken", "idpNm", idpNm)}, data)

				return
			}

			mdeUrlRes, mdeUrlResErr := http.Get(mdeUrl)
			if mdeUrlResErr != nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-empty-response")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: Get response payload",
					slog.String("mdeUrlResErr", mdeUrlResErr.Error()),
					slog.String("idpNm"       , idpNm),
					slog.String("mdeUrl"      , mdeUrl),
				)

				return
			}

			mtdRaw, mtdRawErr := io.ReadAll(mdeUrlRes.Body)
			if mtdRawErr != nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-unreadable-response")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: Get response payload",
					slog.String("mtdRawErr" , mtdRawErr.Error()),
					slog.String("idpNm"     , idpNm),
					slog.String("mdeUrl"    , mdeUrl),
				)

				return
			}

			mtd    := &gosaml2types.EntityDescriptor{}
			mtdErr := xml.Unmarshal(mtdRaw, mtd)
			if mtdErr != nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-read-metadata")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: unmarshal metadata into struct",
					slog.String("mtdErr" , mtdErr.Error()),
					slog.String("idpNm"  , idpNm),
					slog.String("mdeUrl" , mdeUrl),
				)

				return
			}

			if mtd.IDPSSODescriptor == nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-read-metadata")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: the metadata did not contain IDPSSODescriptor elements",
					slog.String("idpNm"  , idpNm),
					slog.String("mdeUrl" , mdeUrl),
				)

				return
			}

			lkd := len(mtd.IDPSSODescriptor.KeyDescriptors)

			ipcCrt    := make([][]byte    , lkd)
			cruNm     := make([]string    , lkd)
			ipcIncTs  := make([]time.Time , lkd)
			ipcExpTs  := make([]time.Time , lkd)

			for i, kds := range mtd.IDPSSODescriptor.KeyDescriptors {
				for _, x5c := range kds.KeyInfo.X509Data.X509Certificates {
					if strings.TrimSpace(x5c.Data) == "" {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-empty-cert")}, data)

						ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: x5c.Data is empty",
							slog.Int   ("i"   , i),
							slog.String("kds" , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					x5d, x5dErr := base64.StdEncoding.DecodeString(strings.TrimSpace(x5c.Data))
					if x5dErr != nil {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-decode-cert")}, data)

						slog.LogAttrs(ctx, slog.LevelError, "Post:: cannot decode x5c.Data",
							slog.String("x5dErr" , x5dErr.Error()),
							slog.Int   ("i"      , i),
							slog.String("kds"    , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					crt, crtErr := x509.ParseCertificate(x5d)
					if crtErr != nil {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.warning-input-gen-cert")}, data)

						slog.LogAttrs(ctx, slog.LevelError, "Post:: cannot parse x5d",
							slog.String("crtErr" , crtErr.Error()),
							slog.Int   ("i"      , i),
							slog.String("kds"    , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					ipcCrt[i]   = x5d
					cruNm[i]    = kds.Use
					ipcIncTs[i] = crt.NotBefore
					ipcExpTs[i] = crt.NotAfter
				}
			}

			sloUrl    := make([]string, len(mtd.IDPSSODescriptor.SingleLogoutServices))
			sloUrlBnd := make([]string, len(mtd.IDPSSODescriptor.SingleLogoutServices))

			for i, slo := range mtd.IDPSSODescriptor.SingleLogoutServices {
				sloUrl[i]    = slo.Location
				sloUrlBnd[i] = slo.Binding
			}

			ssoUrl    := make([]string, len(mtd.IDPSSODescriptor.SingleSignOnServices))
			ssoUrlBnd := make([]string, len(mtd.IDPSSODescriptor.SingleSignOnServices))

			for i, sso := range mtd.IDPSSODescriptor.SingleSignOnServices {
				ssoUrl[i]    = sso.Location
				ssoUrlBnd[i] = sso.Binding
			}

			postErr := PostIdp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm, mtd.EntityID, ipcCrt, cruNm, ipcIncTs, ipcExpTs, &mdeUrl, sloUrl, sloUrlBnd, ssoUrl, ssoUrlBnd, data.User.AurNm, nil)
			if postErr != nil {
				error.IntSrv(ctx, rw, postErr)
				return
			}

			rw.Header().Set("HX-Trigger", `{"mod":{"target":"#s2c-tnt-inf-idp-form"}}`)

			notification.Toast(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-mde-form.message-input-success", "idpNm", idpNm)}, data)

		case "xml":
			mpfErr := r.ParseMultipartForm(200 * 1024) // 200 Kb file upload limit
			if mpfErr != nil {
				error.IntSrv(ctx, rw, mpfErr)
				return
			}

			idpNm := form.VText(r, "s2c-tnt-reg-xml-idp-nm")

			valRs, valRsErr := val.GetIdpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm)
			if valRsErr != nil {
				error.IntSrv(ctx, rw, valRsErr)
				return
			}

			if ! valRs[0].IdpNmOk {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-idp-nm-taken", "idpNm", idpNm)}, data)

				return
			}

			xmlFile, _, xmlFileErr := r.FormFile("s2c-tnt-reg-xml-file")
			if xmlFileErr != nil {
				error.IntSrv(ctx, rw, xmlFileErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from xml form",
				slog.String("idpNm" , idpNm),
			)

			mtdRaw, mtdRawErr := io.ReadAll(xmlFile)
			if mtdRawErr != nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-unreadable-xml")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: Get XML from file",
					slog.String("mtdRawErr" , mtdRawErr.Error()),
					slog.String("idpNm"     , idpNm),
				)

				return
			}
	
			mtd    := &gosaml2types.EntityDescriptor{}
			mtdErr := xml.Unmarshal(mtdRaw, mtd)
			if mtdErr != nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-read-metadata")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: unmarshal metadata into struct",
					slog.String("mtdErr" , mtdErr.Error()),
					slog.String("idpNm"  , idpNm),
				)

				return
			}

			if mtd.IDPSSODescriptor == nil {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-read-metadata")}, data)

				ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: the xml file did not contain IDPSSODescriptor elements",
					slog.String("idpNm" , idpNm),
				)

				return
			}

			lkd := len(mtd.IDPSSODescriptor.KeyDescriptors)

			ipcCrt    := make([][]byte    , lkd)
			cruNm     := make([]string    , lkd)
			ipcIncTs  := make([]time.Time , lkd)
			ipcExpTs  := make([]time.Time , lkd)

			for i, kds := range mtd.IDPSSODescriptor.KeyDescriptors {
				for _, x5c := range kds.KeyInfo.X509Data.X509Certificates {
					if strings.TrimSpace(x5c.Data) == "" {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-empty-cert")}, data)

						ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post:: x5c.Data is empty",
							slog.Int   ("i"   , i),
							slog.String("kds" , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					x5d, x5dErr := base64.StdEncoding.DecodeString(strings.TrimSpace(x5c.Data))
					if x5dErr != nil {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-decode-cert")}, data)

						slog.LogAttrs(ctx, slog.LevelError, "Post:: cannot decode x5c.Data",
							slog.String("x5dErr" , x5dErr.Error()),
							slog.Int   ("i"      , i),
							slog.String("kds"    , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					crt, crtErr := x509.ParseCertificate(x5d)
					if crtErr != nil {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-gen-cert")}, data)

						slog.LogAttrs(ctx, slog.LevelError, "Post:: cannot parse x5d",
							slog.String("crtErr" , crtErr.Error()),
							slog.Int   ("i"      , i),
							slog.String("kds"    , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					mCrt, mCrtErr := json.Marshal(crt)
					if mCrtErr != nil {
						notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.warning-input-marshal-cert")}, data)

						slog.LogAttrs(ctx, slog.LevelError, "Post:: marshal crt",
							slog.String("crtErr" , crtErr.Error()),
							slog.Int   ("i"      , i),
							slog.String("kds"    , fmt.Sprintf("%+v", kds)),
						)

						return
					}

					ipcCrt[i]   = mCrt
					cruNm[i]    = kds.Use
					ipcIncTs[i] = crt.NotBefore
					ipcExpTs[i] = crt.NotAfter
				}
			}

			sloUrl    := make([]string, len(mtd.IDPSSODescriptor.SingleLogoutServices))
			sloUrlBnd := make([]string, len(mtd.IDPSSODescriptor.SingleLogoutServices))

			for i, slo := range mtd.IDPSSODescriptor.SingleLogoutServices {
				sloUrl[i]    = slo.Location
				sloUrlBnd[i] = slo.Binding
			}

			ssoUrl    := make([]string, len(mtd.IDPSSODescriptor.SingleSignOnServices))
			ssoUrlBnd := make([]string, len(mtd.IDPSSODescriptor.SingleSignOnServices))

			for i, sso := range mtd.IDPSSODescriptor.SingleSignOnServices {
				ssoUrl[i]    = sso.Location
				ssoUrlBnd[i] = sso.Binding
			}

			postErr := PostIdp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, idpNm, mtd.EntityID, ipcCrt, cruNm, ipcIncTs, ipcExpTs, nil, sloUrl, sloUrlBnd, ssoUrl, ssoUrlBnd, data.User.AurNm, nil)
			if postErr != nil {
				error.IntSrv(ctx, rw, postErr)
				return
			}

			rw.Header().Set("HX-Trigger", `{"mod":{"target":"#s2c-tnt-inf-idp-form"}}`)

			notification.Toast(ctx, ssd.Logger, rw, r, "success" , &map[string]string{"Message" : data.T("web-core-auth-s2c-tnt-reg-xml-form.message-input-success", "idpNm", idpNm)}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::end")
}
