package xiami

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetAlbum(albumId string) (*provider.Album, error) {
	resp, err := a.GetAlbumRaw(albumId)
	if err != nil {
		return nil, err
	}

	_songs := resp.Data.Data.AlbumDetail.Songs
	n := len(_songs)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongLyric(_songs...)
	songs := resolve(_songs...)
	return &provider.Album{
		Id:     resp.Data.Data.AlbumDetail.AlbumId,
		Name:   strings.TrimSpace(resp.Data.Data.AlbumDetail.AlbumName),
		PicURL: resp.Data.Data.AlbumDetail.AlbumLogo,
		Count:  n,
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(albumId string) (*AlbumResponse, error) {
	token, err := a.getToken(APIGetAlbum)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(albumId)
	if err != nil {
		model["albumStringId"] = albumId
	} else {
		model["albumId"] = albumId
	}
	params := sreq.Params(signPayload(token, model))
	resp := new(AlbumResponse)
	err = a.Request(sreq.MethodGet, APIGetAlbum, sreq.WithQuery(params)).JSON(resp)
	if err != nil {
		return nil, err
	}
	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get album: %w", err)
	}

	return resp, nil
}
