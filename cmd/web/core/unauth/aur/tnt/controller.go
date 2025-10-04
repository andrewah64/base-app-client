package aur

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/password"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/common/core/tenant"
	"github.com/andrewah64/base-app-client/internal/common/core/token"
	"github.com/andrewah64/base-app-client/internal/common/core/validator"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
	"github.com/andrewah64/base-app-client/internal/web/core/passkey"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

import (
	"github.com/andrewah64/base-app-client/cmd/web/core/unauth/aur/tnt/val"
)

import (
	"github.com/pquerna/otp/totp"
)

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
)

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

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_aur_tnt_reg")

	aumRs, aumRsErr := GetAumInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aumRsErr != nil {
		error.IntSrv(ctx, rw, aumRsErr)
		return
	}
	
	rs := make(map[string]any)

	rs["Aum"] = &aumRs

	if aumRs[0].AupcEnabled {
		aupcRs, aupcRsErr := GetAupcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
		if aupcRsErr != nil {
			error.IntSrv(ctx, rw, aupcRsErr)
			return
		}

		rs["Aupc"] = &aupcRs
	}

	if aumRs[0].AukcEnabled {
		aukcRs, aukcRsErr := GetAukcInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
		if aukcRsErr != nil {
			error.IntSrv(ctx, rw, aukcRsErr)
			return
		}

		rs["Aukc"] = &aukcRs
	}

	data.ResultSet = &rs

	html.Tmpl(ctx, ssd.Logger, rw, r, "core/unauth/aur/tnt/content", http.StatusOK, &data)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
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

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr) 
		return
	}

	switch strings.Split(r.PathValue("aum"), "/")[0] {
		case "aupc":
			aurNm  := strings.ToLower(form.VText (r, "aur-tnt-reg-aupc-aur-nm"))
			aurEa  := form.VText (r, "aur-tnt-reg-aupc-aur-ea")
			aurPw  := form.VText (r, "aur-tnt-reg-aupc-aur-pwd")
			aurPw2 := form.VText (r, "aur-tnt-reg-aupc-aur-pwd-2")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
				slog.String("aurNm" , aurNm),
				slog.Any   ("aurEa", aurEa),
			)

			if validator.Blank(aurNm)||validator.Blank(aurEa)||validator.Blank(aurPw)||validator.Blank(aurPw2) {
				ssd.Logger.LogAttrs(ctx, slog.LevelError, "Post::aupc::front end mandatory field checks have failed",
					slog.Bool("validator.Blank(aurNm)" , validator.Blank(aurNm)),
					slog.Bool("validator.Blank(aurEa)" , validator.Blank(aurEa)),
					slog.Bool("validator.Blank(aurPw)" , validator.Blank(aurPw)),
					slog.Bool("validator.Blank(aurPw2)", validator.Blank(aurPw2)),
				)

				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-aupc-tab.error-input-unexpected")}, data)

				return
			}

			aurNmRs, aurNmRsErr := val.GetAurNmInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurNmRsErr != nil {
				error.IntSrv(ctx, rw, aurNmRsErr)
				return
			}

			if !aurNmRs[0].AurNmLenPass {
				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-aupc-tab.error-input-aur-nm-len")}, data)

				return
			}

			if !aurNmRs[0].AurNmAvbPass {
				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-aupc-tab.error-input-aur-nm-avb")}, data)

				return
			}

			aurEaRs, aurEaRsErr := val.GetAurEaInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurEa)
			if aurEaRsErr != nil {
				error.IntSrv(ctx, rw, aurEaRsErr)
				return
			}

			if !aurEaRs[0].AurEaAvbPass {
				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-aupc-tab.error-input-aur-ea-avb")}, data)

				return
			}

			if _, aurEaErr := mail.ParseAddress(aurEa); aurEaErr != nil {
				ssd.Logger.LogAttrs(ctx, slog.LevelError, "Post::front end email validity check has failed",
					slog.String("aurEa"            , aurEa),
					slog.String("aurEaErr.Error()" , aurEaErr.Error()),
				)

				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-aupc-tab.error-input-aur-ea-vld")}, data)

				return
			}

			if aurPw != aurPw2 {
				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-aupc-tab.error-input-aur-pwd-2-vld")}, data)

				return
			}

			aurHshPw, pErr := password.Hash(aurPw)
			if pErr != nil {
				error.IntSrv(ctx, rw, pErr)
				return
			}

			session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_aur_tnt_reg")

			mfaRs, mfaRsErr := GetMfaInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if mfaRsErr != nil {
				error.IntSrv(ctx, rw, mfaRsErr)
				return
			}

			if mfaRs[0].AupcMfaEnabled {
				otpId , otpIdErr := token.Token(32)
				if otpIdErr != nil {
					error.IntSrv(ctx, rw, otpIdErr)
					return
				}

				otpTotpSecret, otpTotpSecretErr := totp.Generate(totp.GenerateOpts{
					Issuer     : tenant.Origin(r),
					AccountName: aurNm,
				})
				if otpTotpSecretErr != nil {
					error.IntSrv(ctx, rw, otpTotpSecretErr)
					return
				}

				otpSecret := otpTotpSecret.Secret()

				regErr := PostPwAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, aurHshPw, aurEa, &otpId, &otpSecret, nil)
				if regErr != nil {
					error.IntSrv(ctx, rw, regErr)
					return
				}

				rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"/web/core/unauth/otp/aur/%v", "target":"#main", "select":"#content"}`, otpId))
			} else {
				regErr := PostPwAur(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, aurHshPw, aurEa, nil, nil, nil)
				if regErr != nil {
					error.IntSrv(ctx, rw, regErr)
					return
				}

				rw.Header().Set("HX-Location", `{"path":"/", "target":"#main", "select":"#content", "values":{"ntf": "web-core-unauth-aur-tnt-aupc-tab.message-success"}}`)
			}
		case "pky-reg-bgn":
			aurNm := strings.ToLower(form.VText(r, "aur-tnt-reg-pky-aur-nm"))

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
				slog.String("aurNm" , aurNm),
			)

			if validator.Blank(aurNm) {
				ssd.Logger.LogAttrs(ctx, slog.LevelError, "Post::pky::front end mandatory field checks have failed",
					slog.Bool("validator.Blank(aurNm)" , validator.Blank(aurNm)),
				)

				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-ssn-aur-reg-pky-form.error-input-unexpected")}, data)

				return
			}

			session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_aur_tnt_reg")

			aurNmRs, aurNmRsErr := val.GetAurNmInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurNmRsErr != nil {
				error.IntSrv(ctx, rw, aurNmRsErr)
				return
			}

			if !aurNmRs[0].AurNmLenPass {
				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-pky-tab.error-input-aur-nm-len")}, data)

				return
			}

			if !aurNmRs[0].AurNmAvbPass {
				notification.Show(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-aur-tnt-pky-tab.error-input-aur-nm-avb")}, data)

				return
			}

			user := &passkey.User{
				Id          : []byte(aurNm),
				Name        : aurNm,
				DisplayName : aurNm,
			}

			aukcRegInfRs, aukcRegInfRsErr := GetAukcRegInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if aukcRegInfRsErr != nil {
				error.IntSrv(ctx, rw, aukcRegInfRsErr)
				return
			}

			var pka protocol.ConveyancePreference = protocol.ConveyancePreference(aukcRegInfRs[0].PkaNm)

			var pkg []protocol.CredentialParameter = make([]protocol.CredentialParameter, len(aukcRegInfRs[0].PkgCd))
			for i, v := range aukcRegInfRs[0].PkgCd {
				pkg[i] = protocol.CredentialParameter{
					Type      : protocol.PublicKeyCredentialType,
					Algorithm : webauthncose.COSEAlgorithmIdentifier(v),
				}
			}

			var pkh []protocol.PublicKeyCredentialHints = make([]protocol.PublicKeyCredentialHints, len(aukcRegInfRs[0].PkhNm))
			for i, v := range aukcRegInfRs[0].PkhNm {
				pkh[i] = protocol.PublicKeyCredentialHints(v)
			}

			rrk := false

			regOpts := []webauthn.RegistrationOption {
				webauthn.WithAuthenticatorSelection(
					protocol.AuthenticatorSelection {
						AuthenticatorAttachment : protocol.AuthenticatorAttachment(aukcRegInfRs[0].PktNm),
						RequireResidentKey      : &rrk,
						ResidentKey             : protocol.ResidentKeyRequirement(aukcRegInfRs[0].PdcNm),
						UserVerification        : protocol.UserVerificationRequirement(aukcRegInfRs[0].PuvNm),
					},
				),
				webauthn.WithConveyancePreference(pka),
				webauthn.WithCredentialParameters(pkg),
				webauthn.WithPublicKeyCredentialHints(pkh),
			}

			c, s, brErr := passkey.WebAuthn(&ctx, ssd.Logger, ssd.TntId).BeginRegistration(user, regOpts...)
			if brErr != nil {
				error.IntSrv(ctx, rw, brErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::Begin passkey registration ceremony",
				slog.String("c" , fmt.Sprintf("%c", c)),
			)

			sd, sdErr := json.Marshal(s)
			if sdErr != nil {
				error.IntSrv(ctx, rw, sdErr)
				return
			}

			regErr := PostPrs(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, sd, nil)
			if regErr != nil {
				error.IntSrv(ctx, rw, regErr)
				return
			}

			rw.Header().Set("Content-Type", "application/json")

			json.NewEncoder(rw).Encode(c)
		case "pky-reg-end":
			aurNm := strings.ToLower(strings.TrimSpace(strings.Split(r.PathValue("aum"), "/")[1]))

			user := &passkey.User{
				Id          : []byte(aurNm),
				Name        : aurNm,
				DisplayName : aurNm,
			}

			session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_aur_tnt_reg")

			prsRs, prsRsErr := GetPrsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if prsRsErr != nil {
				error.IntSrv(ctx, rw, prsRsErr)
				return
			}

			var sd webauthn.SessionData
			sdErr := json.Unmarshal([]byte(prsRs[0].PrsJs), &sd)
			if sdErr != nil {
				error.IntSrv(ctx, rw, sdErr)
				return
			}

			c, brErr := passkey.WebAuthn(&ctx, ssd.Logger, ssd.TntId).FinishRegistration(user, sd,  r)
			if brErr != nil {
				error.IntSrv(ctx, rw, brErr)
				return
			}
			
			var t []string = make([]string, len(c.Transport))
			for i, v := range c.Transport {
				t[i] = string(v)
			}

			regErr := PostPkyAur(
				&ctx,
				ssd.Logger,
				ssd.Conn,
				ssd.TntId,
				aurNm,
				true,
				c.ID,
				c.PublicKey,
				c.AttestationType,
				t,
				c.Flags.UserPresent,
				c.Flags.UserVerified,
				c.Flags.BackupEligible,
				c.Flags.BackupState,
				c.Authenticator.AAGUID,
				int(c.Authenticator.SignCount),
				c.Authenticator.CloneWarning,
				string(c.Authenticator.Attachment),
				c.Attestation.ClientDataJSON,
				c.Attestation.ClientDataHash,
				c.Attestation.AuthenticatorData,
				c.Attestation.PublicKeyAlgorithm,
				c.Attestation.Object,
				nil,
			)
			if regErr != nil {
				error.IntSrv(ctx, rw, regErr)
				return
			}
		}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::end")

	return
}
