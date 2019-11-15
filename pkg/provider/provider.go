package provider

import (
	"context"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
)

var (
	std *sreq.Client
)

type (
	API interface {
		// 搜索歌曲
		SearchSongs(ctx context.Context, keyword string) (*api.SearchSongsResponse, error)
		// 获取单曲
		GetSong(ctx context.Context, songId string) (*api.SongResponse, error)
		// 获取歌手
		GetArtist(ctx context.Context, artistId string) (*api.ArtistResponse, error)
		// 获取专辑
		GetAlbum(ctx context.Context, albumId string) (*api.AlbumResponse, error)
		// 获取歌单
		GetPlaylist(ctx context.Context, playlistId string) (*api.PlaylistResponse, error)
		// 网络请求
		Request(method string, url string, opts ...sreq.RequestOption) *sreq.Response
	}
)

func init() {
	std = sreq.New(nil)
	std.SetDefaultRequestOpts(
		sreq.WithHeaders(sreq.Headers{
			"User-Agent": UserAgent,
		}),
	)
}

func Client() *sreq.Client {
	return std
}
