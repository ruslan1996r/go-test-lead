package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BasicHandler struct{}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *BasicHandler) sendOk(ctx *gin.Context, val any) {
	ctx.JSON(http.StatusOK, val)
}

func (h *BasicHandler) sendInternalServerError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)

	ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
}

func (h *BasicHandler) notFound(ctx *gin.Context, val any) {
	ctx.JSON(http.StatusNotFound, val)
}
