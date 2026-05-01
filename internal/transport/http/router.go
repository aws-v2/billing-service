package http

import (
	"github.com/Qarani-m/billing-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, docsHandler *DocsHandler, billingHandler *BillingHandler) {
	r.Use(middleware.AuthMiddleware())

	docs := r.Group("/api/v1/billing/docs")
	{
		docs.GET("", docsHandler.GetDocsManifest)
		docs.GET("/:slug", docsHandler.GetDoc)
	}

	billing := r.Group("/api/v1/billing")
	{
		billing.POST("/card", billingHandler.AddCard)
		billing.POST("/mpesa/stk-push", billingHandler.InitiateStkPush)
		billing.GET("/billings", billingHandler.GetAllBillings)
		billing.GET("/services/:serviceName", billingHandler.GetBillingByService)
	}
}