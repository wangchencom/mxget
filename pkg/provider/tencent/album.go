package tencent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

func (a *API) GetAlbum(ctx context.Context, albumMid string) (*api.AlbumResponse, error) {
	resp, err := a.GetAlbumRaw(ctx, albumMid)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.GetSongInfo)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	_songs := resp.Data.GetSongInfo
	a.patchSongsURLV1(ctx, _songs...)
	a.patchSongsLyric(ctx, _songs...)
	songs := resolve(_songs...)
	return &api.AlbumResponse{
		Id:     resp.Data.GetAlbumInfo.FAlbumMid,
		Name:   strings.TrimSpace(resp.Data.GetAlbumInfo.FAlbumName),
		PicUrl: fmt.Sprintf(AlbumPicURL, resp.Data.GetAlbumInfo.FAlbumMid),
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumMid string) (*AlbumResponse, error) {
	params := sreq.Params{
		"albummid": albumMid,
	}

	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbum,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get album: %d", resp.Code)
	}

	return resp, nil
}
