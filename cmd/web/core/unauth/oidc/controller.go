package oidc

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/coreos/go-oidc/v3/oidc"
)

import (
	cs "github.com/andrewah64/base-app-client/internal/common/core/session"
	t  "github.com/andrewah64/base-app-client/internal/common/core/token"
	e  "github.com/andrewah64/base-app-client/internal/web/core/error"
	ws "github.com/andrewah64/base-app-client/internal/web/core/session"
)

func callbackCookie(rw http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(rw, c)
}

func Call(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := cs.FromContext(ctx)
	if ! ok {
		e.IntSrv(ctx, rw, fmt.Errorf("Call::get request info"))
		return
	}

	ocpNm := r.PathValue("nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Call::start",
		slog.String("ocpNm", ocpNm),
	)

	cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_oidc_call_inf")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Call::get OIDC provider details",
		slog.Int   ("ssd.TntId" , ssd.TntId),
		slog.String("ocpNm"     , ocpNm),
	)

	callInfRs, callInfRsErr := GetCallInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, ocpNm)
	if callInfRsErr != nil {
		e.IntSrv(ctx, rw, callInfRsErr)
		return
	}

	if len(callInfRs) != 1 {
		e.IntSrv(ctx, rw, fmt.Errorf("%v records of OIDC information about %v were retrieved", len(callInfRs), ocpNm))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Call::get OIDC details",
		slog.String("callInfRs[0].OccClientId" , callInfRs[0].OccClientId),
		slog.String("callInfRs[0].OccCbUrl"    , callInfRs[0].OccCbUrl),
		slog.Any   ("callInfRs[0].OcsNm"       , callInfRs[0].OcsNm),
	)

	stTkn, stTknErr := t.Token(16)
	if stTknErr != nil {
		e.IntSrv(ctx, rw, stTknErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Call::get 'state' value",
		slog.String("stTkn" , stTkn),
	)

	ncTkn, ncTknErr := t.Token(16)
	if ncTknErr != nil {
		e.IntSrv(ctx, rw, ncTknErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Call::get 'nonce' value",
		slog.String("ncTkn" , stTkn),
	)

	callbackCookie(rw, r, "state", stTkn)
	callbackCookie(rw, r, "nonce", ncTkn)

	provider, providerErr := oidc.NewProvider(ctx, callInfRs[0].OccUrl)
	if providerErr != nil {
		e.IntSrv(ctx, rw, providerErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Call::get OIDC provider's endpoint",
		slog.Any("provider.Endpoint()" , provider.Endpoint()),
	)

	config := oauth2.Config{
		ClientID     : callInfRs[0].OccClientId,
		Endpoint     : provider.Endpoint(),
		RedirectURL  : strings.Replace(callInfRs[0].OccCbUrl, "{nm}", ocpNm, 1),
		Scopes       : callInfRs[0].OcsNm,
	}

	http.Redirect(rw, r, config.AuthCodeURL(stTkn, oidc.Nonce(ncTkn)), http.StatusFound)
}

func Callback(rw http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ssd, ok := cs.FromContext(ctx)
	if ! ok {
		e.IntSrv(ctx, rw, fmt.Errorf("Callback::get request info"))
		return
	}

	ocpNm := r.PathValue("nm")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::start",
		slog.String("ocpNm", ocpNm),
	)

	stTkn, stTknErr := r.Cookie("state")
	if stTknErr != nil {
		e.IntSrv(ctx, rw, stTknErr)
		return
	}

	stUrl := r.URL.Query().Get("state")

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get 'state' from (1) the URL and (2) the 'state' cookie",
		slog.String("stUrl", stUrl),
		slog.String("stTkn", stTkn.Value),
	)

	if stUrl != stTkn.Value {
		e.IntSrv(ctx, rw, fmt.Errorf("Callback::URL & cookie states do not match"))
		return
	}

	cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_oidc_callback_mod")

	cbInfRs, cbInfRsErr := GetCallbackInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, ocpNm)
	if cbInfRsErr != nil {
		e.IntSrv(ctx, rw, cbInfRsErr)
		return
	}

	if len(cbInfRs) != 1 {
		e.IntSrv(ctx, rw, fmt.Errorf("%v records of OIDC information about %v were retrieved", len(cbInfRs), ocpNm))
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get OIDC details",
		slog.String("cbInfRs[0].OccClientId" , cbInfRs[0].OccClientId),
		slog.String("cbInfRs[0].OccUrl"      , cbInfRs[0].OccUrl),
		slog.String("cbInfRs[0].OccCbUrl"    , cbInfRs[0].OccCbUrl),
	)

	provider, providerErr := oidc.NewProvider(ctx, cbInfRs[0].OccUrl)
	if providerErr != nil {
		e.IntSrv(ctx, rw, providerErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get OIDC provider's endpoint",
		slog.Any("provider.Endpoint()" , provider.Endpoint()),
	)

	config := oauth2.Config{
		ClientID     : cbInfRs[0].OccClientId,
		ClientSecret : cbInfRs[0].OccClientSecret,
		Endpoint     : provider.Endpoint(),
		RedirectURL  : strings.Replace(cbInfRs[0].OccCbUrl, "{nm}", ocpNm, 1),
	}

	oauth2Tkn, oauth2TknErr := config.Exchange(ctx, r.URL.Query().Get("code"))
	if oauth2TknErr != nil {
		http.Redirect(rw, r, "/", http.StatusFound)
		return
	}

	rawIdTkn, ok := oauth2Tkn.Extra("id_token").(string)
	if !ok {
		e.IntSrv(ctx, rw, fmt.Errorf("Callback::no id_token field in oauth2 token"))
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: cbInfRs[0].OccClientId,
	}

	verifier := provider.Verifier(oidcConfig)

	idTkn, idTknErr := verifier.Verify(ctx, rawIdTkn)
	if idTknErr != nil {
		e.IntSrv(ctx, rw, idTknErr)
		return
	}

	nonce, nonceErr := r.Cookie("nonce")
	if nonceErr != nil {
		e.IntSrv(ctx, rw, nonceErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get 'nonce' from (1) the ID token and (2) the 'nonce' cookie",
		slog.String("nonce.Value" , nonce.Value),
		slog.String("idTkn.Nonce" , idTkn.Nonce),
	)

	if idTkn.Nonce != nonce.Value {
		e.IntSrv(ctx, rw, fmt.Errorf("Callback::URL & cookie nonces do not match"))
		return
	}

	oauth2Tkn.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Tkn, new(json.RawMessage)}

	if claimsErr := idTkn.Claims(&resp.IDTokenClaims); claimsErr != nil {
		e.IntSrv(ctx, rw, claimsErr)
		return
	}

	var idTknMap map[string]interface{}

	jsonErr := json.Unmarshal(*resp.IDTokenClaims, &idTknMap)
	if jsonErr != nil {
		e.IntSrv(ctx, rw, jsonErr)
		return
	}

	aurEa := idTknMap["email"].(string)

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get user's details from OIDC provider",
		slog.String("aurEa" , aurEa),
	)

	aurInfRs, aurInfRsErr := GetAurInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurEa)
	if aurInfRsErr != nil {
		e.IntSrv(ctx, rw, aurInfRsErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get user's details",
		slog.Int("len(aurInfRs)" , len(aurInfRs)),
	)

	switch len(aurInfRs) {
		case 0:
			regErr := RegAur (&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurEa, nil)
			if regErr != nil {
				e.IntSrv(ctx, rw, regErr)
				return
			}

			aurInfRs, aurInfRsErr = GetAurInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurEa)
			if aurInfRsErr != nil {
				e.IntSrv(ctx, rw, aurInfRsErr)
				return
			}

			ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Callback::get user's details after registering them",
				slog.Int("len(aurInfRs)" , len(aurInfRs)),
			)
		case 1:
			//the user was already registered
		default:
			e.IntSrv(ctx, rw, fmt.Errorf("Callback::%v records were returned when only 0 or 1 are expected", len(aurInfRs)))
			return
	}

	cookieExpiry := time.Now().Add(aurInfRs[0].SsnDn)

	cs.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_ssn_aur_reg")

	ssnErr := ws.Begin(&ctx, ssd.Logger, ssd.Conn, rw, aurInfRs[0].AurId, cookieExpiry)
	if ssnErr != nil {
		e.IntSrv(ctx, rw, ssnErr)
		return
	}

	ssd.Logger.LogAttrs(ctx, slog.LevelDebug, "Post::redirect to user's home page",
		slog.String("aurInfRs[0].EppPt", aurInfRs[0].EppPt),
	)

	http.Redirect(rw, r, aurInfRs[0].EppPt, http.StatusFound)
}
