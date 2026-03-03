package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//
// ===== Context Key =====
//

type contextKey struct{}

var UserIDKey = &contextKey{}

//
// ===== Compatibility Helpers =====
//

const Glob_date_layout = "2006-01-02T15:04:05"

func GetTime(sec int) time.Time {
	return time.Now().Add(time.Duration(sec) * time.Second)
}

//
// ===== JWT Config =====
//

const (
	tokenIssuer   = "gowrite-api"
	tokenAudience = "gowrite-client"
)

type Claims struct {
	UserID int `json:"uid"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

//
// ===== Init =====
//

func InitAuth() {

	secret := os.Getenv("JwtSecret")

	if secret == "" {
		log.Fatal("JwtSecret not set")
	}

	if len(secret) < 32 {
		log.Fatal("JwtSecret too short (min 32 bytes)")
	}

	jwtSecret = []byte(secret)
}

//
// ===== Token Creation (compatible) =====
//

func CreateAccessToken(userID int, expiresAt time.Time) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			Subject:   strconv.Itoa(userID),
			Audience:  []string{tokenAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

//
// ===== Token Validation =====
//

func validateToken(r *http.Request) (int, error) {

	header := r.Header.Get("Authorization")
	if header == "" {
		return 0, errors.New("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("invalid authorization format")
	}

	tokenStr := parts[1]

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {

			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return jwtSecret, nil
		},
		jwt.WithAudience(tokenAudience),
		jwt.WithIssuer(tokenIssuer),
	)

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	return claims.UserID, nil
}

//
// ===== Middleware =====
//

func authMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	userID, err := validateToken(r)
	if err != nil || userID <= 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), UserIDKey, userID)

	next(w, r.WithContext(ctx))
}

//
// ===== Router Helper (compatible) =====
//

func Handler(
	route string,
	router *mux.Router,
	next http.HandlerFunc,
	method string,
	permissions string,
) {

	n := negroni.New()
	n.Use(negroni.HandlerFunc(authMiddleware))

	if permissions == "" {
		n.UseHandlerFunc(next)
	} else {
		n.UseHandlerFunc(callWithPermissions(router, next, permissions))
	}

	router.Handle(route, n).Methods(method)
}

//
// ===== Permissions Stub =====
//

func callWithPermissions(
	router *mux.Router,
	next http.HandlerFunc,
	permissions string,
) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// Hier später echte Permission-Prüfung implementieren
		next(w, r)
	}
}

//
// ===== CORS =====
//

func CorsHandler(router *mux.Router) http.Handler {

	originsEnv := os.Getenv("AllowedOrigins")
	if originsEnv == "" {
		log.Fatal("AllowedOrigins not set")
	}

	allowed := strings.Split(originsEnv, "|")

	c := cors.New(cors.Options{
		AllowedOrigins:   allowed,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
	})

	return c.Handler(router)
}

//
// ===== Utils =====
//

func generateJTI() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
