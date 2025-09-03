package otp

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"image/png"
	"log/slog"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/tenant"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

import (
	"github.com/pquerna/otp/totp"
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

	otpId := r.PathValue("id")

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_otp_aur_inf")

	otpInfRs, otpInfRsErr := GetOtpAurInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, otpId)
	if otpInfRsErr != nil {
		error.IntSrv(ctx, rw, otpInfRsErr)
		return
	}

	if len(otpInfRs) != 1 {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::%v records were found", len(otpInfRs)))
		return
	}

	otpKey, otpKeyErr := totp.Generate(totp.GenerateOpts{
		Issuer     : tenant.Origin(r),
		AccountName: otpInfRs[0].AurNm,
		Secret     : []byte(otpInfRs[0].OtpSecret),
	})
	if otpKeyErr != nil {
		error.IntSrv(ctx, rw, otpKeyErr)
		return
	}

	var buf bytes.Buffer
	img, err := otpKey.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	imgStr := base64.StdEncoding.EncodeToString(buf.Bytes())

	data.ResultSet = &map[string]any{
		"AurId"  : &(otpInfRs[0].AurId),
		"OtpImg" : &imgStr,
		"OtpId"  : &otpId,
	}

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/unauth/otp/aur/content", http.StatusOK, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
}

func Post (rw http.ResponseWriter, r *http.Request) {
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

	otpId := r.PathValue("id")
	aurId := form.VInt (r, "otp-aur-reg-aur-id")
	otpCd := form.VText(r, "otp-aur-reg-otp-cd")

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_otp_aur_mod")

	otpInfRs, otpInfRsErr := GetOtpInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId, otpId)
	if otpInfRsErr != nil {
		error.IntSrv(ctx, rw, otpInfRsErr)
		return
	}

	otpSecret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(otpInfRs[0].OtpSecret))

	if totp.Validate(otpCd, otpSecret) {
		postErr := PostOtp (&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurId, otpId, nil)
		if postErr != nil {
			error.IntSrv(ctx, rw, postErr)
			return
		}

		rw.Header().Set("HX-Location", `{"path":"/", "target":"#main", "select":"#content", "values":{"ntf": "web-core-unauth-otp-aur-mod-form.message-otp-cd-success"}}`)
	} else {
		notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-otp-aur-mod-form.error-otp-cd")}, data)
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::end")
}
