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
func (h *DocsHandler) GetPublicManifest(c *gin.Context) {
	data, err := h.service.GetManifest(false)
	if err != nil {
		SendErrorResponse(c, 500, err.Error())
		return
	}
	SendSuccessResponse(c, 200, "Public manifest retrieved successfully", data)
}

func (h *DocsHandler) GetInternalManifest(c *gin.Context) {
	data, err := h.service.GetManifest(true)
	if err != nil {
		SendErrorResponse(c, 500, err.Error())
		return
	}
	SendSuccessResponse(c, 200, "Internal manifest retrieved successfully", data)
}


func (h *DocsHandler) GetPublicDoc(c *gin.Context) {
	slug := c.Param("slug")

	doc, err := h.service.GetDoc(slug, false)
	if err != nil {
		SendErrorResponse(c, 404, "not found")
		return
	}

	SendSuccessResponse(c, 200, "Public document retrieved successfully", doc)
}

func (h *DocsHandler) GetInternalDoc(c *gin.Context) {
	slug := c.Param("slug")

	doc, err := h.service.GetDoc(slug, true)
	if err != nil {
		SendErrorResponse(c, 404, "not found")
		return
	}

	SendSuccessResponse(c, 200, "Internal document retrieved successfully", doc)
}