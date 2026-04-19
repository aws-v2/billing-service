package http

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, docsHandler *DocsHandler) {

	docs := r.Group("/api/v1/billing/docs")
	{
		docs.GET("", docsHandler.GetPublicManifest)
		docs.GET("/:slug", docsHandler.GetPublicDoc)
	}

	internal := r.Group("/api/v1/billing/internal/docs")
	{
		internal.GET("", docsHandler.GetInternalManifest)
		internal.GET("/:slug", docsHandler.GetInternalDoc)
	}
}