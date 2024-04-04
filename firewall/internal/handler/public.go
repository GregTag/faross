package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addPublicAPI(rg *gin.RouterGroup) {

	rg.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Pong!")
	})
}
