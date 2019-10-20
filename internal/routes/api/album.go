package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/provider/kugou"
	"github.com/winterssy/mxget/pkg/provider/kuwo"
	"github.com/winterssy/mxget/pkg/provider/migu"
	"github.com/winterssy/mxget/pkg/provider/netease"
	"github.com/winterssy/mxget/pkg/provider/qq"
)

func GetAlbumFromNetEase(c *gin.Context) {
	getAlbum(c, netease.Client())
}

func GetAlbumFromQQ(c *gin.Context) {
	getAlbum(c, qq.Client())
}

func GetAlbumFromMiGu(c *gin.Context) {
	getAlbum(c, migu.Client())
}

func GetAlbumFromKuGou(c *gin.Context) {
	getAlbum(c, kugou.Client())
}

func GetAlbumFromKuWo(c *gin.Context) {
	getAlbum(c, kuwo.Client())
}

func getAlbum(c *gin.Context, client provider.API) {
	id := strings.TrimSpace(c.Param("id"))
	data, err := client.GetAlbum(id)
	if err != nil {
		c.JSON(500, &provider.Response{
			Code:     500,
			Msg:      err.Error(),
			Platform: client.Platform(),
		})
		return
	}

	c.JSON(200, &provider.Response{
		Code:     200,
		Msg:      "ok",
		Data:     data,
		Platform: client.Platform(),
	})
}
