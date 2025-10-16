package tnt

import (
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	//"strings"
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
	fmt.Printf("\n\n XXX \n\n")

	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::start")
/*
	data, ok := page.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request data"))
		return
	}
*/
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

	if acsInfRs[0].SloUrl == nil {
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
	} else {
		sp = &saml2.SAMLServiceProvider {
			IdentityProviderSSOURL      : acsInfRs[0].SsoUrl,
			IdentityProviderSSOBinding  : acsInfRs[0].SsoBndNm,
			IdentityProviderSLOURL      : *acsInfRs[0].SloUrl,
			IdentityProviderSLOBinding  : *acsInfRs[0].SloBndNm,
			IdentityProviderIssuer      : acsInfRs[0].IdpEntityId,
			ServiceProviderIssuer       : acsInfRs[0].S2cEntityId,
			AssertionConsumerServiceURL : acsInfRs[0].AcsEppPt,
			SignAuthnRequests           : true,
			AudienceURI                 : acsInfRs[0].S2cEntityId,
			IDPCertificateStore         : &cs,
			SPKeyStore                  : dsig.RandomKeyStoreForTest(),
		}
	}

	fmt.Printf("\n\n sp struct :: %+v \n\n", sp)
}
