package tencent

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
	if len(resp.Data.CDList) == 0 || len(resp.Data.CDList[0].SongList) == 0 {
		return nil, errors.New("get playlist: no data")
	}

	playlist := resp.Data.CDList[0]
	if playlist.PicURL == "" {
		playlist.PicURL = playlist.Logo
	}
	_songs := playlist.SongList
	a.patchSongsURLV1(ctx, _songs...)
	a.patchSongsLyric(ctx, _songs...)
	songs := resolve(_songs...)
	return &api.PlaylistResponse{
		Id:     playlist.DissTid,
		Name:   strings.TrimSpace(playlist.DissName),
		PicUrl: playlist.PicURL,
		Count:  uint32(len(songs)),
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId string) (*PlaylistResponse, error) {
	params := sreq.Params{
		"id": playlistId,
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodGet, APIGetPlaylist,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get playlist: %d", resp.Code)
	}

	return resp, nil
}
