package handler

import (
	"firewall/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addNexusAPI(rg *gin.RouterGroup) {
	rg.GET("/api/v2/config", handleConfig)
}

func handleConfig(ctx *gin.Context) {
	property := ctx.Query("property")
	response := gin.H{}

	if property == "quarantinedItemCustomMessage" {
		msg := config.Koanf.String("quarantineMessage")
		response["quarantinedItemCustomMessage"] = msg
	}

	ctx.JSON(http.StatusOK, response)
}
