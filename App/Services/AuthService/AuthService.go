package AuthService

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/rnschulenburg/gowrite-api-go/App/Repositories/AuthRepository"
	"github.com/rnschulenburg/gowrite-api-go/routers/auth"
)

//
// ===== Config =====
//

const (
	AccessTokenLifetime  = 1 * time.Minute
	RefreshTokenLifetime = 30 * 24 * time.Hour
)

//
// ===== Result Struct =====
//

type SessionTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

//
// ===== Utilities =====
//

func generateRefreshToken() (string, error) {

	b := make([]byte, 64)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

//
// ===== LOGIN =====
//

func CreateSession(
	ctx context.Context,
	userID int,
) (*SessionTokens, error) {

	// Access Token
	accessExp := time.Now().Add(AccessTokenLifetime)

	accessToken, err := auth.CreateAccessToken(userID, accessExp)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	hash := hashToken(refreshToken)

	refreshExp := time.Now().Add(RefreshTokenLifetime)

	err = AuthRepository.StoreRefreshToken(
		ctx,
		userID,
		hash,
		refreshExp,
	)
	if err != nil {
		return nil, err
	}

	return &SessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
	}, nil
}

//
// ===== REFRESH =====
//

func RefreshSession(
	ctx context.Context,
	refreshToken string,
) (*SessionTokens, error) {

	hash := hashToken(refreshToken)

	rt, err := AuthRepository.ValidateRefreshToken(ctx, hash)
	if err != nil {
		return nil, err
	}

	err = AuthRepository.RevokeRefreshToken(ctx, hash)
	if err != nil {
		return nil, err
	}

	return CreateSession(ctx, rt.UserID)
}

//
// ===== LOGOUT (single device) =====
//

func Logout(
	ctx context.Context,
	refreshToken string,
) error {

	hash := hashToken(refreshToken)

	return AuthRepository.RevokeRefreshToken(ctx, hash)
}

//
// ===== LOGOUT ALL DEVICES =====
//

func LogoutAll(
	ctx context.Context,
	userID int,
) error {

	return AuthRepository.RevokeAllUserTokens(ctx, userID)
}
