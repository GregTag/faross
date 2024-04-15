package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addPublicAPI(rg *gin.RouterGroup) {
	rg.GET("/status", h.handleGetAll)
	rg.POST("/status", h.handleGetStatus)
	rg.POST("/unquarantine", h.handleUnquarantine)
	rg.POST("/evaluate", h.handleEvaluate)
}

type request struct {
	Purl string `json:"purl" binding:"required"`
}

func (h *Handler) handleGetAll(ctx *gin.Context) {
	resp, err := h.service.GetAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) handleGetStatus(ctx *gin.Context) {
	var body request
	ctx.BindJSON(&body)

	pkg, err := h.service.GetPackage(body.Purl)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if pkg != nil {
		ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(pkg.Report))
	} else {
		ctx.AbortWithStatus(http.StatusNotFound)
	}
}

func (h *Handler) handleUnquarantine(ctx *gin.Context) {
	var body request
	ctx.BindJSON(&body)

	err := h.service.Unquarantine(body.Purl)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) handleEvaluate(ctx *gin.Context) {
	var body request
	ctx.BindJSON(&body)
	pkg, err := h.service.EvaluatePurl(body.Purl)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(pkg.Report))
}
