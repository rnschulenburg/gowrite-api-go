package UserLogin

import (
	"github.com/rnschulenburg/gowrite-api-go/App/Entities/UserEntity"
	"github.com/rnschulenburg/gowrite-api-go/App/Repositories/UserRepository"
)

func ByPassword(nickname, password string) (*UserEntity.User, error) {

	return UserRepository.ByLogin(nickname, password)
}
