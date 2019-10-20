package kuwo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetPlaylist(playlistId string) (*provider.Playlist, error) {
	return std.GetPlaylist(playlistId)
}

func (a *API) GetPlaylist(playlistId string) (*provider.Playlist, error) {
	resp, err := a.GetPlaylistRaw(playlistId, 1, 9999)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.MusicList)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	a.patchSongURL(SongDefaultBR, resp.Data.MusicList...)
	a.patchSongLyric(resp.Data.MusicList...)
	songs := a.resolve(resp.Data.MusicList)
	return &provider.Playlist{
		Name:   strings.TrimSpace(resp.Data.Name),
		PicURL: resp.Data.Img700,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetPlaylistRaw(playlistId string, page int, pageSize int) (*PlaylistResponse, error) {
	return std.GetPlaylistRaw(playlistId, page, pageSize)
}

// 获取歌单，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetPlaylistRaw(playlistId string, page int, pageSize int) (*PlaylistResponse, error) {
	params := sreq.Params{
		"pid": playlistId,
		"pn":  strconv.Itoa(page),
		"rn":  strconv.Itoa(pageSize),
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodGet, GetPlaylistAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get playlist: %s", resp.Msg)
	}

	return resp, nil
}
