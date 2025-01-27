package api

import (
	"net/http"

	"github.com/danielllmuniz/go-http-server/internal/jsonutils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserId") {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"error": "muse be logged in",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
