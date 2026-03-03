package AuthController

import (
	"github.com/rnschulenburg/gowrite-api-go/App/Services/AuthService"
	"github.com/rnschulenburg/gowrite-api-go/Core/Http"
	"net/http"
	"os"
	"time"
)

func Refresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "missing refresh token", 401)
		return
	}

	tokens, err := AuthService.RefreshSession(
		r.Context(),
		cookie.Value,
	)
	if err != nil {
		http.Error(w, "invalid refresh token", 401)
		return
	}

	isProd := os.Getenv("AppEnv") == "prod"

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	response := map[string]interface{}{
		"authToken":  tokens.AccessToken,
		"expiration": tokens.ExpiresAt.Unix(),
	}

	Http.JsonResponse(w, nil, "success", "auth.refresh", response)
}
