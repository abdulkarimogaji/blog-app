package db

import (
	"time"

	"github.com/abdulkarimogaji/blognado/util"
)

type CreateBlogRequest struct {
	AuthorId  int    `json:"author_id" binding:"required,number"`
	Title     string `json:"title" binding:"required"`
	Excerpt   string `json:"excerpt" binding:"required"`
	Thumbnail string `json:"thumbnail" binding:"required,url"`
	Body      string `json:"body" binding:"required"`
	PostedAt  string `json:"posted_at" binding:"required,datetime=2006-01-02 15:04:05"`
}

func (d *DBStruct) CreateBlog(body CreateBlogRequest) (Blog, error) {
	// check if author exists
	var author User
	row := d.DB.QueryRow("SELECT id, first_name, last_name, email, created_at, updated_at  from users WHERE id = ?;", body.AuthorId)
	err := row.Scan(&author.Id, &author.FirstName, &author.LastName, &author.Email, &author.CreatedAt, &author.UpdatedAt)
	if err != nil {
		return Blog{}, err
	}

	stmt, err := d.DB.Prepare("INSERT INTO blogs (author_id, title, slug, excerpt, thumbnail, body, posted_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return Blog{}, err
	}

	createdAt := time.Now()
	updatedAt := time.Now()
	slug := util.ConvertToSlug(body.Title)
	postedAt, err := time.Parse(time.DateTime, body.PostedAt)
	if err != nil {
		return Blog{}, err
	}

	result, err := stmt.Exec(body.AuthorId, body.Title, slug, body.Excerpt, body.Thumbnail, body.Body, postedAt, createdAt, updatedAt)

	if err != nil {
		return Blog{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Blog{}, err
	}

	return Blog{Id: int(id), AuthorId: body.AuthorId, Title: body.Title, Slug: slug, Excerpt: body.Excerpt, Thumbnail: body.Thumbnail, Body: body.Body, PostedAt: postedAt, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}

func (db *DBStruct) GetBlogBySlug(slug string) (Blog, error) {
	return Blog{}, nil
}

func (db *DBStruct) GetBlogById(id int) (Blog, error) {
	return Blog{}, nil
}

func (db *DBStruct) GetBlogs() ([]Blog, error) {
	return []Blog{}, nil
}
