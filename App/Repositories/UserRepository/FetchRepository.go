package UserRepository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rnschulenburg/gowrite-api-go/App/Entities/UserEntity"
	"github.com/rnschulenburg/gowrite-api-go/Package/DbConnection"
)

func FetchProjects(authUserId, userId int) ([]UserEntity.UserProject, error) {

	query := `SELECT
			CONCAT(author.first_name, ' ', author.last_name) AS author,
			author.nick_name  AS author_nick_name,
			p.name            AS project_name,
			p.title           AS title,
			p.created         AS project_start,
			up.permissions
		FROM projects p

		JOIN users_projects up
			ON p.id = up.project_id
		   AND up.user_id = $1

		JOIN users_projects up_author
			ON p.id = up_author.project_id
		   AND up_author.user_id = $2
		   AND up_author.permissions = 'author'

		JOIN users author
			ON author.id = up_author.user_id
		
		ORDER BY up.modified DESC

		LIMIT 50
	`

	rows, err := DbConnection.DB.Query(
		context.Background(),
		query,
		authUserId,
		userId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[UserEntity.UserProject],
	)

	if err != nil {
		return nil, err
	}

	return projects, nil
}
