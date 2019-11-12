package api

import (
	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/pkg/provider"
)

type (
	Response struct {
		Code       int                 `json:"code"`
		Msg        string              `json:"msg,omitempty"`
		Data       interface{}         `json:"data,omitempty"`
		PlatformId provider.PlatformId `json:"platform_id,omitempty"`
	}
)

func response(c *gin.Context, client provider.API, data interface{}, err error) {
	if err != nil {
		c.JSON(500, &Response{
			Code:       500,
			Msg:        err.Error(),
			PlatformId: client.PlatformId(),
		})
		return
	}

	c.JSON(200, &Response{
		Code:       200,
		Data:       data,
		PlatformId: client.PlatformId(),
	})
}
