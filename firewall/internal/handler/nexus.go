package handler

import (
	"firewall/internal/entity"
	"firewall/pkg/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addNexusAPI(rg *gin.RouterGroup) {
	rg.GET("/api/v2/config", handleConfig)
	rg.GET("/rest/integration/applications", handleApplications)
	rg.GET("/rest/config/proprietary", handleProprietary)
	rg.GET("/rest/product/features", handleFeatures)
	rg.GET("/rest/product/version", handleVersion)

	rep := rg.Group("/rest/integration/repositories")
	rep.GET("/evaluate/ignorePatterns", handleIgnorePatterns)
	rep.POST("/:nexus/configureRepositories", h.handleConfigure)
	rep.GET("/:nexus/getConfiguredRepositories", h.handleGetConfigure)

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

func handleApplications(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"applicationSummaries": []any{},
	})
}

func handleProprietary(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"packages": []any{},
		"regexes":  []any{},
	})
}

func handleFeatures(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, []any{})
}

func handleVersion(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"name":      "sonatype-clm-server",
		"version":   "1.174.0-01",
		"timestamp": "202403052139",
		"tag":       "17372cc51bfd53015634b712700829e39c5960f4",
		"build":     "build-number",
	})
}

func handleIgnorePatterns(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) handleConfigure(ctx *gin.Context) {
	instance := ctx.Param("nexus")
	var body struct {
		Repositories                    []entity.RepositoryDTO `json:"repositories" binding:"required"`
		RepositoryManagerProductName    string                 `json:"repositoryManagerProductName" binding:"required"`
		RepositoryManagerProductVersion string                 `json:"repositoryManagerProductVersion" binding:"required"`
	}
	ctx.BindJSON(&body)
	log.Printf("Configure: %+v\n", body)

	// ignore errors
	h.service.ConfigureRepositories(instance, &body.Repositories)
	ctx.Status(http.StatusOK)
}

func (h *Handler) handleGetConfigure(ctx *gin.Context) {
	instance := ctx.Param("nexus")
	sinceStr := ctx.Query("sinceUtcTimestamp")

	repos, err := h.service.GetConfiguredRepositories(instance, sinceStr)
	log.Printf("Sending repos: %+v\n", repos)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, repos)
}
