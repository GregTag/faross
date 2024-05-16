package handler

import (
	"errors"
	"firewall/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addPublicAPI(rg *gin.RouterGroup) {
	rg.GET("/status", h.handleGetAll)
	rg.POST("/report", h.handleGetReport)
	rg.POST("/unquarantine", h.handleUnquarantine)
	rg.POST("/evaluate", h.handleEvaluate)
	rg.PUT("/comment", h.handlePutComment)
}

type request struct {
	Purl string `json:"purl" binding:"required"`
}

func (h *Handler) handleGetAll(ctx *gin.Context) {
	resp, err := h.service.GetAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) handleGetReport(ctx *gin.Context) {
	var body request
	if ctx.BindJSON(&body) != nil {
		return
	}

	pkg, err := h.service.GetPackage(body.Purl)
	if errors.Is(err, entity.ErrPackageNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(pkg.Report))
	}
}

type requestComment struct {
	request
	Comment string `json:"comment" binding:"required"`
}

func (h *Handler) handleUnquarantine(ctx *gin.Context) {
	var body requestComment
	if ctx.BindJSON(&body) != nil {
		return
	}

	err := h.service.Unquarantine(body.Purl, body.Comment)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) handleEvaluate(ctx *gin.Context) {
	var body request
	if ctx.BindJSON(&body) != nil {
		return
	}
	pkg, err := h.service.EvaluatePurl(body.Purl)
	if err != nil {
		if errors.Is(err, entity.ErrPending) {
			ctx.AbortWithStatus(http.StatusAccepted)
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(pkg.Report))
}

func (h *Handler) handlePutComment(ctx *gin.Context) {
	var body requestComment
	if ctx.BindJSON(&body) != nil {
		return
	}

	err := h.service.ChangeComment(body.Purl, body.Comment)
	if errors.Is(err, entity.ErrPackageNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.Status(http.StatusOK)
	}
}
