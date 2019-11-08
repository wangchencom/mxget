package kugou

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/sreq"
)

func GetPlaylist(specialId string) (*provider.Playlist, error) {
	return std.GetPlaylist(specialId)
}

func (a *API) GetPlaylist(specialId string) (*provider.Playlist, error) {
	playlistInfo, err := a.GetPlaylistInfoRaw(specialId)
	if err != nil {
		return nil, err
	}

	playlistSong, err := a.GetPlaylistSongRaw(specialId, 1, -1)
	if err != nil {
		return nil, err
	}

	n := len(playlistSong.Data.Info)
	if n == 0 {
		return nil, errors.New("get playlist song: no data")
	}

	a.patchSongInfo(playlistSong.Data.Info...)
	a.patchAlbumInfo(playlistSong.Data.Info...)
	a.patchSongLyric(playlistSong.Data.Info...)
	songs := resolve(playlistSong.Data.Info...)
	return &provider.Playlist{
		Name:   strings.TrimSpace(playlistInfo.Data.SpecialName),
		PicURL: strings.ReplaceAll(playlistInfo.Data.ImgURL, "{size}", "480"),
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetPlaylistInfoRaw(specialId string) (*PlaylistInfoResponse, error) {
	return std.GetPlaylistInfoRaw(specialId)
}

// 获取歌单信息
func (a *API) GetPlaylistInfoRaw(specialId string) (*PlaylistInfoResponse, error) {
	params := sreq.Params{
		"specialid": specialId,
	}

	resp := new(PlaylistInfoResponse)
	err := a.Request(sreq.MethodGet, APIGetPlaylistInfo,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get playlist info: %s", resp.Error)
	}

	return resp, nil
}

func GetPlaylistSongRaw(specialId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	return std.GetPlaylistSongRaw(specialId, page, pageSize)
}

// 获取歌单歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetPlaylistSongRaw(specialId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	params := sreq.Params{
		"specialid": specialId,
		"page":      strconv.Itoa(page),
		"pagesize":  strconv.Itoa(pageSize),
	}

	resp := new(PlaylistSongsResponse)
	err := a.Request(sreq.MethodGet, APIGetPlaylistSongs,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get playlist song: %s", resp.Error)
	}

	return resp, nil
}
