package http

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, docsHandler *DocsHandler, billingHandler *BillingHandler) {

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

	billing := r.Group("/api/v1/billing")
	{
		billing.POST("/card", billingHandler.AddCard)
		billing.POST("/mpesa/stk-push", billingHandler.InitiateStkPush)
		billing.GET("/billings", billingHandler.GetAllBillings)
		billing.GET("/services/:serviceName", billingHandler.GetBillingByService)
	}
}