package migu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/sreq"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.PlaylistResponse, error) {
	resp, err := a.GetPlaylistRaw(ctx, playlistId)
	if err != nil {
		return nil, err
	}
	if len(resp.Resource) == 0 || len(resp.Resource[0].SongItems) == 0 {
		return nil, errors.New("get playlist: no data")
	}

	a.patchSongsLyric(ctx, resp.Resource[0].SongItems...)
	songs := resolve(resp.Resource[0].SongItems...)
	return &api.PlaylistResponse{
		Id:     resp.Resource[0].MusicListId,
		Name:   strings.TrimSpace(resp.Resource[0].Title),
		PicUrl: resp.Resource[0].ImgItem.Img,
		Count:  uint32(len(songs)),
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId string) (*PlaylistResponse, error) {
	params := sreq.Params{
		"resourceId": playlistId,
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodGet, APIGetPlaylist,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get playlist: %s", resp.errorMessage())
	}

	return resp, nil
}
