package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/provider/baidu"
	"github.com/winterssy/mxget/pkg/provider/kugou"
	"github.com/winterssy/mxget/pkg/provider/kuwo"
	"github.com/winterssy/mxget/pkg/provider/migu"
	"github.com/winterssy/mxget/pkg/provider/netease"
	"github.com/winterssy/mxget/pkg/provider/qq"
	"github.com/winterssy/mxget/pkg/provider/xiami"
)

func GetSongFromNetEase(c *gin.Context) {
	getSong(c, netease.Client())
}

func GetSongFromQQ(c *gin.Context) {
	getSong(c, qq.Client())
}

func GetSongFromMiGu(c *gin.Context) {
	getSong(c, migu.Client())
}

func GetSongFromKuGou(c *gin.Context) {
	getSong(c, kugou.Client())
}

func GetSongFromKuWo(c *gin.Context) {
	getSong(c, kuwo.Client())
}

func GetSongFromXiaMi(c *gin.Context) {
	getSong(c, xiami.Client())
}

func GetSongFromBaiDu(c *gin.Context) {
	getSong(c, baidu.Client())
}

func getSong(c *gin.Context, client provider.API) {
	id := strings.TrimSpace(c.Param("id"))
	data, err := client.GetSong(id)
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
		Data:     data,
		Platform: client.Platform(),
	})
}
