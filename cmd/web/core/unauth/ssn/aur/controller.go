package aur

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

import (
	   "github.com/andrewah64/base-app-client/internal/common/core/password"
	cs "github.com/andrewah64/base-app-client/internal/common/core/session"
	   "github.com/andrewah64/base-app-client/internal/common/core/token"
	   "github.com/andrewah64/base-app-client/internal/web/core/error"
	   "github.com/andrewah64/base-app-client/internal/web/core/passkey"
	ws "github.com/andrewah64/base-app-client/internal/web/core/session"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/html"
	   "github.com/andrewah64/base-app-client/internal/web/core/ui/notification"
)

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

import (
	"github.com/russellhaering/gosaml2"
	"github.com/russellhaering/goxmldsig"
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

	cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_ssn_aur_reg")

	aumRs, aumRsErr := GetAumInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if aumRsErr != nil {
		error.IntSrv(ctx, rw, aumRsErr)
		return
	}

	if aumRs[0].Saml2S2i {
		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::IDP-initiated SAML2 authentication configured")

		rw.WriteHeader(http.StatusForbidden)

		return
	}

	if aumRs[0].Saml2S2s {
		s2sInfRs, s2sRsErr := GetS2sInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
		if s2sRsErr != nil {
			error.IntSrv(ctx, rw, s2sRsErr)

			return
		}

		if len(s2sInfRs) == 0 {
			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::no IDP information retrieved")

			rw.WriteHeader(http.StatusForbidden)

			return
		}

		roots := make([]*x509.Certificate, len(s2sInfRs[0].IpcCrt))

		for i, ipcCrt := range s2sInfRs[0].IpcCrt {
			crt, crtErr := x509.ParseCertificate(ipcCrt)
			if crtErr != nil {
				error.IntSrv(ctx, rw, crtErr)
				return
			}

			roots[i] = crt
		}

		idpCs := dsig.MemoryX509CertificateStore{
			Roots: roots,
		}

		var sp *saml2.SAMLServiceProvider

		sp = &saml2.SAMLServiceProvider {
			IdentityProviderSSOURL      : s2sInfRs[0].SsoUrl,
			IdentityProviderSSOBinding  : s2sInfRs[0].SsoBndNm,
			IdentityProviderIssuer      : s2sInfRs[0].IdpEntityId,
			ServiceProviderIssuer       : s2sInfRs[0].S2cEntityId,
			AssertionConsumerServiceURL : s2sInfRs[0].AcsEppPt,
			SignAuthnRequests           : true,
			AudienceURI                 : s2sInfRs[0].S2cEntityId,
			IDPCertificateStore         : &idpCs,
			SPKeyStore                  : dsig.RandomKeyStoreForTest(),
		}

		idpSsoUrl, idpSsoUrlErr := sp.BuildAuthURL("")
		if idpSsoUrlErr != nil {
			slog.LogAttrs(ctx, slog.LevelError, "Get::unable to generate URL to redirect to for SP-initiated SAML2 auth")

			error.IntSrv(ctx, rw, idpSsoUrlErr)

			return
		}

		http.Redirect(rw, r, idpSsoUrl, http.StatusSeeOther)

		return
	} else {
		p := r.URL.Query()

		if p.Has("ntf"){
			notification.Toast(ctx, ssd.Logger, rw, r, "info", &map[string]string{"Message" : data.T(p.Get("ntf"))} , data)
		}

		rs := make(map[string]any)

		rs["Aum"] = &aumRs

		if aumRs[0].AupcEnabled {
			pwdRs, pwdRsErr := GetPwdInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if pwdRsErr != nil {
				error.IntSrv(ctx, rw, pwdRsErr)
				return
			}

			rs["Pwd"] = &pwdRs
		}

		data.ResultSet = &rs

		html.Tmpl(ctx, ssd.Logger, rw, r, "core/unauth/ssn/aur/content", http.StatusOK, &data)

		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::end")
	}
}

func Post(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := cs.FromContext(ctx)
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
			aurNm  := strings.ToLower(form.VText(r, "ssn-aur-reg-aupc-aur-nm"))
			aurPwd := form.VText(r, "ssn-aur-reg-aupc-aur-pwd")

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
				slog.String("aurNm", aurNm),
			)

			if strings.TrimSpace(aurNm) == ""  || strings.TrimSpace(aurPwd) == "" {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-ssn-aur-reg-aupc-form.error-input-unexpected")}, data)

				return
			}

			cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_ssn_aur_reg")

			aurRs, aurRsErr := GetAurPwdInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurRsErr != nil {
				error.IntSrv(ctx, rw, aurRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::retrieve details of user attempting to authenticate",
				slog.Int("len(aurRs)", len(aurRs)),
			)

			switch len(aurRs) {
			case 1:
				if password.CheckHash(aurPwd, aurRs[0].AurHshPw) {
					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::username/password combination is valid",
						slog.String("aurNm", aurNm),
					)

					if aurRs[0].AupcMfaEnabled {
						nncNonce, nncNonceErr := token.Token(16)
						if nncNonceErr != nil {
							error.IntSrv(ctx, rw, nncNonceErr)
							return
						}

						nncNonce = fmt.Sprintf("%v%v", aurRs[0].AurId, nncNonce)

						PostNnc (&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurRs[0].AurId, nncNonce, time.Now().Add(time.Minute * 5), nil)

						rw.Header().Set("HX-Location", fmt.Sprintf(`{"path":"/web/core/unauth/otp/ssn/aur/%v", "target":"#main", "select":"#content"}`, nncNonce))
					} else {
						cookieExpiry := time.Now().Add(aurRs[0].SsnDn)

						ssnErr := ws.Begin(&ctx, ssd.Logger, ssd.Conn, rw, aurRs[0].AurId, cookieExpiry)
						if ssnErr != nil {
							error.IntSrv(ctx, rw, ssnErr)
							return
						}

						ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::redirect to user's home page",
							slog.String("aurRs[0].EppPt", aurRs[0].EppPt),
						)

						rw.Header().Set("HX-Redirect", aurRs[0].EppPt)
					}
				}

				fallthrough
			default:
				if len(aurRs) <= 1 {
					ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::username/password combination is invalid",
						slog.String("aurNm", aurNm),
					)

					msgs := []string{data.T("web-core-unauth-ssn-aur-reg-aupc-form.error-input-aur-nm-pwd-vld")}

					notification.Vrl(ctx, ssd.Logger, rw, r,
						data.T("web-core-unauth-ssn-aur-reg-page.title"),
						data.T("web-core-unauth-ssn-aur-reg-aupc-form.title-warning-singular", "n", strconv.Itoa(len(msgs))),
						data.T("web-core-unauth-ssn-aur-reg-aupc-form.title-warning-plural"  , "n", strconv.Itoa(len(msgs))),
						&msgs,
						data,
					)

					return
				} else {
					error.IntSrv(ctx, rw, fmt.Errorf("Post::details of more than one user found"))
					return
				}
			}
		case "pky-atn-bgn":
			aurNm := strings.ToLower(form.VText(r, "ssn-aur-reg-pky-aur-nm"))

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::get data from form",
				slog.String("aurNm", aurNm),
			)

			if strings.TrimSpace(aurNm) == "" {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-ssn-aur-reg-pky-form.error-input-unexpected")}, data)

				return
			}

			cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_ssn_aur_reg")

			aurNmRs, aurNmRsErr := GetAurNmInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurNmRsErr != nil {
				error.IntSrv(ctx, rw, aurNmRsErr)
				return
			}

			if !aurNmRs[0].AurNmPass {
				notification.Toast(ctx, ssd.Logger, rw, r, "error" , &map[string]string{"Message" : data.T("web-core-unauth-ssn-aur-reg-pky-form.error-input-aur-nm")}, data)

				return
			}

			pkyAur, pkyAurErr := GetPkyAur(&ctx, ssd.Conn, ssd.TntId, aurNm)
			if pkyAurErr != nil {
				error.IntSrv(ctx, rw, pkyAurErr)
				return
			}

			aukcAtnInfRs, aukcAtnInfRsErr := GetAukcAtnInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if aukcAtnInfRsErr != nil {
				error.IntSrv(ctx, rw, aukcAtnInfRsErr)
				return
			}

			var pkh []protocol.PublicKeyCredentialHints = make([]protocol.PublicKeyCredentialHints, len(aukcAtnInfRs[0].PkhNm))
			for i, v := range aukcAtnInfRs[0].PkhNm {
				pkh[i] = protocol.PublicKeyCredentialHints(v)
			}

			atnOpts := []webauthn.LoginOption {
				webauthn.WithUserVerification(protocol.UserVerificationRequirement(aukcAtnInfRs[0].PuvNm)),
				webauthn.WithAssertionPublicKeyCredentialHints(pkh),
			}

			c, s, blErr := passkey.WebAuthn(&ctx, ssd.Logger, ssd.TntId).BeginLogin(pkyAur, atnOpts...)
			if blErr != nil {
				error.IntSrv(ctx, rw, blErr)
				return
			}

			sd, sdErr := json.Marshal(s)
			if sdErr != nil {
				error.IntSrv(ctx, rw, sdErr)
				return
			}

			regErr := PostPls(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, s.Challenge, sd, nil)
			if regErr != nil {
				error.IntSrv(ctx, rw, regErr)
				return
			}

			rw.Header().Set("Content-Type", "application/json")

			json.NewEncoder(rw).Encode(c)

		case "pky-atn-end":
			aurNm := strings.ToLower(strings.TrimSpace(strings.Split(r.PathValue("aum"), "/")[1]))

			cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_ssn_aur_reg")

			pkyAur, pkyAurErr := GetPkyAur(&ctx, ssd.Conn, ssd.TntId, aurNm)
			if pkyAurErr != nil {
				error.IntSrv(ctx, rw, pkyAurErr)
				return
			}

			pR, pRErr := protocol.ParseCredentialRequestResponse(r);
			if pRErr != nil {
				error.IntSrv(ctx, rw, pRErr)
				return
			}

			plsRs, plsRsErr := GetPlsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm, pR.Response.CollectedClientData.Challenge)
			if plsRsErr != nil {
				error.IntSrv(ctx, rw, plsRsErr)
				return
			}

			var sd webauthn.SessionData
			sdErr := json.Unmarshal([]byte(plsRs[0].PlsJs), &sd)
			if sdErr != nil {
				error.IntSrv(ctx, rw, sdErr)
				return
			}

			_, brErr := passkey.WebAuthn(&ctx, ssd.Logger, ssd.TntId).ValidateLogin(pkyAur, sd, pR)
			if brErr != nil {
				error.IntSrv(ctx, rw, brErr)
				return
			}

			aurRs, aurRsErr := GetAurPkyInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurRsErr != nil {
				error.IntSrv(ctx, rw, aurRsErr)
				return
			}

			cookieExpiry := time.Now().Add(aurRs[0].SsnDn)

			ssnErr := ws.Begin(&ctx, ssd.Logger, ssd.Conn, rw, aurRs[0].AurId, cookieExpiry)
			if ssnErr != nil {
				error.IntSrv(ctx, rw, ssnErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::redirect to user's home page",
				slog.String("aurRs[0].EppPt", aurRs[0].EppPt),
			)

			rw.Header().Set("Content-Type", "application/json")

			json.NewEncoder(rw).Encode(struct {EppPt string `json:"eppPt"`}{EppPt : aurRs[0].EppPt})
	}
}
