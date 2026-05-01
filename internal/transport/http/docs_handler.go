package http

import (
	"github.com/Qarani-m/billing-service/internal/application"

	"github.com/gin-gonic/gin"
)

type DocsHandler struct {
	service *application.DocsService
}

func NewDocsHandler(service *application.DocsService) *DocsHandler {
	return &DocsHandler{service: service}
}

func (h *DocsHandler) GetDocsManifest(c *gin.Context) {
	role, _ := c.Get("role")
	isAdmin := role == "ADMIN"
	
	data, err := h.service.GetUnifiedManifest(isAdmin)
	if err != nil {
		SendErrorResponse(c, 500, err.Error())
		return
	}
	SendSuccessResponse(c, 200, "Manifest retrieved successfully", data)
}

func (h *DocsHandler) GetDoc(c *gin.Context) {
	slug := c.Param("slug")
	role, _ := c.Get("role")
	isAdmin := role == "ADMIN"

	doc, err := h.service.GetUnifiedDoc(slug, isAdmin)
	if err != nil {
		SendErrorResponse(c, 404, "Document not found or unauthorized")
		return
	}

	SendSuccessResponse(c, 200, "Document retrieved successfully", doc)
}