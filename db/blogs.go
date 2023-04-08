package db

import (
	"fmt"
	"log"
	"strings"
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

type GetBlogsFilters struct {
	AuthorId     *int
	Title        *string
	PostedAfter  *string
	PostedBefore *string
	AuthorName   *string
}

type PaginationParams struct {
	Page  int
	Limit int
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

func (d *DBStruct) GetBlogBySlug(slug string) (Blog, error) {
	var blog Blog
	row := d.DB.QueryRow("SELECT blogs.*, users.id, first_name, last_name, email, photo, socials from blogs LEFT JOIN users ON users.id = author_id LEFT JOIN profile ON profile.user_id = author_id WHERE slug = ?;", slug)
	err := row.Scan(&blog.Id, &blog.AuthorId, &blog.Title, &blog.Slug, &blog.Excerpt, &blog.Thumbnail, &blog.Body, &blog.PostedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Author.Id, &blog.Author.FirstName, &blog.Author.LastName, &blog.Author.Email, &blog.Author.Photo, &blog.Author.Socials)
	if err != nil {
		return Blog{}, err
	}
	return blog, nil
}

func (d *DBStruct) GetBlogById(id int) (Blog, error) {
	var blog Blog
	row := d.DB.QueryRow("SELECT blogs.*, users.id, first_name, last_name, email, photo, socials from blogs LEFT JOIN users ON users.id = author_id LEFT JOIN profile ON profile.user_id = author_id WHERE blogs.id = ?;", id)
	err := row.Scan(&blog.Id, &blog.AuthorId, &blog.Title, &blog.Slug, &blog.Excerpt, &blog.Thumbnail, &blog.Body, &blog.PostedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Author.Id, &blog.Author.FirstName, &blog.Author.LastName, &blog.Author.Email, &blog.Author.Photo, &blog.Author.Socials)
	if err != nil {
		return Blog{}, err
	}
	return blog, nil
}

func (d *DBStruct) GetBlogs(filters GetBlogsFilters, params PaginationParams) ([]Blog, int, error) {
	fields := []string{"blogs.id", "author_id", "title", "slug", "excerpt", "thumbnail", "posted_at", "blogs.created_at", "blogs.updated_at", "users.id", "first_name", "last_name", "photo", "socials"}

	where := fmt.Sprintf("%s AND %s AND %s AND %s AND (%s OR %s)", getIntClause("author_id", filters.AuthorId), getLikeClause("title", filters.Title), getTimeClause("posted_at", ">", filters.PostedAfter), getTimeClause("posted_at", "<", filters.PostedBefore), getLikeClause("users.first_name", filters.AuthorName), getLikeClause("users.last_name", filters.AuthorName))

	sql := fmt.Sprintf("SELECT %s FROM blogs LEFT JOIN users ON users.id = author_id LEFT JOIN profile ON profile.user_id = author_id WHERE %s %s", strings.Join(fields, ","), where, getPaginationStr(params))
	log.Printf("\n\n %s \n\n", sql)
	rows, err := d.DB.Query(sql)
	if err != nil {
		return []Blog{}, 0, err
	}

	var total int
	row := d.DB.QueryRow(fmt.Sprintf("SELECT COUNT(blogs.id) as total FROM blogs LEFT JOIN users ON users.id = author_id WHERE %s", where))
	err = row.Scan(&total)
	if err != nil {
		return []Blog{}, 0, err
	}

	defer rows.Close()

	blogs := []Blog{}

	for rows.Next() {
		var tmp Blog
		err = rows.Scan(&tmp.Id, &tmp.AuthorId, &tmp.Title, &tmp.Slug, &tmp.Excerpt, &tmp.Thumbnail, &tmp.PostedAt, &tmp.CreatedAt, &tmp.UpdatedAt, &tmp.Author.Id, &tmp.Author.FirstName, &tmp.Author.LastName, &tmp.Author.Photo, &tmp.Author.Socials)
		if err != nil {
			return blogs, 0, err
		}
		blogs = append(blogs, tmp)
	}

	return blogs, total, nil
}
