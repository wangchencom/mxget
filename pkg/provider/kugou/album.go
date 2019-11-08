package kugou

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
	albumInfo, err := a.GetAlbumInfoRaw(albumId)
	if err != nil {
		return nil, err
	}

	albumSong, err := a.GetAlbumSongRaw(albumId, 1, -1)
	if err != nil {
		return nil, err
	}

	n := len(albumSong.Data.Info)
	if n == 0 {
		return nil, errors.New("get album song: no data")
	}

	a.patchSongInfo(albumSong.Data.Info...)
	a.patchAlbumInfo(albumSong.Data.Info...)
	a.patchSongLyric(albumSong.Data.Info...)
	songs := resolve(albumSong.Data.Info...)
	return &provider.Album{
		Name:   strings.TrimSpace(albumInfo.Data.AlbumName),
		PicURL: strings.ReplaceAll(albumInfo.Data.ImgURL, "{size}", "480"),
		Count:  n,
		Songs:  songs,
	}, nil
}

func GetAlbumInfoRaw(albumId string) (*AlbumInfoResponse, error) {
	return std.GetAlbumInfoRaw(albumId)
}

// 获取专辑信息
func (a *API) GetAlbumInfoRaw(albumId string) (*AlbumInfoResponse, error) {
	params := sreq.Params{
		"albumid": albumId,
	}

	resp := new(AlbumInfoResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbumInfo,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get album info: %s", resp.Error)
	}

	return resp, nil
}

func GetAlbumSongRaw(albumId string, page int, pageSize int) (*AlbumSongsResponse, error) {
	return std.GetAlbumSongRaw(albumId, page, pageSize)
}

// 获取专辑歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetAlbumSongRaw(albumId string, page int, pageSize int) (*AlbumSongsResponse, error) {
	params := sreq.Params{
		"albumid":  albumId,
		"page":     strconv.Itoa(page),
		"pagesize": strconv.Itoa(pageSize),
	}

	resp := new(AlbumSongsResponse)
	err := a.Request(sreq.MethodGet, APIGetAlbumSongs,
		sreq.WithQuery(params),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get album song: %s", resp.Error)
	}

	return resp, nil
}
