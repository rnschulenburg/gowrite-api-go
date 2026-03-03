package AuthRepository

import (
	"context"
	"errors"
	"time"

	"github.com/rnschulenburg/gowrite-api-go/Package/DbConnection"
)

//
// ===== Entity =====
//

type RefreshToken struct {
	ID        int64
	UserID    int
	TokenHash string
	Expires   time.Time
	Revoked   bool
	Created   time.Time
}

//
// ===== Create =====
//

func StoreRefreshToken(
	ctx context.Context,
	userID int,
	tokenHash string,
	expires time.Time,
) error {

	_, err := DbConnection.DB.Exec(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires, revoked)
		 VALUES ($1, $2, $3, false)`,
		userID,
		tokenHash,
		expires,
	)

	return err
}

//
// ===== Read by Token Hash =====
//

func GetRefreshToken(
	ctx context.Context,
	tokenHash string,
) (*RefreshToken, error) {

	var rt RefreshToken

	err := DbConnection.DB.QueryRow(
		ctx,
		`SELECT id, user_id, token, expires, revoked, created
		 FROM refresh_tokens
		 WHERE token = $1`,
		tokenHash,
	).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.Expires,
		&rt.Revoked,
		&rt.Created,
	)

	if err != nil {
		return nil, err
	}

	return &rt, nil
}

//
// ===== Validate Token =====
//

func ValidateRefreshToken(
	ctx context.Context,
	tokenHash string,
) (*RefreshToken, error) {

	rt, err := GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("refresh token not found")
	}

	if rt.Revoked {
		return nil, errors.New("refresh token revoked")
	}

	if time.Now().After(rt.Expires) {
		return nil, errors.New("refresh token expired")
	}

	return rt, nil
}

//
// ===== Revoke One Token =====
//

func RevokeRefreshToken(
	ctx context.Context,
	tokenHash string,
) error {

	_, err := DbConnection.DB.Exec(
		ctx,
		`UPDATE refresh_tokens
		 SET revoked = true
		 WHERE token = $1`,
		tokenHash,
	)

	return err
}

//
// ===== Revoke All Tokens for User (Logout All Devices) =====
//

func RevokeAllUserTokens(
	ctx context.Context,
	userID int,
) error {

	_, err := DbConnection.DB.Exec(
		ctx,
		`UPDATE refresh_tokens
		 SET revoked = true
		 WHERE user_id = $1`,
		userID,
	)

	return err
}

//
// ===== Delete Expired / Revoked Tokens =====
//

func CleanupExpiredTokens(ctx context.Context) error {

	_, err := DbConnection.DB.Exec(
		ctx,
		`DELETE FROM refresh_tokens
		 WHERE revoked = true
		    OR expires < NOW()`,
	)

	return err
}

//
// ===== Optional: Delete Old Tokens for User (Limit Sessions) =====
//

func DeleteUserTokensOlderThan(
	ctx context.Context,
	userID int,
	t time.Time,
) error {

	_, err := DbConnection.DB.Exec(
		ctx,
		`DELETE FROM refresh_tokens
		 WHERE user_id = $1
		   AND created < $2`,
		userID,
		t,
	)

	return err
}
