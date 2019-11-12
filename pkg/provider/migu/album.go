package migu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetAlbum(albumId string) (*provider.Album, error) {
	resp, err := a.GetAlbumRaw(albumId)
	if err != nil {
		return nil, err
	}

	if len(resp.Resource) == 0 || len(resp.Resource[0].SongItems) == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongLyric(resp.Resource[0].SongItems...)
	songs := resolve(resp.Resource[0].SongItems...)
	return &provider.Album{
		Id:     resp.Resource[0].AlbumId,
		Name:   strings.TrimSpace(resp.Resource[0].Title),
		PicURL: picURL(resp.Resource[0].ImgItems),
		Count:  len(songs),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(albumId string) (*AlbumResponse, error) {
	params := sreq.Params{
		"resourceId": albumId,
	}

	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbum,
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
