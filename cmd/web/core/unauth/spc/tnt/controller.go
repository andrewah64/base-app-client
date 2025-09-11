package tnt

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
)

func Get (rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Get::get request info"))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Get::start")

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_spc_tnt_inf")

	spcRs, spcRsErr := GetSpc(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
	if spcRsErr != nil {
		error.IntSrv(ctx, rw, spcRsErr)
		return
	}

	s2cEntityId := spcRs[0].S2cEntityId
	spcSgnCrt   := strings.TrimRight(base64.StdEncoding.EncodeToString(spcRs[0].SpcSgnCrt), "=")
	spcEncCrt   := strings.TrimRight(base64.StdEncoding.EncodeToString(spcRs[0].SpcEncCrt), "=")
	s2cAcsUrl   := spcRs[0].S2cAcsUrl

	metadata := `<?xml version="1.0" encoding="UTF-8"?>
	               <EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="%s">
	                 <SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol" AuthnRequestsSigned="true" WantAssertionsSigned="true">
	                   <KeyDescriptor use="signing">
	                     <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
	                       <ds:X509Data>
	                         <ds:X509Certificate>%s</ds:X509Certificate>
	                       </ds:X509Data>
	                     </ds:KeyInfo>
	                   </KeyDescriptor>
	                   <KeyDescriptor use="encryption">
	                     <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
	                       <ds:X509Data>
	                         <ds:X509Certificate>%s</ds:X509Certificate>
	                       </ds:X509Data>
	                     </ds:KeyInfo>
	                   </KeyDescriptor>
	                   <AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="%s" index="0" isDefault="true"/>
	                 </SPSSODescriptor>
	               </EntityDescriptor>`

	rw.Header().Set("Content-Type", "application/samlmetadata+xml")

	rw.Write([]byte(fmt.Sprintf(metadata, s2cEntityId, spcSgnCrt, spcEncCrt, s2cAcsUrl)))
}
