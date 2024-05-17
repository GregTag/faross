package handler

import (
	"firewall/internal/entity"
	"firewall/internal/service"
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
	rep.POST("/:nexus_id/configureRepositories", h.handleConfigure)
	rep.POST("/:nexus_id/:public_id/enable/:enabled", h.handleAuditEnable)
	rep.POST("/:nexus_id/:public_id/quarantine/:enabled", h.handleQuarantineEnable)
	rep.GET("/:nexus_id/getConfiguredRepositories", h.handleGetConfigure)
	rep.POST("/:nexus_id/:public_id/evaluate/quarantine", h.handleEvalQuarantine)
	rep.GET("/:nexus_id/:public_id/components/unquarantined", h.handleUnquarantined)
	rep.GET("/:nexus_id/:public_id/summary", h.handleGetSummary)
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
	instance := ctx.Param("nexus_id")
	var body struct {
		Repositories                    []entity.RepositoryDTO `json:"repositories" binding:"required"`
		RepositoryManagerProductName    string                 `json:"repositoryManagerProductName" binding:"required"`
		RepositoryManagerProductVersion string                 `json:"repositoryManagerProductVersion" binding:"required"`
	}
	if ctx.BindJSON(&body) != nil {
		return
	}
	log.Printf("Configure: %+v\n", body)

	// ignore errors
	h.service.ConfigureRepositories(instance, body.Repositories)
	ctx.Status(http.StatusOK)
}

func (h *Handler) handleGetConfigure(ctx *gin.Context) {
	instance := ctx.Param("nexus_id")
	sinceStr := ctx.Query("sinceUtcTimestamp")

	repos, err := h.service.GetConfiguredRepositories(instance, sinceStr)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, repos)
}

func (h *Handler) handleAuditEnable(ctx *gin.Context) {
	instance := ctx.Param("nexus_id")
	name := ctx.Param("public_id")
	enabledStr := ctx.Param("enabled")
	enabled := enabledStr == "true"

	repos, err := h.service.SetAuditEnable(instance, name, enabled)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, repos)
}

func (h *Handler) handleQuarantineEnable(ctx *gin.Context) {
	instance := ctx.Param("nexus_id")
	name := ctx.Param("public_id")
	enabledStr := ctx.Param("enabled")
	enabled := enabledStr == "true"

	repos, err := h.service.SetQuarantineEnable(instance, name, enabled)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, repos)
}

func (h *Handler) handleEvalQuarantine(ctx *gin.Context) {
	instance := ctx.Param("nexus_id")
	name := ctx.Param("public_id")
	var body struct {
		Components []service.EvalDataRequest `json:"components" binding:"required"`
		Cause      string                    `json:"cause" binding:"required"`
	}
	if ctx.BindJSON(&body) != nil {
		return
	}

	evalResults, err := h.service.EvalRequest(instance, name, body.Components)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := gin.H{
		"componentEvalResults": evalResults,
	}
	log.Printf("Eval request:\n%+v\nEval response:\n%+v\n", body, response)
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) handleUnquarantined(ctx *gin.Context) {
	instance := ctx.Param("nexus_id")
	name := ctx.Param("public_id")
	sinceStr := ctx.Query("sinceUtcTimestamp")

	list, err := h.service.GetUnquarantined(instance, name, sinceStr)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"pathnames": list,
	})
}

func (h *Handler) handleGetSummary(ctx *gin.Context) {
	instance := ctx.Param("nexus_id")
	name := ctx.Param("public_id")
	response, err := h.service.GetSummary(instance, name)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
