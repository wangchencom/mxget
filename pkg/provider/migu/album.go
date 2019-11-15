package migu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

func (a *API) GetAlbum(ctx context.Context, albumId string) (*api.AlbumResponse, error) {
	resp, err := a.GetAlbumRaw(ctx, albumId)
	if err != nil {
		return nil, err
	}

	if len(resp.Resource) == 0 || len(resp.Resource[0].SongItems) == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsLyric(ctx, resp.Resource[0].SongItems...)
	songs := resolve(resp.Resource[0].SongItems...)
	return &api.AlbumResponse{
		Id:     resp.Resource[0].AlbumId,
		Name:   strings.TrimSpace(resp.Resource[0].Title),
		PicUrl: picURL(resp.Resource[0].ImgItems),
		Count:  uint32(len(songs)),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	params := sreq.Params{
		"resourceId": albumId,
	}

	resp := new(AlbumResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbum,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get album: %s", resp.Info)
	}

	return resp, nil
}
