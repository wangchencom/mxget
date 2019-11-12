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

func GetAlbumFromXiaMi(c *gin.Context) {
	getAlbum(c, xiami.Client())
}

func GetAlbumFromBaiDu(c *gin.Context) {
	getAlbum(c, baidu.Client())
}

func getAlbum(c *gin.Context, client provider.API) {
	id := strings.TrimSpace(c.Param("id"))
	data, err := client.GetAlbum(id)
	response(c, client, data, err)
}
