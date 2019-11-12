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

func GetArtistFromNetEase(c *gin.Context) {
	getArtist(c, netease.Client())
}

func GetArtistFromQQ(c *gin.Context) {
	getArtist(c, qq.Client())
}

func GetArtistFromMiGu(c *gin.Context) {
	getArtist(c, migu.Client())
}

func GetArtistFromKuGou(c *gin.Context) {
	getArtist(c, kugou.Client())
}

func GetArtistFromKuWo(c *gin.Context) {
	getArtist(c, kuwo.Client())
}

func GetArtistFromXiaMi(c *gin.Context) {
	getArtist(c, xiami.Client())
}

func GetArtistFromBaiDu(c *gin.Context) {
	getArtist(c, baidu.Client())
}

func getArtist(c *gin.Context, client provider.API) {
	id := strings.TrimSpace(c.Param("id"))
	data, err := client.GetArtist(id)
	response(c, client, data, err)
}
