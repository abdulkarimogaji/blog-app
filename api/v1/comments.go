package v1

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type getCommentQuery struct {
	Page         string `form:"page"`
	Limit        string `form:"limit"`
	Message      string `form:"message"`
	UserId       string `form:"user_id"`
	BlogId       string `form:"blog_id"`
	PostedAfter  string `form:"posted_after"`
	PostedBefore string `form:"posted_before"`
}

func getCommentsPaginate(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query getCommentQuery
		c.ShouldBindQuery(&query)
		filters, paginationParams := parseCommentQueryParams(query)

		comments, total, err := dbService.GetComments(filters, paginationParams)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
			})
			return
		}

		var totalPages, currentPage, pageSize *int

		if paginationParams.Limit == 0 {
			totalPages = nil
			currentPage = nil
			pageSize = nil
		} else {
			tmp := total / paginationParams.Limit
			tmp2 := 1
			totalPages = &tmp
			if paginationParams.Page > 0 {
				tmp2 = paginationParams.Page
			}
			currentPage = &tmp2
			pageSize = &paginationParams.Limit
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Blog created successfully",
			"data": gin.H{
				"list":        comments,
				"total":       total,
				"page":        currentPage,
				"page_size":   pageSize,
				"total_pages": totalPages,
			},
		})

	}
}

func parseCommentQueryParams(query getCommentQuery) (db.GetCommentsFilters, db.PaginationParams) {
	var limit, page int
	limit, _ = strconv.Atoi(query.Limit)
	page, _ = strconv.Atoi(query.Page)

	// get filters
	var message, postedBefore, postedAfter *string
	var userId, blogId *int

	u_id, _ := strconv.Atoi(query.UserId)
	b_id, _ := strconv.Atoi(query.BlogId)

	if u_id == 0 {
		userId = nil
	} else {
		userId = &u_id
	}

	if b_id == 0 {
		blogId = nil
	} else {
		blogId = &b_id
	}

	if query.Message == "" {
		message = nil
	} else {
		message = &query.Message
	}

	if _, err := time.Parse(time.DateTime, query.PostedAfter); err != nil || query.PostedAfter == "" {
		postedAfter = nil
	} else {
		postedAfter = &query.PostedAfter
	}

	if _, err := time.Parse(time.DateTime, query.PostedBefore); err != nil || query.PostedBefore == "" {
		postedBefore = nil
	} else {
		postedBefore = &query.PostedBefore
	}

	return db.GetCommentsFilters{Message: message,
			PostedAfter:  postedAfter,
			PostedBefore: postedBefore,
			UserId:       userId,
			BlogId:       blogId,
		}, db.PaginationParams{Page: page,
			Limit: limit}
}

func createComment(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body db.CreateCommentRequest
		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			validationResponse(err, c)
			return
		}

		comment, err := dbService.CreateComment(body)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "user does not exist",
					"error":   err.Error(),
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Blog created successfully",
			"data":    comment,
		})
	}
}
