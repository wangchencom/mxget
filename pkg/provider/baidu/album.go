package baidu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetAlbum(albumId string) (*provider.Album, error) {
	return std.GetAlbum(albumId)
}

func (a *API) GetAlbum(albumId string) (*provider.Album, error) {
	resp, err := a.GetAlbumRaw(albumId)
	if err != nil {
		return nil, err
	}

	n := len(resp.SongList)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongURL(resp.SongList...)
	a.patchSongLyric(resp.SongList...)
	songs := resolve(resp.SongList...)
	return &provider.Album{
		Name:   strings.TrimSpace(resp.AlbumInfo.Title),
		PicURL: resp.AlbumInfo.PicBig,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetAlbumRaw(albumId string) (*AlbumResponse, error) {
	return std.GetAlbumRaw(albumId)
}

// 获取专辑
func (a *API) GetAlbumRaw(albumId string) (*AlbumResponse, error) {
	params := sreq.Params{
		"album_id": albumId,
	}

	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbum,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 0 && resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get album: %s", resp.ErrorMessage)
	}

	return resp, nil
}
