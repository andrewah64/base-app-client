package register

import (
	"log/slog"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/api/core/error"
	"github.com/andrewah64/base-app-client/internal/api/core/json"
)

const (
	usernameInput = "username"
)

const (
	credentialsTaken = "api-core-auth-users-register-form.warning-input-credentials-taken"
)

func Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	env := json.Envelope{
		"status"     : "available",
		"system_info": map[string]string{
			"environment": "tenant1",
			"version"    : "v1",
		},
	}

	jsErr := json.Write(&ctx, slog.Default(), rw, http.StatusOK, env, nil)
	if jsErr != nil {
		error.IntSrv(ctx, rw, jsErr)
	}
}
