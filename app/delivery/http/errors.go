package http

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func newErrResponce(err error) errResponse {
	return errResponse{"error", err.Error()}
}

var (
	errClietnEmptyName  = errors.New("client name is empty")
	errInvalidAccountID = errors.New("invalid account ID")
)

func (h *Handler) errResponse(c *gin.Context, code int, err error) {
	c.JSON(code, newErrResponce(err))
}
