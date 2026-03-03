package UserEntity

type User struct {
	Id       int    `json:"id"`
	NickName string `json:"nick_name"`
	Password string `json:"password"`
}
