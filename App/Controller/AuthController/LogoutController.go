package AuthController

import (
	"github.com/rnschulenburg/gowrite-api-go/App/Services/AuthService"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return
	}

	_ = AuthService.Logout(r.Context(), cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}
