package v1

import (
	"github.com/abdulkarimogaji/blognado/api/middleware/auth"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/worker"
	"github.com/gin-gonic/gin"
)

const MYSQL_KEY_EXISTS = 1062

func ConfigureRoutes(router *gin.RouterGroup, db db.DBService, taskDistributor worker.TaskDistributor) {
	router.POST("/signup", signUp(db, taskDistributor))
	router.POST("/login", login(db, taskDistributor))
	router.POST("/refresh-token", refreshToken(db))
	router.POST("/verify-email", verifyEmail(db))
	authRouter := router.Group("/", auth.AuthorizeClient())
	authRouter.POST("/check", checkToken)
	authRouter.POST("/blogs", createBlog(db))
	authRouter.GET("/blogs/:idOrSlug", getBlogByIdOrSlug(db))
	authRouter.GET("/blogs", getBlogsPaginate(db))
	authRouter.GET("/comments", getCommentsPaginate(db))
	authRouter.POST("/comments", createComment(db))

	// TODO: setup a guard here to only allow role = admin to delete
	authRouter.DELETE("/:tableName/:id", deleteHandler(db))
}
