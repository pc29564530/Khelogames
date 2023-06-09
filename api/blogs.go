package api

import (
	"database/sql"
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server HTTP request for blogs post.
type createBlogRequest struct {
	Username string `json:"username" `
	Title    string `json:"title"`
	Content  string `json:"content"`
}

func (server *Server) createBlog(ctx *gin.Context) {
	var req createBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBlogParams{
		Username: req.Username,
		Title:    req.Title,
		Content:  req.Content,
	}

	blog, err := server.store.CreateBlog(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

type getBlogRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBlog(ctx *gin.Context) {
	var req getBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blog, err := server.store.GetBlog(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, blog)
}

type getBlogUsernameRequest struct {
	Username string `json:"username"`
}

//func (server *Server) getBlogUsername(ctx *gin.Context) {
//	var req getBlogUsernameRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusNotFound, errorResponse(err))
//		return
//	}
//
//	blog, err := server.store.()
//}
