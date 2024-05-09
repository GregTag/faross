package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handleHome(ctx *gin.Context) {
	data, err := h.service.PrepareData()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.HTML(http.StatusOK, "table.html", data)
}
