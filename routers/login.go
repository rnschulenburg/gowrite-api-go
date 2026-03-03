package routers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rnschulenburg/gowrite-api-go/App/Services/AuthService"
	"github.com/rnschulenburg/gowrite-api-go/Core/Http"
	"github.com/rnschulenburg/gowrite-api-go/Core/UserLogin"
)

func Login(w http.ResponseWriter, r *http.Request) {

	type LoginRequest struct {
		Password string `json:"password"`
		NickName string `json:"nickName"`
	}

	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Http.JsonErrorResponse(w, 400, "routers.Login")
		return
	}

	user, err := UserLogin.ByPassword(request.NickName, request.Password)
	if err != nil {
		Http.JsonErrorResponse(w, 404, err.Error())
		return
	}

	tokens, err := AuthService.CreateSession(r.Context(), user.Id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	refreshSeconds := 30 * 24 * 3600 // Default: 30 Tage

	if env := os.Getenv("RefreshExpirationSeconds"); env != "" {
		if v, err := strconv.Atoi(env); err == nil && v > 0 {
			refreshSeconds = v
		}
	}

	isProd := os.Getenv("AppEnv") == "prod"

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Duration(refreshSeconds) * time.Second),
	})

	type LoginResponse struct {
		AuthToken  string `json:"authToken"`
		Expiration int64  `json:"expiration"` // Unix Timestamp
	}

	response := LoginResponse{
		AuthToken:  tokens.AccessToken,
		Expiration: tokens.ExpiresAt.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")

	Http.JsonResponse(w, nil, "success", "routers.Login", response)
}
