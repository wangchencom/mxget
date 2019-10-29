package migu

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

	if len(resp.Resource) == 0 || len(resp.Resource[0].SongItems) == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongURL(SongDefaultBR, resp.Resource[0].SongItems...)
	a.patchSongLyric(resp.Resource[0].SongItems...)
	songs := a.resolve(resp.Resource[0].SongItems...)
	return &provider.Album{
		Name:   strings.TrimSpace(resp.Resource[0].Title),
		PicURL: a.picURL(resp.Resource[0].ImgItems),
		Count:  len(songs),
		Songs:  songs,
	}, nil
}

func GetAlbumRaw(albumId string) (*AlbumResponse, error) {
	return std.GetAlbumRaw(albumId)
}

// 获取专辑
func (a *API) GetAlbumRaw(albumId string) (*AlbumResponse, error) {
	params := sreq.Params{
		"resourceId": albumId,
	}

	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodGet, GetAlbumAPI,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get album: %s", resp.Info)
	}

	return resp, nil
}
