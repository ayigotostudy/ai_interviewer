package controller

import (
	"ai_jianli_go/types/resp/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller definition
type Controller[T any] struct {
	Request  *T
	Response *common.Response
	c        *gin.Context
}

// NewCtrl Generic factory function for creating a controller
func NewCtrl[T any](c *gin.Context) *Controller[T] {
	return &Controller[T]{
		Request:  new(T),
		Response: new(common.Response),
		c:        c,
	}
}

// NoDataJSON parse with Nodata to json and return
func (ctrl *Controller[T]) NoDataJSON(code int64) {
	ctrl.Response.SetNoData(code)
	ctrl.c.JSON(http.StatusOK, ctrl.Response)
}

// WithDataJSON parse with data to json and return
func (ctrl *Controller[T]) WithDataJSON(code int64, data any) {
	ctrl.Response.SetWithData(code, data)
	ctrl.c.JSON(http.StatusOK, ctrl.Response)
}
