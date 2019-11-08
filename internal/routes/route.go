package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/internal/routes/api"
)

func Init(router *gin.Engine) {
	r := router.Group("/api")

	// 网易云音乐
	r.GET("/netease/search/:keyword", api.SearchSongsFromNetEase)
	r.GET("/netease/song/:id", api.GetSongFromNetEase)
	r.GET("/netease/artist/:id", api.GetArtistFromNetEase)
	r.GET("/netease/album/:id", api.GetAlbumFromNetEase)
	r.GET("/netease/playlist/:id", api.GetPlaylistFromNetEase)

	// QQ音乐
	r.GET("/qq/search/:keyword", api.SearchSongsFromQQ)
	r.GET("/qq/song/:id", api.GetSongFromQQ)
	r.GET("/qq/artist/:id", api.GetArtistFromQQ)
	r.GET("/qq/album/:id", api.GetAlbumFromQQ)
	r.GET("/qq/playlist/:id", api.GetPlaylistFromQQ)

	// 咪咕音乐
	r.GET("/migu/search/:keyword", api.SearchSongsFromMiGu)
	r.GET("/migu/song/:id", api.GetSongFromMiGu)
	r.GET("/migu/artist/:id", api.GetArtistFromMiGu)
	r.GET("/migu/album/:id", api.GetAlbumFromMiGu)
	r.GET("/migu/playlist/:id", api.GetPlaylistFromMiGu)

	// 酷狗音乐
	r.GET("/kugou/search/:keyword", api.SearchSongsFromKuGou)
	r.GET("/kugou/song/:id", api.GetSongFromKuGou)
	r.GET("/kugou/artist/:id", api.GetArtistFromKuGou)
	r.GET("/kugou/album/:id", api.GetAlbumFromKuGou)
	r.GET("/kugou/playlist/:id", api.GetPlaylistFromKuGou)

	// 酷我音乐
	r.GET("/kuwo/search/:keyword", api.SearchSongsFromKuWo)
	r.GET("/kuwo/song/:id", api.GetSongFromKuWo)
	r.GET("/kuwo/artist/:id", api.GetArtistFromKuWo)
	r.GET("/kuwo/album/:id", api.GetAlbumFromKuWo)
	r.GET("/kuwo/playlist/:id", api.GetPlaylistFromKuWo)
}
