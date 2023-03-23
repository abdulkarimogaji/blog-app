package db

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type GetCommentsFilters struct {
	UserId       *int
	BlogId       *int
	Message      *string
	PostedAfter  *string
	PostedBefore *string
}

type CreateCommentRequest struct {
	UserId   int    `json:"user_id" binding:"required,number"`
	BlogId   int    `json:"blog_id" binding:"required,number"`
	Message  string `json:"message" binding:"required"`
	PostedAt string `json:"posted_at" binding:"required,datetime=2006-01-02 15:04:05"`
}

func (d *DBStruct) GetComments(filters GetCommentsFilters, params PaginationParams) ([]Comment, int, error) {
	fields := []string{"comments.id", "comments.user_id", "blog_id", "message", "posted_at", "comments.created_at", "comments.updated_at", "users.id", "first_name", "last_name", "photo"}

	where := fmt.Sprintf("%s AND %s AND %s AND %s AND %s", getIntClause("comments.user_id", filters.UserId), getIntClause("blog_id", filters.BlogId), getLikeClause("message", filters.Message), getTimeClause("posted_at", ">", filters.PostedAfter), getTimeClause("posted_at", "<", filters.PostedBefore))

	sql := fmt.Sprintf("SELECT %s FROM comments LEFT JOIN users ON users.id = comments.user_id LEFT JOIN profile on profile.user_id = comments.user_id WHERE %s %s", strings.Join(fields, ","), where, getPaginationStr(params))
	log.Printf("\n\n %s \n\n", sql)
	rows, err := d.DB.Query(sql)
	if err != nil {
		return []Comment{}, 0, err
	}

	var total int
	row := d.DB.QueryRow(fmt.Sprintf("SELECT COUNT(id) as total FROM comments WHERE %s", where))
	err = row.Scan(&total)
	if err != nil {
		return []Comment{}, 0, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var tmp Comment
		err = rows.Scan(&tmp.Id, &tmp.UserId, &tmp.BlogId, &tmp.Message, &tmp.PostedAt, &tmp.CreatedAt, &tmp.UpdatedAt, &tmp.User.Id, &tmp.User.FirstName, &tmp.User.LastName, &tmp.User.Photo)
		if err != nil {
			return comments, 0, err
		}
		comments = append(comments, tmp)
	}

	return comments, total, nil
}

func (d *DBStruct) CreateComment(body CreateCommentRequest) (Comment, error) {
	// check if author exists
	var author User
	row := d.DB.QueryRow("SELECT id, first_name, last_name, email, created_at, updated_at  from users WHERE id = ?;", body.UserId)
	err := row.Scan(&author.Id, &author.FirstName, &author.LastName, &author.Email, &author.CreatedAt, &author.UpdatedAt)
	if err != nil {
		return Comment{}, err
	}

	// check if blog exists
	var blog Blog
	row = d.DB.QueryRow("SELECT id from blogs WHERE id = ?;", body.BlogId)
	err = row.Scan(&blog.Id)
	if err != nil {
		return Comment{}, err
	}

	stmt, err := d.DB.Prepare("INSERT INTO comments (user_id, blog_id, message, posted_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return Comment{}, err
	}

	createdAt := time.Now()
	updatedAt := time.Now()
	postedAt, err := time.Parse(time.DateTime, body.PostedAt)
	if err != nil {
		return Comment{}, err
	}

	result, err := stmt.Exec(body.UserId, body.BlogId, body.Message, postedAt, createdAt, updatedAt)

	if err != nil {
		return Comment{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Comment{}, err
	}

	return Comment{Id: int(id), UserId: body.UserId, BlogId: body.BlogId, Message: body.Message, PostedAt: postedAt, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}
