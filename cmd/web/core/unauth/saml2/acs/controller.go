package tnt

import (
	//"encoding/base64"
	"fmt"
	//"log/slog"
	"net/http"
	//"strings"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
)

func Post (rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ssd, ok := session.FromContext(ctx)
	if ! ok {
		error.IntSrv(ctx, rw, fmt.Errorf("Acs::get request info"))
		return
	}

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr) 
		return
	}

	fmt.Printf("\n\n r.Form : %+v \n\n", r.Form)

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_saml2_acs_inf")
}
