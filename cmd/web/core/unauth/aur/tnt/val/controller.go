package val

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/password"
	"github.com/andrewah64/base-app-client/internal/common/core/session"
	"github.com/andrewah64/base-app-client/internal/web/core/error"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/form"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
)

func Get(rw http.ResponseWriter, r *http.Request) {
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

	pfErr := r.ParseForm()
	if pfErr != nil {
		error.IntSrv(ctx, rw, pfErr)
		return
	}

	session.Identity(&ctx, ssd.Logger, ssd.Conn, "role_web_core_unauth_aur_tnt_reg")

	switch r.PathValue("id") {
		case "aupc-aur-ea":
			aurEa := form.VText (r, "aur-tnt-reg-aupc-aur-ea")

			_, aurEaErr := mail.ParseAddress(aurEa)

			aurEaRs, aurEaRsErr := GetAurEaInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurEa)
			if aurEaRsErr != nil {
				error.IntSrv(ctx, rw, aurEaRsErr)
				return
			}

			data.ResultSet = &map[string]any{
				"AurEaAvbPass" : aurEaRs[0].AurEaAvbPass,
				"AurEaVldPass" : aurEaErr == nil,
			}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/unauth/aur/tnt/fragment/valaupcaurea", http.StatusOK, &data)

			return
		case "aupc-aur-nm":
			aurNm := form.VText (r, "aur-tnt-reg-aupc-aur-nm")

			aurNmRs, aurNmRsErr := GetAurNmInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurNmRsErr != nil {
				error.IntSrv(ctx, rw, aurNmRsErr)
				return
			}

			data.ResultSet = &map[string]any{
				"AurNmLenPass" : aurNmRs[0].AurNmLenPass,
				"AurNmAvbPass" : aurNmRs[0].AurNmAvbPass,
			}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/unauth/aur/tnt/fragment/valaupcaurnm", http.StatusOK, &data)

			return
		case "aupc-aur-pwd":
			aurPwd  := form.VText (r, "aur-tnt-reg-aupc-aur-pwd")
			aurPwd2 := form.VText (r, "aur-tnt-reg-aupc-aur-pwd-2")

			pwdRs, pwdRsErr := GetPwdInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId)
			if pwdRsErr != nil {
				error.IntSrv(ctx, rw, pwdRsErr)
				return
			}

			aurPwdLenPass, aurPwdSymPass, aurPwdNumPass := password.Validate(aurPwd, pwdRs[0].AurPwdMinLen, pwdRs[0].AurPwdMaxLen, pwdRs[0].AurPwdIncSym, pwdRs[0].AurPwdIncNum)

			data.ResultSet = &map[string]any{
				"AurPwdLenPass" : aurPwdLenPass,
				"AurPwdSymPass" : aurPwdSymPass,
				"AurPwdNumPass" : aurPwdNumPass,
				"AurPwd2Pass"   : aurPwd == aurPwd2,
			}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/unauth/aur/tnt/fragment/valaupcaurpwd", http.StatusOK, &data)

			return
		case "aupc-aur-pwd-2":
			aurPwd  := form.VText (r, "aur-tnt-reg-aupc-aur-pwd")
			aurPwd2 := form.VText (r, "aur-tnt-reg-aupc-aur-pwd-2")

			data.ResultSet = &map[string]any{
				"AurPwd2Pass" : aurPwd == aurPwd2,
			}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/unauth/aur/tnt/fragment/valaupcaurpwd2", http.StatusOK, &data)
		case "pky-aur-nm":
			aurNm := form.VText (r, "aur-tnt-reg-pky-aur-nm")

			aurNmRs, aurNmRsErr := GetAurNmInf(&ctx, ssd.Logger, ssd.Conn, ssd.TntId, aurNm)
			if aurNmRsErr != nil {
				error.IntSrv(ctx, rw, aurNmRsErr)
				return
			}

			data.ResultSet = &map[string]any{
				"AurNmLenPass" : aurNmRs[0].AurNmLenPass,
				"AurNmAvbPass" : aurNmRs[0].AurNmAvbPass,
			}

			html.Fragment(ctx, ssd.Logger, rw, r, "core/unauth/aur/tnt/fragment/valpkeyaurnm", http.StatusOK, &data)

			return
	}
}
