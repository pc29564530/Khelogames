package handlers

import (
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addAdminRequest struct {
	ContentID int64  `json:"content_id"`
	Admin     string `json:"admin"`
}

func (s *HandlersServer) AddAdminFunc(ctx *gin.Context) {
	var req addAdminRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.AddAdminParams{
		ContentID: req.ContentID,
		Admin:     req.Admin,
	}

	response, err := s.store.AddAdmin(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the admin ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type getAdminRequest struct {
	ContentID int64 `json:"content_id"`
}

func (s *HandlersServer) GetAdminFunc(ctx *gin.Context) {
	var req getAdminRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	response, err := s.store.GetAdmin(ctx, req.ContentID)
	if err != nil {
		s.logger.Error("Failed to get the admin: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type deletAdminRequest struct {
	ContentID int64 `json:"content_id"`
}

func (s *HandlersServer) DeleteAdminFunc(ctx *gin.Context) {
	var req getAdminRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationHeaderKey).(*token.Payload)

	arg := db.DeleteAdminParams{
		ContentID: req.ContentID,
		Admin:     authPayload.Username,
	}

	response, err := s.store.DeleteAdmin(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to delete the admin: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
