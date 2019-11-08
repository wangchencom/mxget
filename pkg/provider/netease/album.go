package netease

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
	id, err := strconv.Atoi(albumId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetAlbumRaw(id)
	if err != nil {
		return nil, err
	}

	n := len(resp.Songs)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongURL(SongDefaultBR, resp.Songs...)
	a.patchSongLyric(resp.Songs...)
	songs := resolve(resp.Songs...)
	return &provider.Album{
		Name:   strings.TrimSpace(resp.Album.Name),
		PicURL: resp.Album.PicURL,
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetAlbumRaw(id int) (*AlbumResponse, error) {
	return std.GetAlbumRaw(id)
}

// 获取专辑
func (a *API) GetAlbumRaw(id int) (*AlbumResponse, error) {
	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodPost, fmt.Sprintf(APIGetAlbum, id),
		sreq.WithForm(weapi(struct{}{})),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get album: %s", resp.Msg)
	}

	return resp, nil
}
