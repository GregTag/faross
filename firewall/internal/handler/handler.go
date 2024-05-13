package handler

import (
	"firewall/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return Handler{service}
}

func (h *Handler) GetRoute() *gin.Engine {
	var router = gin.Default()

	nexus := router.Group("/nexus")
	h.addNexusAPI(nexus)

	public := router.Group("/api")
	h.addPublicAPI(public)

	router.LoadHTMLFiles("templates/table.html")
	router.GET("/", h.handleHome)

	return router
}
