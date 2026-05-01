package http

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard response structure for all API endpoints.
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendSuccessResponse sends a standardized 2xx or 3xx response.
// If data is nil or an empty slice/map, the "data" field is omitted from the JSON.
func SendSuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	if isNilOrEmpty(data) {
		c.JSON(code, APIResponse{
			Code:    code,
			Message: message,
		})
		return
	}
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse sends a standardized 4xx or 5xx response.
func SendErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
	})
}

// isNilOrEmpty checks if the data is nil or an empty slice/map.
func isNilOrEmpty(data interface{}) bool {
	if data == nil {
		return true
	}

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return isNilOrEmpty(v.Elem().Interface())
	default:
		return false
	}
}
