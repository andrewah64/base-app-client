package otp

import (
	"encoding/base32"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

import (
	cs "github.com/andrewah64/base-app-client/internal/common/core/session"
	   "github.com/andrewah64/base-app-client/internal/common/core/tenant"
	   "github.com/andrewah64/base-app-client/internal/common/core/token"
	   "github.com/andrewah64/base-app-client/internal/web/core/error"
	ws "github.com/andrewah64/base-app-client/internal/web/core/session"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

import (
	"github.com/pquerna/otp/totp"
)

func Get(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := cs.FromContext(ctx)
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

	nncNonce := r.PathValue("id")

	cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_otp_ssn_aur_mod")

	nncRs, nncRsErr := GetNncInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, nncNonce)
	if nncRsErr != nil {
		error.IntSrv(ctx, rw, nncRsErr)
		return
	}

	if len(nncRs) != 1 {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::received %v records", len(nncRs)))
		return
	}

	if nncRs[0].NncEnabled && nncRs[0].OtpEnabled {
		data.ResultSet = &map[string]any{"AurId" : nncRs[0].AurId, "NncNonce" : nncNonce}

		html.Tmpl(ctx, ssd.Logger, rw, r, "core/unauth/otp/ssn/aur/content", http.StatusOK, &data)
	} else {
		if ! nncRs[0].NncEnabled {
			rw.Header().Set("HX-Location", `{"path":"/", "target":"#main", "select":"#content", "values":{"ntf": "web-core-unauth-otp-ssn-aur-mod-form.error-timeout"}}`)

			return
		}

		if ! nncRs[0].OtpEnabled {
			otpId , otpIdErr := token.Token(32)
			if otpIdErr != nil {
				error.IntSrv(ctx, rw, otpIdErr)
				return
			}

			otpTotpSecret, otpTotpSecretErr := totp.Generate(totp.GenerateOpts{
				Issuer     : tenant.Origin(r),
				AccountName: nncRs[0].AurNm,
			})

			if otpTotpSecretErr != nil {
				error.IntSrv(ctx, rw, otpTotpSecretErr)
				return
			}

			otpSecret := otpTotpSecret.Secret()

			regErr := PostOtp(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, nncRs[0].AurId, otpId, otpSecret, nil)
			if regErr != nil {
				error.IntSrv(ctx, rw, regErr)
				return
			}

			rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"/web/core/unauth/otp/aur/%v", "target":"#main", "select":"#content"}`, otpId))

			return
		}
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
}

func Post(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := cs.FromContext(ctx)
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

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr) 
		return
	}

	nncNonce := r.PathValue("id")
	aurId    := form.VInt (r, "otp-ssn-aur-mod-aur-id")
	otpCd    := form.VText(r, "otp-ssn-aur-mod-otp-cd")

	cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_otp_ssn_aur_mod")

	aurRs, aurRsErr := GetAurInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId, nncNonce)
	if aurRsErr != nil {
		error.IntSrv(ctx, rw, aurRsErr)
		return
	}

	otpSecret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(aurRs[0].OtpSecret))

	if totp.Validate(otpCd, otpSecret) {
		cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_ssn_aur_reg")

		cookieExpiry := time.Now().Add(aurRs[0].AurSsnDn)

		ssnErr := ws.Begin(&ctx, ssd.Logger, ssd.Conn, rw, aurId, cookieExpiry)
		if ssnErr != nil {
			error.IntSrv(ctx, rw, ssnErr)
			return
		}

		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::redirect to user's home page",
			slog.String("aurRs[0].EppPt", aurRs[0].EppPt),
		)

		rw.Header().Set("HX-Redirect", aurRs[0].EppPt)
	} else {
		notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-otp-ssn-aur-mod-form.error-otp-cd")}, data)
	}
}
