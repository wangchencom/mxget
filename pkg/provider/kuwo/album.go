package kuwo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetAlbum(albumId string) (*provider.Album, error) {
	return std.GetAlbum(albumId)
}

func (a *API) GetAlbum(albumId string) (*provider.Album, error) {
	resp, err := a.GetAlbumRaw(albumId, 1, 9999)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.MusicList)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongURL(SongDefaultBR, resp.Data.MusicList...)
	a.patchSongLyric(resp.Data.MusicList...)
	songs := a.resolve(resp.Data.MusicList)
	return &provider.Album{
		Name:   strings.TrimSpace(resp.Data.Album),
		PicURL: resp.Data.Pic,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetAlbumRaw(albumId string, page int, pageSize int) (*AlbumResponse, error) {
	return std.GetAlbumRaw(albumId, page, pageSize)
}

// 获取专辑，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetAlbumRaw(albumId string, page int, pageSize int) (*AlbumResponse, error) {
	params := sreq.Params{
		"albumId": albumId,
		"pn":      strconv.Itoa(page),
		"rn":      strconv.Itoa(pageSize),
	}
	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodGet, GetAlbumAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get album: %s", resp.Msg)
	}

	return resp, nil
}
