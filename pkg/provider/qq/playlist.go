package qq

import (
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetPlaylist(playlistId string) (*provider.Playlist, error) {
	return std.GetPlaylist(playlistId)
}

func (a *API) GetPlaylist(playlistId string) (*provider.Playlist, error) {
	resp, err := a.GetPlaylistRaw(playlistId)
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
	a.patchSongURL(_songs...)
	a.patchSongLyric(_songs...)
	songs := resolve(_songs...)
	return &provider.Playlist{
		Name:   strings.TrimSpace(playlist.DissName),
		PicURL: playlist.PicURL,
		Count:  len(songs),
		Songs:  songs,
	}, nil
}

func GetPlaylistRaw(playlistId string) (*PlaylistResponse, error) {
	return std.GetPlaylistRaw(playlistId)
}

// 获取歌单
func (a *API) GetPlaylistRaw(playlistId string) (*PlaylistResponse, error) {
	params := sreq.Params{
		"id": playlistId,
	}

	resp := new(PlaylistResponse)
	err := a.Request(sreq.MethodGet, APIGetPlaylist,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get playlist: %d", resp.Code)
	}

	return resp, nil
}
