package UserRepository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/rnschulenburg/gowrite-api-go/App/Entities/UserEntity"
	"github.com/rnschulenburg/gowrite-api-go/Package/DbConnection"
)

var ErrUserNotFound = errors.New("user not found")

func ByLogin(nickname, password string) (*UserEntity.User, error) {

	query := `
		SELECT id, nick_name, password
		FROM users
		WHERE nick_name = $1 AND password = $2
	`

	row := DbConnection.DB.QueryRow(
		context.Background(),
		query,
		nickname,
		password,
	)

	var u UserEntity.User

	err := row.Scan(&u.Id, &u.NickName, &u.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}
