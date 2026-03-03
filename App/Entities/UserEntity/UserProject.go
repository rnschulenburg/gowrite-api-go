package UserEntity

import "time"

type UserProject struct {
	Author         string    `db:"author" json:"author"`
	AuthorNickName string    `db:"author_nick_name" json:"authorNickName"`
	ProjectName    string    `db:"project_name" json:"projectName"`
	Title          string    `db:"title" json:"title"`
	ProjectStart   time.Time `db:"project_start" json:"projectStart"`
	Permissions    string    `db:"permissions" json:"permissions"`
}
