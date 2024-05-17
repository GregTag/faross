package handler

import (
	"firewall/internal/service"
	"net/http"

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
	router.GET("/old", h.handleHome)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/f")
	})
	router.Static("/f", "./static")

	return router
}
