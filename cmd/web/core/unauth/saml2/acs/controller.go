package tnt

import (
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
)

import (
	"github.com/russellhaering/gosaml2"
	"github.com/russellhaering/goxmldsig"
)

func Post (rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Post::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::start")

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_saml2_acs_inf")

	acsInfRs, acsInfRsErr := GetAcsInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if acsInfRsErr != nil {
		error.IntSrv(ctx, rw, acsInfRsErr)
		return
	}

	if len(acsInfRs) == 0 {
		fmt.Printf("\n\n A \n\n")
		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::no acs information retrieved")

		rw.WriteHeader(http.StatusForbidden)

		return
	}

	roots := make([]*x509.Certificate, len(acsInfRs[0].IpcCrt))

	for i, ipcCrt := range acsInfRs[0].IpcCrt {
		crt, crtErr := x509.ParseCertificate(ipcCrt)
		if crtErr != nil {
			error.IntSrv(ctx, rw, crtErr)
			return
		}

		roots[i] = crt
	}

	cs := dsig.MemoryX509CertificateStore{
		Roots: roots,
	}

	var sp *saml2.SAMLServiceProvider

	sp = &saml2.SAMLServiceProvider {
		IdentityProviderSSOURL      : acsInfRs[0].SsoUrl,
		IdentityProviderSSOBinding  : acsInfRs[0].SsoBndNm,
		IdentityProviderIssuer      : acsInfRs[0].IdpEntityId,
		ServiceProviderIssuer       : acsInfRs[0].S2cEntityId,
		AssertionConsumerServiceURL : acsInfRs[0].AcsEppPt,
		SignAuthnRequests           : true,
		AudienceURI                 : acsInfRs[0].S2cEntityId,
		IDPCertificateStore         : &cs,
		SPKeyStore                  : dsig.RandomKeyStoreForTest(),
	}

	astInf, astInfErr := sp.RetrieveAssertionInfo(r.FormValue("SAMLResponse"))
	if astInfErr != nil {
		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::error retrieving assertion information")

		rw.WriteHeader(http.StatusForbidden)

		return
	}

	if astInf.WarningInfo.InvalidTime {
		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::invalid time")

		rw.WriteHeader(http.StatusForbidden)

		return
	}

	if astInf.WarningInfo.NotInAudience {
		ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::not in audience",
			slog.String("acsInfRs[0].S2cEntityId", acsInfRs[0].S2cEntityId),
		)

		rw.WriteHeader(http.StatusForbidden)

		return
	}

	aurEa := astInf.Values["http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"].Values[0].Value

	if _, aurEaErr := mail.ParseAddress(aurEa); aurEaErr != nil {
		ssd.Logger.LogAttrs(ctx, slog.LevelError, "Post::email validity check has failed",
			slog.String("aurEa"            , aurEa),
			slog.String("aurEaErr.Error()" , aurEaErr.Error()),
		)

		rw.WriteHeader(http.StatusForbidden)

		return
	}

	//get user info based on email address - aur_inf

	//if they aren't registered, reg aur, get user info based on email address - reg_aur, aur_inf

	//begin http session

	//redirect them to the app

	fmt.Printf("\n\n Hang on, did we get here? \n\n")
}
