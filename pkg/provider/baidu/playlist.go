package baidu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func (a *API) GetPlaylist(playlistId string) (*provider.Playlist, error) {
	resp, err := a.GetPlaylistRaw(playlistId)
	if err != nil {
		return nil, err
	}

	n := len(resp.Result.SongList)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	a.patchSongURL(resp.Result.SongList...)
	a.patchSongLyric(resp.Result.SongList...)
	songs := resolve(resp.Result.SongList...)
	return &provider.Playlist{
		Name:   strings.TrimSpace(resp.Result.Info.ListTitle),
		PicURL: resp.Result.Info.ListPic,
		Count:  n,
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(playlistId string) (*PlaylistResponse, error) {
	params := sreq.Params{
		"list_id":   playlistId,
		"withcount": "1",
		"withsong":  "1",
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodGet, APIGetPlaylist,
		sreq.WithQuery(signPayload(params)),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get playlist: %s", resp.ErrorMessage)
	}

	return resp, nil
}
