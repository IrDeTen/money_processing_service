package http

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	errClietnEmptyName  = errors.New("client name is empty")
	errInvalidAccountID = errors.New("invalid account ID")
)

func (h *Handler) errResponse(c *gin.Context, code int, err error) {
	c.JSON(code, map[string]interface{}{"status": "error", "error": err.Error()})
}
