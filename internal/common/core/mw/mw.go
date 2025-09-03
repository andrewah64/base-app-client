package mware

import (
	"net/http"
)

import (
	"github.com/justinas/nosurf"
)

func CSRFHandler(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path    : "/",
		Secure  : true,
	})

	return csrfHandler
}

func ResponseHeaders (next http.Handler) http.Handler {
	return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request){
		rw.Header().Set("Strict-Transport-Security", "max-age=31557600; includeSubDomains; preload")
		rw.Header().Set("X-XSS-Protection"         , "0")
		rw.Header().Set("X-Frame-Options"          , "deny")
		rw.Header().Set("X-Content-Type-Options"   , "nosniff")
		rw.Header().Set("Referrer-Policy"          , "origin-when-cross-origin")

		next.ServeHTTP(rw, r)
	})
}
