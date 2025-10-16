package mware

import (
	"fmt"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/web/core/error"
)

func Recover (next http.Handler) http.Handler {
	return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request){
		defer func (){
			if err := recover(); err != nil {
				rw.Header().Set("Connection", "close")
				error.IntSrv(r.Context(), rw,
					fmt.Errorf("web/core/mw/recover middleware :: %s", err),
				)
			}
		} ()
		next.ServeHTTP(rw, r)
	})
}
